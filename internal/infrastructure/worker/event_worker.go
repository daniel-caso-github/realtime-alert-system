// Package worker provides background workers for async processing.
package worker

import (
	"context"

	"github.com/rs/zerolog/log"

	appevent "github.com/daniel-caso-github/realtime-alerting-system/internal/application/event"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/event/handlers"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
)

// EventWorker manages event consumers and handlers.
type EventWorker struct {
	bus            event.Bus
	alertConsumer  *appevent.AlertConsumer
	metricsHandler *handlers.MetricsHandler
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewEventWorker creates a new event worker.
func NewEventWorker(bus event.Bus) *EventWorker {
	ctx, cancel := context.WithCancel(context.Background())

	return &EventWorker{
		bus:    bus,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start starts the event worker and all consumers.
func (w *EventWorker) Start() error {
	log.Info().Msg("Starting event worker...")

	// Create consumers
	w.alertConsumer = appevent.NewAlertConsumer()

	// Create and register handlers
	loggingHandler := handlers.NewLoggingHandler()
	w.metricsHandler = handlers.NewMetricsHandler()

	w.alertConsumer.RegisterHandler(loggingHandler)
	w.alertConsumer.RegisterHandler(w.metricsHandler)

	// Subscribe to streams
	if err := w.bus.Subscribe(w.ctx, event.StreamAlerts, event.GroupAlertProcessors, w.alertConsumer.Handle); err != nil {
		return err
	}

	log.Info().Msg("Event worker started successfully")
	return nil
}

// Stop stops the event worker.
func (w *EventWorker) Stop() error {
	log.Info().Msg("Stopping event worker...")
	w.cancel()

	if err := w.bus.Unsubscribe(); err != nil {
		log.Error().Err(err).Msg("Error unsubscribing from event bus")
		return err
	}

	log.Info().Msg("Event worker stopped")
	return nil
}

// GetMetrics returns the current event metrics.
func (w *EventWorker) GetMetrics() map[string]int64 {
	if w.metricsHandler == nil {
		return make(map[string]int64)
	}
	return w.metricsHandler.GetMetrics()
}
