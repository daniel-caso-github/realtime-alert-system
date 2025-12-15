// Package database provides database implementations for repositories.
package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// PostgresAlertRepository implements AlertRepository using PostgreSQL.
type PostgresAlertRepository struct {
	db *sqlx.DB
}

// NewPostgresAlertRepository creates a new PostgreSQL alert repository.
func NewPostgresAlertRepository(db *PostgresDB) *PostgresAlertRepository {
	return &PostgresAlertRepository{
		db: db.DB,
	}
}

// Create inserts a new alert into the database.
func (r *PostgresAlertRepository) Create(ctx context.Context, alert *entity.Alert) error {
	query := `
		INSERT INTO alerts (id, rule_id, title, message, severity, status, source, metadata, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	metadata, err := json.Marshal(alert.Metadata)
	if err != nil {
		return err
	}

	var ruleID *string
	if alert.RuleID != nil {
		id := alert.RuleID.String()
		ruleID = &id
	}

	_, err = r.db.ExecContext(ctx, query,
		alert.ID.String(),
		ruleID,
		alert.Title,
		alert.Message,
		string(alert.Severity),
		string(alert.Status),
		alert.Source,
		metadata,
		alert.ExpiresAt,
		alert.CreatedAt,
		alert.UpdatedAt,
	)

	return TranslateError(err)
}

// GetByID retrieves an alert by its ID.
func (r *PostgresAlertRepository) GetByID(ctx context.Context, id entity.ID) (*entity.Alert, error) {
	query := `SELECT * FROM alerts WHERE id = $1`

	var model AlertModel
	err := r.db.GetContext(ctx, &model, query, id.String())
	if err != nil {
		return nil, TranslateError(err)
	}

	return model.ToEntity()
}

// Update updates an existing alert.
func (r *PostgresAlertRepository) Update(ctx context.Context, alert *entity.Alert) error {
	query := `
		UPDATE alerts
		SET title = $1, message = $2, severity = $3, status = $4, source = $5, metadata = $6,
		    acknowledged_by = $7, acknowledged_at = $8, resolved_by = $9, resolved_at = $10,
		    expires_at = $11, updated_at = $12
		WHERE id = $13
	`

	metadata, err := json.Marshal(alert.Metadata)
	if err != nil {
		return err
	}

	var ackBy, resBy *string
	if alert.AcknowledgedBy != nil {
		id := alert.AcknowledgedBy.String()
		ackBy = &id
	}
	if alert.ResolvedBy != nil {
		id := alert.ResolvedBy.String()
		resBy = &id
	}

	result, err := r.db.ExecContext(ctx, query,
		alert.Title,
		alert.Message,
		string(alert.Severity),
		string(alert.Status),
		alert.Source,
		metadata,
		ackBy,
		alert.AcknowledgedAt,
		resBy,
		alert.ResolvedAt,
		alert.ExpiresAt,
		alert.UpdatedAt,
		alert.ID.String(),
	)

	if err != nil {
		return TranslateError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// Delete removes an alert from the database.
func (r *PostgresAlertRepository) Delete(ctx context.Context, id entity.ID) error {
	query := `DELETE FROM alerts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return TranslateError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// List retrieves alerts with filtering and pagination.
func (r *PostgresAlertRepository) List(
	ctx context.Context,
	filter valueobject.AlertFilter,
	pagination valueobject.Pagination,
) (*valueobject.PaginatedResult[*entity.Alert], error) {
	where, args := r.buildWhereClause(filter)

	countQuery := "SELECT COUNT(*) FROM alerts" + where
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, TranslateError(err)
	}

	query := fmt.Sprintf(`
		SELECT * FROM alerts %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)

	args = append(args, pagination.PageSize(), pagination.Offset())

	var models []AlertModel
	if err := r.db.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, TranslateError(err)
	}

	alerts, err := r.modelsToEntities(models)
	if err != nil {
		return nil, err
	}

	result := valueobject.NewPaginatedResult(alerts, total, pagination)
	return &result, nil
}

// ListByStatus returns alerts filtered by status.
func (r *PostgresAlertRepository) ListByStatus(
	ctx context.Context,
	status entity.AlertStatus,
	pagination valueobject.Pagination,
) (*valueobject.PaginatedResult[*entity.Alert], error) {
	countQuery := `SELECT COUNT(*) FROM alerts WHERE status = $1`
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, string(status)); err != nil {
		return nil, TranslateError(err)
	}

	query := `
		SELECT * FROM alerts
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var models []AlertModel
	if err := r.db.SelectContext(ctx, &models, query, string(status), pagination.PageSize(), pagination.Offset()); err != nil {
		return nil, TranslateError(err)
	}

	alerts, err := r.modelsToEntities(models)
	if err != nil {
		return nil, err
	}

	result := valueobject.NewPaginatedResult(alerts, total, pagination)
	return &result, nil
}

// ListByRuleID returns alerts generated by a specific rule.
func (r *PostgresAlertRepository) ListByRuleID(
	ctx context.Context,
	ruleID entity.ID,
	pagination valueobject.Pagination,
) (*valueobject.PaginatedResult[*entity.Alert], error) {
	countQuery := `SELECT COUNT(*) FROM alerts WHERE rule_id = $1`
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, ruleID.String()); err != nil {
		return nil, TranslateError(err)
	}

	query := `
		SELECT * FROM alerts
		WHERE rule_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var models []AlertModel
	if err := r.db.SelectContext(ctx, &models, query, ruleID.String(), pagination.PageSize(), pagination.Offset()); err != nil {
		return nil, TranslateError(err)
	}

	alerts, err := r.modelsToEntities(models)
	if err != nil {
		return nil, err
	}

	result := valueobject.NewPaginatedResult(alerts, total, pagination)
	return &result, nil
}

// ListActive retrieves all active alerts (for WebSocket broadcast).
func (r *PostgresAlertRepository) ListActive(ctx context.Context) ([]*entity.Alert, error) {
	query := `SELECT * FROM alerts WHERE status = 'active' ORDER BY severity, created_at DESC`

	var models []AlertModel
	if err := r.db.SelectContext(ctx, &models, query); err != nil {
		return nil, TranslateError(err)
	}

	return r.modelsToEntities(models)
}

// ListExpired retrieves alerts that have expired but not marked as such.
func (r *PostgresAlertRepository) ListExpired(ctx context.Context) ([]*entity.Alert, error) {
	query := `
		SELECT * FROM alerts
		WHERE status NOT IN ('resolved', 'expired')
		AND expires_at IS NOT NULL
		AND expires_at < NOW()
	`

	var models []AlertModel
	if err := r.db.SelectContext(ctx, &models, query); err != nil {
		return nil, TranslateError(err)
	}

	return r.modelsToEntities(models)
}

// Count returns the total number of alerts.
func (r *PostgresAlertRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM alerts`
	var count int64
	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, TranslateError(err)
	}
	return count, nil
}

// CountByStatus returns the number of alerts by status.
func (r *PostgresAlertRepository) CountByStatus(ctx context.Context, status entity.AlertStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM alerts WHERE status = $1`
	var count int64
	if err := r.db.GetContext(ctx, &count, query, string(status)); err != nil {
		return 0, TranslateError(err)
	}
	return count, nil
}

// CountBySeverity returns the number of alerts by severity.
func (r *PostgresAlertRepository) CountBySeverity(ctx context.Context, severity entity.AlertSeverity) (int64, error) {
	query := `SELECT COUNT(*) FROM alerts WHERE severity = $1`
	var count int64
	if err := r.db.GetContext(ctx, &count, query, string(severity)); err != nil {
		return 0, TranslateError(err)
	}
	return count, nil
}

// GetStatistics retrieves alert statistics.
func (r *PostgresAlertRepository) GetStatistics(ctx context.Context) (*repository.AlertStatistics, error) {
	query := `
		SELECT
			COUNT(*) as total_alerts,
			COUNT(*) FILTER (WHERE status = 'active') as active_alerts,
			COUNT(*) FILTER (WHERE status = 'acknowledged') as acknowledged_alerts,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved_alerts
		FROM alerts
	`

	var stats repository.AlertStatistics
	if err := r.db.GetContext(ctx, &stats, query); err != nil {
		return nil, TranslateError(err)
	}

	// Get by severity
	severityQuery := `SELECT severity, COUNT(*) as count FROM alerts GROUP BY severity`
	rows, err := r.db.QueryContext(ctx, severityQuery)
	if err != nil {
		return nil, TranslateError(err)
	}
	defer func() { _ = rows.Close() }()

	stats.BySeverity = make(map[string]int64)
	for rows.Next() {
		var severity string
		var count int64
		if err := rows.Scan(&severity, &count); err != nil {
			return nil, err
		}
		stats.BySeverity[severity] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get by source
	sourceQuery := `SELECT source, COUNT(*) as count FROM alerts WHERE source != '' GROUP BY source`
	rows, err = r.db.QueryContext(ctx, sourceQuery)
	if err != nil {
		return nil, TranslateError(err)
	}
	defer func() { _ = rows.Close() }()

	stats.BySource = make(map[string]int64)
	for rows.Next() {
		var source string
		var count int64
		if err := rows.Scan(&source, &count); err != nil {
			return nil, err
		}
		stats.BySource[source] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &stats, nil
}

// buildWhereClause builds the WHERE clause for filtering alerts.
func (r *PostgresAlertRepository) buildWhereClause(filter valueobject.AlertFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if len(filter.Statuses) > 0 {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, string(status))
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(filter.Severities) > 0 {
		placeholders := make([]string, len(filter.Severities))
		for i, severity := range filter.Severities {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, string(severity))
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("severity IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.Source != nil {
		conditions = append(conditions, fmt.Sprintf("source = $%d", argIndex))
		args = append(args, *filter.Source)
		argIndex++
	}

	if filter.Search != nil && *filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR message ILIKE $%d)", argIndex, argIndex+1))
		searchTerm := "%" + *filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	if filter.FromDate != nil && filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, filter.FromDate, filter.ToDate)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

// modelsToEntities converts a slice of AlertModel to a slice of entity.Alert.
func (r *PostgresAlertRepository) modelsToEntities(models []AlertModel) ([]*entity.Alert, error) {
	alerts := make([]*entity.Alert, 0, len(models))
	for _, model := range models {
		alert, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

// Compile-time interface verification
var _ repository.AlertRepository = (*PostgresAlertRepository)(nil)
