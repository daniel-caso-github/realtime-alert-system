// Package main is the entry point for the alerting system API.
//
//	@title						Real-Time Alerting System API
//	@version					1.0
//	@description				Enterprise-grade distributed real-time alerting system
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.email				support@alerting.local
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
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
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/database"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/messaging"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/tracing"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/worker"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/router"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/websocket"
)

func main() {
	// Load .env file (optional in production)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logger
	setupLogger(cfg)

	log.Info().
		Str("app", cfg.App.Name).
		Str("version", cfg.App.Version).
		Str("env", cfg.App.Env).
		Msg("Starting application...")

	// Initialize PostgreSQL
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	log.Info().Msg("Connected to PostgreSQL")

	// Initialize Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		closeDB(db)
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	log.Info().Msg("Connected to Redis")

	// Initialize tracing (after critical connections, so defer works properly)
	shutdownTracer, err := tracing.InitTracer(tracing.Config{
		ServiceName:    cfg.App.Name,
		ServiceVersion: cfg.App.Version,
		Environment:    cfg.App.Env,
		JaegerEndpoint: cfg.Tracing.JaegerEndpoint,
		Enabled:        cfg.Tracing.Enabled,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to initialize tracing, continuing without it")
	} else {
		log.Info().Msg("Tracing initialized")
		defer func() {
			if err := shutdownTracer(context.Background()); err != nil {
				log.Error().Err(err).Msg("Error shutting down tracer")
			}
		}()
	}

	// Initialize repositories
	userRepo := database.NewPostgresUserRepository(db)
	alertRepo := database.NewPostgresAlertRepository(db)
	cacheRepo := database.NewRedisCacheRepository(redisClient)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	log.Info().Msg("WebSocket hub started")

	// Initialize Event Bus
	eventBus := messaging.NewRedisStreamBus(redisClient.GetClient(), cfg.EventBus.ConsumerID)
	retryConfig := messaging.RetryConfig{
		MaxRetries:     cfg.EventBus.MaxRetries,
		InitialBackoff: cfg.EventBus.InitialBackoff,
		MaxBackoff:     cfg.EventBus.MaxBackoff,
		Multiplier:     cfg.EventBus.Multiplier,
		Jitter:         true,
	}
	retryableBus := messaging.NewRetryableBus(eventBus, retryConfig)
	log.Info().Msg("Event bus initialized")

	// Initialize Event Worker
	eventWorker := worker.NewEventWorker(retryableBus)
	if err := eventWorker.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to start event worker")
	}

	// Initialize Dead Letter Processor
	deadLetterProcessor := worker.NewDeadLetterProcessor(retryableBus, cacheRepo)
	if err := deadLetterProcessor.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to start dead letter processor")
	}

	// Setup router with dependencies
	app := router.Setup(router.Dependencies{
		Config:              cfg,
		UserRepo:            userRepo,
		AlertRepo:           alertRepo,
		CacheRepo:           cacheRepo,
		DBHealthCheck:       db,
		WSHub:               wsHub,
		EventBus:            retryableBus,
		EventWorker:         eventWorker,
		DeadLetterProcessor: deadLetterProcessor,
	})

	// Start server in goroutine
	go func() {
		log.Info().Str("address", cfg.Server.Address()).Msg("HTTP server started")
		if err := app.Listen(cfg.Server.Address()); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop workers
	_ = eventWorker.Stop()
	_ = deadLetterProcessor.Stop()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Error during shutdown")
	}

	// Close connections
	closeRedis(redisClient)
	closeDB(db)

	log.Info().Msg("Server stopped")
}

func setupLogger(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.Logging.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if cfg.App.IsDevelopment() {
		log.Logger = log.With().Caller().Logger()
	}
}

func closeDB(db *database.PostgresDB) {
	if err := db.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing database connection")
	}
}

func closeRedis(client *database.RedisClient) {
	if err := client.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing Redis connection")
	}
}
