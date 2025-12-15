package websocket

import (
	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// AlertPublisher publishes alert events to WebSocket clients.
type AlertPublisher struct {
	hub *Hub
}

// NewAlertPublisher creates a new alert publisher.
func NewAlertPublisher(hub *Hub) *AlertPublisher {
	return &AlertPublisher{
		hub: hub,
	}
}

// PublishAlertCreated broadcasts a new alert to all clients.
func (p *AlertPublisher) PublishAlertCreated(alert *entity.Alert) {
	msg := NewAlertCreatedMessage(dto.AlertFromEntity(alert))
	p.hub.Broadcast(msg)
}

// PublishAlertAcknowledged broadcasts an acknowledged alert to all clients.
func (p *AlertPublisher) PublishAlertAcknowledged(alert *entity.Alert) {
	msg := NewAlertAcknowledgedMessage(dto.AlertFromEntity(alert))
	p.hub.Broadcast(msg)
}

// PublishAlertResolved broadcasts a resolved alert to all clients.
func (p *AlertPublisher) PublishAlertResolved(alert *entity.Alert) {
	msg := NewAlertResolvedMessage(dto.AlertFromEntity(alert))
	p.hub.Broadcast(msg)
}

// PublishAlertDeleted broadcasts a deleted alert to all clients.
func (p *AlertPublisher) PublishAlertDeleted(alertID string) {
	msg := NewAlertDeletedMessage(alertID)
	p.hub.Broadcast(msg)
}
