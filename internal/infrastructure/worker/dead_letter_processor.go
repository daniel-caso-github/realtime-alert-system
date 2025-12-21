package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/event"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
)

// FailedEvent represents a failed event stored for analysis.
type FailedEvent struct {
	ID          string          `json:"id"`
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	Payload     json.RawMessage `json:"payload"`
	Retries     int             `json:"retries"`
	LastError   string          `json:"last_error,omitempty"`
	FailedAt    time.Time       `json:"failed_at"`
	ProcessedAt *time.Time      `json:"processed_at,omitempty"`
	Status      string          `json:"status"` // pending, processed, ignored
}

// DeadLetterProcessor processes events from the dead letter queue.
type DeadLetterProcessor struct {
	bus       event.Bus
	cacheRepo repository.CacheRepository
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewDeadLetterProcessor creates a new dead letter processor.
func NewDeadLetterProcessor(bus event.Bus, cacheRepo repository.CacheRepository) *DeadLetterProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	return &DeadLetterProcessor{
		bus:       bus,
		cacheRepo: cacheRepo,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start starts the dead letter processor.
func (p *DeadLetterProcessor) Start() error {
	log.Info().Msg("Starting dead letter processor...")

	if err := p.bus.Subscribe(p.ctx, event.StreamDeadLetter, event.GroupDeadLetterProcessors, p.handleDeadLetter); err != nil {
		return err
	}

	log.Info().Msg("Dead letter processor started successfully")
	return nil
}

// Stop stops the dead letter processor.
func (p *DeadLetterProcessor) Stop() error {
	log.Info().Msg("Stopping dead letter processor...")
	p.cancel()
	log.Info().Msg("Dead letter processor stopped")
	return nil
}

// handleDeadLetter processes a dead letter event.
func (p *DeadLetterProcessor) handleDeadLetter(ctx context.Context, evt *event.Event) error {
	log.Warn().
		Str("event_id", evt.ID).
		Str("event_type", string(evt.Type)).
		Int("retries", evt.Retries).
		Msg("Processing dead letter event")

	// Store the failed event for later analysis
	failedEvent := FailedEvent{
		ID:        evt.ID,
		EventID:   evt.ID,
		EventType: string(evt.Type),
		Payload:   evt.Payload,
		Retries:   evt.Retries,
		FailedAt:  time.Now().UTC(),
		Status:    "pending",
	}

	// Store in Redis with expiration (30 days)
	key := "failed_event:" + evt.ID
	if err := p.cacheRepo.Set(ctx, key, failedEvent, 30*24*time.Hour); err != nil {
		log.Error().Err(err).Str("event_id", evt.ID).Msg("Failed to store dead letter event")
		return err
	}

	// Add to failed events index
	indexKey := "failed_events:index"
	if err := p.addToIndex(ctx, indexKey, evt.ID); err != nil {
		log.Error().Err(err).Msg("Failed to add event to index")
	}

	// Log detailed information for debugging
	log.Error().
		Str("event_id", evt.ID).
		Str("event_type", string(evt.Type)).
		Int("retries", evt.Retries).
		RawJSON("payload", evt.Payload).
		Msg("Event moved to dead letter queue - manual intervention may be required")

	return nil
}

// addToIndex adds an event ID to the failed events index.
func (p *DeadLetterProcessor) addToIndex(ctx context.Context, indexKey, eventID string) error {
	var index []string
	_ = p.cacheRepo.Get(ctx, indexKey, &index)

	index = append(index, eventID)

	// Keep only the last 1000 entries
	if len(index) > 1000 {
		index = index[len(index)-1000:]
	}

	return p.cacheRepo.Set(ctx, indexKey, index, 30*24*time.Hour)
}

// GetFailedEvents retrieves all failed events.
func (p *DeadLetterProcessor) GetFailedEvents(ctx context.Context) ([]FailedEvent, error) {
	indexKey := "failed_events:index"
	var index []string
	if err := p.cacheRepo.Get(ctx, indexKey, &index); err != nil {
		return nil, err
	}

	events := make([]FailedEvent, 0, len(index))
	for _, eventID := range index {
		var failedEvent FailedEvent
		key := "failed_event:" + eventID
		if err := p.cacheRepo.Get(ctx, key, &failedEvent); err != nil {
			continue
		}
		events = append(events, failedEvent)
	}

	return events, nil
}

// RetryEvent retries a failed event.
func (p *DeadLetterProcessor) RetryEvent(ctx context.Context, eventID string) error {
	key := "failed_event:" + eventID
	var failedEvent FailedEvent
	if err := p.cacheRepo.Get(ctx, key, &failedEvent); err != nil {
		return err
	}

	// Create a new event from the failed event
	evt := &event.Event{
		ID:        failedEvent.EventID,
		Type:      event.Type(failedEvent.EventType),
		Payload:   failedEvent.Payload,
		Timestamp: time.Now().UTC(),
		Version:   1,
		Retries:   0, // Reset retries
	}

	// Publish back to the appropriate stream
	if err := p.bus.Publish(ctx, evt); err != nil {
		return err
	}

	// Mark as processed
	now := time.Now().UTC()
	failedEvent.ProcessedAt = &now
	failedEvent.Status = "retried"
	return p.cacheRepo.Set(ctx, key, failedEvent, 30*24*time.Hour)
}

// IgnoreEvent marks a failed event as ignored.
func (p *DeadLetterProcessor) IgnoreEvent(ctx context.Context, eventID string) error {
	key := "failed_event:" + eventID
	var failedEvent FailedEvent
	if err := p.cacheRepo.Get(ctx, key, &failedEvent); err != nil {
		return err
	}

	now := time.Now().UTC()
	failedEvent.ProcessedAt = &now
	failedEvent.Status = "ignored"
	return p.cacheRepo.Set(ctx, key, failedEvent, 30*24*time.Hour)
}
