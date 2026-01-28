package generators

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// ExtraRandomLetter generates variations by inserting a random letter
// at each position in the name part of the domain.
func extraLetter(domain string) (string, []string) {
	var result []string

	dix := strings.LastIndex(domain, ".")
	if dix == -1 {
		return "extra-random-letter", result
	}

	name := domain[:dix]
	tld := domain[dix:]

	name = strings.TrimPrefix(name, "www.")

	for i := 0; i < len(name); i++ {
		excludedLetters := excludeLettersFromWord(i, name)

		rl := randomLetter(excludedLetters...)
		fuzzyDomain := name[:i] + rl + name[i:] + tld

		if strings.HasPrefix(domain, "www.") {
			fuzzyDomain = fmt.Sprintf("www.%s", fuzzyDomain)
		}

		result = append(result, fuzzyDomain)
	}

	return "extra-random-letter", result
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func randomLetter(exclude ...string) string {
	letters := "abcdefghijklmnopqrstuvwxyz"

	for _, s := range exclude {
		letters = strings.Replace(letters, s, "", 1)
	}

	i := randRange(0, len(letters))

	return string(letters[i])
}

func excludeLettersFromWord(currentIndex int, word string) (excludedLetters []string) {
	var nextLetter string
	var prevLetter string

	if currentIndex-1 > 0 {
		prevLetter = string(word[currentIndex-1])
	}

	if currentIndex+1 < len(word) {
		nextLetter = string(word[currentIndex+1])
	}
	excludedLetters = []string{prevLetter, nextLetter}

	if currentIndex == 0 {
		excludedLetters = append(excludedLetters, string(word[currentIndex]))
	}

	return excludedLetters
}
