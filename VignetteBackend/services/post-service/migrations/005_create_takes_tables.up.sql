-- Create takes table (short-form video content)
CREATE TABLE IF NOT EXISTS takes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    caption TEXT DEFAULT '',
    media_id UUID NOT NULL, -- Single video media
    audio_track_id UUID,
    duration DOUBLE PRECISION NOT NULL,
    thumbnail_url TEXT NOT NULL,
    hashtags JSONB DEFAULT '[]'::jsonb,
    filter_used TEXT,
    location JSONB,
    tagged_user_ids JSONB DEFAULT '[]'::jsonb,
    
    -- Takes-specific
    template_id UUID REFERENCES takes_templates(id) ON DELETE SET NULL,
    trend_id UUID REFERENCES takes_trends(id) ON DELETE SET NULL,
    has_btt BOOLEAN DEFAULT FALSE,
    
    -- Engagement
    views_count BIGINT DEFAULT 0,
    likes_count BIGINT DEFAULT 0,
    comments_count BIGINT DEFAULT 0,
    shares_count BIGINT DEFAULT 0,
    saves_count BIGINT DEFAULT 0,
    remix_count BIGINT DEFAULT 0,
    
    -- Settings
    comments_enabled BOOLEAN DEFAULT TRUE,
    remix_enabled BOOLEAN DEFAULT TRUE,
    
    -- Metadata
    is_sponsored BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT takes_duration_positive CHECK (duration > 0),
    CONSTRAINT takes_counts_non_negative CHECK (
        views_count >= 0 AND likes_count >= 0 AND comments_count >= 0 AND
        shares_count >= 0 AND saves_count >= 0 AND remix_count >= 0
    )
);

-- Create Behind-the-Takes table
CREATE TABLE IF NOT EXISTS behind_the_takes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    take_id UUID NOT NULL UNIQUE REFERENCES takes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    media_ids JSONB NOT NULL,
    description TEXT NOT NULL,
    steps JSONB NOT NULL,
    equipment JSONB DEFAULT '[]'::jsonb,
    software JSONB DEFAULT '[]'::jsonb,
    tips JSONB DEFAULT '[]'::jsonb,
    views_count BIGINT DEFAULT 0,
    likes_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT btt_media_required CHECK (jsonb_array_length(media_ids) > 0),
    CONSTRAINT btt_steps_required CHECK (jsonb_array_length(steps) > 0)
);

-- Create Takes Templates table
CREATE TABLE IF NOT EXISTS takes_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_take_id UUID NOT NULL REFERENCES takes(id) ON DELETE CASCADE,
    creator_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    thumbnail_url TEXT NOT NULL,
    audio_track_id UUID,
    effects JSONB DEFAULT '[]'::jsonb,
    transitions JSONB DEFAULT '[]'::jsonb,
    timing_cues JSONB DEFAULT '[]'::jsonb,
    usage_count BIGINT DEFAULT 0,
    is_public BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT template_usage_non_negative CHECK (usage_count >= 0)
);

-- Create Takes Trends table (with deep-linking to originator)
CREATE TABLE IF NOT EXISTS takes_trends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    keyword VARCHAR(50) NOT NULL UNIQUE, -- Case-insensitive via index
    originator_id UUID NOT NULL, -- Who started the trend
    origin_take_id UUID NOT NULL REFERENCES takes(id) ON DELETE CASCADE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    thumbnail_url TEXT NOT NULL,
    audio_track_id UUID,
    participant_count BIGINT DEFAULT 1, -- Originator counts as first participant
    views_count BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    peak_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT trend_counts_non_negative CHECK (participant_count >= 0 AND views_count >= 0)
);

-- Indexes for Takes
CREATE INDEX idx_takes_user_id ON takes(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_takes_created_desc ON takes(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_takes_template_id ON takes(template_id) WHERE deleted_at IS NULL AND template_id IS NOT NULL;
CREATE INDEX idx_takes_trend_id ON takes(trend_id) WHERE deleted_at IS NULL AND trend_id IS NOT NULL;
CREATE INDEX idx_takes_hashtags ON takes USING GIN(hashtags) WHERE deleted_at IS NULL;

-- Trending Takes algorithm index
CREATE INDEX idx_takes_trending ON takes(
    (views_count / 100 + likes_count * 3 + shares_count * 5 + remix_count * 10) DESC,
    created_at DESC
) WHERE deleted_at IS NULL;

-- Indexes for BTT
CREATE INDEX idx_btt_take_id ON behind_the_takes(take_id);
CREATE INDEX idx_btt_user_id ON behind_the_takes(user_id, created_at DESC);
CREATE INDEX idx_btt_trending ON behind_the_takes(
    (views_count + likes_count * 5) DESC
);

-- Indexes for Templates
CREATE INDEX idx_templates_creator_id ON takes_templates(creator_id, created_at DESC);
CREATE INDEX idx_templates_category ON takes_templates(category, usage_count DESC) WHERE is_public = TRUE;
CREATE INDEX idx_templates_featured ON takes_templates(is_featured, usage_count DESC) WHERE is_public = TRUE;
CREATE INDEX idx_templates_trending ON takes_templates(usage_count DESC, created_at DESC) WHERE is_public = TRUE;
CREATE INDEX idx_templates_search ON takes_templates USING GIN(to_tsvector('english', name || ' ' || description));

-- Indexes for Trends
CREATE UNIQUE INDEX idx_trends_keyword_lower ON takes_trends(LOWER(keyword));
CREATE INDEX idx_trends_active ON takes_trends(is_active, participant_count DESC) WHERE is_active = TRUE;
CREATE INDEX idx_trends_featured ON takes_trends(is_featured, participant_count DESC) WHERE is_featured = TRUE;
CREATE INDEX idx_trends_originator ON takes_trends(originator_id);
CREATE INDEX idx_trends_origin_take ON takes_trends(origin_take_id);

-- Auto-update triggers
CREATE OR REPLACE FUNCTION update_takes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER takes_updated_at_trigger
    BEFORE UPDATE ON takes
    FOR EACH ROW
    EXECUTE FUNCTION update_takes_updated_at();

CREATE TRIGGER btt_updated_at_trigger
    BEFORE UPDATE ON behind_the_takes
    FOR EACH ROW
    EXECUTE FUNCTION update_takes_updated_at();

CREATE TRIGGER templates_updated_at_trigger
    BEFORE UPDATE ON takes_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_takes_updated_at();

CREATE TRIGGER trends_updated_at_trigger
    BEFORE UPDATE ON takes_trends
    FOR EACH ROW
    EXECUTE FUNCTION update_takes_updated_at();

-- Comments
COMMENT ON TABLE takes IS 'Short-form video content (formerly called Reels)';
COMMENT ON COLUMN takes.media_id IS 'Single video media from media service';
COMMENT ON COLUMN takes.template_id IS 'Template used to create this Take';
COMMENT ON COLUMN takes.trend_id IS 'Trend this Take participates in';
COMMENT ON COLUMN takes.has_btt IS 'TRUE if Behind-the-Takes content exists';
COMMENT ON COLUMN takes.remix_count IS 'Times this Take was remixed/used as template';

COMMENT ON TABLE behind_the_takes IS 'Behind-the-scenes content showing how Takes were created';
COMMENT ON COLUMN behind_the_takes.steps IS 'Step-by-step creation process';
COMMENT ON COLUMN behind_the_takes.equipment IS 'Equipment/gear used';
COMMENT ON COLUMN behind_the_takes.software IS 'Apps and tools used';

COMMENT ON TABLE takes_templates IS 'Reusable templates for creating Takes';
COMMENT ON COLUMN takes_templates.effects IS 'Visual effects with timing';
COMMENT ON COLUMN takes_templates.timing_cues IS 'Beat markers and timing cues';
COMMENT ON COLUMN takes_templates.usage_count IS 'Times this template was used';

COMMENT ON TABLE takes_trends IS 'Trending challenges and topics with deep-linking to originators';
COMMENT ON COLUMN takes_trends.keyword IS 'Case-insensitive trend keyword (unique)';
COMMENT ON COLUMN takes_trends.originator_id IS 'User who started this trend (deep-linked)';
COMMENT ON COLUMN takes_trends.origin_take_id IS 'Original Take that started the trend';
COMMENT ON COLUMN takes_trends.participant_count IS 'Number of creators who joined';
