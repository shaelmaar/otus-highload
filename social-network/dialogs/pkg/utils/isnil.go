package utils

import "reflect"

// IsNil проверка значения на nil.
func IsNil(i any) bool {
	if i == nil {
		return true
	}

	defer func() {
		_ = recover()
	}()

	return reflect.ValueOf(i).IsNil()
}
