// Package middleware provides HTTP middleware functions.
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/service"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// AuthMiddleware handles JWT authentication.
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate validates the JWT token and sets user info in context.
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helper.Unauthorized(c, "Missing authorization header")
	}

	// Check Bearer prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return helper.Unauthorized(c, "Invalid authorization header format")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return helper.Unauthorized(c, "Missing token")
	}

	// Validate token
	claims, err := m.authService.ValidateToken(c.Context(), token)
	if err != nil {
		return helper.Unauthorized(c, "Invalid or expired token")
	}

	// Parse user ID
	userID, err := entity.ParseID(claims.UserID)
	if err != nil {
		return helper.Unauthorized(c, "Invalid token claims")
	}

	// Set user info in context for handlers to use
	c.Locals("userID", userID)
	c.Locals("userEmail", claims.Email)
	c.Locals("userRole", claims.Role)
	c.Locals("user", &dto.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	})

	return c.Next()
}

// OptionalAuth validates JWT if present, but allows unauthenticated requests.
func (m *AuthMiddleware) OptionalAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Next()
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return c.Next()
	}

	claims, err := m.authService.ValidateToken(c.Context(), token)
	if err != nil {
		return c.Next()
	}

	userID, err := entity.ParseID(claims.UserID)
	if err != nil {
		return c.Next()
	}

	c.Locals("userID", userID)
	c.Locals("userEmail", claims.Email)
	c.Locals("userRole", claims.Role)
	c.Locals("user", &dto.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	})

	return c.Next()
}
