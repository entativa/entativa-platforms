"""
Feed schemas for Vignette Feed Service
Supports 3 feed types: Home (For You), Circle, Surprise & Delight
"""
from typing import List, Optional, Dict, Any, Literal
from datetime import datetime
from pydantic import BaseModel, Field
from enum import Enum


class FeedType(str, Enum):
    """Vignette-specific feed types"""
    HOME = "home"  # For You feed (TikTok competitor)
    CIRCLE = "circle"  # Friends + nearby
    SURPRISE_DELIGHT = "surprise_delight"  # Balanced exploration


class ContentType(str, Enum):
    """Type of content in feed"""
    POST = "post"
    TAKE = "take"  # Short-form video
    STORY = "story"
    AD = "ad"


class FeedItemSource(str, Enum):
    """Source of feed item"""
    # For Home feed
    VIRAL = "viral"  # High engagement
    INTEREST = "interest"  # User interests
    TRENDING = "trending"  # Trending content
    DISCOVERY = "discovery"  # New content
    
    # For Circle feed
    FOLLOWING = "following"  # Followed users
    FRIENDS = "friends"  # Close friends
    NEARBY = "nearby"  # Location-based
    MUTUAL = "mutual"  # Mutual connections
    
    # For Surprise & Delight
    KNOWN_INTEREST = "known_interest"  # Known preferences (60%)
    EXPLORATION = "exploration"  # Outside interests (30%)
    SURPRISE = "surprise"  # Random viral (10%)


class FeedItem(BaseModel):
    """Single item in a feed"""
    content_id: str
    content_type: ContentType
    user_id: str  # Creator
    username: str
    avatar_url: Optional[str] = None
    
    # Content data
    caption: Optional[str] = None
    media_urls: List[str] = []
    thumbnail_url: Optional[str] = None
    duration_seconds: Optional[int] = None  # For Takes
    
    # Engagement metrics
    likes_count: int = 0
    comments_count: int = 0
    shares_count: int = 0
    saves_count: int = 0
    views_count: int = 0
    
    # Ranking metadata
    score: float  # Ranking score
    source: FeedItemSource
    rank: int  # Position in feed
    
    # Timestamps
    created_at: datetime
    
    # Additional metadata
    is_following: bool = False
    has_liked: bool = False
    has_saved: bool = False
    distance_km: Optional[float] = None  # For nearby content


class FeedRequest(BaseModel):
    """Request for feed"""
    feed_type: FeedType
    user_id: str
    limit: int = Field(default=20, ge=1, le=100)
    offset: int = Field(default=0, ge=0)
    
    # Optional filters
    seen_ids: List[str] = Field(default_factory=list)  # Already seen content
    refresh: bool = False  # Force refresh
    
    # Location (for Circle feed)
    latitude: Optional[float] = None
    longitude: Optional[float] = None


class FeedResponse(BaseModel):
    """Feed response"""
    feed_type: FeedType
    items: List[FeedItem]
    next_offset: int
    has_more: bool
    refresh_token: Optional[str] = None  # For pagination


class UserSignal(BaseModel):
    """User engagement signal"""
    user_id: str
    content_id: str
    content_type: ContentType
    signal_type: Literal["view", "like", "comment", "share", "save", "skip", "hide", "report"]
    
    # Signal metadata
    time_spent_seconds: Optional[float] = None
    completion_rate: Optional[float] = None  # % of video watched
    is_organic: bool = True  # vs promoted
    
    timestamp: datetime = Field(default_factory=datetime.utcnow)


class UserProfile(BaseModel):
    """User profile for personalization"""
    user_id: str
    
    # Interests (extracted from behavior)
    interest_topics: List[str] = []
    interest_hashtags: List[str] = []
    interest_creators: List[str] = []
    interest_sounds: List[str] = []
    
    # Preferences
    preferred_content_types: List[ContentType] = []
    preferred_duration: Optional[str] = None  # short, medium, long
    
    # Social graph
    following_ids: List[str] = []
    follower_ids: List[str] = []
    friend_ids: List[str] = []  # Close friends (mutual + interaction)
    
    # Location
    last_latitude: Optional[float] = None
    last_longitude: Optional[float] = None
    location_city: Optional[str] = None
    location_country: Optional[str] = None
    
    # Engagement patterns
    avg_session_duration_minutes: float = 0.0
    active_hours: List[int] = []  # 0-23
    total_watch_time_hours: float = 0.0
    
    # Last updated
    updated_at: datetime = Field(default_factory=datetime.utcnow)


class ContentFeatures(BaseModel):
    """Content features for ranking"""
    content_id: str
    content_type: ContentType
    creator_id: str
    
    # Engagement metrics
    likes: int = 0
    comments: int = 0
    shares: int = 0
    saves: int = 0
    views: int = 0
    
    # Derived metrics
    engagement_rate: float = 0.0  # (likes + comments + shares) / views
    viral_score: float = 0.0  # Exponential growth rate
    quality_score: float = 0.0  # Based on completion rate, saves, etc
    
    # Content metadata
    hashtags: List[str] = []
    topics: List[str] = []
    mentioned_users: List[str] = []
    sound_id: Optional[str] = None
    duration_seconds: Optional[int] = None
    
    # Creator features
    creator_follower_count: int = 0
    creator_avg_engagement: float = 0.0
    
    # Temporal
    created_at: datetime
    recency_hours: float = 0.0
    
    # Location (if available)
    location_lat: Optional[float] = None
    location_lon: Optional[float] = None


class RankingConfig(BaseModel):
    """Configuration for ranking algorithm"""
    # Weight factors
    engagement_weight: float = 0.3
    recency_weight: float = 0.2
    personalization_weight: float = 0.3
    diversity_weight: float = 0.1
    social_weight: float = 0.1
    
    # Penalties
    seen_penalty: float = 0.8  # Reduce score for already seen
    creator_fatigue_limit: int = 3  # Max items from same creator
    
    # Diversity
    min_creator_diversity: int = 10  # Min unique creators in feed
    max_same_topic: int = 5  # Max consecutive items with same topic


class FeedMetrics(BaseModel):
    """Feed performance metrics"""
    feed_type: FeedType
    user_id: str
    
    # Engagement
    items_shown: int = 0
    items_engaged: int = 0  # Liked, commented, shared
    avg_time_spent_seconds: float = 0.0
    
    # Quality
    skip_rate: float = 0.0
    completion_rate: float = 0.0
    engagement_rate: float = 0.0
    
    # Session
    session_duration_minutes: float = 0.0
    items_per_minute: float = 0.0
    
    timestamp: datetime = Field(default_factory=datetime.utcnow)


class TrendingTopic(BaseModel):
    """Trending topic/hashtag"""
    name: str
    category: str
    
    # Metrics
    mention_count: int = 0
    growth_rate: float = 0.0  # % increase in last hour
    engagement_score: float = 0.0
    
    # Time window
    window_hours: int = 24
    last_updated: datetime = Field(default_factory=datetime.utcnow)
