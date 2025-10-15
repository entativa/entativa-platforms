-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create privacy enum
CREATE TYPE privacy AS ENUM (
    'public',
    'friends',
    'friends_except',
    'specific_friends',
    'only_me',
    'custom'
);

-- Create reaction_type enum
CREATE TYPE reaction_type AS ENUM (
    'like',
    'love',
    'haha',
    'wow',
    'sad',
    'angry',
    'care'
);

-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    media_ids JSONB DEFAULT '[]'::jsonb,
    privacy privacy NOT NULL DEFAULT 'public',
    location TEXT,
    tagged_user_ids JSONB DEFAULT '[]'::jsonb,
    feeling TEXT,
    activity TEXT,
    likes_count BIGINT DEFAULT 0,
    comments_count BIGINT DEFAULT 0,
    shares_count BIGINT DEFAULT 0,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT posts_content_not_empty CHECK (LENGTH(TRIM(content)) > 0 OR jsonb_array_length(media_ids) > 0),
    CONSTRAINT posts_counts_non_negative CHECK (
        likes_count >= 0 AND
        comments_count >= 0 AND
        shares_count >= 0
    )
);

-- Create indexes for high performance
CREATE INDEX idx_posts_user_id ON posts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_created_desc ON posts(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_privacy ON posts(privacy) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_user_created ON posts(user_id, created_at DESC) WHERE deleted_at IS NULL;

-- Partial index for trending posts
CREATE INDEX idx_posts_trending ON posts(
    (likes_count * 2 + comments_count * 3 + shares_count * 5) DESC,
    created_at DESC
) WHERE deleted_at IS NULL AND privacy = 'public';

-- GIN index for media_ids array search
CREATE INDEX idx_posts_media_ids ON posts USING GIN(media_ids) WHERE deleted_at IS NULL;

-- GIN index for tagged users
CREATE INDEX idx_posts_tagged_users ON posts USING GIN(tagged_user_ids) WHERE deleted_at IS NULL;

-- Full text search index for content
CREATE INDEX idx_posts_content_search ON posts USING GIN(to_tsvector('english', content)) WHERE deleted_at IS NULL;

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
COMMENT ON TABLE posts IS 'Social media posts with text and media attachments';
COMMENT ON COLUMN posts.media_ids IS 'Array of media IDs from media service';
COMMENT ON COLUMN posts.tagged_user_ids IS 'Array of user IDs tagged in the post';
COMMENT ON COLUMN posts.privacy IS 'Privacy setting for post visibility';
COMMENT ON COLUMN posts.feeling IS 'User feeling/emotion (e.g., happy, excited)';
COMMENT ON COLUMN posts.activity IS 'User activity (e.g., watching, eating)';
