// Package service provides application business logic and use cases.
package service

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/metrics"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/tracing"
)

// ErrAlertNotFound Alert service errors.
var (
	ErrAlertNotFound = errors.New("alert not found")
)

// AlertEventPublisher defines the interface for publishing alert events.
type AlertEventPublisher interface {
	PublishAlertCreated(alert *entity.Alert)
	PublishAlertAcknowledged(alert *entity.Alert)
	PublishAlertResolved(alert *entity.Alert)
	PublishAlertDeleted(alertID string)
}

// AlertEventProducer defines the interface for publishing alert events to the event bus.
type AlertEventProducer interface {
	PublishAlertCreated(ctx context.Context, alert *entity.Alert)
	PublishAlertAcknowledged(ctx context.Context, alert *entity.Alert)
	PublishAlertResolved(ctx context.Context, alert *entity.Alert)
	PublishAlertDeleted(ctx context.Context, alertID string, deletedBy string)
	PublishAlertExpired(ctx context.Context, alert *entity.Alert)
}

// AlertService handles alert business logic.
type AlertService struct {
	alertRepo     repository.AlertRepository
	cacheRepo     repository.CacheRepository
	wsPublisher   AlertEventPublisher
	eventProducer AlertEventProducer
}

// NewAlertService creates a new alert service.
func NewAlertService(
	alertRepo repository.AlertRepository,
	cacheRepo repository.CacheRepository,
	wsPublisher AlertEventPublisher,
) *AlertService {
	return &AlertService{
		alertRepo:   alertRepo,
		cacheRepo:   cacheRepo,
		wsPublisher: wsPublisher,
	}
}

// SetEventProducer sets the event producer for async event publishing.
func (s *AlertService) SetEventProducer(producer AlertEventProducer) {
	s.eventProducer = producer
}

// CreateAlertInput represents input for creating an alert.
type CreateAlertInput struct {
	Title    string
	Message  string
	Severity entity.AlertSeverity
	Source   string
	Metadata map[string]interface{}
}

// Create creates a new alert.
func (s *AlertService) Create(ctx context.Context, input CreateAlertInput) (*entity.Alert, error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("alert.title", input.Title),
		attribute.String("alert.severity", string(input.Severity)),
		attribute.String("alert.source", input.Source),
	)

	alert, err := entity.NewAlert(input.Title, input.Message, input.Severity, input.Source)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	for key, value := range input.Metadata {
		alert.AddMetadata(key, value)
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	span.SetAttributes(attribute.String("alert.id", alert.ID.String()))

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	// Record metrics
	metrics.AlertsCreatedTotal.WithLabelValues(string(input.Severity), input.Source).Inc()
	metrics.AlertsActiveGauge.Inc()

	// Publish to WebSocket (real-time)
	if s.wsPublisher != nil {
		s.wsPublisher.PublishAlertCreated(alert)
	}

	// Publish to Event Bus (async processing)
	if s.eventProducer != nil {
		s.eventProducer.PublishAlertCreated(ctx, alert)
	}

	tracing.AddEvent(ctx, "alert_created", attribute.String("alert.id", alert.ID.String()))

	return alert, nil
}

// GetByID retrieves an alert by ID.
func (s *AlertService) GetByID(ctx context.Context, id entity.ID) (*entity.Alert, error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.GetByID")
	defer span.End()

	span.SetAttributes(attribute.String("alert.id", id.String()))

	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			span.SetAttributes(attribute.Bool("alert.found", false))
			return nil, ErrAlertNotFound
		}
		tracing.RecordError(ctx, err)
		return nil, err
	}

	span.SetAttributes(attribute.Bool("alert.found", true))
	return alert, nil
}

// ListInput represents input for listing alerts.
type ListInput struct {
	Filter     valueobject.AlertFilter
	Pagination valueobject.Pagination
}

// List retrieves alerts with filters and pagination.
func (s *AlertService) List(ctx context.Context, input ListInput) (*valueobject.PaginatedResult[*entity.Alert], error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.List")
	defer span.End()

	span.SetAttributes(
		attribute.Int("pagination.page", input.Pagination.Page()),
		attribute.Int("pagination.page_size", input.Pagination.PageSize()),
	)

	if len(input.Filter.Statuses) > 0 {
		span.SetAttributes(attribute.Int("filter.statuses_count", len(input.Filter.Statuses)))
	}
	if len(input.Filter.Severities) > 0 {
		span.SetAttributes(attribute.Int("filter.severities_count", len(input.Filter.Severities)))
	}
	if input.Filter.Source != nil {
		span.SetAttributes(attribute.String("filter.source", *input.Filter.Source))
	}

	result, err := s.alertRepo.List(ctx, input.Filter, input.Pagination)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int64("result.total_items", result.TotalItems),
		attribute.Int("result.items_count", len(result.Items)),
	)

	return result, nil
}

// Acknowledge marks an alert as acknowledged.
func (s *AlertService) Acknowledge(ctx context.Context, alertID, userID entity.ID) (*entity.Alert, error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.Acknowledge")
	defer span.End()

	span.SetAttributes(
		attribute.String("alert.id", alertID.String()),
		attribute.String("user.id", userID.String()),
	)

	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		tracing.RecordError(ctx, err)
		return nil, err
	}

	if err := alert.Acknowledge(userID); err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	// Record metrics
	metrics.AlertsAcknowledgedTotal.Inc()

	// Publish to WebSocket (real-time)
	if s.wsPublisher != nil {
		s.wsPublisher.PublishAlertAcknowledged(alert)
	}

	// Publish to Event Bus (async processing)
	if s.eventProducer != nil {
		s.eventProducer.PublishAlertAcknowledged(ctx, alert)
	}

	tracing.AddEvent(ctx, "alert_acknowledged", attribute.String("alert.id", alert.ID.String()))

	return alert, nil
}

// Resolve marks an alert as resolved.
func (s *AlertService) Resolve(ctx context.Context, alertID, userID entity.ID) (*entity.Alert, error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.Resolve")
	defer span.End()

	span.SetAttributes(
		attribute.String("alert.id", alertID.String()),
		attribute.String("user.id", userID.String()),
	)

	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		tracing.RecordError(ctx, err)
		return nil, err
	}

	if err := alert.Resolve(userID); err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	// Record metrics
	metrics.AlertsResolvedTotal.Inc()
	metrics.AlertsActiveGauge.Dec()

	// Publish to WebSocket (real-time)
	if s.wsPublisher != nil {
		s.wsPublisher.PublishAlertResolved(alert)
	}

	// Publish to Event Bus (async processing)
	if s.eventProducer != nil {
		s.eventProducer.PublishAlertResolved(ctx, alert)
	}

	tracing.AddEvent(ctx, "alert_resolved", attribute.String("alert.id", alert.ID.String()))

	return alert, nil
}

// Delete removes an alert.
func (s *AlertService) Delete(ctx context.Context, id entity.ID, deletedBy entity.ID) error {
	ctx, span := tracing.StartSpan(ctx, "AlertService.Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("alert.id", id.String()),
		attribute.String("deleted_by", deletedBy.String()),
	)

	if err := s.alertRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrAlertNotFound
		}
		tracing.RecordError(ctx, err)
		return err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	// Record metrics
	metrics.AlertsDeletedTotal.Inc()

	// Publish to WebSocket (real-time)
	if s.wsPublisher != nil {
		s.wsPublisher.PublishAlertDeleted(id.String())
	}

	// Publish to Event Bus (async processing)
	if s.eventProducer != nil {
		s.eventProducer.PublishAlertDeleted(ctx, id.String(), deletedBy.String())
	}

	tracing.AddEvent(ctx, "alert_deleted", attribute.String("alert.id", id.String()))

	return nil
}

// GetStatistics retrieves alert statistics.
func (s *AlertService) GetStatistics(ctx context.Context) (*repository.AlertStatistics, error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.GetStatistics")
	defer span.End()

	var stats repository.AlertStatistics
	err := s.cacheRepo.Get(ctx, "stats:alerts", &stats)
	if err == nil {
		span.SetAttributes(attribute.Bool("cache.hit", true))
		return &stats, nil
	}

	span.SetAttributes(attribute.Bool("cache.hit", false))

	dbStats, err := s.alertRepo.GetStatistics(ctx)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	_ = s.cacheRepo.Set(ctx, "stats:alerts", dbStats, time.Minute)

	span.SetAttributes(attribute.Int64("stats.total_alerts", dbStats.TotalAlerts))

	return dbStats, nil
}

// GetActiveAlerts retrieves all active alerts.
func (s *AlertService) GetActiveAlerts(ctx context.Context) ([]*entity.Alert, error) {
	ctx, span := tracing.StartSpan(ctx, "AlertService.GetActiveAlerts")
	defer span.End()

	alerts, err := s.alertRepo.ListActive(ctx)
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("alerts.count", len(alerts)))

	return alerts, nil
}
