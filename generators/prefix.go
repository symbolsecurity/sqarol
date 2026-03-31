package generators

import "strings"

// WithPrefix generates variations by prepending common phishing-related
// prefixes to the domain name.
func withPrefix(domain string) (string, []string) {
	prefixes := []string{
		"account",
		"accounts",
		"account",
		"account-verify",
		"auth-access",
		"auth",
		"auth-user",
		"billing",
		"www",
		"client",
		"portal",
		"access",
		"login",
		"signin",
		"settings",
		"management",
		"control",
		"panel",
		"dashboard",
		"user",
		"menu",
	}

	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "prefix", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for _, prefix := range prefixes {
		// Skip www prefix if name already starts with www
		if prefix == "www" && strings.HasPrefix(name, "www") {
			continue
		}
		result = append(result, prefix+"-"+name+tld)
	}

	return "prefix", result
}
