package dto

import "time"

// HealthResponse representa la respuesta del endpoint de health check.
// Contiene el estado general del servicio y el estado de cada dependencia.
type HealthResponse struct {
	Status    string            `json:"status"`    // "healthy", "unhealthy", "degraded"
	Timestamp time.Time         `json:"timestamp"` // Momento de la verificación
	Version   string            `json:"version"`   // Versión de la aplicación
	Services  map[string]string `json:"services"`  // Estado de cada servicio externo
}
