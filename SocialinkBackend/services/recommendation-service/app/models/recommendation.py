"""
Recommendation models for Socialink
"""
from typing import List, Optional, Dict, Any
from datetime import datetime
from pydantic import BaseModel, Field
from enum import Enum


class RecommendationType(str, Enum):
    """Socialink-specific recommendation types"""
    FRIENDS_SUGGESTIONS = "friends_suggestions"  # Friend suggestions (mutual friends priority)
    PEOPLE_YOU_MAY_KNOW = "people_you_may_know"  # People you may know (broader network)
    COMMUNITIES_FOR_YOU = "communities_for_you"  # Community recommendations


class RecommendationSource(str, Enum):
    """Source of recommendation"""
    COLLABORATIVE = "collaborative"  # Similar users
    CONTENT_BASED = "content_based"  # Similar interests
    GRAPH_BASED = "graph_based"  # Social graph
    POPULARITY = "popularity"  # Trending
    HYBRID = "hybrid"  # Combined


class UserRecommendation(BaseModel):
    """Recommended user/creator"""
    user_id: str
    username: str
    display_name: Optional[str] = None
    avatar_url: Optional[str] = None
    bio: Optional[str] = None
    
    # Stats
    follower_count: int = 0
    post_count: int = 0
    is_verified: bool = False
    is_creator: bool = False
    
    # Recommendation metadata
    score: float
    source: RecommendationSource
    reason: str  # Why recommended
    mutual_friends_count: int = 0
    mutual_friends: List[str] = Field(default_factory=list)
    common_interests: List[str] = Field(default_factory=list)
    
    # Engagement prediction
    predicted_follow_probability: Optional[float] = None


class CommunityRecommendation(BaseModel):
    """Recommended community"""
    community_id: str
    name: str
    description: Optional[str] = None
    cover_photo: Optional[str] = None
    category: str
    
    # Stats
    member_count: int = 0
    post_count: int = 0
    is_verified: bool = False
    
    # Recommendation metadata
    score: float
    source: RecommendationSource
    reason: str
    mutual_members_count: int = 0
    mutual_members: List[str] = Field(default_factory=list)
    matching_interests: List[str] = Field(default_factory=list)
    
    # Engagement prediction
    predicted_join_probability: Optional[float] = None


class RecommendationRequest(BaseModel):
    """Request for recommendations"""
    user_id: str
    type: RecommendationType
    limit: int = Field(default=20, ge=1, le=100)
    offset: int = Field(default=0, ge=0)
    
    # Filters
    exclude_ids: List[str] = Field(default_factory=list)
    categories: Optional[List[str]] = None  # For communities
    min_followers: Optional[int] = None  # For users


class RecommendationResponse(BaseModel):
    """Response with recommendations"""
    type: RecommendationType
    users: List[UserRecommendation] = Field(default_factory=list)
    communities: List[CommunityRecommendation] = Field(default_factory=list)
    next_offset: int
    has_more: bool


class UserProfile(BaseModel):
    """User profile for recommendations"""
    user_id: str
    
    # Interests
    interest_topics: List[str] = []
    interest_hashtags: List[str] = []
    interest_categories: List[str] = []
    
    # Social graph
    following_ids: List[str] = []
    follower_ids: List[str] = []
    friend_ids: List[str] = []  # Mutual follows
    blocked_ids: List[str] = []
    
    # Communities
    community_ids: List[str] = []
    
    # Behavior
    avg_engagement_rate: float = 0.0
    active_hours: List[int] = []
    content_preferences: Dict[str, float] = {}
    
    # Location (optional)
    location_city: Optional[str] = None
    location_country: Optional[str] = None


class SocialGraphFeatures(BaseModel):
    """Features from social graph analysis"""
    user_id: str
    target_id: str
    
    # Mutual connections
    mutual_friends_count: int = 0
    mutual_followers_count: int = 0
    common_communities: List[str] = []
    
    # Graph metrics
    friend_of_friend_distance: Optional[int] = None  # Degrees of separation
    social_closeness: float = 0.0  # PageRank-based
    clustering_coefficient: float = 0.0
    
    # Interaction history
    has_interacted: bool = False
    interaction_count: int = 0
    last_interaction: Optional[datetime] = None


class CollaborativeFeatures(BaseModel):
    """Features from collaborative filtering"""
    user_id: str
    target_id: str
    
    # User similarity
    cosine_similarity: float = 0.0
    jaccard_similarity: float = 0.0
    
    # Predicted affinity
    predicted_score: float = 0.0
    confidence: float = 0.0
    
    # Similar users
    similar_users_also_follow: int = 0
    similar_users_engagement: float = 0.0


class ContentFeatures(BaseModel):
    """Content-based features"""
    user_id: str
    target_id: str
    
    # Interest overlap
    topic_overlap: float = 0.0
    hashtag_overlap: float = 0.0
    category_overlap: float = 0.0
    
    # Content similarity (embedding-based)
    content_similarity: float = 0.0
    
    # Behavioral similarity
    engagement_pattern_similarity: float = 0.0


class RecommendationFeedback(BaseModel):
    """User feedback on recommendations"""
    user_id: str
    recommended_id: str
    recommendation_type: RecommendationType
    
    # Action taken
    action: str  # viewed, followed, dismissed, reported
    timestamp: datetime = Field(default_factory=datetime.utcnow)
    
    # Context
    position: int  # Position in list
    source: RecommendationSource
