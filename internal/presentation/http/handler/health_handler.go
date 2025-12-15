package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/helper"
)

// Health status constants.
const (
	statusHealthy       = "healthy"
	statusUnhealthy     = "unhealthy"
	statusDegraded      = "degraded"
	statusNotConfigured = "not configured"
	statusReady         = "ready"
	statusNotReady      = "not ready"
	statusAlive         = "alive"
)

// HealthChecker defines the interface for health checking services.
type HealthChecker interface {
	Health(ctx context.Context) error
}

// CacheHealthChecker defines the interface for cache health checking.
type CacheHealthChecker interface {
	Ping(ctx context.Context) error
}

// WebSocketStats defines the interface for WebSocket statistics.
type WebSocketStats interface {
	ClientCount() int
}

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	config  *config.Config
	db      HealthChecker
	cache   CacheHealthChecker
	wsStats WebSocketStats
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(cfg *config.Config, db HealthChecker, cache CacheHealthChecker, wsStats WebSocketStats) *HealthHandler {
	return &HealthHandler{
		config:  cfg,
		db:      db,
		cache:   cache,
		wsStats: wsStats,
	}
}

// Check handles GET /health
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)
	status := statusHealthy

	// Check PostgreSQL
	if h.db != nil {
		if err := h.db.Health(ctx); err != nil {
			services["postgres"] = statusUnhealthy
			status = statusDegraded
		} else {
			services["postgres"] = statusHealthy
		}
	} else {
		services["postgres"] = statusNotConfigured
	}

	// Check Redis
	if h.cache != nil {
		if err := h.cache.Ping(ctx); err != nil {
			services["redis"] = statusUnhealthy
			status = statusDegraded
		} else {
			services["redis"] = statusHealthy
		}
	} else {
		services["redis"] = statusNotConfigured
	}

	// WebSocket status
	if h.wsStats != nil {
		services["websocket"] = statusHealthy
		services["websocket_clients"] = fmt.Sprintf("%d", h.wsStats.ClientCount())
	}

	response := dto.HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Version:   h.config.App.Version,
		Services:  services,
	}

	if status == statusHealthy {
		return helper.Success(c, response)
	}
	return helper.JSON(c, fiber.StatusServiceUnavailable, response)
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if h.db != nil {
		if err := h.db.Health(ctx); err != nil {
			return helper.JSON(c, fiber.StatusServiceUnavailable, dto.ReadyResponse{Status: statusNotReady})
		}
	}

	return helper.Success(c, dto.ReadyResponse{Status: statusReady})
}

// Live handles GET /live
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return helper.Success(c, dto.LiveResponse{Status: statusAlive})
}
