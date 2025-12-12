// Package dto for mapped the objects
package dto

import "time"

// HealthResponse representa la respuesta del endpoint de health check.
// Contiene el estado general del servicio y el estado de cada dependencia.
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}
