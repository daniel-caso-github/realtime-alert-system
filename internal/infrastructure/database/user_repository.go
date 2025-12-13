// Package database

package database

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// Ensure PostgresUserRepository implements repository.UserRepository
var _ repository.UserRepository = (*PostgresUserRepository)(nil)

// PostgresUserRepository implements UserRepository using PostgreSQL.
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository.
func NewPostgresUserRepository(db *PostgresDB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db.DB,
	}
}

// Create saves a new user to the database.
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, role, is_active, last_login_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.IsActive,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return TranslateError(err)
}

// GetByID finds a user by their ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id entity.ID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, TranslateError(err)
	}

	return &user, nil
}

// GetByEmail finds a user by their email.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, TranslateError(err)
	}

	return &user, nil
}

// Update updates an existing user.
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, name = $4, role = $5, is_active = $6, last_login_at = $7, updated_at = $8
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.IsActive,
		user.LastLoginAt,
		user.UpdatedAt,
	)
	if err != nil {
		return TranslateError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return TranslateError(err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// Delete removes a user by their ID.
func (r *PostgresUserRepository) Delete(ctx context.Context, id entity.ID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return TranslateError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return TranslateError(err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// List returns paginated users.
func (r *PostgresUserRepository) List(ctx context.Context, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.User], error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, TranslateError(err)
	}

	// Get paginated results
	query := `
		SELECT id, email, password_hash, name, role, is_active, last_login_at, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var users []*entity.User
	if err := r.db.SelectContext(ctx, &users, query, pagination.Limit(), pagination.Offset()); err != nil {
		return nil, TranslateError(err)
	}

	// Return empty slice if nil to avoid null in JSON
	if users == nil {
		users = []*entity.User{}
	}

	result := valueobject.NewPaginatedResult(users, total, pagination)
	return &result, nil
}

// ListByRole returns users filtered by role.
func (r *PostgresUserRepository) ListByRole(ctx context.Context, role entity.UserRole, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.User], error) {
	// Get total count for this role
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE role = $1`
	if err := r.db.GetContext(ctx, &total, countQuery, role); err != nil {
		return nil, TranslateError(err)
	}

	// Get paginated results
	query := `
		SELECT id, email, password_hash, name, role, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE role = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var users []*entity.User
	if err := r.db.SelectContext(ctx, &users, query, role, pagination.Limit(), pagination.Offset()); err != nil {
		return nil, TranslateError(err)
	}

	if users == nil {
		users = []*entity.User{}
	}

	result := valueobject.NewPaginatedResult(users, total, pagination)
	return &result, nil
}

// ExistsByEmail checks if a user with that email exists.
func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, email); err != nil {
		return false, TranslateError(err)
	}

	return exists, nil
}

// Count returns the total number of users.
func (r *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, TranslateError(err)
	}

	return count, nil
}

// CountByRole returns the number of users by role.
func (r *PostgresUserRepository) CountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = $1`

	var count int64
	if err := r.db.GetContext(ctx, &count, query, role); err != nil {
		return 0, TranslateError(err)
	}

	return count, nil
}
