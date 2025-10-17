"""
MongoDB client and utilities
"""
from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorDatabase
from typing import Optional
import logging

from app.config import settings

logger = logging.getLogger(__name__)


class MongoDB:
    """MongoDB connection manager"""
    
    client: Optional[AsyncIOMotorClient] = None
    db: Optional[AsyncIOMotorDatabase] = None
    
    @classmethod
    async def connect(cls):
        """Connect to MongoDB"""
        try:
            logger.info(f"Connecting to MongoDB at {settings.MONGODB_URL}")
            cls.client = AsyncIOMotorClient(settings.MONGODB_URL)
            cls.db = cls.client[settings.MONGODB_DB_NAME]
            
            # Ping to verify connection
            await cls.client.admin.command('ping')
            logger.info(f"Connected to MongoDB database: {settings.MONGODB_DB_NAME}")
            
            # Create indexes
            await cls.create_indexes()
            
        except Exception as e:
            logger.error(f"Failed to connect to MongoDB: {e}")
            raise
    
    @classmethod
    async def close(cls):
        """Close MongoDB connection"""
        if cls.client:
            cls.client.close()
            logger.info("MongoDB connection closed")
    
    @classmethod
    async def create_indexes(cls):
        """Create database indexes for performance"""
        if not cls.db:
            return
        
        try:
            # Stories collection indexes
            stories = cls.db.stories
            await stories.create_index("user_id")
            await stories.create_index("created_at")
            await stories.create_index("expires_at")
            await stories.create_index([("user_id", 1), ("created_at", -1)])
            await stories.create_index([("is_expired", 1), ("expires_at", 1)])
            await stories.create_index([("privacy", 1), ("is_expired", 1)])
            
            # Compound index for story feed queries
            await stories.create_index([
                ("is_expired", 1),
                ("is_deleted", 1),
                ("created_at", -1)
            ])
            
            # Highlights collection indexes
            highlights = cls.db.highlights
            await highlights.create_index("user_id")
            await highlights.create_index([("user_id", 1), ("order_index", 1)])
            await highlights.create_index([("user_id", 1), ("is_pinned", -1)])
            
            # Viewers collection indexes (for analytics)
            viewers = cls.db.story_viewers
            await viewers.create_index("story_id")
            await viewers.create_index("user_id")
            await viewers.create_index([("story_id", 1), ("user_id", 1)], unique=True)
            
            # Close friends collection indexes
            close_friends = cls.db.close_friends
            await close_friends.create_index("user_id")
            await close_friends.create_index([("user_id", 1), ("friend_user_id", 1)], unique=True)
            
            logger.info("MongoDB indexes created successfully")
            
        except Exception as e:
            logger.error(f"Failed to create indexes: {e}")
    
    @classmethod
    def get_database(cls) -> AsyncIOMotorDatabase:
        """Get database instance"""
        if not cls.db:
            raise RuntimeError("Database not initialized. Call connect() first.")
        return cls.db


# Convenience function
async def get_db() -> AsyncIOMotorDatabase:
    """Get database instance (for dependency injection)"""
    return MongoDB.get_database()
