package database

import (
	"time"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
)

// AlertModel represents the database model for alerts.
type AlertModel struct {
	ID             string     `db:"id"`
	RuleID         *string    `db:"rule_id"`
	Title          string     `db:"title"`
	Message        string     `db:"message"`
	Severity       string     `db:"severity"`
	Status         string     `db:"status"`
	Source         string     `db:"source"`
	Metadata       JSONMap    `db:"metadata"`
	AcknowledgedBy *string    `db:"acknowledged_by"`
	AcknowledgedAt *time.Time `db:"acknowledged_at"`
	ResolvedBy     *string    `db:"resolved_by"`
	ResolvedAt     *time.Time `db:"resolved_at"`
	ExpiresAt      *time.Time `db:"expires_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// ToEntity converts the database model to a domain entity.
func (m *AlertModel) ToEntity() (*entity.Alert, error) {
	id, err := entity.ParseID(m.ID)
	if err != nil {
		return nil, err
	}

	alert := &entity.Alert{
		ID:             id,
		Title:          m.Title,
		Message:        m.Message,
		Severity:       entity.AlertSeverity(m.Severity),
		Status:         entity.AlertStatus(m.Status),
		Source:         m.Source,
		Metadata:       m.Metadata,
		AcknowledgedAt: m.AcknowledgedAt,
		ResolvedAt:     m.ResolvedAt,
		ExpiresAt:      m.ExpiresAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if m.RuleID != nil {
		ruleID, err := entity.ParseID(*m.RuleID)
		if err != nil {
			return nil, err
		}
		alert.RuleID = &ruleID
	}

	if m.AcknowledgedBy != nil {
		ackBy, err := entity.ParseID(*m.AcknowledgedBy)
		if err != nil {
			return nil, err
		}
		alert.AcknowledgedBy = &ackBy
	}

	if m.ResolvedBy != nil {
		resBy, err := entity.ParseID(*m.ResolvedBy)
		if err != nil {
			return nil, err
		}
		alert.ResolvedBy = &resBy
	}

	return alert, nil
}
