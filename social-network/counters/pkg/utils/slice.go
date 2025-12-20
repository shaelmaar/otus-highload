package utils

func MapSlice[T any, T1 any](in []T, f func(T) T1) []T1 {
	out := make([]T1, len(in))

	for k := range in {
		out[k] = f(in[k])
	}

	return out
}
