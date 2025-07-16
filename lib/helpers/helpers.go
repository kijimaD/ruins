package helpers

// GetPtr は値のポインタを返す
func GetPtr[T any](x T) *T {
	return &x
}
