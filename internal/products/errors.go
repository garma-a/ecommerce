package products

import (
	"errors"
	"fmt"
)

var (
	// ErrValidation indicates input validation failure
	ErrValidation = errors.New("validation error")

	// ErrNotFound indicates requested resource not found
	ErrNotFound = errors.New("product not found")

	// ErrInternal indicates internal service error
	ErrInternal = errors.New("internal service error")
)

// ValidationError wraps validation errors with field details
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
