package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/service"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// AlertManagerWebhook represents the webhook payload from AlertManager.
type AlertManagerWebhook struct {
	Version           string              `json:"version"`
	GroupKey          string              `json:"groupKey"`
	TruncatedAlerts   int                 `json:"truncatedAlerts"`
	Status            string              `json:"status"`
	Receiver          string              `json:"receiver"`
	GroupLabels       map[string]string   `json:"groupLabels"`
	CommonLabels      map[string]string   `json:"commonLabels"`
	CommonAnnotations map[string]string   `json:"commonAnnotations"`
	ExternalURL       string              `json:"externalURL"`
	Alerts            []AlertManagerAlert `json:"alerts"`
}

// AlertManagerAlert represents a single alert from AlertManager.
type AlertManagerAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// WebhookHandler handles webhook endpoints.
type WebhookHandler struct {
	alertService *service.AlertService
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(alertService *service.AlertService) *WebhookHandler {
	return &WebhookHandler{
		alertService: alertService,
	}
}

// AlertManagerWebhookHandler handles POST /api/v1/webhooks/alertmanager
//
//	@Summary		Receive AlertManager webhook
//	@Description	Receives alerts from Prometheus AlertManager
//	@Tags			webhooks
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	AlertManagerWebhook	true	"AlertManager webhook payload"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse
//	@Router			/webhooks/alertmanager [post]
func (h *WebhookHandler) AlertManagerWebhookHandler(c *fiber.Ctx) error {
	var payload AlertManagerWebhook
	if err := c.BodyParser(&payload); err != nil {
		log.Error().Err(err).Msg("Failed to parse AlertManager webhook")
		return helper.BadRequest(c, "Invalid webhook payload")
	}

	log.Info().
		Str("status", payload.Status).
		Str("receiver", payload.Receiver).
		Int("alert_count", len(payload.Alerts)).
		Msg("Received AlertManager webhook")

	for _, alert := range payload.Alerts {
		if err := h.processAlert(c, alert); err != nil {
			log.Error().Err(err).Str("fingerprint", alert.Fingerprint).Msg("Failed to process alert")
		}
	}

	return helper.Success(c, fiber.Map{"status": "received"})
}

// processAlert processes a single AlertManager alert.
func (h *WebhookHandler) processAlert(c *fiber.Ctx, alert AlertManagerAlert) error {
	severity := h.mapSeverity(alert.Labels["severity"])

	title := alert.Labels["alertname"]
	if title == "" {
		title = "AlertManager Alert"
	}

	message := alert.Annotations["description"]
	if message == "" {
		message = alert.Annotations["summary"]
	}
	if message == "" {
		message = "Alert triggered from Prometheus"
	}

	source := "alertmanager"
	if instance, ok := alert.Labels["instance"]; ok {
		source = "alertmanager:" + instance
	}

	// Only create alerts for firing status
	if alert.Status == "firing" {
		input := service.CreateAlertInput{
			Title:    title,
			Message:  message,
			Severity: severity,
			Source:   source,
			Metadata: map[string]interface{}{
				"fingerprint":   alert.Fingerprint,
				"generator_url": alert.GeneratorURL,
				"labels":        alert.Labels,
				"annotations":   alert.Annotations,
				"starts_at":     alert.StartsAt,
			},
		}

		_, err := h.alertService.Create(c.Context(), input)
		if err != nil {
			return err
		}

		log.Info().
			Str("alertname", title).
			Str("severity", string(severity)).
			Str("fingerprint", alert.Fingerprint).
			Msg("Created alert from AlertManager")
	} else {
		log.Info().
			Str("alertname", title).
			Str("status", alert.Status).
			Str("fingerprint", alert.Fingerprint).
			Msg("Alert resolved in AlertManager")
	}

	return nil
}

// mapSeverity maps AlertManager severity to entity severity.
func (h *WebhookHandler) mapSeverity(severity string) entity.AlertSeverity {
	switch severity {
	case "critical":
		return entity.AlertSeverityCritical
	case "warning", "high":
		return entity.AlertSeverityHigh
	case "info", "medium":
		return entity.AlertSeverityMedium
	default:
		return entity.AlertSeverityLow
	}
}
