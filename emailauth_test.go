package sqarol

import "testing"

func TestHasSPFRecord(t *testing.T) {
	cases := []struct {
		name     string
		records  []string
		expected bool
	}{
		{"no records", nil, false},
		{"valid SPF record", []string{"v=spf1 include:_spf.google.com ~all"}, true},
		{"valid SPF record, mixed case", []string{"V=SPF1 include:_spf.google.com ~all"}, true},
		{"SPF record among other TXT records", []string{"some-verification-token=abc123", "v=spf1 -all"}, true},
		{"no SPF record", []string{"some-verification-token=abc123"}, false},
		{"DMARC record is not an SPF record", []string{"v=DMARC1; p=reject;"}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasSPFRecord(tc.records); got != tc.expected {
				t.Errorf("hasSPFRecord(%v) = %v, want %v", tc.records, got, tc.expected)
			}
		})
	}
}

func TestHasDMARCRecord(t *testing.T) {
	cases := []struct {
		name     string
		records  []string
		expected bool
	}{
		{"no records", nil, false},
		{"valid DMARC record", []string{"v=DMARC1; p=reject; rua=mailto:dmarc@example.com"}, true},
		{"valid DMARC record, lowercase", []string{"v=dmarc1; p=none;"}, true},
		{"no DMARC record", []string{"some-verification-token=abc123"}, false},
		{"SPF record is not a DMARC record", []string{"v=spf1 -all"}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasDMARCRecord(tc.records); got != tc.expected {
				t.Errorf("hasDMARCRecord(%v) = %v, want %v", tc.records, got, tc.expected)
			}
		})
	}
}
