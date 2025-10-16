-- Create profiles table for Vignette (Instagram-like profile management)
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    category VARCHAR(50), -- personal, creator, business
    category_type VARCHAR(50), -- photographer, artist, musician, etc.
    gender VARCHAR(50),
    pronouns VARCHAR(50),
    link_in_bio JSONB DEFAULT '[]'::jsonb,
    highlights JSONB DEFAULT '[]'::jsonb,
    pinned_posts JSONB DEFAULT '[]'::jsonb,
    profile_badges JSONB DEFAULT '[]'::jsonb,
    contact_options JSONB,
    creator_insights JSONB,
    business_info JSONB,
    profile_views BIGINT DEFAULT 0,
    profile_views_enabled BOOLEAN DEFAULT TRUE,
    availability JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_category ON profiles(category);
CREATE INDEX idx_profiles_profile_views ON profiles(profile_views DESC);
CREATE INDEX idx_profiles_updated_at ON profiles(updated_at);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_profiles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER profiles_updated_at_trigger
    BEFORE UPDATE ON profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_profiles_updated_at();
