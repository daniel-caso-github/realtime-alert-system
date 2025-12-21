package messaging

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	"time"

	"github.com/rs/zerolog/log"
)

// RetryConfig configures the retry behavior.
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
	Jitter         bool
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		Jitter:         true,
	}
}

// Retries handles retry logic with exponential backoff.
type Retries struct {
	config RetryConfig
}

// NewRetries creates a new retries with the given configuration.
func NewRetries(config RetryConfig) *Retries {
	return &Retries{
		config: config,
	}
}

// RetryableFunc is a function that can be retried.
type RetryableFunc func(ctx context.Context) error

// Do executes the function with retry logic.
func (r *Retries) Do(ctx context.Context, operation string, fn RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := r.calculateBackoff(attempt)
			log.Debug().
				Str("operation", operation).
				Int("attempt", attempt).
				Dur("backoff", backoff).
				Msg("Retrying operation")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := fn(ctx)
		if err == nil {
			if attempt > 0 {
				log.Debug().
					Str("operation", operation).
					Int("attempts", attempt+1).
					Msg("Operation succeeded after retry")
			}
			return nil
		}

		lastErr = err
		log.Warn().
			Err(err).
			Str("operation", operation).
			Int("attempt", attempt+1).
			Int("max_retries", r.config.MaxRetries).
			Msg("Operation failed")
	}

	log.Error().
		Err(lastErr).
		Str("operation", operation).
		Int("attempts", r.config.MaxRetries+1).
		Msg("Operation failed after all retries")

	return lastErr
}

// calculateBackoff calculates the backoff duration for the given attempt.
func (r *Retries) calculateBackoff(attempt int) time.Duration {
	backoff := float64(r.config.InitialBackoff) * math.Pow(r.config.Multiplier, float64(attempt-1))

	if r.config.Jitter {
		// Add random jitter (±25%)
		var b [8]byte
		_, err := rand.Read(b[:])
		if err != nil {
			log.Error().Err(err).Msg("failed to generate secure random number for jitter; proceeding without jitter")
			// If crypto/rand fails, default to no jitter
			// This makes (randVal * 2 - 1) effectively 0
			backoff += backoff * 0.25 * (0.5*2 - 1)
		} else {
			val := binary.BigEndian.Uint64(b[:])
			randVal := float64(val) / (math.MaxUint64 + 1.0)
			jitter := backoff * 0.25 * (randVal*2 - 1)
			backoff += jitter
		}
	}

	if backoff > float64(r.config.MaxBackoff) {
		backoff = float64(r.config.MaxBackoff)
	}

	return time.Duration(backoff)
}

// IsRetryable determines if an error is retryable.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Add specific error types that should be retried
	// For now, retry all errors except context cancellation
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return true
}
