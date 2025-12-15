// Package router configures HTTP routes and middleware.
package router

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberws "github.com/gofiber/websocket/v2"
	swagger "github.com/swaggo/fiber-swagger"

	_ "github.com/daniel-caso-github/realtime-alerting-system/docs" // Swagger docs

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/service"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/handler"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/middleware"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/websocket"
)

// Dependencies holds all dependencies needed by the router.
type Dependencies struct {
	Config        *config.Config
	UserRepo      repository.UserRepository
	AlertRepo     repository.AlertRepository
	CacheRepo     repository.CacheRepository
	DBHealthCheck handler.HealthChecker
	WSHub         *websocket.Hub
}

// Setup configures and returns a Fiber app with all routes.
func Setup(deps Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      deps.Config.App.Name,
		ReadTimeout:  deps.Config.Server.ReadTimeout,
		WriteTimeout: deps.Config.Server.WriteTimeout,
		IdleTimeout:  deps.Config.Server.IdleTimeout,
		ErrorHandler: customErrorHandler,
	})

	setupMiddleware(app, deps.Config)

	// Create publisher for WebSocket events
	alertPublisher := websocket.NewAlertPublisher(deps.WSHub)

	// Create services
	authService := service.NewAuthService(deps.UserRepo, deps.CacheRepo, &deps.Config.JWT)
	alertService := service.NewAlertService(deps.AlertRepo, deps.CacheRepo, alertPublisher)

	// Create handlers
	healthHandler := handler.NewHealthHandler(deps.Config, deps.DBHealthCheck, deps.CacheRepo, deps.WSHub)
	authHandler := handler.NewAuthHandler(authService)
	alertHandler := handler.NewAlertHandler(alertService)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	apiRateLimiter := middleware.APIRateLimiter(deps.CacheRepo)
	loginRateLimiter := middleware.LoginRateLimiter(deps.CacheRepo)

	// WebSocket handler
	wsHandler := websocket.NewHandler(deps.WSHub)

	// Health routes (no auth required)
	app.Get("/health", healthHandler.Check)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/live", healthHandler.Live)

	// API v1 routes
	v1 := app.Group("/api/v1")
	v1.Use(apiRateLimiter.Limit())

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/login", loginRateLimiter.LimitByEndpoint(), authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", authMiddleware.Authenticate, authHandler.Me)

	// Alert routes (protected)
	alerts := v1.Group("/alerts", authMiddleware.Authenticate)
	alerts.Get("/", alertHandler.List)
	alerts.Get("/statistics", alertHandler.GetStatistics)
	alerts.Post("/", middleware.RequireOperator(), alertHandler.Create)
	alerts.Get("/:id", alertHandler.GetByID)
	alerts.Post("/:id/acknowledge", middleware.RequireOperator(), alertHandler.Acknowledge)
	alerts.Post("/:id/resolve", middleware.RequireOperator(), alertHandler.Resolve)
	alerts.Delete("/:id", middleware.RequireAdmin(), alertHandler.Delete)

	// WebSocket route
	app.Use("/ws", wsHandler.Upgrade)
	app.Get("/ws", authMiddleware.OptionalAuth, fiberws.New(wsHandler.Handle))

	// Swagger documentation
	app.Get("/swagger/*", swagger.WrapHandler)

	return app
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.App.IsDevelopment(),
	}))

	app.Use(requestid.New())

	if cfg.App.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
