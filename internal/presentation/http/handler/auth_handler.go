package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/service"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if errors := helper.ValidateStruct(req); len(errors) > 0 {
		return helper.ValidationErrors(c, errors)
	}

	// Authenticate
	tokens, user, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return helper.Unauthorized(c, "Invalid email or password")
		}
		if errors.Is(err, service.ErrUserNotActive) {
			return helper.Forbidden(c, "Account is deactivated")
		}
		return helper.InternalError(c, "Authentication failed")
	}

	response := dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		User:         dto.UserFromEntity(user),
	}

	return helper.Success(c, response)
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if errors := helper.ValidateStruct(req); len(errors) > 0 {
		return helper.ValidationErrors(c, errors)
	}

	// Register user
	tokens, user, err := h.authService.Register(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return helper.Conflict(c, "Email already registered")
		}
		return helper.InternalError(c, "Registration failed")
	}

	response := dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		User:         dto.UserFromEntity(user),
	}

	return helper.Created(c, response)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if errors := helper.ValidateStruct(req); len(errors) > 0 {
		return helper.ValidationErrors(c, errors)
	}

	// Refresh tokens
	tokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			return helper.Unauthorized(c, "Refresh token has expired")
		}
		if errors.Is(err, service.ErrTokenInvalid) {
			return helper.Unauthorized(c, "Invalid refresh token")
		}
		return helper.InternalError(c, "Token refresh failed")
	}

	response := dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}

	return helper.Success(c, response)
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get tokens from request
	accessToken := c.Get("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	var req dto.RefreshTokenRequest
	_ = c.BodyParser(&req)

	// Logout
	if err := h.authService.Logout(c.Context(), accessToken, req.RefreshToken); err != nil {
		return helper.InternalError(c, "Logout failed")
	}

	return helper.NoContent(c)
}

// Me handles GET /api/v1/auth/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user, ok := c.Locals("user").(*dto.UserResponse)
	if !ok {
		return helper.Unauthorized(c, "User not authenticated")
	}

	return helper.Success(c, user)
}
