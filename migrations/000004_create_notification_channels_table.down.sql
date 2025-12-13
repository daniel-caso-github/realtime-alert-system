-- Rollback: Drop notification_channels table

DROP TRIGGER IF EXISTS update_notification_channels_updated_at ON notification_channels;
DROP TABLE IF EXISTS notification_channels;
DROP TYPE IF EXISTS channel_type;
