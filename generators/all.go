// Package generators provides domain typosquatting variation algorithms.
// Each algorithm takes a normalized domain and produces a set of look-alike
// domain strings that could be used for phishing or brand impersonation.
package generators

// Fn is a function that takes a normalized domain string and
// returns the name of the algorithm and a list of generated variants.
type Fn func(string) (string, []string)

// All returns every registered fuzzing algorithm.
func All() []Fn {
	return []Fn{
		omission,
		transposition,
		vowelSwap,
		hyphenation,
		repetition,
		extraLetter,
		withPrefix,
		withSuffix,
		subdomain,
		tldSwap,
		typoTrick,
		phonetic,
		inflect,
		tldCompletion,
	}
}
