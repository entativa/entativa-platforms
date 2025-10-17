-- Settings Service Database Schema
-- Comprehensive app settings + encrypted key backup with PIN/Passphrase

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- USER SETTINGS TABLE
-- ============================================
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    
    -- All settings stored as JSONB for flexibility
    appearance JSONB DEFAULT '{"theme": "auto", "accent_color": "#007AFF", "font_size": "medium", "high_contrast": false, "reduce_motion": false, "compact_mode": false}',
    privacy JSONB DEFAULT '{"profile_visibility": "public", "last_seen": "everyone", "read_receipts": true, "typing_indicator": true, "online_status": true, "blocked_users": [], "allow_tagging": "everyone", "allow_mentions": "everyone", "searchable_by_email": true, "searchable_by_phone": true, "show_activity": true}',
    notifications JSONB DEFAULT '{"push_enabled": true, "email_enabled": true, "sms_enabled": false, "likes": true, "comments": true, "mentions": true, "follows": true, "messages": true, "group_messages": true, "event_invites": true, "event_reminders": true, "live_streams": true, "quiet_hours_enabled": false, "quiet_hours_start": "22:00", "quiet_hours_end": "08:00", "notification_sound": "default", "vibrate": true}',
    chat JSONB DEFAULT '{"key_storage_location": "entativa_server", "encryption_method": "passphrase", "backup_keys_to_server": true, "enter_to_send": false, "auto_download_media": true, "auto_play_videos": true, "auto_play_gifs": true, "save_to_gallery": false, "auto_delete_messages": false, "auto_delete_after_days": 0, "screen_security": false, "incognito_keyboard": false}',
    media JSONB DEFAULT '{"auto_download_photos": true, "auto_download_videos": false, "auto_download_files": false, "upload_quality": "high", "video_quality": "auto", "media_storage_location": "internal", "auto_delete_media": false, "auto_delete_after_days": 0}',
    data_storage JSONB DEFAULT '{"data_saver_mode": false, "low_data_mode": false, "wifi_only": false, "cache_size": 500, "auto_clear_cache": false, "auto_clear_after_days": 30}',
    security JSONB DEFAULT '{"two_factor_enabled": false, "biometric_enabled": false, "app_lock_enabled": false, "app_lock_timeout": 300, "active_sessions": 1, "show_login_alerts": true, "recovery_email": "", "recovery_phone": ""}',
    accessibility JSONB DEFAULT '{"screen_reader": false, "closed_captions": false, "color_blind_mode": "none", "high_contrast_text": false, "large_text": false, "reduce_transparency": false, "voice_control": false}',
    language JSONB DEFAULT '{"app_language": "en", "content_languages": ["en"], "translation_enabled": true, "auto_translate": false}',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_settings_user ON user_settings(user_id);

-- ============================================
-- ENCRYPTED KEY BACKUPS TABLE
-- ============================================
CREATE TABLE encrypted_key_backups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_settings(user_id) ON DELETE CASCADE,
    
    -- Storage configuration
    storage_location VARCHAR(50) NOT NULL CHECK (storage_location IN ('entativa_server', 'local_device', 'icloud', 'google_drive')),
    encryption_method VARCHAR(20) NOT NULL CHECK (encryption_method IN ('pin', 'passphrase')),
    
    -- Encrypted data (CRITICAL: This is double-encrypted)
    -- 1. First encryption: Signal protocol (E2EE)
    -- 2. Second encryption: User's PIN/Passphrase (only user can decrypt)
    encrypted_keys BYTEA NOT NULL,
    keys_hash VARCHAR(64) NOT NULL, -- SHA256 for integrity verification
    
    -- PIN/Passphrase protection (hashed with bcrypt)
    pin_hash TEXT NOT NULL,  -- bcrypt hash of PIN/Passphrase
    salt VARCHAR(64) NOT NULL, -- Random salt for key derivation
    iterations INTEGER DEFAULT 100000, -- PBKDF2 iterations
    
    -- Metadata (ONLY this is visible in plain text)
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    backup_version INTEGER DEFAULT 1,
    last_backup_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, device_id)
);

CREATE INDEX idx_key_backups_user ON encrypted_key_backups(user_id);
CREATE INDEX idx_key_backups_storage ON encrypted_key_backups(storage_location);
CREATE INDEX idx_key_backups_last_backup ON encrypted_key_backups(last_backup_at DESC);

-- ============================================
-- SETTING CHANGE HISTORY (Audit log)
-- ============================================
CREATE TABLE settings_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    setting_category VARCHAR(50) NOT NULL, -- appearance, privacy, notifications, etc
    setting_key VARCHAR(100) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT
);

CREATE INDEX idx_settings_history_user ON settings_history(user_id, changed_at DESC);
CREATE INDEX idx_settings_history_category ON settings_history(setting_category);

-- ============================================
-- KEY BACKUP ACCESS LOG (Security audit)
-- ============================================
CREATE TABLE key_backup_access_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    backup_id UUID NOT NULL REFERENCES encrypted_key_backups(id) ON DELETE CASCADE,
    action VARCHAR(20) NOT NULL CHECK (action IN ('create', 'restore', 'update', 'delete', 'failed_restore')),
    device_id VARCHAR(255),
    ip_address INET,
    success BOOLEAN DEFAULT TRUE,
    failure_reason TEXT,
    accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_backup_access_user ON key_backup_access_log(user_id, accessed_at DESC);
CREATE INDEX idx_backup_access_action ON key_backup_access_log(action, accessed_at DESC);
CREATE INDEX idx_backup_access_failed ON key_backup_access_log(user_id) WHERE NOT success;

-- ============================================
-- TRIGGERS
-- ============================================

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_settings_updated_at
    BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER encrypted_key_backups_updated_at
    BEFORE UPDATE ON encrypted_key_backups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Auto-log settings changes
CREATE OR REPLACE FUNCTION log_settings_changes()
RETURNS TRIGGER AS $$
BEGIN
    -- Log each JSONB field change
    IF OLD.appearance IS DISTINCT FROM NEW.appearance THEN
        INSERT INTO settings_history (user_id, setting_category, setting_key, old_value, new_value)
        VALUES (NEW.user_id, 'appearance', 'full', OLD.appearance::text, NEW.appearance::text);
    END IF;
    
    IF OLD.privacy IS DISTINCT FROM NEW.privacy THEN
        INSERT INTO settings_history (user_id, setting_category, setting_key, old_value, new_value)
        VALUES (NEW.user_id, 'privacy', 'full', OLD.privacy::text, NEW.privacy::text);
    END IF;
    
    IF OLD.chat IS DISTINCT FROM NEW.chat THEN
        INSERT INTO settings_history (user_id, setting_category, setting_key, old_value, new_value)
        VALUES (NEW.user_id, 'chat', 'full', OLD.chat::text, NEW.chat::text);
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_user_settings_changes
    AFTER UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION log_settings_changes();

-- Comments
COMMENT ON TABLE user_settings IS 'Comprehensive user settings (appearance, privacy, notifications, chat, etc)';
COMMENT ON TABLE encrypted_key_backups IS 'ENCRYPTED CHAT KEY BACKUPS - Double-encrypted with Signal + PIN/Passphrase. Authorities only get metadata!';
COMMENT ON COLUMN encrypted_key_backups.encrypted_keys IS 'CRITICAL: Double-encrypted keys (Signal + User PIN). Server cannot decrypt!';
COMMENT ON COLUMN encrypted_key_backups.pin_hash IS 'bcrypt hash of user PIN/Passphrase. Never stored in plain text!';
COMMENT ON COLUMN encrypted_key_backups.salt IS 'Random salt for PBKDF2 key derivation';
COMMENT ON TABLE settings_history IS 'Audit log of all settings changes';
COMMENT ON TABLE key_backup_access_log IS 'Security audit log of key backup access (create/restore/failed attempts)';
