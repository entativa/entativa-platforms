package model

import (
	"time"

	"github.com/google/uuid"
)

// Socialink Friend System Models
// Friend requests + Friends + Follows

const (
	// Friend limit - WAY better than Facebook's 5,000!
	MaxFriendsLimit = 1500 // Reasonable limit to prevent spam
	
	// Request limits to prevent spam
	MaxPendingRequests = 100  // Max pending outgoing requests
	MaxDailyRequests   = 50   // Max friend requests per day
)

// FriendRequest represents a friend request
type FriendRequest struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	SenderID    uuid.UUID  `json:"sender_id" db:"sender_id"`
	ReceiverID  uuid.UUID  `json:"receiver_id" db:"receiver_id"`
	Status      string     `json:"status" db:"status"` // pending, accepted, rejected, cancelled
	Message     *string    `json:"message,omitempty" db:"message"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty" db:"responded_at"`
}

// Friend represents a friendship (bi-directional)
type Friend struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID1   uuid.UUID `json:"user_id_1" db:"user_id_1"` // Lower UUID
	UserID2   uuid.UUID `json:"user_id_2" db:"user_id_2"` // Higher UUID
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Follow (for pages/public figures, separate from friends)
type Follow struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FollowerID  uuid.UUID `json:"follower_id" db:"follower_id"`
	FollowingID uuid.UUID `json:"following_id" db:"following_id"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// FriendStats
type FriendStats struct {
	UserID                 uuid.UUID `json:"user_id"`
	FriendsCount           int       `json:"friends_count"`
	PendingRequestsCount   int       `json:"pending_requests_count"` // Incoming
	SentRequestsCount      int       `json:"sent_requests_count"`     // Outgoing
	FollowersCount         int       `json:"followers_count"`         // For pages
	FollowingCount         int       `json:"following_count"`
	MutualFriendsAvailable int       `json:"mutual_friends_available"` // How many more friends can add
}

// API Request/Response Models

type SendFriendRequestRequest struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Message    *string   `json:"message,omitempty"`
}

type FriendRequestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Request *FriendRequest `json:"request,omitempty"`
}

type RespondToFriendRequestRequest struct {
	Action string `json:"action" binding:"required,oneof=accept reject"` // accept or reject
}

type FriendResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	FriendsCount int    `json:"friends_count,omitempty"`
	IsFriend     bool   `json:"is_friend"`
}

type MutualFriendsResponse struct {
	MutualCount   int         `json:"mutual_count"`
	MutualFriends []uuid.UUID `json:"mutual_friends,omitempty"`
}
