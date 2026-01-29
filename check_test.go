package sqarol_test

import (
	"context"
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
