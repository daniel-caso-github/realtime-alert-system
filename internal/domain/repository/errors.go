package repository

import "errors"

// Common repository errors.
// These errors are generic and can be used by any repository implementation.
var (
	// ErrNotFound indicates that the requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists indicates that a resource with that identifier already exists.
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrDuplicateKey indicates a unique constraint violation (e.g., duplicate email).
	ErrDuplicateKey = errors.New("duplicate key violation")

	// ErrForeignKeyViolation indicates an attempt to reference a non-existent resource.
	ErrForeignKeyViolation = errors.New("foreign key violation")

	// ErrInvalidData indicates that the provided data is invalid.
	ErrInvalidData = errors.New("invalid data")

	// ErrConnection indicates a connection problem with the storage.
	ErrConnection = errors.New("connection error")

	// ErrTimeout indicates that the operation exceeded the time limit.
	ErrTimeout = errors.New("operation timeout")
)
