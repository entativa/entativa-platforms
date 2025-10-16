"""
Recommendation API endpoints for Vignette
"""
from fastapi import APIRouter, Depends, Query, HTTPException
from typing import Optional, List

from app.models.recommendation import (
    RecommendationType, RecommendationRequest, RecommendationResponse,
    RecommendationFeedback
)
from app.services.recommendation_service import RecommendationService


router = APIRouter(prefix="/recommendations", tags=["recommendations"])


@router.get("/creators", response_model=RecommendationResponse)
async def get_creator_recommendations(
    user_id: str = Query(..., description="User ID"),
    limit: int = Query(default=20, ge=1, le=100),
    offset: int = Query(default=0, ge=0),
    exclude_ids: Optional[str] = Query(default=None, description="Comma-separated IDs to exclude"),
    # service: RecommendationService = Depends(get_recommendation_service)
):
    """
    Get creator recommendations (Creators for You)
    
    Returns top creators to follow based on:
    - Creators similar users follow
    - Popular creators
    - Friends of friends who are creators
    """
    exclude_list = exclude_ids.split(",") if exclude_ids else []
    
    request = RecommendationRequest(
        user_id=user_id,
        type=RecommendationType.CREATORS_FOR_YOU,
        limit=limit,
        offset=offset,
        exclude_ids=exclude_list
    )
    
    # In production: Use actual service
    # return await service.get_recommendations(request)
    
    return RecommendationResponse(
        type=RecommendationType.CREATORS_FOR_YOU,
        users=[],
        next_offset=offset + limit,
        has_more=False
    )


@router.get("/users", response_model=RecommendationResponse)
async def get_user_recommendations(
    user_id: str = Query(..., description="User ID"),
    limit: int = Query(default=20, ge=1, le=100),
    offset: int = Query(default=0, ge=0),
    exclude_ids: Optional[str] = Query(default=None),
    # service: RecommendationService = Depends(get_recommendation_service)
):
    """
    Get user recommendations (Suggested for You)
    
    Returns general user recommendations based on:
    - Friends of friends
    - Similar users
    - Popular accounts
    """
    exclude_list = exclude_ids.split(",") if exclude_ids else []
    
    request = RecommendationRequest(
        user_id=user_id,
        type=RecommendationType.SUGGESTED_FOR_YOU,
        limit=limit,
        offset=offset,
        exclude_ids=exclude_list
    )
    
    return RecommendationResponse(
        type=RecommendationType.SUGGESTED_FOR_YOU,
        users=[],
        next_offset=offset + limit,
        has_more=False
    )


@router.get("/communities", response_model=RecommendationResponse)
async def get_community_recommendations(
    user_id: str = Query(..., description="User ID"),
    limit: int = Query(default=20, ge=1, le=100),
    offset: int = Query(default=0, ge=0),
    categories: Optional[str] = Query(default=None, description="Comma-separated categories"),
    exclude_ids: Optional[str] = Query(default=None),
    # service: RecommendationService = Depends(get_recommendation_service)
):
    """
    Get community recommendations (Communities for You)
    
    Returns community recommendations based on:
    - Communities friends are in
    - Communities matching interests
    - Popular communities
    """
    exclude_list = exclude_ids.split(",") if exclude_ids else []
    category_list = categories.split(",") if categories else None
    
    request = RecommendationRequest(
        user_id=user_id,
        type=RecommendationType.COMMUNITIES_FOR_YOU,
        limit=limit,
        offset=offset,
        exclude_ids=exclude_list,
        categories=category_list
    )
    
    return RecommendationResponse(
        type=RecommendationType.COMMUNITIES_FOR_YOU,
        communities=[],
        next_offset=offset + limit,
        has_more=False
    )


@router.post("/feedback")
async def submit_feedback(
    feedback: RecommendationFeedback,
    # service: RecommendationService = Depends(get_recommendation_service)
):
    """
    Submit feedback on recommendations
    
    Helps improve future recommendations by tracking:
    - Which recommendations were followed/dismissed
    - Position in the list
    - Source of recommendation
    """
    # In production: Store feedback for model training
    return {"status": "success", "message": "Feedback recorded"}


@router.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "vignette-recommendation-service"}
