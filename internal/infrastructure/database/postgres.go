package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/jmoiron/sqlx"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
)

// PostgresDB wraps the sqlx.DB connection with additional functionality.
type PostgresDB struct {
	*sqlx.DB
	config *config.DatabaseConfig
}

// NewPostgresDB creates a new PostgreSQL connection.
// It configures connection pooling and verifies connectivity.
func NewPostgresDB(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	// Open connection using pgx driver
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return &PostgresDB{
		DB:     db,
		config: cfg,
	}, nil
}

// Health checks if the database connection is healthy.
func (p *PostgresDB) Health(ctx context.Context) error {
	return p.PingContext(ctx)
}

// Close closes the database connection.
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
