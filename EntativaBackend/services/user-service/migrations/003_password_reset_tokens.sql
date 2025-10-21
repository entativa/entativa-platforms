-- Create password_reset_tokens table
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_token (token),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
);

-- Add index for cleanup queries
CREATE INDEX idx_expired_tokens ON password_reset_tokens(expires_at) WHERE used = false;

-- Add updated_at trigger
CREATE OR REPLACE FUNCTION update_password_reset_tokens_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Comment on table
COMMENT ON TABLE password_reset_tokens IS 'Stores password reset tokens with expiration';
COMMENT ON COLUMN password_reset_tokens.token IS 'Secure random token for password reset';
COMMENT ON COLUMN password_reset_tokens.expires_at IS 'Token expiration time (1 hour from creation)';
COMMENT ON COLUMN password_reset_tokens.used IS 'Whether token has been used';
