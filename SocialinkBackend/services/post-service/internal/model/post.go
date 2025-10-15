package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Post represents an Facebook-style post (image/video required)
type Post struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	UserID        uuid.UUID      `json:"user_id" db:"user_id"`
	Caption       string         `json:"caption" db:"caption"`
	MediaIDs      MediaIDList    `json:"media_ids" db:"media_ids"` // Required for Socialink
	Location      *Location      `json:"location,omitempty" db:"location"`
	TaggedUserIDs UUIDList       `json:"tagged_user_ids" db:"tagged_user_ids"`
	Hashtags      StringList     `json:"hashtags" db:"hashtags"`
	FilterUsed    *string        `json:"filter_used,omitempty" db:"filter_used"`
	IsCarousel    bool           `json:"is_carousel" db:"is_carousel"`
	LikesCount    int64          `json:"likes_count" db:"likes_count"`
	CommentsCount int64          `json:"comments_count" db:"comments_count"`
	ViewsCount    int64          `json:"views_count" db:"views_count"`
	SavesCount    int64          `json:"saves_count" db:"saves_count"`
	SharesCount   int64          `json:"shares_count" db:"shares_count"`
	IsEdited      bool           `json:"is_edited" db:"is_edited"`
	EditedAt      *time.Time     `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
	
	// Facebook-specific
	IsSponsored   bool           `json:"is_sponsored" db:"is_sponsored"`
	IsReels       bool           `json:"is_reels" db:"is_reels"`
	CommentsEnabled bool         `json:"comments_enabled" db:"comments_enabled"`
	LikesVisible  bool           `json:"likes_visible" db:"likes_visible"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	PostID    uuid.UUID   `json:"post_id" db:"post_id"`
	UserID    uuid.UUID   `json:"user_id" db:"user_id"`
	ParentID  *uuid.UUID  `json:"parent_id,omitempty" db:"parent_id"`
	Content   string      `json:"content" db:"content"`
	LikesCount int64      `json:"likes_count" db:"likes_count"`
	IsEdited  bool        `json:"is_edited" db:"is_edited"`
	EditedAt  *time.Time  `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Like represents a like on a post or comment
type Like struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	PostID    *uuid.UUID `json:"post_id,omitempty" db:"post_id"`
	CommentID *uuid.UUID `json:"comment_id,omitempty" db:"comment_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// Save represents a saved post (Facebook bookmark feature)
type Save struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	PostID    uuid.UUID  `json:"post_id" db:"post_id"`
	Collection *string   `json:"collection,omitempty" db:"collection"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// Location for Facebook-style location tagging
type Location struct {
	Name      string   `json:"name"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	PlaceID   *string  `json:"place_id,omitempty"`
}

// Custom types for JSONB fields
type UUIDList []uuid.UUID
type MediaIDList []uuid.UUID
type StringList []string

// Scan implements sql.Scanner
func (u *UUIDList) Scan(value interface{}) error {
	if value == nil {
		*u = []uuid.UUID{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		*u = []uuid.UUID{}
		return nil
	}
	
	return json.Unmarshal(bytes, u)
}

// Value implements driver.Valuer
func (u UUIDList) Value() (driver.Value, error) {
	if len(u) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(u)
}

// Scan implements sql.Scanner for MediaIDList
func (m *MediaIDList) Scan(value interface{}) error {
	if value == nil {
		*m = []uuid.UUID{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		*m = []uuid.UUID{}
		return nil
	}
	
	return json.Unmarshal(bytes, m)
}

// Value implements driver.Valuer for MediaIDList
func (m MediaIDList) Value() (driver.Value, error) {
	if len(m) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(m)
}

// Scan implements sql.Scanner for StringList
func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		*s = []string{}
		return nil
	}
	
	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer for StringList
func (s StringList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(s)
}

// Scan for Location
func (l *Location) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, l)
}

// Value for Location
func (l Location) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// DTOs for API requests/responses

type CreatePostRequest struct {
	Caption          string      `json:"caption" binding:"max=2200"`
	MediaIDs         []uuid.UUID `json:"media_ids" binding:"required,min=1"` // At least 1 media required
	Location         *Location   `json:"location,omitempty"`
	TaggedUserIDs    []uuid.UUID `json:"tagged_user_ids,omitempty"`
	Hashtags         []string    `json:"hashtags,omitempty"`
	FilterUsed       *string     `json:"filter_used,omitempty"`
	CommentsEnabled  bool        `json:"comments_enabled" binding:"required"`
	LikesVisible     bool        `json:"likes_visible" binding:"required"`
}

type UpdatePostRequest struct {
	Caption         *string    `json:"caption,omitempty" binding:"omitempty,max=2200"`
	Location        *Location  `json:"location,omitempty"`
	Hashtags        *[]string  `json:"hashtags,omitempty"`
	CommentsEnabled *bool      `json:"comments_enabled,omitempty"`
	LikesVisible    *bool      `json:"likes_visible,omitempty"`
}

type CreateCommentRequest struct {
	Content  string     `json:"content" binding:"required,max=2200"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,max=2200"`
}

type SavePostRequest struct {
	Collection *string `json:"collection,omitempty"`
}

type PostResponse struct {
	*Post
	User         *UserInfo   `json:"user,omitempty"`
	MediaURLs    []MediaInfo `json:"media_urls,omitempty"`
	TaggedUsers  []UserInfo  `json:"tagged_users,omitempty"`
	IsLiked      bool        `json:"is_liked"`
	IsSaved      bool        `json:"is_saved"`
}

type CommentResponse struct {
	*Comment
	User         *UserInfo `json:"user,omitempty"`
	RepliesCount int64     `json:"replies_count,omitempty"`
	IsLiked      bool      `json:"is_liked"`
}

type UserInfo struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	ProfilePicture *string   `json:"profile_picture,omitempty"`
	IsVerified     bool      `json:"is_verified"`
}

type MediaInfo struct {
	ID           uuid.UUID `json:"id"`
	URL          string    `json:"url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Type         string    `json:"type"`
	Width        int32     `json:"width"`
	Height       int32     `json:"height"`
	Blurhash     *string   `json:"blurhash,omitempty"`
	FilterUsed   *string   `json:"filter_used,omitempty"`
}

type FeedQuery struct {
	Cursor string `form:"cursor"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

type PostListResponse struct {
	Posts      []PostResponse `json:"posts"`
	NextCursor *string        `json:"next_cursor,omitempty"`
	HasMore    bool           `json:"has_more"`
}
