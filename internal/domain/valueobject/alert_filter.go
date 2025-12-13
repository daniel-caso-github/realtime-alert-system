package valueobject

import (
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// AlertFilter represents filtering criteria for querying alerts.
// It uses a fluent builder pattern to construct type-safe queries.
// All filter methods return a new AlertFilter, allowing method chaining.
//
// Example usage:
//
//	filter := NewAlertFilter().
//		WithStatuses(entity.AlertStatusActive).
//		WithSeverities(entity.AlertSeverityCritical, entity.AlertSeverityHigh).
//		WithDateRange(startDate, endDate)
type AlertFilter struct {
	// Statuses filters alerts by their current status (e.g., active, acknowledged, resolved).
	Statuses []entity.AlertStatus
	// Severities filters alerts by severity level (e.g., critical, high, medium, low).
	Severities []entity.AlertSeverity
	// Source filters alerts by their originating source system.
	Source *string
	// RuleID filters alerts by the rule that triggered them.
	RuleID *entity.ID
	// FromDate filters alerts created on or after this timestamp.
	FromDate *time.Time
	// ToDate filters alerts created on or before this timestamp.
	ToDate *time.Time
	// Search performs a text search across alert title and message fields.
	Search *string
}

// NewAlertFilter creates an empty AlertFilter with no criteria set.
// Use the fluent builder methods to add filtering criteria.
func NewAlertFilter() AlertFilter {
	return AlertFilter{}
}

// WithStatuses adds a status filter to include only alerts with the specified statuses.
// Multiple statuses can be provided; alerts matching any of them will be included.
func (f AlertFilter) WithStatuses(statuses ...entity.AlertStatus) AlertFilter {
	f.Statuses = statuses
	return f
}

// WithSeverities adds a severity filter to include only alerts with the specified severities.
// Multiple severities can be provided; alerts matching any of them will be included.
func (f AlertFilter) WithSeverities(severities ...entity.AlertSeverity) AlertFilter {
	f.Severities = severities
	return f
}

// WithSource adds a source filter to include only alerts from the specified source system.
func (f AlertFilter) WithSource(source string) AlertFilter {
	f.Source = &source
	return f
}

// WithRuleID adds a rule filter to include only alerts triggered by the specified rule.
func (f AlertFilter) WithRuleID(ruleID entity.ID) AlertFilter {
	f.RuleID = &ruleID
	return f
}

// WithDateRange adds a date range filter to include only alerts within the specified time period.
// Both from and to dates are inclusive.
func (f AlertFilter) WithDateRange(from, to time.Time) AlertFilter {
	f.FromDate = &from
	f.ToDate = &to
	return f
}

// WithSearch adds a text search filter to find alerts matching the search term.
// The search is performed against alert title and message fields.
// Empty search strings are ignored.
func (f AlertFilter) WithSearch(search string) AlertFilter {
	if search != "" {
		f.Search = &search
	}
	return f
}

// ActiveOnly is a convenience method that filters for alerts with active status only.
// Equivalent to WithStatuses(entity.AlertStatusActive).
func (f AlertFilter) ActiveOnly() AlertFilter {
	return f.WithStatuses(entity.AlertStatusActive)
}

// CriticalOnly is a convenience method that filters for critical severity alerts only.
// Equivalent to WithSeverities(entity.AlertSeverityCritical).
func (f AlertFilter) CriticalOnly() AlertFilter {
	return f.WithSeverities(entity.AlertSeverityCritical)
}

// NeedsAttention is a convenience method that filters for alerts requiring immediate attention.
// Returns active alerts with critical or high severity.
func (f AlertFilter) NeedsAttention() AlertFilter {
	return f.WithStatuses(entity.AlertStatusActive).
		WithSeverities(entity.AlertSeverityCritical, entity.AlertSeverityHigh)
}

// HasStatusFilter returns true if at least one status filter is set.
func (f AlertFilter) HasStatusFilter() bool {
	return len(f.Statuses) > 0
}

// HasSeverityFilter returns true if at least one severity filter is set.
func (f AlertFilter) HasSeverityFilter() bool {
	return len(f.Severities) > 0
}

// HasDateFilter returns true if either FromDate or ToDate is set.
func (f AlertFilter) HasDateFilter() bool {
	return f.FromDate != nil || f.ToDate != nil
}

// HasSearch returns true if a non-empty search term is set.
func (f AlertFilter) HasSearch() bool {
	return f.Search != nil && *f.Search != ""
}

// IsEmpty returns true if no filtering criteria are set.
// Useful to determine if a full table scan would be performed.
func (f AlertFilter) IsEmpty() bool {
	return !f.HasStatusFilter() &&
		!f.HasSeverityFilter() &&
		f.Source == nil &&
		f.RuleID == nil &&
		!f.HasDateFilter() &&
		!f.HasSearch()
}
