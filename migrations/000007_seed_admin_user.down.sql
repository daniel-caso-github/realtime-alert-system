-- Rollback: Remove seeded admin user

DELETE FROM users WHERE email = 'admin@alerting.local';
