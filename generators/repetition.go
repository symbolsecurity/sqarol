package generators

import "strings"

// Repetition generates variations by doubling each character one at a
// time in the name part of the domain.
func repetition(domain string) (string, []string) {
	var result []string
	seen := make(map[string]bool)

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "repetition", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for i := 0; i < len(name); i++ {
		fuzzy := name[:i] + string(name[i]) + name[i:] + tld
		if seen[fuzzy] || strings.Contains(fuzzy, "..") {
			continue
		}

		seen[fuzzy] = true
		result = append(result, fuzzy)
	}

	return "repetition", result
}
