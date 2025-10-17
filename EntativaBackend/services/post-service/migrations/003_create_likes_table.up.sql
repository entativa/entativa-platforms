-- Create likes table
CREATE TABLE IF NOT EXISTS likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    reaction_type reaction_type NOT NULL DEFAULT 'like',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT likes_target_check CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR
        (post_id IS NULL AND comment_id IS NOT NULL)
    )
);

-- Create unique indexes to prevent duplicate likes
CREATE UNIQUE INDEX idx_likes_user_post ON likes(user_id, post_id) 
    WHERE post_id IS NOT NULL AND comment_id IS NULL;

CREATE UNIQUE INDEX idx_likes_user_comment ON likes(user_id, comment_id) 
    WHERE comment_id IS NOT NULL AND post_id IS NULL;

-- Indexes for querying
CREATE INDEX idx_likes_post_id ON likes(post_id, created_at DESC) WHERE post_id IS NOT NULL;
CREATE INDEX idx_likes_comment_id ON likes(comment_id, created_at DESC) WHERE comment_id IS NOT NULL;
CREATE INDEX idx_likes_user_id ON likes(user_id, created_at DESC);
CREATE INDEX idx_likes_reaction_type ON likes(reaction_type, created_at DESC);

-- Add table comments
COMMENT ON TABLE likes IS 'Likes and reactions on posts and comments';
COMMENT ON COLUMN likes.reaction_type IS 'Type of reaction (like, love, haha, wow, sad, angry, care)';
COMMENT ON CONSTRAINT likes_target_check ON likes IS 'Ensures like targets either post or comment, not both';
