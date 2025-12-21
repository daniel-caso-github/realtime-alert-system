package event

import "time"

// AlertPayload represents the payload for alert events.
type AlertPayload struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Severity       string                 `json:"severity"`
	Status         string                 `json:"status"`
	Source         string                 `json:"source"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	AcknowledgedBy *string                `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
	ResolvedBy     *string                `json:"resolved_by,omitempty"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

// AlertDeletedPayload represents the payload for alert deleted events.
type AlertDeletedPayload struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy string    `json:"deleted_by,omitempty"`
}

// UserPayload represents the payload for user events.
type UserPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// NotificationPayload represents the payload for notification events.
type NotificationPayload struct {
	Channel   string                 `json:"channel"`
	Recipient string                 `json:"recipient"`
	Subject   string                 `json:"subject"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
