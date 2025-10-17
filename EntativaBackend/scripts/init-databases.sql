-- Initialize all databases for Vignette microservices

-- User Service
CREATE DATABASE vignette_users;

-- Post Service  
CREATE DATABASE vignette_posts;

-- Messaging Service
CREATE DATABASE vignette_messaging;

-- Settings Service
CREATE DATABASE vignette_settings;

-- Story Service uses MongoDB (not needed here)

-- Search Service uses Elasticsearch (not needed here)

-- Notification Service
CREATE DATABASE vignette_notifications;

-- Feed Service uses MongoDB (not needed here)

-- Community Service
CREATE DATABASE vignette_communities;

-- Recommendation Service uses MongoDB (not needed here)

-- Streaming Service
CREATE DATABASE vignette_streaming;

-- Creator Service
CREATE DATABASE vignette_creator;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE vignette_users TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_posts TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_messaging TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_settings TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_notifications TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_communities TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_streaming TO postgres;
GRANT ALL PRIVILEGES ON DATABASE vignette_creator TO postgres;
