package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/tracing"
)

// TracingMiddleware adds distributed tracing to HTTP requests.
func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract trace context from incoming request
		ctx := c.UserContext()
		propagator := propagation.TraceContext{}
		ctx = propagator.Extract(ctx, &headerCarrier{c: c})

		// Start span
		ctx, span := tracing.StartSpan(ctx, c.Method()+" "+c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Set span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.route", c.Route().Path),
			attribute.String("http.url", c.OriginalURL()),
			attribute.String("http.user_agent", c.Get("User-Agent")),
			attribute.String("http.host", c.Hostname()),
			attribute.String("http.request_id", c.Locals("requestid").(string)),
			attribute.String("net.peer.ip", c.IP()),
		)

		// Set trace ID in response header
		traceID := span.SpanContext().TraceID().String()
		c.Set("X-Trace-ID", traceID)

		// Update context
		c.SetUserContext(ctx)

		// Process request
		err := c.Next()

		// Set response status
		span.SetAttributes(attribute.Int("http.status_code", c.Response().StatusCode()))

		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

// headerCarrier adapts Fiber context to propagation.TextMapCarrier.
type headerCarrier struct {
	c *fiber.Ctx
}

func (h *headerCarrier) Get(key string) string {
	return h.c.Get(key)
}

func (h *headerCarrier) Set(key, value string) {
	h.c.Set(key, value)
}

func (h *headerCarrier) Keys() []string {
	headers := h.c.GetReqHeaders()
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	return keys
}
