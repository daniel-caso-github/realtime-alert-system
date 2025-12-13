-- Rollback: Drop alert_rules table

DROP TRIGGER IF EXISTS update_alert_rules_updated_at ON alert_rules;
DROP TABLE IF EXISTS alert_rules;
DROP TYPE IF EXISTS alert_severity;
