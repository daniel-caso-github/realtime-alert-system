// Package service implements the application layer services following hexagonal architecture.
// Services orchestrate domain logic and coordinate between repositories and other infrastructure.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// Alert service errors define domain-specific error types for the alert service.
var (
	ErrAlertNotFound = errors.New("alert not found") // Returned when an alert cannot be found by ID
)

// AlertService handles alert business logic and orchestrates operations between
// the alert repository and cache. It implements the application use cases for
// alert management including creation, retrieval, acknowledgment, and resolution.
type AlertService struct {
	alertRepo repository.AlertRepository // Primary storage for alerts
	cacheRepo repository.CacheRepository // Cache for statistics and frequently accessed data
}

// NewAlertService creates a new AlertService with the required dependencies.
// Both repositories are required for proper operation.
func NewAlertService(
	alertRepo repository.AlertRepository, // Repository for alert persistence
	cacheRepo repository.CacheRepository, // Repository for caching operations
) *AlertService {
	return &AlertService{
		alertRepo: alertRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateAlertInput represents the input parameters for creating a new alert.
// This struct is used to pass validated data from the handler layer to the service.
type CreateAlertInput struct {
	Title    string                 // Alert title (required)
	Message  string                 // Alert message body (required)
	Severity entity.AlertSeverity   // Severity level (critical, high, medium, low, info)
	Source   string                 // Source system that generated the alert
	Metadata map[string]interface{} // Additional key-value metadata
}

// Create creates a new alert and persists it to the database.
// It creates a domain entity, attaches any metadata, saves to the repository,
// and invalidates the statistics cache to ensure consistency.
// Returns the created alert entity or an error if validation or persistence fails.
func (s *AlertService) Create(ctx context.Context, input CreateAlertInput) (*entity.Alert, error) {
	// Create the domain entity with validation
	alert, err := entity.NewAlert(input.Title, input.Message, input.Severity, input.Source)
	if err != nil {
		return nil, err
	}

	// Attach optional metadata to the alert
	for key, value := range input.Metadata {
		alert.AddMetadata(key, value)
	}

	// Persist to the database
	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, err
	}

	// Invalidate statistics cache to reflect the new alert
	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	return alert, nil
}

// GetByID retrieves a single alert by its unique identifier.
// Returns ErrAlertNotFound if no alert exists with the given ID.
func (s *AlertService) GetByID(ctx context.Context, id entity.ID) (*entity.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		// Convert repository-level not found error to service-level error
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}
	return alert, nil
}

// ListInput represents the input parameters for listing alerts with filters.
// It combines filtering criteria with pagination settings.
type ListInput struct {
	Filter     valueobject.AlertFilter // Criteria for filtering alerts (status, severity, date range, etc.)
	Pagination valueobject.Pagination  // Pagination settings (page, page size)
}

// List retrieves alerts matching the specified filters with pagination.
// Returns a paginated result containing the alerts and pagination metadata.
func (s *AlertService) List(ctx context.Context, input ListInput) (*valueobject.PaginatedResult[*entity.Alert], error) {
	return s.alertRepo.List(ctx, input.Filter, input.Pagination)
}

// Acknowledge marks an alert as acknowledged by the specified user.
// This indicates that a user has seen and is aware of the alert.
// The operation is atomic: it retrieves the alert, updates its status, and persists the change.
// Returns ErrAlertNotFound if the alert doesn't exist, or a domain error if
// the alert cannot be acknowledged (e.g., already resolved).
func (s *AlertService) Acknowledge(ctx context.Context, alertID, userID entity.ID) (*entity.Alert, error) {
	// Retrieve the alert from storage
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}

	// Apply the domain operation to mark as acknowledged
	if err := alert.Acknowledge(userID); err != nil {
		return nil, err
	}

	// Persist the updated alert
	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	// Invalidate statistics cache to reflect the status change
	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	return alert, nil
}

// Resolve marks an alert as resolved by the specified user.
// This indicates that the issue that triggered the alert has been addressed.
// The operation is atomic: it retrieves the alert, updates its status, and persists the change.
// Returns ErrAlertNotFound if the alert doesn't exist, or a domain error if
// the alert cannot be resolved (e.g., already resolved).
func (s *AlertService) Resolve(ctx context.Context, alertID, userID entity.ID) (*entity.Alert, error) {
	// Retrieve the alert from storage
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlertNotFound
		}
		return nil, err
	}

	// Apply the domain operation to mark as resolved
	if err := alert.Resolve(userID); err != nil {
		return nil, err
	}

	// Persist the updated alert
	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, err
	}

	// Invalidate statistics cache to reflect the status change
	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	return alert, nil
}

// Delete permanently removes an alert from the system.
// Returns ErrAlertNotFound if no alert exists with the given ID.
// Use with caution: this operation cannot be undone.
func (s *AlertService) Delete(ctx context.Context, id entity.ID) error {
	if err := s.alertRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrAlertNotFound
		}
		return err
	}

	// Invalidate statistics cache to reflect the deletion
	_ = s.cacheRepo.Delete(ctx, "stats:alerts")

	return nil
}

// GetStatistics retrieves aggregated alert statistics for dashboards.
// It implements a cache-aside pattern: first checking the cache for recent statistics,
// and falling back to the database if not cached. Results are cached for 1 minute
// to balance freshness with database load.
func (s *AlertService) GetStatistics(ctx context.Context) (*repository.AlertStatistics, error) {
	// Try to retrieve from cache first (cache-aside pattern)
	var stats repository.AlertStatistics
	err := s.cacheRepo.Get(ctx, "stats:alerts", &stats)
	if err == nil {
		return &stats, nil // Cache hit
	}

	// Cache miss: fetch from database
	dbStats, err := s.alertRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache with 1 minute TTL for subsequent requests
	_ = s.cacheRepo.Set(ctx, "stats:alerts", dbStats, time.Minute)

	return dbStats, nil
}

// GetActiveAlerts retrieves all alerts with status "active".
// This is useful for real-time dashboards and notification systems
// that need to display current unresolved alerts.
func (s *AlertService) GetActiveAlerts(ctx context.Context) ([]*entity.Alert, error) {
	return s.alertRepo.ListActive(ctx)
}
