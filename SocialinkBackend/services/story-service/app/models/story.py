"""
Story models for Socialink
"""
from datetime import datetime, timedelta
from typing import Optional, List, Dict, Any
from enum import Enum
from pydantic import BaseModel, Field, validator
from bson import ObjectId


class PyObjectId(ObjectId):
    """Custom ObjectId type for Pydantic"""
    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        if not ObjectId.is_valid(v):
            raise ValueError("Invalid ObjectId")
        return ObjectId(v)

    @classmethod
    def __modify_schema__(cls, field_schema):
        field_schema.update(type="string")


class StoryType(str, Enum):
    """Story content types"""
    IMAGE = "image"
    VIDEO = "video"
    TEXT = "text"


class StickerType(str, Enum):
    """Interactive sticker types"""
    POLL = "poll"
    QUESTION = "question"
    QUIZ = "quiz"
    COUNTDOWN = "countdown"
    SLIDER = "slider"
    MENTION = "mention"
    LOCATION = "location"
    HASHTAG = "hashtag"
    MUSIC = "music"
    LINK = "link"


class StoryPrivacy(str, Enum):
    """Story privacy levels"""
    PUBLIC = "public"
    FOLLOWERS = "followers"
    CLOSE_FRIENDS = "close_friends"
    CUSTOM = "custom"


class PollOption(BaseModel):
    """Poll option model"""
    text: str = Field(..., max_length=50)
    votes: int = Field(default=0, ge=0)
    voter_ids: List[str] = Field(default_factory=list)
    percentage: float = Field(default=0.0, ge=0, le=100)


class Poll(BaseModel):
    """Poll sticker"""
    question: str = Field(..., max_length=200)
    options: List[PollOption] = Field(..., min_items=2, max_items=4)
    total_votes: int = Field(default=0, ge=0)
    created_at: datetime = Field(default_factory=datetime.utcnow)


class QuizOption(BaseModel):
    """Quiz option model"""
    text: str = Field(..., max_length=50)
    is_correct: bool = False
    selected_count: int = Field(default=0, ge=0)


class Quiz(BaseModel):
    """Quiz sticker"""
    question: str = Field(..., max_length=200)
    options: List[QuizOption] = Field(..., min_items=2, max_items=4)
    correct_answer_index: int = Field(..., ge=0)
    total_attempts: int = Field(default=0, ge=0)
    correct_attempts: int = Field(default=0, ge=0)


class Question(BaseModel):
    """Question sticker"""
    text: str = Field(..., max_length=200)
    responses: List[Dict[str, Any]] = Field(default_factory=list)  # {user_id, answer, timestamp}
    response_count: int = Field(default=0, ge=0)


class Countdown(BaseModel):
    """Countdown sticker"""
    title: str = Field(..., max_length=100)
    end_time: datetime
    is_over: bool = False
    follower_count: int = Field(default=0, ge=0)  # Users following countdown
    followers: List[str] = Field(default_factory=list)


class Slider(BaseModel):
    """Slider sticker (e.g., "How much do you like this?")"""
    question: str = Field(..., max_length=200)
    emoji: str = Field(default="😍")
    min_value: int = 0
    max_value: int = 100
    responses: List[Dict[str, Any]] = Field(default_factory=list)  # {user_id, value, timestamp}
    average_value: float = Field(default=50.0, ge=0, le=100)


class Mention(BaseModel):
    """User mention sticker"""
    user_id: str
    username: str
    x: float = Field(..., ge=0, le=1)  # Normalized position
    y: float = Field(..., ge=0, le=1)
    width: float = Field(default=0.3, ge=0, le=1)
    height: float = Field(default=0.1, ge=0, le=1)


class Location(BaseModel):
    """Location sticker"""
    name: str
    latitude: Optional[float] = None
    longitude: Optional[float] = None
    place_id: Optional[str] = None


class Music(BaseModel):
    """Music sticker"""
    track_id: str
    track_name: str
    artist_name: str
    start_time: float = 0.0  # In seconds
    duration: float = 15.0


class InteractiveSticker(BaseModel):
    """Base interactive sticker"""
    sticker_type: StickerType
    x: float = Field(..., ge=0, le=1)  # Normalized position (0-1)
    y: float = Field(..., ge=0, le=1)
    width: float = Field(default=0.5, ge=0, le=1)
    height: float = Field(default=0.3, ge=0, le=1)
    rotation: float = Field(default=0, ge=-180, le=180)
    
    # Sticker data (one of these will be populated based on type)
    poll: Optional[Poll] = None
    quiz: Optional[Quiz] = None
    question: Optional[Question] = None
    countdown: Optional[Countdown] = None
    slider: Optional[Slider] = None
    mention: Optional[Mention] = None
    location: Optional[Location] = None
    music: Optional[Music] = None


class StoryViewer(BaseModel):
    """Story viewer model"""
    user_id: str
    username: str
    avatar_url: Optional[str] = None
    viewed_at: datetime = Field(default_factory=datetime.utcnow)
    interaction_type: Optional[str] = None  # liked, replied, voted, etc.


class StoryReply(BaseModel):
    """Reply to a story"""
    user_id: str
    username: str
    message: str = Field(..., max_length=500)
    replied_at: datetime = Field(default_factory=datetime.utcnow)
    is_read: bool = False


class Story(BaseModel):
    """Main story model"""
    id: Optional[PyObjectId] = Field(alias="_id", default=None)
    user_id: str = Field(..., description="Owner of the story")
    
    # Content
    story_type: StoryType
    media_id: Optional[str] = None  # From media service
    media_url: Optional[str] = None  # Generated URL
    thumbnail_url: Optional[str] = None
    text_content: Optional[str] = Field(None, max_length=2000)
    background_color: Optional[str] = None
    background_gradient: Optional[List[str]] = None
    
    # Interactive elements
    stickers: List[InteractiveSticker] = Field(default_factory=list)
    
    # Metadata
    duration: int = Field(default=5, ge=1, le=15)  # Seconds to display
    music_track_id: Optional[str] = None
    
    # Privacy
    privacy: StoryPrivacy = StoryPrivacy.FOLLOWERS
    close_friends_only: bool = False
    allowed_viewer_ids: List[str] = Field(default_factory=list)  # For CUSTOM privacy
    hidden_from_ids: List[str] = Field(default_factory=list)
    
    # Engagement
    viewers: List[StoryViewer] = Field(default_factory=list)
    view_count: int = Field(default=0, ge=0)
    replies: List[StoryReply] = Field(default_factory=list)
    reply_count: int = Field(default=0, ge=0)
    
    # Status
    is_archived: bool = False
    is_deleted: bool = False
    is_expired: bool = False
    
    # Timestamps
    created_at: datetime = Field(default_factory=datetime.utcnow)
    expires_at: datetime = Field(default=None)
    updated_at: datetime = Field(default_factory=datetime.utcnow)
    
    class Config:
        populate_by_name = True
        arbitrary_types_allowed = True
        json_encoders = {ObjectId: str, datetime: lambda v: v.isoformat()}
    
    @validator("expires_at", always=True, pre=True)
    def set_expires_at(cls, v, values):
        """Auto-set expiry to 24 hours from creation"""
        if v is None and "created_at" in values:
            return values["created_at"] + timedelta(hours=24)
        return v


class StoryCreateRequest(BaseModel):
    """Request model for creating a story"""
    story_type: StoryType
    media_id: Optional[str] = None
    text_content: Optional[str] = Field(None, max_length=2000)
    background_color: Optional[str] = None
    background_gradient: Optional[List[str]] = None
    duration: int = Field(default=5, ge=1, le=15)
    music_track_id: Optional[str] = None
    privacy: StoryPrivacy = StoryPrivacy.FOLLOWERS
    close_friends_only: bool = False
    allowed_viewer_ids: List[str] = Field(default_factory=list)
    hidden_from_ids: List[str] = Field(default_factory=list)
    stickers: List[InteractiveSticker] = Field(default_factory=list)


class StoryResponse(BaseModel):
    """Response model for a story"""
    id: str
    user_id: str
    story_type: StoryType
    media_url: Optional[str] = None
    thumbnail_url: Optional[str] = None
    text_content: Optional[str] = None
    background_color: Optional[str] = None
    background_gradient: Optional[List[str]] = None
    duration: int
    stickers: List[InteractiveSticker] = Field(default_factory=list)
    privacy: StoryPrivacy
    close_friends_only: bool
    view_count: int
    reply_count: int
    has_viewed: bool = False  # Whether current user has viewed
    created_at: datetime
    expires_at: datetime
    time_remaining: Optional[int] = None  # Seconds until expiry


class StoryBucket(BaseModel):
    """Story bucket (user's story collection)"""
    user_id: str
    username: str
    avatar_url: Optional[str] = None
    is_close_friend: bool = False
    story_count: int = 0
    latest_story: Optional[StoryResponse] = None
    has_unseen: bool = False  # Has stories user hasn't viewed
    created_at: datetime  # Time of oldest story
    last_updated: datetime  # Time of newest story


class StickerInteractionRequest(BaseModel):
    """Request for interacting with stickers"""
    story_id: str
    sticker_index: int
    
    # For polls
    option_index: Optional[int] = None
    
    # For questions
    answer: Optional[str] = Field(None, max_length=500)
    
    # For quiz
    selected_option: Optional[int] = None
    
    # For sliders
    slider_value: Optional[float] = Field(None, ge=0, le=100)
    
    # For countdowns
    follow_countdown: Optional[bool] = None
