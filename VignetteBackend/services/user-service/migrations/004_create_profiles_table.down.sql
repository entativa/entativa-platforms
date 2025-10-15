-- Drop trigger and function
DROP TRIGGER IF EXISTS profiles_updated_at_trigger ON profiles;
DROP FUNCTION IF EXISTS update_profiles_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_profiles_updated_at;
DROP INDEX IF EXISTS idx_profiles_profile_views;
DROP INDEX IF EXISTS idx_profiles_category;
DROP INDEX IF EXISTS idx_profiles_user_id;

-- Drop table
DROP TABLE IF EXISTS profiles;
