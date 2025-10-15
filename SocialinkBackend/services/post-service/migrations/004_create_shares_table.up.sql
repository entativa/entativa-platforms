-- Create shares table
CREATE TABLE IF NOT EXISTS shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    original_post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    caption TEXT,
    privacy privacy NOT NULL DEFAULT 'public',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT shares_no_duplicate UNIQUE(user_id, original_post_id)
);

-- Create indexes
CREATE INDEX idx_shares_user_id ON shares(user_id, created_at DESC);
CREATE INDEX idx_shares_original_post ON shares(original_post_id, created_at DESC);
CREATE INDEX idx_shares_created_desc ON shares(created_at DESC);

-- Add table comments
COMMENT ON TABLE shares IS 'Shared posts on user timelines';
COMMENT ON COLUMN shares.caption IS 'Optional caption added when sharing';
COMMENT ON CONSTRAINT shares_no_duplicate ON shares IS 'Prevents duplicate shares of same post by same user';
