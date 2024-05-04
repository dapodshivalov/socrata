package util

func Min[T int | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}
