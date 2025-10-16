package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"socialink/search-service/internal/elasticsearch"
	"socialink/search-service/internal/model"

	"github.com/go-redis/redis/v8"
)

type AutocompleteService struct {
	es    *elasticsearch.Client
	redis *redis.Client
}

func NewAutocompleteService(es *elasticsearch.Client, redis *redis.Client) *AutocompleteService {
	return &AutocompleteService{
		es:    es,
		redis: redis,
	}
}

// Autocomplete provides search suggestions as user types
func (s *AutocompleteService) Autocomplete(ctx context.Context, req *model.AutocompleteRequest) (*model.AutocompleteResponse, error) {
	start := time.Now()

	// Validate
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit > 20 {
		req.Limit = 20
	}

	// Normalize query
	query := strings.TrimSpace(strings.ToLower(req.Query))
	if len(query) < 2 {
		return &model.AutocompleteResponse{
			Query:       req.Query,
			Suggestions: []model.AutocompleteSuggestion{},
			Took:        time.Since(start).Milliseconds(),
		}, nil
	}

	// Try cache first
	cacheKey := fmt.Sprintf("autocomplete:%s:%s:%d", req.Type, query, req.Limit)
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		cached.Took = time.Since(start).Milliseconds()
		return cached, nil
	}

	// Get suggestions from Elasticsearch
	suggestions, err := s.getSuggestions(ctx, query, req.Type, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("error getting suggestions: %w", err)
	}

	// Include popular searches
	if len(suggestions) < req.Limit {
		popularSuggestions := s.getPopularSuggestions(ctx, query, req.Type, req.Limit-len(suggestions))
		suggestions = append(suggestions, popularSuggestions...)
	}

	response := &model.AutocompleteResponse{
		Query:       req.Query,
		Suggestions: suggestions,
		Took:        time.Since(start).Milliseconds(),
	}

	// Cache results
	s.cacheResults(ctx, cacheKey, response)

	return response, nil
}

// getSuggestions gets suggestions from Elasticsearch completion suggester
func (s *AutocompleteService) getSuggestions(ctx context.Context, query string, searchType model.SearchType, limit int) ([]model.AutocompleteSuggestion, error) {
	var suggestions []model.AutocompleteSuggestion

	// Determine which fields to use for suggestion
	fieldMap := map[model.SearchType][]string{
		model.SearchTypeUser:     {"username.suggest"},
		model.SearchTypeHashtag:  {"tag.suggest"},
		model.SearchTypeLocation: {"name.suggest"},
		model.SearchTypeAll: {
			"username.suggest",
			"tag.suggest",
			"name.suggest",
		},
	}

	fields, ok := fieldMap[searchType]
	if !ok {
		fields = []string{"username.suggest", "tag.suggest"}
	}

	// Query each field
	for _, field := range fields {
		// Determine index
		var index string
		switch {
		case strings.Contains(field, "username"):
			index = elasticsearch.IndexUsers
		case strings.Contains(field, "tag"):
			index = elasticsearch.IndexHashtags
		case strings.Contains(field, "name"):
			index = elasticsearch.IndexLocations
		default:
			continue
		}

		// Build suggest query
		suggestQuery := elasticsearch.BuildAutocompleteQuery(query, field)

		queryJSON, err := json.Marshal(suggestQuery)
		if err != nil {
			log.Printf("Error marshaling suggest query: %v", err)
			continue
		}

		// Execute
		res, err := s.es.GetClient().Search(
			s.es.GetClient().Search.WithContext(ctx),
			s.es.GetClient().Search.WithIndex(index),
			s.es.GetClient().Search.WithBody(strings.NewReader(string(queryJSON))),
		)
		if err != nil {
			log.Printf("Error executing suggest query: %v", err)
			continue
		}

		// Parse response
		var suggestResp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&suggestResp); err != nil {
			res.Body.Close()
			log.Printf("Error parsing suggest response: %v", err)
			continue
		}
		res.Body.Close()

		// Extract suggestions
		if suggest, ok := suggestResp["suggest"].(map[string]interface{}); ok {
			if autocomplete, ok := suggest["autocomplete"].([]interface{}); ok {
				for _, item := range autocomplete {
					itemMap := item.(map[string]interface{})
					if options, ok := itemMap["options"].([]interface{}); ok {
						for _, opt := range options {
							optMap := opt.(map[string]interface{})
							
							suggestion := model.AutocompleteSuggestion{
								Text:  optMap["text"].(string),
								Score: optMap["_score"].(float64),
							}

							// Set type
							switch index {
							case elasticsearch.IndexUsers:
								suggestion.Type = model.SearchTypeUser
							case elasticsearch.IndexHashtags:
								suggestion.Type = model.SearchTypeHashtag
							case elasticsearch.IndexLocations:
								suggestion.Type = model.SearchTypeLocation
							}

							// Add metadata
							if source, ok := optMap["_source"].(map[string]interface{}); ok {
								suggestion.Metadata = source
							}

							suggestions = append(suggestions, suggestion)

							if len(suggestions) >= limit {
								break
							}
						}
					}
					if len(suggestions) >= limit {
						break
					}
				}
			}
		}

		if len(suggestions) >= limit {
			break
		}
	}

	return suggestions, nil
}

// getPopularSuggestions gets popular/trending suggestions matching the query
func (s *AutocompleteService) getPopularSuggestions(ctx context.Context, query string, searchType model.SearchType, limit int) []model.AutocompleteSuggestion {
	if limit <= 0 {
		return nil
	}

	// Get trending searches
	trendingKey := fmt.Sprintf("search:trending:%s", searchType)
	results, err := s.redis.ZRevRangeWithScores(ctx, trendingKey, 0, 49).Result()
	if err != nil {
		return nil
	}

	var suggestions []model.AutocompleteSuggestion
	for _, result := range results {
		term := result.Member.(string)
		
		// Check if term matches query (prefix or contains)
		if strings.HasPrefix(strings.ToLower(term), query) || 
		   strings.Contains(strings.ToLower(term), query) {
			
			suggestions = append(suggestions, model.AutocompleteSuggestion{
				Text:  term,
				Type:  searchType,
				Score: result.Score,
				Metadata: map[string]interface{}{
					"trending": true,
					"count":    int64(result.Score),
				},
			})

			if len(suggestions) >= limit {
				break
			}
		}
	}

	return suggestions
}

// getFromCache retrieves autocomplete results from cache
func (s *AutocompleteService) getFromCache(ctx context.Context, key string) (*model.AutocompleteResponse, error) {
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response model.AutocompleteResponse
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// cacheResults caches autocomplete results
func (s *AutocompleteService) cacheResults(ctx context.Context, key string, response *model.AutocompleteResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling autocomplete response for cache: %v", err)
		return
	}

	// Cache for 15 minutes (autocomplete queries are popular)
	if err := s.redis.Set(ctx, key, data, 15*time.Minute).Err(); err != nil {
		log.Printf("Error caching autocomplete results: %v", err)
	}
}

// GetRecentSearches gets user's recent searches for autocomplete
func (s *AutocompleteService) GetRecentSearches(ctx context.Context, userID string, limit int) ([]model.AutocompleteSuggestion, error) {
	if limit == 0 {
		limit = 10
	}

	key := fmt.Sprintf("search:history:%s", userID)
	results, err := s.redis.LRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	var suggestions []model.AutocompleteSuggestion
	seen := make(map[string]bool)

	for _, data := range results {
		var history model.SearchHistory
		if err := json.Unmarshal([]byte(data), &history); err != nil {
			continue
		}

		// Deduplicate
		if seen[history.Query] {
			continue
		}
		seen[history.Query] = true

		suggestions = append(suggestions, model.AutocompleteSuggestion{
			Text: history.Query,
			Type: history.Type,
			Metadata: map[string]interface{}{
				"recent":       true,
				"result_count": history.ResultCount,
				"searched_at":  history.SearchedAt,
			},
		})
	}

	return suggestions, nil
}
