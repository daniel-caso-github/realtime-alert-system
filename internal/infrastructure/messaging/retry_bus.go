package messaging

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// RetryableBus wraps an event bus with retry logic.
type RetryableBus struct {
	bus     event.Bus
	retries *Retries
}

// NewRetryableBus creates a new retryable event bus.
func NewRetryableBus(bus event.Bus, config RetryConfig) *RetryableBus {
	return &RetryableBus{
		bus:     bus,
		retries: NewRetries(config),
	}
}

// Publish publishes an event with retry logic.
func (b *RetryableBus) Publish(ctx context.Context, evt *event.Event) error {
	return b.retries.Do(ctx, "publish_event", func(ctx context.Context) error {
		return b.bus.Publish(ctx, evt)
	})
}

// PublishToStream publishes an event to a specific stream with retry logic.
func (b *RetryableBus) PublishToStream(ctx context.Context, stream string, evt *event.Event) error {
	return b.retries.Do(ctx, "publish_to_stream", func(ctx context.Context) error {
		return b.bus.PublishToStream(ctx, stream, evt)
	})
}

// Subscribe subscribes to a stream (no retry needed as it maintains connection).
func (b *RetryableBus) Subscribe(ctx context.Context, stream string, group string, handler event.Handler) error {
	return b.bus.Subscribe(ctx, stream, group, handler)
}

// Unsubscribe unsubscribes from all streams.
func (b *RetryableBus) Unsubscribe() error {
	return b.bus.Unsubscribe()
}

// Compile-time interface verification.
var _ event.Bus = (*RetryableBus)(nil)
