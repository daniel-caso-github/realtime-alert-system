// Package messaging provides event bus implementations.
package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// RedisStreamBus implements event.Bus using Redis Streams.
type RedisStreamBus struct {
	client     *redis.Client
	handlers   map[string]event.Handler
	mu         sync.RWMutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
	consumerID string
}

// NewRedisStreamBus creates a new Redis Streams event bus.
func NewRedisStreamBus(client *redis.Client, consumerID string) *RedisStreamBus {
	return &RedisStreamBus{
		client:     client,
		handlers:   make(map[string]event.Handler),
		stopCh:     make(chan struct{}),
		consumerID: consumerID,
	}
}

// Publish publishes an event to the default stream based on event type.
func (b *RedisStreamBus) Publish(ctx context.Context, evt *event.Event) error {
	stream := b.getStreamForEventType(evt.Type)
	return b.PublishToStream(ctx, stream, evt)
}

// PublishToStream publishes an event to a specific stream.
func (b *RedisStreamBus) PublishToStream(ctx context.Context, stream string, evt *event.Event) error {
	args := &redis.XAddArgs{
		Stream: stream,
		Values: evt.ToMap(),
	}

	_, err := b.client.XAdd(ctx, args).Result()
	if err != nil {
		log.Error().Err(err).Str("stream", stream).Str("event_type", string(evt.Type)).Msg("Failed to publish event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Debug().Str("stream", stream).Str("event_id", evt.ID).Str("event_type", string(evt.Type)).Msg("Event published")
	return nil
}

// Subscribe subscribes to a stream with a consumer group.
func (b *RedisStreamBus) Subscribe(ctx context.Context, stream string, group string, handler event.Handler) error {
	// Create consumer group if it doesn't exist
	err := b.client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	b.mu.Lock()
	key := fmt.Sprintf("%s:%s", stream, group)
	b.handlers[key] = handler
	b.mu.Unlock()

	b.wg.Add(1)
	go b.consume(ctx, stream, group, handler)

	log.Info().Str("stream", stream).Str("group", group).Str("consumer", b.consumerID).Msg("Subscribed to stream")
	return nil
}

// consume reads messages from the stream.
func (b *RedisStreamBus) consume(ctx context.Context, stream string, group string, handler event.Handler) {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopCh:
			return
		case <-ctx.Done():
			return
		default:
			b.readMessages(ctx, stream, group, handler)
		}
	}
}

// readMessages reads and processes messages from the stream.
func (b *RedisStreamBus) readMessages(ctx context.Context, stream string, group string, handler event.Handler) {
	streams, err := b.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: b.consumerID,
		Streams:  []string{stream, ">"},
		Count:    10,
		Block:    time.Second * 5,
	}).Result()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Error().Err(err).Str("stream", stream).Msg("Error reading from stream")
		}
		return
	}

	for _, s := range streams {
		for _, msg := range s.Messages {
			b.processMessage(ctx, stream, group, msg, handler)
		}
	}
}

// processMessage processes a single message.
func (b *RedisStreamBus) processMessage(ctx context.Context, stream string, group string, msg redis.XMessage, handler event.Handler) {
	evt, err := event.FromMap(msg.Values)
	if err != nil {
		log.Error().Err(err).Str("message_id", msg.ID).Msg("Failed to parse event")
		b.acknowledgeMessage(ctx, stream, group, msg.ID)
		return
	}

	if err := handler(ctx, evt); err != nil {
		log.Error().Err(err).Str("event_id", evt.ID).Str("event_type", string(evt.Type)).Msg("Failed to handle event")
		b.handleFailedEvent(ctx, evt, err)
	}

	b.acknowledgeMessage(ctx, stream, group, msg.ID)
}

// acknowledgeMessage acknowledges a message.
func (b *RedisStreamBus) acknowledgeMessage(ctx context.Context, stream string, group string, messageID string) {
	if err := b.client.XAck(ctx, stream, group, messageID).Err(); err != nil {
		log.Error().Err(err).Str("message_id", messageID).Msg("Failed to acknowledge message")
	}
}

// handleFailedEvent moves failed events to the dead letter queue.
func (b *RedisStreamBus) handleFailedEvent(ctx context.Context, evt *event.Event, _ error) {
	evt.Retries++

	if evt.Retries >= 3 {
		// Move to dead letter queue
		if err := b.PublishToStream(ctx, event.StreamDeadLetter, evt); err != nil {
			log.Error().Err(err).Str("event_id", evt.ID).Msg("Failed to move event to dead letter queue")
		}
		log.Warn().Str("event_id", evt.ID).Int("retries", evt.Retries).Msg("Event moved to dead letter queue")
		return
	}

	// Re-publish for retry
	stream := b.getStreamForEventType(evt.Type)
	if err := b.PublishToStream(ctx, stream, evt); err != nil {
		log.Error().Err(err).Str("event_id", evt.ID).Msg("Failed to re-publish event for retry")
	}
	log.Debug().Str("event_id", evt.ID).Int("retries", evt.Retries).Msg("Event re-published for retry")
}

// Unsubscribe stops all consumers.
func (b *RedisStreamBus) Unsubscribe() error {
	close(b.stopCh)
	b.wg.Wait()
	return nil
}

// getStreamForEventType returns the stream name for an event type.
func (b *RedisStreamBus) getStreamForEventType(eventType event.Type) string {
	switch eventType {
	case event.AlertCreated, event.AlertAcknowledged, event.AlertResolved, event.AlertDeleted, event.AlertExpired:
		return event.StreamAlerts
	case event.UserCreated, event.UserUpdated:
		return event.StreamNotifications
	default:
		return event.StreamAlerts
	}
}

// Compile-time interface verification.
var _ event.Bus = (*RedisStreamBus)(nil)
