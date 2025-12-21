package handlers

import (
	"context"
	"sync/atomic"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// MetricsHandler collects metrics from alert events.
type MetricsHandler struct {
	alertsCreated      int64
	alertsAcknowledged int64
	alertsResolved     int64
	alertsDeleted      int64
	alertsExpired      int64
}

// NewMetricsHandler creates a new metrics handler.
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// HandleAlertCreated increments the alerts created counter.
func (h *MetricsHandler) HandleAlertCreated(_ context.Context, _ event.AlertPayload) error {
	atomic.AddInt64(&h.alertsCreated, 1)
	return nil
}

// HandleAlertAcknowledged increments the alerts acknowledged counter.
func (h *MetricsHandler) HandleAlertAcknowledged(_ context.Context, _ event.AlertPayload) error {
	atomic.AddInt64(&h.alertsAcknowledged, 1)
	return nil
}

// HandleAlertResolved increments the alerts resolved counter.
func (h *MetricsHandler) HandleAlertResolved(_ context.Context, _ event.AlertPayload) error {
	atomic.AddInt64(&h.alertsResolved, 1)
	return nil
}

// HandleAlertDeleted increments the alerts deleted counter.
func (h *MetricsHandler) HandleAlertDeleted(_ context.Context, _ event.AlertDeletedPayload) error {
	atomic.AddInt64(&h.alertsDeleted, 1)
	return nil
}

// HandleAlertExpired increments the alerts expired counter.
func (h *MetricsHandler) HandleAlertExpired(_ context.Context, _ event.AlertPayload) error {
	atomic.AddInt64(&h.alertsExpired, 1)
	return nil
}

// GetMetrics returns the current metrics.
func (h *MetricsHandler) GetMetrics() map[string]int64 {
	return map[string]int64{
		"alerts_created":      atomic.LoadInt64(&h.alertsCreated),
		"alerts_acknowledged": atomic.LoadInt64(&h.alertsAcknowledged),
		"alerts_resolved":     atomic.LoadInt64(&h.alertsResolved),
		"alerts_deleted":      atomic.LoadInt64(&h.alertsDeleted),
		"alerts_expired":      atomic.LoadInt64(&h.alertsExpired),
	}
}
