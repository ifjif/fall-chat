package utils

func ToAnySlice[T any](data []T) []any {
	result := make([]any, len(data))

	for i, v := range data {
		result[i] = v
	}

	return result
}
