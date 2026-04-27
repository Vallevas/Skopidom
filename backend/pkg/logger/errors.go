// Package logger defines sentinel errors used across all application layers.
package logger

import "errors"

// Sentinel errors — use errors.Is() to check for these in handlers.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists is returned when a unique constraint would be violated.
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrDisposed is returned when attempting to mutate a disposed item.
	ErrDisposed = errors.New("item is disposed and cannot be modified")

	// ErrForbidden is returned when a user lacks permission for an action.
	ErrForbidden = errors.New("insufficient permissions")

	// ErrInvalidInput is returned when request data fails domain validation.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when a request is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrConflict is returned when a delete operation would violate referential integrity.
	ErrConflict = errors.New("conflict")
)

// BusinessError represents a domain-level validation or business rule violation.
type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

// NewBusinessError creates a new business error with the given message.
func NewBusinessError(message string) error {
	return &BusinessError{Message: message}
}
