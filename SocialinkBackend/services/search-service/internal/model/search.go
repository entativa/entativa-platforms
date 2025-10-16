package model

import "time"

// SearchType represents the type of search entity
type SearchType string

const (
	SearchTypeUser     SearchType = "user"
	SearchTypePost     SearchType = "post"
	SearchTypeTake     SearchType = "take"
	SearchTypeHashtag  SearchType = "hashtag"
	SearchTypeLocation SearchType = "location"
	SearchTypeAll      SearchType = "all"
)

// SearchRequest represents a search query
type SearchRequest struct {
	Query      string       `json:"query" binding:"required,min=1,max=100"`
	Type       SearchType   `json:"type" binding:"required"`
	Limit      int          `json:"limit" binding:"omitempty,min=1,max=50"`
	Offset     int          `json:"offset" binding:"omitempty,min=0"`
	Filters    SearchFilter `json:"filters,omitempty"`
	UserID     string       `json:"-"` // From auth header
}

// SearchFilter represents advanced search filters
type SearchFilter struct {
	// User filters
	Verified   *bool  `json:"verified,omitempty"`
	Location   string `json:"location,omitempty"`
	MinFollowers int   `json:"min_followers,omitempty"`
	
	// Post/Take filters
	HasMedia   *bool     `json:"has_media,omitempty"`
	MediaType  string    `json:"media_type,omitempty"` // image, video
	DateFrom   time.Time `json:"date_from,omitempty"`
	DateTo     time.Time `json:"date_to,omitempty"`
	MinLikes   int       `json:"min_likes,omitempty"`
	MinViews   int       `json:"min_views,omitempty"`
	
	// Hashtag filters
	TrendingOnly bool `json:"trending_only,omitempty"`
	
	// Location filters
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Radius     float64 `json:"radius,omitempty"` // In km
}

// SearchResponse represents search results
type SearchResponse struct {
	Query       string        `json:"query"`
	Type        SearchType    `json:"type"`
	TotalHits   int64         `json:"total_hits"`
	Results     []SearchResult `json:"results"`
	Took        int64         `json:"took_ms"`
	Suggestions []string      `json:"suggestions,omitempty"`
}

// SearchResult represents a single search result
type SearchResult struct {
	ID          string                 `json:"id"`
	Type        SearchType             `json:"type"`
	Score       float64                `json:"score"`
	Highlight   map[string][]string    `json:"highlight,omitempty"`
	Data        map[string]interface{} `json:"data"`
	Snippet     string                 `json:"snippet,omitempty"`
}

// AutocompleteRequest represents an autocomplete query
type AutocompleteRequest struct {
	Query  string     `json:"query" binding:"required,min=1,max=50"`
	Type   SearchType `json:"type" binding:"required"`
	Limit  int        `json:"limit" binding:"omitempty,min=1,max=20"`
	UserID string     `json:"-"`
}

// AutocompleteResponse represents autocomplete suggestions
type AutocompleteResponse struct {
	Query       string               `json:"query"`
	Suggestions []AutocompleteSuggestion `json:"suggestions"`
	Took        int64                `json:"took_ms"`
}

// AutocompleteSuggestion represents a single suggestion
type AutocompleteSuggestion struct {
	Text        string     `json:"text"`
	Type        SearchType `json:"type"`
	Score       float64    `json:"score"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TrendingSearch represents a trending search query
type TrendingSearch struct {
	Query       string     `json:"query"`
	Type        SearchType `json:"type"`
	Count       int64      `json:"count"`
	GrowthRate  float64    `json:"growth_rate"`
	Rank        int        `json:"rank"`
}

// SearchHistory represents a user's search history
type SearchHistory struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Query      string     `json:"query"`
	Type       SearchType `json:"type"`
	ResultCount int       `json:"result_count"`
	SearchedAt time.Time  `json:"searched_at"`
}

// SearchAnalytics represents search analytics data
type SearchAnalytics struct {
	TotalSearches    int64                     `json:"total_searches"`
	UniqueUsers      int64                     `json:"unique_users"`
	TopSearches      []TrendingSearch          `json:"top_searches"`
	SearchesByType   map[SearchType]int64      `json:"searches_by_type"`
	AvgResultsPerSearch float64                `json:"avg_results_per_search"`
	NoResultsRate    float64                   `json:"no_results_rate"`
}
