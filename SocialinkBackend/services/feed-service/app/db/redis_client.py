"""
Redis connection
"""
import redis.asyncio as redis
from redis.exceptions import ConnectionError
import logging

from app.config import get_settings


logger = logging.getLogger(__name__)
settings = get_settings()


class RedisClient:
    """Redis connection manager"""
    
    def __init__(self):
        self.client: redis.Redis = None
    
    async def connect(self):
        """Connect to Redis"""
        try:
            logger.info(f"Connecting to Redis: {settings.REDIS_URL}")
            self.client = await redis.from_url(
                settings.REDIS_URL,
                db=settings.REDIS_DB,
                encoding="utf-8",
                decode_responses=True
            )
            
            # Test connection
            await self.client.ping()
            logger.info(f"✅ Connected to Redis")
        
        except ConnectionError as e:
            logger.error(f"❌ Redis connection failed: {e}")
            raise
    
    async def disconnect(self):
        """Disconnect from Redis"""
        if self.client:
            await self.client.close()
            logger.info("👋 Disconnected from Redis")
    
    def get_client(self):
        """Get Redis client"""
        return self.client


# Global instance
redis_client = RedisClient()
