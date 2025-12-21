package event

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// AlertConsumer consumes and processes alert events.
type AlertConsumer struct {
	handlers []AlertEventHandler
}

// NewAlertConsumer creates a new alert consumer.
func NewAlertConsumer() *AlertConsumer {
	return &AlertConsumer{
		handlers: make([]AlertEventHandler, 0),
	}
}

// RegisterHandler registers an event handler.
func (c *AlertConsumer) RegisterHandler(handler AlertEventHandler) {
	c.handlers = append(c.handlers, handler)
}

// Handle processes an event from the event bus.
func (c *AlertConsumer) Handle(ctx context.Context, evt *event.Event) error {
	log.Debug().
		Str("event_id", evt.ID).
		Str("event_type", string(evt.Type)).
		Int("retries", evt.Retries).
		Msg("Processing event")

	switch evt.Type {
	case event.AlertCreated:
		return c.handleAlertCreated(ctx, evt)
	case event.AlertAcknowledged:
		return c.handleAlertAcknowledged(ctx, evt)
	case event.AlertResolved:
		return c.handleAlertResolved(ctx, evt)
	case event.AlertDeleted:
		return c.handleAlertDeleted(ctx, evt)
	case event.AlertExpired:
		return c.handleAlertExpired(ctx, evt)
	default:
		log.Warn().Str("event_type", string(evt.Type)).Msg("Unknown event type")
		return nil
	}
}

func (c *AlertConsumer) handleAlertCreated(ctx context.Context, evt *event.Event) error {
	var payload event.AlertPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal alert created payload")
		return err
	}

	for _, handler := range c.handlers {
		if err := handler.HandleAlertCreated(ctx, payload); err != nil {
			log.Error().Err(err).Str("alert_id", payload.ID).Msg("Handler failed for alert.created")
			return err
		}
	}

	return nil
}

func (c *AlertConsumer) handleAlertAcknowledged(ctx context.Context, evt *event.Event) error {
	var payload event.AlertPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal alert acknowledged payload")
		return err
	}

	for _, handler := range c.handlers {
		if err := handler.HandleAlertAcknowledged(ctx, payload); err != nil {
			log.Error().Err(err).Str("alert_id", payload.ID).Msg("Handler failed for alert.acknowledged")
			return err
		}
	}

	return nil
}

func (c *AlertConsumer) handleAlertResolved(ctx context.Context, evt *event.Event) error {
	var payload event.AlertPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal alert resolved payload")
		return err
	}

	for _, handler := range c.handlers {
		if err := handler.HandleAlertResolved(ctx, payload); err != nil {
			log.Error().Err(err).Str("alert_id", payload.ID).Msg("Handler failed for alert.resolved")
			return err
		}
	}

	return nil
}

func (c *AlertConsumer) handleAlertDeleted(ctx context.Context, evt *event.Event) error {
	var payload event.AlertDeletedPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal alert deleted payload")
		return err
	}

	for _, handler := range c.handlers {
		if err := handler.HandleAlertDeleted(ctx, payload); err != nil {
			log.Error().Err(err).Str("alert_id", payload.ID).Msg("Handler failed for alert.deleted")
			return err
		}
	}

	return nil
}

func (c *AlertConsumer) handleAlertExpired(ctx context.Context, evt *event.Event) error {
	var payload event.AlertPayload
	if err := evt.UnmarshalPayload(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal alert expired payload")
		return err
	}

	for _, handler := range c.handlers {
		if err := handler.HandleAlertExpired(ctx, payload); err != nil {
			log.Error().Err(err).Str("alert_id", payload.ID).Msg("Handler failed for alert.expired")
			return err
		}
	}

	return nil
}
