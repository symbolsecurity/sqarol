// Package sqarol generates domain typosquatting variations for a given
// domain name. It produces look-alike domains annotated with the
// generation technique, Levenshtein edit distance, and a visual
// deceptiveness score.
package sqarol

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/symbolsecurity/sqarol/attributes"
	"github.com/symbolsecurity/sqarol/generators"
)

// Generate normalizes the given domain and runs every registered
// fuzzing algorithm on it, returning all produced variations or an
// error if the domain is not valid.
func Generate(domain string) ([]Variation, error) {
	domain, err := normalize(domain)
	if err != nil {
		return nil, err
	}

	var result []Variation
	for _, g := range generators.All() {
		kind, variants := g(domain)
		for _, variant := range variants {
			result = append(result, Variation{
				Original: domain,
				Variant:  variant,
				Kind:     kind,
			})
		}
	}

	for i := range result {
		result[i].Distance = attributes.Levenshtein(result[i].Original, result[i].Variant)
		result[i].Effectiveness = attributes.Effectiveness(result[i].Original, result[i].Variant, result[i].Distance)
	}

	return result, nil
}

// normalize validates if the domain passed as parameter
// matches the Fully Qualified Domain Name(FQDN) convention.
// It returns the hostname part if domain describes a full
// URI resource or same domain if it is a valid domain name,
// otherwise it returns an error.
func normalize(domain string) (string, error) {
	fqdn := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

	ndomain := domain

	if strings.Contains(domain, "://") {
		u, err := url.Parse(domain)
		if err != nil {
			return "", errors.New("malformed url for domain name")
		}
		ndomain = u.Hostname()
	}

	if len(ndomain) > 253 {
		return "", errors.New("length for domain name is greater than 253")
	}

	if !isASCII(ndomain) {
		return "", errors.New("invalid domain name")
	}

	if !fqdn.MatchString(ndomain) {
		return "", errors.New("invalid domain name")
	}

	return ndomain, nil
}

// isASCII reports whether s contains only ASCII characters.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}
