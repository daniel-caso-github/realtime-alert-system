package router

import (
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
	// Crear la aplicación Fiber con configuración personalizada
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		// En producción, no queremos mostrar errores detallados
		ErrorHandler: customErrorHandler,
	})

	// =========================================================================
	// MIDDLEWARES GLOBALES
	// =========================================================================
	// Se ejecutan en CADA request, en el orden que se registran.

	// Recover: Captura panics y evita que la app se caiga
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.App.IsDevelopment(),
	}))

	// Request ID: Genera un ID único para cada request (útil para tracing)
	app.Use(requestid.New())

	// Logger: Registra cada request en los logs
	if cfg.App.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
		}))
	}

	// CORS: Permite requests desde otros dominios (necesario para frontends)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // En producción, especificar dominios permitidos
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// =========================================================================
	// HEALTH CHECK ROUTES (sin autenticación)
	// =========================================================================
	healthHandler := handler.NewHealthHandler(cfg)

	app.Get("/health", healthHandler.Check) // Estado completo
	app.Get("/ready", healthHandler.Ready)  // Readiness probe
	app.Get("/live", healthHandler.Live)    // Liveness probe

	// =========================================================================
	// API V1 ROUTES
	// =========================================================================
	// Agrupamos las rutas bajo /api/v1 para versionamiento
	v1 := app.Group("/api/v1")

	// TODO: Agregar rutas de la API en las siguientes fases
	// v1.Get("/alerts", alertHandler.List)
	// v1.Post("/alerts", alertHandler.Create)
	// etc.

	// Ruta de prueba para verificar que la API funciona
	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	return app
}

// customErrorHandler maneja errores de forma consistente.
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Código de error por defecto
	code := fiber.StatusInternalServerError

	// Si es un error de Fiber, usar su código
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Retornar error en formato JSON
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
