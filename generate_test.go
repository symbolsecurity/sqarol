package sqarol_test

import (
	"github.com/symbolsecurity/sqarol"
	"testing"
)

func TestGenerate(t *testing.T) {
	domain := "symbolsecurity.com"

	vars, err := sqarol.Generate(domain)
	if err != nil {
		t.Fatal("there was an error generating variations:", err)
	}

	if len(vars) == 0 {
		t.Fatal("expected variations but got none")
	}

	// Every variation must have Original, Variant, Kind, and Distance set.
	for i, v := range vars {
		if v.Original != domain {
			t.Errorf("variation[%d]: expected Original %q, got %q", i, domain, v.Original)
		}

		if v.Variant == "" {
			t.Errorf("variation[%d]: Variant is empty", i)
		}

		if v.Kind == "" {
			t.Errorf("variation[%d]: Kind is empty", i)
		}

		if v.Distance <= 0 {
			t.Errorf("variation[%d] (%s): expected Distance > 0, got %d", i, v.Variant, v.Distance)
		}

		if v.Effectiveness < 0 || v.Effectiveness > 1 {
			t.Errorf("variation[%d] (%s): expected Effectiveness in [0, 1], got %f", i, v.Variant, v.Effectiveness)
		}
	}
}

func TestGenerate_InvalidDomain(t *testing.T) {
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
			_, err := sqarol.Generate(tc.domain)
			if err == nil {
				t.Fatal("expected error but got nil")
			}
		})
	}
}

func TestGenerate_URL(t *testing.T) {
	vars, err := sqarol.Generate("https://symbolsecurity.com/path")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(vars) == 0 {
		t.Fatal("expected variations but got none")
	}

	// Original should be the normalized domain, not the full URL.
	for _, v := range vars {
		if v.Original != "symbolsecurity.com" {
			t.Fatalf("expected Original %q, got %q", "symbolsecurity.com", v.Original)
		}
	}
}

func TestGenerate_KindsCoverage(t *testing.T) {
	// Use "examplenet.com" so the complete-with-tld algorithm fires
	// ("examplenet" ends with "net", a known TLD).
	vars, err := sqarol.Generate("examplenet.com")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	expectedKinds := []string{
		"omission",
		"transposition",
		"vowel-swap",
		"hyphenation",
		"repetition",
		"extra-random-letter",
		"prefix",
		"suffix",
		"subdomain",
		"tld-swap",
		"typo-trick",
		"phonetic",
		"inflect",
		"complete-with-tld",
	}

	kindSet := make(map[string]bool)
	for _, v := range vars {
		kindSet[v.Kind] = true
	}

	for _, kind := range expectedKinds {
		if !kindSet[kind] {
			t.Errorf("expected kind %q to be present in variations", kind)
		}
	}
}
