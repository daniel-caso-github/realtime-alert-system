package entity

import (
	"errors"
	"regexp"
	"time"
)

// UserRole defines the possible roles a user can have in the system.
// It uses a custom type for type-safety instead of raw strings.
type UserRole string

// User role constants define the available roles in the system.
const (
	// UserRoleAdmin has full system access including user management.
	UserRoleAdmin UserRole = "admin"
	// UserRoleOperator can manage alerts and notification channels.
	UserRoleOperator UserRole = "operator"
	// UserRoleViewer has read-only access to the system.
	UserRoleViewer UserRole = "viewer"
)

// IsValid checks if the role is a valid system role.
// Returns true if the role matches one of the defined constants.
func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin, UserRoleOperator, UserRoleViewer:
		return true
	default:
		return false
	}
}

// User represents a user in the alerting system.
// It contains authentication data, profile information, and role-based access control.
type User struct {
	// ID is the unique identifier for the user.
	ID ID `json:"id" db:"id"`
	// Email is the user's email address, used for login and notifications.
	Email string `json:"email" db:"email"`
	// PasswordHash stores the hashed password (excluded from JSON serialization).
	PasswordHash string `json:"-" db:"password_hash"`
	// Name is the user's display name.
	Name string `json:"name" db:"name"`
	// Role defines the user's permissions level.
	Role UserRole `json:"role" db:"role"`
	// IsActive indicates whether the user account is enabled.
	IsActive bool `json:"is_active" db:"is_active"`
	// LastLoginAt records the timestamp of the user's last login (nil if never logged in).
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	// Timestamps embeds creation and update audit fields.
	Timestamps
}

// User validation errors.
// Defined as variables to allow comparison using errors.Is().
var (
	ErrUserInvalidEmail     = errors.New("invalid email format")
	ErrUserEmailRequired    = errors.New("email is required")
	ErrUserNameRequired     = errors.New("name is required")
	ErrUserNameTooShort     = errors.New("name must be at least 2 characters")
	ErrUserInvalidRole      = errors.New("invalid user role")
	ErrUserPasswordRequired = errors.New("password hash is required")
)

// emailRegex is the regular expression pattern for validating email format.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewUser creates a new user with the provided data and validates it.
// This is the only entry point for creating users, ensuring they are always valid.
// Returns an error if validation fails.
func NewUser(email, passwordHash, name string, role UserRole) (*User, error) {
	user := &User{
		ID:           NewID(),
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		Role:         role,
		IsActive:     true,
		LastLoginAt:  nil,
		Timestamps:   NewTimestamps(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate checks that all user fields contain valid data.
// Can be called at any time to verify the entity's state.
// Returns the first validation error encountered, or nil if valid.
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrUserEmailRequired
	}

	if !emailRegex.MatchString(u.Email) {
		return ErrUserInvalidEmail
	}

	if u.Name == "" {
		return ErrUserNameRequired
	}

	if len(u.Name) < 2 {
		return ErrUserNameTooShort
	}

	if !u.Role.IsValid() {
		return ErrUserInvalidRole
	}

	if u.PasswordHash == "" {
		return ErrUserPasswordRequired
	}

	return nil
}

// Deactivate disables the user account, preventing login.
// Automatically updates the UpdatedAt timestamp.
func (u *User) Deactivate() {
	u.IsActive = false
	u.Touch()
}

// Activate enables the user account, allowing login.
// Automatically updates the UpdatedAt timestamp.
func (u *User) Activate() {
	u.IsActive = true
	u.Touch()
}

// UpdateLastLogin records the current time as the user's last login.
// Automatically updates the UpdatedAt timestamp.
func (u *User) UpdateLastLogin() {
	now := time.Now().UTC()
	u.LastLoginAt = &now
	u.Touch()
}

// ChangeRole updates the user's role after validating the new role.
// Returns ErrUserInvalidRole if the role is not valid.
// Automatically updates the UpdatedAt timestamp on success.
func (u *User) ChangeRole(newRole UserRole) error {
	if !newRole.IsValid() {
		return ErrUserInvalidRole
	}
	u.Role = newRole
	u.Touch()
	return nil
}

// IsAdmin checks if the user has administrator privileges.
// Returns true if the user's role is UserRoleAdmin.
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// CanManageAlerts checks if the user has permission to manage alerts.
// Returns true if the user is an admin or operator.
func (u *User) CanManageAlerts() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleOperator
}
