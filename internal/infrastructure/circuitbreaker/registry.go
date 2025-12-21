package circuitbreaker

import "sync"

// Registry manages multiple circuit breakers.
type Registry struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewRegistry creates a new circuit breaker registry.
func NewRegistry() *Registry {
	return &Registry{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Get gets or creates a circuit breaker by name.
func (r *Registry) Get(name string) *CircuitBreaker {
	r.mu.RLock()
	cb, exists := r.breakers[name]
	r.mu.RUnlock()

	if exists {
		return cb
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if cb, exists = r.breakers[name]; exists {
		return cb
	}

	cb = New(DefaultConfig(name))
	r.breakers[name] = cb
	return cb
}

// GetWithConfig gets or creates a circuit breaker with custom config.
func (r *Registry) GetWithConfig(config Config) *CircuitBreaker {
	r.mu.RLock()
	cb, exists := r.breakers[config.Name]
	r.mu.RUnlock()

	if exists {
		return cb
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if cb, exists = r.breakers[config.Name]; exists {
		return cb
	}

	cb = New(config)
	r.breakers[config.Name] = cb
	return cb
}

// Stats returns statistics for all circuit breakers.
func (r *Registry) Stats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make(map[string]interface{})
	for name, cb := range r.breakers {
		stats[name] = cb.Stats()
	}
	return stats
}
