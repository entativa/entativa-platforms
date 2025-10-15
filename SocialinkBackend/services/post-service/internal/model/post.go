package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Post represents a social media post (Facebook-like)
type Post struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	UserID        uuid.UUID      `json:"user_id" db:"user_id"`
	Content       string         `json:"content" db:"content"`
	MediaIDs      MediaIDList    `json:"media_ids" db:"media_ids"`
	Privacy       Privacy        `json:"privacy" db:"privacy"`
	Location      *string        `json:"location,omitempty" db:"location"`
	TaggedUserIDs UUIDList       `json:"tagged_user_ids" db:"tagged_user_ids"`
	Feeling       *string        `json:"feeling,omitempty" db:"feeling"`
	Activity      *string        `json:"activity,omitempty" db:"activity"`
	LikesCount    int64          `json:"likes_count" db:"likes_count"`
	CommentsCount int64          `json:"comments_count" db:"comments_count"`
	SharesCount   int64          `json:"shares_count" db:"shares_count"`
	IsEdited      bool           `json:"is_edited" db:"is_edited"`
	EditedAt      *time.Time     `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	PostID    uuid.UUID  `json:"post_id" db:"post_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"` // For nested comments
	Content   string     `json:"content" db:"content"`
	MediaID   *uuid.UUID `json:"media_id,omitempty" db:"media_id"`
	LikesCount int64     `json:"likes_count" db:"likes_count"`
	IsEdited  bool       `json:"is_edited" db:"is_edited"`
	EditedAt  *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Like represents a like/reaction on a post or comment
type Like struct {
	ID         uuid.UUID    `json:"id" db:"id"`
	UserID     uuid.UUID    `json:"user_id" db:"user_id"`
	PostID     *uuid.UUID   `json:"post_id,omitempty" db:"post_id"`
	CommentID  *uuid.UUID   `json:"comment_id,omitempty" db:"comment_id"`
	ReactionType ReactionType `json:"reaction_type" db:"reaction_type"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
}

// Share represents a shared post
type Share struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	OriginalPostID uuid.UUID `json:"original_post_id" db:"original_post_id"`
	Caption     *string    `json:"caption,omitempty" db:"caption"`
	Privacy     Privacy    `json:"privacy" db:"privacy"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// Privacy levels for posts
type Privacy string

const (
	PrivacyPublic      Privacy = "public"
	PrivacyFriends     Privacy = "friends"
	PrivacyFriendsExcept Privacy = "friends_except"
	PrivacySpecificFriends Privacy = "specific_friends"
	PrivacyOnlyMe      Privacy = "only_me"
	PrivacyCustom      Privacy = "custom"
)

// ReactionType for Facebook-like reactions
type ReactionType string

const (
	ReactionLike  ReactionType = "like"
	ReactionLove  ReactionType = "love"
	ReactionHaha  ReactionType = "haha"
	ReactionWow   ReactionType = "wow"
	ReactionSad   ReactionType = "sad"
	ReactionAngry ReactionType = "angry"
	ReactionCare  ReactionType = "care"
)

// Custom types for JSONB fields
type UUIDList []uuid.UUID
type MediaIDList []uuid.UUID

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

// DTOs for API requests/responses

type CreatePostRequest struct {
	Content       string      `json:"content" binding:"required,max=10000"`
	MediaIDs      []uuid.UUID `json:"media_ids,omitempty"`
	Privacy       Privacy     `json:"privacy" binding:"required"`
	Location      *string     `json:"location,omitempty"`
	TaggedUserIDs []uuid.UUID `json:"tagged_user_ids,omitempty"`
	Feeling       *string     `json:"feeling,omitempty"`
	Activity      *string     `json:"activity,omitempty"`
}

type UpdatePostRequest struct {
	Content  *string     `json:"content,omitempty" binding:"omitempty,max=10000"`
	Privacy  *Privacy    `json:"privacy,omitempty"`
	Location *string     `json:"location,omitempty"`
	Feeling  *string     `json:"feeling,omitempty"`
	Activity *string     `json:"activity,omitempty"`
}

type CreateCommentRequest struct {
	Content  string     `json:"content" binding:"required,max=2000"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	MediaID  *uuid.UUID `json:"media_id,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

type LikeRequest struct {
	ReactionType ReactionType `json:"reaction_type" binding:"required"`
}

type SharePostRequest struct {
	Caption *string `json:"caption,omitempty" binding:"omitempty,max=1000"`
	Privacy Privacy `json:"privacy" binding:"required"`
}

type PostResponse struct {
	*Post
	User      *UserInfo     `json:"user,omitempty"`
	MediaURLs []MediaInfo   `json:"media_urls,omitempty"`
	TaggedUsers []UserInfo  `json:"tagged_users,omitempty"`
	UserReaction *ReactionType `json:"user_reaction,omitempty"`
	IsLiked   bool          `json:"is_liked"`
}

type CommentResponse struct {
	*Comment
	User         *UserInfo     `json:"user,omitempty"`
	MediaURL     *string       `json:"media_url,omitempty"`
	RepliesCount int64         `json:"replies_count,omitempty"`
	UserReaction *ReactionType `json:"user_reaction,omitempty"`
}

type UserInfo struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	ProfilePicture  *string   `json:"profile_picture,omitempty"`
}

type MediaInfo struct {
	ID           uuid.UUID `json:"id"`
	URL          string    `json:"url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Type         string    `json:"type"`
	Width        int32     `json:"width"`
	Height       int32     `json:"height"`
	Blurhash     *string   `json:"blurhash,omitempty"`
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
