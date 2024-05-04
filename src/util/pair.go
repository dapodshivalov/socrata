package util

type Pair[A any, B any] struct {
	First  A
	Second B
}

func MakePair[A any, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{First: a, Second: b}
}
