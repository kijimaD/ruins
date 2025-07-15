package errors

import "fmt"

// KeyNotFoundError represents an error when a key is not found in a collection
type KeyNotFoundError struct {
	Key        string
	Collection string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: %s in %s", e.Key, e.Collection)
}

// NewKeyNotFoundError creates a new KeyNotFoundError
func NewKeyNotFoundError(key, collection string) error {
	return KeyNotFoundError{Key: key, Collection: collection}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an error with formatted additional context
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
