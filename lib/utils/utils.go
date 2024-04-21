package utils

func GetPtr[T any](x T) *T {
	return &x
}
