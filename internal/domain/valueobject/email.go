// Package valueobject contains immutable value objects for the domain layer.
// Value objects are identified by their attributes rather than by a unique identity.
package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Email validation errors.
var (
	// ErrEmailEmpty is returned when attempting to create an Email with an empty string.
	ErrEmailEmpty = errors.New("email cannot be empty")
	// ErrEmailInvalid is returned when the email format does not match the expected pattern.
	ErrEmailInvalid = errors.New("invalid email format")
	// ErrEmailTooLong is returned when the email exceeds 254 characters.
	ErrEmailTooLong = errors.New("email must be less than 255 characters")
)

// emailRegex is the regular expression pattern used to validate email format.
// This pattern covers most valid cases according to RFC 5322.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email represents a validated email address.
// It is an immutable value object: once created, its value cannot be changed.
// The email is automatically normalized (trimmed and lowercased) upon creation.
type Email struct {
	value string
}

// NewEmail creates a new validated Email value object.
// The input is normalized by trimming whitespace and converting to lowercase.
//
// Validation rules:
//   - Email cannot be empty
//   - Email must not exceed 254 characters
//   - Email must match the RFC 5322 format pattern
//
// Returns the Email and nil on success, or a zero Email and an error if validation fails.
func NewEmail(value string) (Email, error) {
	// Normalize: trim spaces and convert to lowercase
	normalized := strings.ToLower(strings.TrimSpace(value))

	if normalized == "" {
		return Email{}, ErrEmailEmpty
	}

	if len(normalized) > 254 {
		return Email{}, ErrEmailTooLong
	}

	if !emailRegex.MatchString(normalized) {
		return Email{}, ErrEmailInvalid
	}

	return Email{value: normalized}, nil
}

// String returns the email value as a string.
// Implements fmt.Stringer interface for proper string representation.
func (e Email) String() string {
	return e.value
}

// Value returns the internal email value.
// Useful for database persistence and serialization.
func (e Email) Value() string {
	return e.value
}

// Equals compares two Email value objects for equality.
// Two emails are considered equal if they have the same normalized value.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Domain extracts and returns the domain part of the email address.
// Example: "user@gmail.com" returns "gmail.com".
// Returns an empty string if the email format is invalid.
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart extracts and returns the local part of the email address (before the @).
// Example: "user@gmail.com" returns "user".
// Returns an empty string if the email format is invalid.
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// IsEmpty checks whether the email value is empty.
// Returns true if the email was not properly initialized or has no value.
func (e Email) IsEmpty() bool {
	return e.value == ""
}
