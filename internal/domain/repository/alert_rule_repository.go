// Package repository provides interfaces for data persistence operations related to alert rules.
package repository

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// AlertRuleRepository defines the persistence operations for alert rules.
type AlertRuleRepository interface {
	// Create saves a new rule.
	Create(ctx context.Context, rule *entity.AlertRule) error

	// GetByID finds a rule by its ID.
	// Returns ErrNotFound if it doesn't exist.
	GetByID(ctx context.Context, id entity.ID) (*entity.AlertRule, error)

	// Update updates an existing rule.
	// Returns ErrNotFound if it doesn't exist.
	Update(ctx context.Context, rule *entity.AlertRule) error

	// Delete removes a rule by its ID.
	// Returns ErrNotFound if it doesn't exist.
	Delete(ctx context.Context, id entity.ID) error

	// List returns paginated rules.
	List(ctx context.Context, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.AlertRule], error)

	// ListEnabled returns only enabled rules.
	// Useful for the rule evaluation engine.
	ListEnabled(ctx context.Context) ([]*entity.AlertRule, error)

	// ListByCreator returns rules created by a specific user.
	ListByCreator(ctx context.Context, userID entity.ID, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.AlertRule], error)

	// ExistsByName checks if a rule with that name exists.
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Count returns the total number of rules.
	Count(ctx context.Context) (int64, error)

	// CountEnabled returns the number of enabled rules.
	CountEnabled(ctx context.Context) (int64, error)
}
