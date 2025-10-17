"""
Story API endpoints
"""
from fastapi import APIRouter, HTTPException, Header, Query
from typing import Optional, List
import logging

from app.models.story import (
    StoryCreateRequest, StoryResponse, StoryBucket,
    StickerInteractionRequest
)
from app.services.story_service import StoryService

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/stories", tags=["stories"])
story_service = StoryService()


@router.post("/", response_model=StoryResponse, status_code=201)
async def create_story(
    request: StoryCreateRequest,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """
    Create a new story
    
    Stories expire after 24 hours.
    Supports:
    - Images, videos, text
    - Interactive stickers (polls, quizzes, questions, etc.)
    - Privacy controls
    - Close friends only mode
    """
    try:
        story = await story_service.create_story(x_user_id, request)
        return story_service._to_story_response(story)
    except Exception as e:
        logger.error(f"Error creating story: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/feed", response_model=List[StoryBucket])
async def get_story_feed(
    x_user_id: str = Header(..., description="User ID from auth"),
    limit: int = Query(50, ge=1, le=100)
):
    """
    Get story feed (story buckets from followed users)
    
    Returns story buckets sorted by:
    1. Close friends first
    2. Unseen stories
    3. Most recent
    """
    try:
        buckets = await story_service.get_story_feed(x_user_id, limit)
        return buckets
    except Exception as e:
        logger.error(f"Error getting story feed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/{story_id}", response_model=StoryResponse)
async def get_story(
    story_id: str,
    x_user_id: Optional[str] = Header(None, description="User ID from auth")
):
    """Get a single story by ID"""
    try:
        story = await story_service.get_story(story_id, x_user_id)
        if not story:
            raise HTTPException(status_code=404, detail="Story not found")
        
        return story_service._to_story_response(story)
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error getting story: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/user/{user_id}", response_model=List[StoryResponse])
async def get_user_stories(
    user_id: str,
    x_user_id: Optional[str] = Header(None, description="User ID from auth"),
    include_expired: bool = Query(False)
):
    """Get all stories by a specific user"""
    try:
        stories = await story_service.get_user_stories(
            user_id,
            x_user_id,
            include_expired
        )
        return [story_service._to_story_response(s) for s in stories]
    except Exception as e:
        logger.error(f"Error getting user stories: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/{story_id}/view")
async def mark_story_viewed(
    story_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Mark a story as viewed"""
    try:
        result = await story_service.mark_as_viewed(story_id, x_user_id)
        return {"success": result, "message": "Story marked as viewed" if result else "Already viewed"}
    except Exception as e:
        logger.error(f"Error marking story as viewed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/{story_id}/reply")
async def reply_to_story(
    story_id: str,
    message: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """
    Reply to a story
    
    Replies are private messages to the story owner
    """
    try:
        if not message or len(message) > 500:
            raise HTTPException(status_code=400, detail="Invalid message length")
        
        result = await story_service.add_reply(story_id, x_user_id, message)
        if not result:
            raise HTTPException(status_code=404, detail="Story not found")
        
        return {"success": True, "message": "Reply sent"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error replying to story: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/{story_id}/interact")
async def interact_with_sticker(
    story_id: str,
    interaction: StickerInteractionRequest,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """
    Interact with story stickers
    
    Supports:
    - Poll voting
    - Quiz answers
    - Question responses
    - Slider values
    - Countdown follows
    """
    try:
        result = await story_service.interact_with_sticker(
            story_id,
            x_user_id,
            interaction
        )
        
        if not result:
            raise HTTPException(
                status_code=400,
                detail="Interaction failed (already interacted or invalid)"
            )
        
        return {"success": True, "message": "Interaction recorded"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error interacting with sticker: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{story_id}")
async def delete_story(
    story_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Delete a story (soft delete)"""
    try:
        result = await story_service.delete_story(story_id, x_user_id)
        if not result:
            raise HTTPException(
                status_code=404,
                detail="Story not found or not owned by user"
            )
        
        return {"success": True, "message": "Story deleted"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error deleting story: {e}")
        raise HTTPException(status_code=500, detail=str(e))
