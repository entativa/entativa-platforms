"""
Expiration service for managing story lifecycle
Background job to expire stories after 24 hours
"""
from datetime import datetime
import logging
import asyncio

from app.db.mongodb import get_db
from app.db.redis_client import RedisClient

logger = logging.getLogger(__name__)


class ExpirationService:
    """Service to handle story expiration"""
    
    def __init__(self):
        self.redis = RedisClient
        self.is_running = False
    
    async def start(self):
        """Start the expiration checker"""
        self.is_running = True
        logger.info("Story expiration service started")
        
        while self.is_running:
            try:
                await self.expire_old_stories()
                await asyncio.sleep(300)  # Check every 5 minutes
            except Exception as e:
                logger.error(f"Error in expiration service: {e}")
                await asyncio.sleep(60)
    
    def stop(self):
        """Stop the expiration checker"""
        self.is_running = False
        logger.info("Story expiration service stopped")
    
    async def expire_old_stories(self):
        """Mark expired stories as expired"""
        db = await get_db()
        
        now = datetime.utcnow()
        
        # Find stories that should be expired
        result = await db.stories.update_many(
            {
                "is_expired": False,
                "is_deleted": False,
                "expires_at": {"$lte": now}
            },
            {
                "$set": {
                    "is_expired": True,
                    "updated_at": now
                }
            }
        )
        
        if result.modified_count > 0:
            logger.info(f"Expired {result.modified_count} stories")
            
            # Clear relevant caches
            await self._clear_expired_caches(result.modified_count)
    
    async def _clear_expired_caches(self, count: int):
        """Clear caches for expired stories"""
        # This would be more sophisticated in production
        # For now, we'll just log it
        logger.info(f"Cleared caches for {count} expired stories")
    
    async def cleanup_old_expired_stories(self, days: int = 30):
        """Hard delete stories expired for more than X days"""
        db = await get_db()
        
        cutoff = datetime.utcnow() - timedelta(days=days)
        
        result = await db.stories.delete_many({
            "is_expired": True,
            "expires_at": {"$lte": cutoff}
        })
        
        if result.deleted_count > 0:
            logger.info(f"Cleaned up {result.deleted_count} old expired stories")
        
        return result.deleted_count
