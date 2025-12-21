package notification

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/notification"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/circuitbreaker"
)

// ResilientNotifier wraps a notifier with circuit breaker protection.
type ResilientNotifier struct {
	notifier notification.Notifier
	cb       *circuitbreaker.CircuitBreaker
}

// NewResilientNotifier creates a new resilient notifier.
func NewResilientNotifier(notifier notification.Notifier, cb *circuitbreaker.CircuitBreaker) *ResilientNotifier {
	return &ResilientNotifier{
		notifier: notifier,
		cb:       cb,
	}
}

// Send sends a notification with circuit breaker protection.
func (n *ResilientNotifier) Send(ctx context.Context, msg notification.Message) error {
	return n.cb.Execute(ctx, func(ctx context.Context) error {
		return n.notifier.Send(ctx, msg)
	})
}

// Name returns the notifier name.
func (n *ResilientNotifier) Name() string {
	return n.notifier.Name()
}

// IsEnabled returns whether the notifier is enabled.
func (n *ResilientNotifier) IsEnabled() bool {
	return n.notifier.IsEnabled()
}

// Stats returns circuit breaker statistics.
func (n *ResilientNotifier) Stats() map[string]interface{} {
	return n.cb.Stats()
}

// Compile-time interface verification.
var _ notification.Notifier = (*ResilientNotifier)(nil)
