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
		result = append(result, name[:i]+"-"+name[i:]+tld)
	}

	return "hyphenation", result
}
