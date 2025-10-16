-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create posts table (Facebook-style)
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    caption TEXT DEFAULT '',
    media_ids JSONB NOT NULL, -- Required for Facebook
    location JSONB,
    tagged_user_ids JSONB DEFAULT '[]'::jsonb,
    hashtags JSONB DEFAULT '[]'::jsonb,
    filter_used TEXT,
    is_carousel BOOLEAN DEFAULT FALSE,
    likes_count BIGINT DEFAULT 0,
    comments_count BIGINT DEFAULT 0,
    views_count BIGINT DEFAULT 0,
    saves_count BIGINT DEFAULT 0,
    shares_count BIGINT DEFAULT 0,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Facebook-specific
    is_sponsored BOOLEAN DEFAULT FALSE,
    is_reels BOOLEAN DEFAULT FALSE,
    comments_enabled BOOLEAN DEFAULT TRUE,
    likes_visible BOOLEAN DEFAULT TRUE,
    
    CONSTRAINT posts_media_required CHECK (jsonb_array_length(media_ids) > 0),
    CONSTRAINT posts_counts_non_negative CHECK (
        likes_count >= 0 AND
        comments_count >= 0 AND
        views_count >= 0 AND
        saves_count >= 0 AND
        shares_count >= 0
    )
);

-- Create indexes for high performance
CREATE INDEX idx_posts_user_id ON posts(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_created_desc ON posts(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_reels ON posts(is_reels, created_at DESC) WHERE deleted_at IS NULL AND is_reels = TRUE;

-- Trending/Explore algorithm index
CREATE INDEX idx_posts_trending ON posts(
    (likes_count + views_count/10 + saves_count * 2 + shares_count * 3 + comments_count * 2) DESC,
    created_at DESC
) WHERE deleted_at IS NULL;

-- GIN index for media_ids array search
CREATE INDEX idx_posts_media_ids ON posts USING GIN(media_ids) WHERE deleted_at IS NULL;

-- GIN index for tagged users
CREATE INDEX idx_posts_tagged_users ON posts USING GIN(tagged_user_ids) WHERE deleted_at IS NULL;

-- GIN index for hashtags (Facebook's core feature)
CREATE INDEX idx_posts_hashtags ON posts USING GIN(hashtags) WHERE deleted_at IS NULL;

-- Full text search for caption
CREATE INDEX idx_posts_caption_search ON posts USING GIN(to_tsvector('english', caption)) WHERE deleted_at IS NULL;

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_posts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_updated_at_trigger
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_posts_updated_at();

-- Add table comments for documentation
COMMENT ON TABLE posts IS 'Facebook-style posts with required media attachments';
COMMENT ON COLUMN posts.media_ids IS 'Array of media IDs (at least 1 required)';
COMMENT ON COLUMN posts.hashtags IS 'Extracted hashtags from caption';
COMMENT ON COLUMN posts.is_carousel IS 'TRUE if post has multiple images';
COMMENT ON COLUMN posts.is_reels IS 'TRUE if post is a Reel (short video)';
COMMENT ON COLUMN posts.filter_used IS 'Facebook filter applied to media';
COMMENT ON COLUMN posts.views_count IS 'Number of views (important for Reels)';
COMMENT ON COLUMN posts.saves_count IS 'Number of saves/bookmarks';
