-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    notification_type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    
    -- Actor (who triggered it)
    actor_id UUID,
    actor_username VARCHAR(255),
    actor_avatar_url TEXT,
    
    -- Related entities
    post_id UUID,
    take_id UUID,
    comment_id UUID,
    story_id UUID,
    trend_id UUID,
    
    -- Metadata
    data JSONB,
    image_url TEXT,
    deep_link TEXT,
    
    -- Status
    is_read BOOLEAN DEFAULT FALSE,
    is_delivered BOOLEAN DEFAULT FALSE,
    delivery_channels TEXT[],
    priority VARCHAR(20) DEFAULT 'Normal',
    
    -- Grouping
    group_key VARCHAR(255),
    group_count INT DEFAULT 1,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    read_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for notifications
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE NOT is_read;
CREATE INDEX idx_notifications_group_key ON notifications(group_key) WHERE group_key IS NOT NULL;
CREATE INDEX idx_notifications_expires ON notifications(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_notifications_type ON notifications(notification_type, created_at DESC);
CREATE INDEX idx_notifications_actor ON notifications(actor_id) WHERE actor_id IS NOT NULL;

-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    device_token TEXT NOT NULL UNIQUE,
    platform VARCHAR(20) NOT NULL,
    device_name VARCHAR(255),
    device_model VARCHAR(255),
    os_version VARCHAR(50),
    app_version VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for devices
CREATE INDEX idx_devices_user ON devices(user_id);
CREATE INDEX idx_devices_token ON devices(device_token);
CREATE INDEX idx_devices_active ON devices(user_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_devices_platform ON devices(platform, is_active);

-- Create notification preferences table
CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id UUID PRIMARY KEY,
    enable_push BOOLEAN DEFAULT TRUE,
    enable_email BOOLEAN DEFAULT TRUE,
    enable_sms BOOLEAN DEFAULT FALSE,
    
    -- Fine-grained preferences
    notify_on_like BOOLEAN DEFAULT TRUE,
    notify_on_comment BOOLEAN DEFAULT TRUE,
    notify_on_follow BOOLEAN DEFAULT TRUE,
    notify_on_mention BOOLEAN DEFAULT TRUE,
    notify_on_share BOOLEAN DEFAULT TRUE,
    notify_on_take_remix BOOLEAN DEFAULT TRUE,
    notify_on_story_reply BOOLEAN DEFAULT TRUE,
    notify_on_tagged BOOLEAN DEFAULT TRUE,
    
    -- Quiet hours
    quiet_hours_enabled BOOLEAN DEFAULT FALSE,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create notification templates table
CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    subject TEXT NOT NULL,
    body_text TEXT NOT NULL,
    body_html TEXT,
    category VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Auto-update trigger
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notifications_updated_at
    BEFORE UPDATE ON notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER devices_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER preferences_updated_at
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- Comments
COMMENT ON TABLE notifications IS 'User notifications with multi-channel delivery';
COMMENT ON COLUMN notifications.group_key IS 'Key for grouping similar notifications (e.g., "like:post-uuid")';
COMMENT ON COLUMN notifications.group_count IS 'Number of grouped notifications';
COMMENT ON COLUMN notifications.delivery_channels IS 'Array of delivery channels (inapp, push, email, sms, websocket)';
COMMENT ON COLUMN notifications.priority IS 'Notification priority (Low, Normal, High, Urgent)';

COMMENT ON TABLE devices IS 'Registered devices for push notifications';
COMMENT ON COLUMN devices.device_token IS 'FCM token (Android) or APN token (iOS)';
COMMENT ON COLUMN devices.platform IS 'Device platform (iOS, Android, Web)';

COMMENT ON TABLE notification_preferences IS 'User notification preferences';
COMMENT ON COLUMN notification_preferences.quiet_hours_start IS 'Start time for quiet hours (HH:MM format)';
COMMENT ON COLUMN notification_preferences.quiet_hours_end IS 'End time for quiet hours (HH:MM format)';
