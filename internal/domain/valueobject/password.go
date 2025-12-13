package valueobject

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Password validation errors.
var (
	// ErrPasswordEmpty is returned when attempting to create a password hash with an empty string.
	ErrPasswordEmpty = errors.New("password cannot be empty")
	// ErrPasswordTooShort is returned when the password has fewer than 8 characters.
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	// ErrPasswordTooLong is returned when the password exceeds 72 characters (bcrypt limit).
	ErrPasswordTooLong = errors.New("password must be less than 72 characters")
	// ErrPasswordNoUppercase is returned when the password lacks an uppercase letter.
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	// ErrPasswordNoLowercase is returned when the password lacks a lowercase letter.
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	// ErrPasswordNoNumber is returned when the password lacks a numeric digit.
	ErrPasswordNoNumber = errors.New("password must contain at least one number")
	// ErrPasswordHashFailed is returned when bcrypt fails to generate the hash.
	ErrPasswordHashFailed = errors.New("failed to hash password")
	// ErrPasswordInvalid is returned when password verification fails.
	ErrPasswordInvalid = errors.New("invalid password")
)

// PasswordHash represents a securely hashed password.
// It is an immutable value object that never stores the plain text password.
// The hash is generated using bcrypt with the default cost factor.
type PasswordHash struct {
	hash string
}

// NewPasswordHash creates a new PasswordHash from a plain text password.
// It validates password strength requirements before hashing.
//
// Password strength requirements:
//   - Must not be empty
//   - Must be at least 8 characters long
//   - Must not exceed 72 characters (bcrypt limitation)
//   - Must contain at least one uppercase letter
//   - Must contain at least one lowercase letter
//   - Must contain at least one number
//
// Returns the PasswordHash and nil on success, or a zero PasswordHash and an error
// if validation fails or hashing encounters an error.
func NewPasswordHash(plainPassword string) (PasswordHash, error) {
	// Validate password strength before hashing
	if err := validatePasswordStrength(plainPassword); err != nil {
		return PasswordHash{}, err
	}

	// Hash with bcrypt (default cost of 10 provides a good security/performance balance)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return PasswordHash{}, ErrPasswordHashFailed
	}

	return PasswordHash{hash: string(hashedBytes)}, nil
}

// NewPasswordHashFromHash creates a PasswordHash from an existing hash string.
// This is useful when loading a previously hashed password from a database.
// No validation is performed on the hash; it is assumed to be valid.
func NewPasswordHashFromHash(hash string) PasswordHash {
	return PasswordHash{hash: hash}
}

// validatePasswordStrength checks that the password meets minimum security requirements.
// This function is internal and is called by NewPasswordHash before hashing.
func validatePasswordStrength(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if len(password) > 72 {
		return ErrPasswordTooLong
	}

	var hasUpper, hasLower, hasNumber bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}

	if !hasLower {
		return ErrPasswordNoLowercase
	}

	if !hasNumber {
		return ErrPasswordNoNumber
	}

	return nil
}

// Verify compares a plain text password against the stored hash.
// It uses bcrypt's constant-time comparison to prevent timing attacks.
// Returns true if the password matches the hash, false otherwise.
func (p PasswordHash) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainPassword))
	return err == nil
}

// String returns the hash string representation (NOT the original password).
// Implements fmt.Stringer interface. Useful for database persistence.
func (p PasswordHash) String() string {
	return p.hash
}

// Value returns the internal hash value as a string.
// Useful for database persistence and serialization.
func (p PasswordHash) Value() string {
	return p.hash
}

// IsEmpty checks whether the password hash is empty.
// Returns true if no hash has been set or if the PasswordHash was not properly initialized.
func (p PasswordHash) IsEmpty() bool {
	return p.hash == ""
}
