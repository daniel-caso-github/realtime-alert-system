package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	applogger "github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/logger"
)

// RequestLogger returns a middleware that logs HTTP requests.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get or generate request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = c.Locals("requestid").(string)
		}

		// Add request ID to context
		ctx := applogger.WithRequestID(c.Context(), requestID)
		c.SetUserContext(ctx)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response status
		status := c.Response().StatusCode()

		// Build log event
		event := log.Info()
		if status >= 500 {
			event = log.Error()
		} else if status >= 400 {
			event = log.Warn()
		}

		// Log request details
		event.
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Dur("duration", duration).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("HTTP request")

		return err
	}
}

// AddUserToContext adds user ID to the logging context after authentication.
func AddUserToContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get user ID from locals (set by auth middleware)
		if userID := c.Locals("userID"); userID != nil {
			ctx := applogger.WithUserID(c.UserContext(), userID.(string))
			c.SetUserContext(ctx)
		}
		return c.Next()
	}
}
