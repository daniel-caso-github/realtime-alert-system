package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

func TestNewEmail_Success(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"test@example.com", "test@example.com"},
		{"TEST@EXAMPLE.COM", "test@example.com"},
		{"  test@example.com  ", "test@example.com"},
		{"user.name+tag@domain.co.uk", "user.name+tag@domain.co.uk"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			email, err := valueobject.NewEmail(tc.input)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, email.Value())
		})
	}
}

func TestNewEmail_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"empty", "", valueobject.ErrEmailEmpty},
		{"only spaces", "   ", valueobject.ErrEmailEmpty},
		{"no @", "invalidemail", valueobject.ErrEmailInvalid},
		{"no domain", "test@", valueobject.ErrEmailInvalid},
		{"no local part", "@example.com", valueobject.ErrEmailInvalid},
		{"double @", "test@@example.com", valueobject.ErrEmailInvalid},
		{"no tld", "test@example", valueobject.ErrEmailInvalid},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			email, err := valueobject.NewEmail(tc.input)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.True(t, email.IsEmpty())
		})
	}
}

func TestEmail_Domain(t *testing.T) {
	email, _ := valueobject.NewEmail("user@gmail.com")
	assert.Equal(t, "gmail.com", email.Domain())
}

func TestEmail_LocalPart(t *testing.T) {
	email, _ := valueobject.NewEmail("user@gmail.com")
	assert.Equal(t, "user", email.LocalPart())
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := valueobject.NewEmail("test@example.com")
	email2, _ := valueobject.NewEmail("TEST@EXAMPLE.COM")
	email3, _ := valueobject.NewEmail("other@example.com")

	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}
