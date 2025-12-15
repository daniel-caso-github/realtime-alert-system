package websocket

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// Handler handles WebSocket connections.
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler.
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// Upgrade is middleware that checks if the request is a WebSocket upgrade request.
func (h *Handler) Upgrade(c *fiber.Ctx) error {
	if fiberws.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Handle handles WebSocket connections.
func (h *Handler) Handle(c *fiberws.Conn) {
	var userID *entity.ID
	var userRole string

	if id, ok := c.Locals("userID").(entity.ID); ok {
		userID = &id
	}
	if role, ok := c.Locals("userRole").(string); ok {
		userRole = role
	}

	client := NewClient(h.hub, c.Conn, userID, userRole)
	h.hub.Register(client)

	log.Debug().
		Bool("authenticated", userID != nil).
		Str("role", userRole).
		Msg("New WebSocket connection")

	go client.WritePump()
	client.ReadPump()
}

// GetHub returns the hub instance.
func (h *Handler) GetHub() *Hub {
	return h.hub
}
