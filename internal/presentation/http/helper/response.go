// Package helper provides HTTP utility functions for handlers.
package helper

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
)

// JSON sends a JSON response with the given status code.
func JSON(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}

// Success sends a success response (200 OK).
func Success(c *fiber.Ctx, data interface{}) error {
	return JSON(c, fiber.StatusOK, data)
}

// Created sends a created response (201 Created).
func Created(c *fiber.Ctx, data interface{}) error {
	return JSON(c, fiber.StatusCreated, data)
}

// NoContent sends a no content response (204 No Content).
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error sends an error response.
func Error(c *fiber.Ctx, status int, message string, code string) error {
	requestID, _ := c.Locals("requestid").(string)
	return JSON(c, status, dto.ErrorResponse{
		Error:     message,
		Code:      code,
		Timestamp: time.Now().UTC(),
		RequestID: requestID,
	})
}

// BadRequest sends a 400 Bad Request response.
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message, "BAD_REQUEST")
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, "UNAUTHORIZED")
}

// Forbidden sends a 403 Forbidden response.
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, "FORBIDDEN")
}

// NotFound sends a 404 Not Found response.
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, "NOT_FOUND")
}

// Conflict sends a 409 Conflict response.
func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message, "CONFLICT")
}

// UnprocessableEntity sends a 422 Unprocessable Entity response.
func UnprocessableEntity(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnprocessableEntity, message, "UNPROCESSABLE_ENTITY")
}

// InternalError sends a 500 Internal Server Error response.
func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message, "INTERNAL_ERROR")
}

// ValidationErrors sends a 422 response with field-level errors.
func ValidationErrors(c *fiber.Ctx, errors []ValidationError) error {
	fields := make(map[string]string)
	for _, e := range errors {
		fields[e.Field] = e.Message
	}

	return JSON(c, fiber.StatusUnprocessableEntity, dto.ValidationErrorResponse{
		Error:     "Validation failed",
		Code:      "VALIDATION_ERROR",
		Fields:    fields,
		Timestamp: time.Now().UTC(),
	})
}
