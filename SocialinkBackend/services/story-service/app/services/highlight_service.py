"""
Highlight service for permanent story collections
"""
from datetime import datetime
from typing import List, Optional
from bson import ObjectId
import logging

from app.models.highlight import (
    Highlight, HighlightCreateRequest, HighlightUpdateRequest,
    HighlightResponse, HighlightStoriesRequest
)
from app.models.story import Story, StoryResponse
from app.db.mongodb import get_db
from app.db.redis_client import RedisClient

logger = logging.getLogger(__name__)


class HighlightService:
    """Service for managing story highlights"""
    
    def __init__(self):
        self.redis = RedisClient
    
    async def create_highlight(
        self,
        user_id: str,
        request: HighlightCreateRequest
    ) -> Highlight:
        """Create a new highlight"""
        db = await get_db()
        
        # Get next order index
        max_order = await db.highlights.find_one(
            {"user_id": user_id},
            sort=[("order_index", -1)]
        )
        order_index = (max_order["order_index"] + 1) if max_order else 0
        
        highlight = Highlight(
            user_id=user_id,
            title=request.title,
            description=request.description,
            story_ids=request.story_ids,
            story_count=len(request.story_ids),
            cover=request.cover or Highlight().cover,
            is_pinned=request.is_pinned,
            order_index=order_index
        )
        
        result = await db.highlights.insert_one(highlight.dict(by_alias=True, exclude={"id"}))
        highlight.id = result.inserted_id
        
        # Invalidate cache
        await self.redis.delete(f"user_highlights:{user_id}")
        
        logger.info(f"Created highlight {highlight.id} for user {user_id}")
        return highlight
    
    async def get_highlight(self, highlight_id: str) -> Optional[Highlight]:
        """Get a highlight by ID"""
        db = await get_db()
        
        highlight_doc = await db.highlights.find_one({
            "_id": ObjectId(highlight_id),
            "is_archived": False
        })
        
        if not highlight_doc:
            return None
        
        return Highlight(**highlight_doc)
    
    async def get_user_highlights(
        self,
        user_id: str,
        include_archived: bool = False
    ) -> List[Highlight]:
        """Get all highlights for a user"""
        db = await get_db()
        
        # Try cache
        cache_key = f"user_highlights:{user_id}"
        if not include_archived:
            cached = await self.redis.get(cache_key)
            if cached:
                return [Highlight(**h) for h in cached]
        
        query = {"user_id": user_id}
        if not include_archived:
            query["is_archived"] = False
        
        cursor = db.highlights.find(query).sort([
            ("is_pinned", -1),
            ("order_index", 1)
        ])
        
        highlights = []
        async for highlight_doc in cursor:
            highlights.append(Highlight(**highlight_doc))
        
        # Cache it
        if not include_archived:
            await self.redis.set(
                cache_key,
                [h.dict(by_alias=True) for h in highlights],
                ttl=3600
            )
        
        return highlights
    
    async def update_highlight(
        self,
        highlight_id: str,
        user_id: str,
        request: HighlightUpdateRequest
    ) -> bool:
        """Update a highlight"""
        db = await get_db()
        
        update_data = {k: v for k, v in request.dict().items() if v is not None}
        if not update_data:
            return False
        
        update_data["updated_at"] = datetime.utcnow()
        
        result = await db.highlights.update_one(
            {
                "_id": ObjectId(highlight_id),
                "user_id": user_id
            },
            {"$set": update_data}
        )
        
        if result.modified_count > 0:
            # Invalidate cache
            await self.redis.delete(f"user_highlights:{user_id}")
            logger.info(f"Updated highlight {highlight_id}")
            return True
        
        return False
    
    async def add_stories_to_highlight(
        self,
        highlight_id: str,
        user_id: str,
        story_ids: List[str]
    ) -> bool:
        """Add stories to a highlight"""
        db = await get_db()
        
        result = await db.highlights.update_one(
            {
                "_id": ObjectId(highlight_id),
                "user_id": user_id
            },
            {
                "$addToSet": {"story_ids": {"$each": story_ids}},
                "$inc": {"story_count": len(story_ids)},
                "$set": {"updated_at": datetime.utcnow()}
            }
        )
        
        if result.modified_count > 0:
            await self.redis.delete(f"user_highlights:{user_id}")
            logger.info(f"Added {len(story_ids)} stories to highlight {highlight_id}")
            return True
        
        return False
    
    async def remove_stories_from_highlight(
        self,
        highlight_id: str,
        user_id: str,
        story_ids: List[str]
    ) -> bool:
        """Remove stories from a highlight"""
        db = await get_db()
        
        result = await db.highlights.update_one(
            {
                "_id": ObjectId(highlight_id),
                "user_id": user_id
            },
            {
                "$pullAll": {"story_ids": story_ids},
                "$inc": {"story_count": -len(story_ids)},
                "$set": {"updated_at": datetime.utcnow()}
            }
        )
        
        if result.modified_count > 0:
            await self.redis.delete(f"user_highlights:{user_id}")
            logger.info(f"Removed {len(story_ids)} stories from highlight {highlight_id}")
            return True
        
        return False
    
    async def delete_highlight(self, highlight_id: str, user_id: str) -> bool:
        """Delete a highlight (hard delete)"""
        db = await get_db()
        
        result = await db.highlights.delete_one({
            "_id": ObjectId(highlight_id),
            "user_id": user_id
        })
        
        if result.deleted_count > 0:
            await self.redis.delete(f"user_highlights:{user_id}")
            logger.info(f"Deleted highlight {highlight_id}")
            return True
        
        return False
    
    async def reorder_highlights(
        self,
        user_id: str,
        highlight_ids: List[str]
    ) -> bool:
        """Reorder user's highlights"""
        db = await get_db()
        
        # Update order_index for each highlight
        for index, highlight_id in enumerate(highlight_ids):
            await db.highlights.update_one(
                {
                    "_id": ObjectId(highlight_id),
                    "user_id": user_id
                },
                {"$set": {"order_index": index}}
            )
        
        # Invalidate cache
        await self.redis.delete(f"user_highlights:{user_id}")
        logger.info(f"Reordered highlights for user {user_id}")
        return True
    
    async def get_highlight_with_stories(
        self,
        highlight_id: str
    ) -> Optional[HighlightResponse]:
        """Get highlight with its stories"""
        db = await get_db()
        
        # Get highlight
        highlight = await self.get_highlight(highlight_id)
        if not highlight:
            return None
        
        # Get stories
        if highlight.story_ids:
            story_docs = await db.stories.find({
                "_id": {"$in": [ObjectId(sid) for sid in highlight.story_ids]}
            }).to_list(length=100)
            
            stories = [Story(**doc) for doc in story_docs]
            story_responses = [
                StoryResponse(
                    id=str(s.id),
                    user_id=s.user_id,
                    story_type=s.story_type,
                    media_url=s.media_url,
                    thumbnail_url=s.thumbnail_url,
                    text_content=s.text_content,
                    background_color=s.background_color,
                    background_gradient=s.background_gradient,
                    duration=s.duration,
                    stickers=s.stickers,
                    privacy=s.privacy,
                    close_friends_only=s.close_friends_only,
                    view_count=s.view_count,
                    reply_count=s.reply_count,
                    has_viewed=False,
                    created_at=s.created_at,
                    expires_at=s.expires_at
                )
                for s in stories
            ]
        else:
            story_responses = []
        
        return HighlightResponse(
            id=str(highlight.id),
            user_id=highlight.user_id,
            title=highlight.title,
            description=highlight.description,
            cover=highlight.cover,
            story_count=highlight.story_count,
            is_pinned=highlight.is_pinned,
            created_at=highlight.created_at,
            updated_at=highlight.updated_at,
            stories=story_responses
        )
