package websocket

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/metrics"
)

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Clients indexed by user ID for targeted messages
	userClients map[entity.ID]map[*Client]bool

	// Inbound messages from clients to broadcast
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[entity.ID]map[*Client]bool),
		broadcast:   make(chan []byte, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

// Run starts the hub's main loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient adds a client to the hub.
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// Add to user-specific map if authenticated
	if client.userID != nil {
		if h.userClients[*client.userID] == nil {
			h.userClients[*client.userID] = make(map[*Client]bool)
		}
		h.userClients[*client.userID][client] = true
	}

	// Update Prometheus metrics
	metrics.WebSocketConnectionsTotal.Inc()
	metrics.WebSocketConnectionsActive.Set(float64(len(h.clients)))

	log.Info().
		Int("total_clients", len(h.clients)).
		Msg("WebSocket client connected")
}

// unregisterClient removes a client from the hub.
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; !ok {
		return
	}

	delete(h.clients, client)

	// Remove from user-specific map
	if client.userID != nil {
		if clients, ok := h.userClients[*client.userID]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.userClients, *client.userID)
			}
		}
	}

	// Update Prometheus metrics
	metrics.WebSocketConnectionsActive.Set(float64(len(h.clients)))

	log.Info().
		Int("total_clients", len(h.clients)).
		Msg("WebSocket client disconnected")
}

// broadcastMessage sends a message to all connected clients.
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		client.Send(message)
	}

	// Update messages sent metric
	metrics.WebSocketMessagesSent.Add(float64(len(h.clients)))
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal broadcast message")
		return
	}

	h.broadcast <- data
}

// BroadcastToUser sends a message to all connections of a specific user.
func (h *Hub) BroadcastToUser(userID entity.ID, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.userClients[userID]
	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal user message")
		return
	}

	for client := range clients {
		client.Send(data)
	}

	// Update messages sent metric
	metrics.WebSocketMessagesSent.Add(float64(len(clients)))
}

// BroadcastToRole sends a message to all users with a specific role.
func (h *Hub) BroadcastToRole(role string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal role message")
		return
	}

	count := 0
	for client := range h.clients {
		if client.userRole == role {
			client.Send(data)
			count++
		}
	}

	// Update messages sent metric
	metrics.WebSocketMessagesSent.Add(float64(count))
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Register adds a client to the hub.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
