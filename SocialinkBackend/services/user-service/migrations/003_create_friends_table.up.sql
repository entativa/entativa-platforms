-- Socialink Friend System (Facebook-style but better!)
-- Friend requests must be accepted
-- Reasonable limit (not 5,000 like Facebook!)

CREATE TABLE IF NOT EXISTS friend_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL,      -- Who sent the request
    receiver_id UUID NOT NULL,    -- Who receives the request
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'cancelled')),
    message TEXT,                 -- Optional message with request
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    responded_at TIMESTAMP WITH TIME ZONE,
    
    -- Prevent duplicate requests
    UNIQUE(sender_id, receiver_id),
    
    -- Prevent self-requests
    CHECK (sender_id != receiver_id)
);

CREATE TABLE IF NOT EXISTS friends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id_1 UUID NOT NULL,      -- First user (lower ID)
    user_id_2 UUID NOT NULL,      -- Second user (higher ID)
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'removed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure user_id_1 < user_id_2 for consistency
    CHECK (user_id_1 < user_id_2),
    
    -- Prevent duplicate friendships
    UNIQUE(user_id_1, user_id_2)
);

CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID NOT NULL,    -- Who is following
    following_id UUID NOT NULL,   -- Who is being followed (pages, public figures)
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'removed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

-- Indexes for friend requests
CREATE INDEX idx_friend_requests_sender ON friend_requests(sender_id, status);
CREATE INDEX idx_friend_requests_receiver ON friend_requests(receiver_id, status);
CREATE INDEX idx_friend_requests_created ON friend_requests(created_at DESC);

-- Indexes for friends
CREATE INDEX idx_friends_user1 ON friends(user_id_1, status);
CREATE INDEX idx_friends_user2 ON friends(user_id_2, status);
CREATE INDEX idx_friends_created ON friends(created_at DESC);

-- Indexes for follows
CREATE INDEX idx_follows_follower ON follows(follower_id, status);
CREATE INDEX idx_follows_following ON follows(following_id, status);

-- Comments
COMMENT ON TABLE friend_requests IS 'Socialink friend requests (must be accepted, not instant like Instagram)';
COMMENT ON TABLE friends IS 'Socialink friend relationships (bi-directional, limited to prevent spam)';
COMMENT ON TABLE follows IS 'Socialink follows for pages/public figures (separate from friends)';
