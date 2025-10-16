"""
Configuration for Vignette Feed Service
"""
from pydantic_settings import BaseSettings
from functools import lru_cache
from typing import List


class Settings(BaseSettings):
    """Application settings"""
    
    # Service
    SERVICE_NAME: str = "vignette-feed-service"
    HOST: str = "0.0.0.0"
    PORT: int = 8085
    DEBUG: bool = False
    
    # Database
    MONGODB_URL: str = "mongodb://localhost:27017"
    MONGODB_DB: str = "vignette_feed"
    
    POSTGRES_URL: str = "postgresql://postgres:postgres@localhost:5432/vignette_feed"
    
    # Redis
    REDIS_URL: str = "redis://localhost:6379"
    REDIS_DB: int = 0
    
    # Elasticsearch (optional)
    ELASTICSEARCH_URL: str = "http://localhost:9200"
    ELASTICSEARCH_INDEX: str = "vignette_content"
    
    # CORS
    CORS_ORIGINS: List[str] = ["*"]
    
    # JWT (for auth)
    JWT_SECRET: str = "your-secret-key-change-in-production"
    JWT_ALGORITHM: str = "HS256"
    
    # Feed settings
    DEFAULT_FEED_SIZE: int = 20
    MAX_FEED_SIZE: int = 100
    CANDIDATE_POOL_SIZE: int = 500
    
    # Cache TTL (seconds)
    USER_PROFILE_TTL: int = 3600  # 1 hour
    VIRAL_CONTENT_TTL: int = 300  # 5 minutes
    CONTENT_DATA_TTL: int = 3600  # 1 hour
    
    # Ranking weights (can be adjusted)
    HOME_VIRAL_WEIGHT: float = 0.35
    HOME_INTEREST_WEIGHT: float = 0.30
    HOME_QUALITY_WEIGHT: float = 0.15
    HOME_RECENCY_WEIGHT: float = 0.15
    HOME_DIVERSITY_WEIGHT: float = 0.05
    
    CIRCLE_SOCIAL_WEIGHT: float = 0.50
    CIRCLE_PROXIMITY_WEIGHT: float = 0.25
    CIRCLE_ENGAGEMENT_WEIGHT: float = 0.15
    CIRCLE_RECENCY_WEIGHT: float = 0.10
    
    SURPRISE_KNOWN_RATIO: float = 0.60
    SURPRISE_EXPLORE_RATIO: float = 0.30
    SURPRISE_RANDOM_RATIO: float = 0.10
    
    class Config:
        env_file = ".env"
        case_sensitive = True


@lru_cache()
def get_settings() -> Settings:
    """Get cached settings instance"""
    return Settings()
