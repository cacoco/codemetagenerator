package utils

type Predicate[T any] func(T) bool

func Filter[T any](input []T, predicate Predicate[T]) []T {
	out := []T{}
	for _, elem := range input {
		if predicate(elem) {
			out = append(out, elem)
		}
	}
	return out
}
