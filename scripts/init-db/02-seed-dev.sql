-- ============================================================================
-- Development Seed Data
-- ============================================================================
-- This script adds sample data for development and testing

-- Sample notification channels
INSERT INTO notification_channels (id, name, type, config, is_enabled) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'Slack General', 'slack',
     '{"webhook_url": "https://hooks.slack.com/services/xxx/yyy/zzz", "channel": "#alerts"}', true),
    ('a0000000-0000-0000-0000-000000000002', 'Email Operations', 'email',
     '{"smtp_host": "smtp.example.com", "recipients": ["ops@example.com"]}', true),
    ('a0000000-0000-0000-0000-000000000003', 'PagerDuty Webhook', 'webhook',
     '{"url": "https://events.pagerduty.com/v2/enqueue", "method": "POST"}', false)
ON CONFLICT DO NOTHING;

-- Sample alert rules
INSERT INTO alert_rules (id, name, description, condition, severity, is_enabled, cooldown_minutes) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'High CPU Usage',
     'Triggers when CPU usage exceeds 90%',
     '{"metric": "cpu_usage", "operator": ">", "threshold": 90}',
     'high', true, 10),
    ('b0000000-0000-0000-0000-000000000002', 'Memory Critical',
     'Triggers when memory usage exceeds 95%',
     '{"metric": "memory_usage", "operator": ">", "threshold": 95}',
     'critical', true, 5),
    ('b0000000-0000-0000-0000-000000000003', 'Disk Space Low',
     'Triggers when disk space falls below 10%',
     '{"metric": "disk_free", "operator": "<", "threshold": 10}',
     'medium', true, 30),
    ('b0000000-0000-0000-0000-000000000004', 'Service Health Check Failed',
     'Triggers when health check fails 3 times',
     '{"metric": "health_check", "operator": "==", "threshold": 0, "consecutive": 3}',
     'critical', true, 1)
ON CONFLICT DO NOTHING;

-- Link rules to channels
INSERT INTO alert_rule_channels (rule_id, channel_id) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001'),
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001'),
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000002'),
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001'),
    ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001'),
    ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000002')
ON CONFLICT DO NOTHING;

-- Sample alerts
INSERT INTO alerts (id, rule_id, title, message, severity, status, source, metadata) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
     'High CPU on web-server-01', 'CPU usage at 94% for the last 5 minutes',
     'high', 'active', 'web-server-01', '{"cpu_percent": 94, "duration_minutes": 5}'),
    ('c0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000003',
     'Low disk space on db-server-01', 'Only 8% disk space remaining',
     'medium', 'acknowledged', 'db-server-01', '{"disk_free_percent": 8, "disk_total_gb": 500}'),
    ('c0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000002',
     'Critical memory on api-server-02', 'Memory usage at 97%',
     'critical', 'resolved', 'api-server-02', '{"memory_percent": 97}')
ON CONFLICT DO NOTHING;

RAISE NOTICE 'Development seed data inserted successfully!';
