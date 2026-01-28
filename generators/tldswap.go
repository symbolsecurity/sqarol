package generators

import "strings"

// TLDSwap generates variations by replacing the TLD with every other
// TLD from the catalog.
func tldSwap(domain string) (string, []string) {
	var result []string
	parts := strings.Split(domain, ".")

	if len(parts) < 2 {
		tld := parts[len(parts)-2]
		for _, i := range tldCatalog {
			if i == tld {
				firstPart := parts[:len(parts)-2]
				finalPart := parts[len(parts)-1]
				for i := range tldCatalog {
					secondPart := tldCatalog[i]
					if secondPart == parts[len(parts)-2] {
						continue
					}
					result = append(result, strings.Join(firstPart, ".")+"."+secondPart+"."+finalPart)
				}
				return "tld-swap", result
			}
		}
	}

	firstPart := parts[:len(parts)-1]
	for i := range tldCatalog {
		secondPart := tldCatalog[i]
		if secondPart != parts[len(parts)-1] {
			result = append(result, strings.Join(firstPart, ".")+"."+secondPart)
		}
	}

	return "tld-swap", result
}
