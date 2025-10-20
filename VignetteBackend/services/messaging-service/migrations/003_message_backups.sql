-- E2EE Message Backup Tables
-- Users can backup their encrypted messages to our servers (most secure)
-- Or to Google Drive/iCloud (with warnings about potential decryption)

-- Backup encryption keys (derived from user's PIN/passphrase)
CREATE TABLE IF NOT EXISTS backup_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Key derivation
    salt BYTEA NOT NULL,                    -- Random salt for PBKDF2/Argon2
    iterations INTEGER NOT NULL,             -- Key derivation iterations
    algorithm VARCHAR(50) NOT NULL,          -- 'argon2id' or 'pbkdf2-sha256'
    
    -- Encrypted master backup key
    encrypted_backup_key BYTEA NOT NULL,     -- AES-256 encrypted with derived key
    backup_key_nonce BYTEA NOT NULL,         -- Nonce for backup key encryption
    
    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_backup_at TIMESTAMP,
    
    UNIQUE(user_id),
    INDEX idx_backup_keys_user (user_id)
);

-- Encrypted message backups (stored on our servers)
CREATE TABLE IF NOT EXISTS message_backups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    backup_key_id UUID NOT NULL REFERENCES backup_keys(id) ON DELETE CASCADE,
    
    -- Backup data
    encrypted_data BYTEA NOT NULL,           -- Encrypted backup blob
    backup_nonce BYTEA NOT NULL,             -- Nonce for backup encryption
    backup_type VARCHAR(50) NOT NULL,        -- 'full' or 'incremental'
    
    -- Metadata
    messages_count INTEGER NOT NULL,
    conversations_count INTEGER NOT NULL,
    backup_size BIGINT NOT NULL,             -- Size in bytes
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,                    -- Optional expiration
    
    INDEX idx_message_backups_user (user_id),
    INDEX idx_message_backups_created (created_at DESC)
);

-- Backup settings per user
CREATE TABLE IF NOT EXISTS backup_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Backup location
    backup_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    backup_location VARCHAR(50) NOT NULL DEFAULT 'our_servers', -- 'our_servers', 'google_drive', 'icloud'
    
    -- Auto-backup settings
    auto_backup_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    auto_backup_frequency VARCHAR(50) DEFAULT 'daily', -- 'daily', 'weekly', 'monthly'
    auto_backup_wifi_only BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Backup retention
    keep_backups_count INTEGER DEFAULT 7,   -- Keep last N backups
    
    -- Third-party backup (if not our_servers)
    third_party_account_id VARCHAR(255),    -- Google/Apple account identifier
    third_party_file_id VARCHAR(255),       -- File ID on Google Drive/iCloud
    
    -- Warnings acknowledged
    third_party_warning_acknowledged BOOLEAN DEFAULT FALSE,
    third_party_warning_acknowledged_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Backup activity log
CREATE TABLE IF NOT EXISTS backup_activity_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    action VARCHAR(50) NOT NULL,             -- 'created', 'restored', 'deleted', 'failed'
    backup_location VARCHAR(50) NOT NULL,
    backup_type VARCHAR(50),                 -- 'full' or 'incremental'
    
    -- Result
    success BOOLEAN NOT NULL,
    error_message TEXT,
    
    -- Metadata
    messages_backed_up INTEGER,
    backup_size BIGINT,
    duration_ms INTEGER,                     -- Time taken
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_backup_activity_user (user_id),
    INDEX idx_backup_activity_created (created_at DESC)
);

-- Backup restoration tokens (one-time use)
CREATE TABLE IF NOT EXISTS backup_restoration_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    token_hash VARCHAR(255) NOT NULL,        -- Hashed restoration token
    backup_id UUID REFERENCES message_backups(id) ON DELETE CASCADE,
    
    -- Expiry and usage
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    is_used BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_restoration_tokens_user (user_id),
    INDEX idx_restoration_tokens_expires (expires_at)
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_backup_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_backup_settings_updated_at
    BEFORE UPDATE ON backup_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_backup_settings_updated_at();

-- Function to cleanup old backups
CREATE OR REPLACE FUNCTION cleanup_old_backups()
RETURNS void AS $$
BEGIN
    -- Delete backups older than retention count
    DELETE FROM message_backups mb
    WHERE mb.id IN (
        SELECT mb2.id
        FROM message_backups mb2
        JOIN backup_settings bs ON mb2.user_id = bs.user_id
        WHERE mb2.id NOT IN (
            SELECT id
            FROM message_backups mb3
            WHERE mb3.user_id = mb2.user_id
            ORDER BY created_at DESC
            LIMIT bs.keep_backups_count
        )
    );
    
    -- Delete expired restoration tokens
    DELETE FROM backup_restoration_tokens
    WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE backup_keys IS 'User backup encryption keys derived from PIN/passphrase';
COMMENT ON TABLE message_backups IS 'Encrypted message backups stored on our servers';
COMMENT ON TABLE backup_settings IS 'User backup preferences and settings';
COMMENT ON TABLE backup_activity_log IS 'Audit log for all backup operations';
COMMENT ON TABLE backup_restoration_tokens IS 'One-time tokens for backup restoration';

COMMENT ON COLUMN backup_keys.salt IS 'Random salt for key derivation (PBKDF2/Argon2)';
COMMENT ON COLUMN backup_keys.encrypted_backup_key IS 'Master backup key encrypted with user PIN/passphrase';
COMMENT ON COLUMN backup_settings.backup_location IS 'our_servers (most secure), google_drive, or icloud (with warnings)';
COMMENT ON COLUMN backup_settings.third_party_warning_acknowledged IS 'User acknowledged that Google/Apple can decrypt backups';

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_backup_keys_last_backup ON backup_keys(last_backup_at DESC);
CREATE INDEX IF NOT EXISTS idx_message_backups_user_created ON message_backups(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_backup_activity_success ON backup_activity_log(user_id, success);
