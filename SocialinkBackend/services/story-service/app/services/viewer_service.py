"""
Viewer service for story analytics and insights
"""
from datetime import datetime, timedelta
from typing import List, Dict
from bson import ObjectId
import logging

from app.models.viewer import ViewerInsight, StoryInsights
from app.models.story import Story, StoryViewer
from app.db.mongodb import get_db
from app.db.redis_client import RedisClient

logger = logging.getLogger(__name__)


class ViewerService:
    """Service for story viewer analytics"""
    
    def __init__(self):
        self.redis = RedisClient
    
    async def get_story_viewers(
        self,
        story_id: str,
        user_id: str
    ) -> List[StoryViewer]:
        """Get list of viewers for a story"""
        db = await get_db()
        
        story = await db.stories.find_one({
            "_id": ObjectId(story_id),
            "user_id": user_id
        })
        
        if not story:
            return []
        
        return [StoryViewer(**v) for v in story.get("viewers", [])]
    
    async def get_story_insights(
        self,
        story_id: str,
        user_id: str
    ) -> Optional[StoryInsights]:
        """Get detailed insights for a story"""
        db = await get_db()
        
        story = await db.stories.find_one({
            "_id": ObjectId(story_id),
            "user_id": user_id
        })
        
        if not story:
            return None
        
        story_obj = Story(**story)
        
        # Calculate metrics
        unique_viewers = len(set(v.user_id for v in story_obj.viewers))
        total_interactions = sum(
            self._count_sticker_interactions(sticker)
            for sticker in story_obj.stickers
        )
        
        interaction_rate = (total_interactions / unique_viewers * 100) if unique_viewers > 0 else 0
        
        # Views by hour
        views_by_hour: Dict[int, int] = {}
        for viewer in story_obj.viewers:
            hour = viewer.viewed_at.hour
            views_by_hour[hour] = views_by_hour.get(hour, 0) + 1
        
        peak_hour = max(views_by_hour, key=views_by_hour.get) if views_by_hour else None
        
        # Sticker performance
        sticker_interactions = {}
        for i, sticker in enumerate(story_obj.stickers):
            count = self._count_sticker_interactions(sticker)
            if count > 0:
                sticker_interactions[f"sticker_{i}_{sticker.sticker_type}"] = count
        
        # Top viewers (would need more data from user service)
        top_viewers = []
        
        insights = StoryInsights(
            story_id=story_id,
            total_views=story_obj.view_count,
            unique_viewers=unique_viewers,
            reach_percentage=0.0,  # Would calculate from follower count
            total_interactions=total_interactions,
            interaction_rate=round(interaction_rate, 2),
            replies=story_obj.reply_count,
            top_viewers=top_viewers,
            views_by_hour=views_by_hour,
            peak_viewing_hour=peak_hour,
            sticker_interactions=sticker_interactions
        )
        
        return insights
    
    async def get_user_story_insights(
        self,
        user_id: str,
        days: int = 7
    ) -> Dict:
        """Get aggregated insights for user's stories"""
        db = await get_db()
        
        # Get stories from last N days
        since = datetime.utcnow() - timedelta(days=days)
        
        pipeline = [
            {
                "$match": {
                    "user_id": user_id,
                    "created_at": {"$gte": since},
                    "is_deleted": False
                }
            },
            {
                "$group": {
                    "_id": None,
                    "total_stories": {"$sum": 1},
                    "total_views": {"$sum": "$view_count"},
                    "total_replies": {"$sum": "$reply_count"},
                    "avg_views_per_story": {"$avg": "$view_count"}
                }
            }
        ]
        
        result = await db.stories.aggregate(pipeline).to_list(length=1)
        
        if not result:
            return {
                "total_stories": 0,
                "total_views": 0,
                "total_replies": 0,
                "avg_views_per_story": 0
            }
        
        return result[0]
    
    def _count_sticker_interactions(self, sticker) -> int:
        """Count total interactions on a sticker"""
        count = 0
        
        if sticker.poll:
            count += sticker.poll.total_votes
        
        if sticker.quiz:
            count += sticker.quiz.total_attempts
        
        if sticker.question:
            count += sticker.question.response_count
        
        if sticker.slider and sticker.slider.responses:
            count += len(sticker.slider.responses)
        
        if sticker.countdown:
            count += sticker.countdown.follower_count
        
        return count
