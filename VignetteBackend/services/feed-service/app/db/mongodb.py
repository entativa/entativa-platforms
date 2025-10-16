"""
MongoDB connection
"""
from motor.motor_asyncio import AsyncIOMotorClient
from pymongo.errors import ConnectionFailure
import logging

from app.config import get_settings


logger = logging.getLogger(__name__)
settings = get_settings()


class MongoDB:
    """MongoDB connection manager"""
    
    def __init__(self):
        self.client: AsyncIOMotorClient = None
        self.db = None
    
    async def connect(self):
        """Connect to MongoDB"""
        try:
            logger.info(f"Connecting to MongoDB: {settings.MONGODB_URL}")
            self.client = AsyncIOMotorClient(settings.MONGODB_URL)
            
            # Test connection
            await self.client.admin.command('ping')
            
            self.db = self.client[settings.MONGODB_DB]
            logger.info(f"✅ Connected to MongoDB database: {settings.MONGODB_DB}")
        
        except ConnectionFailure as e:
            logger.error(f"❌ MongoDB connection failed: {e}")
            raise
    
    async def disconnect(self):
        """Disconnect from MongoDB"""
        if self.client:
            self.client.close()
            logger.info("👋 Disconnected from MongoDB")
    
    def get_db(self):
        """Get database instance"""
        return self.db
    
    def get_collection(self, name: str):
        """Get collection"""
        return self.db[name]


# Global instance
mongodb = MongoDB()
