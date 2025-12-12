// Package router sets up and configures the HTTP router for the application.
package router

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/handler"
)

// Setup configura y retorna una instancia de Fiber con todas las rutas.
func Setup(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: customErrorHandler,
	})

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

	healthHandler := handler.NewHealthHandler(cfg)

	app.Get("/health", healthHandler.Check)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/live", healthHandler.Live)

	v1 := app.Group("/api/v1")

	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	return app
}

// customErrorHandler maneja errores de forma consistente.
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}
