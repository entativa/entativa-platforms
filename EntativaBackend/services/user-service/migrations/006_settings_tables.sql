-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- User settings table
CREATE TABLE IF NOT EXISTS user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Privacy settings
    is_private_account BOOLEAN DEFAULT FALSE,
    show_activity_status BOOLEAN DEFAULT TRUE,
    read_receipts BOOLEAN DEFAULT TRUE,
    allow_message_requests BOOLEAN DEFAULT TRUE,
    posts_visibility VARCHAR(20) DEFAULT 'everyone', -- everyone, friends, only_me
    comments_allowed VARCHAR(20) DEFAULT 'everyone', -- everyone, friends, no_one
    mentions_allowed VARCHAR(20) DEFAULT 'everyone', -- everyone, following, off
    story_sharing VARCHAR(20) DEFAULT 'everyone',    -- everyone, following, off
    similar_account_suggestions BOOLEAN DEFAULT TRUE,
    include_in_recommendations BOOLEAN DEFAULT TRUE,
    
    -- Notification settings
    notify_likes BOOLEAN DEFAULT TRUE,
    notify_comments BOOLEAN DEFAULT TRUE,
    notify_followers BOOLEAN DEFAULT TRUE,
    notify_messages BOOLEAN DEFAULT TRUE,
    notify_friend_requests BOOLEAN DEFAULT TRUE,
    notify_video_views BOOLEAN DEFAULT TRUE,
    notify_live_videos BOOLEAN DEFAULT TRUE,
    email_weekly_summary BOOLEAN DEFAULT TRUE,
    email_product_updates BOOLEAN DEFAULT FALSE,
    email_tips BOOLEAN DEFAULT FALSE,
    notification_sound BOOLEAN DEFAULT TRUE,
    notification_vibration BOOLEAN DEFAULT TRUE,
    show_badge_count BOOLEAN DEFAULT TRUE,
    
    -- Data settings
    upload_quality VARCHAR(20) DEFAULT 'high', -- high, medium, low, normal, basic
    autoplay_settings VARCHAR(20) DEFAULT 'wifi', -- always, wifi, never
    data_saver_mode BOOLEAN DEFAULT FALSE,
    use_less_data BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Blocked users table
CREATE TABLE IF NOT EXISTS blocked_users (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, blocked_user_id),
    
    INDEX idx_blocked_users_user (user_id),
    INDEX idx_blocked_users_blocked (blocked_user_id)
);

-- Muted users table (for Instagram-style mute)
CREATE TABLE IF NOT EXISTS muted_users (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    muted_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    muted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    mute_stories BOOLEAN DEFAULT TRUE,
    mute_posts BOOLEAN DEFAULT TRUE,
    
    PRIMARY KEY (user_id, muted_user_id),
    
    INDEX idx_muted_users_user (user_id)
);

-- Restricted users table (Instagram-style restrict)
CREATE TABLE IF NOT EXISTS restricted_users (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    restricted_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    restricted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, restricted_user_id),
    
    INDEX idx_restricted_users_user (user_id)
);

-- Add deletion columns to users table if they don't exist
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add device_name and location to sessions table for login activity
ALTER TABLE sessions
ADD COLUMN IF NOT EXISTS device_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS location VARCHAR(255);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_settings_updated_at
    BEFORE UPDATE ON user_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_settings_updated_at();

-- Comments
COMMENT ON TABLE user_settings IS 'Stores all user preferences and settings';
COMMENT ON TABLE blocked_users IS 'Tracks blocked user relationships';
COMMENT ON TABLE muted_users IS 'Tracks muted user relationships (Instagram-style)';
COMMENT ON TABLE restricted_users IS 'Tracks restricted user relationships (Instagram-style)';

COMMENT ON COLUMN user_settings.is_private_account IS 'If true, user must approve followers';
COMMENT ON COLUMN user_settings.posts_visibility IS 'Who can see user posts: everyone, friends, only_me';
COMMENT ON COLUMN user_settings.upload_quality IS 'Quality for uploading media: high, medium, low';
COMMENT ON COLUMN user_settings.autoplay_settings IS 'When to autoplay videos: always, wifi, never';

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_settings_privacy ON user_settings(is_private_account);
CREATE INDEX IF NOT EXISTS idx_user_settings_activity ON user_settings(show_activity_status);
CREATE INDEX IF NOT EXISTS idx_users_deleted ON users(deleted_at) WHERE deleted_at IS NOT NULL;
