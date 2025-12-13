package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// Cache TTL constants
const (
	userCacheTTL = 15 * time.Minute
)

// Ensure CachedUserRepository implements repository.UserRepository
var _ repository.UserRepository = (*CachedUserRepository)(nil)

// CachedUserRepository wraps PostgresUserRepository with Redis caching.
type CachedUserRepository struct {
	postgres *PostgresUserRepository
	cache    repository.CacheRepository
	keys     *CacheKey
}

// NewCachedUserRepository creates a new cached user repository.
func NewCachedUserRepository(postgres *PostgresUserRepository, cache repository.CacheRepository) *CachedUserRepository {
	return &CachedUserRepository{
		postgres: postgres,
		cache:    cache,
		keys:     NewCacheKey(),
	}
}

// Create saves a new user (no caching on create).
func (r *CachedUserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.postgres.Create(ctx, user)
}

// GetByID finds a user by ID, using cache when available.
func (r *CachedUserRepository) GetByID(ctx context.Context, id entity.ID) (*entity.User, error) {
	cacheKey := r.keys.User(id)

	// Try cache first
	var user entity.User
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// Cache miss - get from database
	dbUser, err := r.postgres.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache (ignore errors - cache is optional)
	if cacheErr := r.cache.Set(ctx, cacheKey, dbUser, userCacheTTL); cacheErr != nil {
		log.Warn().Err(cacheErr).Str("key", cacheKey).Msg("Failed to cache user")
	}

	return dbUser, nil
}

// GetByEmail finds a user by email, using cache when available.
func (r *CachedUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	cacheKey := r.keys.UserByEmail(email)

	// Try cache first
	var user entity.User
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// Cache miss - get from database
	dbUser, err := r.postgres.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if cacheErr := r.cache.Set(ctx, cacheKey, dbUser, userCacheTTL); cacheErr != nil {
		log.Warn().Err(cacheErr).Str("key", cacheKey).Msg("Failed to cache user")
	}

	return dbUser, nil
}

// Update updates a user and invalidates cache.
func (r *CachedUserRepository) Update(ctx context.Context, user *entity.User) error {
	// Update in database first
	if err := r.postgres.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate cache
	r.invalidateUserCache(ctx, user)

	return nil
}

// Delete removes a user and invalidates cache.
func (r *CachedUserRepository) Delete(ctx context.Context, id entity.ID) error {
	// Get user first to invalidate email cache
	user, err := r.postgres.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from database
	if err := r.postgres.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	r.invalidateUserCache(ctx, user)

	return nil
}

// invalidateUserCache removes user from all cache keys.
func (r *CachedUserRepository) invalidateUserCache(ctx context.Context, user *entity.User) {
	// Delete by ID
	if err := r.cache.Delete(ctx, r.keys.User(user.ID)); err != nil {
		log.Warn().Err(err).Msg("Failed to invalidate user cache by ID")
	}

	// Delete by email
	if err := r.cache.Delete(ctx, r.keys.UserByEmail(user.Email)); err != nil {
		log.Warn().Err(err).Msg("Failed to invalidate user cache by email")
	}
}

// List returns paginated users (not cached - lists change frequently).
func (r *CachedUserRepository) List(ctx context.Context, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.User], error) {
	return r.postgres.List(ctx, pagination)
}

// ListByRole returns users by role (not cached).
func (r *CachedUserRepository) ListByRole(ctx context.Context, role entity.UserRole, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.User], error) {
	return r.postgres.ListByRole(ctx, role, pagination)
}

// ExistsByEmail checks if email exists (not cached - needs to be real-time).
func (r *CachedUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.postgres.ExistsByEmail(ctx, email)
}

// Count returns total users (not cached).
func (r *CachedUserRepository) Count(ctx context.Context) (int64, error) {
	return r.postgres.Count(ctx)
}

// CountByRole returns users count by role (not cached).
func (r *CachedUserRepository) CountByRole(ctx context.Context, role entity.UserRole) (int64, error) {
	return r.postgres.CountByRole(ctx, role)
}
