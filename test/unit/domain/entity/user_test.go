package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

func TestNewUser_Success(t *testing.T) {
	// Arrange
	email := "test@example.com"
	passwordHash := "$2a$10$validhashhere"
	name := "John Doe"
	role := entity.UserRoleOperator

	// Act
	user, err := entity.NewUser(email, passwordHash, name, role)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, entity.ID{}, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, passwordHash, user.PasswordHash)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, role, user.Role)
	assert.True(t, user.IsActive)
	assert.Nil(t, user.LastLoginAt)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestNewUser_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		passwordHash string
		userName     string
		role         entity.UserRole
		expectedErr  error
	}{
		{
			name:         "empty email",
			email:        "",
			passwordHash: "hash",
			userName:     "John",
			role:         entity.UserRoleViewer,
			expectedErr:  entity.ErrUserEmailRequired,
		},
		{
			name:         "invalid email format",
			email:        "invalid-email",
			passwordHash: "hash",
			userName:     "John",
			role:         entity.UserRoleViewer,
			expectedErr:  entity.ErrUserInvalidEmail,
		},
		{
			name:         "empty name",
			email:        "test@example.com",
			passwordHash: "hash",
			userName:     "",
			role:         entity.UserRoleViewer,
			expectedErr:  entity.ErrUserNameRequired,
		},
		{
			name:         "name too short",
			email:        "test@example.com",
			passwordHash: "hash",
			userName:     "J",
			role:         entity.UserRoleViewer,
			expectedErr:  entity.ErrUserNameTooShort,
		},
		{
			name:         "invalid role",
			email:        "test@example.com",
			passwordHash: "hash",
			userName:     "John",
			role:         entity.UserRole("invalid"),
			expectedErr:  entity.ErrUserInvalidRole,
		},
		{
			name:         "empty password hash",
			email:        "test@example.com",
			passwordHash: "",
			userName:     "John",
			role:         entity.UserRoleViewer,
			expectedErr:  entity.ErrUserPasswordRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := entity.NewUser(tc.email, tc.passwordHash, tc.userName, tc.role)

			assert.Nil(t, user)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestUserRole_IsValid(t *testing.T) {
	testCases := []struct {
		role     entity.UserRole
		expected bool
	}{
		{entity.UserRoleAdmin, true},
		{entity.UserRoleOperator, true},
		{entity.UserRoleViewer, true},
		{entity.UserRole("invalid"), false},
		{entity.UserRole(""), false},
	}

	for _, tc := range testCases {
		t.Run(string(tc.role), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.role.IsValid())
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	// Arrange
	user, _ := entity.NewUser("test@example.com", "hash", "John", entity.UserRoleViewer)
	originalUpdatedAt := user.UpdatedAt

	// Act
	user.Deactivate()

	// Assert
	assert.False(t, user.IsActive)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt) || user.UpdatedAt.Equal(originalUpdatedAt))
}

func TestUser_Activate(t *testing.T) {
	// Arrange
	user, _ := entity.NewUser("test@example.com", "hash", "John", entity.UserRoleViewer)
	user.Deactivate()

	// Act
	user.Activate()

	// Assert
	assert.True(t, user.IsActive)
}

func TestUser_UpdateLastLogin(t *testing.T) {
	// Arrange
	user, _ := entity.NewUser("test@example.com", "hash", "John", entity.UserRoleViewer)
	assert.Nil(t, user.LastLoginAt)

	// Act
	user.UpdateLastLogin()

	// Assert
	assert.NotNil(t, user.LastLoginAt)
}

func TestUser_ChangeRole(t *testing.T) {
	// Arrange
	user, _ := entity.NewUser("test@example.com", "hash", "John", entity.UserRoleViewer)

	// Act & Assert - valid role
	err := user.ChangeRole(entity.UserRoleAdmin)
	assert.NoError(t, err)
	assert.Equal(t, entity.UserRoleAdmin, user.Role)

	// Act & Assert - invalid role
	err = user.ChangeRole(entity.UserRole("invalid"))
	assert.ErrorIs(t, err, entity.ErrUserInvalidRole)
}

func TestUser_IsAdmin(t *testing.T) {
	admin, _ := entity.NewUser("admin@example.com", "hash", "Admin", entity.UserRoleAdmin)
	viewer, _ := entity.NewUser("viewer@example.com", "hash", "Viewer", entity.UserRoleViewer)

	assert.True(t, admin.IsAdmin())
	assert.False(t, viewer.IsAdmin())
}

func TestUser_CanManageAlerts(t *testing.T) {
	admin, _ := entity.NewUser("admin@example.com", "hash", "Admin", entity.UserRoleAdmin)
	operator, _ := entity.NewUser("operator@example.com", "hash", "Operator", entity.UserRoleOperator)
	viewer, _ := entity.NewUser("viewer@example.com", "hash", "Viewer", entity.UserRoleViewer)

	assert.True(t, admin.CanManageAlerts())
	assert.True(t, operator.CanManageAlerts())
	assert.False(t, viewer.CanManageAlerts())
}
