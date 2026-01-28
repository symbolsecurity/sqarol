package generators

import (
	"fmt"
	"strings"
)

// WithSuffix generates variations by appending common phishing-related
// suffixes to the domain name.
func withSuffix(domain string) (string, []string) {
	suffixes := []string{
		"app", "platform", "info", "us", "site",
		"online", "secure", "support", "help", "account", "login", "verify",
		"service", "update", "management", "auth", "verification", "banking",
		"securelogin", "pay", "checkout", "payment", "portal", "billing",
		"email", "customer", "contact", "notification", "alerts", "safety",
		"access", "user", "signin", "password", "web", "net", "protection",
		"recover", "accountsecure", "reset",
	}

	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "suffix", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for _, suffix := range suffixes {
		result = append(result, fmt.Sprintf("%s-%s%s", name, suffix, tld))
	}

	return "suffix", result
}
