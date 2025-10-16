-- Vignette Follow System (Instagram-style)
-- Simple follow/unfollow, no friend requests

CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID NOT NULL,  -- Who is following
    following_id UUID NOT NULL, -- Who is being followed
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'removed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Prevent duplicate follows
    UNIQUE(follower_id, following_id),
    
    -- Prevent self-follows
    CHECK (follower_id != following_id)
);

-- Indexes for performance
CREATE INDEX idx_follows_follower ON follows(follower_id, status);
CREATE INDEX idx_follows_following ON follows(following_id, status);
CREATE INDEX idx_follows_created ON follows(created_at DESC);

-- Composite index for quick lookups
CREATE INDEX idx_follows_relationship ON follows(follower_id, following_id) WHERE status = 'active';

-- Comments
COMMENT ON TABLE follows IS 'Vignette follow relationships (Instagram-style, no requests needed)';
COMMENT ON COLUMN follows.follower_id IS 'User who is following';
COMMENT ON COLUMN follows.following_id IS 'User being followed';
