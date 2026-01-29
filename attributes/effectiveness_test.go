package attributes

import (
	"math"
	"testing"
)

func TestEffectiveness(t *testing.T) {
	cases := []struct {
		name     string
		original string
		variant  string
		distance int
		minScore float64
		maxScore float64
	}{
		{
			name:     "rn for m is highly effective",
			original: "modern.com",
			variant:  "rnodern.com",
			distance: Levenshtein("modern.com", "rnodern.com"),
			minScore: 0.6,
			maxScore: 1.0,
		},
		{
			name:     "vv for w is highly effective",
			original: "web.com",
			variant:  "vveb.com",
			distance: Levenshtein("web.com", "vveb.com"),
			minScore: 0.6,
			maxScore: 1.0,
		},
		{
			name:     "single char omission is moderately effective",
			original: "example.com",
			variant:  "exmple.com",
			distance: 1,
			minScore: 0.5,
			maxScore: 1.0,
		},
		{
			name:     "digit substitution less effective than letter-only",
			original: "google.com",
			variant:  "g00gle.com",
			distance: Levenshtein("google.com", "g00gle.com"),
			minScore: 0.3,
			maxScore: 0.85,
		},
		{
			name:     "long prefix addition is less effective",
			original: "example.com",
			variant:  "login-example.com",
			distance: Levenshtein("example.com", "login-example.com"),
			minScore: 0.0,
			maxScore: 0.45,
		},
		{
			name:     "tld swap has moderate effectiveness",
			original: "example.com",
			variant:  "example.net",
			distance: Levenshtein("example.com", "example.net"),
			minScore: 0.1,
			maxScore: 0.6,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			score := Effectiveness(tc.original, tc.variant, tc.distance)

			if score < tc.minScore || score > tc.maxScore {
				t.Errorf("Effectiveness(%q, %q, %d) = %.4f, expected in [%.2f, %.2f]",
					tc.original, tc.variant, tc.distance, score, tc.minScore, tc.maxScore)
			}
		})
	}
}

func Test_editDistanceScore(t *testing.T) {
	cases := []struct {
		distance int
		expected float64
	}{
		{0, 0.0},
		{1, 1.0},
		{2, 0.5},
		{3, 1.0 / 3.0},
		{10, 0.1},
	}

	for _, tc := range cases {
		got := editDistanceScore(tc.distance)
		if math.Abs(got-tc.expected) > 0.001 {
			t.Errorf("editDistanceScore(%d) = %f, expected %f", tc.distance, got, tc.expected)
		}
	}
}

func Test_lengthRatioScore(t *testing.T) {
	cases := []struct {
		a, b     string
		expected float64
	}{
		{"abc", "abc", 1.0},
		{"abc", "abcd", 0.75},
		{"", "abc", 0.0},
		{"abcdefghij", "abcde", 0.5},
	}

	for _, tc := range cases {
		got := lengthRatioScore(tc.a, tc.b)
		if math.Abs(got-tc.expected) > 0.001 {
			t.Errorf("lengthRatioScore(%q, %q) = %f, expected %f", tc.a, tc.b, got, tc.expected)
		}
	}
}

func Test_charClassScore(t *testing.T) {
	cases := []struct {
		name     string
		original string
		variant  string
		expected float64
	}{
		{"no new classes", "example.com", "exampel.com", 1.0},
		{"digit introduced", "example.com", "3xample.com", 0.5},
		{"original has digits", "ex4mple.com", "3x4mple.com", 1.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := charClassScore(tc.original, tc.variant)
			if math.Abs(got-tc.expected) > 0.001 {
				t.Errorf("charClassScore(%q, %q) = %f, expected %f",
					tc.original, tc.variant, got, tc.expected)
			}
		})
	}
}

func Test_visualSimilarityScore(t *testing.T) {
	cases := []struct {
		name     string
		original string
		variant  string
		minScore float64
		maxScore float64
	}{
		{
			name:     "rn for m exact match",
			original: "modern.com",
			variant:  "rnodern.com",
			minScore: 0.9,
			maxScore: 1.0,
		},
		{
			name:     "no confusable pair",
			original: "example.com",
			variant:  "exampla.com",
			minScore: 0.2,
			maxScore: 0.4,
		},
		{
			name:     "identical strings",
			original: "example.com",
			variant:  "example.com",
			minScore: 0.0,
			maxScore: 0.0,
		},
		{
			name:     "completely different length",
			original: "example.com",
			variant:  "login-example.com",
			minScore: 0.0,
			maxScore: 0.1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := visualSimilarityScore(tc.original, tc.variant)
			if got < tc.minScore || got > tc.maxScore {
				t.Errorf("visualSimilarityScore(%q, %q) = %.4f, expected in [%.2f, %.2f]",
					tc.original, tc.variant, got, tc.minScore, tc.maxScore)
			}
		})
	}
}
