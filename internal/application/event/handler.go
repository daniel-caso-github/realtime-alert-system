package event

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// AlertEventHandler handles alert events from the event bus.
type AlertEventHandler interface {
	HandleAlertCreated(ctx context.Context, payload event.AlertPayload) error
	HandleAlertAcknowledged(ctx context.Context, payload event.AlertPayload) error
	HandleAlertResolved(ctx context.Context, payload event.AlertPayload) error
	HandleAlertDeleted(ctx context.Context, payload event.AlertDeletedPayload) error
	HandleAlertExpired(ctx context.Context, payload event.AlertPayload) error
}
