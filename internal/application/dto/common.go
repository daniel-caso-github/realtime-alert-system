package dto

import "time"

// ===============================================
// COMMON RESPONSES
// ===============================================

// PaginatedResponse wraps paginated data with metadata.
type PaginatedResponse[T any] struct {
	Items       []T   `json:"items"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error     string            `json:"error"`
	Code      string            `json:"code,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	RequestID string            `json:"request_id,omitempty"`
}

// NewErrorResponse creates a new error response.
func NewErrorResponse(err string, code string, requestID string) ErrorResponse {
	return ErrorResponse{
		Error:     err,
		Code:      code,
		Timestamp: time.Now().UTC(),
		RequestID: requestID,
	}
}

// ValidationErrorResponse represents validation errors.
type ValidationErrorResponse struct {
	Error     string            `json:"error"`
	Code      string            `json:"code"`
	Fields    map[string]string `json:"fields"`
	Timestamp time.Time         `json:"timestamp"`
}

// ===============================================
// HEALTH RESPONSES
// ===============================================

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// ReadyResponse represents the readiness check response.
type ReadyResponse struct {
	Status string `json:"status"`
}

// LiveResponse represents the liveness check response.
type LiveResponse struct {
	Status string `json:"status"`
}
