// Package main provides the entry point for the Real-Time Alerting System API.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/router"
)

func main() {
	// Load environment variables from .env file (development only).
	// In production, environment variables should be set by the container orchestrator.
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load envs")
	}

	// Load application configuration from config files and environment variables.
	// The empty string argument uses the default config path.
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Configure the global zerolog logger based on application settings.
	setupLogger(cfg)

	log.Info().
		Str("app", cfg.App.Name).
		Str("version", cfg.App.Version).
		Str("env", cfg.App.Env).
		Msg("🚀 Starting Real-Time Alerting System...")

	// Create the Fiber application with all routes configured.
	app := router.Setup(cfg)

	// Start the HTTP server in a separate goroutine to allow
	// graceful shutdown handling in the main thread.
	go func() {
		log.Info().
			Str("address", cfg.Server.Address()).
			Msg("✅ HTTP server started")

		if err := app.Listen(cfg.Server.Address()); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Wait for termination signals (Ctrl+C or kill command).
	// This blocks until a signal is received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("🛑 Shutting down server...")

	// Create a context with timeout for graceful shutdown.
	// This gives in-flight requests up to 10 seconds to complete.
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the Fiber server gracefully, waiting for active connections to close.
	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	}

	log.Info().Msg("👋 Server stopped gracefully")
}

// setupLogger configures the global zerolog logger based on application configuration.
// It sets the log level, output format, and optionally adds caller information.
//
// The function supports the following configurations:
//   - Log level: Parsed from cfg.Logging.Level (defaults to debug if invalid)
//   - Output format: "console" for human-readable output, JSON otherwise
//   - Caller info: Enabled in development mode to show file:line in logs
func setupLogger(cfg *config.Config) {
	// Parsear el nivel de log desde la configuración
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	// En desarrollo, usar formato legible para humanos
	if cfg.Logging.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// En desarrollo, agregar información del caller (archivo:línea)
	if cfg.App.IsDevelopment() {
		log.Logger = log.With().Caller().Logger()
	}
}
