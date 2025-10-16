"""
Feed Service for Vignette
Orchestrates 3 feed types: Home, Circle, Surprise & Delight
"""
from typing import List, Optional, Set
from datetime import datetime, timedelta
import logging

from app.schemas.feed import (
    FeedType, FeedRequest, FeedResponse, FeedItem,
    ContentFeatures, UserProfile, FeedItemSource,
    ContentType, UserSignal
)
from app.ml.recommendation_engine import RecommendationEngine
from app.services.ranking_service import RankingService
from app.services.personalization_service import PersonalizationService


logger = logging.getLogger(__name__)


class FeedService:
    """
    Main feed service - generates personalized feeds
    """
    
    def __init__(
        self,
        recommendation_engine: RecommendationEngine,
        ranking_service: RankingService,
        personalization_service: PersonalizationService
    ):
        self.rec_engine = recommendation_engine
        self.ranking_service = ranking_service
        self.personalization_service = personalization_service
    
    async def generate_feed(
        self, request: FeedRequest
    ) -> FeedResponse:
        """
        Generate personalized feed based on feed type
        """
        logger.info(
            f"Generating {request.feed_type} feed for user {request.user_id}"
        )
        
        # Get user profile
        user_profile = await self.personalization_service.get_user_profile(
            request.user_id
        )
        
        # Get candidate content
        candidates = await self._get_candidates(request, user_profile)
        
        # Get seen content
        seen_ids = set(request.seen_ids)
        
        # Rank based on feed type
        if request.feed_type == FeedType.HOME:
            ranked_items = await self._generate_home_feed(
                user_profile, candidates, seen_ids
            )
        
        elif request.feed_type == FeedType.CIRCLE:
            ranked_items = await self._generate_circle_feed(
                user_profile, candidates, seen_ids,
                request.latitude, request.longitude
            )
        
        elif request.feed_type == FeedType.SURPRISE_DELIGHT:
            ranked_items = await self._generate_surprise_delight_feed(
                user_profile, candidates, seen_ids
            )
        
        else:
            raise ValueError(f"Unknown feed type: {request.feed_type}")
        
        # Convert to feed items
        feed_items = await self._convert_to_feed_items(
            ranked_items, request.user_id, request.offset
        )
        
        # Paginate
        start_idx = request.offset
        end_idx = start_idx + request.limit
        paginated_items = feed_items[start_idx:end_idx]
        
        # Build response
        response = FeedResponse(
            feed_type=request.feed_type,
            items=paginated_items,
            next_offset=end_idx,
            has_more=end_idx < len(feed_items)
        )
        
        logger.info(
            f"Generated {len(paginated_items)} items for {request.feed_type} feed"
        )
        
        return response
    
    async def _get_candidates(
        self, request: FeedRequest, user_profile: UserProfile
    ) -> List[ContentFeatures]:
        """
        Get candidate content for ranking
        Uses different strategies per feed type
        """
        if request.feed_type == FeedType.HOME:
            # Home: Mix of viral, trending, personalized
            candidates = await self.ranking_service.get_home_candidates(
                user_profile, limit=500
            )
        
        elif request.feed_type == FeedType.CIRCLE:
            # Circle: Social graph + location
            candidates = await self.ranking_service.get_circle_candidates(
                user_profile,
                latitude=request.latitude,
                longitude=request.longitude,
                limit=300
            )
        
        elif request.feed_type == FeedType.SURPRISE_DELIGHT:
            # Surprise: Diverse mix
            candidates = await self.ranking_service.get_surprise_candidates(
                user_profile, limit=400
            )
        
        else:
            candidates = []
        
        return candidates
    
    async def _generate_home_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str]
    ) -> List:
        """
        Generate Home Feed (For You)
        """
        logger.info(f"Ranking {len(candidates)} candidates for Home feed")
        
        ranked_items = self.rec_engine.rank_home_feed(
            user_profile, candidates, seen_ids
        )
        
        logger.info(f"Ranked {len(ranked_items)} items for Home feed")
        return ranked_items
    
    async def _generate_circle_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str],
        latitude: Optional[float],
        longitude: Optional[float]
    ) -> List:
        """
        Generate Circle Feed (Friends + Nearby)
        """
        logger.info(f"Ranking {len(candidates)} candidates for Circle feed")
        
        ranked_items = self.rec_engine.rank_circle_feed(
            user_profile, candidates, seen_ids, latitude, longitude
        )
        
        logger.info(f"Ranked {len(ranked_items)} items for Circle feed")
        return ranked_items
    
    async def _generate_surprise_delight_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str]
    ) -> List:
        """
        Generate Surprise & Delight Feed
        """
        logger.info(f"Ranking {len(candidates)} candidates for Surprise & Delight feed")
        
        ranked_items = self.rec_engine.rank_surprise_delight_feed(
            user_profile, candidates, seen_ids
        )
        
        logger.info(f"Ranked {len(ranked_items)} items for Surprise & Delight feed")
        return ranked_items
    
    async def _convert_to_feed_items(
        self,
        ranked_items: List,
        user_id: str,
        offset: int
    ) -> List[FeedItem]:
        """
        Convert ranked content to feed items with full data
        """
        feed_items = []
        
        for idx, item in enumerate(ranked_items):
            # Unpack based on feed type
            if len(item) == 3:
                # Home, Surprise & Delight: (content, score, source)
                content, score, source = item
                distance_km = None
            else:
                # Circle: (content, score, source, distance)
                content, score, source, distance_km = item
            
            # Get full content data (from DB/cache)
            full_content = await self.ranking_service.get_full_content(
                content.content_id, content.content_type
            )
            
            if not full_content:
                continue
            
            # Check user interactions
            is_following = content.creator_id in (
                await self.personalization_service.get_user_following(user_id)
            )
            has_liked = await self.ranking_service.has_user_liked(
                user_id, content.content_id
            )
            has_saved = await self.ranking_service.has_user_saved(
                user_id, content.content_id
            )
            
            # Build feed item
            feed_item = FeedItem(
                content_id=content.content_id,
                content_type=content.content_type,
                user_id=content.creator_id,
                username=full_content.get("username", ""),
                avatar_url=full_content.get("avatar_url"),
                caption=full_content.get("caption"),
                media_urls=full_content.get("media_urls", []),
                thumbnail_url=full_content.get("thumbnail_url"),
                duration_seconds=full_content.get("duration_seconds"),
                likes_count=content.likes,
                comments_count=content.comments,
                shares_count=content.shares,
                saves_count=content.saves,
                views_count=content.views,
                score=score,
                source=source,
                rank=offset + idx,
                created_at=content.created_at,
                is_following=is_following,
                has_liked=has_liked,
                has_saved=has_saved,
                distance_km=distance_km
            )
            
            feed_items.append(feed_item)
        
        return feed_items
    
    async def track_signal(self, signal: UserSignal):
        """
        Track user engagement signal for learning
        """
        logger.info(
            f"Tracking signal: {signal.signal_type} for user {signal.user_id} "
            f"on content {signal.content_id}"
        )
        
        # Update personalization profile
        await self.personalization_service.process_signal(signal)
        
        # Update content features (engagement metrics)
        await self.ranking_service.update_content_metrics(signal)
        
        logger.info(f"Signal tracked successfully")
    
    async def get_feed_metrics(
        self, user_id: str, feed_type: FeedType, hours: int = 24
    ):
        """
        Get feed performance metrics for a user
        """
        since = datetime.utcnow() - timedelta(hours=hours)
        
        metrics = await self.personalization_service.get_feed_metrics(
            user_id, feed_type, since
        )
        
        return metrics
