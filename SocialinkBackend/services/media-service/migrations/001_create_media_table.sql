-- Create media_type enum
CREATE TYPE media_type AS ENUM ('image', 'video', 'audio', 'document');

-- Create processing_status enum
CREATE TYPE processing_status AS ENUM ('pending', 'processing', 'completed', 'failed', 'cancelled');

-- Create media table
CREATE TABLE IF NOT EXISTS media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    media_type media_type NOT NULL,
    file_size BIGINT NOT NULL,
    storage_path TEXT NOT NULL,
    storage_provider VARCHAR(50) NOT NULL,
    url TEXT NOT NULL,
    cdn_url TEXT,
    thumbnail_url TEXT,
    width INTEGER,
    height INTEGER,
    duration DOUBLE PRECISION,
    aspect_ratio DOUBLE PRECISION,
    orientation SMALLINT,
    hash VARCHAR(128) NOT NULL,
    blurhash VARCHAR(100),
    processing_status processing_status DEFAULT 'pending',
    variants JSONB DEFAULT '[]'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    is_processed BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    view_count BIGINT DEFAULT 0,
    download_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT media_file_size_positive CHECK (file_size > 0),
    CONSTRAINT media_dimensions_positive CHECK (
        (width IS NULL OR width > 0) AND
        (height IS NULL OR height > 0)
    )
);

-- Create indexes for high performance
CREATE INDEX idx_media_user_id ON media(user_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_media_type ON media(media_type) WHERE is_deleted = FALSE;
CREATE INDEX idx_media_status ON media(processing_status) WHERE is_deleted = FALSE;
CREATE INDEX idx_media_created_desc ON media(created_at DESC) WHERE is_deleted = FALSE;
CREATE INDEX idx_media_hash ON media(hash) WHERE is_deleted = FALSE;
CREATE INDEX idx_media_is_public ON media(is_public) WHERE is_deleted = FALSE;

-- Create composite index for common queries
CREATE INDEX idx_media_user_type_created ON media(user_id, media_type, created_at DESC) 
    WHERE is_deleted = FALSE;

-- Create partial index for processing queue
CREATE INDEX idx_media_processing_queue ON media(created_at ASC)
    WHERE processing_status = 'pending' AND is_deleted = FALSE;

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_media_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER media_updated_at_trigger
    BEFORE UPDATE ON media
    FOR EACH ROW
    EXECUTE FUNCTION update_media_updated_at();

-- Create table for upload sessions (multipart uploads)
CREATE TABLE IF NOT EXISTS upload_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id UUID,
    user_id UUID NOT NULL,
    filename VARCHAR(255) NOT NULL,
    total_size BIGINT NOT NULL,
    chunk_size BIGINT NOT NULL,
    total_chunks INTEGER NOT NULL,
    uploaded_chunks JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    
    CONSTRAINT upload_total_size_positive CHECK (total_size > 0),
    CONSTRAINT upload_chunk_size_positive CHECK (chunk_size > 0)
);

CREATE INDEX idx_upload_sessions_user ON upload_sessions(user_id);
CREATE INDEX idx_upload_sessions_status ON upload_sessions(status);
CREATE INDEX idx_upload_sessions_expires ON upload_sessions(expires_at);

-- Create media analytics table
CREATE TABLE IF NOT EXISTS media_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- view, download, share
    user_id UUID,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_media_analytics_media ON media_analytics(media_id, created_at DESC);
CREATE INDEX idx_media_analytics_event ON media_analytics(event_type);
CREATE INDEX idx_media_analytics_created ON media_analytics(created_at DESC);

-- Create function to clean up expired upload sessions
CREATE OR REPLACE FUNCTION cleanup_expired_upload_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM upload_sessions WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE media IS 'Stores all media files uploaded by users';
COMMENT ON COLUMN media.hash IS 'BLAKE3 hash for deduplication and integrity';
COMMENT ON COLUMN media.blurhash IS 'BlurHash for progressive loading placeholders';
COMMENT ON COLUMN media.variants IS 'Different sizes/formats of the same media';
COMMENT ON COLUMN media.metadata IS 'EXIF, video codec info, audio metadata, etc.';
