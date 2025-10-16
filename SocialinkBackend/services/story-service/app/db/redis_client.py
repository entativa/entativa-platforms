"""
Redis client for caching and real-time features
"""
import redis.asyncio as redis
from typing import Optional, Any
import json
import logging

from app.config import settings

logger = logging.getLogger(__name__)


class RedisClient:
    """Redis connection manager"""
    
    client: Optional[redis.Redis] = None
    
    @classmethod
    async def connect(cls):
        """Connect to Redis"""
        try:
            logger.info(f"Connecting to Redis at {settings.REDIS_URL}")
            cls.client = await redis.from_url(
                settings.REDIS_URL,
                db=settings.REDIS_DB,
                decode_responses=True
            )
            
            # Ping to verify connection
            await cls.client.ping()
            logger.info("Connected to Redis successfully")
            
        except Exception as e:
            logger.error(f"Failed to connect to Redis: {e}")
            raise
    
    @classmethod
    async def close(cls):
        """Close Redis connection"""
        if cls.client:
            await cls.client.close()
            logger.info("Redis connection closed")
    
    @classmethod
    async def get(cls, key: str) -> Optional[Any]:
        """Get value from Redis"""
        if not cls.client:
            return None
        
        try:
            value = await cls.client.get(key)
            if value:
                return json.loads(value)
            return None
        except Exception as e:
            logger.error(f"Redis GET error for key {key}: {e}")
            return None
    
    @classmethod
    async def set(cls, key: str, value: Any, ttl: Optional[int] = None):
        """Set value in Redis"""
        if not cls.client:
            return
        
        try:
            serialized = json.dumps(value, default=str)
            if ttl:
                await cls.client.setex(key, ttl, serialized)
            else:
                await cls.client.set(key, serialized)
        except Exception as e:
            logger.error(f"Redis SET error for key {key}: {e}")
    
    @classmethod
    async def delete(cls, key: str):
        """Delete key from Redis"""
        if not cls.client:
            return
        
        try:
            await cls.client.delete(key)
        except Exception as e:
            logger.error(f"Redis DELETE error for key {key}: {e}")
    
    @classmethod
    async def exists(cls, key: str) -> bool:
        """Check if key exists"""
        if not cls.client:
            return False
        
        try:
            return await cls.client.exists(key) > 0
        except Exception as e:
            logger.error(f"Redis EXISTS error for key {key}: {e}")
            return False
    
    @classmethod
    async def increment(cls, key: str, amount: int = 1) -> int:
        """Increment a counter"""
        if not cls.client:
            return 0
        
        try:
            return await cls.client.incrby(key, amount)
        except Exception as e:
            logger.error(f"Redis INCR error for key {key}: {e}")
            return 0
    
    @classmethod
    async def sadd(cls, key: str, *values):
        """Add to set"""
        if not cls.client:
            return
        
        try:
            await cls.client.sadd(key, *values)
        except Exception as e:
            logger.error(f"Redis SADD error for key {key}: {e}")
    
    @classmethod
    async def sismember(cls, key: str, value: str) -> bool:
        """Check if value in set"""
        if not cls.client:
            return False
        
        try:
            return await cls.client.sismember(key, value)
        except Exception as e:
            logger.error(f"Redis SISMEMBER error for key {key}: {e}")
            return False
    
    @classmethod
    async def smembers(cls, key: str) -> set:
        """Get all members of set"""
        if not cls.client:
            return set()
        
        try:
            return await cls.client.smembers(key)
        except Exception as e:
            logger.error(f"Redis SMEMBERS error for key {key}: {e}")
            return set()
    
    @classmethod
    async def zadd(cls, key: str, mapping: dict):
        """Add to sorted set"""
        if not cls.client:
            return
        
        try:
            await cls.client.zadd(key, mapping)
        except Exception as e:
            logger.error(f"Redis ZADD error for key {key}: {e}")
    
    @classmethod
    async def zrange(cls, key: str, start: int, end: int, desc: bool = False) -> list:
        """Get range from sorted set"""
        if not cls.client:
            return []
        
        try:
            if desc:
                return await cls.client.zrevrange(key, start, end)
            return await cls.client.zrange(key, start, end)
        except Exception as e:
            logger.error(f"Redis ZRANGE error for key {key}: {e}")
            return []
    
    @classmethod
    async def expire(cls, key: str, seconds: int):
        """Set expiry on key"""
        if not cls.client:
            return
        
        try:
            await cls.client.expire(key, seconds)
        except Exception as e:
            logger.error(f"Redis EXPIRE error for key {key}: {e}")


# Convenience function
async def get_redis() -> redis.Redis:
    """Get Redis client (for dependency injection)"""
    if not RedisClient.client:
        raise RuntimeError("Redis not initialized. Call connect() first.")
    return RedisClient.client
