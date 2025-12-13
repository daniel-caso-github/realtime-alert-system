package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

func TestNewPasswordHash_Success(t *testing.T) {
	// Arrange
	validPassword := "SecurePass123"

	// Act
	hash, err := valueobject.NewPasswordHash(validPassword)

	// Assert
	require.NoError(t, err)
	assert.False(t, hash.IsEmpty())
	assert.NotEqual(t, validPassword, hash.Value())
}

func TestNewPasswordHash_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{"empty", "", valueobject.ErrPasswordEmpty},
		{"too short", "Short1", valueobject.ErrPasswordTooShort},
		{"no uppercase", "lowercase123", valueobject.ErrPasswordNoUppercase},
		{"no lowercase", "UPPERCASE123", valueobject.ErrPasswordNoLowercase},
		{"no number", "NoNumbersHere", valueobject.ErrPasswordNoNumber},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := valueobject.NewPasswordHash(tc.password)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.True(t, hash.IsEmpty())
		})
	}
}

func TestPasswordHash_Verify(t *testing.T) {
	// Arrange
	password := "SecurePass123"
	hash, _ := valueobject.NewPasswordHash(password)

	// Assert
	assert.True(t, hash.Verify(password))
	assert.False(t, hash.Verify("WrongPassword123"))
	assert.False(t, hash.Verify(""))
}

func TestNewPasswordHashFromHash(t *testing.T) {
	existingHash := "$2a$10$someexistinghashfromdatabase"

	hash := valueobject.NewPasswordHashFromHash(existingHash)

	assert.Equal(t, existingHash, hash.Value())
	assert.False(t, hash.IsEmpty())
}
