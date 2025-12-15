package dto

import (
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// ===============================================
// USER REQUESTS
// ===============================================

// CreateUserRequest represents the request to create a user.
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
	Role     string `json:"role" validate:"required,oneof=admin operator viewer"`
}

// UpdateUserRequest represents the request to update a user.
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin operator viewer"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// ChangePasswordRequest represents the request to change password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ===============================================
// USER RESPONSES
// ===============================================

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserFromEntity converts a domain entity to a response DTO.
func UserFromEntity(u *entity.User) UserResponse {
	return UserResponse{
		ID:          u.ID.String(),
		Email:       u.Email,
		Name:        u.Name,
		Role:        string(u.Role),
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// UsersFromEntities converts a slice of entities to response DTOs.
func UsersFromEntities(users []*entity.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = UserFromEntity(u)
	}
	return result
}
