"""
Feed API endpoints for Vignette
"""
from fastapi import APIRouter, Depends, HTTPException, Query
from typing import Optional, List

from app.schemas.feed import (
    FeedType, FeedRequest, FeedResponse,
    UserSignal, FeedMetrics
)
from app.services.feed_service import FeedService
from app.dependencies import get_feed_service, get_current_user


router = APIRouter(prefix="/feed", tags=["feed"])


@router.get("/{feed_type}", response_model=FeedResponse)
async def get_feed(
    feed_type: FeedType,
    limit: int = Query(default=20, ge=1, le=100),
    offset: int = Query(default=0, ge=0),
    refresh: bool = Query(default=False),
    seen_ids: Optional[str] = Query(default=None),
    latitude: Optional[float] = Query(default=None),
    longitude: Optional[float] = Query(default=None),
    current_user: dict = Depends(get_current_user),
    feed_service: FeedService = Depends(get_feed_service)
):
    """
    Get personalized feed
    
    Feed types:
    - **home**: For You feed (TikTok-style discovery)
    - **circle**: Friends + nearby content
    - **surprise_delight**: Balanced exploration (60% interests, 30% explore, 10% surprise)
    """
    user_id = current_user["user_id"]
    
    # Parse seen_ids
    seen_list = seen_ids.split(",") if seen_ids else []
    
    request = FeedRequest(
        feed_type=feed_type,
        user_id=user_id,
        limit=limit,
        offset=offset,
        seen_ids=seen_list,
        refresh=refresh,
        latitude=latitude,
        longitude=longitude
    )
    
    response = await feed_service.generate_feed(request)
    
    return response


@router.post("/signal")
async def track_signal(
    signal: UserSignal,
    current_user: dict = Depends(get_current_user),
    feed_service: FeedService = Depends(get_feed_service)
):
    """
    Track user engagement signal
    
    Signal types:
    - **view**: Content viewed
    - **like**: Content liked
    - **comment**: User commented
    - **share**: Content shared
    - **save**: Content saved
    - **skip**: Content skipped
    - **hide**: Content hidden
    - **report**: Content reported
    """
    # Verify user owns signal
    if signal.user_id != current_user["user_id"]:
        raise HTTPException(status_code=403, detail="Cannot track signal for another user")
    
    await feed_service.track_signal(signal)
    
    return {"status": "success", "message": "Signal tracked"}


@router.get("/metrics/{feed_type}", response_model=FeedMetrics)
async def get_feed_metrics(
    feed_type: FeedType,
    hours: int = Query(default=24, ge=1, le=168),
    current_user: dict = Depends(get_current_user),
    feed_service: FeedService = Depends(get_feed_service)
):
    """
    Get feed performance metrics for current user
    """
    user_id = current_user["user_id"]
    
    metrics = await feed_service.get_feed_metrics(user_id, feed_type, hours)
    
    return metrics


@router.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "vignette-feed-service"}
