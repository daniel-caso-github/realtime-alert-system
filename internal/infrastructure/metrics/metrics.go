// Package metrics provides Prometheus metrics for the application.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP metrics.
var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)

// Alert metrics.
var (
	AlertsCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alerts_created_total",
			Help: "Total number of alerts created",
		},
		[]string{"severity", "source"},
	)

	AlertsAcknowledgedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "alerts_acknowledged_total",
			Help: "Total number of alerts acknowledged",
		},
	)

	AlertsResolvedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "alerts_resolved_total",
			Help: "Total number of alerts resolved",
		},
	)

	AlertsDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "alerts_deleted_total",
			Help: "Total number of alerts deleted",
		},
	)

	AlertsActiveGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "alerts_active",
			Help: "Current number of active alerts",
		},
	)
)

// Event bus metrics.
var (
	EventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_published_total",
			Help: "Total number of events published",
		},
		[]string{"event_type", "stream"},
	)

	EventsConsumedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_consumed_total",
			Help: "Total number of events consumed",
		},
		[]string{"event_type", "status"},
	)

	EventsFailedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_failed_total",
			Help: "Total number of events that failed processing",
		},
		[]string{"event_type"},
	)

	EventProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_processing_duration_seconds",
			Help:    "Event processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_type"},
	)
)

// WebSocket metrics.
var (
	WebSocketConnectionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_connections_total",
			Help: "Total number of WebSocket connections",
		},
	)

	WebSocketConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Current number of active WebSocket connections",
		},
	)

	WebSocketMessagesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_total",
			Help: "Total number of WebSocket messages sent",
		},
	)
)

// Database metrics.
var (
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Current number of active database connections",
		},
	)
)

// Cache metrics.
var (
	CacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

// Circuit breaker metrics.
var (
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
		},
		[]string{"name"},
	)

	CircuitBreakerFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_failures_total",
			Help: "Total number of circuit breaker failures",
		},
		[]string{"name"},
	)
)

// Authentication metrics.
var (
	AuthLoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"},
	)

	AuthTokensIssued = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_tokens_issued_total",
			Help: "Total number of tokens issued",
		},
	)
)
