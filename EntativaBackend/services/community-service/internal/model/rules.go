package model

import (
	"time"

	"github.com/google/uuid"
)

// CommunityRule represents a community rule/guideline
type CommunityRule struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Position    int       `json:"position" db:"position"` // Display order
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateRuleRequest for creating a rule
type CreateRuleRequest struct {
	Title       string `json:"title" binding:"required,min=5,max=200"`
	Description string `json:"description" binding:"required,min=10,max=1000"`
	Position    int    `json:"position"`
}

// UpdateRuleRequest for updating a rule
type UpdateRuleRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=5,max=200"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=10,max=1000"`
	Position    *int    `json:"position,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// ModerationAction represents a moderation action taken
type ModerationAction struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	CommunityID uuid.UUID        `json:"community_id" db:"community_id"`
	ModeratorID uuid.UUID        `json:"moderator_id" db:"moderator_id"`
	TargetID    uuid.UUID        `json:"target_id" db:"target_id"` // User or content ID
	TargetType  string           `json:"target_type" db:"target_type"` // user, post, comment
	Action      ModerationType   `json:"action" db:"action"`
	Reason      string           `json:"reason" db:"reason"`
	Details     *string          `json:"details,omitempty" db:"details"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
}

type ModerationType string

const (
	ActionBan           ModerationType = "ban"
	ActionUnban         ModerationType = "unban"
	ActionMute          ModerationType = "mute"
	ActionUnmute        ModerationType = "unmute"`
	ActionRemovePost    ModerationType = "remove_post"
	ActionApprovePost   ModerationType = "approve_post"
	ActionRemoveComment ModerationType = "remove_comment"
	ActionWarn          ModerationType = "warn"
)

// ReportedContent represents user-reported content
type ReportedContent struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	ContentID   uuid.UUID `json:"content_id" db:"content_id"` // Post or comment ID
	ContentType string    `json:"content_type" db:"content_type"` // post, comment
	ReporterID  uuid.UUID `json:"reporter_id" db:"reporter_id"`
	Reason      string    `json:"reason" db:"reason"`
	Details     *string   `json:"details,omitempty" db:"details"`
	Status      string    `json:"status" db:"status"` // pending, reviewed, action_taken, dismissed
	ReviewedBy  *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
}

// CreateReportRequest for reporting content
type CreateReportRequest struct {
	ContentID   uuid.UUID `json:"content_id" binding:"required"`
	ContentType string    `json:"content_type" binding:"required,oneof=post comment"`
	Reason      string    `json:"reason" binding:"required,oneof=spam harassment inappropriate violence misinformation other"`
	Details     *string   `json:"details,omitempty"`
}
