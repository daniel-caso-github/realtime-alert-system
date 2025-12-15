// Package service provides application business logic and use cases.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// Alert service errors.
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

// AlertService handles alert business logic.
type AlertService struct {
	alertRepo repository.AlertRepository
	cacheRepo repository.CacheRepository
	publisher AlertEventPublisher
}

// NewAlertService creates a new alert service.
func NewAlertService(
	alertRepo repository.AlertRepository,
	cacheRepo repository.CacheRepository,
	publisher AlertEventPublisher,
) *AlertService {
	return &AlertService{
		alertRepo: alertRepo,
		cacheRepo: cacheRepo,
		publisher: publisher,
	}
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
	alert, err := entity.NewAlert(input.Title, input.Message, input.Severity, input.Source)
	if err != nil {
		return nil, err
	}

	for key, value := range input.Metadata {
		alert.AddMetadata(key, value)
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	if s.publisher != nil {
		s.publisher.PublishAlertCreated(alert)
	}

	return alert, nil
}

// GetByID retrieves an alert by ID.
func (s *AlertService) GetByID(ctx context.Context, id entity.ID) (*entity.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}
	return alert, nil
}

// ListInput represents input for listing alerts.
type ListInput struct {
	Filter     valueobject.AlertFilter
	Pagination valueobject.Pagination
}

// List retrieves alerts with filters and pagination.
func (s *AlertService) List(ctx context.Context, input ListInput) (*valueobject.PaginatedResult[*entity.Alert], error) {
	return s.alertRepo.List(ctx, input.Filter, input.Pagination)
}

// Acknowledge marks an alert as acknowledged.
func (s *AlertService) Acknowledge(ctx context.Context, alertID, userID entity.ID) (*entity.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}

	if err := alert.Acknowledge(userID); err != nil {
		return nil, err
	}

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	if s.publisher != nil {
		s.publisher.PublishAlertAcknowledged(alert)
	}

	return alert, nil
}

// Resolve marks an alert as resolved.
func (s *AlertService) Resolve(ctx context.Context, alertID, userID entity.ID) (*entity.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}

	if err := alert.Resolve(userID); err != nil {
		return nil, err
	}

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	if s.publisher != nil {
		s.publisher.PublishAlertResolved(alert)
	}

	return alert, nil
}

// Delete removes an alert.
func (s *AlertService) Delete(ctx context.Context, id entity.ID) error {
	if err := s.alertRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrAlertNotFound
		}
		return err
	}

	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	if s.publisher != nil {
		s.publisher.PublishAlertDeleted(id.String())
	}

	return nil
}

// GetStatistics retrieves alert statistics.
func (s *AlertService) GetStatistics(ctx context.Context) (*repository.AlertStatistics, error) {
	var stats repository.AlertStatistics
	err := s.cacheRepo.Get(ctx, "stats:alerts", &stats)
	if err == nil {
		return &stats, nil
	}

	dbStats, err := s.alertRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.cacheRepo.Set(ctx, "stats:alerts", dbStats, time.Minute)

	return dbStats, nil
}

// GetActiveAlerts retrieves all active alerts.
func (s *AlertService) GetActiveAlerts(ctx context.Context) ([]*entity.Alert, error) {
	return s.alertRepo.ListActive(ctx)
}
