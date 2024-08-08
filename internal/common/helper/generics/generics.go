package generics

func Overwrite[T any](current T, new T) T {
	if isTruthy(new) {
		return new
	}
	return current
}

func isTruthy[T any](value T) bool {
	switch typedValue := any(value).(type) {
	case bool:
		return typedValue
	case int:
		return typedValue != 0
	case int64:
		return typedValue != 0
	case float64:
		return typedValue != 0.0
	case string:
		return typedValue != ""
	}
	return false
}
