package utils

import "context"

func MapSlice[T any, T1 any](in []T, f func(T) T1) []T1 {
	out := make([]T1, len(in))

	for k := range in {
		out[k] = f(in[k])
	}

	return out
}

func SafeSliceRange[T any](in []T, start, end int) []T {
	switch {
	case end < start:
		return []T{}
	case start > len(in) || start < 0:
		return []T{}
	case end > len(in):
		end = len(in)
	}

	return in[start:end]
}

// SliceToMapAsKeys возвращает мапу с ключами из элементов слайса и значениями struct{}.
func SliceToMapAsKeys[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{}, len(slice))

	for _, elem := range slice {
		m[elem] = struct{}{}
	}

	return m
}

func ChunkSlice[T any](ctx context.Context, slice []T, size int) <-chan []T {
	ch := make(chan []T)

	go func() {
		defer close(ch)

		for i := 0; i < len(slice); i += size {
			rIndex := i + size
			if rIndex > len(slice) {
				rIndex = len(slice)
			}

			select {
			case ch <- slice[i:rIndex]:
			case <-ctx.Done():
				// если контекст обработчика завершен - можно закончить разбиение на части.
				return
			}
		}
	}()

	return ch
}
