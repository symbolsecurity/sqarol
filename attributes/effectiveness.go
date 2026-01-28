package attributes

import (
	"math"
	"strings"
	"unicode"
)

// Effectiveness computes a score from 0.0 to 1.0 representing how
// hard it is for a human to visually distinguish the variant from
// the original domain. Higher values mean the variation is more
// deceptive and therefore more effective as a typosquat.
//
// The score is a weighted combination of four signals:
//
//   - Edit distance (40%): fewer edits = harder to spot.
//   - Length ratio  (20%): same-length domains are harder to distinguish.
//   - Visual similarity (30%): changes involving visually confusable
//     character pairs (rn/m, 0/o, 1/l, vv/w, etc.) are subtler.
//   - Character class preservation (10%): variants that stay
//     all-alphabetic (no digits or symbols introduced) are subtler.
func Effectiveness(original, variant string, distance int) float64 {
	ed := EditDistanceScore(distance)
	lr := LengthRatioScore(original, variant)
	vs := VisualSimilarityScore(original, variant)
	cc := CharClassScore(original, variant)

	score := ed*0.40 + lr*0.20 + vs*0.30 + cc*0.10

	// Clamp to [0, 1].
	return math.Max(0, math.Min(1, score))
}

// EditDistanceScore returns a value in (0, 1] that is higher when
// the Levenshtein distance is small. A distance of 1 yields 1.0;
// larger distances decay proportionally.
func EditDistanceScore(distance int) float64 {
	if distance <= 0 {
		return 0
	}

	return 1.0 / float64(distance)
}

// LengthRatioScore returns 1.0 when both strings have the same length
// and decreases as the length difference grows.
func LengthRatioScore(a, b string) float64 {
	la := float64(len(a))
	lb := float64(len(b))

	if la == 0 || lb == 0 {
		return 0
	}

	return math.Min(la, lb) / math.Max(la, lb)
}

// VisualSimilarityScore returns a value in [0, 1] that is higher when
// the differences between original and variant involve visually
// confusable character pairs.
func VisualSimilarityScore(original, variant string) float64 {
	// confusables is a list of character-pair substitutions that look
	// nearly identical in common fonts. Each entry is [original, replacement].
	confusables := [][]string{
		{"rn", "m"},
		{"m", "rn"},
		{"vv", "w"},
		{"w", "vv"},
		{"cl", "d"},
		{"d", "cl"},
		{"cj", "g"},
		{"g", "cj"},
		{"ci", "a"},
		{"a", "ci"},
		{"nn", "m"},
		{"o", "0"},
		{"0", "o"},
		{"l", "1"},
		{"1", "l"},
		{"I", "l"},
		{"l", "I"},
	}

	if original == variant {
		return 0
	}

	// Check how many confusable substitutions can explain the difference.
	// We test whether applying confusable replacements to the original
	// can produce the variant.
	transformed := original
	matchCount := 0

	for _, pair := range confusables {
		from := pair[0]
		to := pair[1]

		// If the original contains "from" and the variant contains "to"
		// at the corresponding position, this is a confusable change.
		if strings.Contains(original, from) && strings.Contains(variant, to) {
			candidate := strings.Replace(transformed, from, to, -1)
			if candidate == variant {
				return 1.0
			}

			// Partial match: the confusable explains part of the difference.
			if candidate != transformed {
				transformed = candidate
				matchCount++
			}
		}
	}

	if matchCount == 0 {
		// No confusable pairs found. Check if the edit is a single
		// character swap between visually similar characters
		// (e.g., vowel swaps are somewhat similar).
		if absDiff(len(original), len(variant)) == 0 {
			diffCount := 0
			for i := 0; i < len(original) && i < len(variant); i++ {
				if original[i] != variant[i] {
					diffCount++
				}
			}

			// Single character difference at same length is still mildly
			// confusing even without a known confusable pair.
			if diffCount == 1 {
				return 0.3
			}
		}

		return 0
	}

	// Scale based on how many confusable pairs contributed.
	// Cap at 0.9 for partial matches (full match already returned 1.0).
	score := float64(matchCount) * 0.45
	if score > 0.9 {
		score = 0.9
	}

	return score
}

// CharClassScore returns 1.0 if the variant preserves the same
// character classes as the original (no new digits or symbols
// introduced), and 0.0 if new character classes appear.
func CharClassScore(original, variant string) float64 {
	origHasDigit := false
	origHasSymbol := false

	for _, r := range original {
		if unicode.IsDigit(r) {
			origHasDigit = true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' {
			origHasSymbol = true
		}
	}

	newDigit := false
	newSymbol := false

	for _, r := range variant {
		if unicode.IsDigit(r) && !origHasDigit {
			newDigit = true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' && !origHasSymbol {
			newSymbol = true
		}
	}

	if newDigit && newSymbol {
		return 0.0
	}

	if newDigit || newSymbol {
		return 0.5
	}

	return 1.0
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}

	return b - a
}
