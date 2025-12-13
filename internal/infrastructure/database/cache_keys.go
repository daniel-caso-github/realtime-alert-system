package database

import (
	"fmt"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// CacheKey provides consistent cache key generation.
// Format: {prefix}:{entity}:{identifier}
type CacheKey struct{}

// NewCacheKey creates a new CacheKey helper.
func NewCacheKey() *CacheKey {
	return &CacheKey{}
}

// User returns the cache key for a user by ID.
func (c *CacheKey) User(id entity.ID) string {
	return fmt.Sprintf("user:%s", id.String())
}

// UserByEmail returns the cache key for a user by email.
func (c *CacheKey) UserByEmail(email string) string {
	return fmt.Sprintf("user:email:%s", email)
}

// Alert returns the cache key for an alert by ID.
func (c *CacheKey) Alert(id entity.ID) string {
	return fmt.Sprintf("alert:%s", id.String())
}

// AlertRule returns the cache key for an alert rule by ID.
func (c *CacheKey) AlertRule(id entity.ID) string {
	return fmt.Sprintf("rule:%s", id.String())
}

// AlertRulesEnabled returns the cache key for all enabled rules.
func (c *CacheKey) AlertRulesEnabled() string {
	return "rules:enabled"
}

// NotificationChannel returns the cache key for a channel by ID.
func (c *CacheKey) NotificationChannel(id entity.ID) string {
	return fmt.Sprintf("channel:%s", id.String())
}

// RateLimitUser returns the cache key for user rate limiting.
func (c *CacheKey) RateLimitUser(userID entity.ID, endpoint string) string {
	return fmt.Sprintf("ratelimit:%s:%s", userID.String(), endpoint)
}

// Session returns the cache key for a user session.
func (c *CacheKey) Session(token string) string {
	return fmt.Sprintf("session:%s", token)
}

// BlacklistedToken returns the cache key for a blacklisted JWT.
func (c *CacheKey) BlacklistedToken(tokenID string) string {
	return fmt.Sprintf("blacklist:%s", tokenID)
}

// AlertStatistics returns the cache key for alert statistics.
func (c *CacheKey) AlertStatistics() string {
	return "stats:alerts"
}

// Pattern returns a pattern for matching multiple keys.
// Example: Pattern("user", "*") returns "user:*"
func (c *CacheKey) Pattern(parts ...string) string {
	if len(parts) == 0 {
		return "*"
	}

	key := parts[0]
	for i := 1; i < len(parts); i++ {
		key += ":" + parts[i]
	}

	return key
}
