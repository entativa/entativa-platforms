"""
Highlights API endpoints
"""
from fastapi import APIRouter, HTTPException, Header, Query
from typing import List
import logging

from app.models.highlight import (
    HighlightCreateRequest, HighlightUpdateRequest,
    HighlightResponse, HighlightStoriesRequest
)
from app.services.highlight_service import HighlightService

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/highlights", tags=["highlights"])
highlight_service = HighlightService()


@router.post("/", response_model=HighlightResponse, status_code=201)
async def create_highlight(
    request: HighlightCreateRequest,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """
    Create a new highlight
    
    Highlights are permanent collections of stories.
    Perfect for:
    - Travel memories
    - Life events
    - Product showcases
    - Tutorials
    """
    try:
        highlight = await highlight_service.create_highlight(x_user_id, request)
        
        return HighlightResponse(
            id=str(highlight.id),
            user_id=highlight.user_id,
            title=highlight.title,
            description=highlight.description,
            cover=highlight.cover,
            story_count=highlight.story_count,
            is_pinned=highlight.is_pinned,
            created_at=highlight.created_at,
            updated_at=highlight.updated_at
        )
    except Exception as e:
        logger.error(f"Error creating highlight: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/user/{user_id}", response_model=List[HighlightResponse])
async def get_user_highlights(
    user_id: str,
    include_archived: bool = Query(False)
):
    """Get all highlights for a user"""
    try:
        highlights = await highlight_service.get_user_highlights(
            user_id,
            include_archived
        )
        
        return [
            HighlightResponse(
                id=str(h.id),
                user_id=h.user_id,
                title=h.title,
                description=h.description,
                cover=h.cover,
                story_count=h.story_count,
                is_pinned=h.is_pinned,
                created_at=h.created_at,
                updated_at=h.updated_at
            )
            for h in highlights
        ]
    except Exception as e:
        logger.error(f"Error getting user highlights: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/{highlight_id}", response_model=HighlightResponse)
async def get_highlight(highlight_id: str):
    """Get a single highlight with its stories"""
    try:
        highlight = await highlight_service.get_highlight_with_stories(highlight_id)
        if not highlight:
            raise HTTPException(status_code=404, detail="Highlight not found")
        
        return highlight
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error getting highlight: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.patch("/{highlight_id}")
async def update_highlight(
    highlight_id: str,
    request: HighlightUpdateRequest,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Update a highlight"""
    try:
        result = await highlight_service.update_highlight(
            highlight_id,
            x_user_id,
            request
        )
        
        if not result:
            raise HTTPException(
                status_code=404,
                detail="Highlight not found or not owned by user"
            )
        
        return {"success": True, "message": "Highlight updated"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error updating highlight: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/{highlight_id}/stories")
async def add_stories_to_highlight(
    highlight_id: str,
    request: HighlightStoriesRequest,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Add stories to a highlight"""
    try:
        result = await highlight_service.add_stories_to_highlight(
            highlight_id,
            x_user_id,
            request.story_ids
        )
        
        if not result:
            raise HTTPException(
                status_code=404,
                detail="Highlight not found or not owned by user"
            )
        
        return {"success": True, "message": f"Added {len(request.story_ids)} stories"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error adding stories to highlight: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{highlight_id}/stories")
async def remove_stories_from_highlight(
    highlight_id: str,
    request: HighlightStoriesRequest,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Remove stories from a highlight"""
    try:
        result = await highlight_service.remove_stories_from_highlight(
            highlight_id,
            x_user_id,
            request.story_ids
        )
        
        if not result:
            raise HTTPException(
                status_code=404,
                detail="Highlight not found or not owned by user"
            )
        
        return {"success": True, "message": f"Removed {len(request.story_ids)} stories"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error removing stories from highlight: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.delete("/{highlight_id}")
async def delete_highlight(
    highlight_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Delete a highlight"""
    try:
        result = await highlight_service.delete_highlight(highlight_id, x_user_id)
        if not result:
            raise HTTPException(
                status_code=404,
                detail="Highlight not found or not owned by user"
            )
        
        return {"success": True, "message": "Highlight deleted"}
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error deleting highlight: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/reorder")
async def reorder_highlights(
    highlight_ids: List[str],
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Reorder user's highlights"""
    try:
        result = await highlight_service.reorder_highlights(x_user_id, highlight_ids)
        return {"success": True, "message": "Highlights reordered"}
    except Exception as e:
        logger.error(f"Error reordering highlights: {e}")
        raise HTTPException(status_code=500, detail=str(e))
