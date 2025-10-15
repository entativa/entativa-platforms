package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a Socialink user account (Facebook-like)
type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	FirstName         string     `json:"first_name" db:"first_name"`
	LastName          string     `json:"last_name" db:"last_name"`
	Email             string     `json:"email" db:"email"`
	Username          string     `json:"username" db:"username"`
	Password          string     `json:"-" db:"password_hash"` // Never expose password in JSON
	Birthday          time.Time  `json:"birthday" db:"birthday"`
	Gender            string     `json:"gender" db:"gender"`
	PhoneNumber       *string    `json:"phone_number,omitempty" db:"phone_number"`
	Bio               *string    `json:"bio,omitempty" db:"bio"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
	CoverPhotoURL     *string    `json:"cover_photo_url,omitempty" db:"cover_photo_url"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	IsDeleted         bool       `json:"is_deleted" db:"is_deleted"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// SignupRequest represents the request body for user registration
type SignupRequest struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=128"`
	Birthday  string `json:"birthday" binding:"required"` // Format: YYYY-MM-DD
	Gender    string `json:"gender" binding:"required,oneof=male female other prefer_not_to_say"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
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
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	Username          string     `json:"username"`
	Birthday          time.Time  `json:"birthday"`
	Gender            string     `json:"gender"`
	PhoneNumber       *string    `json:"phone_number,omitempty"`
	Bio               *string    `json:"bio,omitempty"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty"`
	CoverPhotoURL     *string    `json:"cover_photo_url,omitempty"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
}

// ToUserResponse converts a User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:                u.ID,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Email:             u.Email,
		Username:          u.Username,
		Birthday:          u.Birthday,
		Gender:            u.Gender,
		PhoneNumber:       u.PhoneNumber,
		Bio:               u.Bio,
		ProfilePictureURL: u.ProfilePictureURL,
		CoverPhotoURL:     u.CoverPhotoURL,
		IsActive:          u.IsActive,
		CreatedAt:         u.CreatedAt,
		LastLoginAt:       u.LastLoginAt,
	}
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FirstName         *string `json:"first_name,omitempty" binding:"omitempty,min=1,max=50"`
	LastName          *string `json:"last_name,omitempty" binding:"omitempty,min=1,max=50"`
	Bio               *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	PhoneNumber       *string `json:"phone_number,omitempty"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	CoverPhotoURL     *string `json:"cover_photo_url,omitempty"`
}
