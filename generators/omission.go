package generators

import "strings"

// Omission generates variations by removing one character at a time
// from the name part of the domain.
func omission(domain string) (string, []string) {
	added := make(map[string]bool)
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "omission", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for i := 0; i < len(name); i++ {
		fuzzy := name[:i] + name[i+1:] + tld
		if added[fuzzy] {
			continue
		}

		result = append(result, fuzzy)
		added[fuzzy] = true
	}

	return "omission", result
}
