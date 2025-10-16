"""
Main Recommendation Service for Vignette
Orchestrates all recommendation algorithms
"""
from typing import List, Dict, Set, Optional
import logging
from collections import defaultdict

from app.models.recommendation import (
    RecommendationType, RecommendationRequest, RecommendationResponse,
    UserRecommendation, CommunityRecommendation,
    RecommendationSource, UserProfile
)
from app.ml.collaborative_filtering import CollaborativeFilter, ItemBasedCollaborativeFilter
from app.ml.graph_based import SocialGraphRecommender


logger = logging.getLogger(__name__)


class RecommendationService:
    """
    Main recommendation service for Vignette
    
    Vignette-specific recommendations:
    - Creators for You: Top creators to follow
    - Suggested for You: General user recommendations
    - Communities for You: Community recommendations
    """
    
    def __init__(
        self,
        db,
        redis_client,
        collaborative_filter: CollaborativeFilter,
        graph_recommender: SocialGraphRecommender
    ):
        self.db = db
        self.redis = redis_client
        self.collaborative = collaborative_filter
        self.graph = graph_recommender
    
    async def get_recommendations(
        self, request: RecommendationRequest
    ) -> RecommendationResponse:
        """
        Get recommendations based on type
        """
        logger.info(f"Generating {request.type} recommendations for user {request.user_id}")
        
        # Get user profile
        user_profile = await self._get_user_profile(request.user_id)
        
        if request.type == RecommendationType.CREATORS_FOR_YOU:
            return await self._recommend_creators(request, user_profile)
        
        elif request.type == RecommendationType.SUGGESTED_FOR_YOU:
            return await self._recommend_users(request, user_profile)
        
        elif request.type == RecommendationType.COMMUNITIES_FOR_YOU:
            return await self._recommend_communities(request, user_profile)
        
        else:
            raise ValueError(f"Unknown recommendation type: {request.type}")
    
    async def _recommend_creators(
        self, request: RecommendationRequest, user_profile: UserProfile
    ) -> RecommendationResponse:
        """
        Creators for You (Vignette)
        
        Algorithm:
        1. Find similar users (collaborative)
        2. Get creators they follow
        3. Add popular creators (high follower count)
        4. Score by engagement potential
        5. Mix and rank
        """
        logger.info(f"Generating creator recommendations for {request.user_id}")
        
        # Get user's following
        following_ids = set(user_profile.following_ids)
        exclude_ids = set(request.exclude_ids) | following_ids | {request.user_id}
        
        # 1. Collaborative filtering (60% weight)
        collab_recs = self.collaborative.recommend_users(
            request.user_id,
            following_ids,
            n_recommendations=40,
            exclude_ids=exclude_ids
        )
        
        # 2. Graph-based (Friends of Friends) (20% weight)
        fof_recs = self.graph.friends_of_friends(
            request.user_id,
            max_distance=2,
            exclude_ids=exclude_ids
        )
        
        # 3. Popular creators (20% weight)
        popular_recs = self.graph.recommend_by_popularity(
            request.user_id,
            following_ids,
            top_k=20,
            exclude_ids=exclude_ids
        )
        
        # Combine and score
        combined_scores = defaultdict(lambda: {
            "score": 0.0,
            "sources": [],
            "mutuals": 0
        })
        
        # Add collaborative recommendations
        for user_id, score in collab_recs:
            combined_scores[user_id]["score"] += score * 0.6
            combined_scores[user_id]["sources"].append(RecommendationSource.COLLABORATIVE)
        
        # Add FoF recommendations
        for user_id, distance, mutuals in fof_recs[:30]:
            score = 1.0 / distance  # Closer = better
            combined_scores[user_id]["score"] += score * 0.2
            combined_scores[user_id]["sources"].append(RecommendationSource.GRAPH_BASED)
            combined_scores[user_id]["mutuals"] = mutuals
        
        # Add popular recommendations
        for user_id, score in popular_recs:
            combined_scores[user_id]["score"] += score * 0.2
            combined_scores[user_id]["sources"].append(RecommendationSource.POPULARITY)
        
        # Sort by score
        sorted_candidates = sorted(
            combined_scores.items(),
            key=lambda x: x[1]["score"],
            reverse=True
        )
        
        # Get user details and build recommendations
        user_recs = []
        for user_id, data in sorted_candidates[request.offset:request.offset + request.limit]:
            user_info = await self._get_user_info(user_id)
            if not user_info:
                continue
            
            # Filter: only creators (high follower count or verified)
            if user_info.get("follower_count", 0) < 1000 and not user_info.get("is_verified"):
                continue
            
            # Build recommendation
            source = data["sources"][0] if data["sources"] else RecommendationSource.HYBRID
            reason = self._generate_reason(source, data["mutuals"], user_info)
            
            # Get mutual friends
            mutual_friends = list(self.graph.get_mutual_friends(request.user_id, user_id))[:5]
            
            rec = UserRecommendation(
                user_id=user_id,
                username=user_info.get("username", ""),
                display_name=user_info.get("display_name"),
                avatar_url=user_info.get("avatar_url"),
                bio=user_info.get("bio"),
                follower_count=user_info.get("follower_count", 0),
                post_count=user_info.get("post_count", 0),
                is_verified=user_info.get("is_verified", False),
                is_creator=True,
                score=data["score"],
                source=source,
                reason=reason,
                mutual_friends_count=data["mutuals"],
                mutual_friends=mutual_friends,
                common_interests=self._get_common_interests(user_profile, user_id)
            )
            
            user_recs.append(rec)
        
        logger.info(f"Generated {len(user_recs)} creator recommendations")
        
        return RecommendationResponse(
            type=RecommendationType.CREATORS_FOR_YOU,
            users=user_recs,
            next_offset=request.offset + len(user_recs),
            has_more=len(sorted_candidates) > request.offset + request.limit
        )
    
    async def _recommend_users(
        self, request: RecommendationRequest, user_profile: UserProfile
    ) -> RecommendationResponse:
        """
        Suggested for You (Vignette)
        
        General user recommendations (not just creators)
        Mix of friends of friends, similar users, and popular accounts
        """
        logger.info(f"Generating user recommendations for {request.user_id}")
        
        following_ids = set(user_profile.following_ids)
        exclude_ids = set(request.exclude_ids) | following_ids | {request.user_id}
        
        # 1. Friends of Friends (50% weight - social graph priority)
        fof_recs = self.graph.friends_of_friends(
            request.user_id,
            max_distance=2,
            exclude_ids=exclude_ids
        )
        
        # 2. Collaborative filtering (30% weight)
        collab_recs = self.collaborative.recommend_users(
            request.user_id,
            following_ids,
            n_recommendations=30,
            exclude_ids=exclude_ids
        )
        
        # 3. Popular users (20% weight)
        popular_recs = self.graph.recommend_by_popularity(
            request.user_id,
            following_ids,
            top_k=20,
            exclude_ids=exclude_ids
        )
        
        # Combine
        combined_scores = defaultdict(lambda: {
            "score": 0.0,
            "sources": [],
            "mutuals": 0
        })
        
        for user_id, distance, mutuals in fof_recs[:50]:
            score = 1.0 / distance
            combined_scores[user_id]["score"] += score * 0.5
            combined_scores[user_id]["sources"].append(RecommendationSource.GRAPH_BASED)
            combined_scores[user_id]["mutuals"] = mutuals
        
        for user_id, score in collab_recs:
            combined_scores[user_id]["score"] += score * 0.3
            combined_scores[user_id]["sources"].append(RecommendationSource.COLLABORATIVE)
        
        for user_id, score in popular_recs:
            combined_scores[user_id]["score"] += score * 0.2
            combined_scores[user_id]["sources"].append(RecommendationSource.POPULARITY)
        
        # Sort and build recommendations
        sorted_candidates = sorted(
            combined_scores.items(),
            key=lambda x: x[1]["score"],
            reverse=True
        )
        
        user_recs = []
        for user_id, data in sorted_candidates[request.offset:request.offset + request.limit]:
            user_info = await self._get_user_info(user_id)
            if not user_info:
                continue
            
            source = data["sources"][0] if data["sources"] else RecommendationSource.HYBRID
            reason = self._generate_reason(source, data["mutuals"], user_info)
            mutual_friends = list(self.graph.get_mutual_friends(request.user_id, user_id))[:5]
            
            rec = UserRecommendation(
                user_id=user_id,
                username=user_info.get("username", ""),
                display_name=user_info.get("display_name"),
                avatar_url=user_info.get("avatar_url"),
                bio=user_info.get("bio"),
                follower_count=user_info.get("follower_count", 0),
                post_count=user_info.get("post_count", 0),
                is_verified=user_info.get("is_verified", False),
                is_creator=user_info.get("follower_count", 0) > 1000,
                score=data["score"],
                source=source,
                reason=reason,
                mutual_friends_count=data["mutuals"],
                mutual_friends=mutual_friends,
                common_interests=self._get_common_interests(user_profile, user_id)
            )
            
            user_recs.append(rec)
        
        logger.info(f"Generated {len(user_recs)} user recommendations")
        
        return RecommendationResponse(
            type=RecommendationType.SUGGESTED_FOR_YOU,
            users=user_recs,
            next_offset=request.offset + len(user_recs),
            has_more=len(sorted_candidates) > request.offset + request.limit
        )
    
    async def _recommend_communities(
        self, request: RecommendationRequest, user_profile: UserProfile
    ) -> RecommendationResponse:
        """
        Communities for You
        
        Algorithm:
        1. Get communities friends are in
        2. Get communities matching interests
        3. Add popular communities
        4. Score and rank
        """
        logger.info(f"Generating community recommendations for {request.user_id}")
        
        user_communities = set(user_profile.community_ids)
        exclude_ids = set(request.exclude_ids) | user_communities
        
        # 1. Communities friends are in (50% weight)
        friend_communities = await self._get_friend_communities(
            user_profile.friend_ids,
            exclude_ids
        )
        
        # 2. Interest-based communities (30% weight)
        interest_communities = await self._get_interest_communities(
            user_profile.interest_topics,
            user_profile.interest_categories,
            exclude_ids
        )
        
        # 3. Popular communities (20% weight)
        popular_communities = await self._get_popular_communities(
            request.categories,
            exclude_ids
        )
        
        # Combine scores
        combined_scores = defaultdict(lambda: {
            "score": 0.0,
            "sources": [],
            "mutual_members": 0
        })
        
        for comm_id, mutual_count in friend_communities.items():
            combined_scores[comm_id]["score"] += (mutual_count / 10.0) * 0.5
            combined_scores[comm_id]["sources"].append(RecommendationSource.GRAPH_BASED)
            combined_scores[comm_id]["mutual_members"] = mutual_count
        
        for comm_id, match_score in interest_communities.items():
            combined_scores[comm_id]["score"] += match_score * 0.3
            combined_scores[comm_id]["sources"].append(RecommendationSource.CONTENT_BASED)
        
        for comm_id, pop_score in popular_communities.items():
            combined_scores[comm_id]["score"] += pop_score * 0.2
            combined_scores[comm_id]["sources"].append(RecommendationSource.POPULARITY)
        
        # Sort and build
        sorted_candidates = sorted(
            combined_scores.items(),
            key=lambda x: x[1]["score"],
            reverse=True
        )
        
        community_recs = []
        for comm_id, data in sorted_candidates[request.offset:request.offset + request.limit]:
            comm_info = await self._get_community_info(comm_id)
            if not comm_info:
                continue
            
            source = data["sources"][0] if data["sources"] else RecommendationSource.HYBRID
            reason = self._generate_community_reason(source, data["mutual_members"], comm_info)
            
            rec = CommunityRecommendation(
                community_id=comm_id,
                name=comm_info.get("name", ""),
                description=comm_info.get("description"),
                cover_photo=comm_info.get("cover_photo"),
                category=comm_info.get("category", ""),
                member_count=comm_info.get("member_count", 0),
                post_count=comm_info.get("post_count", 0),
                is_verified=comm_info.get("is_verified", False),
                score=data["score"],
                source=source,
                reason=reason,
                mutual_members_count=data["mutual_members"],
                matching_interests=self._get_matching_interests(user_profile, comm_info)
            )
            
            community_recs.append(rec)
        
        logger.info(f"Generated {len(community_recs)} community recommendations")
        
        return RecommendationResponse(
            type=RecommendationType.COMMUNITIES_FOR_YOU,
            communities=community_recs,
            next_offset=request.offset + len(community_recs),
            has_more=len(sorted_candidates) > request.offset + request.limit
        )
    
    # Helper methods
    
    async def _get_user_profile(self, user_id: str) -> UserProfile:
        """Get or build user profile"""
        # In production: Query from DB/cache
        return UserProfile(user_id=user_id)
    
    async def _get_user_info(self, user_id: str) -> Optional[Dict]:
        """Get user information"""
        # In production: Query user service
        return {
            "user_id": user_id,
            "username": f"user_{user_id}",
            "follower_count": 5000,
            "post_count": 100,
            "is_verified": False
        }
    
    async def _get_community_info(self, community_id: str) -> Optional[Dict]:
        """Get community information"""
        # In production: Query community service
        return {
            "community_id": community_id,
            "name": f"Community {community_id}",
            "category": "technology",
            "member_count": 1000
        }
    
    async def _get_friend_communities(
        self, friend_ids: List[str], exclude_ids: Set[str]
    ) -> Dict[str, int]:
        """Get communities that friends are in"""
        # In production: Query DB
        return {}
    
    async def _get_interest_communities(
        self, topics: List[str], categories: List[str], exclude_ids: Set[str]
    ) -> Dict[str, float]:
        """Get communities matching interests"""
        # In production: Query DB with matching
        return {}
    
    async def _get_popular_communities(
        self, categories: Optional[List[str]], exclude_ids: Set[str]
    ) -> Dict[str, float]:
        """Get popular communities"""
        # In production: Query DB sorted by member_count
        return {}
    
    def _generate_reason(
        self, source: RecommendationSource, mutuals: int, user_info: Dict
    ) -> str:
        """Generate recommendation reason"""
        if source == RecommendationSource.GRAPH_BASED and mutuals > 0:
            return f"{mutuals} mutual friend{'s' if mutuals > 1 else ''}"
        elif source == RecommendationSource.COLLABORATIVE:
            return "Similar to creators you follow"
        elif source == RecommendationSource.POPULARITY:
            return "Popular creator"
        else:
            return "Suggested for you"
    
    def _generate_community_reason(
        self, source: RecommendationSource, mutuals: int, comm_info: Dict
    ) -> str:
        """Generate community recommendation reason"""
        if mutuals > 0:
            return f"{mutuals} friend{'s' if mutuals > 1 else ''} are members"
        elif source == RecommendationSource.CONTENT_BASED:
            return "Matches your interests"
        else:
            return "Popular in your area"
    
    def _get_common_interests(self, user_profile: UserProfile, target_id: str) -> List[str]:
        """Get common interests"""
        # In production: Query and compare
        return []
    
    def _get_matching_interests(self, user_profile: UserProfile, comm_info: Dict) -> List[str]:
        """Get matching interests for community"""
        # In production: Compare topics/tags
        return []
