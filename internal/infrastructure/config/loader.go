package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load reads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/alerting/")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, will use env vars and defaults
	}

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables explicitly
	bindEnvVars(v)

	// Set defaults
	setDefaults(v)

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func bindEnvVars(v *viper.Viper) {
	// App
	_ = v.BindEnv("app.name", "APP_NAME")
	_ = v.BindEnv("app.env", "APP_ENV")
	_ = v.BindEnv("app.version", "APP_VERSION")

	// Server
	_ = v.BindEnv("server.host", "SERVER_HOST")
	_ = v.BindEnv("server.port", "SERVER_PORT")

	// Database
	_ = v.BindEnv("database.host", "DATABASE_HOST")
	_ = v.BindEnv("database.port", "DATABASE_PORT")
	_ = v.BindEnv("database.user", "DATABASE_USER")
	_ = v.BindEnv("database.password", "DATABASE_PASSWORD")
	_ = v.BindEnv("database.name", "DATABASE_NAME")
	_ = v.BindEnv("database.ssl_mode", "DATABASE_SSL_MODE")

	// Redis
	_ = v.BindEnv("redis.host", "REDIS_HOST")
	_ = v.BindEnv("redis.port", "REDIS_PORT")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = v.BindEnv("redis.db", "REDIS_DB")

	// JWT
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("jwt.expiration", "JWT_EXPIRATION")

	// Logging
	_ = v.BindEnv("logging.level", "LOG_LEVEL")
	_ = v.BindEnv("logging.format", "LOG_FORMAT")
}

func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "realtime-alerting-system")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.version", "1.0.0")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")
	v.SetDefault("server.idle_timeout", "120s")

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "alerting_db")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)

	// JWT defaults
	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.expiration", "15m")
	v.SetDefault("jwt.refresh_expiration", "168h")
	v.SetDefault("jwt.issuer", "realtime-alerting-system")

	// Logging defaults
	v.SetDefault("logging.level", "debug")
	v.SetDefault("logging.format", "console")

	// WebSocket defaults
	v.SetDefault("websocket.read_buffer_size", 1024)
	v.SetDefault("websocket.write_buffer_size", 1024)
	v.SetDefault("websocket.ping_interval", "30s")
	v.SetDefault("websocket.pong_timeout", "60s")

	// Event Bus defaults
	viper.SetDefault("event_bus.consumer_id", "api-server-1")
	viper.SetDefault("event_bus.max_retries", 3)
	viper.SetDefault("event_bus.retry_backoff", "1s")
}
