"""
Highlight models for Socialink
Highlights are permanent story collections
"""
from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel, Field
from bson import ObjectId

from .story import PyObjectId, StoryResponse


class HighlightCover(BaseModel):
    """Highlight cover settings"""
    story_id: Optional[str] = None
    media_url: Optional[str] = None
    emoji: Optional[str] = None
    color: str = "#FF6B6B"  # Default gradient color


class Highlight(BaseModel):
    """Highlight collection model"""
    id: Optional[PyObjectId] = Field(alias="_id", default=None)
    user_id: str = Field(..., description="Owner of the highlight")
    
    # Basic info
    title: str = Field(..., min_length=1, max_length=50)
    description: Optional[str] = Field(None, max_length=200)
    
    # Cover
    cover: HighlightCover = Field(default_factory=HighlightCover)
    
    # Stories
    story_ids: List[str] = Field(default_factory=list)
    story_count: int = Field(default=0, ge=0)
    
    # Metadata
    is_archived: bool = False
    is_pinned: bool = False
    order_index: int = Field(default=0)  # For custom ordering
    
    # Timestamps
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)
    
    class Config:
        populate_by_name = True
        arbitrary_types_allowed = True
        json_encoders = {ObjectId: str, datetime: lambda v: v.isoformat()}


class HighlightCreateRequest(BaseModel):
    """Request model for creating a highlight"""
    title: str = Field(..., min_length=1, max_length=50)
    description: Optional[str] = Field(None, max_length=200)
    story_ids: List[str] = Field(default_factory=list)
    cover: Optional[HighlightCover] = None
    is_pinned: bool = False


class HighlightUpdateRequest(BaseModel):
    """Request model for updating a highlight"""
    title: Optional[str] = Field(None, min_length=1, max_length=50)
    description: Optional[str] = Field(None, max_length=200)
    cover: Optional[HighlightCover] = None
    is_pinned: Optional[bool] = None
    is_archived: Optional[bool] = None


class HighlightResponse(BaseModel):
    """Response model for a highlight"""
    id: str
    user_id: str
    title: str
    description: Optional[str] = None
    cover: HighlightCover
    story_count: int
    is_pinned: bool
    created_at: datetime
    updated_at: datetime
    
    # Optionally include stories
    stories: Optional[List[StoryResponse]] = None


class HighlightStoriesRequest(BaseModel):
    """Request to add/remove stories from highlight"""
    story_ids: List[str]
