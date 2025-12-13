package repository

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// UserRepository defines the persistence operations for users.
// This interface is implemented by the infrastructure layer (PostgreSQL).
type UserRepository interface {
	// Create saves a new user.
	// Returns ErrDuplicateKey if the email already exists.
	Create(ctx context.Context, user *entity.User) error

	// GetByID finds a user by their ID.
	// Returns ErrNotFound if it doesn't exist.
	GetByID(ctx context.Context, id entity.ID) (*entity.User, error)

	// GetByEmail finds a user by their email.
	// Returns ErrNotFound if it doesn't exist.
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates an existing user.
	// Returns ErrNotFound if it doesn't exist.
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user by their ID.
	// Returns ErrNotFound if it doesn't exist.
	Delete(ctx context.Context, id entity.ID) error

	// List returns paginated users.
	List(ctx context.Context, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.User], error)

	// ListByRole returns users filtered by role.
	ListByRole(ctx context.Context, role entity.UserRole, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.User], error)

	// ExistsByEmail checks if a user with that email exists.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Count returns the total number of users.
	Count(ctx context.Context) (int64, error)

	// CountByRole returns the number of users by role.
	CountByRole(ctx context.Context, role entity.UserRole) (int64, error)
}
