package generators

import "strings"

// Transposition generates variations by swapping adjacent characters
// in the name part of the domain.
func transposition(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "transposition", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for i := 0; i < len(name)-1; i++ {
		trans := name[:i] + string(name[i+1]) + string(name[i]) + name[i+2:] + tld
		if trans == domain {
			continue
		}

		result = append(result, trans)
	}

	return "transposition", result
}
