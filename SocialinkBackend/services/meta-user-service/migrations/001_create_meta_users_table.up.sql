-- Create meta_users table
CREATE TABLE IF NOT EXISTS meta_users (
    id UUID PRIMARY KEY,
    meta_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    phone_number VARCHAR(50),
    phone_verified BOOLEAN DEFAULT FALSE,
    password_hash TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    trust_score DOUBLE PRECISION DEFAULT 0.5,
    risk_level VARCHAR(50) DEFAULT 'medium',
    account_tier VARCHAR(50) DEFAULT 'basic',
    platform_links JSONB,
    security_profile JSONB,
    privacy_settings JSONB,
    data_rights JSONB,
    device_fingerprints JSONB DEFAULT '[]'::jsonb,
    biometric_tokens JSONB DEFAULT '[]'::jsonb,
    federated_identities JSONB DEFAULT '[]'::jsonb,
    session_management JSONB,
    anomaly_detection JSONB,
    cross_platform_activity JSONB,
    compliance_data JSONB,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX idx_meta_users_meta_id ON meta_users(meta_id);
CREATE INDEX idx_meta_users_email ON meta_users(email);
CREATE INDEX idx_meta_users_status ON meta_users(status);
CREATE INDEX idx_meta_users_trust_score ON meta_users(trust_score);
CREATE INDEX idx_meta_users_risk_level ON meta_users(risk_level);
CREATE INDEX idx_meta_users_deleted_at ON meta_users(deleted_at);
CREATE INDEX idx_meta_users_created_at ON meta_users(created_at);

-- Create event sourcing table
CREATE TABLE IF NOT EXISTS meta_user_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES meta_users(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_meta_user_events_user_id ON meta_user_events(user_id);
CREATE INDEX idx_meta_user_events_event_type ON meta_user_events(event_type);
CREATE INDEX idx_meta_user_events_created_at ON meta_user_events(created_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for auto-updating updated_at
CREATE TRIGGER update_meta_users_updated_at BEFORE UPDATE ON meta_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
