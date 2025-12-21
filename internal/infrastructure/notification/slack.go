// Package notification provides notification implementations.
package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/notification"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/infrastructure/config"
)

// SlackNotifier sends notifications to Slack.
type SlackNotifier struct {
	webhookURL string
	channel    string
	username   string
	enabled    bool
	client     *http.Client
}

// slackMessage represents a Slack message payload.
type slackMessage struct {
	Channel     string            `json:"channel,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Attachments []slackAttachment `json:"attachments"`
}

// slackAttachment represents a Slack message attachment.
type slackAttachment struct {
	Color     string       `json:"color"`
	Title     string       `json:"title"`
	Text      string       `json:"text"`
	Fields    []slackField `json:"fields,omitempty"`
	Footer    string       `json:"footer,omitempty"`
	Timestamp int64        `json:"ts,omitempty"`
}

// slackField represents a field in a Slack attachment.
type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// NewSlackNotifier creates a new Slack notifier.
func NewSlackNotifier(cfg config.SlackConfig, timeout time.Duration) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: cfg.WebhookURL,
		channel:    cfg.Channel,
		username:   cfg.Username,
		enabled:    cfg.Enabled && cfg.WebhookURL != "",
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Send sends a notification to Slack.
func (n *SlackNotifier) Send(ctx context.Context, msg notification.Message) error {
	if !n.enabled {
		log.Debug().Msg("Slack notifications disabled, skipping")
		return nil
	}

	slackMsg := n.buildMessage(msg)

	payload, err := json.Marshal(slackMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned non-200 status: %d", resp.StatusCode)
	}

	log.Debug().
		Str("alert_id", msg.AlertID).
		Str("severity", msg.Severity).
		Msg("Slack notification sent")

	return nil
}

// Name returns the notifier name.
func (n *SlackNotifier) Name() string {
	return "slack"
}

// IsEnabled returns whether the notifier is enabled.
func (n *SlackNotifier) IsEnabled() bool {
	return n.enabled
}

// buildMessage builds a Slack message from a notification message.
func (n *SlackNotifier) buildMessage(msg notification.Message) slackMessage {
	color := n.severityToColor(msg.Severity)
	emoji := n.severityToEmoji(msg.Severity)

	fields := make([]slackField, 0)
	fields = append(fields, slackField{
		Title: "Severity",
		Value: fmt.Sprintf("%s %s", emoji, msg.Severity),
		Short: true,
	})

	if msg.Source != "" {
		fields = append(fields, slackField{
			Title: "Source",
			Value: msg.Source,
			Short: true,
		})
	}

	if msg.AlertID != "" {
		fields = append(fields, slackField{
			Title: "Alert ID",
			Value: msg.AlertID,
			Short: true,
		})
	}

	for key, value := range msg.Fields {
		fields = append(fields, slackField{
			Title: key,
			Value: value,
			Short: true,
		})
	}

	return slackMessage{
		Channel:   n.channel,
		Username:  n.username,
		IconEmoji: ":rotating_light:",
		Attachments: []slackAttachment{
			{
				Color:     color,
				Title:     msg.Title,
				Text:      msg.Text,
				Fields:    fields,
				Footer:    "Real-Time Alerting System",
				Timestamp: time.Now().Unix(),
			},
		},
	}
}

// severityToColor maps severity to Slack attachment color.
func (n *SlackNotifier) severityToColor(severity string) string {
	switch severity {
	case notification.SeverityCritical:
		return "#dc3545" // Red
	case notification.SeverityHigh:
		return "#fd7e14" // Orange
	case notification.SeverityMedium:
		return "#ffc107" // Yellow
	case notification.SeverityLow:
		return "#17a2b8" // Blue
	case notification.SeverityInfo:
		return "#6c757d" // Gray
	default:
		return "#6c757d"
	}
}

// severityToEmoji maps severity to emoji.
func (n *SlackNotifier) severityToEmoji(severity string) string {
	switch severity {
	case notification.SeverityCritical:
		return "🔴"
	case notification.SeverityHigh:
		return "🟠"
	case notification.SeverityMedium:
		return "🟡"
	case notification.SeverityLow:
		return "🔵"
	case notification.SeverityInfo:
		return "⚪"
	default:
		return "⚪"
	}
}

// Compile-time interface verification.
var _ notification.Notifier = (*SlackNotifier)(nil)
