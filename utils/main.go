package utils

func Map[T any, K any](x []T, f func(T) K) []K {
	y := make([]K, len(x))
	for i, val := range x {
		y[i] = f(val)
	}
	return y
}

func Reduce[T any, K any](x []T, init K, combine func(K, T) K) K {
	var acc K = init
	for _, v := range x {
		acc = combine(acc, v)
	}
	return acc
}
