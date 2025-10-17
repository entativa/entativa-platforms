package model

import "time"

// Document represents a searchable document in Elasticsearch
type Document struct {
	ID        string                 `json:"id"`
	Type      SearchType             `json:"type"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Data      map[string]interface{} `json:"data"`
}

// UserDocument represents a user in search index
type UserDocument struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	Verified    bool      `json:"verified"`
	AvatarURL   string    `json:"avatar_url"`
	FollowerCount int     `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
	PostCount   int       `json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// PostDocument represents a post in search index
type PostDocument struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Caption     string    `json:"caption"`
	Content     string    `json:"content"` // For text posts
	MediaIDs    []string  `json:"media_ids"`
	MediaType   string    `json:"media_type"` // image, video, text
	HasMedia    bool      `json:"has_media"`
	Hashtags    []string  `json:"hashtags"`
	Location    string    `json:"location"`
	LikesCount  int64     `json:"likes_count"`
	CommentsCount int64   `json:"comments_count"`
	SharesCount int64     `json:"shares_count"`
	ViewsCount  int64     `json:"views_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// TakeDocument represents a Take in search index
type TakeDocument struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Caption     string    `json:"caption"`
	MediaID     string    `json:"media_id"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Duration    float64   `json:"duration"`
	Hashtags    []string  `json:"hashtags"`
	AudioTrackID string   `json:"audio_track_id,omitempty"`
	FilterUsed  string    `json:"filter_used,omitempty"`
	TrendID     string    `json:"trend_id,omitempty"`
	TemplateID  string    `json:"template_id,omitempty"`
	ViewsCount  int64     `json:"views_count"`
	LikesCount  int64     `json:"likes_count"`
	CommentsCount int64   `json:"comments_count"`
	RemixCount  int64     `json:"remix_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// HashtagDocument represents a hashtag in search index
type HashtagDocument struct {
	Tag         string    `json:"tag"`
	DisplayTag  string    `json:"display_tag"` // With # prefix
	UsageCount  int64     `json:"usage_count"`
	PostCount   int64     `json:"post_count"`
	TakeCount   int64     `json:"take_count"`
	GrowthRate  float64   `json:"growth_rate"`
	IsTrending  bool      `json:"is_trending"`
	Category    string    `json:"category,omitempty"`
	FirstUsed   time.Time `json:"first_used"`
	LastUsed    time.Time `json:"last_used"`
}

// LocationDocument represents a location in search index
type LocationDocument struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	PostCount   int64     `json:"post_count"`
	TakeCount   int64     `json:"take_count"`
	CheckinCount int64    `json:"checkin_count"`
	Category    string    `json:"category,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// IndexingRequest represents a request to index a document
type IndexingRequest struct {
	Action     string     `json:"action"` // index, update, delete
	DocumentType SearchType `json:"document_type"`
	DocumentID string     `json:"document_id"`
	Data       interface{} `json:"data,omitempty"`
}

// IndexingResponse represents the response of an indexing operation
type IndexingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	DocumentID string `json:"document_id,omitempty"`
}

// BulkIndexingRequest represents a bulk indexing request
type BulkIndexingRequest struct {
	Documents []IndexingRequest `json:"documents"`
}

// BulkIndexingResponse represents the response of a bulk indexing operation
type BulkIndexingResponse struct {
	Success      bool   `json:"success"`
	TotalDocs    int    `json:"total_docs"`
	IndexedDocs  int    `json:"indexed_docs"`
	FailedDocs   int    `json:"failed_docs"`
	Errors       []string `json:"errors,omitempty"`
}
