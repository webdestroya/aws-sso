package utils

func ToAnySlice[T any](input []T) []any {
	out := make([]any, len(input))
	for i := range input {
		out[i] = input[i]
	}
	return out
}
