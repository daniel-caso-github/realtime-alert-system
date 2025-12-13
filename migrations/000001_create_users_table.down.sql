-- Rollback: Drop users table

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
-- Note: We don't drop the function or extension as other tables may use them
