"""
Close friends service for managing close friends lists
"""
from datetime import datetime
from typing import List
from bson import ObjectId
import logging

from app.models.viewer import CloseFriend
from app.db.mongodb import get_db
from app.db.redis_client import RedisClient

logger = logging.getLogger(__name__)


class CloseFriendsService:
    """Service for managing close friends"""
    
    def __init__(self):
        self.redis = RedisClient
    
    async def add_close_friend(
        self,
        user_id: str,
        friend_user_id: str
    ) -> bool:
        """Add a user to close friends list"""
        db = await get_db()
        
        # Create close friend document
        close_friend = {
            "user_id": user_id,
            "friend_user_id": friend_user_id,
            "added_at": datetime.utcnow(),
            "is_mutual": False  # Would check if reciprocal
        }
        
        try:
            await db.close_friends.insert_one(close_friend)
            
            # Add to Redis set for fast lookup
            await self.redis.sadd(f"close_friends:{user_id}", friend_user_id)
            
            logger.info(f"User {user_id} added {friend_user_id} to close friends")
            return True
        
        except Exception as e:
            # Duplicate key error
            logger.warning(f"Close friend relationship already exists: {e}")
            return False
    
    async def remove_close_friend(
        self,
        user_id: str,
        friend_user_id: str
    ) -> bool:
        """Remove a user from close friends list"""
        db = await get_db()
        
        result = await db.close_friends.delete_one({
            "user_id": user_id,
            "friend_user_id": friend_user_id
        })
        
        if result.deleted_count > 0:
            # Remove from Redis
            await self.redis.delete(f"close_friends:{user_id}")
            
            logger.info(f"User {user_id} removed {friend_user_id} from close friends")
            return True
        
        return False
    
    async def get_close_friends(self, user_id: str) -> List[str]:
        """Get user's close friends list"""
        # Try Redis cache first
        cached = await self.redis.smembers(f"close_friends:{user_id}")
        if cached:
            return list(cached)
        
        # Get from database
        db = await get_db()
        
        cursor = db.close_friends.find({"user_id": user_id})
        friend_ids = []
        
        async for doc in cursor:
            friend_ids.append(doc["friend_user_id"])
        
        # Cache it
        if friend_ids:
            await self.redis.sadd(f"close_friends:{user_id}", *friend_ids)
            await self.redis.expire(f"close_friends:{user_id}", 3600)  # 1 hour
        
        return friend_ids
    
    async def is_close_friend(
        self,
        user_id: str,
        friend_user_id: str
    ) -> bool:
        """Check if a user is in close friends list"""
        return await self.redis.sismember(
            f"close_friends:{user_id}",
            friend_user_id
        )
    
    async def get_close_friends_count(self, user_id: str) -> int:
        """Get count of close friends"""
        db = await get_db()
        
        count = await db.close_friends.count_documents({"user_id": user_id})
        return count
