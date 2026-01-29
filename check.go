package sqarol

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// whoisSem limits concurrent WHOIS TCP connections to avoid
// rate-limiting from WHOIS servers.
var whoisSem = make(chan struct{}, 5)

// parkingNameservers contains hostname suffixes of known domain parking
// nameserver providers. A domain whose NS records match any of these
// suffixes is likely parked.
var parkingNameservers = []string{
	"sedoparking.com",
	"bodis.com",
	"parkingcrew.net",
	"above.com",
	"pendingrenewaldeletion.com",
	"parklogic.com",
	"parkitonline.com",
	"domainparking.com",
	"hugedomains.com",
	"afternic.com",
	"undeveloped.com",
	"dan.com",
	"uniregistry.com",
	"domaincontrol.com",
	"registrar-servers.com",
}

// parkingIPPrefixes contains IP address prefixes associated with known
// domain parking services. A domain whose A record matches any of these
// prefixes is likely parked.
var parkingIPPrefixes = []string{
	// Sedo
	"52.119.124.",
	// GoDaddy parking
	"34.102.136.",
	"184.168.131.",
	// Bodis
	"199.59.242.",
	"199.59.243.",
	// ParkingCrew / Freenom
	"104.219.248.",
	"104.219.249.",
	// Sedoparking
	"91.195.240.",
	"91.195.241.",
	// Above.com
	"66.96.149.",
	// HugeDomains
	"65.55.72.",
	// Team Internet / ParkLogic
	"185.53.178.",
	"185.53.179.",
}

// DomainCheck holds the results of checking a domain's registration
// and DNS status.
type DomainCheck struct {
	// Domain is the queried domain name.
	Domain string `json:"domain"`
	// IsRegistered indicates whether the domain is registered.
	// Determined by the presence of NS records for the domain.
	IsRegistered bool `json:"is_registered"`
	// Owner is the registrant organization or name extracted from WHOIS,
	// if available.
	Owner string `json:"owner,omitempty"`
	// HasARecords indicates whether the domain has IPv4 address records.
	HasARecords bool `json:"has_a_records"`
	// HasMXRecords indicates whether the domain has mail exchange records.
	HasMXRecords bool `json:"has_mx_records"`
	// IsParked indicates whether the domain appears to be parked,
	// based on known parking nameservers and IP addresses.
	IsParked bool `json:"is_parked"`
}

// Check queries DNS and WHOIS to determine whether a domain is
// registered, who owns it, and what A and MX records it has.
//
// Registration is determined by the presence of NS records, which is
// the most reliable signal — a domain has NS records if and only if it
// has been delegated by the registry. WHOIS is queried separately to
// extract the registrant owner, with automatic referral following for
// richer results.
//
// It uses the provided context for cancellation and timeouts.
func Check(ctx context.Context, domain string) (*DomainCheck, error) {
	domain, err := normalize(domain)
	if err != nil {
		return nil, err
	}

	check := &DomainCheck{Domain: domain}
	resolver := net.DefaultResolver

	// Look up NS records to determine registration.
	var nsHosts []string
	nss, err := resolver.LookupNS(ctx, domain)
	if err == nil && len(nss) > 0 {
		check.IsRegistered = true
		for _, ns := range nss {
			nsHosts = append(nsHosts, strings.TrimSuffix(ns.Host, "."))
		}
	}

	// Look up A records.
	var ipv4s []string
	ips, err := resolver.LookupHost(ctx, domain)
	if err == nil {
		for _, ip := range ips {
			if parsed := net.ParseIP(ip); parsed != nil && parsed.To4() != nil {
				check.HasARecords = true
				ipv4s = append(ipv4s, ip)
			}
		}
	}

	// Look up MX records.
	mxs, err := resolver.LookupMX(ctx, domain)
	if err == nil && len(mxs) > 0 {
		check.HasMXRecords = true
	}

	// Detect parking from NS hostnames and A record IPs.
	check.IsParked = isParkingNS(nsHosts) || isParkingIP(ipv4s)

	// Query WHOIS to extract owner information.
	check.Owner = whoisOwner(ctx, domain)

	return check, nil
}

// whoisOwner queries WHOIS for a domain and returns the registrant
// owner, following referrals to the registrar's WHOIS server for
// richer results. Returns an empty string if the owner cannot be
// determined.
func whoisOwner(ctx context.Context, domain string) string {
	server := whoisServer(domain)
	resp, err := queryWhois(ctx, server, domain)
	if err != nil || resp == "" {
		return ""
	}

	// Check for a referral to the registrar's WHOIS server.
	if referral := extractReferralServer(resp); referral != "" && referral != server {
		referralResp, err := queryWhois(ctx, referral, domain)
		if err == nil && referralResp != "" {
			// Prefer the referral response for owner extraction
			// since it typically has richer registrant data.
			if owner := extractWhoisOwner(referralResp); owner != "" {
				return owner
			}
		}
	}

	return extractWhoisOwner(resp)
}

// extractReferralServer parses a WHOIS response for a referral to
// the registrar's own WHOIS server. Returns an empty string if no
// referral is found.
func extractReferralServer(response string) string {
	fields := []string{
		"Registrar WHOIS Server:",
		"ReferralServer:",
		"Whois Server:",
	}

	for line := range strings.SplitSeq(response, "\n") {
		trimmed := strings.TrimSpace(line)
		for _, field := range fields {
			if strings.HasPrefix(trimmed, field) {
				value := strings.TrimSpace(trimmed[len(field):])
				// Strip any protocol prefix (some responses include "whois://").
				value = strings.TrimPrefix(value, "whois://")
				value = strings.TrimPrefix(value, "http://")
				value = strings.TrimPrefix(value, "https://")
				// Strip any trailing path or port.
				if idx := strings.IndexAny(value, ":/"); idx != -1 {
					value = value[:idx]
				}
				if value != "" {
					return value
				}
			}
		}
	}

	return ""
}

// whoisServer returns the appropriate WHOIS server for the given domain's TLD.
func whoisServer(domain string) string {
	tld := domain
	if idx := strings.LastIndex(domain, "."); idx != -1 {
		tld = domain[idx+1:]
	}

	servers := map[string]string{
		"com":  "whois.verisign-grs.com",
		"net":  "whois.verisign-grs.com",
		"org":  "whois.pir.org",
		"info": "whois.afilias.net",
		"io":   "whois.nic.io",
		"co":   "whois.nic.co",
		"biz":  "whois.biz",
		"us":   "whois.nic.us",
		"me":   "whois.nic.me",
		"uk":   "whois.nic.uk",
		"de":   "whois.denic.de",
		"fr":   "whois.nic.fr",
		"eu":   "whois.eu",
		"ru":   "whois.tcinet.ru",
		"au":   "whois.auda.org.au",
		"ca":   "whois.cira.ca",
		"br":   "whois.registro.br",
		"in":   "whois.registry.in",
		"nl":   "whois.sidn.nl",
		"be":   "whois.dns.be",
		"at":   "whois.nic.at",
		"ch":   "whois.nic.ch",
		"it":   "whois.nic.it",
		"se":   "whois.iis.se",
		"no":   "whois.norid.no",
		"dk":   "whois.dk-hostmaster.dk",
		"fi":   "whois.fi",
		"pl":   "whois.dns.pl",
		"cz":   "whois.nic.cz",
		"jp":   "whois.jprs.jp",
		"kr":   "whois.kr",
		"cn":   "whois.cnnic.cn",
		"tw":   "whois.twnic.net.tw",
		"xyz":  "whois.nic.xyz",
		"app":  "whois.nic.google",
		"dev":  "whois.nic.google",
	}

	if server, ok := servers[strings.ToLower(tld)]; ok {
		return server
	}

	// Fallback: try whois.nic.<tld> which works for many newer TLDs.
	return "whois.nic." + strings.ToLower(tld)
}

// queryWhois performs a raw TCP WHOIS query against the given server.
// It acquires the WHOIS semaphore to limit concurrent connections.
func queryWhois(ctx context.Context, server, domain string) (string, error) {
	// Acquire semaphore to throttle concurrent WHOIS connections.
	select {
	case whoisSem <- struct{}{}:
		defer func() { <-whoisSem }()
	case <-ctx.Done():
		return "", ctx.Err()
	}

	addr := server + ":43"

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(10 * time.Second)
	}

	dialer := net.Dialer{Deadline: deadline}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return "", fmt.Errorf("whois dial %s: %w", addr, err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(deadline)

	_, err = fmt.Fprintf(conn, "%s\r\n", domain)
	if err != nil {
		return "", fmt.Errorf("whois write: %w", err)
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return sb.String(), fmt.Errorf("whois read: %w", err)
	}

	return sb.String(), nil
}

// isParkingNS reports whether any of the given nameserver hostnames
// match a known parking provider suffix.
func isParkingNS(nsHosts []string) bool {
	for _, host := range nsHosts {
		lower := strings.ToLower(host)
		for _, suffix := range parkingNameservers {
			if strings.HasSuffix(lower, suffix) {
				return true
			}
		}
	}
	return false
}

// isParkingIP reports whether any of the given IPv4 addresses match
// a known parking provider IP prefix.
func isParkingIP(ips []string) bool {
	for _, ip := range ips {
		for _, prefix := range parkingIPPrefixes {
			if strings.HasPrefix(ip, prefix) {
				return true
			}
		}
	}
	return false
}

// extractWhoisOwner attempts to extract the registrant organization or
// name from a WHOIS response. Returns an empty string if not found.
func extractWhoisOwner(response string) string {
	// Fields to look for, in priority order.
	fields := []string{
		"Registrant Organization:",
		"Registrant Name:",
		"registrant:",
		"org-name:",
		"Organisation:",
		"Organization:",
		"Registrant:",
		"holder:",
	}

	for line := range strings.SplitSeq(response, "\n") {
		trimmed := strings.TrimSpace(line)
		for _, field := range fields {
			if strings.HasPrefix(strings.ToLower(trimmed), strings.ToLower(field)) {
				value := strings.TrimSpace(trimmed[len(field):])
				if value != "" && !strings.HasPrefix(value, "REDACTED") {
					return value
				}
			}
		}
	}

	return ""
}
