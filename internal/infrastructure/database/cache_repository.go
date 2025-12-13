package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
)

// Ensure RedisCacheRepository implements repository.CacheRepository
var _ repository.CacheRepository = (*RedisCacheRepository)(nil)

// RedisCacheRepository implements CacheRepository using Redis.
type RedisCacheRepository struct {
	client *redis.Client
}

// NewRedisCacheRepository creates a new Redis cache repository.
func NewRedisCacheRepository(redisClient *RedisClient) *RedisCacheRepository {
	return &RedisCacheRepository{
		client: redisClient.Client(),
	}
}

// Set stores a value with optional TTL.
// The value is serialized to JSON before storing.
func (r *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return translateRedisError(err)
	}

	return nil
}

// Get retrieves a value by its key.
// The value is deserialized from JSON into the destination.
func (r *RedisCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return translateRedisError(err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key.
func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return translateRedisError(err)
	}

	return nil
}

// Exists checks if a key exists.
func (r *RedisCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, translateRedisError(err)
	}

	return result > 0, nil
}

// SetNX stores only if the key doesn't exist (Set if Not eXists).
// Returns true if the key was set, false if it already existed.
func (r *RedisCacheRepository) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := r.client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return false, translateRedisError(err)
	}

	return result, nil
}

// Increment increments a counter.
// If the key doesn't exist, it creates it with value 1.
func (r *RedisCacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, translateRedisError(err)
	}

	return result, nil
}

// Decrement decrements a counter.
func (r *RedisCacheRepository) Decrement(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, translateRedisError(err)
	}

	return result, nil
}

// Expire sets TTL on an existing key.
func (r *RedisCacheRepository) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		return translateRedisError(err)
	}

	return nil
}

// TTL returns the remaining time to live of a key.
// Returns -1 if the key has no TTL, -2 if it doesn't exist.
func (r *RedisCacheRepository) TTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, translateRedisError(err)
	}

	return result, nil
}

// Keys returns all keys matching a pattern.
// Warning: Use carefully in production, can be slow with many keys.
func (r *RedisCacheRepository) Keys(ctx context.Context, pattern string) ([]string, error) {
	result, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, translateRedisError(err)
	}

	return result, nil
}

// DeleteByPattern deletes all keys matching a pattern.
// Uses SCAN internally to avoid blocking Redis.
func (r *RedisCacheRepository) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		var batch []string

		// SCAN is non-blocking unlike KEYS
		batch, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return translateRedisError(err)
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			return translateRedisError(err)
		}
	}

	return nil
}

// Ping verifies the connection with Redis.
func (r *RedisCacheRepository) Ping(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return translateRedisError(err)
	}

	return nil
}

// Close closes the connection.
func (r *RedisCacheRepository) Close() error {
	return r.client.Close()
}

// translateRedisError converts Redis errors to domain errors.
func translateRedisError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return repository.ErrNotFound
	}

	// Check for connection errors
	if isRedisConnectionError(err) {
		return repository.ErrConnection
	}

	return err
}

// isRedisConnectionError checks if the error is a connection error.
func isRedisConnectionError(err error) bool {
	if err == nil {
		return false
	}

	// go-redis wraps connection errors
	errMsg := err.Error()
	connectionKeywords := []string{
		"connection refused",
		"connection reset",
		"i/o timeout",
		"network is unreachable",
	}

	for _, keyword := range connectionKeywords {
		if contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// contains checks if s contains substr (case-insensitive).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsLower(toLower(s), toLower(substr))
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func containsLower(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
