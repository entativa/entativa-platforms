-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    media_id UUID,
    likes_count BIGINT DEFAULT 0,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT comments_content_not_empty CHECK (LENGTH(TRIM(content)) > 0),
    CONSTRAINT comments_likes_non_negative CHECK (likes_count >= 0)
);

-- Create indexes for performance
CREATE INDEX idx_comments_post_id ON comments(post_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_comments_user_id ON comments(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_comments_parent_id ON comments(parent_id, created_at ASC) WHERE deleted_at IS NULL AND parent_id IS NOT NULL;
CREATE INDEX idx_comments_created_desc ON comments(created_at DESC) WHERE deleted_at IS NULL;

-- Full text search for comments
CREATE INDEX idx_comments_content_search ON comments USING GIN(to_tsvector('english', content)) WHERE deleted_at IS NULL;

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_comments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER comments_updated_at_trigger
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_comments_updated_at();

-- Add table comments
COMMENT ON TABLE comments IS 'Comments on posts, supporting nested replies';
COMMENT ON COLUMN comments.parent_id IS 'Parent comment ID for nested replies';
COMMENT ON COLUMN comments.media_id IS 'Optional media attachment (image/GIF)';
