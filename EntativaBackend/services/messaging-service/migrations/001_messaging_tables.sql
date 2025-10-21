-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL DEFAULT 'direct', -- 'direct' or 'group'
    name VARCHAR(255), -- For group chats
    avatar_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_message_at TIMESTAMP,
    
    INDEX idx_conversations_updated (updated_at DESC)
);

-- Conversation participants
CREATE TABLE IF NOT EXISTS conversation_participants (
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    role VARCHAR(20) DEFAULT 'member', -- 'admin', 'member'
    muted BOOLEAN DEFAULT FALSE,
    last_read_at TIMESTAMP,
    
    PRIMARY KEY (conversation_id, user_id),
    INDEX idx_participants_user (user_id),
    INDEX idx_participants_conversation (conversation_id)
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    
    -- E2EE encrypted content
    encrypted_content TEXT NOT NULL, -- Base64 encoded ciphertext
    content_type VARCHAR(20) NOT NULL DEFAULT 'text', -- 'text', 'image', 'video', 'audio', 'file'
    
    -- Metadata (not encrypted)
    media_url TEXT, -- For media messages (encrypted on storage)
    media_key TEXT, -- Encryption key for media (encrypted per recipient)
    thumbnail_url TEXT,
    
    -- Message status
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    
    -- Features
    reply_to UUID REFERENCES messages(id) ON DELETE SET NULL,
    forward_from UUID REFERENCES messages(id) ON DELETE SET NULL,
    expires_at TIMESTAMP, -- For disappearing messages
    deleted_at TIMESTAMP, -- Soft delete
    edited_at TIMESTAMP,
    
    INDEX idx_messages_conversation (conversation_id, sent_at DESC),
    INDEX idx_messages_sender (sender_id),
    INDEX idx_messages_status (conversation_id, delivered_at, read_at)
);

-- Message receipts (for group chats)
CREATE TABLE IF NOT EXISTS message_receipts (
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    
    PRIMARY KEY (message_id, user_id),
    INDEX idx_receipts_message (message_id),
    INDEX idx_receipts_user (user_id)
);

-- E2EE: Identity keys (long-term Curve25519 keys)
CREATE TABLE IF NOT EXISTS identity_keys (
    user_id UUID PRIMARY KEY,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_identity_user (user_id)
);

-- E2EE: Signed prekeys
CREATE TABLE IF NOT EXISTS signed_prekeys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    key_id INTEGER NOT NULL,
    public_key BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (user_id, key_id),
    INDEX idx_signed_prekeys_user (user_id)
);

-- E2EE: One-time prekeys
CREATE TABLE IF NOT EXISTS one_time_prekeys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    key_id INTEGER NOT NULL,
    public_key BYTEA NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (user_id, key_id),
    INDEX idx_prekeys_user (user_id),
    INDEX idx_prekeys_unused (user_id, used) WHERE used = FALSE
);

-- E2EE: Session state (Double Ratchet sessions)
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    partner_id UUID NOT NULL,
    
    -- Session state (encrypted, stored as JSON)
    root_key BYTEA NOT NULL,
    sending_chain_key BYTEA,
    receiving_chain_key BYTEA,
    sending_ratchet_key BYTEA,
    receiving_ratchet_key BYTEA,
    send_counter INTEGER DEFAULT 0,
    receive_counter INTEGER DEFAULT 0,
    previous_counter INTEGER DEFAULT 0,
    
    -- Skipped message keys (for out-of-order messages)
    skipped_messages JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (user_id, partner_id),
    INDEX idx_sessions_user (user_id),
    INDEX idx_sessions_partner (partner_id)
);

-- Typing indicators (ephemeral, can use Redis in production)
CREATE TABLE IF NOT EXISTS typing_indicators (
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (conversation_id, user_id),
    INDEX idx_typing_conversation (conversation_id)
);

-- Message reactions
CREATE TABLE IF NOT EXISTS message_reactions (
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    emoji VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (message_id, user_id),
    INDEX idx_reactions_message (message_id)
);

-- Call history
CREATE TABLE IF NOT EXISTS calls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    caller_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'audio', 'video'
    status VARCHAR(20) NOT NULL, -- 'missed', 'completed', 'declined', 'failed'
    duration INTEGER, -- Duration in seconds
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    
    INDEX idx_calls_conversation (conversation_id),
    INDEX idx_calls_caller (caller_id)
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_update_sessions_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- Function to clean up old typing indicators (run periodically)
CREATE OR REPLACE FUNCTION cleanup_old_typing_indicators()
RETURNS void AS $$
BEGIN
    DELETE FROM typing_indicators
    WHERE started_at < NOW() - INTERVAL '10 seconds';
END;
$$ LANGUAGE plpgsql;

-- Function to clean up expired messages
CREATE OR REPLACE FUNCTION cleanup_expired_messages()
RETURNS void AS $$
BEGIN
    UPDATE messages
    SET deleted_at = NOW(),
        encrypted_content = '[Message deleted]'
    WHERE expires_at IS NOT NULL
      AND expires_at < NOW()
      AND deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE conversations IS 'Stores conversation metadata (direct or group chats)';
COMMENT ON TABLE messages IS 'Stores E2EE encrypted messages';
COMMENT ON TABLE identity_keys IS 'User long-term identity keys for Signal Protocol';
COMMENT ON TABLE signed_prekeys IS 'Signed prekeys for X3DH key exchange';
COMMENT ON TABLE one_time_prekeys IS 'One-time prekeys for perfect forward secrecy';
COMMENT ON TABLE sessions IS 'Double Ratchet session state between users';
COMMENT ON COLUMN messages.encrypted_content IS 'Encrypted with Signal Protocol Double Ratchet';
COMMENT ON COLUMN messages.expires_at IS 'Timestamp for disappearing messages';

-- Insert some sample data for testing
-- Note: In production, keys should be properly generated using the Signal Protocol library

-- Sample identity keys (these are just examples, not real keys)
INSERT INTO identity_keys (user_id, public_key) VALUES
(uuid_generate_v4(), decode('0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef', 'hex'))
ON CONFLICT (user_id) DO NOTHING;
