-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add founder flags to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_founder BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS admin_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS is_banned BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS banned_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS ban_reason TEXT,
ADD COLUMN IF NOT EXISTS ban_expires_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS is_shadowbanned BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS shadowban_reason TEXT,
ADD COLUMN IF NOT EXISTS is_muted BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS muted_until TIMESTAMP,
ADD COLUMN IF NOT EXISTS is_suspended BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS suspended_until TIMESTAMP,
ADD COLUMN IF NOT EXISTS suspension_reason TEXT,
ADD COLUMN IF NOT EXISTS force_password_reset BOOLEAN DEFAULT FALSE;

-- Set Neo Qiss as founder (update when account exists)
-- UPDATE users SET is_founder = TRUE, admin_level = 10 WHERE username = 'neoqiss';

-- Admin audit logs table
CREATE TABLE IF NOT EXISTS admin_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    target_type VARCHAR(50), -- 'user', 'content', 'platform', 'system'
    target_id UUID,
    reason TEXT,
    severity VARCHAR(20) DEFAULT 'normal', -- 'normal', 'high', 'critical'
    
    -- Metadata
    metadata JSONB,
    
    -- Context
    ip_address VARCHAR(45),
    device_id VARCHAR(255),
    user_agent TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_audit_admin (admin_id),
    INDEX idx_audit_target (target_type, target_id),
    INDEX idx_audit_action (action),
    INDEX idx_audit_created (created_at DESC),
    INDEX idx_audit_severity (severity)
);

-- Whitelisted devices for founder access
CREATE TABLE IF NOT EXISTS admin_whitelisted_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255),
    device_type VARCHAR(50), -- 'ios', 'android', 'web'
    device_fingerprint TEXT,
    
    first_seen TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    
    UNIQUE (user_id, device_id),
    INDEX idx_devices_user (user_id),
    INDEX idx_devices_active (user_id, is_active)
);

-- Impersonation sessions (for user impersonation with step-up auth)
CREATE TABLE IF NOT EXISTS impersonation_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES users(id),
    target_user_id UUID NOT NULL REFERENCES users(id),
    reason TEXT NOT NULL,
    
    impersonation_token VARCHAR(500) NOT NULL,
    
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL, -- 10-minute timeout
    ended_at TIMESTAMP,
    
    -- Audit trail
    ip_address VARCHAR(45),
    device_id VARCHAR(255),
    
    INDEX idx_impersonation_admin (admin_id),
    INDEX idx_impersonation_target (target_user_id),
    INDEX idx_impersonation_active (expires_at) WHERE ended_at IS NULL
);

-- Feature flags (for A/B testing and gradual rollouts)
CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    feature_name VARCHAR(100) NOT NULL UNIQUE,
    is_enabled BOOLEAN DEFAULT FALSE,
    rollout_percentage INTEGER DEFAULT 0, -- 0-100
    description TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES users(id),
    
    INDEX idx_features_enabled (is_enabled)
);

-- Platform kill switches (emergency controls)
CREATE TABLE IF NOT EXISTS kill_switches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    switch_name VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT FALSE,
    reason TEXT,
    
    activated_at TIMESTAMP,
    activated_by UUID REFERENCES users(id),
    deactivated_at TIMESTAMP,
    
    INDEX idx_switches_active (is_active)
);

-- Maintenance mode configuration
CREATE TABLE IF NOT EXISTS maintenance_mode (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN DEFAULT FALSE,
    message TEXT,
    estimated_end TIMESTAMP,
    
    enabled_at TIMESTAMP,
    enabled_by UUID REFERENCES users(id),
    disabled_at TIMESTAMP,
    
    INDEX idx_maintenance_active (is_active)
);

-- IP blocks (for banning IP ranges)
CREATE TABLE IF NOT EXISTS ip_blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address VARCHAR(45),
    ip_range_start VARCHAR(45),
    ip_range_end VARCHAR(45),
    reason TEXT NOT NULL,
    
    blocked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    blocked_by UUID NOT NULL REFERENCES users(id),
    expires_at TIMESTAMP,
    
    INDEX idx_ip_blocks_address (ip_address),
    INDEX idx_ip_blocks_range (ip_range_start, ip_range_end)
);

-- Broadcast notifications (platform-wide announcements)
CREATE TABLE IF NOT EXISTS broadcast_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    notification_type VARCHAR(50) DEFAULT 'info', -- 'info', 'warning', 'critical', 'announcement'
    
    target_audience VARCHAR(50) DEFAULT 'all', -- 'all', 'premium', 'creators'
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL REFERENCES users(id),
    scheduled_for TIMESTAMP,
    sent_at TIMESTAMP,
    
    recipient_count INTEGER DEFAULT 0,
    
    INDEX idx_broadcasts_created (created_at DESC),
    INDEX idx_broadcasts_scheduled (scheduled_for)
);

-- User admin data (additional metadata for moderation)
CREATE TABLE IF NOT EXISTS user_admin_data (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Account flags
    is_verified BOOLEAN DEFAULT FALSE,
    is_creator BOOLEAN DEFAULT FALSE,
    is_premium BOOLEAN DEFAULT FALSE,
    
    -- Moderation history
    ban_count INTEGER DEFAULT 0,
    last_ban_at TIMESTAMP,
    shadowban_count INTEGER DEFAULT 0,
    warning_count INTEGER DEFAULT 0,
    
    -- Trust score
    trust_score DECIMAL(3,2) DEFAULT 1.00, -- 0.00 to 1.00
    
    -- Notes
    admin_notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default feature flags
INSERT INTO feature_flags (feature_name, is_enabled, description) VALUES
('posting', TRUE, 'Allow users to create posts'),
('commenting', TRUE, 'Allow users to comment'),
('messaging', TRUE, 'Allow direct messaging'),
('takes', TRUE, 'Allow Takes (short videos)'),
('stories', TRUE, 'Allow Stories'),
('live_streaming', FALSE, 'Live video streaming'),
('voice_rooms', FALSE, 'Audio rooms/Spaces'),
('marketplace', FALSE, 'Buy/sell marketplace')
ON CONFLICT (feature_name) DO NOTHING;

-- Insert default kill switches
INSERT INTO kill_switches (switch_name) VALUES
('disable_posting'),
('disable_commenting'),
('disable_messaging'),
('disable_algorithm'),
('emergency_lockdown')
ON CONFLICT (switch_name) DO NOTHING;

-- Insert single maintenance mode record
INSERT INTO maintenance_mode (is_active) VALUES (FALSE)
ON CONFLICT DO NOTHING;

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_admin_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_feature_flags_updated_at
    BEFORE UPDATE ON feature_flags
    FOR EACH ROW
    EXECUTE FUNCTION update_admin_updated_at();

CREATE TRIGGER trigger_update_user_admin_data_updated_at
    BEFORE UPDATE ON user_admin_data
    FOR EACH ROW
    EXECUTE FUNCTION update_admin_updated_at();

-- Function to auto-expire temporary bans
CREATE OR REPLACE FUNCTION cleanup_expired_bans()
RETURNS void AS $$
BEGIN
    UPDATE users
    SET 
        is_banned = FALSE,
        ban_expires_at = NULL
    WHERE is_banned = TRUE
      AND ban_expires_at IS NOT NULL
      AND ban_expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Function to auto-expire impersonation sessions
CREATE OR REPLACE FUNCTION cleanup_expired_impersonations()
RETURNS void AS $$
BEGIN
    UPDATE impersonation_sessions
    SET ended_at = NOW()
    WHERE ended_at IS NULL
      AND expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE admin_audit_logs IS 'Complete audit trail of all admin actions';
COMMENT ON TABLE admin_whitelisted_devices IS 'Devices authorized for founder admin access';
COMMENT ON TABLE impersonation_sessions IS 'Tracks user impersonation sessions with 10-min timeout';
COMMENT ON TABLE feature_flags IS 'Platform feature toggles for gradual rollouts';
COMMENT ON TABLE kill_switches IS 'Emergency controls to disable platform features';
COMMENT ON TABLE maintenance_mode IS 'Platform maintenance mode configuration';
COMMENT ON TABLE ip_blocks IS 'Blocked IP addresses and ranges';
COMMENT ON TABLE broadcast_notifications IS 'Platform-wide announcements';
COMMENT ON TABLE user_admin_data IS 'Additional metadata for user moderation';

COMMENT ON COLUMN users.is_founder IS 'Reserved for @neoqiss - Supreme admin level 10';
COMMENT ON COLUMN users.admin_level IS '0=normal, 10=founder (Neo Qiss only)';
COMMENT ON COLUMN admin_audit_logs.severity IS 'normal, high, or critical';
COMMENT ON COLUMN impersonation_sessions.expires_at IS 'Auto-terminates after 10 minutes';

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_founder ON users(is_founder) WHERE is_founder = TRUE;
CREATE INDEX IF NOT EXISTS idx_users_banned ON users(is_banned) WHERE is_banned = TRUE;
CREATE INDEX IF NOT EXISTS idx_users_shadowbanned ON users(is_shadowbanned) WHERE is_shadowbanned = TRUE;
