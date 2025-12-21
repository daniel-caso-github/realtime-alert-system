// Package logger provides structured logging utilities.
package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ContextKey is a type for context keys.
type ContextKey string

// Context keys for logging.
const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	TraceIDKey   ContextKey = "trace_id"
	SpanIDKey    ContextKey = "span_id"
)

// Config holds logger configuration.
type Config struct {
	Level      string
	Format     string // "json" or "console"
	TimeFormat string
	Caller     bool
}

// Setup initializes the global logger.
func Setup(cfg Config) {
	// Set time format
	if cfg.TimeFormat != "" {
		zerolog.TimeFieldFormat = cfg.TimeFormat
	} else {
		zerolog.TimeFieldFormat = time.RFC3339
	}

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output format
	if cfg.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.Kitchen,
		})
	}

	// Add caller information
	if cfg.Caller {
		log.Logger = log.With().Caller().Logger()
	}
}

// WithContext returns a logger with context values.
func WithContext(ctx context.Context) zerolog.Logger {
	logger := log.Logger

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		logger = logger.With().Str("request_id", requestID).Logger()
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		logger = logger.With().Str("user_id", userID).Logger()
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		logger = logger.With().Str("trace_id", traceID).Logger()
	}

	if spanID, ok := ctx.Value(SpanIDKey).(string); ok && spanID != "" {
		logger = logger.With().Str("span_id", spanID).Logger()
	}

	return logger
}

// WithRequestID adds request ID to context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds user ID to context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithTraceID adds trace ID to context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithSpanID adds span ID to context.
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, SpanIDKey, spanID)
}

// Info logs an info message with context.
func Info(ctx context.Context, msg string) {
	l := WithContext(ctx)
	l.Info().Msg(msg)
}

// Error logs an error message with context.
func Error(ctx context.Context, err error, msg string) {
	l := WithContext(ctx)
	l.Error().Err(err).Msg(msg)
}

// Debug logs a debug message with context.
func Debug(ctx context.Context, msg string) {
	l := WithContext(ctx)
	l.Debug().Msg(msg)
}

// Warn logs a warning message with context.
func Warn(ctx context.Context, msg string) {
	l := WithContext(ctx)
	l.Warn().Msg(msg)
}
