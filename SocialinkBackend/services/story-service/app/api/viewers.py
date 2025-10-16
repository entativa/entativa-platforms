"""
Viewers API endpoints (analytics)
"""
from fastapi import APIRouter, HTTPException, Header, Query
from typing import List
import logging

from app.models.viewer import StoryInsights
from app.models.story import StoryViewer
from app.services.viewer_service import ViewerService

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/viewers", tags=["viewers"])
viewer_service = ViewerService()


@router.get("/story/{story_id}", response_model=List[StoryViewer])
async def get_story_viewers(
    story_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """
    Get list of viewers for a story
    
    Only story owner can see viewers
    """
    try:
        viewers = await viewer_service.get_story_viewers(story_id, x_user_id)
        return viewers
    except Exception as e:
        logger.error(f"Error getting story viewers: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/story/{story_id}/insights", response_model=StoryInsights)
async def get_story_insights(
    story_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """
    Get detailed analytics for a story
    
    Includes:
    - View counts and reach
    - Interaction rates
    - Peak viewing times
    - Sticker performance
    - Top viewers
    """
    try:
        insights = await viewer_service.get_story_insights(story_id, x_user_id)
        if not insights:
            raise HTTPException(status_code=404, detail="Story not found or not owned by user")
        
        return insights
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error getting story insights: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/user/insights")
async def get_user_story_insights(
    x_user_id: str = Header(..., description="User ID from auth"),
    days: int = Query(7, ge=1, le=90)
):
    """
    Get aggregated insights for user's stories
    
    Shows performance over the last N days
    """
    try:
        insights = await viewer_service.get_user_story_insights(x_user_id, days)
        return insights
    except Exception as e:
        logger.error(f"Error getting user insights: {e}")
        raise HTTPException(status_code=500, detail=str(e))
