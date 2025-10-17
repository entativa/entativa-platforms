"""
Graph-Based Recommendations
Uses social graph analysis (Friends of Friends, PageRank, etc.)
"""
import networkx as nx
from typing import List, Dict, Set, Tuple
from collections import defaultdict, deque
import logging

logger = logging.getLogger(__name__)


class SocialGraphRecommender:
    """
    Graph-based recommendation using social network analysis
    """
    
    def __init__(self):
        self.graph = nx.DiGraph()
        self.pagerank_scores = {}
        self.community_map = {}
    
    def build_graph(
        self,
        user_follows: Dict[str, Set[str]],
        mutual_only: bool = False
    ):
        """
        Build social graph from follow relationships
        """
        self.graph.clear()
        
        for user_id, follows in user_follows.items():
            self.graph.add_node(user_id)
            for followed_id in follows:
                self.graph.add_node(followed_id)
                self.graph.add_edge(user_id, followed_id)
                
                # Add reverse edge for mutual follows
                if mutual_only and followed_id in user_follows:
                    if user_id in user_follows[followed_id]:
                        self.graph.add_edge(followed_id, user_id)
        
        logger.info(f"Built graph: {self.graph.number_of_nodes()} nodes, {self.graph.number_of_edges()} edges")
    
    def compute_pagerank(self, alpha: float = 0.85):
        """
        Compute PageRank scores for all users
        Higher score = more influential
        """
        try:
            self.pagerank_scores = nx.pagerank(self.graph, alpha=alpha)
            logger.info(f"Computed PageRank for {len(self.pagerank_scores)} users")
        except Exception as e:
            logger.error(f"PageRank failed: {e}")
            self.pagerank_scores = {}
    
    def friends_of_friends(
        self,
        user_id: str,
        max_distance: int = 2,
        exclude_ids: Set[str] = None
    ) -> List[Tuple[str, int, int]]:
        """
        Find friends of friends (FoF)
        
        Returns list of (user_id, distance, mutual_friends_count)
        
        Algorithm:
        1. BFS from user up to max_distance
        2. Count mutual friends
        3. Score by distance and mutuals
        """
        if exclude_ids is None:
            exclude_ids = set()
        
        if user_id not in self.graph:
            return []
        
        # Get direct follows
        direct_follows = set(self.graph.successors(user_id))
        
        # BFS to find FoF
        visited = {user_id}
        queue = deque([(user_id, 0)])
        candidates = defaultdict(lambda: {"distance": float('inf'), "mutuals": set()})
        
        while queue:
            current, distance = queue.popleft()
            
            if distance >= max_distance:
                continue
            
            # Get neighbors
            for neighbor in self.graph.successors(current):
                if neighbor in visited:
                    continue
                
                visited.add(neighbor)
                new_distance = distance + 1
                
                # Track mutual friends
                if neighbor != user_id and neighbor not in direct_follows:
                    if neighbor not in exclude_ids:
                        if new_distance < candidates[neighbor]["distance"]:
                            candidates[neighbor]["distance"] = new_distance
                        candidates[neighbor]["mutuals"].add(current)
                
                if new_distance < max_distance:
                    queue.append((neighbor, new_distance))
        
        # Convert to list with scores
        recommendations = [
            (user, data["distance"], len(data["mutuals"]))
            for user, data in candidates.items()
        ]
        
        # Sort by: fewer hops, more mutuals
        recommendations.sort(key=lambda x: (x[1], -x[2]))
        
        logger.info(f"Found {len(recommendations)} FoF recommendations for {user_id}")
        
        return recommendations
    
    def get_mutual_friends(
        self,
        user_id: str,
        target_id: str
    ) -> Set[str]:
        """
        Get mutual friends between two users
        """
        if user_id not in self.graph or target_id not in self.graph:
            return set()
        
        user_follows = set(self.graph.successors(user_id))
        target_follows = set(self.graph.successors(target_id))
        
        return user_follows & target_follows
    
    def get_common_communities(
        self,
        user_id: str,
        target_id: str,
        community_members: Dict[str, Set[str]]
    ) -> List[str]:
        """
        Get communities both users are in
        """
        user_communities = [
            comm_id for comm_id, members in community_members.items()
            if user_id in members
        ]
        target_communities = [
            comm_id for comm_id, members in community_members.items()
            if target_id in members
        ]
        
        return list(set(user_communities) & set(target_communities))
    
    def recommend_by_popularity(
        self,
        user_id: str,
        user_follows: Set[str],
        top_k: int = 20,
        exclude_ids: Set[str] = None
    ) -> List[Tuple[str, float]]:
        """
        Recommend popular users (high PageRank)
        """
        if exclude_ids is None:
            exclude_ids = set()
        
        if not self.pagerank_scores:
            self.compute_pagerank()
        
        # Filter and sort
        candidates = [
            (user, score)
            for user, score in self.pagerank_scores.items()
            if user != user_id
            and user not in user_follows
            and user not in exclude_ids
        ]
        
        candidates.sort(key=lambda x: x[1], reverse=True)
        
        return candidates[:top_k]
    
    def compute_social_closeness(
        self,
        user_id: str,
        target_id: str
    ) -> float:
        """
        Compute social closeness score
        
        Based on:
        - Shortest path distance
        - Common neighbors
        - PageRank influence
        """
        if user_id not in self.graph or target_id not in self.graph:
            return 0.0
        
        score = 0.0
        
        # 1. Shortest path (inverse distance)
        try:
            path_length = nx.shortest_path_length(self.graph, user_id, target_id)
            score += 1.0 / path_length if path_length > 0 else 0
        except (nx.NetworkXNoPath, nx.NodeNotFound):
            pass
        
        # 2. Common neighbors
        user_neighbors = set(self.graph.successors(user_id))
        target_neighbors = set(self.graph.successors(target_id))
        common = len(user_neighbors & target_neighbors)
        union = len(user_neighbors | target_neighbors)
        if union > 0:
            score += common / union
        
        # 3. PageRank influence
        if self.pagerank_scores:
            target_rank = self.pagerank_scores.get(target_id, 0)
            score += target_rank * 10  # Scale up
        
        # Normalize to 0-1
        return min(score / 3.0, 1.0)
    
    def detect_communities(self):
        """
        Detect communities using Louvain algorithm
        """
        try:
            # Convert to undirected for community detection
            undirected = self.graph.to_undirected()
            
            # Use Louvain algorithm
            import community as community_louvain
            partition = community_louvain.best_partition(undirected)
            
            # Store community map
            self.community_map = defaultdict(set)
            for user, comm_id in partition.items():
                self.community_map[comm_id].add(user)
            
            logger.info(f"Detected {len(self.community_map)} communities")
        except Exception as e:
            logger.error(f"Community detection failed: {e}")
            self.community_map = {}
    
    def get_user_community(self, user_id: str) -> Optional[int]:
        """Get community ID for a user"""
        for comm_id, members in self.community_map.items():
            if user_id in members:
                return comm_id
        return None
