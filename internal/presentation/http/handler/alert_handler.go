// Package handler provides HTTP request handlers for the API.
package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/service"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"

	"github.com/rs/zerolog/log"
)

// AlertHandler handles alert-related HTTP requests.
type AlertHandler struct {
	alertService *service.AlertService
}

// NewAlertHandler creates a new alert handler.
func NewAlertHandler(alertService *service.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

// Create handles POST /api/v1/alerts
//
//	@Summary		Create alert
//	@Description	Create a new alert
//	@Tags			alerts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateAlertRequest	true	"Alert data"
//	@Success		201		{object}	dto.AlertResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		422		{object}	dto.ValidationErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts [post]
func (h *AlertHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if errors := helper.ValidateStruct(req); len(errors) > 0 {
		return helper.ValidationErrors(c, errors)
	}

	// Create alert
	input := service.CreateAlertInput{
		Title:    req.Title,
		Message:  req.Message,
		Severity: entity.AlertSeverity(req.Severity),
		Source:   req.Source,
		Metadata: req.Metadata,
	}

	alert, err := h.alertService.Create(c.Context(), input)
	if err != nil {
		return helper.InternalError(c, "Failed to create alert")
	}

	return helper.Created(c, dto.AlertFromEntity(alert))
}

// GetByID handles GET /api/v1/alerts/:id
//
//	@Summary		Get alert by ID
//	@Description	Retrieve a specific alert
//	@Tags			alerts
//	@Produce		json
//	@Param			id	path		string	true	"Alert ID"
//	@Success		200	{object}	dto.AlertResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts/{id} [get]
func (h *AlertHandler) GetByID(c *fiber.Ctx) error {
	id, err := entity.ParseID(c.Params("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid alert ID")
	}

	alert, err := h.alertService.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			return helper.NotFound(c, "Alert not found")
		}
		return helper.InternalError(c, "Failed to get alert")
	}

	return helper.Success(c, dto.AlertFromEntity(alert))
}

// List handles GET /api/v1/alerts
//
//	@Summary		List alerts
//	@Description	Retrieve paginated list of alerts with optional filters
//	@Tags			alerts
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	default(20)
//	@Param			status		query		[]string	false	"Filter by status"
//	@Param			severity	query		[]string	false	"Filter by severity"
//	@Param			source		query		string	false	"Filter by source"
//	@Param			search		query		string	false	"Search in title/message"
//	@Success		200			{object}	dto.PaginatedAlertResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts [get]
func (h *AlertHandler) List(c *fiber.Ctx) error {
	var req dto.ListAlertsRequest
	if err := c.QueryParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid query parameters")
	}

	// Build filter
	filter := valueobject.NewAlertFilter()

	if len(req.Status) > 0 {
		statuses := make([]entity.AlertStatus, len(req.Status))
		for i, s := range req.Status {
			statuses[i] = entity.AlertStatus(s)
		}
		filter = filter.WithStatuses(statuses...)
	}

	if len(req.Severity) > 0 {
		severities := make([]entity.AlertSeverity, len(req.Severity))
		for i, s := range req.Severity {
			severities[i] = entity.AlertSeverity(s)
		}
		filter = filter.WithSeverities(severities...)
	}

	if req.Source != "" {
		filter = filter.WithSource(req.Source)
	}

	if req.Search != "" {
		filter = filter.WithSearch(req.Search)
	}

	filter = applyDateFilter(filter, req.FromDate, req.ToDate)

	// Build pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	pagination := valueobject.NewPagination(page, pageSize)

	// Get alerts
	result, err := h.alertService.List(c.Context(), service.ListInput{
		Filter:     filter,
		Pagination: pagination,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create alert")
		return helper.InternalError(c, "Failed to create alert")
	}

	// Build response
	response := dto.PaginatedResponse[dto.AlertResponse]{
		Items:       dto.AlertsFromEntities(result.Items),
		TotalItems:  result.TotalItems,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		PageSize:    result.PageSize,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}

	return helper.Success(c, response)
}

// Acknowledge handles POST /api/v1/alerts/:id/acknowledge
//
//	@Summary		Acknowledge alert
//	@Description	Mark an alert as acknowledged
//	@Tags			alerts
//	@Produce		json
//	@Param			id	path		string	true	"Alert ID"
//	@Success		200	{object}	dto.AlertResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts/{id}/acknowledge [post]
func (h *AlertHandler) Acknowledge(c *fiber.Ctx) error {
	alertID, err := entity.ParseID(c.Params("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid alert ID")
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(entity.ID)
	if !ok {
		return helper.Unauthorized(c, "User not authenticated")
	}

	alert, err := h.alertService.Acknowledge(c.Context(), alertID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			return helper.NotFound(c, "Alert not found")
		}
		if errors.Is(err, entity.ErrAlertAlreadyAcknowledged) {
			return helper.Conflict(c, "Alert is already acknowledged")
		}
		if errors.Is(err, entity.ErrAlertAlreadyResolved) {
			return helper.Conflict(c, "Alert is already resolved")
		}
		return helper.InternalError(c, "Failed to acknowledge alert")
	}

	return helper.Success(c, dto.AlertFromEntity(alert))
}

// Resolve handles POST /api/v1/alerts/:id/resolve
//
//	@Summary		Resolve alert
//	@Description	Mark an alert as resolved
//	@Tags			alerts
//	@Produce		json
//	@Param			id	path		string	true	"Alert ID"
//	@Success		200	{object}	dto.AlertResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts/{id}/resolve [post]
func (h *AlertHandler) Resolve(c *fiber.Ctx) error {
	alertID, err := entity.ParseID(c.Params("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid alert ID")
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(entity.ID)
	if !ok {
		return helper.Unauthorized(c, "User not authenticated")
	}

	alert, err := h.alertService.Resolve(c.Context(), alertID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			return helper.NotFound(c, "Alert not found")
		}
		if errors.Is(err, entity.ErrAlertAlreadyResolved) {
			return helper.Conflict(c, "Alert is already resolved")
		}
		return helper.InternalError(c, "Failed to resolve alert")
	}

	return helper.Success(c, dto.AlertFromEntity(alert))
}

// Delete handles DELETE /api/v1/alerts/:id
//
//	@Summary		Delete alert
//	@Description	Remove an alert (admin only)
//	@Tags			alerts
//	@Param			id	path	string	true	"Alert ID"
//	@Success		204
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts/{id} [delete]
func (h *AlertHandler) Delete(c *fiber.Ctx) error {
	id, err := entity.ParseID(c.Params("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid alert ID")
	}

	if err := h.alertService.Delete(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			return helper.NotFound(c, "Alert not found")
		}
		return helper.InternalError(c, "Failed to delete alert")
	}

	return helper.NoContent(c)
}

// GetStatistics handles GET /api/v1/alerts/statistics
//
//	@Summary		Get alert statistics
//	@Description	Retrieve aggregated alert statistics
//	@Tags			alerts
//	@Produce		json
//	@Success		200	{object}	dto.AlertStatisticsResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/alerts/statistics [get]
func (h *AlertHandler) GetStatistics(c *fiber.Ctx) error {
	stats, err := h.alertService.GetStatistics(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get statistics")
		return helper.InternalError(c, "Failed to get statistics")
	}

	response := dto.AlertStatisticsResponse{
		TotalAlerts:        stats.TotalAlerts,
		ActiveAlerts:       stats.ActiveAlerts,
		AcknowledgedAlerts: stats.AcknowledgedAlerts,
		ResolvedAlerts:     stats.ResolvedAlerts,
		BySeverity:         stats.BySeverity,
		BySource:           stats.BySource,
	}

	return helper.Success(c, response)
}

// applyDateFilter applies date range filter if valid dates are provided.
func applyDateFilter(filter valueobject.AlertFilter, fromDate, toDate string) valueobject.AlertFilter {
	if fromDate == "" {
		return filter
	}

	from, err := time.Parse(time.RFC3339, fromDate)
	if err != nil {
		return filter
	}

	if toDate == "" {
		return filter
	}

	to, err := time.Parse(time.RFC3339, toDate)
	if err != nil {
		return filter
	}

	return filter.WithDateRange(from, to)
}
