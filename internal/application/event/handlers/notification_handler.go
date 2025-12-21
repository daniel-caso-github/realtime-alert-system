package handlers

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/service"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/notification"
)

// NotificationHandler sends notifications for alert events.
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new notification handler.
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// HandleAlertCreated sends notification for new alerts.
func (h *NotificationHandler) HandleAlertCreated(ctx context.Context, payload event.AlertPayload) error {
	msg := notification.Message{
		Title:    "🚨 New Alert: " + payload.Title,
		Text:     payload.Message,
		Severity: payload.Severity,
		AlertID:  payload.ID,
		Source:   payload.Source,
		Fields:   make(map[string]string),
	}

	return h.notificationService.Notify(ctx, msg)
}

// HandleAlertAcknowledged sends notification when alert is acknowledged.
func (h *NotificationHandler) HandleAlertAcknowledged(ctx context.Context, payload event.AlertPayload) error {
	acknowledgedBy := "unknown"
	if payload.AcknowledgedBy != nil {
		acknowledgedBy = *payload.AcknowledgedBy
	}

	msg := notification.Message{
		Title:    "✅ Alert Acknowledged: " + payload.Title,
		Text:     "Alert has been acknowledged",
		Severity: payload.Severity,
		AlertID:  payload.ID,
		Source:   payload.Source,
		Fields: map[string]string{
			"Acknowledged By": acknowledgedBy,
		},
	}

	return h.notificationService.Notify(ctx, msg)
}

// HandleAlertResolved sends notification when alert is resolved.
func (h *NotificationHandler) HandleAlertResolved(ctx context.Context, payload event.AlertPayload) error {
	resolvedBy := "unknown"
	if payload.ResolvedBy != nil {
		resolvedBy = *payload.ResolvedBy
	}

	msg := notification.Message{
		Title:    "✔️ Alert Resolved: " + payload.Title,
		Text:     "Alert has been resolved",
		Severity: payload.Severity,
		AlertID:  payload.ID,
		Source:   payload.Source,
		Fields: map[string]string{
			"Resolved By": resolvedBy,
		},
	}

	return h.notificationService.Notify(ctx, msg)
}

// HandleAlertDeleted does not send notification (optional).
func (h *NotificationHandler) HandleAlertDeleted(_ context.Context, _ event.AlertDeletedPayload) error {
	// No notification for deleted alerts
	return nil
}

// HandleAlertExpired sends notification when alert expires.
func (h *NotificationHandler) HandleAlertExpired(ctx context.Context, payload event.AlertPayload) error {
	msg := notification.Message{
		Title:    "⏰ Alert Expired: " + payload.Title,
		Text:     "Alert has expired without resolution",
		Severity: payload.Severity,
		AlertID:  payload.ID,
		Source:   payload.Source,
	}

	return h.notificationService.Notify(ctx, msg)
}
