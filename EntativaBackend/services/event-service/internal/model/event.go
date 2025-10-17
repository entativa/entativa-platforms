package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event types
const (
	EventTypeInPerson EventType = "in_person"
	EventTypeOnline   EventType = "online"
)

// Event privacy levels
const (
	PrivacyPublic  EventPrivacy = "public"
	PrivacyPrivate EventPrivacy = "private"
	PrivacyFriends EventPrivacy = "friends"
)

// Event categories (Facebook-style)
const (
	CategorySocial        EventCategory = "social"
	CategoryBusiness      EventCategory = "business"
	CategoryEntertainment EventCategory = "entertainment"
	CategorySports        EventCategory = "sports"
	CategoryEducation     EventCategory = "education"
	CategoryReligious     EventCategory = "religious"
	CategoryCommunity     EventCategory = "community"
	CategoryCauses        EventCategory = "causes"
	CategoryHealth        EventCategory = "health"
	CategoryArts          EventCategory = "arts"
	CategoryOther         EventCategory = "other"
)

// RSVP statuses
const (
	RSVPGoing      RSVPStatus = "going"
	RSVPInterested RSVPStatus = "interested"
	RSVPNotGoing   RSVPStatus = "not_going"
)

type EventType string
type EventPrivacy string
type EventCategory string
type RSVPStatus string

// Event represents a Facebook-style event
type Event struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	CreatorID   uuid.UUID     `json:"creator_id" db:"creator_id"`
	Title       string        `json:"title" db:"title"`
	Description string        `json:"description" db:"description"`
	CoverPhoto  *string       `json:"cover_photo,omitempty" db:"cover_photo"`
	
	// Event details
	Type        EventType     `json:"type" db:"type"`
	Category    EventCategory `json:"category" db:"category"`
	Privacy     EventPrivacy  `json:"privacy" db:"privacy"`
	
	// Location
	LocationName *string   `json:"location_name,omitempty" db:"location_name"`
	Address      *string   `json:"address,omitempty" db:"address"`
	City         *string   `json:"city,omitempty" db:"city"`
	Country      *string   `json:"country,omitempty" db:"country"`
	Latitude     *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude    *float64  `json:"longitude,omitempty" db:"longitude"`
	OnlineLink   *string   `json:"online_link,omitempty" db:"online_link"` // For virtual events
	
	// Time
	StartTime   time.Time  `json:"start_time" db:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty" db:"end_time"`
	Timezone    string     `json:"timezone" db:"timezone"`
	
	// Recurring
	IsRecurring       bool       `json:"is_recurring" db:"is_recurring"`
	RecurrenceRule    *string    `json:"recurrence_rule,omitempty" db:"recurrence_rule"` // iCal RRULE format
	RecurrenceEndDate *time.Time `json:"recurrence_end_date,omitempty" db:"recurrence_end_date"`
	
	// Settings
	AllowGuestInvites bool `json:"allow_guest_invites" db:"allow_guest_invites"`
	RequireApproval   bool `json:"require_approval" db:"require_approval"`
	MaxAttendees      *int `json:"max_attendees,omitempty" db:"max_attendees"`
	
	// Co-hosts
	CoHosts UUIDArray `json:"co_hosts" db:"co_hosts"`
	
	// Stats (denormalized for performance)
	GoingCount      int `json:"going_count" db:"going_count"`
	InterestedCount int `json:"interested_count" db:"interested_count"`
	ViewCount       int `json:"view_count" db:"view_count"`
	
	// Status
	IsCancelled bool      `json:"is_cancelled" db:"is_cancelled"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// EventRSVP represents user's response to event
type EventRSVP struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	EventID   uuid.UUID  `json:"event_id" db:"event_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Status    RSVPStatus `json:"status" db:"status"`
	GuestCount int       `json:"guest_count" db:"guest_count"` // +1, +2, etc.
	
	// Check-in
	CheckedIn   bool       `json:"checked_in" db:"checked_in"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty" db:"checked_in_at"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// EventInvite for invited users
type EventInvite struct {
	ID        uuid.UUID `json:"id" db:"id"`
	EventID   uuid.UUID `json:"event_id" db:"event_id"`
	InviterID uuid.UUID `json:"inviter_id" db:"inviter_id"`
	InviteeID uuid.UUID `json:"invitee_id" db:"invitee_id"`
	Status    string    `json:"status" db:"status"` // pending, accepted, declined
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EventDiscussion for event posts/discussions
type EventDiscussion struct {
	ID        uuid.UUID `json:"id" db:"id"`
	EventID   uuid.UUID `json:"event_id" db:"event_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	MediaURLs StringArray `json:"media_urls" db:"media_urls"`
	IsPinned  bool      `json:"is_pinned" db:"is_pinned"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EventReminder for reminders
type EventReminder struct {
	ID         uuid.UUID `json:"id" db:"id"`
	EventID    uuid.UUID `json:"event_id" db:"event_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	RemindAt   time.Time `json:"remind_at" db:"remind_at"`
	IsSent     bool      `json:"is_sent" db:"is_sent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// API Request/Response Models

type CreateEventRequest struct {
	Title       string         `json:"title" binding:"required,min=5,max=200"`
	Description string         `json:"description" binding:"required,max=5000"`
	CoverPhoto  *string        `json:"cover_photo"`
	Type        EventType      `json:"type" binding:"required,oneof=in_person online"`
	Category    EventCategory  `json:"category" binding:"required"`
	Privacy     EventPrivacy   `json:"privacy" binding:"required,oneof=public private friends"`
	
	// Location (required for in_person)
	LocationName *string  `json:"location_name"`
	Address      *string  `json:"address"`
	City         *string  `json:"city"`
	Country      *string  `json:"country"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	OnlineLink   *string  `json:"online_link"` // Required for online events
	
	// Time
	StartTime time.Time  `json:"start_time" binding:"required"`
	EndTime   *time.Time `json:"end_time"`
	Timezone  string     `json:"timezone" binding:"required"`
	
	// Recurring
	IsRecurring       bool       `json:"is_recurring"`
	RecurrenceRule    *string    `json:"recurrence_rule"`
	RecurrenceEndDate *time.Time `json:"recurrence_end_date"`
	
	// Settings
	AllowGuestInvites bool `json:"allow_guest_invites"`
	RequireApproval   bool `json:"require_approval"`
	MaxAttendees      *int `json:"max_attendees"`
	
	// Co-hosts
	CoHosts []uuid.UUID `json:"co_hosts"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title,omitempty" binding:"omitempty,min=5,max=200"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=5000"`
	CoverPhoto  *string    `json:"cover_photo,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	OnlineLink  *string    `json:"online_link,omitempty"`
}

type RSVPRequest struct {
	Status     RSVPStatus `json:"status" binding:"required,oneof=going interested not_going"`
	GuestCount int        `json:"guest_count" binding:"min=0,max=10"`
}

type InviteUsersRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required,min=1"`
}

type EventListResponse struct {
	Events []*Event `json:"events"`
	Total  int      `json:"total"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
}

type EventStats struct {
	EventID         uuid.UUID `json:"event_id"`
	GoingCount      int       `json:"going_count"`
	InterestedCount int       `json:"interested_count"`
	NotGoingCount   int       `json:"not_going_count"`
	ViewCount       int       `json:"view_count"`
	CheckedInCount  int       `json:"checked_in_count"`
	DiscussionCount int       `json:"discussion_count"`
}

// Helper types for PostgreSQL array support

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return json.Marshal(a)
}

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
