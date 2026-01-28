package sqarol

// Variation represents a single domain variation produced
// by one of the fuzzing algorithms.
type Variation struct {
	// Original is the normalized input domain.
	Original string `json:"original"`
	// Variant is the generated typo-squatted domain.
	Variant string `json:"variant"`
	// Kind identifies the algorithm that produced this variation.
	Kind string `json:"kind"`
	// Distance is the Levenshtein distance between Original and Variant.
	Distance int `json:"distance"`
	// Effectiveness is a score from 0.0 to 1.0 representing how hard
	// it is for a human to visually detect the difference. Higher values
	// mean the variation is more deceptive.
	Effectiveness float64 `json:"effectiveness"`
}
