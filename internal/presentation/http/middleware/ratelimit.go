package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// RateLimitConfig holds rate limiter configuration.
type RateLimitConfig struct {
	// Max requests allowed in the window
	Max int
	// Time window duration
	Window time.Duration
	// Key prefix for Redis
	KeyPrefix string
	// Message to show when rate limited
	Message string
}

// DefaultRateLimitConfig returns default rate limit configuration.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:       100,
		Window:    time.Minute,
		KeyPrefix: "ratelimit",
		Message:   "Too many requests, please try again later",
	}
}

// RateLimiter handles request rate limiting using Redis.
type RateLimiter struct {
	cache  repository.CacheRepository
	config RateLimitConfig
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(cache repository.CacheRepository, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		cache:  cache,
		config: config,
	}
}

// Limit returns a middleware that limits requests based on IP.
func (r *RateLimiter) Limit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use IP as identifier
		key := fmt.Sprintf("%s:ip:%s", r.config.KeyPrefix, c.IP())
		return r.checkLimit(c, key)
	}
}

// LimitByUser returns a middleware that limits requests based on user ID.
func (r *RateLimiter) LimitByUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context
		userID, ok := c.Locals("userID").(entity.ID)
		if !ok {
			// Fall back to IP if not authenticated
			key := fmt.Sprintf("%s:ip:%s", r.config.KeyPrefix, c.IP())
			return r.checkLimit(c, key)
		}

		key := fmt.Sprintf("%s:user:%s", r.config.KeyPrefix, userID.String())
		return r.checkLimit(c, key)
	}
}

// LimitByEndpoint returns a middleware that limits requests per endpoint.
func (r *RateLimiter) LimitByEndpoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Combine IP and endpoint
		key := fmt.Sprintf("%s:endpoint:%s:%s:%s", r.config.KeyPrefix, c.IP(), c.Method(), c.Path())
		return r.checkLimit(c, key)
	}
}

// checkLimit checks if the request should be rate limited.
func (r *RateLimiter) checkLimit(c *fiber.Ctx, key string) error {
	ctx := c.Context()

	// Increment counter
	count, err := r.cache.Increment(ctx, key)
	if err != nil {
		// If Redis fails, allow the request (fail open)
		return c.Next()
	}

	// Set expiry on first request
	if count == 1 {
		_ = r.cache.Expire(ctx, key, r.config.Window)
	}

	// Get remaining TTL
	ttl, _ := r.cache.TTL(ctx, key)

	// Set rate limit headers
	c.Set("X-RateLimit-Limit", strconv.Itoa(r.config.Max))
	c.Set("X-RateLimit-Remaining", strconv.Itoa(max(0, r.config.Max-int(count))))
	c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

	// Check if over limit
	if int(count) > r.config.Max {
		c.Set("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
		return helper.Error(c, fiber.StatusTooManyRequests, r.config.Message, "RATE_LIMITED")
	}

	return c.Next()
}

// LoginRateLimiter creates a rate limiter specifically for login attempts.
func LoginRateLimiter(cache repository.CacheRepository) *RateLimiter {
	return NewRateLimiter(cache, RateLimitConfig{
		Max:       5,
		Window:    15 * time.Minute,
		KeyPrefix: "ratelimit:login",
		Message:   "Too many login attempts, please try again in 15 minutes",
	})
}

// APIRateLimiter creates a rate limiter for general API requests.
func APIRateLimiter(cache repository.CacheRepository) *RateLimiter {
	return NewRateLimiter(cache, RateLimitConfig{
		Max:       100,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:api",
		Message:   "Too many requests, please slow down",
	})
}

// AlertCreationRateLimiter creates a rate limiter for creating alerts.
func AlertCreationRateLimiter(cache repository.CacheRepository) *RateLimiter {
	return NewRateLimiter(cache, RateLimitConfig{
		Max:       30,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:alerts",
		Message:   "Too many alerts created, please try again later",
	})
}
