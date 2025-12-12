// Package handler provides HTTP handlers for the application.
package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
)

// HealthHandler maneja los endpoints relacionados con el estado del servicio.
type HealthHandler struct {
	config *config.Config
}

// NewHealthHandler crea una nueva instancia del handler.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

// Check handles the health check endpoint, returning the overall status of the service.
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	response := dto.HealthResponse{
		Status:    "healthy",
		Timestamp: c.Context().Time(),
		Version:   h.config.App.Version,
		Services: map[string]string{
			"postgres": "healthy",
			"redis":    "healthy",
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Ready verifica si la aplicación está lista para recibir tráfico.
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ready",
	})
}

// Live verifica si la aplicación está viva (no está colgada).
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "alive",
	})
}
