package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/metrics"
)

// PrometheusMiddleware collects HTTP metrics.
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()
		path := c.Route().Path // Use route path to avoid high cardinality

		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
