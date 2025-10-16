package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"vignette/search-service/internal/elasticsearch"
	"vignette/search-service/internal/model"

	"github.com/go-redis/redis/v8"
)

type HashtagService struct {
	es    *elasticsearch.Client
	redis *redis.Client
}

func NewHashtagService(es *elasticsearch.Client, redis *redis.Client) *HashtagService {
	return &HashtagService{
		es:    es,
		redis: redis,
	}
}

// GetTrendingHashtags returns trending hashtags
func (s *HashtagService) GetTrendingHashtags(ctx context.Context, req *model.TrendingHashtagsRequest) (*model.TrendingHashtagsResponse, error) {
	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 50 {
		req.Limit = 50
	}

	// Try cache first
	cacheKey := fmt.Sprintf("hashtags:trending:%s:%d", req.Category, req.Limit)
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	// Build query
	query := elasticsearch.BuildTrendingHashtagsQuery(req.Limit)

	// Add category filter if specified
	if req.Category != "" {
		if queryMap, ok := query["query"].(map[string]interface{}); ok {
			if boolMap, ok := queryMap["bool"].(map[string]interface{}); ok {
				if filters, ok := boolMap["filter"].([]map[string]interface{}); ok {
					filters = append(filters, map[string]interface{}{
						"term": map[string]interface{}{
							"category": req.Category,
						},
					})
					boolMap["filter"] = filters
				}
			}
		}
	}

	// Execute search
	res, err := s.es.Search(ctx, []string{elasticsearch.IndexHashtags}, query)
	if err != nil {
		return nil, fmt.Errorf("error searching trending hashtags: %w", err)
	}
	defer res.Body.Close()

	// Parse response
	var searchResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// Extract hashtags
	hits := searchResp["hits"].(map[string]interface{})["hits"].([]interface{})
	var hashtags []model.HashtagStats
	
	for i, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		hashtag := model.HashtagStats{
			Tag:        source["tag"].(string),
			DisplayTag: source["display_tag"].(string),
			Rank:       i + 1,
			IsTrending: true,
		}

		if usageCount, ok := source["usage_count"].(float64); ok {
			hashtag.UsageCount = int64(usageCount)
		}
		if postCount, ok := source["post_count"].(float64); ok {
			hashtag.PostCount = int64(postCount)
		}
		if takeCount, ok := source["take_count"].(float64); ok {
			hashtag.TakeCount = int64(takeCount)
		}
		if growthRate, ok := source["growth_rate"].(float64); ok {
			hashtag.GrowthRate = growthRate
		}
		if category, ok := source["category"].(string); ok {
			hashtag.Category = category
		}

		hashtags = append(hashtags, hashtag)
	}

	response := &model.TrendingHashtagsResponse{
		Hashtags:  hashtags,
		UpdatedAt: time.Now(),
	}

	// Cache for 5 minutes
	s.cacheResults(ctx, cacheKey, response, 5*time.Minute)

	return response, nil
}

// GetRelatedHashtags returns hashtags related to a given hashtag
func (s *HashtagService) GetRelatedHashtags(ctx context.Context, req *model.RelatedHashtagsRequest) (*model.RelatedHashtagsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit > 20 {
		req.Limit = 20
	}

	// Try cache first
	cacheKey := fmt.Sprintf("hashtags:related:%s:%d", req.Tag, req.Limit)
	if cached, err := s.getRelatedFromCache(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	// Build More Like This query
	query := elasticsearch.BuildRelatedHashtagsQuery(req.Tag, req.Limit)

	// Execute search
	res, err := s.es.Search(ctx, []string{elasticsearch.IndexHashtags}, query)
	if err != nil {
		return nil, fmt.Errorf("error searching related hashtags: %w", err)
	}
	defer res.Body.Close()

	// Parse response
	var searchResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// Extract related hashtags
	hits := searchResp["hits"].(map[string]interface{})["hits"].([]interface{})
	var relatedHashtags []model.HashtagStats

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		// Skip if it's the same hashtag
		if tag, ok := source["tag"].(string); ok && tag == req.Tag {
			continue
		}

		hashtag := model.HashtagStats{
			Tag:        source["tag"].(string),
			DisplayTag: source["display_tag"].(string),
		}

		if usageCount, ok := source["usage_count"].(float64); ok {
			hashtag.UsageCount = int64(usageCount)
		}
		if postCount, ok := source["post_count"].(float64); ok {
			hashtag.PostCount = int64(postCount)
		}
		if takeCount, ok := source["take_count"].(float64); ok {
			hashtag.TakeCount = int64(takeCount)
		}
		if growthRate, ok := source["growth_rate"].(float64); ok {
			hashtag.GrowthRate = growthRate
		}
		if isTrending, ok := source["is_trending"].(bool); ok {
			hashtag.IsTrending = isTrending
		}

		relatedHashtags = append(relatedHashtags, hashtag)
	}

	response := &model.RelatedHashtagsResponse{
		Tag:             req.Tag,
		RelatedHashtags: relatedHashtags,
	}

	// Cache for 10 minutes
	s.cacheRelatedResults(ctx, cacheKey, response, 10*time.Minute)

	return response, nil
}

// SearchHashtags searches for hashtags matching a query
func (s *HashtagService) SearchHashtags(ctx context.Context, req *model.HashtagSearchRequest) ([]model.HashtagStats, error) {
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 50 {
		req.Limit = 50
	}

	// Build search query
	query := map[string]interface{}{
		"size": req.Limit,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"tag": map[string]interface{}{
					"query":     req.Query,
					"fuzziness": "AUTO",
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"usage_count": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	// Execute search
	res, err := s.es.Search(ctx, []string{elasticsearch.IndexHashtags}, query)
	if err != nil {
		return nil, fmt.Errorf("error searching hashtags: %w", err)
	}
	defer res.Body.Close()

	// Parse response
	var searchResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// Extract hashtags
	hits := searchResp["hits"].(map[string]interface{})["hits"].([]interface{})
	var hashtags []model.HashtagStats

	for i, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		hashtag := model.HashtagStats{
			Tag:        source["tag"].(string),
			DisplayTag: source["display_tag"].(string),
			Rank:       i + 1,
		}

		if usageCount, ok := source["usage_count"].(float64); ok {
			hashtag.UsageCount = int64(usageCount)
		}
		if postCount, ok := source["post_count"].(float64); ok {
			hashtag.PostCount = int64(postCount)
		}
		if takeCount, ok := source["take_count"].(float64); ok {
			hashtag.TakeCount = int64(takeCount)
		}
		if growthRate, ok := source["growth_rate"].(float64); ok {
			hashtag.GrowthRate = growthRate
		}
		if isTrending, ok := source["is_trending"].(bool); ok {
			hashtag.IsTrending = isTrending
		}
		if category, ok := source["category"].(string); ok {
			hashtag.Category = category
		}

		hashtags = append(hashtags, hashtag)
	}

	return hashtags, nil
}

// getFromCache retrieves trending hashtags from cache
func (s *HashtagService) getFromCache(ctx context.Context, key string) (*model.TrendingHashtagsResponse, error) {
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response model.TrendingHashtagsResponse
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// cacheResults caches trending hashtags
func (s *HashtagService) cacheResults(ctx context.Context, key string, response *model.TrendingHashtagsResponse, ttl time.Duration) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling trending hashtags for cache: %v", err)
		return
	}

	if err := s.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Printf("Error caching trending hashtags: %v", err)
	}
}

// getRelatedFromCache retrieves related hashtags from cache
func (s *HashtagService) getRelatedFromCache(ctx context.Context, key string) (*model.RelatedHashtagsResponse, error) {
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response model.RelatedHashtagsResponse
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// cacheRelatedResults caches related hashtags
func (s *HashtagService) cacheRelatedResults(ctx context.Context, key string, response *model.RelatedHashtagsResponse, ttl time.Duration) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling related hashtags for cache: %v", err)
		return
	}

	if err := s.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Printf("Error caching related hashtags: %v", err)
	}
}
