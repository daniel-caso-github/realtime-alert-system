// Package database provides PostgreSQL-backed implementations of repository interfaces.
package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/entity"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/repository"
	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

// Ensure PostgresAlertRepository implements repository.AlertRepository
var _ repository.AlertRepository = (*PostgresAlertRepository)(nil)

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

// Create saves a new alert to the database.
func (r *PostgresAlertRepository) Create(ctx context.Context, alert *entity.Alert) error {
	query := `
		INSERT INTO alerts (
			id, rule_id, title, message, severity, status, source, metadata,
			acknowledged_by, acknowledged_at, resolved_by, resolved_at, expires_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, query,
		alert.ID,
		alert.RuleID,
		alert.Title,
		alert.Message,
		alert.Severity,
		alert.Status,
		alert.Source,
		alert.Metadata,
		alert.AcknowledgedBy,
		alert.AcknowledgedAt,
		alert.ResolvedBy,
		alert.ResolvedAt,
		alert.ExpiresAt,
		alert.CreatedAt,
		alert.UpdatedAt,
	)

	return TranslateError(err)
}

// GetByID finds an alert by its ID.
func (r *PostgresAlertRepository) GetByID(ctx context.Context, id entity.ID) (*entity.Alert, error) {
	query := `
		SELECT id, rule_id, title, message, severity, status, source, metadata,
			   acknowledged_by, acknowledged_at, resolved_by, resolved_at, expires_at,
			   created_at, updated_at
		FROM alerts
		WHERE id = $1
	`

	var alert entity.Alert
	err := r.db.GetContext(ctx, &alert, query, id)
	if err != nil {
		return nil, TranslateError(err)
	}

	return &alert, nil
}

// Update updates an existing alert.
func (r *PostgresAlertRepository) Update(ctx context.Context, alert *entity.Alert) error {
	query := `
		UPDATE alerts
		SET rule_id = $2, title = $3, message = $4, severity = $5, status = $6,
			source = $7, metadata = $8, acknowledged_by = $9, acknowledged_at = $10,
			resolved_by = $11, resolved_at = $12, expires_at = $13, updated_at = $14
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		alert.ID,
		alert.RuleID,
		alert.Title,
		alert.Message,
		alert.Severity,
		alert.Status,
		alert.Source,
		alert.Metadata,
		alert.AcknowledgedBy,
		alert.AcknowledgedAt,
		alert.ResolvedBy,
		alert.ResolvedAt,
		alert.ExpiresAt,
		alert.UpdatedAt,
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

// Delete removes an alert by its ID.
func (r *PostgresAlertRepository) Delete(ctx context.Context, id entity.ID) error {
	query := `DELETE FROM alerts WHERE id = $1`

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

// List returns paginated alerts with optional filters.
func (r *PostgresAlertRepository) List(ctx context.Context, filter valueobject.AlertFilter, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.Alert], error) {
	// Build WHERE clause dynamically
	whereClause, args := r.buildWhereClause(filter)

	// Get total count
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM alerts %s`, whereClause)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, TranslateError(err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, rule_id, title, message, severity, status, source, metadata,
			   acknowledged_by, acknowledged_at, resolved_by, resolved_at, expires_at,
			   created_at, updated_at
		FROM alerts
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pagination.Limit(), pagination.Offset())

	var alerts []*entity.Alert
	if err := r.db.SelectContext(ctx, &alerts, query, args...); err != nil {
		return nil, TranslateError(err)
	}

	if alerts == nil {
		alerts = []*entity.Alert{}
	}

	result := valueobject.NewPaginatedResult(alerts, total, pagination)
	return &result, nil
}

// buildWhereClause constructs the WHERE clause based on filters.
func (r *PostgresAlertRepository) buildWhereClause(filter valueobject.AlertFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.HasStatusFilter() {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ", ")))
	}

	if filter.HasSeverityFilter() {
		placeholders := make([]string, len(filter.Severities))
		for i, severity := range filter.Severities {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, severity)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("severity IN (%s)", strings.Join(placeholders, ", ")))
	}

	if filter.Source != nil {
		conditions = append(conditions, fmt.Sprintf("source = $%d", argIndex))
		args = append(args, *filter.Source)
		argIndex++
	}

	if filter.RuleID != nil {
		conditions = append(conditions, fmt.Sprintf("rule_id = $%d", argIndex))
		args = append(args, *filter.RuleID)
		argIndex++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.FromDate)
		argIndex++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.ToDate)
		argIndex++
	}

	if filter.HasSearch() {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR message ILIKE $%d)", argIndex, argIndex+1))
		searchPattern := "%" + *filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

// ListByStatus returns alerts filtered by status.
func (r *PostgresAlertRepository) ListByStatus(ctx context.Context, status entity.AlertStatus, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.Alert], error) {
	filter := valueobject.NewAlertFilter().WithStatuses(status)
	return r.List(ctx, filter, pagination)
}

// ListByRuleID returns alerts generated by a specific rule.
func (r *PostgresAlertRepository) ListByRuleID(ctx context.Context, ruleID entity.ID, pagination valueobject.Pagination) (*valueobject.PaginatedResult[*entity.Alert], error) {
	filter := valueobject.NewAlertFilter().WithRuleID(ruleID)
	return r.List(ctx, filter, pagination)
}

// ListActive returns all active alerts.
func (r *PostgresAlertRepository) ListActive(ctx context.Context) ([]*entity.Alert, error) {
	query := `
		SELECT id, rule_id, title, message, severity, status, source, metadata,
			   acknowledged_by, acknowledged_at, resolved_by, resolved_at, expires_at,
			   created_at, updated_at
		FROM alerts
		WHERE status = $1
		ORDER BY severity ASC, created_at DESC
	`

	var alerts []*entity.Alert
	if err := r.db.SelectContext(ctx, &alerts, query, entity.AlertStatusActive); err != nil {
		return nil, TranslateError(err)
	}

	if alerts == nil {
		alerts = []*entity.Alert{}
	}

	return alerts, nil
}

// ListExpired returns alerts that have expired but are still active.
func (r *PostgresAlertRepository) ListExpired(ctx context.Context) ([]*entity.Alert, error) {
	query := `
		SELECT id, rule_id, title, message, severity, status, source, metadata,
			   acknowledged_by, acknowledged_at, resolved_by, resolved_at, expires_at,
			   created_at, updated_at
		FROM alerts
		WHERE status = $1 AND expires_at IS NOT NULL AND expires_at < NOW()
	`

	var alerts []*entity.Alert
	if err := r.db.SelectContext(ctx, &alerts, query, entity.AlertStatusActive); err != nil {
		return nil, TranslateError(err)
	}

	if alerts == nil {
		alerts = []*entity.Alert{}
	}

	return alerts, nil
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
	if err := r.db.GetContext(ctx, &count, query, status); err != nil {
		return 0, TranslateError(err)
	}

	return count, nil
}

// CountBySeverity returns the number of alerts by severity.
func (r *PostgresAlertRepository) CountBySeverity(ctx context.Context, severity entity.AlertSeverity) (int64, error) {
	query := `SELECT COUNT(*) FROM alerts WHERE severity = $1`

	var count int64
	if err := r.db.GetContext(ctx, &count, query, severity); err != nil {
		return 0, TranslateError(err)
	}

	return count, nil
}

// GetStatistics returns aggregated alert statistics.
func (r *PostgresAlertRepository) GetStatistics(ctx context.Context) (*repository.AlertStatistics, error) {
	stats := &repository.AlertStatistics{
		BySeverity: make(map[string]int64),
		BySource:   make(map[string]int64),
	}

	// Get total and status counts in one query
	statusQuery := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'active') as active,
			COUNT(*) FILTER (WHERE status = 'acknowledged') as acknowledged,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved
		FROM alerts
	`

	var statusStats struct {
		Total        int64 `db:"total"`
		Active       int64 `db:"active"`
		Acknowledged int64 `db:"acknowledged"`
		Resolved     int64 `db:"resolved"`
	}

	if err := r.db.GetContext(ctx, &statusStats, statusQuery); err != nil {
		return nil, TranslateError(err)
	}

	stats.TotalAlerts = statusStats.Total
	stats.ActiveAlerts = statusStats.Active
	stats.AcknowledgedAlerts = statusStats.Acknowledged
	stats.ResolvedAlerts = statusStats.Resolved

	// Get counts by severity
	severityQuery := `SELECT severity, COUNT(*) as count FROM alerts GROUP BY severity`
	var severityCounts []struct {
		Severity string `db:"severity"`
		Count    int64  `db:"count"`
	}

	if err := r.db.SelectContext(ctx, &severityCounts, severityQuery); err != nil {
		return nil, TranslateError(err)
	}

	for _, sc := range severityCounts {
		stats.BySeverity[sc.Severity] = sc.Count
	}

	// Get counts by source (top 10)
	sourceQuery := `
		SELECT source, COUNT(*) as count
		FROM alerts
		WHERE source IS NOT NULL AND source != ''
		GROUP BY source
		ORDER BY count DESC
		LIMIT 10
	`
	var sourceCounts []struct {
		Source string `db:"source"`
		Count  int64  `db:"count"`
	}

	if err := r.db.SelectContext(ctx, &sourceCounts, sourceQuery); err != nil {
		return nil, TranslateError(err)
	}

	for _, sc := range sourceCounts {
		stats.BySource[sc.Source] = sc.Count
	}

	return stats, nil
}
