package utils

// GroupBy is a generic helper to group items by a key function
func GroupBy[T any, K comparable](items []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range items {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}
