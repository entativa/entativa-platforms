"""
Personalization Service - User profiling and learning
"""
from typing import List, Optional, Dict
from datetime import datetime, timedelta
from collections import Counter, defaultdict
import logging

from app.schemas.feed import (
    UserProfile, UserSignal, ContentFeatures,
    FeedType, FeedMetrics, ContentType
)


logger = logging.getLogger(__name__)


class PersonalizationService:
    """
    Handles user profiling, interest extraction, and signal processing
    """
    
    def __init__(self, db, redis_client):
        self.db = db
        self.redis = redis_client
    
    async def get_user_profile(self, user_id: str) -> UserProfile:
        """
        Get or build user profile
        """
        # Check cache
        cache_key = f"user_profile:{user_id}"
        cached = await self.redis.get(cache_key)
        if cached:
            logger.debug(f"Cache hit for user profile {user_id}")
            return UserProfile.model_validate_json(cached)
        
        # Build from DB
        profile = await self._build_user_profile(user_id)
        
        # Cache for 1 hour
        await self.redis.setex(cache_key, 3600, profile.model_dump_json())
        
        return profile
    
    async def _build_user_profile(self, user_id: str) -> UserProfile:
        """
        Build user profile from historical data
        """
        # Get user's engagement history (last 30 days)
        signals = await self._get_user_signals(user_id, days=30)
        
        # Extract interests
        interest_topics = await self._extract_interest_topics(signals)
        interest_hashtags = await self._extract_interest_hashtags(signals)
        interest_creators = await self._extract_interest_creators(signals)
        interest_sounds = await self._extract_interest_sounds(signals)
        
        # Get social graph
        following_ids = await self._get_user_following(user_id)
        follower_ids = await self._get_user_followers(user_id)
        friend_ids = await self._get_user_friends(user_id, following_ids, follower_ids)
        
        # Get location
        last_lat, last_lon = await self._get_last_location(user_id)
        
        # Get engagement patterns
        avg_session = await self._calculate_avg_session(user_id)
        active_hours = await self._extract_active_hours(signals)
        total_watch_time = await self._calculate_watch_time(signals)
        
        # Get preferred content types
        preferred_types = await self._extract_preferred_types(signals)
        
        profile = UserProfile(
            user_id=user_id,
            interest_topics=interest_topics,
            interest_hashtags=interest_hashtags,
            interest_creators=interest_creators,
            interest_sounds=interest_sounds,
            preferred_content_types=preferred_types,
            following_ids=following_ids,
            follower_ids=follower_ids,
            friend_ids=friend_ids,
            last_latitude=last_lat,
            last_longitude=last_lon,
            avg_session_duration_minutes=avg_session,
            active_hours=active_hours,
            total_watch_time_hours=total_watch_time,
            updated_at=datetime.utcnow()
        )
        
        logger.info(f"Built profile for user {user_id}: {len(interest_topics)} topics, {len(interest_creators)} creators")
        
        return profile
    
    async def process_signal(self, signal: UserSignal):
        """
        Process user engagement signal and update profile
        """
        logger.debug(f"Processing signal: {signal.signal_type} for user {signal.user_id}")
        
        # Store signal in DB (for batch processing)
        await self._store_signal(signal)
        
        # Update Redis counters (real-time)
        await self._update_signal_counters(signal)
        
        # Invalidate user profile cache (will rebuild on next request)
        cache_key = f"user_profile:{signal.user_id}"
        await self.redis.delete(cache_key)
        
        # Track for feed metrics
        await self._track_feed_metric(signal)
        
        logger.debug(f"Signal processed successfully")
    
    async def _store_signal(self, signal: UserSignal):
        """Store signal in database"""
        # In production: INSERT INTO user_signals ...
        # For now, just log
        logger.debug(f"Storing signal: {signal.signal_type}")
    
    async def _update_signal_counters(self, signal: UserSignal):
        """Update Redis counters for real-time stats"""
        # User signal counter
        user_key = f"user_signals:{signal.user_id}:{signal.signal_type}"
        await self.redis.incr(user_key)
        
        # Content signal counter
        content_key = f"content_signals:{signal.content_id}:{signal.signal_type}"
        await self.redis.incr(content_key)
        
        # Track liked content
        if signal.signal_type == "like":
            await self.redis.sadd(f"user_likes:{signal.user_id}", signal.content_id)
        elif signal.signal_type == "save":
            await self.redis.sadd(f"user_saves:{signal.user_id}", signal.content_id)
    
    async def _track_feed_metric(self, signal: UserSignal):
        """Track signal for feed metrics"""
        # Store in time-series for analytics
        metric_key = f"feed_metrics:{signal.user_id}:{datetime.utcnow().date()}"
        
        await self.redis.hincrby(metric_key, "total_signals", 1)
        await self.redis.hincrby(metric_key, f"signal_{signal.signal_type}", 1)
        
        if signal.time_spent_seconds:
            await self.redis.hincrbyfloat(
                metric_key, "total_time_spent", signal.time_spent_seconds
            )
    
    async def _get_user_signals(self, user_id: str, days: int) -> List[UserSignal]:
        """Get user's recent signals"""
        # In production: Query from DB
        # SELECT * FROM user_signals WHERE user_id = ... AND created_at > NOW() - days
        # For now, return empty
        return []
    
    async def _extract_interest_topics(self, signals: List[UserSignal]) -> List[str]:
        """Extract top interest topics from signals"""
        # Count topics from liked/saved content
        # In production: JOIN with content to get topics
        # For now, return mock
        return ["travel", "food", "fitness", "music", "fashion"]
    
    async def _extract_interest_hashtags(self, signals: List[UserSignal]) -> List[str]:
        """Extract top interest hashtags"""
        return ["#vignette", "#takes", "#viral", "#trending"]
    
    async def _extract_interest_creators(self, signals: List[UserSignal]) -> List[str]:
        """Extract favorite creators (most engaged with)"""
        # Count creator IDs from signals
        creator_counter = Counter()
        for signal in signals:
            if signal.signal_type in ["like", "save", "comment"]:
                # In production: Get creator_id from content
                creator_counter["mock_creator"] += 1
        
        # Return top 20
        return [creator for creator, count in creator_counter.most_common(20)]
    
    async def _extract_interest_sounds(self, signals: List[UserSignal]) -> List[str]:
        """Extract favorite sounds/music"""
        return []
    
    async def _get_user_following(self, user_id: str) -> List[str]:
        """Get users that this user follows"""
        # Check cache
        cache_key = f"user_following:{user_id}"
        cached = await self.redis.smembers(cache_key)
        if cached:
            return list(cached)
        
        # Query from user service or DB
        # For now, return mock
        following = [f"user_{i}" for i in range(10)]
        
        # Cache for 1 hour
        if following:
            await self.redis.sadd(cache_key, *following)
            await self.redis.expire(cache_key, 3600)
        
        return following
    
    async def get_user_following(self, user_id: str) -> List[str]:
        """Public method to get following list"""
        return await self._get_user_following(user_id)
    
    async def _get_user_followers(self, user_id: str) -> List[str]:
        """Get users that follow this user"""
        # Similar to _get_user_following
        return []
    
    async def _get_user_friends(
        self, user_id: str, following: List[str], followers: List[str]
    ) -> List[str]:
        """
        Get close friends (mutual follows + high interaction)
        """
        # Find mutual follows
        following_set = set(following)
        followers_set = set(followers)
        mutual = list(following_set & followers_set)
        
        # In production: Also filter by interaction frequency
        # For now, just return mutuals
        return mutual
    
    async def _get_last_location(self, user_id: str) -> tuple:
        """Get user's last known location"""
        # In production: Query from location_history table
        return None, None
    
    async def _calculate_avg_session(self, user_id: str) -> float:
        """Calculate average session duration"""
        # In production: Query session data
        return 15.0  # 15 minutes average
    
    async def _extract_active_hours(self, signals: List[UserSignal]) -> List[int]:
        """Extract hours of day when user is most active"""
        if not signals:
            return []
        
        hour_counter = Counter()
        for signal in signals:
            hour = signal.timestamp.hour
            hour_counter[hour] += 1
        
        # Return top 5 hours
        return [hour for hour, count in hour_counter.most_common(5)]
    
    async def _calculate_watch_time(self, signals: List[UserSignal]) -> float:
        """Calculate total watch time in hours"""
        total_seconds = sum(
            signal.time_spent_seconds or 0
            for signal in signals
            if signal.signal_type == "view"
        )
        return total_seconds / 3600.0
    
    async def _extract_preferred_types(self, signals: List[UserSignal]) -> List[ContentType]:
        """Extract preferred content types"""
        if not signals:
            return []
        
        type_counter = Counter()
        for signal in signals:
            if signal.signal_type in ["like", "save"]:
                type_counter[signal.content_type] += 1
        
        # Return top 3
        return [content_type for content_type, count in type_counter.most_common(3)]
    
    async def get_feed_metrics(
        self, user_id: str, feed_type: FeedType, since: datetime
    ) -> FeedMetrics:
        """
        Get feed performance metrics
        """
        # Query from Redis or DB
        date_key = since.date()
        metric_key = f"feed_metrics:{user_id}:{date_key}"
        
        metrics_data = await self.redis.hgetall(metric_key)
        
        if not metrics_data:
            return FeedMetrics(
                feed_type=feed_type,
                user_id=user_id
            )
        
        total_signals = int(metrics_data.get("total_signals", 0))
        signal_like = int(metrics_data.get("signal_like", 0))
        signal_comment = int(metrics_data.get("signal_comment", 0))
        signal_share = int(metrics_data.get("signal_share", 0))
        signal_skip = int(metrics_data.get("signal_skip", 0))
        total_time = float(metrics_data.get("total_time_spent", 0))
        
        items_engaged = signal_like + signal_comment + signal_share
        skip_rate = signal_skip / max(total_signals, 1)
        engagement_rate = items_engaged / max(total_signals, 1)
        avg_time_spent = total_time / max(total_signals, 1)
        
        return FeedMetrics(
            feed_type=feed_type,
            user_id=user_id,
            items_shown=total_signals,
            items_engaged=items_engaged,
            avg_time_spent_seconds=avg_time_spent,
            skip_rate=skip_rate,
            engagement_rate=engagement_rate,
            timestamp=datetime.utcnow()
        )
