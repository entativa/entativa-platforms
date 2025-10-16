package model

import (
	"time"

	"github.com/google/uuid"
)

// Follow represents a follow relationship (Vignette - Instagram style)
type Follow struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FollowerID  uuid.UUID `json:"follower_id" db:"follower_id"`   // Who is following
	FollowingID uuid.UUID `json:"following_id" db:"following_id"` // Who is being followed
	Status      string    `json:"status" db:"status"`             // active, removed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// FollowStats represents follow statistics for a user
type FollowStats struct {
	UserID         uuid.UUID `json:"user_id"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
}

// FollowRequest for API
type FollowRequest struct {
	FollowingID uuid.UUID `json:"following_id" binding:"required"`
}

// FollowResponse for API
type FollowResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	IsFollowing   bool   `json:"is_following"`
	FollowerCount int    `json:"follower_count,omitempty"`
}

// MutualFollowsResponse
type MutualFollowsResponse struct {
	MutualCount int         `json:"mutual_count"`
	MutualUsers []uuid.UUID `json:"mutual_users,omitempty"`
}
