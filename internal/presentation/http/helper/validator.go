// Package helper provides HTTP utility functions for handlers.
package helper

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Validator is a shared validator instance.
var Validator = validator.New()

// ValidationError represents a field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns formatted errors.
func ValidateStruct(s interface{}) []ValidationError {
	var validationErrors []ValidationError

	err := Validator.Struct(s)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				validationErrors = append(validationErrors, ValidationError{
					Field:   toSnakeCase(fe.Field()),
					Message: getErrorMessage(fe),
				})
			}
		}
	}

	return validationErrors
}

// getErrorMessage returns a human-readable error message.
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum: " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + err.Param() + ")"
	case "oneof":
		return "Value must be one of: " + err.Param()
	default:
		return "Invalid value"
	}
}

// toSnakeCase converts CamelCase to snake_case.
func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c+'a'-'A'))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
