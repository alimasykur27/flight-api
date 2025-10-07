package util

// RemoveDuplicate removes duplicate values from a slice.
func RemoveDuplicate[T comparable](input []T) []T {
	uniqueMap := make(map[T]struct{})
	for _, item := range input {
		uniqueMap[item] = struct{}{}
	}

	uniqueSlice := make([]T, 0, len(uniqueMap))
	for item := range uniqueMap {
		uniqueSlice = append(uniqueSlice, item)
	}

	return uniqueSlice
}
