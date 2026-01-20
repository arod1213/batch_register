package utils

import (
	"cmp"
	"testing"
)

func gt[T cmp.Ordered](t *testing.T, a T, b T) bool {
	if a > b {
		return true
	}
	t.Fatalf("%v not greater than %v", a, b)
	return false
}

func lt[T cmp.Ordered](t *testing.T, a T, b T) bool {
	if a < b {
		return true
	}
	t.Fatalf("%v not greater than %v", a, b)
	return false
}

func eq[T comparable](t *testing.T, a T, b T) bool {
	if a == b {
		return true
	}
	t.Fatalf("%v not greater than %v", a, b)
	return false
}

func TestFuzz(t *testing.T) {
	var matchThresh float32 = 75.0

	a := "Drunk Tank"
	b := "Drunk Tank (Remix)"

	score := FuzzyScore(a, b)
	gt(t, score, matchThresh)

	a = "Drunk Tank"
	b = "Drunk Tank (Remix)"

	score = FuzzyScore(a, b)
	gt(t, score, matchThresh)

	a = "Drunk Tank"
	b = "(Drunk Tank)"

	score = FuzzyScore(a, b)
	eq(t, score, 100)
}
