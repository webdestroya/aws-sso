package util

func Coalesce[T any](values ...*T) *T {
	for _, value := range values {
		if value != nil {
			return value
		}
	}

	return nil
}

func CoalesceString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
