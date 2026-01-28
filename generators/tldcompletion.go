package generators

import (
	"fmt"
	"strings"
)

// CompleteWithTLD checks if the name part ends with a known TLD string
// and, if so, splits it to create a new domain. For example,
// "examplenet.com" becomes "example.net".
func tldCompletion(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "complete-with-tld", result
	}

	name := domain[:dix]

	for _, sub := range tldCatalog {
		if !strings.HasSuffix(name, sub) {
			continue
		}

		index := strings.LastIndex(name, sub)
		newDomainName := name[:index]
		tld := name[index:]
		fuzzyDomain := fmt.Sprintf("%s.%s", newDomainName, tld)

		result = append(result, fuzzyDomain)
	}

	return "complete-with-tld", result
}
