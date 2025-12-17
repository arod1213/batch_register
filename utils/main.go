package utils

func Map[T any, K any](x []T, f func(T) K) []K {
	y := make([]K, len(x))
	for i, val := range x {
		y[i] = f(val)
	}
	return y
}
