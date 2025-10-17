-- Create saves table (Facebook bookmark feature)
CREATE TABLE IF NOT EXISTS saves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    collection TEXT DEFAULT 'all',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT saves_no_duplicate UNIQUE(user_id, post_id)
);

-- Create indexes
CREATE INDEX idx_saves_user_id ON saves(user_id, created_at DESC);
CREATE INDEX idx_saves_post_id ON saves(post_id, created_at DESC);
CREATE INDEX idx_saves_collection ON saves(user_id, collection, created_at DESC);

-- Add table comments
COMMENT ON TABLE saves IS 'Saved/bookmarked posts by users (Facebook feature)';
COMMENT ON COLUMN saves.collection IS 'Optional collection name for organizing saved posts';
