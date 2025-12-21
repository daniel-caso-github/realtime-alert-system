// Package event provides event producers for the application layer.
package event

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// AlertProducer publishes alert-related events.
type AlertProducer struct {
	bus event.Publisher
}

// NewAlertProducer creates a new alert event producer.
func NewAlertProducer(bus event.Publisher) *AlertProducer {
	return &AlertProducer{
		bus: bus,
	}
}

// PublishAlertCreated publishes an alert created event.
func (p *AlertProducer) PublishAlertCreated(ctx context.Context, alert *entity.Alert) {
	payload := p.alertToPayload(alert)

	evt, err := event.NewEvent(event.AlertCreated, payload)
	if err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to create alert.created event")
		return
	}

	if err := p.bus.Publish(ctx, evt); err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to publish alert.created event")
	}
}

// PublishAlertAcknowledged publishes an alert acknowledged event.
func (p *AlertProducer) PublishAlertAcknowledged(ctx context.Context, alert *entity.Alert) {
	payload := p.alertToPayload(alert)

	evt, err := event.NewEvent(event.AlertAcknowledged, payload)
	if err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to create alert.acknowledged event")
		return
	}

	if err := p.bus.Publish(ctx, evt); err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to publish alert.acknowledged event")
	}
}

// PublishAlertResolved publishes an alert resolved event.
func (p *AlertProducer) PublishAlertResolved(ctx context.Context, alert *entity.Alert) {
	payload := p.alertToPayload(alert)

	evt, err := event.NewEvent(event.AlertResolved, payload)
	if err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to create alert.resolved event")
		return
	}

	if err := p.bus.Publish(ctx, evt); err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to publish alert.resolved event")
	}
}

// PublishAlertDeleted publishes an alert deleted event.
func (p *AlertProducer) PublishAlertDeleted(ctx context.Context, alertID string, deletedBy string) {
	payload := event.AlertDeletedPayload{
		ID:        alertID,
		DeletedAt: time.Now().UTC(),
		DeletedBy: deletedBy,
	}

	evt, err := event.NewEvent(event.AlertDeleted, payload)
	if err != nil {
		log.Error().Err(err).Str("alert_id", alertID).Msg("Failed to create alert.deleted event")
		return
	}

	if err := p.bus.Publish(ctx, evt); err != nil {
		log.Error().Err(err).Str("alert_id", alertID).Msg("Failed to publish alert.deleted event")
	}
}

// PublishAlertExpired publishes an alert expired event.
func (p *AlertProducer) PublishAlertExpired(ctx context.Context, alert *entity.Alert) {
	payload := p.alertToPayload(alert)

	evt, err := event.NewEvent(event.AlertExpired, payload)
	if err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to create alert.expired event")
		return
	}

	if err := p.bus.Publish(ctx, evt); err != nil {
		log.Error().Err(err).Str("alert_id", alert.ID.String()).Msg("Failed to publish alert.expired event")
	}
}

// alertToPayload converts an alert entity to an event payload.
func (p *AlertProducer) alertToPayload(alert *entity.Alert) event.AlertPayload {
	payload := event.AlertPayload{
		ID:        alert.ID.String(),
		Title:     alert.Title,
		Message:   alert.Message,
		Severity:  string(alert.Severity),
		Status:    string(alert.Status),
		Source:    alert.Source,
		Metadata:  alert.Metadata,
		CreatedAt: alert.CreatedAt,
	}

	if alert.AcknowledgedBy != nil {
		ackBy := alert.AcknowledgedBy.String()
		payload.AcknowledgedBy = &ackBy
	}
	if alert.AcknowledgedAt != nil {
		payload.AcknowledgedAt = alert.AcknowledgedAt
	}
	if alert.ResolvedBy != nil {
		resBy := alert.ResolvedBy.String()
		payload.ResolvedBy = &resBy
	}
	if alert.ResolvedAt != nil {
		payload.ResolvedAt = alert.ResolvedAt
	}

	return payload
}
