package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
)

// HealthHandler maneja los endpoints relacionados con el estado del servicio.
// Contiene las dependencias necesarias para verificar la salud del sistema.
type HealthHandler struct {
	config *config.Config
	// Aquí agregaremos más adelante:
	// db    *database.PostgresDB
	// redis *database.RedisClient
}

// NewHealthHandler crea una nueva instancia del handler.
// Recibe la configuración como dependencia (Dependency Injection).
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

// Check verifica el estado de la aplicación y sus dependencias.
// GET /health
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	// Por ahora retornamos healthy ya que no tenemos conexiones a DB todavía.
	// En la Fase 3 agregaremos las verificaciones reales de PostgreSQL y Redis.

	response := dto.HealthResponse{
		Status:    "healthy",
		Timestamp: c.Context().Time(),
		Version:   h.config.App.Version,
		Services: map[string]string{
			"postgres": "healthy", // TODO: verificación real
			"redis":    "healthy", // TODO: verificación real
		},
	}

	// Retornamos JSON con status 200 OK
	return c.Status(fiber.StatusOK).JSON(response)
}

// Ready verifica si la aplicación está lista para recibir tráfico.
// GET /ready
// Kubernetes usa esto para el readinessProbe.
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Por ahora siempre estamos listos.
	// Más adelante verificaremos que las conexiones a DB estén establecidas.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ready",
	})
}

// Live verifica si la aplicación está viva (no está colgada).
// GET /live
// Kubernetes usa esto para el livenessProbe.
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	// Si podemos responder, estamos vivos.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "alive",
	})
}
