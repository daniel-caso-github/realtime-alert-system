-- Migration: Seed initial admin user
-- Description: Create default admin user for first login
-- Password: Admin123! (bcrypt hash)

INSERT INTO users (email, password_hash, name, role, is_active)
VALUES (
    'admin@alerting.local',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.aWz2DWFDBaENthCveu',
    'System Administrator',
    'admin',
    true
)
ON CONFLICT (email) DO NOTHING;
