-- Creator Service Database Schema
-- Instagram-style creator tools with analytics, monetization, insights

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- CREATOR PROFILES TABLE
-- ============================================
CREATE TABLE creator_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    account_type VARCHAR(20) NOT NULL DEFAULT 'personal' CHECK (account_type IN ('personal', 'business', 'creator')),
    
    -- Creator info
    display_name VARCHAR(100) NOT NULL,
    bio TEXT,
    category VARCHAR(50) NOT NULL,
    badges JSONB DEFAULT '[]',
    
    -- Contact
    email VARCHAR(255),
    phone VARCHAR(50),
    website TEXT,
    
    -- Monetization
    monetization_enabled BOOLEAN DEFAULT FALSE,
    monetization_status VARCHAR(20) DEFAULT 'pending' CHECK (monetization_status IN ('pending', 'approved', 'rejected', 'suspended')),
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_creator_profiles_user ON creator_profiles(user_id);
CREATE INDEX idx_creator_profiles_category ON creator_profiles(category);
CREATE INDEX idx_creator_profiles_monetization ON creator_profiles(monetization_enabled, monetization_status);

-- ============================================
-- CREATOR ANALYTICS TABLE (Daily aggregates)
-- ============================================
CREATE TABLE creator_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES creator_profiles(user_id) ON DELETE CASCADE,
    date DATE NOT NULL,
    
    -- Followers
    followers_count INTEGER DEFAULT 0,
    followers_gained INTEGER DEFAULT 0,
    followers_lost INTEGER DEFAULT 0,
    
    -- Engagement
    total_likes INTEGER DEFAULT 0,
    total_comments INTEGER DEFAULT 0,
    total_shares INTEGER DEFAULT 0,
    total_views INTEGER DEFAULT 0,
    engagement_rate DOUBLE PRECISION DEFAULT 0,
    
    -- Content
    posts_count INTEGER DEFAULT 0,
    takes_count INTEGER DEFAULT 0,
    stories_count INTEGER DEFAULT 0,
    
    -- Reach
    accounts_reached INTEGER DEFAULT 0,
    accounts_engaged INTEGER DEFAULT 0,
    
    -- Demographics (JSONB for flexibility)
    age_gender_breakdown JSONB DEFAULT '{}',
    top_locations JSONB DEFAULT '{}',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, date)
);

CREATE INDEX idx_analytics_user_date ON creator_analytics(user_id, date DESC);
CREATE INDEX idx_analytics_date ON creator_analytics(date DESC);

-- ============================================
-- CONTENT INSIGHTS TABLE
-- ============================================
CREATE TABLE content_insights (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL,
    content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('post', 'take', 'story')),
    
    -- Performance
    impressions INTEGER DEFAULT 0,
    reach INTEGER DEFAULT 0,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    saves INTEGER DEFAULT 0,
    engagement INTEGER DEFAULT 0,
    engagement_rate DOUBLE PRECISION DEFAULT 0,
    
    -- Traffic sources
    from_home INTEGER DEFAULT 0,
    from_explore INTEGER DEFAULT 0,
    from_profile INTEGER DEFAULT 0,
    from_hashtags INTEGER DEFAULT 0,
    from_other INTEGER DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(content_id)
);

CREATE INDEX idx_insights_content ON content_insights(content_id);
CREATE INDEX idx_insights_type ON content_insights(content_type, engagement_rate DESC);
CREATE INDEX idx_insights_created ON content_insights(created_at DESC);

-- ============================================
-- MONETIZATION APPLICATIONS TABLE
-- ============================================
CREATE TABLE monetization_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES creator_profiles(user_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'suspended')),
    
    -- Requirements check
    followers_count INTEGER NOT NULL,
    posts_count INTEGER NOT NULL,
    meets_requirements BOOLEAN DEFAULT FALSE,
    
    -- Tax and payout
    tax_id VARCHAR(100),
    payout_method VARCHAR(50),
    
    -- Review
    reviewed_at TIMESTAMP WITH TIME ZONE,
    reviewed_by UUID,
    rejection_reason TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_monetization_user ON monetization_applications(user_id);
CREATE INDEX idx_monetization_status ON monetization_applications(status, created_at DESC);

-- ============================================
-- CREATOR EARNINGS TABLE
-- ============================================
CREATE TABLE creator_earnings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES creator_profiles(user_id) ON DELETE CASCADE,
    month DATE NOT NULL,
    
    -- Revenue streams
    ads_revenue DECIMAL(10, 2) DEFAULT 0,
    tips_revenue DECIMAL(10, 2) DEFAULT 0,
    brand_deals_revenue DECIMAL(10, 2) DEFAULT 0,
    other_revenue DECIMAL(10, 2) DEFAULT 0,
    total_revenue DECIMAL(10, 2) DEFAULT 0,
    
    -- Payout status
    is_paid BOOLEAN DEFAULT FALSE,
    paid_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, month)
);

CREATE INDEX idx_earnings_user ON creator_earnings(user_id, month DESC);
CREATE INDEX idx_earnings_unpaid ON creator_earnings(is_paid, month) WHERE NOT is_paid;

-- ============================================
-- CREATOR BADGES TABLE (Badge history)
-- ============================================
CREATE TABLE creator_badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES creator_profiles(user_id) ON DELETE CASCADE,
    badge VARCHAR(50) NOT NULL CHECK (badge IN ('verified', 'partner', 'top_creator', 'trendsetter', 'rising')),
    awarded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_badges_user ON creator_badges(user_id, is_active);
CREATE INDEX idx_badges_active ON creator_badges(user_id) WHERE is_active = TRUE;

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

CREATE TRIGGER creator_profiles_updated_at
    BEFORE UPDATE ON creator_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER content_insights_updated_at
    BEFORE UPDATE ON content_insights
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER monetization_applications_updated_at
    BEFORE UPDATE ON monetization_applications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Auto-calculate engagement rate
CREATE OR REPLACE FUNCTION calculate_engagement_rate()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.reach > 0 THEN
        NEW.engagement_rate = (NEW.likes + NEW.comments + NEW.shares + NEW.saves)::DOUBLE PRECISION / NEW.reach * 100;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER content_insights_engagement_rate
    BEFORE INSERT OR UPDATE ON content_insights
    FOR EACH ROW EXECUTE FUNCTION calculate_engagement_rate();

-- Auto-calculate total revenue
CREATE OR REPLACE FUNCTION calculate_total_revenue()
RETURNS TRIGGER AS $$
BEGIN
    NEW.total_revenue = NEW.ads_revenue + NEW.tips_revenue + NEW.brand_deals_revenue + NEW.other_revenue;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER creator_earnings_total
    BEFORE INSERT OR UPDATE ON creator_earnings
    FOR EACH ROW EXECUTE FUNCTION calculate_total_revenue();

-- Comments
COMMENT ON TABLE creator_profiles IS 'Instagram-style creator professional accounts';
COMMENT ON TABLE creator_analytics IS 'Daily analytics aggregates for creators';
COMMENT ON TABLE content_insights IS 'Performance insights for individual posts/takes/stories';
COMMENT ON TABLE monetization_applications IS 'Applications for creator monetization program';
COMMENT ON TABLE creator_earnings IS 'Monthly earnings tracking for creators';
COMMENT ON TABLE creator_badges IS 'Creator badge awards (verified, partner, etc)';
