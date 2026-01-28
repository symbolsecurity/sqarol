// Package generators provides domain typosquatting variation algorithms.
// Each algorithm takes a normalized domain and produces a set of look-alike
// domain strings that could be used for phishing or brand impersonation.
package generators

// All is the registry of every variation generator
// algorithm.
var All = []func(string) (string, []string){
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
