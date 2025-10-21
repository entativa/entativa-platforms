-- Create takes table
CREATE TABLE IF NOT EXISTS takes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_url TEXT NOT NULL,
    thumbnail_url TEXT,
    caption TEXT,
    audio_name VARCHAR(255),
    audio_url TEXT,
    duration INTEGER NOT NULL DEFAULT 0,
    likes_count INTEGER NOT NULL DEFAULT 0,
    comments_count INTEGER NOT NULL DEFAULT 0,
    shares_count INTEGER NOT NULL DEFAULT 0,
    views_count INTEGER NOT NULL DEFAULT 0,
    hashtags JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_takes_user_id (user_id),
    INDEX idx_takes_created_at (created_at DESC),
    INDEX idx_takes_hashtags USING GIN (hashtags)
);

-- Create take_likes table
CREATE TABLE IF NOT EXISTS take_likes (
    take_id UUID NOT NULL REFERENCES takes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (take_id, user_id),
    INDEX idx_take_likes_user_id (user_id),
    INDEX idx_take_likes_created_at (created_at DESC)
);

-- Create take_comments table
CREATE TABLE IF NOT EXISTS take_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    take_id UUID NOT NULL REFERENCES takes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    likes_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_take_comments_take_id (take_id),
    INDEX idx_take_comments_user_id (user_id),
    INDEX idx_take_comments_created_at (created_at DESC)
);

-- Create take_saves table
CREATE TABLE IF NOT EXISTS take_saves (
    take_id UUID NOT NULL REFERENCES takes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (take_id, user_id),
    INDEX idx_take_saves_user_id (user_id)
);

-- Add comments
COMMENT ON TABLE takes IS 'Stores takes/reels videos';
COMMENT ON TABLE take_likes IS 'Tracks take likes';
COMMENT ON TABLE take_comments IS 'Stores comments on takes';
COMMENT ON TABLE take_saves IS 'Tracks saved takes';

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_takes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_takes_updated_at
    BEFORE UPDATE ON takes
    FOR EACH ROW
    EXECUTE FUNCTION update_takes_updated_at();

-- Insert sample data for testing
INSERT INTO takes (id, user_id, video_url, thumbnail_url, caption, audio_name, duration, likes_count, comments_count, shares_count, views_count, hashtags) VALUES
('550e8400-e29b-41d4-a716-446655440001', (SELECT id FROM users LIMIT 1), 'https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4', 'https://sample-videos.com/img/Sample-jpg-image-50kb.jpg', 'Check out this amazing transformation! 💪 #fitness #motivation', 'Original Audio', 30, 45200, 892, 1234, 234500, '["fitness", "motivation"]'),
('550e8400-e29b-41d4-a716-446655440002', (SELECT id FROM users LIMIT 1), 'https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_2mb.mp4', 'https://sample-videos.com/img/Sample-jpg-image-50kb.jpg', 'Best pasta recipe ever! 🍝 Try it and let me know what you think!', 'Cooking Vibes - Sound Library', 45, 78300, 1456, 2890, 456700, '["cooking", "food"]'),
('550e8400-e29b-41d4-a716-446655440003', (SELECT id FROM users LIMIT 1), 'https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_5mb.mp4', 'https://sample-videos.com/img/Sample-jpg-image-50kb.jpg', 'Hidden gems in Bali you NEED to visit! 🌴✨ #travel #bali', 'Tropical Summer - Music Mix', 60, 123400, 3421, 5678, 890200, '["travel", "bali"]');
