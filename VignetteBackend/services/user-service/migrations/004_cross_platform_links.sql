-- Create cross_platform_links table for linking Entativa and Vignette accounts
CREATE TABLE IF NOT EXISTS cross_platform_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL, -- 'vignette' or 'entativa'
    platform_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, platform),
    INDEX idx_user_platform (user_id, platform),
    INDEX idx_platform_user (platform, platform_user_id)
);

-- Add comments
COMMENT ON TABLE cross_platform_links IS 'Links user accounts across Entativa and Vignette platforms';
COMMENT ON COLUMN cross_platform_links.platform IS 'The other platform name (vignette or entativa)';
COMMENT ON COLUMN cross_platform_links.platform_user_id IS 'The user ID on the other platform';

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_cross_platform_links_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_cross_platform_links_updated_at
    BEFORE UPDATE ON cross_platform_links
    FOR EACH ROW
    EXECUTE FUNCTION update_cross_platform_links_updated_at();
