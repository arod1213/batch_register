package utils

import (
	"strings"
)

func cleanString(a string) string {
	symbols := []string{"-", "|", "_", "%", "'", "(", ")"}
	var x string = a
	for _, sym := range symbols {
		x = strings.ReplaceAll(x, sym, "")
	}
	return x
}

// Returns out of 100 confidence score
func FuzzyScore(a string, b string) float32 {
	a_words := strings.Split(a, " ")
	b_words := strings.Split(b, " ")

	minLen := min(len(a_words), len(b_words))

	var matches uint = 0
	for _, a_word := range a_words {
		a_lower := cleanString(strings.ToLower(a_word))
		for _, b_word := range b_words {
			b_lower := cleanString(strings.ToLower(b_word))
			if a_lower == b_lower {
				matches++
			}
		}
	}
	matched := float32(matches) / float32(minLen) * 100.0
	return matched
}
