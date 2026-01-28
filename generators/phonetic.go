package generators

import "strings"

// Phonetic generates variations using phonetically similar letter
// group substitutions, applied both in isolation and in pairwise
// combinations.
func phonetic(domain string) (string, []string) {
	var result []string

	catalog := map[string][]string{
		"F":  {"ph"},
		"C":  {"k", "q"},
		"X":  {"ex"},
		"Ch": {"sh"},
		"G":  {"j"},
		"T":  {"d", "th"},
		"P":  {"b"},
		"V":  {"w"},
		"S":  {"z"},
		"Th": {"t", "d"},
		// Reverse the map
		"Ph": {"f"},
		"K":  {"c"},
		"Q":  {"c"},
		"Ex": {"x"},
		"Sh": {"ch"},
		"J":  {"g"},
		"D":  {"t", "th"},
		"B":  {"p"},
		"W":  {"v"},
		"Z":  {"s"},
	}

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "phonetic", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	// Apply one by one in isolation
	for key, value := range catalog {
		key = strings.ToLower(key)
		if !strings.Contains(name, key) {
			continue
		}

		for _, v := range value {
			result = append(result, strings.Replace(name, key, v, -1)+tld)
		}
	}

	// Apply in combination (2 together)
	for key, value := range catalog {
		key = strings.ToLower(key)
		if !strings.Contains(name, key) {
			continue
		}

		for _, v := range value {
			for key2, value2 := range catalog {
				key2 = strings.ToLower(key2)
				if !strings.Contains(name, key2) {
					continue
				}

				for _, v2 := range value2 {
					nx := strings.Replace(name, key, v, -1)
					nx = strings.Replace(nx, key2, v2, -1)

					result = append(result, nx+tld)
				}
			}
		}
	}

	return "phonetic", result
}
