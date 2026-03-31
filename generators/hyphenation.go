package generators

import "strings"

// Hyphenation generates variations by inserting a hyphen between every
// pair of adjacent characters in the name part of the domain.
func hyphenation(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "hyphenation", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for i := 1; i < len(name); i++ {
		// Skip if hyphen would be at start of a label (after a dot)
		if i > 0 && name[i-1] == '.' {
			continue
		}
		// Skip if hyphen would be at end of a label (before a dot)
		if i < len(name) && name[i] == '.' {
			continue
		}
		result = append(result, name[:i]+"-"+name[i:]+tld)
	}

	return "hyphenation", result
}
