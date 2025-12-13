// Package entity contains domain entities for the real-time alerting system.
package entity

import (
	"errors"
)

// ChannelType represents the supported notification channel types.
// It defines the different delivery mechanisms available for sending alerts.
type ChannelType string

// Supported notification channel types.
const (
	// ChannelTypeSlack represents a Slack webhook notification channel.
	ChannelTypeSlack ChannelType = "slack"
	// ChannelTypeEmail represents an email notification channel.
	ChannelTypeEmail ChannelType = "email"
	// ChannelTypeSMS represents an SMS notification channel.
	ChannelTypeSMS ChannelType = "sms"
	// ChannelTypeWebhook represents a generic webhook notification channel.
	ChannelTypeWebhook ChannelType = "webhook"
)

// IsValid checks whether the channel type is a valid supported type.
// Returns true if the type is one of: slack, email, sms, or webhook.
func (t ChannelType) IsValid() bool {
	switch t {
	case ChannelTypeSlack, ChannelTypeEmail, ChannelTypeSMS, ChannelTypeWebhook:
		return true
	default:
		return false
	}
}

// NotificationChannel represents a channel through which notifications are sent.
// It contains the configuration needed to deliver alerts via different mechanisms
// such as Slack, email, SMS, or generic webhooks.
type NotificationChannel struct {
	// ID is the unique identifier for the notification channel.
	ID ID `json:"id" db:"id"`
	// Name is the human-readable name of the channel.
	Name string `json:"name" db:"name"`
	// Type specifies the delivery mechanism (slack, email, sms, webhook).
	Type ChannelType `json:"type" db:"type"`
	// Config holds channel-specific configuration as key-value pairs.
	// Required keys depend on the channel type:
	//   - slack: requires "webhook_url"
	//   - email: requires "recipients"
	Config map[string]interface{} `json:"config" db:"config"`
	// IsEnabled indicates whether the channel is active and can receive notifications.
	IsEnabled bool `json:"is_enabled" db:"is_enabled"`
	// CreatedBy is the optional ID of the user who created this channel.
	CreatedBy *ID `json:"created_by,omitempty" db:"created_by"`
	// Timestamps embeds creation and update timestamps.
	Timestamps
}

// Channel validation errors.
var (
	// ErrChannelNameRequired is returned when the channel name is empty.
	ErrChannelNameRequired = errors.New("channel name is required")
	// ErrChannelNameTooLong is returned when the channel name exceeds 255 characters.
	ErrChannelNameTooLong = errors.New("channel name must be less than 256 characters")
	// ErrChannelInvalidType is returned when the channel type is not supported.
	ErrChannelInvalidType = errors.New("invalid channel type")
	// ErrChannelConfigRequired is returned when the channel config is nil.
	ErrChannelConfigRequired = errors.New("channel config is required")
	// ErrChannelMissingWebhook is returned when a Slack channel lacks webhook_url in config.
	ErrChannelMissingWebhook = errors.New("slack channel requires webhook_url in config")
	// ErrChannelMissingEmail is returned when an email channel lacks recipients in config.
	ErrChannelMissingEmail = errors.New("email channel requires recipients in config")
)

// NewNotificationChannel creates a new notification channel with the given parameters.
// It generates a new unique ID, sets the channel as enabled by default,
// and validates the channel before returning it.
//
// Parameters:
//   - name: the human-readable name for the channel (required, max 255 chars)
//   - channelType: the type of notification channel (must be valid)
//   - config: channel-specific configuration map (required, contents depend on type)
//   - createdBy: optional pointer to the ID of the creating user
//
// Returns the created NotificationChannel and nil on success,
// or nil and an error if validation fails.
func NewNotificationChannel(name string, channelType ChannelType, config map[string]interface{}, createdBy *ID) (*NotificationChannel, error) {
	channel := &NotificationChannel{
		ID:         NewID(),
		Name:       name,
		Type:       channelType,
		Config:     config,
		IsEnabled:  true,
		CreatedBy:  createdBy,
		Timestamps: NewTimestamps(),
	}

	if err := channel.Validate(); err != nil {
		return nil, err
	}

	return channel, nil
}

// Validate checks that the notification channel has valid data.
// It performs the following validations:
//   - Name must not be empty
//   - Name must not exceed 255 characters
//   - Type must be a valid ChannelType
//   - Config must not be nil
//   - Channel-specific validations (e.g., Slack requires webhook_url)
//
// Returns nil if valid, or an appropriate error otherwise.
func (c *NotificationChannel) Validate() error {
	if c.Name == "" {
		return ErrChannelNameRequired
	}

	if len(c.Name) > 255 {
		return ErrChannelNameTooLong
	}

	if !c.Type.IsValid() {
		return ErrChannelInvalidType
	}

	if c.Config == nil {
		return ErrChannelConfigRequired
	}

	// Channel-specific validations
	switch c.Type {
	case ChannelTypeSlack:
		if _, ok := c.Config["webhook_url"]; !ok {
			return ErrChannelMissingWebhook
		}
	case ChannelTypeEmail:
		if _, ok := c.Config["recipients"]; !ok {
			return ErrChannelMissingEmail
		}
	}

	return nil
}

// Enable activates the notification channel, allowing it to receive and send notifications.
// It also updates the channel's timestamp to reflect the modification.
func (c *NotificationChannel) Enable() {
	c.IsEnabled = true
	c.Touch()
}

// Disable deactivates the notification channel, preventing it from sending notifications.
// It also updates the channel's timestamp to reflect the modification.
func (c *NotificationChannel) Disable() {
	c.IsEnabled = false
	c.Touch()
}

// UpdateConfig replaces the channel's configuration with the provided config map.
// It updates the timestamp and validates the new configuration.
// Returns an error if the new configuration is invalid for the channel type.
func (c *NotificationChannel) UpdateConfig(config map[string]interface{}) error {
	c.Config = config
	c.Touch()
	return c.Validate()
}

// GetWebhookURL retrieves the webhook URL from the channel's configuration.
// This method is intended for Slack and webhook channel types.
// Returns the webhook URL as a string, or an empty string if not configured
// or if the value is not a string.
func (c *NotificationChannel) GetWebhookURL() string {
	if url, ok := c.Config["webhook_url"].(string); ok {
		return url
	}
	return ""
}

// GetRecipients retrieves the list of email recipients from the channel's configuration.
// This method is intended for email channel types.
// Returns a slice of recipient email addresses, or nil if not configured
// or if the recipients are not in the expected format.
func (c *NotificationChannel) GetRecipients() []string {
	if recipients, ok := c.Config["recipients"].([]interface{}); ok {
		result := make([]string, 0, len(recipients))
		for _, r := range recipients {
			if s, ok := r.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}
