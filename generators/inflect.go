package generators

import (
	"strings"

	"github.com/jinzhu/inflection"
)

// Inflect generates plural and singular forms of the domain name.
func inflect(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "inflect", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	domainPlural := inflection.Plural(name) + tld
	domainSingular := inflection.Singular(name) + tld

	if domainPlural != domain {
		result = append(result, domainPlural)
	}

	if domainSingular != domain {
		result = append(result, domainSingular)
	}

	return "inflect", result
}
