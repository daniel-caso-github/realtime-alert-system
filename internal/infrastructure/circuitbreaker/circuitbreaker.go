// Package circuitbreaker provides circuit breaker implementation for external services.
package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// State represents the circuit breaker state.
type State int

// Circuit breaker states.
const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// String returns the string representation of the state.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Errors.
var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyFailures = errors.New("too many failures")
)

// Config holds circuit breaker configuration.
type Config struct {
	Name             string
	MaxFailures      int
	Timeout          time.Duration
	HalfOpenRequests int
}

// DefaultConfig returns the default circuit breaker configuration.
func DefaultConfig(name string) Config {
	return Config{
		Name:             name,
		MaxFailures:      5,
		Timeout:          30 * time.Second,
		HalfOpenRequests: 3,
	}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	config           Config
	state            State
	failures         int
	successes        int
	lastFailure      time.Time
	halfOpenRequests int
	mu               sync.RWMutex
}

// New creates a new circuit breaker.
func New(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute executes the given function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if !cb.canExecute() {
		log.Warn().
			Str("circuit", cb.config.Name).
			Str("state", cb.state.String()).
			Msg("Circuit breaker rejected request")
		return ErrCircuitOpen
	}

	err := fn(ctx)

	cb.recordResult(err)

	return err
}

// canExecute checks if the circuit breaker allows execution.
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// Check if timeout has passed
		if time.Since(cb.lastFailure) > cb.config.Timeout {
			cb.toHalfOpen()
			return true
		}
		return false

	case StateHalfOpen:
		// Allow limited requests in half-open state
		if cb.halfOpenRequests < cb.config.HalfOpenRequests {
			cb.halfOpenRequests++
			return true
		}
		return false

	default:
		return false
	}
}

// recordResult records the result of an execution.
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

// onFailure handles a failure.
func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailure = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.MaxFailures {
			cb.toOpen()
		}

	case StateHalfOpen:
		cb.toOpen()
	}
}

// onSuccess handles a success.
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		cb.failures = 0

	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.HalfOpenRequests {
			cb.toClosed()
		}
	}
}

// toOpen transitions to open state.
func (cb *CircuitBreaker) toOpen() {
	cb.state = StateOpen
	cb.successes = 0
	cb.halfOpenRequests = 0
	log.Warn().
		Str("circuit", cb.config.Name).
		Int("failures", cb.failures).
		Msg("Circuit breaker opened")
}

// toHalfOpen transitions to half-open state.
func (cb *CircuitBreaker) toHalfOpen() {
	cb.state = StateHalfOpen
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
	log.Info().
		Str("circuit", cb.config.Name).
		Msg("Circuit breaker half-opened")
}

// toClosed transitions to closed state.
func (cb *CircuitBreaker) toClosed() {
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
	log.Info().
		Str("circuit", cb.config.Name).
		Msg("Circuit breaker closed")
}

// State returns the current state.
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Stats returns the current statistics.
func (cb *CircuitBreaker) Stats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"name":      cb.config.Name,
		"state":     cb.state.String(),
		"failures":  cb.failures,
		"successes": cb.successes,
	}
}
