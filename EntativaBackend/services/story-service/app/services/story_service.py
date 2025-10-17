"""
Core story service with interactive stickers
"""
from datetime import datetime, timedelta
from typing import List, Optional, Dict, Any
from bson import ObjectId
import logging

from app.models.story import (
    Story, StoryCreateRequest, StoryResponse, StoryBucket,
    StoryViewer, StoryReply, StickerInteractionRequest,
    StickerType, InteractiveSticker, PollOption
)
from app.db.mongodb import get_db
from app.db.redis_client import RedisClient
from app.config import settings

logger = logging.getLogger(__name__)


class StoryService:
    """Story service with full interactive features"""
    
    def __init__(self):
        self.redis = RedisClient
    
    async def create_story(
        self,
        user_id: str,
        request: StoryCreateRequest
    ) -> Story:
        """Create a new story"""
        db = await get_db()
        
        # Create story
        story = Story(
            user_id=user_id,
            story_type=request.story_type,
            media_id=request.media_id,
            text_content=request.text_content,
            background_color=request.background_color,
            background_gradient=request.background_gradient,
            duration=request.duration,
            music_track_id=request.music_track_id,
            privacy=request.privacy,
            close_friends_only=request.close_friends_only,
            allowed_viewer_ids=request.allowed_viewer_ids,
            hidden_from_ids=request.hidden_from_ids,
            stickers=request.stickers,
        )
        
        # Set expiry (24 hours)
        story.expires_at = story.created_at + timedelta(hours=settings.STORY_EXPIRY_HOURS)
        
        # Insert into database
        result = await db.stories.insert_one(story.dict(by_alias=True, exclude={"id"}))
        story.id = result.inserted_id
        
        # Cache user's story count
        cache_key = f"user_stories:{user_id}"
        await self.redis.delete(cache_key)
        
        # Add to story ring (sorted set by timestamp)
        await self.redis.zadd(
            "story_ring",
            {str(story.id): datetime.utcnow().timestamp()}
        )
        
        # Schedule expiration
        await self._schedule_expiration(story)
        
        logger.info(f"Created story {story.id} for user {user_id}")
        return story
    
    async def get_story(self, story_id: str, viewer_user_id: Optional[str] = None) -> Optional[Story]:
        """Get a single story"""
        db = await get_db()
        
        # Try cache first
        cache_key = f"story:{story_id}"
        cached = await self.redis.get(cache_key)
        if cached:
            story = Story(**cached)
        else:
            # Get from database
            story_doc = await db.stories.find_one({
                "_id": ObjectId(story_id),
                "is_deleted": False
            })
            
            if not story_doc:
                return None
            
            story = Story(**story_doc)
            
            # Cache it
            await self.redis.set(cache_key, story.dict(by_alias=True), ttl=settings.REDIS_TTL)
        
        # Check if expired
        if story.expires_at < datetime.utcnow():
            story.is_expired = True
            await db.stories.update_one(
                {"_id": ObjectId(story_id)},
                {"$set": {"is_expired": True}}
            )
        
        # Mark as viewed if viewer provided
        if viewer_user_id and not story.is_expired:
            await self.mark_as_viewed(story_id, viewer_user_id)
        
        return story
    
    async def get_user_stories(
        self,
        user_id: str,
        viewer_user_id: Optional[str] = None,
        include_expired: bool = False
    ) -> List[Story]:
        """Get all stories by a user"""
        db = await get_db()
        
        query: Dict[str, Any] = {
            "user_id": user_id,
            "is_deleted": False
        }
        
        if not include_expired:
            query["is_expired"] = False
            query["expires_at"] = {"$gt": datetime.utcnow()}
        
        cursor = db.stories.find(query).sort("created_at", -1)
        stories = []
        
        async for story_doc in cursor:
            story = Story(**story_doc)
            stories.append(story)
        
        return stories
    
    async def get_story_feed(
        self,
        viewer_user_id: str,
        limit: int = 50
    ) -> List[StoryBucket]:
        """Get story feed (story buckets from followed users)"""
        db = await get_db()
        
        # Get list of followed user IDs (this would come from user service)
        # For now, we'll fetch all active stories
        # In production, you'd filter by followed users
        
        pipeline = [
            {
                "$match": {
                    "is_expired": False,
                    "is_deleted": False,
                    "expires_at": {"$gt": datetime.utcnow()}
                }
            },
            {
                "$sort": {"created_at": -1}
            },
            {
                "$group": {
                    "_id": "$user_id",
                    "stories": {"$push": "$$ROOT"},
                    "story_count": {"$sum": 1},
                    "latest_created": {"$max": "$created_at"},
                    "oldest_created": {"$min": "$created_at"}
                }
            },
            {
                "$limit": limit
            }
        ]
        
        buckets = []
        async for bucket_doc in db.stories.aggregate(pipeline):
            user_id = bucket_doc["_id"]
            stories = [Story(**s) for s in bucket_doc["stories"]]
            
            # Check if user has unseen stories
            has_unseen = await self._has_unseen_stories(user_id, viewer_user_id, stories)
            
            # Get latest story for preview
            latest_story = stories[0] if stories else None
            
            bucket = StoryBucket(
                user_id=user_id,
                username="",  # Would fetch from user service
                avatar_url=None,
                is_close_friend=False,  # Would check close friends
                story_count=bucket_doc["story_count"],
                latest_story=self._to_story_response(latest_story) if latest_story else None,
                has_unseen=has_unseen,
                created_at=bucket_doc["oldest_created"],
                last_updated=bucket_doc["latest_created"]
            )
            buckets.append(bucket)
        
        # Sort: close friends first, then by has_unseen, then by last_updated
        buckets.sort(
            key=lambda b: (not b.is_close_friend, not b.has_unseen, -b.last_updated.timestamp())
        )
        
        return buckets
    
    async def mark_as_viewed(self, story_id: str, viewer_user_id: str) -> bool:
        """Mark story as viewed by user"""
        db = await get_db()
        
        # Check if already viewed
        cache_key = f"story_viewed:{story_id}:{viewer_user_id}"
        if await self.redis.exists(cache_key):
            return False
        
        # Add viewer
        viewer = StoryViewer(
            user_id=viewer_user_id,
            username="",  # Would fetch from user service
            viewed_at=datetime.utcnow()
        )
        
        result = await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            {
                "$addToSet": {"viewers": viewer.dict()},
                "$inc": {"view_count": 1}
            }
        )
        
        if result.modified_count > 0:
            # Cache the view
            await self.redis.set(cache_key, True, ttl=86400)  # 24 hours
            
            # Add to viewer set
            await self.redis.sadd(f"story_viewers:{story_id}", viewer_user_id)
            
            # Increment view count
            await self.redis.increment(f"story_views:{story_id}")
            
            logger.info(f"User {viewer_user_id} viewed story {story_id}")
            return True
        
        return False
    
    async def add_reply(
        self,
        story_id: str,
        user_id: str,
        message: str
    ) -> bool:
        """Add a reply to a story"""
        db = await get_db()
        
        reply = StoryReply(
            user_id=user_id,
            username="",  # Would fetch from user service
            message=message,
            replied_at=datetime.utcnow()
        )
        
        result = await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            {
                "$push": {"replies": reply.dict()},
                "$inc": {"reply_count": 1}
            }
        )
        
        if result.modified_count > 0:
            # Invalidate cache
            await self.redis.delete(f"story:{story_id}")
            logger.info(f"User {user_id} replied to story {story_id}")
            return True
        
        return False
    
    async def interact_with_sticker(
        self,
        story_id: str,
        user_id: str,
        interaction: StickerInteractionRequest
    ) -> bool:
        """Handle sticker interactions (polls, questions, quizzes, etc.)"""
        db = await get_db()
        
        # Get story
        story = await self.get_story(story_id)
        if not story or story.is_expired:
            return False
        
        # Validate sticker index
        if interaction.sticker_index >= len(story.stickers):
            return False
        
        sticker = story.stickers[interaction.sticker_index]
        
        # Handle different sticker types
        if sticker.sticker_type == StickerType.POLL:
            return await self._handle_poll_vote(db, story_id, interaction, sticker, user_id)
        
        elif sticker.sticker_type == StickerType.QUIZ:
            return await self._handle_quiz_answer(db, story_id, interaction, sticker, user_id)
        
        elif sticker.sticker_type == StickerType.QUESTION:
            return await self._handle_question_response(db, story_id, interaction, sticker, user_id)
        
        elif sticker.sticker_type == StickerType.SLIDER:
            return await self._handle_slider_response(db, story_id, interaction, sticker, user_id)
        
        elif sticker.sticker_type == StickerType.COUNTDOWN:
            return await self._handle_countdown_follow(db, story_id, interaction, sticker, user_id)
        
        return False
    
    async def _handle_poll_vote(
        self,
        db,
        story_id: str,
        interaction: StickerInteractionRequest,
        sticker: InteractiveSticker,
        user_id: str
    ) -> bool:
        """Handle poll voting"""
        if not sticker.poll or interaction.option_index is None:
            return False
        
        # Check if already voted
        cache_key = f"poll_voted:{story_id}:{interaction.sticker_index}:{user_id}"
        if await self.redis.exists(cache_key):
            return False
        
        # Validate option
        if interaction.option_index >= len(sticker.poll.options):
            return False
        
        # Update vote
        result = await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            {
                "$inc": {
                    f"stickers.{interaction.sticker_index}.poll.options.{interaction.option_index}.votes": 1,
                    f"stickers.{interaction.sticker_index}.poll.total_votes": 1
                },
                "$push": {
                    f"stickers.{interaction.sticker_index}.poll.options.{interaction.option_index}.voter_ids": user_id
                }
            }
        )
        
        if result.modified_count > 0:
            # Mark as voted
            await self.redis.set(cache_key, True, ttl=86400)
            
            # Recalculate percentages
            await self._recalculate_poll_percentages(db, story_id, interaction.sticker_index)
            
            logger.info(f"User {user_id} voted on poll in story {story_id}")
            return True
        
        return False
    
    async def _handle_quiz_answer(
        self,
        db,
        story_id: str,
        interaction: StickerInteractionRequest,
        sticker: InteractiveSticker,
        user_id: str
    ) -> bool:
        """Handle quiz answer"""
        if not sticker.quiz or interaction.selected_option is None:
            return False
        
        # Check if already answered
        cache_key = f"quiz_answered:{story_id}:{interaction.sticker_index}:{user_id}"
        if await self.redis.exists(cache_key):
            return False
        
        # Validate option
        if interaction.selected_option >= len(sticker.quiz.options):
            return False
        
        is_correct = interaction.selected_option == sticker.quiz.correct_answer_index
        
        # Update quiz
        update_ops = {
            "$inc": {
                f"stickers.{interaction.sticker_index}.quiz.options.{interaction.selected_option}.selected_count": 1,
                f"stickers.{interaction.sticker_index}.quiz.total_attempts": 1
            }
        }
        
        if is_correct:
            update_ops["$inc"][f"stickers.{interaction.sticker_index}.quiz.correct_attempts"] = 1
        
        result = await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            update_ops
        )
        
        if result.modified_count > 0:
            # Mark as answered
            await self.redis.set(cache_key, True, ttl=86400)
            logger.info(f"User {user_id} answered quiz in story {story_id} - Correct: {is_correct}")
            return True
        
        return False
    
    async def _handle_question_response(
        self,
        db,
        story_id: str,
        interaction: StickerInteractionRequest,
        sticker: InteractiveSticker,
        user_id: str
    ) -> bool:
        """Handle question response"""
        if not sticker.question or not interaction.answer:
            return False
        
        response = {
            "user_id": user_id,
            "answer": interaction.answer,
            "timestamp": datetime.utcnow()
        }
        
        result = await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            {
                "$push": {
                    f"stickers.{interaction.sticker_index}.question.responses": response
                },
                "$inc": {
                    f"stickers.{interaction.sticker_index}.question.response_count": 1
                }
            }
        )
        
        if result.modified_count > 0:
            logger.info(f"User {user_id} responded to question in story {story_id}")
            return True
        
        return False
    
    async def _handle_slider_response(
        self,
        db,
        story_id: str,
        interaction: StickerInteractionRequest,
        sticker: InteractiveSticker,
        user_id: str
    ) -> bool:
        """Handle slider response"""
        if not sticker.slider or interaction.slider_value is None:
            return False
        
        # Check if already responded
        cache_key = f"slider_responded:{story_id}:{interaction.sticker_index}:{user_id}"
        if await self.redis.exists(cache_key):
            return False
        
        response = {
            "user_id": user_id,
            "value": interaction.slider_value,
            "timestamp": datetime.utcnow()
        }
        
        result = await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            {
                "$push": {
                    f"stickers.{interaction.sticker_index}.slider.responses": response
                }
            }
        )
        
        if result.modified_count > 0:
            # Recalculate average
            await self._recalculate_slider_average(db, story_id, interaction.sticker_index)
            
            # Mark as responded
            await self.redis.set(cache_key, True, ttl=86400)
            
            logger.info(f"User {user_id} responded to slider in story {story_id}")
            return True
        
        return False
    
    async def _handle_countdown_follow(
        self,
        db,
        story_id: str,
        interaction: StickerInteractionRequest,
        sticker: InteractiveSticker,
        user_id: str
    ) -> bool:
        """Handle countdown follow/unfollow"""
        if not sticker.countdown or interaction.follow_countdown is None:
            return False
        
        if interaction.follow_countdown:
            # Follow countdown
            result = await db.stories.update_one(
                {"_id": ObjectId(story_id)},
                {
                    "$addToSet": {
                        f"stickers.{interaction.sticker_index}.countdown.followers": user_id
                    },
                    "$inc": {
                        f"stickers.{interaction.sticker_index}.countdown.follower_count": 1
                    }
                }
            )
        else:
            # Unfollow countdown
            result = await db.stories.update_one(
                {"_id": ObjectId(story_id)},
                {
                    "$pull": {
                        f"stickers.{interaction.sticker_index}.countdown.followers": user_id
                    },
                    "$inc": {
                        f"stickers.{interaction.sticker_index}.countdown.follower_count": -1
                    }
                }
            )
        
        if result.modified_count > 0:
            logger.info(f"User {user_id} {'followed' if interaction.follow_countdown else 'unfollowed'} countdown in story {story_id}")
            return True
        
        return False
    
    async def delete_story(self, story_id: str, user_id: str) -> bool:
        """Delete a story (soft delete)"""
        db = await get_db()
        
        result = await db.stories.update_one(
            {
                "_id": ObjectId(story_id),
                "user_id": user_id
            },
            {
                "$set": {"is_deleted": True, "updated_at": datetime.utcnow()}
            }
        )
        
        if result.modified_count > 0:
            # Invalidate cache
            await self.redis.delete(f"story:{story_id}")
            await self.redis.delete(f"user_stories:{user_id}")
            logger.info(f"Deleted story {story_id}")
            return True
        
        return False
    
    # Helper methods
    
    async def _has_unseen_stories(
        self,
        story_user_id: str,
        viewer_user_id: str,
        stories: List[Story]
    ) -> bool:
        """Check if user has unseen stories from this user"""
        for story in stories:
            cache_key = f"story_viewed:{str(story.id)}:{viewer_user_id}"
            if not await self.redis.exists(cache_key):
                return True
        return False
    
    def _to_story_response(self, story: Story) -> StoryResponse:
        """Convert Story to StoryResponse"""
        time_remaining = None
        if not story.is_expired:
            time_remaining = int((story.expires_at - datetime.utcnow()).total_seconds())
        
        return StoryResponse(
            id=str(story.id),
            user_id=story.user_id,
            story_type=story.story_type,
            media_url=story.media_url,
            thumbnail_url=story.thumbnail_url,
            text_content=story.text_content,
            background_color=story.background_color,
            background_gradient=story.background_gradient,
            duration=story.duration,
            stickers=story.stickers,
            privacy=story.privacy,
            close_friends_only=story.close_friends_only,
            view_count=story.view_count,
            reply_count=story.reply_count,
            has_viewed=False,  # Would check viewer
            created_at=story.created_at,
            expires_at=story.expires_at,
            time_remaining=time_remaining
        )
    
    async def _schedule_expiration(self, story: Story):
        """Schedule story expiration"""
        # In production, use Celery or APScheduler
        # For now, we'll just set a Redis key with TTL
        expiry_key = f"story_expiry:{str(story.id)}"
        ttl = int((story.expires_at - datetime.utcnow()).total_seconds())
        await self.redis.set(expiry_key, str(story.id), ttl=ttl)
    
    async def _recalculate_poll_percentages(self, db, story_id: str, sticker_index: int):
        """Recalculate poll percentages after vote"""
        story = await db.stories.find_one({"_id": ObjectId(story_id)})
        if not story:
            return
        
        sticker = story["stickers"][sticker_index]
        if "poll" not in sticker:
            return
        
        poll = sticker["poll"]
        total = poll.get("total_votes", 0)
        
        if total == 0:
            return
        
        # Calculate percentages
        for i, option in enumerate(poll["options"]):
            percentage = (option["votes"] / total) * 100
            await db.stories.update_one(
                {"_id": ObjectId(story_id)},
                {
                    "$set": {
                        f"stickers.{sticker_index}.poll.options.{i}.percentage": round(percentage, 1)
                    }
                }
            )
    
    async def _recalculate_slider_average(self, db, story_id: str, sticker_index: int):
        """Recalculate slider average after response"""
        story = await db.stories.find_one({"_id": ObjectId(story_id)})
        if not story:
            return
        
        sticker = story["stickers"][sticker_index]
        if "slider" not in sticker:
            return
        
        slider = sticker["slider"]
        responses = slider.get("responses", [])
        
        if not responses:
            return
        
        # Calculate average
        total = sum(r["value"] for r in responses)
        average = total / len(responses)
        
        await db.stories.update_one(
            {"_id": ObjectId(story_id)},
            {
                "$set": {
                    f"stickers.{sticker_index}.slider.average_value": round(average, 1)
                }
            }
        )
