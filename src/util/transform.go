package util

import "sort"

func Transform[T any, E any](slice []T, mapper func(T) E) []E {
	var result []E
	for _, v := range slice {
		result = append(result, mapper(v))
	}
	return result
}

func SliceToMap[V any, K comparable](slice []V, keyExtractor func(V) K) map[K]V {
	result := make(map[K]V)
	for _, v := range slice {
		result[keyExtractor(v)] = v
	}
	return result
}

func SortedMapKeys[K comparable, V any](data map[K]V, less func(a, b K) bool) []K {
	keys := make([]K, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return less(keys[i], keys[j])
	})
	return keys
}
