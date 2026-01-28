package attributes

import "testing"

func TestLevenshtein(t *testing.T) {
	cases := []struct {
		name     string
		a, b     string
		expected int
	}{
		{"identical", "example.com", "example.com", 0},
		{"one deletion", "example.com", "exmple.com", 1},
		{"one insertion", "example.com", "examplee.com", 1},
		{"one substitution", "example.com", "exampla.com", 1},
		{"transposition counts as 2", "ab.com", "ba.com", 2},
		{"empty first", "", "abc", 3},
		{"empty second", "abc", "", 3},
		{"both empty", "", "", 0},
		{"tld swap", "example.com", "example.net", 3},
		{"prefix added", "example.com", "login-example.com", 6},
		{"hyphen inserted", "example.com", "exam-ple.com", 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Levenshtein(tc.a, tc.b)
			if got != tc.expected {
				t.Fatalf("Levenshtein(%q, %q) = %d, expected %d", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}
