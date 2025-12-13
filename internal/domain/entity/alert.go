package entity

import (
	"errors"
	"time"
)

// AlertSeverity defines the severity levels for alerts.
// Used to prioritize and categorize alerts by their impact level.
type AlertSeverity string

// Alert severity constants ordered from most to least critical.
const (
	// AlertSeverityCritical indicates a system-critical issue requiring immediate action.
	AlertSeverityCritical AlertSeverity = "critical"
	// AlertSeverityHigh indicates a high-priority issue that needs prompt attention.
	AlertSeverityHigh AlertSeverity = "high"
	// AlertSeverityMedium indicates a moderate issue that should be addressed soon.
	AlertSeverityMedium AlertSeverity = "medium"
	// AlertSeverityLow indicates a low-priority issue that can be addressed later.
	AlertSeverityLow AlertSeverity = "low"
	// AlertSeverityInfo indicates an informational alert with no immediate action required.
	AlertSeverityInfo AlertSeverity = "info"
)

// IsValid checks if the severity is a valid AlertSeverity value.
// Returns true if the severity matches one of the defined constants.
func (s AlertSeverity) IsValid() bool {
	switch s {
	case AlertSeverityCritical, AlertSeverityHigh, AlertSeverityMedium, AlertSeverityLow, AlertSeverityInfo:
		return true
	default:
		return false
	}
}

// Priority returns a numeric value for sorting alerts by severity.
// Lower number indicates higher priority (1 = critical, 5 = info).
func (s AlertSeverity) Priority() int {
	switch s {
	case AlertSeverityCritical:
		return 1
	case AlertSeverityHigh:
		return 2
	case AlertSeverityMedium:
		return 3
	case AlertSeverityLow:
		return 4
	case AlertSeverityInfo:
		return 5
	default:
		return 99
	}
}

// AlertStatus defines the possible states of an alert in its lifecycle.
type AlertStatus string

// Alert status constants representing the alert lifecycle stages.
const (
	// AlertStatusActive indicates a new alert that has not been addressed.
	AlertStatusActive AlertStatus = "active"
	// AlertStatusAcknowledged indicates someone is working on the alert.
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	// AlertStatusResolved indicates the alert has been resolved.
	AlertStatusResolved AlertStatus = "resolved"
	// AlertStatusExpired indicates the alert has passed its expiration time.
	AlertStatusExpired AlertStatus = "expired"
)

// IsValid checks if the status is a valid AlertStatus value.
// Returns true if the status matches one of the defined constants.
func (s AlertStatus) IsValid() bool {
	switch s {
	case AlertStatusActive, AlertStatusAcknowledged, AlertStatusResolved, AlertStatusExpired:
		return true
	default:
		return false
	}
}

// Alert represents an alert in the real-time alerting system.
// It tracks the alert lifecycle from creation through resolution or expiration.
type Alert struct {
	// ID is the unique identifier for the alert.
	ID ID `json:"id" db:"id"`
	// RuleID references the alert rule that triggered this alert (nil if manually created).
	RuleID *ID `json:"rule_id,omitempty" db:"rule_id"`
	// Title is a brief description of the alert (max 255 characters).
	Title string `json:"title" db:"title"`
	// Message contains the detailed alert description.
	Message string `json:"message" db:"message"`
	// Severity indicates the alert's priority level.
	Severity AlertSeverity `json:"severity" db:"severity"`
	// Status indicates the current state of the alert.
	Status AlertStatus `json:"status" db:"status"`
	// Source identifies where the alert originated from.
	Source string `json:"source,omitempty" db:"source"`
	// Metadata stores additional key-value data associated with the alert.
	Metadata map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	// AcknowledgedBy is the ID of the user who acknowledged the alert.
	AcknowledgedBy *ID `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	// AcknowledgedAt is the timestamp when the alert was acknowledged.
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	// ResolvedBy is the ID of the user who resolved the alert.
	ResolvedBy *ID `json:"resolved_by,omitempty" db:"resolved_by"`
	// ResolvedAt is the timestamp when the alert was resolved.
	ResolvedAt *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	// ExpiresAt is the optional expiration time for the alert.
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	// Timestamps embeds creation and update audit fields.
	Timestamps
}

// Alert validation errors.
// Defined as variables to allow comparison using errors.Is().
var (
	ErrAlertTitleRequired       = errors.New("alert title is required")
	ErrAlertTitleTooLong        = errors.New("alert title must be less than 256 characters")
	ErrAlertMessageRequired     = errors.New("alert message is required")
	ErrAlertInvalidSeverity     = errors.New("invalid alert severity")
	ErrAlertInvalidStatus       = errors.New("invalid alert status")
	ErrAlertAlreadyAcknowledged = errors.New("alert is already acknowledged")
	ErrAlertAlreadyResolved     = errors.New("alert is already resolved")
	ErrAlertNotActive           = errors.New("alert is not active")
)

// NewAlert creates a new alert with the provided data and validates it.
// The alert is created with Active status and an empty metadata map.
// Returns an error if validation fails.
func NewAlert(title, message string, severity AlertSeverity, source string) (*Alert, error) {
	alert := &Alert{
		ID:         NewID(),
		Title:      title,
		Message:    message,
		Severity:   severity,
		Status:     AlertStatusActive,
		Source:     source,
		Metadata:   make(map[string]interface{}),
		Timestamps: NewTimestamps(),
	}

	if err := alert.Validate(); err != nil {
		return nil, err
	}

	return alert, nil
}

// Validate checks that all alert fields contain valid data.
// Returns the first validation error encountered, or nil if valid.
func (a *Alert) Validate() error {
	if a.Title == "" {
		return ErrAlertTitleRequired
	}

	if len(a.Title) > 255 {
		return ErrAlertTitleTooLong
	}

	if a.Message == "" {
		return ErrAlertMessageRequired
	}

	if !a.Severity.IsValid() {
		return ErrAlertInvalidSeverity
	}

	if !a.Status.IsValid() {
		return ErrAlertInvalidStatus
	}

	return nil
}

// Acknowledge marks the alert as acknowledged by a user.
// This indicates someone is actively working on the alert.
// Returns an error if the alert is not in Active status.
func (a *Alert) Acknowledge(userID ID) error {
	if a.Status == AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	if a.Status == AlertStatusAcknowledged {
		return ErrAlertAlreadyAcknowledged
	}

	if a.Status != AlertStatusActive {
		return ErrAlertNotActive
	}

	now := time.Now().UTC()
	a.Status = AlertStatusAcknowledged
	a.AcknowledgedBy = &userID
	a.AcknowledgedAt = &now
	a.Touch()

	return nil
}

// Resolve marks the alert as resolved by a user.
// Can be called from any status except already resolved.
// Returns ErrAlertAlreadyResolved if the alert is already resolved.
func (a *Alert) Resolve(userID ID) error {
	if a.Status == AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	now := time.Now().UTC()
	a.Status = AlertStatusResolved
	a.ResolvedBy = &userID
	a.ResolvedAt = &now
	a.Touch()

	return nil
}

// Expire marks the alert as expired.
// Typically called by a background job when the alert passes its expiration time.
func (a *Alert) Expire() {
	a.Status = AlertStatusExpired
	a.Touch()
}

// SetExpiration sets the expiration time for the alert.
// After this time, the alert should be marked as expired.
func (a *Alert) SetExpiration(expiresAt time.Time) {
	a.ExpiresAt = &expiresAt
	a.Touch()
}

// IsExpired checks if the alert has passed its expiration time.
// Returns false if no expiration time is set.
func (a *Alert) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*a.ExpiresAt)
}

// AddMetadata adds a key-value pair to the alert's metadata.
// Creates the metadata map if it doesn't exist.
func (a *Alert) AddMetadata(key string, value interface{}) {
	if a.Metadata == nil {
		a.Metadata = make(map[string]interface{})
	}
	a.Metadata[key] = value
	a.Touch()
}

// IsCritical checks if the alert has critical severity.
// Returns true if the severity is AlertSeverityCritical.
func (a *Alert) IsCritical() bool {
	return a.Severity == AlertSeverityCritical
}

// NeedsImmediateAttention checks if the alert requires immediate attention.
// Returns true if the alert is active and has critical or high severity.
func (a *Alert) NeedsImmediateAttention() bool {
	return a.Status == AlertStatusActive &&
		(a.Severity == AlertSeverityCritical || a.Severity == AlertSeverityHigh)
}
