package generators

import "strings"

// VowelSwap generates variations by replacing each vowel in the name
// part with every other vowel.
func vowelSwap(domain string) (string, []string) {
	var result []string
	vowels := "aeiou"

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "vowel-swap", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	for i := 0; i < len(name); i++ {
		if !strings.ContainsRune(vowels, rune(name[i])) {
			continue
		}

		for _, vowel := range vowels {
			if vowel == rune(name[i]) {
				continue
			}
			result = append(result, name[:i]+string(vowel)+name[i+1:]+tld)
		}
	}

	return "vowel-swap", result
}
