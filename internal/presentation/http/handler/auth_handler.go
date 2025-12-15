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
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	dto.LoginResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		422		{object}	dto.ValidationErrorResponse
//	@Router			/auth/login [post]
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
//
//	@Summary		Register new user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Registration data"
//	@Success		201		{object}	dto.LoginResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse
//	@Failure		422		{object}	dto.ValidationErrorResponse
//	@Router			/auth/register [post]
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
//
//	@Summary		Refresh tokens
//	@Description	Get new access token using refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token"
//	@Success		200		{object}	dto.TokenResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Router			/auth/refresh [post]
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
//
//	@Summary		Logout user
//	@Description	Invalidate user tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		500	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/logout [post]
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
//
//	@Summary		Get current user
//	@Description	Get authenticated user information
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	dto.UserResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user, ok := c.Locals("user").(*dto.UserResponse)
	if !ok {
		return helper.Unauthorized(c, "User not authenticated")
	}

	return helper.Success(c, user)
}
