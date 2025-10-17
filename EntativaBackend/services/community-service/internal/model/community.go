package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Community represents a user community
type Community struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	CoverPhoto  string         `json:"cover_photo" db:"cover_photo"` // ONLY cover photo, NO profile photo
	Category    string         `json:"category" db:"category"`
	
	// Privacy & Visibility
	Privacy     CommunityPrivacy `json:"privacy" db:"privacy"`
	Visibility  CommunityVisibility `json:"visibility" db:"visibility"`
	
	// Settings
	IsVerified    bool           `json:"is_verified" db:"is_verified"`
	AllowPosts    bool           `json:"allow_posts" db:"allow_posts"`
	RequireApproval bool         `json:"require_approval" db:"require_approval"` // Posts need approval
	
	// Ownership
	CreatorID   uuid.UUID      `json:"creator_id" db:"creator_id"`
	
	// Stats
	MemberCount int            `json:"member_count" db:"member_count"`
	PostCount   int            `json:"post_count" db:"post_count"`
	
	// Metadata
	Tags        StringArray    `json:"tags" db:"tags"`
	Location    *string        `json:"location,omitempty" db:"location"`
	Website     *string        `json:"website,omitempty" db:"website"`
	
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

type CommunityPrivacy string

const (
	PrivacyPublic  CommunityPrivacy = "public"   // Anyone can see and join
	PrivacyPrivate CommunityPrivacy = "private"  // Visible, join by approval
	PrivacyHidden  CommunityPrivacy = "hidden"   // Invite-only, not searchable
)

type CommunityVisibility string

const (
	VisibilityListed   CommunityVisibility = "listed"   // Searchable
	VisibilityUnlisted CommunityVisibility = "unlisted" // Direct link only
)

// CreateCommunityRequest for creating a new community
type CreateCommunityRequest struct {
	Name            string         `json:"name" binding:"required,min=3,max=100"`
	Description     string         `json:"description" binding:"max=500"`
	CoverPhoto      string         `json:"cover_photo"`
	Category        string         `json:"category" binding:"required"`
	Privacy         CommunityPrivacy `json:"privacy" binding:"required,oneof=public private hidden"`
	Visibility      CommunityVisibility `json:"visibility" binding:"required,oneof=listed unlisted"`
	AllowPosts      bool           `json:"allow_posts"`
	RequireApproval bool           `json:"require_approval"`
	Tags            []string       `json:"tags"`
	Location        *string        `json:"location"`
	Website         *string        `json:"website"`
}

// UpdateCommunityRequest for updating a community
type UpdateCommunityRequest struct {
	Name            *string        `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Description     *string        `json:"description,omitempty" binding:"omitempty,max=500"`
	CoverPhoto      *string        `json:"cover_photo,omitempty"`
	Category        *string        `json:"category,omitempty"`
	Privacy         *CommunityPrivacy `json:"privacy,omitempty"`
	Visibility      *CommunityVisibility `json:"visibility,omitempty"`
	AllowPosts      *bool          `json:"allow_posts,omitempty"`
	RequireApproval *bool          `json:"require_approval,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	Location        *string        `json:"location,omitempty"`
	Website         *string        `json:"website,omitempty"`
}

// CommunityMember represents a member of a community
type CommunityMember struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	CommunityID uuid.UUID          `json:"community_id" db:"community_id"`
	UserID      uuid.UUID          `json:"user_id" db:"user_id"`
	Role        CommunityRole      `json:"role" db:"role"`
	Permissions CommunityPermissions `json:"permissions" db:"permissions"`
	
	// Member status
	Status      MemberStatus       `json:"status" db:"status"`
	IsMuted     bool               `json:"is_muted" db:"is_muted"`
	MutedUntil  *time.Time         `json:"muted_until,omitempty" db:"muted_until"`
	
	// Activity
	PostCount   int                `json:"post_count" db:"post_count"`
	CommentCount int               `json:"comment_count" db:"comment_count"`
	
	// Metadata
	InvitedBy   *uuid.UUID         `json:"invited_by,omitempty" db:"invited_by"`
	JoinedAt    time.Time          `json:"joined_at" db:"joined_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

type CommunityRole string

const (
	RoleOwner     CommunityRole = "owner"     // Creator, full control
	RoleAdmin     CommunityRole = "admin"     // Full management
	RoleModerator CommunityRole = "moderator" // Content moderation
	RoleMember    CommunityRole = "member"    // Regular member
)

type MemberStatus string

const (
	StatusActive  MemberStatus = "active"
	StatusPending MemberStatus = "pending"  // Join request pending
	StatusBanned  MemberStatus = "banned"
)

// CommunityPermissions - Granular permission system!
type CommunityPermissions struct {
	// Content
	CanPost        bool `json:"can_post"`
	CanComment     bool `json:"can_comment"`
	CanUploadMedia bool `json:"can_upload_media"`
	
	// Moderation
	CanModerate    bool `json:"can_moderate"`     // Approve/remove posts
	CanBanMembers  bool `json:"can_ban_members"`
	CanMuteMembers bool `json:"can_mute_members"`
	
	// Management
	CanInviteMembers  bool `json:"can_invite_members"`
	CanRemoveMembers  bool `json:"can_remove_members"`
	CanManageRoles    bool `json:"can_manage_roles"`
	CanEditCommunity  bool `json:"can_edit_community"`
	CanDeleteCommunity bool `json:"can_delete_community"`
	
	// Settings
	CanManageRules     bool `json:"can_manage_rules"`
	CanManageSettings  bool `json:"can_manage_settings"`
	CanViewAnalytics   bool `json:"can_view_analytics"`
	CanManageEvents    bool `json:"can_manage_events"`
}

// Default permissions by role
func GetDefaultPermissions(role CommunityRole) CommunityPermissions {
	switch role {
	case RoleOwner:
		return CommunityPermissions{
			CanPost: true, CanComment: true, CanUploadMedia: true,
			CanModerate: true, CanBanMembers: true, CanMuteMembers: true,
			CanInviteMembers: true, CanRemoveMembers: true, CanManageRoles: true,
			CanEditCommunity: true, CanDeleteCommunity: true,
			CanManageRules: true, CanManageSettings: true, CanViewAnalytics: true,
			CanManageEvents: true,
		}
	case RoleAdmin:
		return CommunityPermissions{
			CanPost: true, CanComment: true, CanUploadMedia: true,
			CanModerate: true, CanBanMembers: true, CanMuteMembers: true,
			CanInviteMembers: true, CanRemoveMembers: true, CanManageRoles: true,
			CanEditCommunity: true, CanDeleteCommunity: false,
			CanManageRules: true, CanManageSettings: true, CanViewAnalytics: true,
			CanManageEvents: true,
		}
	case RoleModerator:
		return CommunityPermissions{
			CanPost: true, CanComment: true, CanUploadMedia: true,
			CanModerate: true, CanBanMembers: false, CanMuteMembers: true,
			CanInviteMembers: true, CanRemoveMembers: false, CanManageRoles: false,
			CanEditCommunity: false, CanDeleteCommunity: false,
			CanManageRules: false, CanManageSettings: false, CanViewAnalytics: false,
			CanManageEvents: false,
		}
	case RoleMember:
		return CommunityPermissions{
			CanPost: true, CanComment: true, CanUploadMedia: true,
			CanModerate: false, CanBanMembers: false, CanMuteMembers: false,
			CanInviteMembers: false, CanRemoveMembers: false, CanManageRoles: false,
			CanEditCommunity: false, CanDeleteCommunity: false,
			CanManageRules: false, CanManageSettings: false, CanViewAnalytics: false,
			CanManageEvents: false,
		}
	default:
		return CommunityPermissions{}
	}
}

// Scan implements sql.Scanner for CommunityPermissions
func (p *CommunityPermissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Value implements driver.Valuer for CommunityPermissions
func (p CommunityPermissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// JoinRequest represents a request to join a private community
type JoinRequest struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Message     string    `json:"message" db:"message"`
	Status      string    `json:"status" db:"status"` // pending, approved, rejected
	ReviewedBy  *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
}

// MemberInvite represents an invitation to join a community
type MemberInvite struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	InvitedUserID uuid.UUID `json:"invited_user_id" db:"invited_user_id"`
	InvitedBy   uuid.UUID `json:"invited_by" db:"invited_by"`
	Status      string    `json:"status" db:"status"` // pending, accepted, declined, expired
	Message     *string   `json:"message,omitempty" db:"message"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty" db:"responded_at"`
}

// BannedMember represents a banned user
type BannedMember struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	BannedBy    uuid.UUID `json:"banned_by" db:"banned_by"`
	Reason      string    `json:"reason" db:"reason"`
	IsPermanent bool      `json:"is_permanent" db:"is_permanent"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// StringArray for PostgreSQL array support
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return json.Marshal(a)
}
