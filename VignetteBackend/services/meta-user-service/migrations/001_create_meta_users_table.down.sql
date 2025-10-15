-- Drop tables and functions
DROP TRIGGER IF EXISTS update_meta_users_updated_at ON meta_users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS meta_user_events;
DROP TABLE IF EXISTS meta_users;
