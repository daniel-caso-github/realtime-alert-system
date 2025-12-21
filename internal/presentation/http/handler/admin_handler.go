package handler

import (
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/circuitbreaker"
	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/worker"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// AdminHandler handles admin endpoints.
type AdminHandler struct {
	deadLetterProcessor *worker.DeadLetterProcessor
	eventWorker         *worker.EventWorker
	cbRegistry          *circuitbreaker.Registry
}

// NewAdminHandler creates a new admin handler.
func NewAdminHandler(dlp *worker.DeadLetterProcessor, ew *worker.EventWorker, cbRegistry *circuitbreaker.Registry) *AdminHandler {
	return &AdminHandler{
		deadLetterProcessor: dlp,
		eventWorker:         ew,
		cbRegistry:          cbRegistry,
	}
}

// Add this method:

// GetCircuitBreakerStats handles GET /api/v1/admin/circuit-breakers
//
//	@Summary		Get circuit breaker stats
//	@Description	Retrieve circuit breaker statistics
//	@Tags			admin
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/circuit-breakers [get]
func (h *AdminHandler) GetCircuitBreakerStats(c *fiber.Ctx) error {
	if h.cbRegistry == nil {
		return helper.Success(c, map[string]interface{}{})
	}

	return helper.Success(c, h.cbRegistry.Stats())
}

// GetFailedEvents handles GET /api/v1/admin/failed-events
//
//	@Summary		Get failed events
//	@Description	Retrieve all events in the dead letter queue
//	@Tags			admin
//	@Produce		json
//	@Success		200	{array}		map[string]interface{}
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/failed-events [get]
func (h *AdminHandler) GetFailedEvents(c *fiber.Ctx) error {
	if h.deadLetterProcessor == nil {
		return helper.Success(c, []worker.FailedEvent{})
	}

	events, err := h.deadLetterProcessor.GetFailedEvents(c.Context())
	if err != nil {
		return helper.InternalError(c, "Failed to retrieve failed events")
	}

	return helper.Success(c, events)
}

// RetryFailedEvent handles POST /api/v1/admin/failed-events/:id/retry
//
//	@Summary		Retry failed event
//	@Description	Retry a failed event from the dead letter queue
//	@Tags			admin
//	@Param			id	path	string	true	"Event ID"
//	@Success		204
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/failed-events/{id}/retry [post]
func (h *AdminHandler) RetryFailedEvent(c *fiber.Ctx) error {
	if h.deadLetterProcessor == nil {
		return helper.NotFound(c, "Dead letter processor not available")
	}

	eventID := c.Params("id")
	if err := h.deadLetterProcessor.RetryEvent(c.Context(), eventID); err != nil {
		return helper.InternalError(c, "Failed to retry event")
	}

	return helper.NoContent(c)
}

// IgnoreFailedEvent handles POST /api/v1/admin/failed-events/:id/ignore
//
//	@Summary		Ignore failed event
//	@Description	Mark a failed event as ignored
//	@Tags			admin
//	@Param			id	path	string	true	"Event ID"
//	@Success		204
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/failed-events/{id}/ignore [post]
func (h *AdminHandler) IgnoreFailedEvent(c *fiber.Ctx) error {
	if h.deadLetterProcessor == nil {
		return helper.NotFound(c, "Dead letter processor not available")
	}

	eventID := c.Params("id")
	if err := h.deadLetterProcessor.IgnoreEvent(c.Context(), eventID); err != nil {
		return helper.InternalError(c, "Failed to ignore event")
	}

	return helper.NoContent(c)
}

// GetEventMetrics handles GET /api/v1/admin/metrics/events
//
//	@Summary		Get event metrics
//	@Description	Retrieve event processing metrics
//	@Tags			admin
//	@Produce		json
//	@Success		200	{object}	map[string]int64
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/admin/metrics/events [get]
func (h *AdminHandler) GetEventMetrics(c *fiber.Ctx) error {
	if h.eventWorker == nil {
		return helper.Success(c, map[string]int64{})
	}

	return helper.Success(c, h.eventWorker.GetMetrics())
}
