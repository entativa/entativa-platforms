package model

import "time"

// HashtagStats represents hashtag statistics
type HashtagStats struct {
	Tag         string    `json:"tag"`
	DisplayTag  string    `json:"display_tag"`
	UsageCount  int64     `json:"usage_count"`
	PostCount   int64     `json:"post_count"`
	TakeCount   int64     `json:"take_count"`
	GrowthRate  float64   `json:"growth_rate"`
	Rank        int       `json:"rank"`
	IsTrending  bool      `json:"is_trending"`
	Category    string    `json:"category,omitempty"`
}

// TrendingHashtagsRequest represents a request for trending hashtags
type TrendingHashtagsRequest struct {
	Limit      int    `json:"limit" binding:"omitempty,min=1,max=50"`
	Category   string `json:"category,omitempty"`
	TimeWindow string `json:"time_window,omitempty"` // 1h, 24h, 7d, 30d
}

// TrendingHashtagsResponse represents trending hashtags response
type TrendingHashtagsResponse struct {
	Hashtags   []HashtagStats `json:"hashtags"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// HashtagSearchRequest represents a hashtag search request
type HashtagSearchRequest struct {
	Query  string `json:"query" binding:"required,min=1"`
	Limit  int    `json:"limit" binding:"omitempty,min=1,max=50"`
}

// RelatedHashtagsRequest represents a request for related hashtags
type RelatedHashtagsRequest struct {
	Tag    string `json:"tag" binding:"required"`
	Limit  int    `json:"limit" binding:"omitempty,min=1,max=20"`
}

// RelatedHashtagsResponse represents related hashtags response
type RelatedHashtagsResponse struct {
	Tag            string         `json:"tag"`
	RelatedHashtags []HashtagStats `json:"related_hashtags"`
}
