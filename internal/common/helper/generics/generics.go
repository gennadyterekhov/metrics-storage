package generics

func Overwrite[T any](current T, new T) T {
	if isTruthy(new) {
		return new
	}
	return current
}

func isTruthy[T any](value T) bool {
	switch typedValue := any(value).(type) {
	case string:
		return typedValue != ""
	case int64:
		return typedValue != 0
	case bool:
		return typedValue
	case float64:
		return typedValue != 0.0
	}
	return false
}
