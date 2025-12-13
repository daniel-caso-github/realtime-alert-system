// Package entity defines the core domain entities and shared types
// for the real-time alerting system.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// ID is a type alias for uuid.UUID representing a universally unique identifier.
// It is used as the primary key type for all domain entities.
type ID = uuid.UUID

// NewID generates and returns a new unique identifier using UUID v4.
func NewID() ID {
	return uuid.New()
}

// ParseID parses a string representation into an ID.
// Returns an error if the string is not a valid UUID format.
func ParseID(s string) (ID, error) {
	return uuid.Parse(s)
}

// Timestamps contains common audit fields that should be embedded in all domain entities.
// It provides automatic tracking of creation and modification times.
type Timestamps struct {
	// CreatedAt is the timestamp when the entity was created.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt is the timestamp when the entity was last modified.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewTimestamps creates a new Timestamps instance with both CreatedAt and UpdatedAt
// set to the current UTC time. Should be called when creating new entities.
func NewTimestamps() Timestamps {
	now := time.Now().UTC()
	return Timestamps{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Touch updates the UpdatedAt field to the current UTC time.
// Should be called whenever an entity is modified.
func (t *Timestamps) Touch() {
	t.UpdatedAt = time.Now().UTC()
}
