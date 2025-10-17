"""
Recommendation Engine for Socialink Feed Service
Single unified feed algorithm with friends priority
Mix: 70% friends, 20% interests, 10% discovery
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
    Socialink recommendation engine
    Single unified feed that learns from behavior and prioritizes friends
    """
    
    def __init__(self):
        self.config = RankingConfig()
    
    # ============================================
    # SOCIALINK UNIFIED FEED
    # Friends-first with intelligent discovery
    # ============================================
    
    def rank_home_feed(
        self,
        user_profile: UserProfile,
        candidates: List[ContentFeatures],
        seen_ids: Set[str]
    ) -> List[Tuple[ContentFeatures, float, FeedItemSource]]:
        """
        Socialink Unified Feed Algorithm
        
        Prioritization:
        - 70% Friends/Following (social graph priority)
        - 20% Interests (pages, groups, topics user likes)
        - 10% Discovery (suggested new content)
        
        Learns from user behavior:
        - Engagement signals (likes, comments, shares)
        - Time spent on content
        - Content types preferred
        - Friends vs pages balance
        """
        # Categorize candidates by source
        friends_content = []
        following_content = []
        groups_content = []
        pages_content = []
        interest_content = []
        discovery_content = []
        
        for content in candidates:
            if content.content_id in seen_ids:
                continue
            
            # Categorize by creator relationship
            source = self._determine_content_source(content, user_profile)
            
            if source == FeedItemSource.FRIENDS:
                score = self._score_friends_content(content, user_profile)
                friends_content.append((content, score, source))
            
            elif source == FeedItemSource.FOLLOWING:
                score = self._score_following_content(content, user_profile)
                following_content.append((content, score, source))
            
            elif source == FeedItemSource.GROUPS:
                score = self._score_groups_content(content, user_profile)
                groups_content.append((content, score, source))
            
            elif source == FeedItemSource.PAGES:
                score = self._score_pages_content(content, user_profile)
                pages_content.append((content, score, source))
            
            elif source == FeedItemSource.INTEREST:
                score = self._score_interest_content(content, user_profile)
                interest_content.append((content, score, source))
            
            else:  # Discovery
                score = self._score_discovery_content(content, user_profile)
                discovery_content.append((content, score, source))
        
        # Sort each category
        friends_content.sort(key=lambda x: x[1], reverse=True)
        following_content.sort(key=lambda x: x[1], reverse=True)
        groups_content.sort(key=lambda x: x[1], reverse=True)
        pages_content.sort(key=lambda x: x[1], reverse=True)
        interest_content.sort(key=lambda x: x[1], reverse=True)
        discovery_content.sort(key=lambda x: x[1], reverse=True)
        
        # Mix according to ratios
        # 70% social (friends + following + groups)
        # 20% interests (pages + topics)
        # 10% discovery
        total_items = 30  # Target feed size
        n_social = int(total_items * 0.70)  # 21 items
        n_interest = int(total_items * 0.20)  # 6 items
        n_discovery = int(total_items * 0.10)  # 3 items
        
        # Within social (70%):
        # - Friends get 60% (highest priority)
        # - Following get 25%
        # - Groups get 15%
        n_friends = int(n_social * 0.60)  # ~13 items
        n_following = int(n_social * 0.25)  # ~5 items
        n_groups = int(n_social * 0.15)  # ~3 items
        
        # Within interests (20%):
        # - Pages: 60%
        # - Topics: 40%
        n_pages = int(n_interest * 0.60)  # ~4 items
        n_topics = int(n_interest * 0.40)  # ~2 items
        
        # Interleave for natural flow
        mixed_feed = self._interleave_socialink_feed(
            friends_content[:n_friends],
            following_content[:n_following],
            groups_content[:n_groups],
            pages_content[:n_pages],
            interest_content[:n_topics],
            discovery_content[:n_discovery]
        )
        
        return mixed_feed
    
    def _determine_content_source(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> FeedItemSource:
        """
        Determine source category for content
        """
        creator_id = content.creator_id
        
        # Close friends (mutual + high interaction)
        if creator_id in user_profile.friend_ids:
            return FeedItemSource.FRIENDS
        
        # Following (users/creators)
        if creator_id in user_profile.following_ids:
            return FeedItemSource.FOLLOWING
        
        # Groups (in production: check if content is from user's groups)
        # For now: Check if mentioned in topics
        if any("group" in topic.lower() for topic in content.topics):
            return FeedItemSource.GROUPS
        
        # Pages (in production: check if from pages user liked)
        # For now: Check creator follower count (pages have many followers)
        if content.creator_follower_count > 10000:
            return FeedItemSource.PAGES
        
        # Interest match
        interest_score = self._calculate_interest_score(content, user_profile)
        if interest_score > 0.4:
            return FeedItemSource.INTEREST
        
        # Discovery (new content)
        return FeedItemSource.SUGGESTED
    
    def _score_friends_content(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score friends' content (highest priority)
        Emphasize: Recent, engagement, meaningful interactions
        """
        # Recency (friends' content should be fresh)
        recency_score = self._calculate_recency_score(content)
        
        # Engagement (quality filter)
        engagement_score = self._calculate_engagement_score(content)
        
        # Interaction history (in production: use past interactions with this friend)
        interaction_score = 1.0  # Default high
        
        # Meaningful content (long captions, photos = more meaningful than quick posts)
        meaningful_score = self._calculate_meaningful_score(content)
        
        # Weighted combination (emphasize recency and meaningfulness)
        return (
            0.40 * recency_score +
            0.25 * engagement_score +
            0.20 * interaction_score +
            0.15 * meaningful_score
        )
    
    def _score_following_content(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score content from followed users/creators
        """
        recency_score = self._calculate_recency_score(content)
        engagement_score = self._calculate_engagement_score(content)
        interest_score = self._calculate_interest_score(content, user_profile)
        
        return (
            0.40 * recency_score +
            0.35 * engagement_score +
            0.25 * interest_score
        )
    
    def _score_groups_content(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score content from groups
        """
        recency_score = self._calculate_recency_score(content)
        engagement_score = self._calculate_engagement_score(content)
        relevance_score = self._calculate_interest_score(content, user_profile)
        
        return (
            0.35 * recency_score +
            0.35 * engagement_score +
            0.30 * relevance_score
        )
    
    def _score_pages_content(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score content from pages user liked
        """
        quality_score = self._calculate_quality_score(content)
        interest_score = self._calculate_interest_score(content, user_profile)
        engagement_score = self._calculate_engagement_score(content)
        
        return (
            0.40 * quality_score +
            0.35 * interest_score +
            0.25 * engagement_score
        )
    
    def _score_interest_content(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score content matching user interests
        """
        interest_score = self._calculate_interest_score(content, user_profile)
        quality_score = self._calculate_quality_score(content)
        viral_score = self._calculate_viral_score(content)
        
        return (
            0.50 * interest_score +
            0.30 * quality_score +
            0.20 * viral_score
        )
    
    def _score_discovery_content(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Score discovery content (suggested new)
        """
        viral_score = self._calculate_viral_score(content)
        quality_score = self._calculate_quality_score(content)
        novelty_score = self._calculate_novelty_score(content, user_profile)
        
        return (
            0.45 * viral_score +
            0.35 * quality_score +
            0.20 * novelty_score
        )
    
    # ============================================
    # HELPER SCORING FUNCTIONS
    # ============================================
    
    def _calculate_recency_score(self, content: ContentFeatures) -> float:
        """
        Recency score (favor fresh content)
        """
        hours_old = content.recency_hours
        
        # Exponential decay
        # Fresh (0-2h): ~0.9-1.0
        # Recent (2-6h): ~0.7-0.9
        # Older (6-24h): ~0.3-0.7
        # Old (24h+): <0.3
        decay_factor = math.exp(-hours_old / 12)
        
        return decay_factor
    
    def _calculate_engagement_score(self, content: ContentFeatures) -> float:
        """
        Engagement score
        """
        if content.views == 0:
            return 0.0
        
        engagement_rate = (
            content.likes + content.comments + content.shares
        ) / max(content.views, 1)
        
        # Normalize (assume 12% is excellent for friends content)
        return min(engagement_rate / 0.12, 1.0)
    
    def _calculate_interest_score(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Interest matching score
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
        
        # Creator preference
        if content.creator_id in user_profile.interest_creators:
            score += 1.0
            matches += 1
        
        return score / max(matches, 1)
    
    def _calculate_quality_score(self, content: ContentFeatures) -> float:
        """
        Content quality score
        """
        save_rate = content.saves / max(content.views, 1)
        comment_rate = content.comments / max(content.views, 1)
        
        normalized_save = min(save_rate / 0.05, 1.0)
        normalized_comment = min(comment_rate / 0.02, 1.0)
        
        return 0.6 * normalized_save + 0.4 * normalized_comment
    
    def _calculate_viral_score(self, content: ContentFeatures) -> float:
        """
        Viral potential score
        """
        if content.views == 0:
            return 0.0
        
        engagement_rate = (
            content.likes + content.comments + content.shares
        ) / max(content.views, 1)
        
        normalized_engagement = min(engagement_rate / 0.15, 1.0)
        viral_velocity = min(content.viral_score, 1.0)
        
        return 0.6 * normalized_engagement + 0.4 * viral_velocity
    
    def _calculate_meaningful_score(self, content: ContentFeatures) -> float:
        """
        Score for meaningful content (photos with long captions > quick posts)
        """
        # Longer captions = more meaningful
        # In production: Check caption length, media type, etc.
        # For now: Use engagement as proxy
        return self._calculate_engagement_score(content)
    
    def _calculate_novelty_score(
        self, content: ContentFeatures, user_profile: UserProfile
    ) -> float:
        """
        Novelty score (how new/different is this content)
        """
        # Check if topics are different from user's usual
        user_topics = set(user_profile.interest_topics)
        content_topics = set(content.topics)
        
        # Low overlap = high novelty
        if not user_topics:
            return 1.0
        
        overlap = len(user_topics & content_topics)
        novelty = 1.0 - (overlap / len(user_topics))
        
        return max(novelty, 0.0)
    
    def _interleave_socialink_feed(
        self,
        friends: List[Tuple],
        following: List[Tuple],
        groups: List[Tuple],
        pages: List[Tuple],
        interests: List[Tuple],
        discovery: List[Tuple]
    ) -> List[Tuple]:
        """
        Interleave categories for natural Facebook-like feed flow
        Pattern: F F F P F F D F I G F F P ...
        (F=Friend, P=Following, D=Discovery, I=Interest, G=Group)
        """
        result = []
        f_idx = p_idx = g_idx = pg_idx = i_idx = d_idx = 0
        
        # Pattern repeats every 10 items:
        # 3 friends, 1 following, 2 friends, 1 discovery, 1 friend, 1 interest, 1 group
        while (f_idx < len(friends) or p_idx < len(following) or 
               g_idx < len(groups) or pg_idx < len(pages) or
               i_idx < len(interests) or d_idx < len(discovery)):
            
            # 3 friends
            for _ in range(3):
                if f_idx < len(friends):
                    result.append(friends[f_idx])
                    f_idx += 1
            
            # 1 following/page
            if p_idx < len(following):
                result.append(following[p_idx])
                p_idx += 1
            elif pg_idx < len(pages):
                result.append(pages[pg_idx])
                pg_idx += 1
            
            # 2 friends
            for _ in range(2):
                if f_idx < len(friends):
                    result.append(friends[f_idx])
                    f_idx += 1
            
            # 1 discovery
            if d_idx < len(discovery):
                result.append(discovery[d_idx])
                d_idx += 1
            
            # 1 friend
            if f_idx < len(friends):
                result.append(friends[f_idx])
                f_idx += 1
            
            # 1 interest
            if i_idx < len(interests):
                result.append(interests[i_idx])
                i_idx += 1
            
            # 1 group
            if g_idx < len(groups):
                result.append(groups[g_idx])
                g_idx += 1
        
        return result
