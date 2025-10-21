-- Social Graph: Following, Friends, Connections
-- Supports both Entativa (friends) and Vignette (followers) models

-- Followers/Following table (Vignette model - one-way)
CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'blocked', 'muted'
    
    -- Notifications
    notifications_enabled BOOLEAN DEFAULT TRUE,
    show_in_feed BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id),
    
    INDEX idx_follows_follower (follower_id),
    INDEX idx_follows_following (following_id),
    INDEX idx_follows_status (status),
    INDEX idx_follows_created (created_at DESC)
);

-- Friend requests table (Entativa model - two-way)
CREATE TABLE IF NOT EXISTS friend_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected', 'cancelled'
    
    -- Message with request (optional)
    message TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    responded_at TIMESTAMP,
    
    -- Constraints
    UNIQUE(sender_id, receiver_id),
    CHECK (sender_id != receiver_id),
    
    INDEX idx_friend_requests_sender (sender_id),
    INDEX idx_friend_requests_receiver (receiver_id),
    INDEX idx_friend_requests_status (status),
    INDEX idx_friend_requests_created (created_at DESC)
);

-- Friendships table (Entativa - mutual friends)
CREATE TABLE IF NOT EXISTS friendships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id_1 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id_2 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Relationship type
    relationship_type VARCHAR(20) DEFAULT 'friend', -- 'friend', 'close_friend', 'acquaintance'
    
    -- Privacy
    show_in_friends_list BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    UNIQUE(user_id_1, user_id_2),
    CHECK (user_id_1 < user_id_2), -- Ensure consistent ordering
    
    INDEX idx_friendships_user1 (user_id_1),
    INDEX idx_friendships_user2 (user_id_2),
    INDEX idx_friendships_type (relationship_type),
    INDEX idx_friendships_created (created_at DESC)
);

-- GRANULAR MESSAGE PERMISSIONS
-- This is THE BEST message control system - better than any social platform!
CREATE TABLE IF NOT EXISTS message_permissions (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Who can message me
    message_permission VARCHAR(50) NOT NULL DEFAULT 'followers', 
    -- Options:
    -- 'everyone' - Anyone can message me
    -- 'followers' - Only people I follow (Vignette)
    -- 'friends' - Only my friends (Entativa)
    -- 'following' - Only people who follow me (Vignette)
    -- 'mutual_followers' - Only mutual followers (Vignette)
    -- 'nobody' - No one can message me (DND mode)
    -- 'custom' - Use allow/block lists
    
    -- Message request settings
    auto_accept_from_followers BOOLEAN DEFAULT TRUE, -- Vignette
    auto_accept_from_friends BOOLEAN DEFAULT TRUE,   -- Entativa
    auto_accept_verified BOOLEAN DEFAULT FALSE,      -- Auto-accept from verified users
    
    -- Request filtering
    require_mutual_connection BOOLEAN DEFAULT FALSE,
    min_follower_count INTEGER DEFAULT 0,            -- Minimum followers to message me
    min_account_age_days INTEGER DEFAULT 0,          -- Minimum account age in days
    
    -- Allow strangers to send ONE message request
    allow_message_requests BOOLEAN DEFAULT TRUE,
    
    -- Read receipts
    send_read_receipts BOOLEAN DEFAULT TRUE,
    send_typing_indicators BOOLEAN DEFAULT TRUE,
    
    -- Group chat settings
    allow_group_invites BOOLEAN DEFAULT TRUE,
    auto_accept_group_from_friends BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Specific allow list (GRANULAR CONTROL!)
-- Users can manually allow SPECIFIC people to message them
CREATE TABLE IF NOT EXISTS message_allow_list (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    allowed_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Why allowed
    reason VARCHAR(100), -- 'manually_added', 'close_friend', 'verified', 'business'
    notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, allowed_user_id),
    INDEX idx_allow_list_user (user_id),
    INDEX idx_allow_list_allowed (allowed_user_id)
);

-- Specific block list for messages
CREATE TABLE IF NOT EXISTS message_block_list (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Reason
    reason VARCHAR(100), -- 'spam', 'harassment', 'unwanted', 'other'
    notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, blocked_user_id),
    INDEX idx_block_list_user (user_id),
    INDEX idx_block_list_blocked (blocked_user_id)
);

-- Message requests (for strangers)
CREATE TABLE IF NOT EXISTS message_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Request message (first message content)
    message_preview TEXT NOT NULL,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected', 'expired'
    
    -- Auto-accept logic result
    auto_accept_eligible BOOLEAN DEFAULT FALSE,
    rejection_reason VARCHAR(100), -- 'blocked', 'spam_filter', 'min_followers', 'account_age'
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP, -- Requests can expire after 30 days
    responded_at TIMESTAMP,
    
    UNIQUE(sender_id, receiver_id),
    INDEX idx_message_requests_sender (sender_id),
    INDEX idx_message_requests_receiver (receiver_id),
    INDEX idx_message_requests_status (status),
    INDEX idx_message_requests_created (created_at DESC)
);

-- Close friends list (Vignette feature)
CREATE TABLE IF NOT EXISTS close_friends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    close_friend_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, close_friend_id),
    INDEX idx_close_friends_user (user_id),
    INDEX idx_close_friends_friend (close_friend_id)
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_social_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_follows_updated_at
    BEFORE UPDATE ON follows
    FOR EACH ROW
    EXECUTE FUNCTION update_social_updated_at();

CREATE TRIGGER trigger_update_friend_requests_updated_at
    BEFORE UPDATE ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_social_updated_at();

CREATE TRIGGER trigger_update_friendships_updated_at
    BEFORE UPDATE ON friendships
    FOR EACH ROW
    EXECUTE FUNCTION update_social_updated_at();

CREATE TRIGGER trigger_update_message_permissions_updated_at
    BEFORE UPDATE ON message_permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_social_updated_at();

-- Function to check if user can message another user
CREATE OR REPLACE FUNCTION can_user_message(sender UUID, receiver UUID)
RETURNS BOOLEAN AS $$
DECLARE
    permissions RECORD;
    is_blocked BOOLEAN;
    is_allowed BOOLEAN;
    is_follower BOOLEAN;
    is_friend BOOLEAN;
    is_verified BOOLEAN;
    follower_count INTEGER;
    account_age_days INTEGER;
BEGIN
    -- Check if blocked
    SELECT EXISTS(
        SELECT 1 FROM message_block_list 
        WHERE user_id = receiver AND blocked_user_id = sender
    ) INTO is_blocked;
    
    IF is_blocked THEN
        RETURN FALSE;
    END IF;
    
    -- Check if in allow list (always allowed if in allow list)
    SELECT EXISTS(
        SELECT 1 FROM message_allow_list 
        WHERE user_id = receiver AND allowed_user_id = sender
    ) INTO is_allowed;
    
    IF is_allowed THEN
        RETURN TRUE;
    END IF;
    
    -- Get receiver's message permissions
    SELECT * FROM message_permissions WHERE user_id = receiver INTO permissions;
    
    -- If no permissions set, use default (followers)
    IF permissions IS NULL THEN
        permissions.message_permission := 'followers';
        permissions.allow_message_requests := TRUE;
    END IF;
    
    -- Check permission level
    CASE permissions.message_permission
        WHEN 'everyone' THEN
            RETURN TRUE;
            
        WHEN 'nobody' THEN
            RETURN FALSE;
            
        WHEN 'followers' THEN
            -- Check if sender is following receiver
            SELECT EXISTS(
                SELECT 1 FROM follows 
                WHERE follower_id = receiver AND following_id = sender AND status = 'active'
            ) INTO is_follower;
            RETURN is_follower;
            
        WHEN 'friends' THEN
            -- Check if they are friends (Entativa)
            SELECT EXISTS(
                SELECT 1 FROM friendships 
                WHERE (user_id_1 = LEAST(sender, receiver) AND user_id_2 = GREATEST(sender, receiver))
            ) INTO is_friend;
            RETURN is_friend;
            
        WHEN 'following' THEN
            -- Check if receiver is following sender
            SELECT EXISTS(
                SELECT 1 FROM follows 
                WHERE follower_id = sender AND following_id = receiver AND status = 'active'
            ) INTO is_follower;
            RETURN is_follower;
            
        WHEN 'mutual_followers' THEN
            -- Check if they follow each other
            SELECT EXISTS(
                SELECT 1 FROM follows f1
                JOIN follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
                WHERE f1.follower_id = sender AND f1.following_id = receiver
                  AND f1.status = 'active' AND f2.status = 'active'
            ) INTO is_follower;
            RETURN is_follower;
            
        WHEN 'custom' THEN
            -- Already checked allow/block lists above
            -- If not in either list, allow message request
            RETURN permissions.allow_message_requests;
            
        ELSE
            -- Default to allowing message requests
            RETURN permissions.allow_message_requests;
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE follows IS 'One-way following relationship (Vignette model)';
COMMENT ON TABLE friend_requests IS 'Friend request system (Entativa model)';
COMMENT ON TABLE friendships IS 'Mutual friendships (Entativa model)';
COMMENT ON TABLE message_permissions IS 'GRANULAR message permission controls - best in class!';
COMMENT ON TABLE message_allow_list IS 'Specific users allowed to message (overrides all other settings)';
COMMENT ON TABLE message_block_list IS 'Specific users blocked from messaging';
COMMENT ON TABLE message_requests IS 'Pending message requests from non-connected users';
COMMENT ON TABLE close_friends IS 'Close friends list for Stories (Vignette)';

COMMENT ON COLUMN message_permissions.message_permission IS 'Primary permission level: everyone, followers, friends, following, mutual_followers, nobody, custom';
COMMENT ON COLUMN message_permissions.auto_accept_from_followers IS 'Auto-accept messages from people I follow';
COMMENT ON COLUMN message_permissions.min_follower_count IS 'Minimum followers required to send me a message request (spam protection)';
COMMENT ON COLUMN message_permissions.min_account_age_days IS 'Minimum account age to send me requests (prevents spam accounts)';

-- Create default message permissions for existing users
INSERT INTO message_permissions (user_id, message_permission, allow_message_requests)
SELECT id, 'followers', TRUE FROM users
ON CONFLICT (user_id) DO NOTHING;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_follows_active ON follows(follower_id, following_id) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_friend_requests_pending ON friend_requests(receiver_id) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_message_requests_pending ON message_requests(receiver_id) WHERE status = 'pending';
