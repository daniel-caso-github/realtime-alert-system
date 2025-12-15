// Package dto provides Data Transfer Objects for the application layer.
// DTOs are used to transfer data between the API handlers and the service layer,
// decoupling the external API representation from the internal domain model.
package dto

import (
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// CreateAlertRequest represents the request payload for creating a new alert.
// It contains all required fields for alert creation with validation tags.
type CreateAlertRequest struct {
	Title    string                 `json:"title" validate:"required,max=255"`
	Message  string                 `json:"message" validate:"required"`
	Severity string                 `json:"severity" validate:"required,oneof=critical high medium low info"`
	Source   string                 `json:"source,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateAlertRequest represents the request payload for updating an existing alert.
// All fields are optional (pointers) to support partial updates.
type UpdateAlertRequest struct {
	Title    *string                `json:"title,omitempty" validate:"omitempty,max=255"`
	Message  *string                `json:"message,omitempty"`
	Severity *string                `json:"severity,omitempty" validate:"omitempty,oneof=critical high medium low info"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AcknowledgeAlertRequest represents the request payload for acknowledging an alert.
// Acknowledging an alert indicates that a user has seen and is aware of the alert.
type AcknowledgeAlertRequest struct {
	Note string `json:"note,omitempty"` // Optional note explaining the acknowledgment
}

// ResolveAlertRequest represents the request payload for resolving an alert.
// Resolving an alert marks it as handled and no longer requiring attention.
type ResolveAlertRequest struct {
	Resolution string `json:"resolution,omitempty"` // Optional description of how the alert was resolved
}

// ListAlertsRequest represents query parameters for listing and filtering alerts.
// It supports pagination, filtering by status/severity/source, date range queries,
// text search, and sorting options.
type ListAlertsRequest struct {
	Page      int      `query:"page" validate:"omitempty,min=1"`
	PageSize  int      `query:"page_size" validate:"omitempty,min=1,max=100"`
	Status    []string `query:"status" validate:"omitempty,dive,oneof=active acknowledged resolved expired"`
	Severity  []string `query:"severity" validate:"omitempty,dive,oneof=critical high medium low info"`
	Source    string   `query:"source"`
	Search    string   `query:"search"`
	FromDate  string   `query:"from_date"`
	ToDate    string   `query:"to_date"`
	SortBy    string   `query:"sort_by" validate:"omitempty,oneof=created_at severity status"`
	SortOrder string   `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// AlertResponse represents the API response format for an alert.
// It converts the internal domain entity to a client-friendly JSON structure.
type AlertResponse struct {
	ID             string                 `json:"id"`
	RuleID         *string                `json:"rule_id,omitempty"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Severity       string                 `json:"severity"`
	Status         string                 `json:"status"`
	Source         string                 `json:"source,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	AcknowledgedBy *string                `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
	ResolvedBy     *string                `json:"resolved_by,omitempty"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AlertFromEntity converts a domain Alert entity to an AlertResponse DTO.
// It handles the conversion of internal types (UUIDs, enums) to string representations
// and properly handles optional fields (acknowledged/resolved information).
func AlertFromEntity(a *entity.Alert) AlertResponse {
	response := AlertResponse{
		ID:        a.ID.String(),
		Title:     a.Title,
		Message:   a.Message,
		Severity:  string(a.Severity),
		Status:    string(a.Status),
		Source:    a.Source,
		Metadata:  a.Metadata,
		ExpiresAt: a.ExpiresAt,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if a.RuleID != nil {
		ruleID := a.RuleID.String()
		response.RuleID = &ruleID
	}

	if a.AcknowledgedBy != nil {
		ackBy := a.AcknowledgedBy.String()
		response.AcknowledgedBy = &ackBy
		response.AcknowledgedAt = a.AcknowledgedAt
	}

	if a.ResolvedBy != nil {
		resBy := a.ResolvedBy.String()
		response.ResolvedBy = &resBy
		response.ResolvedAt = a.ResolvedAt
	}

	return response
}

// AlertsFromEntities converts a slice of Alert entities to AlertResponse DTOs.
// It is a convenience function for batch conversion of alert lists.
func AlertsFromEntities(alerts []*entity.Alert) []AlertResponse {
	result := make([]AlertResponse, len(alerts))
	for i, a := range alerts {
		result[i] = AlertFromEntity(a)
	}
	return result
}

// AlertStatisticsResponse represents aggregated alert statistics for dashboards.
// It provides counts by status and breakdowns by severity and source.
type AlertStatisticsResponse struct {
	TotalAlerts        int64            `json:"total_alerts"`        // Total number of alerts in the system
	ActiveAlerts       int64            `json:"active_alerts"`       // Alerts that are currently active
	AcknowledgedAlerts int64            `json:"acknowledged_alerts"` // Alerts that have been acknowledged
	ResolvedAlerts     int64            `json:"resolved_alerts"`     // Alerts that have been resolved
	BySeverity         map[string]int64 `json:"by_severity"`         // Count of alerts grouped by severity level
	BySource           map[string]int64 `json:"by_source"`           // Count of alerts grouped by source
}
