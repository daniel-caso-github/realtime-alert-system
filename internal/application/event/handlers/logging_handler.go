// Package handlers provides event handler implementations.
package handlers

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// LoggingHandler logs all alert events for auditing.
type LoggingHandler struct{}

// NewLoggingHandler creates a new logging handler.
func NewLoggingHandler() *LoggingHandler {
	return &LoggingHandler{}
}

// HandleAlertCreated logs alert created events.
func (h *LoggingHandler) HandleAlertCreated(_ context.Context, payload event.AlertPayload) error {
	log.Info().
		Str("alert_id", payload.ID).
		Str("title", payload.Title).
		Str("severity", payload.Severity).
		Str("source", payload.Source).
		Msg("Alert created event processed")
	return nil
}

// HandleAlertAcknowledged logs alert acknowledged events.
func (h *LoggingHandler) HandleAlertAcknowledged(_ context.Context, payload event.AlertPayload) error {
	acknowledgedBy := ""
	if payload.AcknowledgedBy != nil {
		acknowledgedBy = *payload.AcknowledgedBy
	}

	log.Info().
		Str("alert_id", payload.ID).
		Str("title", payload.Title).
		Str("acknowledged_by", acknowledgedBy).
		Msg("Alert acknowledged event processed")
	return nil
}

// HandleAlertResolved logs alert resolved events.
func (h *LoggingHandler) HandleAlertResolved(_ context.Context, payload event.AlertPayload) error {
	resolvedBy := ""
	if payload.ResolvedBy != nil {
		resolvedBy = *payload.ResolvedBy
	}

	log.Info().
		Str("alert_id", payload.ID).
		Str("title", payload.Title).
		Str("resolved_by", resolvedBy).
		Msg("Alert resolved event processed")
	return nil
}

// HandleAlertDeleted logs alert deleted events.
func (h *LoggingHandler) HandleAlertDeleted(_ context.Context, payload event.AlertDeletedPayload) error {
	log.Info().
		Str("alert_id", payload.ID).
		Str("deleted_by", payload.DeletedBy).
		Time("deleted_at", payload.DeletedAt).
		Msg("Alert deleted event processed")
	return nil
}

// HandleAlertExpired logs alert expired events.
func (h *LoggingHandler) HandleAlertExpired(_ context.Context, payload event.AlertPayload) error {
	log.Info().
		Str("alert_id", payload.ID).
		Str("title", payload.Title).
		Msg("Alert expired event processed")
	return nil
}
