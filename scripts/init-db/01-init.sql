-- ============================================================================
-- Database Initialization Script
-- ============================================================================

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- ENUM TYPES
-- ============================================================================

-- Alert severity levels
CREATE TYPE alert_severity AS ENUM ('critical', 'high', 'medium', 'low', 'info');

-- Alert status
CREATE TYPE alert_status AS ENUM ('active', 'acknowledged', 'resolved', 'expired');

-- Notification channel types
CREATE TYPE channel_type AS ENUM ('slack', 'email', 'sms', 'webhook');

-- User roles
CREATE TYPE user_role AS ENUM ('admin', 'operator', 'viewer');

-- ============================================================================
-- TABLES
-- ============================================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'viewer',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Alert rules table
CREATE TABLE IF NOT EXISTS alert_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    condition JSONB NOT NULL,
    severity alert_severity NOT NULL DEFAULT 'medium',
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    cooldown_minutes INTEGER NOT NULL DEFAULT 5,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_id UUID REFERENCES alert_rules(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    severity alert_severity NOT NULL,
    status alert_status NOT NULL DEFAULT 'active',
    source VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Notification channels table
CREATE TABLE IF NOT EXISTS notification_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type channel_type NOT NULL,
    config JSONB NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Alert-Channel mapping (many-to-many)
CREATE TABLE IF NOT EXISTS alert_rule_channels (
    rule_id UUID NOT NULL REFERENCES alert_rules(id) ON DELETE CASCADE,
    channel_id UUID NOT NULL REFERENCES notification_channels(id) ON DELETE CASCADE,
    PRIMARY KEY (rule_id, channel_id)
);

-- Notification history table
CREATE TABLE IF NOT EXISTS notification_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    alert_id UUID NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    channel_id UUID NOT NULL REFERENCES notification_channels(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    response JSONB,
    sent_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Users indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Alerts indexes
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_rule_id ON alerts(rule_id);

-- Alert rules indexes
CREATE INDEX idx_alert_rules_is_enabled ON alert_rules(is_enabled);
CREATE INDEX idx_alert_rules_severity ON alert_rules(severity);

-- Notification history indexes
CREATE INDEX idx_notification_history_alert_id ON notification_history(alert_id);
CREATE INDEX idx_notification_history_sent_at ON notification_history(sent_at DESC);

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ============================================================================
-- TRIGGERS
-- ============================================================================

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alerts_updated_at
    BEFORE UPDATE ON alerts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alert_rules_updated_at
    BEFORE UPDATE ON alert_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notification_channels_updated_at
    BEFORE UPDATE ON notification_channels
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA (Development Only)
-- ============================================================================

-- Insert default admin user (password: admin123)
INSERT INTO users (email, password_hash, name, role) VALUES
    ('admin@alerting.local', '$2a$10$rQEY7xQwXqzTvXqPqXqXqOqXqXqXqXqXqXqXqXqXqXqXqXqXqXqXq', 'Admin User', 'admin')
ON CONFLICT (email) DO NOTHING;

RAISE NOTICE 'Database initialization completed successfully!';
