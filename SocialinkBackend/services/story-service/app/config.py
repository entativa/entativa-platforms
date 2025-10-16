"""
Configuration settings for Socialink Story Service
"""
from pydantic_settings import BaseSettings
from typing import Optional


class Settings(BaseSettings):
    """Application settings"""
    
    # App
    APP_NAME: str = "Socialink Story Service"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = False
    
    # Server
    HOST: str = "0.0.0.0"
    PORT: int = 8085
    
    # MongoDB
    MONGODB_URL: str = "mongodb://localhost:27017"
    MONGODB_DB_NAME: str = "socialink_stories"
    
    # Redis
    REDIS_URL: str = "redis://localhost:6379"
    REDIS_DB: int = 0
    REDIS_TTL: int = 3600  # 1 hour cache
    
    # Story settings
    STORY_EXPIRY_HOURS: int = 24
    MAX_STORIES_PER_USER: int = 100
    MAX_STORY_DURATION: int = 15  # seconds
    MAX_HIGHLIGHT_STORIES: int = 100
    
    # Media service (gRPC)
    MEDIA_SERVICE_HOST: str = "localhost"
    MEDIA_SERVICE_PORT: int = 50051
    
    # Authentication
    JWT_SECRET: str = "your-secret-key-change-in-production"
    JWT_ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 30
    
    # Rate limiting
    RATE_LIMIT_STORIES_PER_DAY: int = 50
    RATE_LIMIT_VIEWS_PER_MINUTE: int = 100
    
    # Feature flags
    ENABLE_CLOSE_FRIENDS: bool = True
    ENABLE_STORY_REPLIES: bool = True
    ENABLE_INTERACTIVE_STICKERS: bool = True
    ENABLE_STORY_INSIGHTS: bool = True
    
    class Config:
        env_file = ".env"
        case_sensitive = True


settings = Settings()
