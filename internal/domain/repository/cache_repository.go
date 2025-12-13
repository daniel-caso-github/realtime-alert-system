package repository

import (
	"context"
	"time"
)

// CacheRepository defines the cache operations.
// Implemented by Redis in the infrastructure layer.
type CacheRepository interface {
	// Set stores a value with optional TTL.
	// If ttl is 0, the value does not expire.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Get retrieves a value by its key.
	// Returns ErrNotFound if the key doesn't exist or has expired.
	Get(ctx context.Context, key string, dest interface{}) error

	// Delete removes a key.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// SetNX stores only if the key doesn't exist (Set if Not eXists).
	// Useful for distributed locks.
	// Returns true if stored, false if it already existed.
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)

	// Increment increments a counter.
	// If the key doesn't exist, it creates it with value 1.
	Increment(ctx context.Context, key string) (int64, error)

	// Decrement decrements a counter.
	Decrement(ctx context.Context, key string) (int64, error)

	// Expire sets TTL on an existing key.
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// TTL returns the remaining time to live of a key.
	// Returns -1 if the key has no TTL, -2 if it doesn't exist.
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Keys returns all keys matching a pattern.
	// Example: "alert:*" returns "alert:1", "alert:2", etc.
	// ⚠️ Use carefully in production, can be slow.
	Keys(ctx context.Context, pattern string) ([]string, error)

	// DeleteByPattern deletes all keys matching a pattern.
	DeleteByPattern(ctx context.Context, pattern string) error

	// Ping verifies the connection with the cache server.
	Ping(ctx context.Context) error

	// Close closes the connection.
	Close() error
}
