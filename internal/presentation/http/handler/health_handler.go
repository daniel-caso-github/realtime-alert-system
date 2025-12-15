package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	config *config.Config
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

// Check handles GET /health
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

	return helper.Success(c, response)
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return helper.Success(c, dto.ReadyResponse{Status: "ready"})
}

// Live handles GET /live
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return helper.Success(c, dto.LiveResponse{Status: "alive"})
}
