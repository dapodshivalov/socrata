package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAverage(t *testing.T) {
	valueExtractorInt := func(x int) float64 { return float64(x) }
	valueExtractorFloat := func(x float64) float64 { return x }

	check(t, 5.5, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, valueExtractorInt)
	check(t, 5.5, []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}, valueExtractorFloat)
	check(t, 7.125, []float64{1, 8.5, 9, 10}, valueExtractorFloat)
	check(t, -2.375, []float64{1, 8.5, -9, -10}, valueExtractorFloat)
	check(t, 0, []float64{}, valueExtractorFloat)
	check(t, 0, []int{0, 0, 0, 0, 0}, valueExtractorInt)
}

func check[T any](t *testing.T, expected float64, elements []T, valueExtractor func(T) float64) {
	average := Average(elements, valueExtractor)
	assert.Equal(t, expected, average)
}
