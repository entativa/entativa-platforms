"""
Viewer models for story analytics
"""
from datetime import datetime
from typing import Optional, Dict, Any
from pydantic import BaseModel, Field


class ViewerInsight(BaseModel):
    """Detailed viewer insights"""
    user_id: str
    username: str
    avatar_url: Optional[str] = None
    
    # Viewing behavior
    view_count: int = 0
    last_viewed_at: datetime
    average_view_duration: float = 0.0  # Seconds
    
    # Interactions
    total_interactions: int = 0
    poll_votes: int = 0
    quiz_attempts: int = 0
    question_responses: int = 0
    slider_responses: int = 0
    replies: int = 0
    
    # Engagement
    is_close_friend: bool = False
    engagement_score: float = 0.0  # Calculated score


class StoryInsights(BaseModel):
    """Analytics for a story"""
    story_id: str
    
    # Views
    total_views: int = 0
    unique_viewers: int = 0
    reach_percentage: float = 0.0
    
    # Engagement
    total_interactions: int = 0
    interaction_rate: float = 0.0
    replies: int = 0
    
    # Demographics (if available)
    top_viewers: list[ViewerInsight] = Field(default_factory=list)
    
    # Time-based
    views_by_hour: Dict[int, int] = Field(default_factory=dict)
    peak_viewing_hour: Optional[int] = None
    
    # Sticker performance
    sticker_interactions: Dict[str, int] = Field(default_factory=dict)


class CloseFriend(BaseModel):
    """Close friend model"""
    user_id: str
    added_at: datetime = Field(default_factory=datetime.utcnow)
    is_mutual: bool = False
