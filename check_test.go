package sqarol_test

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/symbolsecurity/sqarol"
)

func TestCheck_InvalidDomain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cases := []struct {
		name   string
		domain string
	}{
		{"empty string", ""},
		{"no tld", "symbolsecurity"},
		{"too long", string(make([]byte, 254)) + ".com"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sqarol.Check(ctx, tc.domain)
			if err == nil {
				t.Fatal("expected error but got nil")
			}
		})
	}
}

func TestCheck_KnownDomain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	check, err := sqarol.Check(ctx, "google.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if check.Domain != "google.com" {
		t.Errorf("expected Domain %q, got %q", "google.com", check.Domain)
	}

	if !check.IsRegistered {
		t.Error("expected google.com to be registered (NS records should exist)")
	}

	if !check.HasARecords {
		t.Error("expected google.com to have A records")
	}

	if !check.HasMXRecords {
		t.Error("expected google.com to have MX records")
	}

	// Verify ARecords slice is populated
	if len(check.ARecords) == 0 {
		t.Error("expected google.com to have non-empty ARecords slice")
	}

	// Verify each A record is a valid IPv4 address
	for i, ip := range check.ARecords {
		parsed := net.ParseIP(ip)
		if parsed == nil {
			t.Errorf("ARecords[%d] %q is not a valid IP address", i, ip)
		} else if parsed.To4() == nil {
			t.Errorf("ARecords[%d] %q is not a valid IPv4 address", i, ip)
		}
	}

	// Verify MXRecords slice is populated
	if len(check.MXRecords) == 0 {
		t.Error("expected google.com to have non-empty MXRecords slice")
	}

	// Verify each MX record has valid properties
	for i, mx := range check.MXRecords {
		if mx.Host == "" {
			t.Errorf("MXRecords[%d].Host is empty", i)
		}
		if mx.Host[len(mx.Host)-1] == '.' {
			t.Errorf("MXRecords[%d].Host %q ends with a dot (should be trimmed)", i, mx.Host)
		}
		if mx.Pref == 0 && len(check.MXRecords) > 1 {
			// Only flag if there are multiple records and this one is zero
			// (single record with Pref=0 is technically valid)
			t.Logf("MXRecords[%d].Pref is 0 (may be valid for single record)", i)
		}
	}

	// Verify MXRecords are sorted by Pref ascending
	for i := 0; i < len(check.MXRecords)-1; i++ {
		if check.MXRecords[i].Pref > check.MXRecords[i+1].Pref {
			t.Errorf("MXRecords not sorted by Pref: MXRecords[%d].Pref=%d > MXRecords[%d].Pref=%d",
				i, check.MXRecords[i].Pref, i+1, check.MXRecords[i+1].Pref)
		}
	}

	// Verify at least one MX record has resolved IPs
	hasResolvedMX := false
	for _, mx := range check.MXRecords {
		if len(mx.IPs) > 0 {
			hasResolvedMX = true
			break
		}
	}
	if !hasResolvedMX {
		t.Error("expected at least one MX record to have resolved IPs")
	}
}

func TestCheck_URL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	check, err := sqarol.Check(ctx, "https://google.com/search")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if check.Domain != "google.com" {
		t.Errorf("expected Domain %q, got %q", "google.com", check.Domain)
	}
}

func TestCheck_UnregisteredDomain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a domain that is very unlikely to be registered.
	check, err := sqarol.Check(ctx, "xyzthisdomaindoesnotexist123456.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if check.IsRegistered {
		t.Error("expected unregistered domain to show IsRegistered=false")
	}

	if check.HasARecords {
		t.Error("expected no A records for unregistered domain")
	}

	if len(check.ARecords) != 0 {
		t.Errorf("expected empty ARecords slice for unregistered domain, got %d records", len(check.ARecords))
	}

	if len(check.MXRecords) != 0 {
		t.Errorf("expected empty MXRecords slice for unregistered domain, got %d records", len(check.MXRecords))
	}
}

func TestCheck_KnownDomainNotParked(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// google.com is a real active domain, not parked.
	check, err := sqarol.Check(ctx, "google.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if check.IsParked {
		t.Error("expected google.com to not be parked")
	}
}

func TestCheck_UnregisteredDomainNotParked(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	check, err := sqarol.Check(ctx, "xyzthisdomaindoesnotexist123456.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if check.IsParked {
		t.Error("expected unregistered domain to not be parked")
	}
}

func TestCheck_RegisteredWithoutWebPresence(t *testing.T) {
	// Verifies that NS-based registration detection works for domains
	// that are registered but may not have A or MX records (e.g. parked).
	// symbolsecurity.com is known to be registered and should have NS records.
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	check, err := sqarol.Check(ctx, "symbolsecurity.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !check.IsRegistered {
		t.Error("expected symbolsecurity.com to be registered via NS lookup")
	}
}

func TestCheck_JSONSerialization(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with a known domain that has A and MX records
	check, err := sqarol.Check(ctx, "google.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(check)
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	// Unmarshal into map to verify key presence
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify a_records key exists and is non-empty array
	aRecordsVal, hasARecords := result["a_records"]
	if !hasARecords {
		t.Error("expected 'a_records' key in JSON output")
	} else {
		aRecordsArr, ok := aRecordsVal.([]interface{})
		if !ok {
			t.Error("expected 'a_records' to be an array")
		} else if len(aRecordsArr) == 0 {
			t.Error("expected 'a_records' to be non-empty for google.com")
		}
	}

	// Verify mx_records key exists and is non-empty array
	mxRecordsVal, hasMXRecords := result["mx_records"]
	if !hasMXRecords {
		t.Error("expected 'mx_records' key in JSON output")
	} else {
		mxRecordsArr, ok := mxRecordsVal.([]interface{})
		if !ok {
			t.Error("expected 'mx_records' to be an array")
		} else if len(mxRecordsArr) == 0 {
			t.Error("expected 'mx_records' to be non-empty for google.com")
		}
	}

	// Verify backward compatibility: has_a_records key exists
	if _, hasKey := result["has_a_records"]; !hasKey {
		t.Error("expected 'has_a_records' key in JSON output (backward compatibility)")
	}

	// Verify backward compatibility: has_mx_records key exists
	if _, hasKey := result["has_mx_records"]; !hasKey {
		t.Error("expected 'has_mx_records' key in JSON output (backward compatibility)")
	}

	// Test with unregistered domain (empty slices should be omitted)
	check2, err := sqarol.Check(ctx, "xyzthisdomaindoesnotexist123456.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	jsonData2, err := json.Marshal(check2)
	if err != nil {
		t.Fatalf("failed to marshal unregistered domain to JSON: %v", err)
	}

	var result2 map[string]interface{}
	err = json.Unmarshal(jsonData2, &result2)
	if err != nil {
		t.Fatalf("failed to unmarshal unregistered domain JSON: %v", err)
	}

	// Verify a_records key is absent (omitempty on empty slice)
	if _, hasARecords := result2["a_records"]; hasARecords {
		t.Error("expected 'a_records' key to be absent for unregistered domain (omitempty)")
	}

	// Verify mx_records key is absent (omitempty on empty slice)
	if _, hasMXRecords := result2["mx_records"]; hasMXRecords {
		t.Error("expected 'mx_records' key to be absent for unregistered domain (omitempty)")
	}
}
