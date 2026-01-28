package generators

import "strings"

// inflect generates plural and singular forms of the domain name.
func inflect(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "inflect", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	domainPlural := pluralize(name) + tld
	domainSingular := singularize(name) + tld

	if domainPlural != domain {
		result = append(result, domainPlural)
	}

	if domainSingular != domain {
		result = append(result, domainSingular)
	}

	return "inflect", result
}

// pluralize applies common English pluralization rules to a word.
func pluralize(word string) string {
	if word == "" {
		return word
	}

	lower := strings.ToLower(word)

	// Words ending in s, x, z, ch, sh → add "es"
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z") || strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "sh") {
		return word + "es"
	}

	// Words ending in consonant + y → replace y with "ies"
	if strings.HasSuffix(lower, "y") && len(word) > 1 && !isVowel(word[len(word)-2]) {
		return word[:len(word)-1] + "ies"
	}

	// Words ending in f → replace with "ves"
	if strings.HasSuffix(lower, "f") {
		return word[:len(word)-1] + "ves"
	}

	// Words ending in fe → replace with "ves"
	if strings.HasSuffix(lower, "fe") {
		return word[:len(word)-2] + "ves"
	}

	// Default: add "s"
	return word + "s"
}

// singularize applies common English singularization rules to a word.
func singularize(word string) string {
	if word == "" {
		return word
	}

	lower := strings.ToLower(word)

	// Words ending in "ies" → replace with "y"
	if strings.HasSuffix(lower, "ies") && len(word) > 3 {
		return word[:len(word)-3] + "y"
	}

	// Words ending in "ves" → replace with "f"
	if strings.HasSuffix(lower, "ves") && len(word) > 3 {
		return word[:len(word)-3] + "f"
	}

	// Words ending in "ses", "xes", "zes", "ches", "shes" → strip "es"
	if strings.HasSuffix(lower, "shes") || strings.HasSuffix(lower, "ches") {
		return word[:len(word)-2]
	}

	if (strings.HasSuffix(lower, "ses") || strings.HasSuffix(lower, "xes") ||
		strings.HasSuffix(lower, "zes")) && len(word) > 3 {
		return word[:len(word)-2]
	}

	// Words ending in "s" (but not "ss") → strip "s"
	if strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") && len(word) > 1 {
		return word[:len(word)-1]
	}

	return word
}

func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}
