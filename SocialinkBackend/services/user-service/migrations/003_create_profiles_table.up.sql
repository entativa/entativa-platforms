-- Create profiles table for Socialink (Facebook-like profile management)
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    hometown VARCHAR(100),
    current_city VARCHAR(100),
    relationship_status VARCHAR(50),
    languages JSONB DEFAULT '[]'::jsonb,
    interested_in JSONB DEFAULT '[]'::jsonb,
    work JSONB DEFAULT '[]'::jsonb,
    education JSONB DEFAULT '[]'::jsonb,
    contact_info JSONB,
    about TEXT,
    favorite_quotes TEXT,
    hobbies JSONB DEFAULT '[]'::jsonb,
    website VARCHAR(255),
    social_links JSONB,
    featured_photos JSONB DEFAULT '[]'::jsonb,
    visibility JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_hometown ON profiles(hometown);
CREATE INDEX idx_profiles_current_city ON profiles(current_city);
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
