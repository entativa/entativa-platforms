package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"vignette/search-service/internal/elasticsearch"
	"vignette/search-service/internal/model"

	"github.com/go-redis/redis/v8"
)

type SearchService struct {
	es    *elasticsearch.Client
	redis *redis.Client
}

func NewSearchService(es *elasticsearch.Client, redis *redis.Client) *SearchService {
	return &SearchService{
		es:    es,
		redis: redis,
	}
}

// Search performs a search across specified indices
func (s *SearchService) Search(ctx context.Context, req *model.SearchRequest) (*model.SearchResponse, error) {
	start := time.Now()

	// Validate and set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 50 {
		req.Limit = 50
	}

	// Get cache key
	cacheKey := s.getCacheKey(req)
	
	// Try cache first
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		cached.Took = time.Since(start).Milliseconds()
		return cached, nil
	}

	// Determine which indices to search
	indices := s.getIndices(req.Type)

	// Build query
	query := elasticsearch.BuildSearchQuery(req)

	// Execute search
	res, err := s.es.Search(ctx, indices, query)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}
	defer res.Body.Close()

	// Parse response
	var searchResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error parsing search response: %w", err)
	}

	// Extract results
	response := s.parseSearchResponse(searchResp, req.Query)
	response.Took = time.Since(start).Milliseconds()

	// Get search suggestions if few results
	if response.TotalHits < 5 {
		suggestions := s.getSearchSuggestions(ctx, req.Query, req.Type)
		response.Suggestions = suggestions
	}

	// Cache results
	s.cacheResults(ctx, cacheKey, response)

	// Record search (async)
	go s.recordSearch(context.Background(), req, response.TotalHits)

	return response, nil
}

// parseSearchResponse parses Elasticsearch response
func (s *SearchService) parseSearchResponse(resp map[string]interface{}, query string) *model.SearchResponse {
	hits := resp["hits"].(map[string]interface{})
	totalHits := int64(hits["total"].(map[string]interface{})["value"].(float64))

	var results []model.SearchResult
	for _, hit := range hits["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		
		result := model.SearchResult{
			ID:    hitMap["_id"].(string),
			Type:  model.SearchType(hitMap["_index"].(string)),
			Score: hitMap["_score"].(float64),
			Data:  hitMap["_source"].(map[string]interface{}),
		}

		// Extract highlights
		if highlight, ok := hitMap["highlight"].(map[string]interface{}); ok {
			result.Highlight = make(map[string][]string)
			for field, fragments := range highlight {
				fragArray := fragments.([]interface{})
				strFrags := make([]string, len(fragArray))
				for i, frag := range fragArray {
					strFrags[i] = frag.(string)
				}
				result.Highlight[field] = strFrags
			}
		}

		// Generate snippet
		result.Snippet = s.generateSnippet(result.Data, result.Highlight)

		results = append(results, result)
	}

	return &model.SearchResponse{
		Query:     query,
		TotalHits: totalHits,
		Results:   results,
	}
}

// generateSnippet generates a text snippet for the result
func (s *SearchService) generateSnippet(data map[string]interface{}, highlight map[string][]string) string {
	// Use highlighted text if available
	if len(highlight) > 0 {
		for _, fragments := range highlight {
			if len(fragments) > 0 {
				return fragments[0]
			}
		}
	}

	// Otherwise, use caption/content/bio
	if caption, ok := data["caption"].(string); ok && caption != "" {
		if len(caption) > 150 {
			return caption[:150] + "..."
		}
		return caption
	}

	if content, ok := data["content"].(string); ok && content != "" {
		if len(content) > 150 {
			return content[:150] + "..."
		}
		return content
	}

	if bio, ok := data["bio"].(string); ok && bio != "" {
		if len(bio) > 150 {
			return bio[:150] + "..."
		}
		return bio
	}

	return ""
}

// getIndices returns the indices to search based on search type
func (s *SearchService) getIndices(searchType model.SearchType) []string {
	switch searchType {
	case model.SearchTypeUser:
		return []string{elasticsearch.IndexUsers}
	case model.SearchTypePost:
		return []string{elasticsearch.IndexPosts}
	case model.SearchTypeTake:
		return []string{elasticsearch.IndexTakes}
	case model.SearchTypeHashtag:
		return []string{elasticsearch.IndexHashtags}
	case model.SearchTypeLocation:
		return []string{elasticsearch.IndexLocations}
	case model.SearchTypeAll:
		return []string{
			elasticsearch.IndexUsers,
			elasticsearch.IndexPosts,
			elasticsearch.IndexTakes,
			elasticsearch.IndexHashtags,
			elasticsearch.IndexLocations,
		}
	default:
		return []string{
			elasticsearch.IndexUsers,
			elasticsearch.IndexPosts,
			elasticsearch.IndexTakes,
		}
	}
}

// getSearchSuggestions gets search suggestions for query with few results
func (s *SearchService) getSearchSuggestions(ctx context.Context, query string, searchType model.SearchType) []string {
	// Split query into words
	words := strings.Fields(query)
	if len(words) == 0 {
		return nil
	}

	// Try fuzzy search with each word
	var suggestions []string
	seen := make(map[string]bool)

	for _, word := range words {
		if len(word) < 3 {
			continue
		}

		// Build fuzzy query
		fuzzyQuery := map[string]interface{}{
			"size": 5,
			"query": elasticsearch.BuildFuzzyQuery("username", word, "AUTO"),
		}

		indices := s.getIndices(searchType)
		res, err := s.es.Search(ctx, indices, fuzzyQuery)
		if err != nil {
			continue
		}

		var resp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			res.Body.Close()
			continue
		}
		res.Body.Close()

		hits := resp["hits"].(map[string]interface{})["hits"].([]interface{})
		for _, hit := range hits {
			hitMap := hit.(map[string]interface{})
			source := hitMap["_source"].(map[string]interface{})
			
			var suggestion string
			if username, ok := source["username"].(string); ok {
				suggestion = username
			} else if tag, ok := source["tag"].(string); ok {
				suggestion = tag
			} else if name, ok := source["name"].(string); ok {
				suggestion = name
			}

			if suggestion != "" && !seen[suggestion] {
				suggestions = append(suggestions, suggestion)
				seen[suggestion] = true
			}

			if len(suggestions) >= 5 {
				break
			}
		}

		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions
}

// getCacheKey generates a cache key for the search request
func (s *SearchService) getCacheKey(req *model.SearchRequest) string {
	return fmt.Sprintf("search:%s:%s:%d:%d", req.Type, req.Query, req.Limit, req.Offset)
}

// getFromCache retrieves results from cache
func (s *SearchService) getFromCache(ctx context.Context, key string) (*model.SearchResponse, error) {
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response model.SearchResponse
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// cacheResults caches search results
func (s *SearchService) cacheResults(ctx context.Context, key string, response *model.SearchResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling search response for cache: %v", err)
		return
	}

	// Cache for 5 minutes
	if err := s.redis.Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		log.Printf("Error caching search results: %v", err)
	}
}

// recordSearch records a search for analytics
func (s *SearchService) recordSearch(ctx context.Context, req *model.SearchRequest, resultCount int64) {
	// Increment search counter
	key := fmt.Sprintf("search:count:%s", req.Type)
	s.redis.Incr(ctx, key)

	// Add to trending searches (sorted set)
	trendingKey := fmt.Sprintf("search:trending:%s", req.Type)
	s.redis.ZIncrBy(ctx, trendingKey, 1, req.Query)

	// Expire after 7 days
	s.redis.Expire(ctx, trendingKey, 7*24*time.Hour)

	// Record in user's search history
	if req.UserID != "" {
		history := model.SearchHistory{
			ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
			UserID:      req.UserID,
			Query:       req.Query,
			Type:        req.Type,
			ResultCount: int(resultCount),
			SearchedAt:  time.Now(),
		}

		historyKey := fmt.Sprintf("search:history:%s", req.UserID)
		data, _ := json.Marshal(history)
		
		// Add to list (keep last 50)
		s.redis.LPush(ctx, historyKey, data)
		s.redis.LTrim(ctx, historyKey, 0, 49)
		s.redis.Expire(ctx, historyKey, 30*24*time.Hour)
	}
}

// GetSearchHistory gets user's search history
func (s *SearchService) GetSearchHistory(ctx context.Context, userID string, limit int) ([]model.SearchHistory, error) {
	if limit == 0 {
		limit = 20
	}

	key := fmt.Sprintf("search:history:%s", userID)
	results, err := s.redis.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	var history []model.SearchHistory
	for _, data := range results {
		var item model.SearchHistory
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			continue
		}
		history = append(history, item)
	}

	return history, nil
}

// DeleteSearchHistory deletes user's search history
func (s *SearchService) DeleteSearchHistory(ctx context.Context, userID string) error {
	key := fmt.Sprintf("search:history:%s", userID)
	return s.redis.Del(ctx, key).Err()
}

// GetTrendingSearches gets trending searches
func (s *SearchService) GetTrendingSearches(ctx context.Context, searchType model.SearchType, limit int) ([]model.TrendingSearch, error) {
	if limit == 0 {
		limit = 10
	}

	key := fmt.Sprintf("search:trending:%s", searchType)
	results, err := s.redis.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	var trending []model.TrendingSearch
	for i, result := range results {
		trending = append(trending, model.TrendingSearch{
			Query:      result.Member.(string),
			Type:       searchType,
			Count:      int64(result.Score),
			GrowthRate: 0, // Would calculate from historical data
			Rank:       i + 1,
		})
	}

	return trending, nil
}
