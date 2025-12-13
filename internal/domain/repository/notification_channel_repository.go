package repository

import (
	"context"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// NotificationChannelRepository defines the persistence operations for notification channels.
type NotificationChannelRepository interface {
	// Create saves a new channel.
	Create(ctx context.Context, channel *entity.NotificationChannel) error

	// GetByID finds a channel by its ID.
	// Returns ErrNotFound if it doesn't exist.
	GetByID(ctx context.Context, id entity.ID) (*entity.NotificationChannel, error)

	// Update updates an existing channel.
	// Returns ErrNotFound if it doesn't exist.
	Update(ctx context.Context, channel *entity.NotificationChannel) error

	// Delete removes a channel by its ID.
	// Returns ErrNotFound if it doesn't exist.
	Delete(ctx context.Context, id entity.ID) error

	// List returns paginated channels.
	List(ctx context.Context, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.NotificationChannel], error)

	// ListEnabled returns only enabled channels.
	ListEnabled(ctx context.Context) ([]*entity.NotificationChannel, error)

	// ListByType returns channels filtered by type.
	ListByType(ctx context.Context, channelType entity.ChannelType, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.NotificationChannel], error)

	// GetChannelsForRule returns the channels associated with a rule.
	GetChannelsForRule(ctx context.Context, ruleID entity.ID) ([]*entity.NotificationChannel, error)

	// AssociateWithRule associates a channel with a rule.
	AssociateWithRule(ctx context.Context, channelID, ruleID entity.ID) error

	// DisassociateFromRule removes the association between a channel and a rule.
	DisassociateFromRule(ctx context.Context, channelID, ruleID entity.ID) error

	// Count returns the total number of channels.
	Count(ctx context.Context) (int64, error)
}
