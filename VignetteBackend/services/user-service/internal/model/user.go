package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a Vignette user account (Instagram-like)
type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Username          string     `json:"username" db:"username"`
	Email             string     `json:"email" db:"email"`
	FullName          string     `json:"full_name" db:"full_name"`
	Password          string     `json:"-" db:"password_hash"` // Never expose password in JSON
	PhoneNumber       *string    `json:"phone_number,omitempty" db:"phone_number"`
	Bio               *string    `json:"bio,omitempty" db:"bio"`
	Website           *string    `json:"website,omitempty" db:"website"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
	IsPrivate         bool       `json:"is_private" db:"is_private"`
	IsVerified        bool       `json:"is_verified" db:"is_verified"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	IsDeleted         bool       `json:"is_deleted" db:"is_deleted"`
	FollowersCount    int        `json:"followers_count" db:"followers_count"`
	FollowingCount    int        `json:"following_count" db:"following_count"`
	PostsCount        int        `json:"posts_count" db:"posts_count"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// SignupRequest represents the request body for user registration
type SignupRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=1,max=100"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	User        *UserResponse `json:"user"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int64         `json:"expires_in"` // seconds
}

// UserResponse represents user data returned to clients (sanitized)
type UserResponse struct {
	ID                uuid.UUID  `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	FullName          string     `json:"full_name"`
	PhoneNumber       *string    `json:"phone_number,omitempty"`
	Bio               *string    `json:"bio,omitempty"`
	Website           *string    `json:"website,omitempty"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty"`
	IsPrivate         bool       `json:"is_private"`
	IsVerified        bool       `json:"is_verified"`
	IsActive          bool       `json:"is_active"`
	FollowersCount    int        `json:"followers_count"`
	FollowingCount    int        `json:"following_count"`
	PostsCount        int        `json:"posts_count"`
	CreatedAt         time.Time  `json:"created_at"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
}

// ToUserResponse converts a User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:                u.ID,
		Username:          u.Username,
		Email:             u.Email,
		FullName:          u.FullName,
		PhoneNumber:       u.PhoneNumber,
		Bio:               u.Bio,
		Website:           u.Website,
		ProfilePictureURL: u.ProfilePictureURL,
		IsPrivate:         u.IsPrivate,
		IsVerified:        u.IsVerified,
		IsActive:          u.IsActive,
		FollowersCount:    u.FollowersCount,
		FollowingCount:    u.FollowingCount,
		PostsCount:        u.PostsCount,
		CreatedAt:         u.CreatedAt,
		LastLoginAt:       u.LastLoginAt,
	}
}

// PublicUserResponse represents public user data (for non-followers)
type PublicUserResponse struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Bio               *string   `json:"bio,omitempty"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	IsPrivate         bool      `json:"is_private"`
	IsVerified        bool      `json:"is_verified"`
	FollowersCount    int       `json:"followers_count"`
	FollowingCount    int       `json:"following_count"`
	PostsCount        int       `json:"posts_count"`
}

// ToPublicUserResponse converts a User to PublicUserResponse
func (u *User) ToPublicUserResponse() *PublicUserResponse {
	return &PublicUserResponse{
		ID:                u.ID,
		Username:          u.Username,
		FullName:          u.FullName,
		Bio:               u.Bio,
		ProfilePictureURL: u.ProfilePictureURL,
		IsPrivate:         u.IsPrivate,
		IsVerified:        u.IsVerified,
		FollowersCount:    u.FollowersCount,
		FollowingCount:    u.FollowingCount,
		PostsCount:        u.PostsCount,
	}
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FullName          *string `json:"full_name,omitempty" binding:"omitempty,min=1,max=100"`
	Bio               *string `json:"bio,omitempty" binding:"omitempty,max=150"`
	Website           *string `json:"website,omitempty"`
	PhoneNumber       *string `json:"phone_number,omitempty"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	IsPrivate         *bool   `json:"is_private,omitempty"`
}
