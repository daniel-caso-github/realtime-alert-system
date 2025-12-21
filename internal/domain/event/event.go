// Package event defines domain events for the alerting system.
package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// Type represents the type of domain event.
type Type string

// Event types.
const (
	AlertCreated      Type = "alert.created"
	AlertAcknowledged Type = "alert.acknowledged"
	AlertResolved     Type = "alert.resolved"
	AlertDeleted      Type = "alert.deleted"
	AlertExpired      Type = "alert.expired"
	UserCreated       Type = "user.created"
	UserUpdated       Type = "user.updated"
)

// Event represents a domain event.
type Event struct {
	ID        string          `json:"id"`
	Type      Type            `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
	Version   int             `json:"version"`
	Retries   int             `json:"retries"`
}

// NewEvent creates a new event with the given type and payload.
func NewEvent(eventType Type, payload interface{}) (*Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:        entity.NewID().String(),
		Type:      eventType,
		Payload:   data,
		Timestamp: time.Now().UTC(),
		Version:   1,
		Retries:   0,
	}, nil
}

// UnmarshalPayload unmarshals the event payload into the given target.
func (e *Event) UnmarshalPayload(target interface{}) error {
	return json.Unmarshal(e.Payload, target)
}

// ToMap converts the event to a map for Redis Streams.
func (e *Event) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        e.ID,
		"type":      string(e.Type),
		"payload":   string(e.Payload),
		"timestamp": e.Timestamp.Format(time.RFC3339Nano),
		"version":   e.Version,
		"retries":   e.Retries,
	}
}

// FromMap creates an event from a Redis Streams map.
func FromMap(data map[string]interface{}) (*Event, error) {
	timestamp, err := time.Parse(time.RFC3339Nano, data["timestamp"].(string))
	if err != nil {
		return nil, err
	}

	version := 1
	if v, ok := data["version"]; ok {
		if vi, ok := v.(int64); ok {
			version = int(vi)
		} else if vs, ok := v.(string); ok {
			var vn int
			if _, err := fmt.Sscanf(vs, "%d", &vn); err == nil {
				version = vn
			}
		}
	}

	retries := 0
	if r, ok := data["retries"]; ok {
		if ri, ok := r.(int64); ok {
			retries = int(ri)
		} else if rs, ok := r.(string); ok {
			var rn int
			if _, err := fmt.Sscanf(rs, "%d", &rn); err == nil {
				retries = rn
			}
		}
	}

	return &Event{
		ID:        data["id"].(string),
		Type:      Type(data["type"].(string)),
		Payload:   json.RawMessage(data["payload"].(string)),
		Timestamp: timestamp,
		Version:   version,
		Retries:   retries,
	}, nil
}
