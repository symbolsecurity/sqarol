package attributes

// Levenshtein computes the Levenshtein distance between two strings,
// which is the minimum number of single-character edits (insertions,
// deletions, or substitutions) required to change one into the other.
func Levenshtein(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)

	la := len(ra)
	lb := len(rb)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Two-row approach to save memory.
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}

			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost

			curr[j] = min(del, ins, sub)
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}

func min(vals ...int) int {
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
