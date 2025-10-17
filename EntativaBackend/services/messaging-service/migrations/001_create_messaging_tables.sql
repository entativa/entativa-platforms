-- Socialink Messaging Service Database Schema
-- Signal-level E2EE messaging with libsignal + MLS

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- DEVICES & KEYS (Identity, Pre-keys)
-- ============================================================================

-- Registered devices
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    registration_id INTEGER NOT NULL,
    identity_key BYTEA NOT NULL,           -- Ed25519 public key (32 bytes)
    is_active BOOLEAN DEFAULT TRUE,
    last_seen TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, device_id)
);

CREATE INDEX idx_devices_user ON devices(user_id);
CREATE INDEX idx_devices_active ON devices(user_id, is_active) WHERE is_active = TRUE;

-- Signed pre-keys (medium-term)
CREATE TABLE IF NOT EXISTS signed_prekeys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    prekey_id INTEGER NOT NULL,
    public_key BYTEA NOT NULL,             -- X25519 public key (32 bytes)
    signature BYTEA NOT NULL,              -- Signed by identity key (64 bytes)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, device_id, prekey_id)
);

CREATE INDEX idx_signed_prekeys_device ON signed_prekeys(user_id, device_id);

-- One-time pre-keys (single-use)
CREATE TABLE IF NOT EXISTS onetime_prekeys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    prekey_id INTEGER NOT NULL,
    public_key BYTEA NOT NULL,             -- X25519 public key (32 bytes)
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, device_id, prekey_id)
);

CREATE INDEX idx_onetime_prekeys_unused ON onetime_prekeys(user_id, device_id, is_used) WHERE is_used = FALSE;

-- ============================================================================
-- CONVERSATIONS & MESSAGES
-- ============================================================================

-- Conversations (1:1 or group)
CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_type JSONB NOT NULL,       -- OneToOne, Group, NoteToSelf
    name VARCHAR(255),                      -- For groups
    avatar_url TEXT,
    created_by UUID NOT NULL,
    last_message_id UUID,
    is_archived BOOLEAN DEFAULT FALSE,
    is_muted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_conversations_creator ON conversations(created_by);
CREATE INDEX idx_conversations_updated ON conversations(updated_at DESC);

-- Conversation participants
CREATE TABLE IF NOT EXISTS conversation_participants (
    conversation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_conv_participants_user ON conversation_participants(user_id);

-- Encrypted messages (server cannot decrypt!)
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    sender_device_id VARCHAR(255) NOT NULL,
    recipient_id UUID,                      -- For 1:1, NULL for groups
    message_type JSONB NOT NULL,            -- Text, Media, Audio, File, Location, Contact, Poll, Event, Call, System
    
    -- ENCRYPTED CONTENT (E2EE!)
    ciphertext BYTEA NOT NULL,              -- Encrypted message body
    ephemeral_key BYTEA,                    -- For key agreement (1:1 only)
    
    -- METADATA (visible to server)
    sequence_number BIGINT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    is_group BOOLEAN DEFAULT FALSE,
    group_epoch BIGINT,                     -- MLS epoch for groups
    
    -- Delivery tracking
    status JSONB NOT NULL,                  -- Sending, Sent, Delivered, Read, Failed
    delivered_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,
    
    -- Features
    is_self_destructing BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMPTZ,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id, sequence_number DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id, created_at DESC);
CREATE INDEX idx_messages_recipient ON messages(recipient_id, created_at DESC) WHERE recipient_id IS NOT NULL;
CREATE INDEX idx_messages_expires ON messages(expires_at) WHERE expires_at IS NOT NULL;

-- Deleted messages (per-user soft delete for E2EE)
CREATE TABLE IF NOT EXISTS deleted_messages (
    user_id UUID NOT NULL,
    message_id UUID NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, message_id)
);

-- ============================================================================
-- GROUP CHATS (MLS)
-- ============================================================================

-- Group chat metadata
CREATE TABLE IF NOT EXISTS group_chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    avatar_url TEXT,
    created_by UUID NOT NULL,
    conversation_id UUID,                   -- Link to conversation
    member_count INTEGER DEFAULT 1,
    max_members INTEGER DEFAULT 1500,
    is_public BOOLEAN DEFAULT FALSE,
    invite_link VARCHAR(255),
    
    -- MLS specific
    mls_group_id BYTEA NOT NULL,
    current_epoch BIGINT DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_groups_creator ON group_chats(created_by);
CREATE INDEX idx_groups_name ON group_chats(name);

-- Group members
CREATE TABLE IF NOT EXISTS group_members (
    group_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role JSONB NOT NULL,                    -- Owner, Admin, Member
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    added_by UUID,
    is_muted BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX idx_group_members_user ON group_members(user_id);

-- MLS group state (encrypted ratchet tree)
CREATE TABLE IF NOT EXISTS mls_group_states (
    group_id UUID PRIMARY KEY,
    epoch BIGINT NOT NULL,
    mls_group_state BYTEA NOT NULL,         -- Serialized MLS state
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- MLS Welcome messages (for new members)
CREATE TABLE IF NOT EXISTS mls_welcome_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    user_id UUID NOT NULL,
    welcome_message BYTEA NOT NULL,         -- Encrypted epoch secrets
    is_consumed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_welcome_messages_user ON mls_welcome_messages(user_id, is_consumed) WHERE is_consumed = FALSE;

-- ============================================================================
-- PRESENCE & TYPING
-- ============================================================================

-- User presence status (cached in Redis, persisted in DB)
CREATE TABLE IF NOT EXISTS user_presence (
    user_id UUID PRIMARY KEY,
    status VARCHAR(20) NOT NULL,            -- Online, Away, Busy, Offline
    last_seen TIMESTAMPTZ DEFAULT NOW(),
    custom_status VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_presence_status ON user_presence(status);

-- Typing indicators (ephemeral, stored briefly in Redis only)

-- ============================================================================
-- READ RECEIPTS
-- ============================================================================

-- Read receipts (who read what and when)
CREATE TABLE IF NOT EXISTS read_receipts (
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    read_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id)
);

CREATE INDEX idx_read_receipts_message ON read_receipts(message_id);
CREATE INDEX idx_read_receipts_user ON read_receipts(user_id, read_at DESC);

-- ============================================================================
-- CALLS (WebRTC)
-- ============================================================================

-- Audio/video calls
CREATE TABLE IF NOT EXISTS calls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL,
    caller_id UUID NOT NULL,
    call_type JSONB NOT NULL,               -- Audio, Video
    status JSONB NOT NULL,                  -- Ringing, Answered, Declined, Missed, Ended, Failed
    
    -- WebRTC signaling (encrypted)
    sdp_offer TEXT,
    sdp_answer TEXT,
    
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    duration_seconds INTEGER,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_calls_conversation ON calls(conversation_id, created_at DESC);
CREATE INDEX idx_calls_caller ON calls(caller_id, created_at DESC);

-- ICE candidates (WebRTC)
CREATE TABLE IF NOT EXISTS call_ice_candidates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    call_id UUID NOT NULL,
    user_id UUID NOT NULL,
    candidate TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ice_candidates_call ON call_ice_candidates(call_id);

-- ============================================================================
-- MEDIA (Encrypted files)
-- ============================================================================

-- Encrypted media metadata
CREATE TABLE IF NOT EXISTS encrypted_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL,
    media_type VARCHAR(50) NOT NULL,        -- Photo, Video, VoiceNote, Audio, Document
    encrypted_url TEXT NOT NULL,            -- URL to encrypted blob
    encryption_key BYTEA NOT NULL,          -- Key to decrypt media (32 bytes)
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100),
    duration_ms INTEGER,                    -- For audio/video
    width INTEGER,
    height INTEGER,
    thumbnail BYTEA,                        -- Encrypted thumbnail
    blurhash VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_encrypted_media_message ON encrypted_media(message_id);

-- ============================================================================
-- TRIGGERS & FUNCTIONS
-- ============================================================================

-- Auto-update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER group_chats_updated_at
    BEFORE UPDATE ON group_chats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE devices IS 'Registered devices for E2EE key exchange';
COMMENT ON TABLE signed_prekeys IS 'Medium-term pre-keys for X3DH';
COMMENT ON TABLE onetime_prekeys IS 'Single-use pre-keys for X3DH';
COMMENT ON TABLE messages IS 'Encrypted messages (server cannot decrypt!)';
COMMENT ON TABLE group_chats IS 'Group chat metadata';
COMMENT ON TABLE mls_group_states IS 'MLS group state (ratchet tree)';
COMMENT ON TABLE encrypted_media IS 'Encrypted media files';

COMMENT ON COLUMN messages.ciphertext IS 'Encrypted message content (E2EE)';
COMMENT ON COLUMN messages.ephemeral_key IS 'Ephemeral key for X3DH (1:1 only)';
COMMENT ON COLUMN messages.group_epoch IS 'MLS epoch number (groups only)';
COMMENT ON COLUMN encrypted_media.encryption_key IS 'AES-256 key for media decryption';
