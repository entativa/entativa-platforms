"""
Collaborative Filtering for User Recommendations
Uses user-user and item-item similarity
"""
import numpy as np
from typing import List, Dict, Tuple, Set
from collections import defaultdict
from scipy.sparse import csr_matrix
from sklearn.metrics.pairwise import cosine_similarity
import logging

logger = logging.getLogger(__name__)


class CollaborativeFilter:
    """
    Collaborative filtering recommendation engine
    Finds similar users and recommends what they follow
    """
    
    def __init__(self):
        self.user_similarity_cache = {}
        self.interaction_matrix = None
        self.user_index = {}
        self.index_to_user = {}
    
    def build_interaction_matrix(
        self,
        user_follows: Dict[str, Set[str]]
    ) -> csr_matrix:
        """
        Build user-user interaction matrix
        Rows = users, Cols = users they follow
        """
        # Create user index
        all_users = set(user_follows.keys())
        for follows in user_follows.values():
            all_users.update(follows)
        
        self.user_index = {user: idx for idx, user in enumerate(sorted(all_users))}
        self.index_to_user = {idx: user for user, idx in self.user_index.items()}
        
        # Build sparse matrix
        rows, cols, data = [], [], []
        for user, follows in user_follows.items():
            user_idx = self.user_index[user]
            for followed_user in follows:
                if followed_user in self.user_index:
                    followed_idx = self.user_index[followed_user]
                    rows.append(user_idx)
                    cols.append(followed_idx)
                    data.append(1.0)
        
        n_users = len(self.user_index)
        matrix = csr_matrix((data, (rows, cols)), shape=(n_users, n_users))
        
        logger.info(f"Built interaction matrix: {n_users} users, {len(data)} interactions")
        self.interaction_matrix = matrix
        return matrix
    
    def compute_user_similarity(
        self,
        user_id: str,
        top_k: int = 100
    ) -> List[Tuple[str, float]]:
        """
        Compute similar users using cosine similarity
        Returns list of (user_id, similarity_score)
        """
        # Check cache
        cache_key = f"{user_id}:{top_k}"
        if cache_key in self.user_similarity_cache:
            return self.user_similarity_cache[cache_key]
        
        if self.interaction_matrix is None:
            return []
        
        if user_id not in self.user_index:
            return []
        
        user_idx = self.user_index[user_id]
        user_vector = self.interaction_matrix[user_idx].toarray()
        
        # Compute cosine similarity with all users
        similarities = cosine_similarity(user_vector, self.interaction_matrix)[0]
        
        # Get top-k similar users (excluding self)
        similar_indices = np.argsort(similarities)[::-1][1:top_k+1]
        similar_users = [
            (self.index_to_user[idx], similarities[idx])
            for idx in similar_indices
            if similarities[idx] > 0
        ]
        
        # Cache results
        self.user_similarity_cache[cache_key] = similar_users
        
        return similar_users
    
    def recommend_users(
        self,
        user_id: str,
        user_follows: Set[str],
        n_recommendations: int = 20,
        exclude_ids: Set[str] = None
    ) -> List[Tuple[str, float]]:
        """
        Recommend users based on what similar users follow
        
        Algorithm:
        1. Find similar users
        2. Get users they follow
        3. Score by similarity * count
        4. Filter out already following
        """
        if exclude_ids is None:
            exclude_ids = set()
        
        # Find similar users
        similar_users = self.compute_user_similarity(user_id, top_k=100)
        
        if not similar_users:
            return []
        
        # Aggregate recommendations
        candidate_scores = defaultdict(float)
        candidate_counts = defaultdict(int)
        
        for similar_user_id, similarity in similar_users:
            # Get who this similar user follows
            similar_user_idx = self.user_index.get(similar_user_id)
            if similar_user_idx is None:
                continue
            
            # Get their follows
            follows_indices = self.interaction_matrix[similar_user_idx].nonzero()[1]
            for followed_idx in follows_indices:
                followed_id = self.index_to_user[followed_idx]
                
                # Skip if already following or excluded
                if followed_id in user_follows or followed_id in exclude_ids:
                    continue
                if followed_id == user_id:
                    continue
                
                # Score = similarity * weight
                candidate_scores[followed_id] += similarity
                candidate_counts[followed_id] += 1
        
        # Sort by score
        recommendations = sorted(
            candidate_scores.items(),
            key=lambda x: (x[1], candidate_counts[x[0]]),
            reverse=True
        )[:n_recommendations]
        
        logger.info(f"Generated {len(recommendations)} collaborative recommendations for {user_id}")
        
        return recommendations
    
    def get_jaccard_similarity(
        self,
        user_id: str,
        target_id: str,
        user_follows: Set[str],
        target_follows: Set[str]
    ) -> float:
        """
        Compute Jaccard similarity between two users
        J(A,B) = |A ∩ B| / |A ∪ B|
        """
        if not user_follows or not target_follows:
            return 0.0
        
        intersection = len(user_follows & target_follows)
        union = len(user_follows | target_follows)
        
        return intersection / union if union > 0 else 0.0
    
    def predict_affinity(
        self,
        user_id: str,
        target_id: str
    ) -> Tuple[float, float]:
        """
        Predict affinity score between user and target
        Returns (score, confidence)
        """
        if self.interaction_matrix is None:
            return 0.0, 0.0
        
        if user_id not in self.user_index or target_id not in self.user_index:
            return 0.0, 0.0
        
        user_idx = self.user_index[user_id]
        target_idx = self.user_index[target_id]
        
        # Check if already following
        if self.interaction_matrix[user_idx, target_idx] > 0:
            return 1.0, 1.0
        
        # Find similar users who follow target
        similar_users = self.compute_user_similarity(user_id, top_k=50)
        
        score = 0.0
        supporters = 0
        
        for similar_user_id, similarity in similar_users:
            similar_idx = self.user_index[similar_user_id]
            if self.interaction_matrix[similar_idx, target_idx] > 0:
                score += similarity
                supporters += 1
        
        # Normalize
        if similar_users:
            score /= len(similar_users)
            confidence = supporters / len(similar_users)
        else:
            confidence = 0.0
        
        return score, confidence


class ItemBasedCollaborativeFilter:
    """
    Item-based collaborative filtering
    "Users who followed A also followed B"
    """
    
    def __init__(self):
        self.item_similarity_cache = {}
    
    def compute_item_similarity(
        self,
        interaction_matrix: csr_matrix,
        item_id: str,
        item_index: Dict[str, int],
        top_k: int = 50
    ) -> List[Tuple[str, float]]:
        """
        Compute similar items (users) based on co-occurrence
        """
        if item_id not in item_index:
            return []
        
        item_idx = item_index[item_id]
        
        # Get users who follow this item
        item_vector = interaction_matrix[:, item_idx].toarray().flatten()
        
        # Compute similarity with all items
        similarities = cosine_similarity(
            item_vector.reshape(1, -1),
            interaction_matrix.T.toarray()
        )[0]
        
        # Get top-k
        similar_indices = np.argsort(similarities)[::-1][1:top_k+1]
        index_to_item = {idx: item for item, idx in item_index.items()}
        
        similar_items = [
            (index_to_item[idx], similarities[idx])
            for idx in similar_indices
            if similarities[idx] > 0
        ]
        
        return similar_items
