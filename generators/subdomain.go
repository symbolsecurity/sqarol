package generators

import "strings"

// Subdomain generates variations by inserting a dot at every position
// in the name part to simulate subdomain creation.
func subdomain(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "subdomain", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for i := 1; i < len(name)-1; i++ {
		sub := name[:i] + "." + name[i:] + tld
		if strings.Contains(sub, "..") {
			continue
		}

		result = append(result, sub)
	}

	return "subdomain", result
}
