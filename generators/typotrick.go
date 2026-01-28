package generators

import "strings"

// TypoTrick generates variations using visually similar character
// substitutions (leet-speak, digraphs, etc.), applied both in
// isolation and in pairwise combinations.
func typoTrick(domain string) (string, []string) {
	var result []string

	typos := [][]string{
		{"o", "0"},
		{"l", "1"},
		{"e", "3"},
		{"a", "4"},
		{"s", "5"},
		{"b", "6"},
		{"t", "7"},
		{"b", "8"},
		{"g", "9"},
		{"ks", "x"},
		{"s", "z"},
		{"d", "cl"},
		{"w", "vv"},
		{"u", "ii"},
		{"t", "ll"},
		{"m", "rn"},
		{"y", "i"},
		{"g", "cj"},
		{"d", "cl"},
		{"a", "ci"},
		{"fi", "a"},
	}

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "typo-trick", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	// Apply one by one in isolation
	for _, comb := range typos {
		base := comb[0]
		trick := comb[1]

		if !strings.Contains(name, base) {
			continue
		}

		result = append(result, strings.Replace(name, base, trick, -1)+tld)
	}

	// Apply in combination (2 together)
	for _, t := range typos {
		base := t[0]
		trick := t[1]

		if !strings.Contains(name, base) {
			continue
		}

		for _, t2 := range typos {
			base2 := t2[0]
			trick2 := t2[1]

			if !strings.Contains(name, base2) {
				continue
			}

			nx := strings.Replace(name, base, trick, -1)
			nx = strings.Replace(nx, base2, trick2, -1)

			result = append(result, nx+tld)
		}
	}

	return "typo-trick", result
}
