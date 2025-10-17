-- Community Service Database Schema
-- Comprehensive tables for robust community management

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- COMMUNITIES TABLE
-- ============================================
CREATE TABLE communities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cover_photo TEXT, -- ONLY cover photo, NO profile photo!
    category VARCHAR(50) NOT NULL,
    
    -- Privacy & Visibility
    privacy VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (privacy IN ('public', 'private', 'hidden')),
    visibility VARCHAR(20) NOT NULL DEFAULT 'listed' CHECK (visibility IN ('listed', 'unlisted')),
    
    -- Settings
    is_verified BOOLEAN DEFAULT FALSE,
    allow_posts BOOLEAN DEFAULT TRUE,
    require_approval BOOLEAN DEFAULT FALSE,
    
    -- Ownership
    creator_id UUID NOT NULL,
    
    -- Stats (denormalized for performance)
    member_count INTEGER DEFAULT 1,
    post_count INTEGER DEFAULT 0,
    
    -- Metadata
    tags TEXT[] DEFAULT '{}',
    location VARCHAR(255),
    website VARCHAR(255),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_communities_creator ON communities(creator_id);
CREATE INDEX idx_communities_category ON communities(category);
CREATE INDEX idx_communities_privacy ON communities(privacy);
CREATE INDEX idx_communities_created_at ON communities(created_at DESC);
CREATE INDEX idx_communities_member_count ON communities(member_count DESC);
CREATE INDEX idx_communities_tags ON communities USING GIN(tags);

-- Full-text search on name and description
CREATE INDEX idx_communities_search ON communities USING GIN(to_tsvector('english', name || ' ' || COALESCE(description, '')));

-- ============================================
-- COMMUNITY MEMBERS TABLE
-- ============================================
CREATE TABLE community_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'moderator', 'member')),
    permissions JSONB NOT NULL DEFAULT '{}',
    
    -- Member status
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'pending', 'banned')),
    is_muted BOOLEAN DEFAULT FALSE,
    muted_until TIMESTAMP WITH TIME ZONE,
    
    -- Activity counters
    post_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    
    -- Metadata
    invited_by UUID,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(community_id, user_id)
);

CREATE INDEX idx_members_community ON community_members(community_id);
CREATE INDEX idx_members_user ON community_members(user_id);
CREATE INDEX idx_members_role ON community_members(community_id, role);
CREATE INDEX idx_members_status ON community_members(status);
CREATE INDEX idx_members_joined_at ON community_members(joined_at DESC);

-- ============================================
-- JOIN REQUESTS TABLE (for private communities)
-- ============================================
CREATE TABLE join_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    message TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    
    UNIQUE(community_id, user_id, status) -- Prevent duplicate pending requests
);

CREATE INDEX idx_join_requests_community ON join_requests(community_id, status);
CREATE INDEX idx_join_requests_user ON join_requests(user_id);
CREATE INDEX idx_join_requests_created ON join_requests(created_at DESC);

-- ============================================
-- MEMBER INVITES TABLE
-- ============================================
CREATE TABLE member_invites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    invited_user_id UUID NOT NULL,
    invited_by UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'expired')),
    message TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    responded_at TIMESTAMP WITH TIME ZONE,
    
    UNIQUE(community_id, invited_user_id, status) -- Prevent duplicate pending invites
);

CREATE INDEX idx_invites_community ON member_invites(community_id);
CREATE INDEX idx_invites_user ON member_invites(invited_user_id, status);
CREATE INDEX idx_invites_expires ON member_invites(expires_at) WHERE status = 'pending';

-- ============================================
-- BANNED MEMBERS TABLE
-- ============================================
CREATE TABLE banned_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    banned_by UUID NOT NULL,
    reason TEXT NOT NULL,
    is_permanent BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(community_id, user_id)
);

CREATE INDEX idx_banned_community ON banned_members(community_id);
CREATE INDEX idx_banned_user ON banned_members(user_id);
CREATE INDEX idx_banned_expires ON banned_members(expires_at) WHERE NOT is_permanent;

-- ============================================
-- COMMUNITY RULES TABLE
-- ============================================
CREATE TABLE community_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    position INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_rules_community ON community_rules(community_id, position);
CREATE INDEX idx_rules_active ON community_rules(community_id, is_active);

-- ============================================
-- MODERATION ACTIONS TABLE
-- ============================================
CREATE TABLE moderation_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    moderator_id UUID NOT NULL,
    target_id UUID NOT NULL, -- User or content ID
    target_type VARCHAR(20) NOT NULL CHECK (target_type IN ('user', 'post', 'comment')),
    action VARCHAR(50) NOT NULL CHECK (action IN ('ban', 'unban', 'mute', 'unmute', 'remove_post', 'approve_post', 'remove_comment', 'warn')),
    reason TEXT NOT NULL,
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_mod_actions_community ON moderation_actions(community_id, created_at DESC);
CREATE INDEX idx_mod_actions_moderator ON moderation_actions(moderator_id);
CREATE INDEX idx_mod_actions_target ON moderation_actions(target_id, target_type);

-- ============================================
-- REPORTED CONTENT TABLE
-- ============================================
CREATE TABLE reported_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    content_id UUID NOT NULL,
    content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('post', 'comment')),
    reporter_id UUID NOT NULL,
    reason VARCHAR(50) NOT NULL CHECK (reason IN ('spam', 'harassment', 'inappropriate', 'violence', 'misinformation', 'other')),
    details TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'action_taken', 'dismissed')),
    reviewed_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reviewed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_reports_community ON reported_content(community_id, status);
CREATE INDEX idx_reports_content ON reported_content(content_id, content_type);
CREATE INDEX idx_reports_reporter ON reported_content(reporter_id);
CREATE INDEX idx_reports_status ON reported_content(status, created_at DESC);

-- ============================================
-- COMMUNITY ANALYTICS TABLE
-- ============================================
CREATE TABLE community_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    
    -- Growth metrics
    new_members INTEGER DEFAULT 0,
    left_members INTEGER DEFAULT 0,
    total_members INTEGER DEFAULT 0,
    
    -- Engagement metrics
    new_posts INTEGER DEFAULT 0,
    total_comments INTEGER DEFAULT 0,
    total_likes INTEGER DEFAULT 0,
    total_shares INTEGER DEFAULT 0,
    
    -- Activity metrics
    active_members INTEGER DEFAULT 0,
    engagement_rate DECIMAL(5,4) DEFAULT 0,
    avg_posts_per_member DECIMAL(10,2) DEFAULT 0,
    
    -- Moderation metrics
    reports_received INTEGER DEFAULT 0,
    actions_taken INTEGER DEFAULT 0,
    members_banned INTEGER DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(community_id, date)
);

CREATE INDEX idx_analytics_community ON community_analytics(community_id, date DESC);
CREATE INDEX idx_analytics_date ON community_analytics(date DESC);

-- ============================================
-- TRIGGERS
-- ============================================

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER communities_updated_at
    BEFORE UPDATE ON communities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER members_updated_at
    BEFORE UPDATE ON community_members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER rules_updated_at
    BEFORE UPDATE ON community_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Auto-increment member_count
CREATE OR REPLACE FUNCTION update_member_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.status = 'active' THEN
        UPDATE communities SET member_count = member_count + 1 WHERE id = NEW.community_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.status != 'active' AND NEW.status = 'active' THEN
        UPDATE communities SET member_count = member_count + 1 WHERE id = NEW.community_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.status = 'active' AND NEW.status != 'active' THEN
        UPDATE communities SET member_count = member_count - 1 WHERE id = NEW.community_id;
    ELSIF TG_OP = 'DELETE' AND OLD.status = 'active' THEN
        UPDATE communities SET member_count = member_count - 1 WHERE id = OLD.community_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER member_count_trigger
    AFTER INSERT OR UPDATE OR DELETE ON community_members
    FOR EACH ROW EXECUTE FUNCTION update_member_count();

-- Comments for documentation
COMMENT ON TABLE communities IS 'Main communities table - NOTE: Only cover_photo, NO profile_photo!';
COMMENT ON COLUMN communities.cover_photo IS 'Cover photo URL - communities only have cover photos, not profile photos';
COMMENT ON TABLE community_members IS 'Community membership with granular role-based permissions';
COMMENT ON COLUMN community_members.permissions IS 'Granular permissions JSON - allows custom per-user permissions';
COMMENT ON TABLE moderation_actions IS 'Audit log for all moderation actions';
COMMENT ON TABLE community_analytics IS 'Daily analytics snapshots for performance tracking';
