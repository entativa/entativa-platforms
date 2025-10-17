"""
Configuration for Socialink Recommendation Service
"""
from pydantic_settings import BaseSettings
from functools import lru_cache
from typing import List


class Settings(BaseSettings):
    """Application settings"""
    
    # Service
    SERVICE_NAME: str = "socialink-recommendation-service"
    HOST: str = "0.0.0.0"
    PORT: int = 8095
    DEBUG: bool = False
    
    # Database
    MONGODB_URL: str = "mongodb://localhost:27017"
    MONGODB_DB: str = "socialink_recommendations"
    
    POSTGRES_URL: str = "postgresql://postgres:postgres@localhost:5432/socialink_recommendations"
    
    # Redis
    REDIS_URL: str = "redis://localhost:6379"
    REDIS_DB: int = 0
    
    # CORS
    CORS_ORIGINS: List[str] = ["*"]
    
    # ML Settings
    SIMILARITY_THRESHOLD: float = 0.1
    MIN_INTERACTIONS: int = 5
    COLLABORATIVE_WEIGHT: float = 0.6
    GRAPH_WEIGHT: float = 0.2
    POPULARITY_WEIGHT: float = 0.2
    
    # Cache TTL (seconds)
    USER_PROFILE_TTL: int = 3600
    RECOMMENDATIONS_TTL: int = 1800
    SIMILARITY_TTL: int = 7200
    
    class Config:
        env_file = ".env"
        case_sensitive = True


@lru_cache()
def get_settings() -> Settings:
    """Get cached settings instance"""
    return Settings()
