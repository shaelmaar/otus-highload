package utils

func Ptr[T any](v T) *T {
	return &v
}

func UnPtr[T any](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}
