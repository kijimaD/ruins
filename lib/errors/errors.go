package errors

import (
	"errors"
	"fmt"
)

// Common error types for the ruins game
var (
	// Generic errors
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrOutOfBounds   = errors.New("out of bounds")
	
	// Entity/Component errors
	ErrEntityNotFound    = errors.New("entity not found")
	ErrComponentNotFound = errors.New("component not found")
	ErrInvalidComponent  = errors.New("invalid component")
	
	// Game state errors
	ErrInvalidState     = errors.New("invalid state")
	ErrStateTransition  = errors.New("state transition failed")
	
	// Resource errors
	ErrResourceNotFound = errors.New("resource not found")
	ErrResourceInvalid  = errors.New("resource invalid")
	
	// Data loading errors
	ErrDataCorrupted    = errors.New("data corrupted")
	ErrDataMissing      = errors.New("data missing")
	ErrParsingFailed    = errors.New("parsing failed")
)

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

// ValidationError represents an error during validation
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field %s (value: %v): %s", e.Field, e.Value, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field string, value interface{}, message string) error {
	return ValidationError{Field: field, Value: value, Message: message}
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