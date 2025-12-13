package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
)

// PostgreSQL error codes
// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
	pgErrCheckViolation      = "23514"
	pgErrNotNullViolation    = "23502"
)

// TranslateError converts PostgreSQL-specific errors to domain errors.
// This keeps the domain layer independent of the database implementation.
func TranslateError(err error) error {
	if err == nil {
		return nil
	}

	// Check for "no rows" error
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrNotFound
	}

	// Check for PostgreSQL-specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgErrUniqueViolation:
			return repository.ErrDuplicateKey
		case pgErrForeignKeyViolation:
			return repository.ErrForeignKeyViolation
		case pgErrCheckViolation, pgErrNotNullViolation:
			return repository.ErrInvalidData
		}
	}

	// Check for connection errors
	if isConnectionError(err) {
		return repository.ErrConnection
	}

	// Return original error if no translation found
	return err
}

// isConnectionError checks if the error is related to connection issues.
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	connectionKeywords := []string{
		"connection refused",
		"connection reset",
		"no connection",
		"connection timed out",
		"network is unreachable",
	}

	for _, keyword := range connectionKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}
