package sqarol

import (
	"strings"
	"testing"
)

func TestNoInvalidVariations(t *testing.T) {
	domains := []string{"www.brrflow.com", "billcotecpa.com", "wawand.co", "www.crc-cpa.com"}

	for _, d := range domains {
		vars, err := Generate(d)
		if err != nil {
			t.Fatalf("Error for %s: %v", d, err)
		}

		for _, v := range vars {
			// Check for double dots
			if strings.Contains(v.Variant, "..") {
				t.Errorf("Double dots in variant: %s (kind: %s) from %s", v.Variant, v.Kind, d)
			}

			// Check for labels ending with hyphen
			if strings.Contains(v.Variant, "-.") {
				t.Errorf("Label ending with hyphen in variant: %s (kind: %s) from %s", v.Variant, v.Kind, d)
			}

			// Check for labels starting with hyphen (after a dot)
			if strings.Contains(v.Variant, ".-") {
				t.Errorf("Label starting with hyphen in variant: %s (kind: %s) from %s", v.Variant, v.Kind, d)
			}

			// Check for www- prefix
			if strings.HasPrefix(v.Variant, "www-.") || strings.Contains(v.Variant, ".www-.") {
				t.Errorf("Invalid www- prefix in variant: %s (kind: %s) from %s", v.Variant, v.Kind, d)
			}
		}
	}
}
