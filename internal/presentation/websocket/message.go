package websocket

import (
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
)

// MessageType represents the type of WebSocket message.
type MessageType string

// WebSocket message types for client-server communication.
const (
	// Client -> Server
	MessageTypePing        MessageType = "ping"
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"

	// Server -> Client
	MessageTypePong         MessageType = "pong"
	MessageTypeSubscribed   MessageType = "subscribed"
	MessageTypeUnsubscribed MessageType = "unsubscribed"
	MessageTypeError        MessageType = "error"

	// Alert events
	MessageTypeAlertCreated      MessageType = "alert.created"
	MessageTypeAlertUpdated      MessageType = "alert.updated"
	MessageTypeAlertAcknowledged MessageType = "alert.acknowledged"
	MessageTypeAlertResolved     MessageType = "alert.resolved"
	MessageTypeAlertDeleted      MessageType = "alert.deleted"

	// Statistics
	MessageTypeStatsUpdate MessageType = "stats.update"
)

// Message represents a WebSocket message.
type Message struct {
	Type      MessageType `json:"type"`
	Channel   string      `json:"channel,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewAlertCreatedMessage creates a new alert created message.
func NewAlertCreatedMessage(alert dto.AlertResponse) Message {
	return Message{
		Type:      MessageTypeAlertCreated,
		Payload:   alert,
		Timestamp: time.Now().UTC(),
	}
}

// NewAlertUpdatedMessage creates a new alert updated message.
func NewAlertUpdatedMessage(alert dto.AlertResponse) Message {
	return Message{
		Type:      MessageTypeAlertUpdated,
		Payload:   alert,
		Timestamp: time.Now().UTC(),
	}
}

// NewAlertAcknowledgedMessage creates a new alert acknowledged message.
func NewAlertAcknowledgedMessage(alert dto.AlertResponse) Message {
	return Message{
		Type:      MessageTypeAlertAcknowledged,
		Payload:   alert,
		Timestamp: time.Now().UTC(),
	}
}

// NewAlertResolvedMessage creates a new alert resolved message.
func NewAlertResolvedMessage(alert dto.AlertResponse) Message {
	return Message{
		Type:      MessageTypeAlertResolved,
		Payload:   alert,
		Timestamp: time.Now().UTC(),
	}
}

// NewAlertDeletedMessage creates a new alert deleted message.
func NewAlertDeletedMessage(alertID string) Message {
	return Message{
		Type: MessageTypeAlertDeleted,
		Payload: map[string]string{
			"id": alertID,
		},
		Timestamp: time.Now().UTC(),
	}
}

// NewStatsUpdateMessage creates a new statistics update message.
func NewStatsUpdateMessage(stats dto.AlertStatisticsResponse) Message {
	return Message{
		Type:      MessageTypeStatsUpdate,
		Payload:   stats,
		Timestamp: time.Now().UTC(),
	}
}

// NewErrorMessage creates a new error message.
func NewErrorMessage(err string) Message {
	return Message{
		Type: MessageTypeError,
		Payload: map[string]string{
			"error": err,
		},
		Timestamp: time.Now().UTC(),
	}
}
