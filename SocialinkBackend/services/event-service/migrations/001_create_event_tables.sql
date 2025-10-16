-- Event Service Database Schema
-- Facebook-style events with RSVP, invites, discussions

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis"; -- For location queries

-- ============================================
-- EVENTS TABLE
-- ============================================
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    cover_photo TEXT,
    
    -- Event details
    type VARCHAR(20) NOT NULL CHECK (type IN ('in_person', 'online')),
    category VARCHAR(30) NOT NULL CHECK (category IN ('social', 'business', 'entertainment', 'sports', 'education', 'religious', 'community', 'causes', 'health', 'arts', 'other')),
    privacy VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (privacy IN ('public', 'private', 'friends')),
    
    -- Location
    location_name VARCHAR(255),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location GEOGRAPHY(POINT, 4326), -- PostGIS for spatial queries
    online_link TEXT,
    
    -- Time
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    
    -- Recurring
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_rule TEXT, -- iCal RRULE format
    recurrence_end_date TIMESTAMP WITH TIME ZONE,
    
    -- Settings
    allow_guest_invites BOOLEAN DEFAULT TRUE,
    require_approval BOOLEAN DEFAULT FALSE,
    max_attendees INTEGER,
    
    -- Co-hosts
    co_hosts JSONB DEFAULT '[]',
    
    -- Stats (denormalized)
    going_count INTEGER DEFAULT 0,
    interested_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    
    -- Status
    is_cancelled BOOLEAN DEFAULT FALSE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_events_creator ON events(creator_id, created_at DESC);
CREATE INDEX idx_events_start_time ON events(start_time) WHERE NOT is_cancelled;
CREATE INDEX idx_events_category ON events(category, start_time) WHERE NOT is_cancelled;
CREATE INDEX idx_events_type ON events(type, start_time);
CREATE INDEX idx_events_privacy ON events(privacy);
CREATE INDEX idx_events_location ON events USING GIST(location) WHERE location IS NOT NULL;
CREATE INDEX idx_events_upcoming ON events(start_time) WHERE start_time > NOW() AND NOT is_cancelled;

-- Full-text search on title and description
CREATE INDEX idx_events_search ON events USING GIN(to_tsvector('english', title || ' ' || description));

-- ============================================
-- EVENT RSVPS TABLE
-- ============================================
CREATE TABLE event_rsvps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('going', 'interested', 'not_going')),
    guest_count INTEGER DEFAULT 0 CHECK (guest_count >= 0 AND guest_count <= 10),
    
    -- Check-in
    checked_in BOOLEAN DEFAULT FALSE,
    checked_in_at TIMESTAMP WITH TIME ZONE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_rsvps_event ON event_rsvps(event_id, status);
CREATE INDEX idx_rsvps_user ON event_rsvps(user_id, created_at DESC);
CREATE INDEX idx_rsvps_going ON event_rsvps(event_id) WHERE status = 'going';
CREATE INDEX idx_rsvps_interested ON event_rsvps(event_id) WHERE status = 'interested';

-- ============================================
-- EVENT INVITES TABLE
-- ============================================
CREATE TABLE event_invites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    inviter_id UUID NOT NULL,
    invitee_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(event_id, invitee_id)
);

CREATE INDEX idx_invites_event ON event_invites(event_id, status);
CREATE INDEX idx_invites_invitee ON event_invites(invitee_id, status);
CREATE INDEX idx_invites_pending ON event_invites(invitee_id) WHERE status = 'pending';

-- ============================================
-- EVENT DISCUSSIONS TABLE
-- ============================================
CREATE TABLE event_discussions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    content TEXT NOT NULL CHECK (length(content) <= 5000),
    media_urls JSONB DEFAULT '[]',
    is_pinned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_discussions_event ON event_discussions(event_id, created_at DESC);
CREATE INDEX idx_discussions_user ON event_discussions(user_id);
CREATE INDEX idx_discussions_pinned ON event_discussions(event_id, is_pinned) WHERE is_pinned = TRUE;

-- ============================================
-- EVENT REMINDERS TABLE
-- ============================================
CREATE TABLE event_reminders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    remind_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(event_id, user_id, remind_at)
);

CREATE INDEX idx_reminders_pending ON event_reminders(remind_at) WHERE NOT is_sent;
CREATE INDEX idx_reminders_user ON event_reminders(user_id);

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

CREATE TRIGGER events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER rsvps_updated_at
    BEFORE UPDATE ON event_rsvps
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Auto-update RSVP counts
CREATE OR REPLACE FUNCTION update_rsvp_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.status = 'going' THEN
            UPDATE events SET going_count = going_count + 1 WHERE id = NEW.event_id;
        ELSIF NEW.status = 'interested' THEN
            UPDATE events SET interested_count = interested_count + 1 WHERE id = NEW.event_id;
        END IF;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Decrement old status
        IF OLD.status = 'going' THEN
            UPDATE events SET going_count = GREATEST(0, going_count - 1) WHERE id = OLD.event_id;
        ELSIF OLD.status = 'interested' THEN
            UPDATE events SET interested_count = GREATEST(0, interested_count - 1) WHERE id = OLD.event_id;
        END IF;
        -- Increment new status
        IF NEW.status = 'going' THEN
            UPDATE events SET going_count = going_count + 1 WHERE id = NEW.event_id;
        ELSIF NEW.status = 'interested' THEN
            UPDATE events SET interested_count = interested_count + 1 WHERE id = NEW.event_id;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.status = 'going' THEN
            UPDATE events SET going_count = GREATEST(0, going_count - 1) WHERE id = OLD.event_id;
        ELSIF OLD.status = 'interested' THEN
            UPDATE events SET interested_count = GREATEST(0, interested_count - 1) WHERE id = OLD.event_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER rsvp_counts_trigger
    AFTER INSERT OR UPDATE OR DELETE ON event_rsvps
    FOR EACH ROW EXECUTE FUNCTION update_rsvp_counts();

-- Auto-create location point
CREATE OR REPLACE FUNCTION update_location_point()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.latitude IS NOT NULL AND NEW.longitude IS NOT NULL THEN
        NEW.location = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326)::geography;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER events_location_point
    BEFORE INSERT OR UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_location_point();

-- Comments
COMMENT ON TABLE events IS 'Facebook-style events with locations, RSVP, and recurring support';
COMMENT ON COLUMN events.recurrence_rule IS 'iCal RRULE format for recurring events';
COMMENT ON COLUMN events.location IS 'PostGIS geography for spatial queries';
COMMENT ON TABLE event_rsvps IS 'User responses to events (Going, Interested, Not Going)';
COMMENT ON TABLE event_invites IS 'Event invitations from users';
COMMENT ON TABLE event_discussions IS 'Event wall posts and discussions';
COMMENT ON TABLE event_reminders IS 'User-set reminders for events';
