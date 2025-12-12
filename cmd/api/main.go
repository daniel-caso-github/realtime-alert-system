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
	// =========================================================================
	// CONFIGURACIÓN INICIAL
	// =========================================================================

	// Cargar variables de entorno desde .env (solo en desarrollo)
	if err := godotenv.Load(); err != nil {
		// No es error si no existe .env, usaremos config.yaml o variables de entorno
	}

	// Cargar configuración desde archivo y variables de entorno
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Configurar el logger según la configuración
	setupLogger(cfg)

	// =========================================================================
	// INICIO DE LA APLICACIÓN
	// =========================================================================

	log.Info().
		Str("app", cfg.App.Name).
		Str("version", cfg.App.Version).
		Str("env", cfg.App.Env).
		Msg("🚀 Starting Real-Time Alerting System...")

	// =========================================================================
	// CONFIGURAR SERVIDOR HTTP
	// =========================================================================

	// Crear la aplicación Fiber con todas las rutas
	app := router.Setup(cfg)

	// =========================================================================
	// INICIAR SERVIDOR EN GOROUTINE
	// =========================================================================

	// Iniciamos el servidor en una goroutine separada para poder
	// manejar el graceful shutdown en el hilo principal.
	go func() {
		log.Info().
			Str("address", cfg.Server.Address()).
			Msg("✅ HTTP server started")

		if err := app.Listen(cfg.Server.Address()); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// =========================================================================
	// GRACEFUL SHUTDOWN
	// =========================================================================

	// Esperamos señales de terminación (Ctrl+C o kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("🛑 Shutting down server...")

	// Crear contexto con timeout para el shutdown
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cerrar el servidor Fiber de forma ordenada
	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	}

	// TODO: Cerrar conexiones a base de datos
	// TODO: Cerrar conexiones a Redis

	log.Info().Msg("👋 Server stopped gracefully")
}

// setupLogger configura zerolog según la configuración de la aplicación.
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
