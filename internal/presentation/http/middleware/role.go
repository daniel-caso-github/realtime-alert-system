package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// RequireRole returns a middleware that checks if user has one of the required roles.
func RequireRole(roles ...entity.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return helper.Unauthorized(c, "User not authenticated")
		}

		// Check if user has one of the required roles
		for _, role := range roles {
			if string(role) == userRole {
				return c.Next()
			}
		}

		return helper.Forbidden(c, "Insufficient permissions")
	}
}

// RequireAdmin is a shortcut for RequireRole(UserRoleAdmin).
func RequireAdmin() fiber.Handler {
	return RequireRole(entity.UserRoleAdmin)
}

// RequireOperator is a shortcut for RequireRole(UserRoleAdmin, UserRoleOperator).
func RequireOperator() fiber.Handler {
	return RequireRole(entity.UserRoleAdmin, entity.UserRoleOperator)
}
