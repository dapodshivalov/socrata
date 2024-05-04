package metrics

func Average[T any](elements []T, valueExtractor func(T) float64) float64 {
	if len(elements) == 0 {
		return 0
	}
	result := 0.0
	for _, element := range elements {
		value := valueExtractor(element)
		result += value
	}
	return result / float64(len(elements))
}
