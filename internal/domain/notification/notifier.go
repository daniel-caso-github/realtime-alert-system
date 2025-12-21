// Package notification defines notification interfaces and types.
package notification

import "context"

// Severity levels for notifications.
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

// Message represents a notification message.
type Message struct {
	Title    string
	Text     string
	Severity string
	Fields   map[string]string
	AlertID  string
	Source   string
}

// Notifier defines the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, msg Message) error
	Name() string
	IsEnabled() bool
}

// SeverityPriority returns the priority of a severity level (lower is higher priority).
func SeverityPriority(severity string) int {
	switch severity {
	case SeverityCritical:
		return 1
	case SeverityHigh:
		return 2
	case SeverityMedium:
		return 3
	case SeverityLow:
		return 4
	case SeverityInfo:
		return 5
	default:
		return 99
	}
}

// ShouldNotify returns true if the severity meets the minimum threshold.
func ShouldNotify(severity, minSeverity string) bool {
	return SeverityPriority(severity) <= SeverityPriority(minSeverity)
}
