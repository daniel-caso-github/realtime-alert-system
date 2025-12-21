package service

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/notification"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
)

// NotificationService manages notifications across multiple channels.
type NotificationService struct {
	notifiers   []notification.Notifier
	minSeverity string
	rateLimit   int
	mu          sync.Mutex
	sentCount   map[string]int
	lastReset   time.Time
}

// NewNotificationService creates a new notification service.
func NewNotificationService(cfg config.NotificationConfig, notifiers ...notification.Notifier) *NotificationService {
	activeNotifiers := make([]notification.Notifier, 0)
	for _, n := range notifiers {
		if n.IsEnabled() {
			activeNotifiers = append(activeNotifiers, n)
			log.Info().Str("notifier", n.Name()).Msg("Notification channel enabled")
		}
	}

	return &NotificationService{
		notifiers:   activeNotifiers,
		minSeverity: cfg.MinSeverity,
		rateLimit:   cfg.RateLimitPerMinute,
		sentCount:   make(map[string]int),
		lastReset:   time.Now(),
	}
}

// Notify sends a notification through all enabled channels.
func (s *NotificationService) Notify(ctx context.Context, msg notification.Message) error {
	// Check severity threshold
	if !notification.ShouldNotify(msg.Severity, s.minSeverity) {
		log.Debug().
			Str("severity", msg.Severity).
			Str("min_severity", s.minSeverity).
			Msg("Notification skipped due to severity threshold")
		return nil
	}

	// Check rate limit
	if !s.checkRateLimit(msg.AlertID) {
		log.Warn().
			Str("alert_id", msg.AlertID).
			Msg("Notification rate limited")
		return nil
	}

	// Send to all notifiers
	var lastErr error
	for _, notifier := range s.notifiers {
		if err := notifier.Send(ctx, msg); err != nil {
			log.Error().
				Err(err).
				Str("notifier", notifier.Name()).
				Str("alert_id", msg.AlertID).
				Msg("Failed to send notification")
			lastErr = err
		}
	}

	return lastErr
}

// checkRateLimit checks if we can send a notification (rate limiting).
func (s *NotificationService) checkRateLimit(alertID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Reset counter every minute
	if time.Since(s.lastReset) > time.Minute {
		s.sentCount = make(map[string]int)
		s.lastReset = time.Now()
	}

	// Check global rate limit
	total := 0
	for _, count := range s.sentCount {
		total += count
	}
	if total >= s.rateLimit {
		return false
	}

	s.sentCount[alertID]++
	return true
}

// GetActiveNotifiers returns the list of active notifier names.
func (s *NotificationService) GetActiveNotifiers() []string {
	names := make([]string, len(s.notifiers))
	for i, n := range s.notifiers {
		names[i] = n.Name()
	}
	return names
}
