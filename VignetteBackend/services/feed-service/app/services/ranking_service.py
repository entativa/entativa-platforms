"""
Ranking Service - Candidate retrieval and content scoring
"""
from typing import List, Optional, Dict, Any
from datetime import datetime, timedelta
import logging

from app.schemas.feed import ContentFeatures, UserProfile, ContentType, UserSignal


logger = logging.getLogger(__name__)


class RankingService:
    """
    Handles candidate retrieval and content data fetching
    """
    
    def __init__(self, db, redis_client, elasticsearch_client=None):
        self.db = db
        self.redis = redis_client
        self.es = elasticsearch_client
    
    async def get_home_candidates(
        self, user_profile: UserProfile, limit: int = 500
    ) -> List[ContentFeatures]:
        """
        Get candidates for Home feed (For You)
        Mix of viral, trending, personalized
        """
        candidates = []
        
        # 1. Viral content (high engagement, last 24h)
        viral_items = await self._get_viral_content(limit=200)
        candidates.extend(viral_items)
        
        # 2. Trending topics (user interests)
        trending_items = await self._get_trending_by_interests(
            user_profile.interest_topics,
            user_profile.interest_hashtags,
            limit=150
        )
        candidates.extend(trending_items)
        
        # 3. Similar to liked content (collaborative filtering)
        similar_items = await self._get_similar_content(
            user_profile.user_id, limit=100
        )
        candidates.extend(similar_items)
        
        # 4. Fresh content (last 6 hours)
        fresh_items = await self._get_fresh_content(hours=6, limit=50)
        candidates.extend(fresh_items)
        
        # Deduplicate
        seen_ids = set()
        unique_candidates = []
        for item in candidates:
            if item.content_id not in seen_ids:
                unique_candidates.append(item)
                seen_ids.add(item.content_id)
        
        logger.info(f"Retrieved {len(unique_candidates)} candidates for Home feed")
        return unique_candidates[:limit]
    
    async def get_circle_candidates(
        self,
        user_profile: UserProfile,
        latitude: Optional[float],
        longitude: Optional[float],
        limit: int = 300
    ) -> List[ContentFeatures]:
        """
        Get candidates for Circle feed
        Social graph + location
        """
        candidates = []
        
        # 1. Friends' content (close connections)
        if user_profile.friend_ids:
            friend_content = await self._get_content_by_users(
                user_profile.friend_ids, limit=150
            )
            candidates.extend(friend_content)
        
        # 2. Following content
        if user_profile.following_ids:
            following_content = await self._get_content_by_users(
                user_profile.following_ids, limit=100
            )
            candidates.extend(following_content)
        
        # 3. Nearby content (if location provided)
        if latitude and longitude:
            nearby_content = await self._get_nearby_content(
                latitude, longitude, radius_km=50, limit=50
            )
            candidates.extend(nearby_content)
        
        # Deduplicate
        seen_ids = set()
        unique_candidates = []
        for item in candidates:
            if item.content_id not in seen_ids:
                unique_candidates.append(item)
                seen_ids.add(item.content_id)
        
        logger.info(f"Retrieved {len(unique_candidates)} candidates for Circle feed")
        return unique_candidates[:limit]
    
    async def get_surprise_candidates(
        self, user_profile: UserProfile, limit: int = 400
    ) -> List[ContentFeatures]:
        """
        Get candidates for Surprise & Delight feed
        Diverse mix
        """
        candidates = []
        
        # 1. Known interests (60%)
        interest_content = await self._get_trending_by_interests(
            user_profile.interest_topics,
            user_profile.interest_hashtags,
            limit=240
        )
        candidates.extend(interest_content)
        
        # 2. Exploration (adjacent topics) (30%)
        exploration_content = await self._get_exploration_content(
            user_profile, limit=120
        )
        candidates.extend(exploration_content)
        
        # 3. Surprise (viral, random) (10%)
        surprise_content = await self._get_viral_content(limit=40)
        candidates.extend(surprise_content)
        
        # Deduplicate
        seen_ids = set()
        unique_candidates = []
        for item in candidates:
            if item.content_id not in seen_ids:
                unique_candidates.append(item)
                seen_ids.add(item.content_id)
        
        logger.info(f"Retrieved {len(unique_candidates)} candidates for Surprise feed")
        return unique_candidates[:limit]
    
    async def _get_viral_content(self, limit: int) -> List[ContentFeatures]:
        """
        Get viral content (high engagement rate, recent)
        """
        # Check cache first
        cache_key = f"viral_content:{limit}"
        cached = await self.redis.get(cache_key)
        if cached:
            logger.debug(f"Cache hit for viral content")
            return self._deserialize_features(cached)
        
        # Query from DB
        # In production: SELECT * FROM content WHERE 
        #   created_at > NOW() - INTERVAL '24 hours' AND
        #   engagement_rate > 0.10
        #   ORDER BY engagement_rate DESC, views DESC
        #   LIMIT limit
        
        # Mock implementation
        viral_items = await self._mock_get_content(
            filter_type="viral", limit=limit
        )
        
        # Cache for 5 minutes
        await self.redis.setex(cache_key, 300, self._serialize_features(viral_items))
        
        return viral_items
    
    async def _get_trending_by_interests(
        self, topics: List[str], hashtags: List[str], limit: int
    ) -> List[ContentFeatures]:
        """
        Get trending content matching user interests
        """
        if not topics and not hashtags:
            return []
        
        # Use Elasticsearch for fast matching (if available)
        if self.es:
            # Query ES for content with matching topics/hashtags
            # Sort by engagement, recency
            results = await self._es_search_by_interests(topics, hashtags, limit)
            return results
        
        # Fallback to DB
        trending_items = await self._mock_get_content(
            filter_type="trending", topics=topics, hashtags=hashtags, limit=limit
        )
        
        return trending_items
    
    async def _get_similar_content(
        self, user_id: str, limit: int
    ) -> List[ContentFeatures]:
        """
        Get content similar to user's liked content
        Collaborative filtering
        """
        # In production: Use collaborative filtering
        # - Get user's liked content
        # - Find users with similar likes
        # - Get their liked content
        
        similar_items = await self._mock_get_content(
            filter_type="similar", user_id=user_id, limit=limit
        )
        
        return similar_items
    
    async def _get_fresh_content(
        self, hours: int, limit: int
    ) -> List[ContentFeatures]:
        """
        Get fresh content (recent)
        """
        cache_key = f"fresh_content:{hours}:{limit}"
        cached = await self.redis.get(cache_key)
        if cached:
            return self._deserialize_features(cached)
        
        # Query recent content
        fresh_items = await self._mock_get_content(
            filter_type="fresh", hours=hours, limit=limit
        )
        
        await self.redis.setex(cache_key, 120, self._serialize_features(fresh_items))
        
        return fresh_items
    
    async def _get_content_by_users(
        self, user_ids: List[str], limit: int
    ) -> List[ContentFeatures]:
        """
        Get content from specific users (following, friends)
        """
        if not user_ids:
            return []
        
        # Query DB for content by these users
        # ORDER BY created_at DESC
        content = await self._mock_get_content(
            filter_type="by_users", user_ids=user_ids, limit=limit
        )
        
        return content
    
    async def _get_nearby_content(
        self, latitude: float, longitude: float, radius_km: float, limit: int
    ) -> List[ContentFeatures]:
        """
        Get content near a location
        """
        # Use PostGIS or geo-indexing
        # SELECT * FROM content WHERE 
        #   ST_Distance_Sphere(
        #     point(longitude, latitude),
        #     point(content.longitude, content.latitude)
        #   ) / 1000 < radius_km
        #   ORDER BY created_at DESC
        #   LIMIT limit
        
        nearby_items = await self._mock_get_content(
            filter_type="nearby", lat=latitude, lon=longitude,
            radius=radius_km, limit=limit
        )
        
        return nearby_items
    
    async def _get_exploration_content(
        self, user_profile: UserProfile, limit: int
    ) -> List[ContentFeatures]:
        """
        Get content for exploration (adjacent to interests)
        """
        # Get content with partial topic overlap
        # Or use "More Like This" query in Elasticsearch
        
        exploration_items = await self._mock_get_content(
            filter_type="exploration", user_profile=user_profile, limit=limit
        )
        
        return exploration_items
    
    async def get_full_content(
        self, content_id: str, content_type: ContentType
    ) -> Optional[Dict[str, Any]]:
        """
        Get full content data (caption, media URLs, etc.)
        """
        # Check cache
        cache_key = f"content:{content_id}"
        cached = await self.redis.get(cache_key)
        if cached:
            return self._deserialize_dict(cached)
        
        # Query from appropriate service
        # - Posts/Takes: Post service API
        # - Stories: Story service API
        
        # Mock implementation
        full_content = {
            "content_id": content_id,
            "username": "mock_user",
            "avatar_url": "https://example.com/avatar.jpg",
            "caption": "Mock content",
            "media_urls": ["https://example.com/media.mp4"],
            "thumbnail_url": "https://example.com/thumb.jpg",
            "duration_seconds": 30
        }
        
        # Cache for 1 hour
        await self.redis.setex(cache_key, 3600, self._serialize_dict(full_content))
        
        return full_content
    
    async def has_user_liked(self, user_id: str, content_id: str) -> bool:
        """Check if user has liked content"""
        # Check in Redis set or DB
        key = f"user_likes:{user_id}"
        return await self.redis.sismember(key, content_id)
    
    async def has_user_saved(self, user_id: str, content_id: str) -> bool:
        """Check if user has saved content"""
        key = f"user_saves:{user_id}"
        return await self.redis.sismember(key, content_id)
    
    async def update_content_metrics(self, signal: UserSignal):
        """
        Update content engagement metrics based on signal
        """
        content_key = f"content_metrics:{signal.content_id}"
        
        # Increment appropriate counter
        if signal.signal_type == "like":
            await self.redis.hincrby(content_key, "likes", 1)
        elif signal.signal_type == "comment":
            await self.redis.hincrby(content_key, "comments", 1)
        elif signal.signal_type == "share":
            await self.redis.hincrby(content_key, "shares", 1)
        elif signal.signal_type == "save":
            await self.redis.hincrby(content_key, "saves", 1)
        elif signal.signal_type == "view":
            await self.redis.hincrby(content_key, "views", 1)
        
        # Recalculate engagement rate
        await self._recalculate_engagement_rate(signal.content_id)
    
    async def _recalculate_engagement_rate(self, content_id: str):
        """Recalculate and cache engagement rate"""
        metrics_key = f"content_metrics:{content_id}"
        metrics = await self.redis.hgetall(metrics_key)
        
        if not metrics:
            return
        
        likes = int(metrics.get("likes", 0))
        comments = int(metrics.get("comments", 0))
        shares = int(metrics.get("shares", 0))
        views = int(metrics.get("views", 1))
        
        engagement_rate = (likes + comments + shares) / max(views, 1)
        
        await self.redis.hset(metrics_key, "engagement_rate", engagement_rate)
    
    # Mock implementations (replace with real DB queries)
    
    async def _mock_get_content(self, filter_type: str, **kwargs) -> List[ContentFeatures]:
        """
        Mock content retrieval
        In production, replace with real DB queries
        """
        limit = kwargs.get("limit", 10)
        
        # Generate mock content features
        mock_items = []
        for i in range(min(limit, 20)):
            mock_items.append(ContentFeatures(
                content_id=f"content_{filter_type}_{i}",
                content_type=ContentType.TAKE if i % 2 == 0 else ContentType.POST,
                creator_id=f"user_{i % 10}",
                likes=100 + i * 10,
                comments=20 + i * 2,
                shares=10 + i,
                saves=5 + i,
                views=1000 + i * 50,
                engagement_rate=0.13 + i * 0.01,
                viral_score=0.5 + i * 0.03,
                quality_score=0.6 + i * 0.02,
                hashtags=[f"#{filter_type}", "#vignette"],
                topics=[f"topic_{i % 5}"],
                created_at=datetime.utcnow() - timedelta(hours=i),
                recency_hours=float(i),
                creator_follower_count=5000 + i * 100
            ))
        
        return mock_items
    
    async def _es_search_by_interests(
        self, topics: List[str], hashtags: List[str], limit: int
    ) -> List[ContentFeatures]:
        """Elasticsearch search by interests"""
        # In production: Use ES multi_match query
        return await self._mock_get_content("interest", limit=limit)
    
    def _serialize_features(self, features: List[ContentFeatures]) -> str:
        """Serialize features for caching"""
        import json
        return json.dumps([f.model_dump() for f in features])
    
    def _deserialize_features(self, data: str) -> List[ContentFeatures]:
        """Deserialize features from cache"""
        import json
        items = json.loads(data)
        return [ContentFeatures(**item) for item in items]
    
    def _serialize_dict(self, data: Dict) -> str:
        """Serialize dict for caching"""
        import json
        return json.dumps(data)
    
    def _deserialize_dict(self, data: str) -> Dict:
        """Deserialize dict from cache"""
        import json
        return json.loads(data)
