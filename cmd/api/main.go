// package main
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load .env file (optional, for development)
	if err := godotenv.Load(); err != nil {
		// .env file not found is okay, we'll use config.yaml or env vars
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Configure logger based on config
	setupLogger(cfg)

	log.Info().
		Str("app", cfg.App.Name).
		Str("version", cfg.App.Version).
		Str("env", cfg.App.Env).
		Msg("Starting Real-Time Alerting System...")

	log.Debug().
		Str("server_address", cfg.Server.Address()).
		Str("database_host", cfg.Database.Host).
		Str("redis_host", cfg.Redis.Host).
		Msg("Configuration loaded")

	// TODO: Initialize database
	// TODO: Initialize Redis
	// TODO: Initialize HTTP server
	// TODO: Initialize WebSocket

	log.Info().
		Str("address", cfg.Server.Address()).
		Msg("Server ready")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
}

func setupLogger(cfg *config.Config) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set output format
	if cfg.Logging.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Add caller info in development
	if cfg.App.IsDevelopment() {
		log.Logger = log.With().Caller().Logger()
	}
}
