-- Live Streaming Service Database Schema
-- YouTube-quality live streaming with follower thresholds

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- LIVE STREAMS TABLE
-- ============================================
CREATE TABLE live_streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    streamer_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    thumbnail_url TEXT,
    
    -- Stream configuration
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'live', 'ended', 'cancelled')),
    quality VARCHAR(10) NOT NULL DEFAULT '720p' CHECK (quality IN ('144p', '240p', '360p', '480p', '720p', '1080p', '1440p', '2160p')),
    is_private BOOLEAN DEFAULT FALSE,
    category VARCHAR(50) NOT NULL,
    tags TEXT[] DEFAULT '{}',
    
    -- Technical details
    stream_key VARCHAR(255) NOT NULL UNIQUE,
    rtmp_url TEXT NOT NULL,
    hls_url TEXT NOT NULL,
    webrtc_url TEXT NOT NULL,
    
    -- Analytics (denormalized for performance)
    viewer_count INTEGER DEFAULT 0,
    peak_viewers INTEGER DEFAULT 0,
    total_views INTEGER DEFAULT 0,
    likes_count INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    shares_count INTEGER DEFAULT 0,
    
    -- Recording
    record_stream BOOLEAN DEFAULT TRUE,
    recording_url TEXT,
    
    -- Timestamps
    scheduled_for TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER DEFAULT 0, -- Seconds
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_streams_streamer ON live_streams(streamer_id, status);
CREATE INDEX idx_streams_status ON live_streams(status, created_at DESC);
CREATE INDEX idx_streams_category ON live_streams(category, status);
CREATE INDEX idx_streams_created ON live_streams(created_at DESC);
CREATE INDEX idx_streams_live ON live_streams(status) WHERE status = 'live';

-- ============================================
-- STREAM VIEWERS TABLE
-- ============================================
CREATE TABLE stream_viewers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    viewer_id UUID NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    watch_time INTEGER DEFAULT 0, -- Seconds
    is_active BOOLEAN DEFAULT TRUE,
    
    UNIQUE(stream_id, viewer_id)
);

CREATE INDEX idx_viewers_stream ON stream_viewers(stream_id, is_active);
CREATE INDEX idx_viewers_user ON stream_viewers(viewer_id);
CREATE INDEX idx_viewers_active ON stream_viewers(stream_id) WHERE is_active = TRUE;

-- ============================================
-- STREAM COMMENTS TABLE
-- ============================================
CREATE TABLE stream_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    content TEXT NOT NULL CHECK (length(content) <= 500),
    is_pinned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_comments_stream ON stream_comments(stream_id, created_at DESC);
CREATE INDEX idx_comments_user ON stream_comments(user_id);
CREATE INDEX idx_comments_pinned ON stream_comments(stream_id, is_pinned) WHERE is_pinned = TRUE;

-- ============================================
-- STREAM REACTIONS TABLE
-- ============================================
CREATE TABLE stream_reactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('like', 'love', 'fire', 'clap', 'wow')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(stream_id, user_id) -- One reaction per user per stream
);

CREATE INDEX idx_reactions_stream ON stream_reactions(stream_id, type);
CREATE INDEX idx_reactions_user ON stream_reactions(user_id);

-- ============================================
-- STREAM ANALYTICS TABLE (Snapshot)
-- ============================================
CREATE TABLE stream_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    snapshot_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Metrics at time of snapshot
    current_viewers INTEGER DEFAULT 0,
    total_views INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    reactions_count INTEGER DEFAULT 0,
    
    UNIQUE(stream_id, snapshot_at)
);

CREATE INDEX idx_analytics_stream ON stream_analytics(stream_id, snapshot_at DESC);

-- ============================================
-- TRIGGERS
-- ============================================

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER live_streams_updated_at
    BEFORE UPDATE ON live_streams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Auto-increment comment count
CREATE OR REPLACE FUNCTION increment_comment_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE live_streams SET comments_count = comments_count + 1 WHERE id = NEW.stream_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER stream_comments_count
    AFTER INSERT ON stream_comments
    FOR EACH ROW EXECUTE FUNCTION increment_comment_count();

-- Auto-increment viewer count on join
CREATE OR REPLACE FUNCTION update_viewer_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.is_active THEN
        UPDATE live_streams 
        SET viewer_count = viewer_count + 1,
            total_views = total_views + 1,
            peak_viewers = GREATEST(peak_viewers, viewer_count + 1)
        WHERE id = NEW.stream_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.is_active AND NOT NEW.is_active THEN
        UPDATE live_streams 
        SET viewer_count = GREATEST(0, viewer_count - 1)
        WHERE id = NEW.stream_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER stream_viewers_count
    AFTER INSERT OR UPDATE ON stream_viewers
    FOR EACH ROW EXECUTE FUNCTION update_viewer_count();

-- Comments
COMMENT ON TABLE live_streams IS 'Live streaming sessions with YouTube-quality support';
COMMENT ON COLUMN live_streams.stream_key IS 'Secret key for RTMP streaming (never expose in API)';
COMMENT ON COLUMN live_streams.quality IS 'Stream quality: 144p to 4K (2160p)';
COMMENT ON TABLE stream_viewers IS 'Track viewers watching streams in real-time';
COMMENT ON TABLE stream_comments IS 'Real-time comments during live streams';
COMMENT ON TABLE stream_reactions IS 'Real-time reactions (like, love, fire, etc)';
