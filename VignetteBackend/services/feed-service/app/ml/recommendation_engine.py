"""
Recommendation Engine for Vignette Feed Service
Implements 3 sophisticated feed algorithms:
1. Home Feed (For You) - TikTok competitor
2. Circle Feed - Friends + nearby
3. Surprise & Delight - Balanced exploration
"""
import numpy as np
from typing import List, Dict, Tuple, Optional, Set
from datetime import datetime, timedelta
from collections import defaultdict
import math

from app.schemas.feed import (
    FeedType, ContentFeatures, UserProfile, FeedItemSource,
    RankingConfig, ContentType
)


class RecommendationEngine:
    """
    Core recommendation engine with multiple feed algorithms
    """
    
    def __init__(self):
        self.config = RankingConfig()
    
    # ============================================
    # HOME FEED (FOR YOU) - TikTok Competitor
    # ============================================
    
    def rank_home_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str]
    ) -> List[Tuple[ContentFeatures, float, FeedItemSource]]:
        """
        Home Feed Algorithm - Discovery-focused like TikTok's For You
        
        Prioritizes:
        - Viral content (high engagement rate)
        - User interests (personalization)
        - Trending topics
        - Content diversity
        - Fresh content
        """
        scored_items = []
        
        for content in candidates:
            if content.content_id in seen_ids:
                continue
            
            # Calculate component scores
            viral_score = self._calculate_viral_score(content)
            interest_score = self._calculate_interest_score(content, user_profile)
            quality_score = self._calculate_quality_score(content)
            recency_score = self._calculate_recency_score(content)
            diversity_bonus = self._calculate_diversity_bonus(content, scored_items)
            
            # Weighted combination
            final_score = (
                0.35 * viral_score +          # Viral potential
                0.30 * interest_score +        # Personalization
                0.15 * quality_score +         # Content quality
                0.15 * recency_score +         # Freshness
                0.05 * diversity_bonus         # Diversity
            )
            
            # Determine source
            source = self._determine_home_source(
                viral_score, interest_score, content
            )
            
            scored_items.append((content, final_score, source))
        
        # Sort by score descending
        scored_items.sort(key=lambda x: x[1], reverse=True)
        
        # Apply diversity constraints
        scored_items = self._apply_diversity_constraints(scored_items)
        
        return scored_items
    
    def _calculate_viral_score(self, content: ContentFeatures) -> float:
        """
        Calculate viral potential score
        Based on engagement rate, growth rate, and velocity
        """
        if content.views == 0:
            return 0.0
        
        # Engagement rate (likes + comments + shares) / views
        engagement_rate = (
            content.likes + content.comments + content.shares
        ) / max(content.views, 1)
        
        # Normalize to 0-1 (assume 15% is excellent)
        normalized_engagement = min(engagement_rate / 0.15, 1.0)
        
        # Viral velocity (exponential growth)
        viral_velocity = min(content.viral_score, 1.0)
        
        # Combine
        return (0.6 * normalized_engagement + 0.4 * viral_velocity)
    
    def _calculate_interest_score(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Calculate personalization score based on user interests
        """
        score = 0.0
        matches = 0
        
        # Hashtag matching
        user_hashtags = set(user_profile.interest_hashtags)
        content_hashtags = set(content.hashtags)
        if user_hashtags and content_hashtags:
            overlap = len(user_hashtags & content_hashtags)
            score += overlap / len(user_hashtags)
            matches += 1
        
        # Topic matching
        user_topics = set(user_profile.interest_topics)
        content_topics = set(content.topics)
        if user_topics and content_topics:
            overlap = len(user_topics & content_topics)
            score += overlap / len(user_topics)
            matches += 1
        
        # Creator matching
        if content.creator_id in user_profile.interest_creators:
            score += 1.0
            matches += 1
        
        # Sound matching
        if content.sound_id and content.sound_id in user_profile.interest_sounds:
            score += 0.8
            matches += 1
        
        # Average if any matches
        return score / max(matches, 1)
    
    def _calculate_quality_score(self, content: ContentFeatures) -> float:
        """
        Calculate content quality score
        """
        # Quality indicators
        save_rate = content.saves / max(content.views, 1)
        comment_rate = content.comments / max(content.views, 1)
        
        # Normalize (assume 5% save rate is excellent, 2% comment rate)
        normalized_save = min(save_rate / 0.05, 1.0)
        normalized_comment = min(comment_rate / 0.02, 1.0)
        
        # Creator authority
        creator_score = min(content.creator_follower_count / 10000, 1.0)
        
        # Combine
        return (
            0.4 * normalized_save +
            0.3 * normalized_comment +
            0.3 * creator_score
        )
    
    def _calculate_recency_score(self, content: ContentFeatures) -> float:
        """
        Calculate recency score (favor fresh content)
        """
        hours_old = content.recency_hours
        
        # Exponential decay: score = e^(-hours/24)
        # Fresh content (0-6 hours): ~0.8-1.0
        # Recent content (6-24 hours): ~0.5-0.8
        # Older content (24+ hours): <0.5
        decay_factor = math.exp(-hours_old / 24)
        
        return decay_factor
    
    def _calculate_diversity_bonus(
        self, content: ContentFeatures, current_items: List
    ) -> float:
        """
        Bonus for diversity (new topics, creators)
        """
        if not current_items:
            return 1.0
        
        # Check creator diversity
        recent_creators = [item[0].creator_id for item in current_items[-5:]]
        if content.creator_id in recent_creators:
            return 0.3  # Penalty for same creator
        
        # Check topic diversity
        recent_topics = set()
        for item in current_items[-3:]:
            recent_topics.update(item[0].topics)
        
        content_topics = set(content.topics)
        if content_topics & recent_topics:
            return 0.7  # Small penalty for same topic
        
        return 1.0  # Full bonus for diversity
    
    def _determine_home_source(
        self, viral_score: float, interest_score: float, content: ContentFeatures
    ) -> FeedItemSource:
        """
        Determine source label for Home feed item
        """
        if viral_score > 0.7:
            return FeedItemSource.VIRAL
        elif interest_score > 0.6:
            return FeedItemSource.INTEREST
        elif content.recency_hours < 6:
            return FeedItemSource.TRENDING
        else:
            return FeedItemSource.DISCOVERY
    
    # ============================================
    # CIRCLE FEED - Friends + Nearby
    # ============================================
    
    def rank_circle_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str],
        user_lat: Optional[float] = None,
        user_lon: Optional[float] = None
    ) -> List[Tuple[ContentFeatures, float, FeedItemSource, Optional[float]]]:
        """
        Circle Feed Algorithm - Social graph + location
        
        Prioritizes:
        - Friends (close connections)
        - Following
        - Nearby content (geo-based)
        - Mutual connections
        """
        scored_items = []
        
        for content in candidates:
            if content.content_id in seen_ids:
                continue
            
            # Calculate component scores
            social_score = self._calculate_social_score(content, user_profile)
            proximity_score = self._calculate_proximity_score(
                content, user_lat, user_lon
            )
            engagement_score = self._calculate_engagement_score(content)
            recency_score = self._calculate_recency_score(content)
            
            # Weighted combination (social and proximity dominant)
            final_score = (
                0.50 * social_score +          # Social graph
                0.25 * proximity_score +       # Location
                0.15 * engagement_score +      # Engagement
                0.10 * recency_score           # Freshness
            )
            
            # Determine source
            source = self._determine_circle_source(
                social_score, proximity_score, content, user_profile
            )
            
            # Calculate distance (for display)
            distance_km = self._calculate_distance(
                user_lat, user_lon, content.location_lat, content.location_lon
            )
            
            scored_items.append((content, final_score, source, distance_km))
        
        # Sort by score descending
        scored_items.sort(key=lambda x: x[1], reverse=True)
        
        return scored_items
    
    def _calculate_social_score(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Calculate social graph score
        """
        creator_id = content.creator_id
        
        # Close friends (highest priority)
        if creator_id in user_profile.friend_ids:
            return 1.0
        
        # Following (high priority)
        if creator_id in user_profile.following_ids:
            return 0.8
        
        # Mutual connections (medium priority)
        # (In production, check if creator follows user back)
        if creator_id in user_profile.follower_ids:
            return 0.6
        
        # No connection
        return 0.1
    
    def _calculate_proximity_score(
        self,
        content: ContentFeatures,
        user_lat: Optional[float],
        user_lon: Optional[float]
    ) -> float:
        """
        Calculate location proximity score
        """
        if not all([user_lat, user_lon, content.location_lat, content.location_lon]):
            return 0.0
        
        distance_km = self._calculate_distance(
            user_lat, user_lon, content.location_lat, content.location_lon
        )
        
        # Scoring:
        # < 1km: 1.0
        # 1-5km: 0.8
        # 5-25km: 0.5
        # 25-100km: 0.2
        # > 100km: 0.0
        if distance_km < 1:
            return 1.0
        elif distance_km < 5:
            return 0.8
        elif distance_km < 25:
            return 0.5
        elif distance_km < 100:
            return 0.2
        else:
            return 0.0
    
    def _calculate_distance(
        self,
        lat1: Optional[float],
        lon1: Optional[float],
        lat2: Optional[float],
        lon2: Optional[float]
    ) -> Optional[float]:
        """
        Calculate haversine distance in kilometers
        """
        if not all([lat1, lon1, lat2, lon2]):
            return None
        
        # Haversine formula
        R = 6371  # Earth radius in km
        
        lat1_rad = math.radians(lat1)
        lat2_rad = math.radians(lat2)
        delta_lat = math.radians(lat2 - lat1)
        delta_lon = math.radians(lon2 - lon1)
        
        a = (
            math.sin(delta_lat / 2) ** 2 +
            math.cos(lat1_rad) * math.cos(lat2_rad) *
            math.sin(delta_lon / 2) ** 2
        )
        c = 2 * math.atan2(math.sqrt(a), math.sqrt(1 - a))
        
        return R * c
    
    def _calculate_engagement_score(self, content: ContentFeatures) -> float:
        """
        Simple engagement score for Circle feed
        """
        if content.views == 0:
            return 0.0
        
        engagement_rate = (
            content.likes + content.comments + content.shares
        ) / max(content.views, 1)
        
        return min(engagement_rate / 0.10, 1.0)
    
    def _determine_circle_source(
        self,
        social_score: float,
        proximity_score: float,
        content: ContentFeatures,
        user_profile: UserProfile
    ) -> FeedItemSource:
        """
        Determine source label for Circle feed item
        """
        creator_id = content.creator_id
        
        if creator_id in user_profile.friend_ids:
            return FeedItemSource.FRIENDS
        elif creator_id in user_profile.following_ids:
            return FeedItemSource.FOLLOWING
        elif proximity_score > 0.5:
            return FeedItemSource.NEARBY
        else:
            return FeedItemSource.MUTUAL
    
    # ============================================
    # SURPRISE & DELIGHT FEED - Balanced Exploration
    # ============================================
    
    def rank_surprise_delight_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str]
    ) -> List[Tuple[ContentFeatures, float, FeedItemSource]]:
        """
        Surprise & Delight Algorithm - Balanced exploration
        
        Mix:
        - 60% Known interests (personalized)
        - 30% Exploration (outside interests)
        - 10% Surprises (random viral)
        
        Less chaotic, more curated than Home feed
        """
        # Categorize candidates
        known_interest = []
        exploration = []
        surprise = []
        
        for content in candidates:
            if content.content_id in seen_ids:
                continue
            
            interest_score = self._calculate_interest_score(content, user_profile)
            viral_score = self._calculate_viral_score(content)
            
            # Categorize
            if interest_score > 0.6:
                # Known interest
                score = self._score_known_interest(content, user_profile)
                known_interest.append((content, score, FeedItemSource.KNOWN_INTEREST))
            
            elif interest_score > 0.2 or self._is_adjacent_interest(content, user_profile):
                # Exploration (adjacent topics)
                score = self._score_exploration(content, user_profile)
                exploration.append((content, score, FeedItemSource.EXPLORATION))
            
            elif viral_score > 0.7:
                # Surprise (viral content)
                score = self._score_surprise(content)
                surprise.append((content, score, FeedItemSource.SURPRISE))
        
        # Sort each category
        known_interest.sort(key=lambda x: x[1], reverse=True)
        exploration.sort(key=lambda x: x[1], reverse=True)
        surprise.sort(key=lambda x: x[1], reverse=True)
        
        # Mix according to ratios: 60% / 30% / 10%
        total_items = 20  # Target feed size
        n_known = int(total_items * 0.6)  # 12 items
        n_explore = int(total_items * 0.3)  # 6 items
        n_surprise = int(total_items * 0.1)  # 2 items
        
        # Interleave items for smooth experience
        mixed_feed = self._interleave_feed(
            known_interest[:n_known],
            exploration[:n_explore],
            surprise[:n_surprise]
        )
        
        return mixed_feed
    
    def _score_known_interest(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score for known interest items (curated, high quality)
        """
        interest_score = self._calculate_interest_score(content, user_profile)
        quality_score = self._calculate_quality_score(content)
        recency_score = self._calculate_recency_score(content)
        
        # Emphasize quality and personalization
        return (
            0.50 * interest_score +
            0.35 * quality_score +
            0.15 * recency_score
        )
    
    def _score_exploration(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score for exploration items (adjacent interests)
        """
        adjacency_score = self._calculate_adjacency_score(content, user_profile)
        quality_score = self._calculate_quality_score(content)
        viral_score = self._calculate_viral_score(content)
        
        # Balance quality and discovery
        return (
            0.35 * adjacency_score +
            0.35 * quality_score +
            0.30 * viral_score
        )
    
    def _score_surprise(self, content: ContentFeatures) -> float:
        """
        Score for surprise items (viral, random)
        """
        viral_score = self._calculate_viral_score(content)
        quality_score = self._calculate_quality_score(content)
        
        # Emphasize viral and quality
        return 0.6 * viral_score + 0.4 * quality_score
    
    def _is_adjacent_interest(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> bool:
        """
        Check if content is adjacent to user interests
        (Related topics, not exact match)
        """
        # In production, use topic embeddings or knowledge graph
        # For now, simple heuristic: share at least one hashtag or topic
        user_tags = set(user_profile.interest_hashtags) | set(user_profile.interest_topics)
        content_tags = set(content.hashtags) | set(content.topics)
        
        return len(user_tags & content_tags) > 0
    
    def _calculate_adjacency_score(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Calculate how adjacent content is to user interests
        """
        # Partial overlap with interests (not full match)
        interest_score = self._calculate_interest_score(content, user_profile)
        
        # Ideal adjacency: 0.2 - 0.6 (some overlap, not too much)
        if 0.2 <= interest_score <= 0.6:
            return 1.0
        elif interest_score < 0.2:
            return interest_score / 0.2  # Too far
        else:
            return (1.0 - interest_score) / 0.4  # Too close
    
    def _interleave_feed(
        self,
        known: List[Tuple],
        explore: List[Tuple],
        surprise: List[Tuple]
    ) -> List[Tuple]:
        """
        Interleave three categories for smooth experience
        Pattern: K K K E K K E S K K E ...
        """
        result = []
        k_idx = e_idx = s_idx = 0
        
        # Pattern: 3 known, 1 explore, repeat, occasional surprise
        while k_idx < len(known) or e_idx < len(explore) or s_idx < len(surprise):
            # Add 3 known interest items
            for _ in range(3):
                if k_idx < len(known):
                    result.append(known[k_idx])
                    k_idx += 1
            
            # Add 1 exploration item
            if e_idx < len(explore):
                result.append(explore[e_idx])
                e_idx += 1
            
            # Every 10 items, add a surprise
            if len(result) % 10 == 0 and s_idx < len(surprise):
                result.append(surprise[s_idx])
                s_idx += 1
        
        return result
    
    # ============================================
    # HELPER METHODS
    # ============================================
    
    def _apply_diversity_constraints(
        self, scored_items: List[Tuple]
    ) -> List[Tuple]:
        """
        Apply diversity constraints to prevent fatigue
        """
        result = []
        creator_count = defaultdict(int)
        topic_count = defaultdict(int)
        
        for item in scored_items:
            content = item[0]
            
            # Check creator fatigue
            if creator_count[content.creator_id] >= self.config.creator_fatigue_limit:
                continue
            
            # Check topic fatigue (consecutive items)
            if result:
                recent_topics = set()
                for prev_item in result[-3:]:
                    recent_topics.update(prev_item[0].topics)
                
                content_topics = set(content.topics)
                if len(content_topics & recent_topics) >= 2:
                    # Too similar to recent items
                    continue
            
            result.append(item)
            creator_count[content.creator_id] += 1
            for topic in content.topics:
                topic_count[topic] += 1
        
        return result
