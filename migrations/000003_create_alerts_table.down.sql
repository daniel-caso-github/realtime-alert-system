-- Rollback: Drop alerts table

DROP TRIGGER IF EXISTS update_alerts_updated_at ON alerts;
DROP TABLE IF EXISTS alerts;
DROP TYPE IF EXISTS alert_status;
