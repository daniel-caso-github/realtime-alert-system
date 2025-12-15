package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/application/dto"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/database"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/http/router"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/presentation/websocket"
)

// TestApp holds the test application and its dependencies.
type TestApp struct {
	App       *fiber.App
	Config    *config.Config
	DB        *database.PostgresDB
	Redis     *database.RedisClient
	UserRepo  repository.UserRepository
	AlertRepo repository.AlertRepository
	CacheRepo repository.CacheRepository
}

// SetupTestApp creates a test application with real database connections.
func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()

	// Load test configuration
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "test-app",
			Env:     "test",
			Version: "1.0.0",
		},
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8081,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "alerting_db",
			SSLMode:  "disable",
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       1, // Use different DB for tests
		},
		JWT: config.JWTConfig{
			Secret:            "test-secret-key-for-testing-only",
			Expiration:        15 * time.Minute,
			RefreshExpiration: 24 * time.Hour,
			Issuer:            "test-app",
		},
	}

	// Connect to database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}

	// Connect to Redis
	redis, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		_ = db.Close()
		t.Skipf("Skipping integration test: %v", err)
	}

	// Clear rate limiting keys before each test
	clearRateLimiting(redis)

	// Create repositories
	userRepo := database.NewPostgresUserRepository(db)
	alertRepo := database.NewPostgresAlertRepository(db)
	cacheRepo := database.NewRedisCacheRepository(redis)

	// Create WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Setup router
	app := router.Setup(router.Dependencies{
		Config:        cfg,
		UserRepo:      userRepo,
		AlertRepo:     alertRepo,
		CacheRepo:     cacheRepo,
		DBHealthCheck: db,
		WSHub:         wsHub,
	})

	return &TestApp{
		App:       app,
		Config:    cfg,
		DB:        db,
		Redis:     redis,
		UserRepo:  userRepo,
		AlertRepo: alertRepo,
		CacheRepo: cacheRepo,
	}
}

// clearRateLimiting clears all rate limiting keys from Redis.
func clearRateLimiting(redis *database.RedisClient) {
	ctx := context.Background()
	// FlushDB clears all keys in the test database (DB 1)
	_ = redis.FlushDB(ctx)
}

// Cleanup cleans up test resources.
func (ta *TestApp) Cleanup(t *testing.T) {
	t.Helper()

	// Clear test data
	ctx := context.Background()
	_, _ = ta.DB.ExecContext(ctx, "DELETE FROM alerts WHERE title LIKE 'Test%'")
	_, _ = ta.DB.ExecContext(ctx, "DELETE FROM users WHERE email LIKE 'test%'")

	// Clear rate limiting
	clearRateLimiting(ta.Redis)

	// Close connections
	_ = ta.Redis.Close()
	_ = ta.DB.Close()
}

// MakeRequest makes an HTTP request to the test app.
func (ta *TestApp) MakeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, _ := ta.App.Test(req, -1)
	defer func() { _ = resp.Body.Close() }()

	// Convert to ResponseRecorder for compatibility
	recorder := httptest.NewRecorder()
	recorder.Code = resp.StatusCode

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	recorder.Body = buf

	return recorder
}

// Login logs in a user and returns the access token.
func (ta *TestApp) Login(t *testing.T, email, password string) string {
	t.Helper()

	resp := ta.MakeRequest("POST", "/api/v1/auth/login", dto.LoginRequest{
		Email:    email,
		Password: password,
	}, "")

	if resp.Code != 200 {
		t.Fatalf("Login failed with status %d: %s", resp.Code, resp.Body.String())
	}

	var loginResp dto.LoginResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	return loginResp.AccessToken
}
