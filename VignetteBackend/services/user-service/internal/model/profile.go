package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Profile represents extended user profile information for Vignette (Instagram-like)
type Profile struct {
	ID                  uuid.UUID         `json:"id" db:"id"`
	UserID              uuid.UUID         `json:"user_id" db:"user_id"`
	Category            *string           `json:"category,omitempty" db:"category"` // personal, creator, business
	CategoryType        *string           `json:"category_type,omitempty" db:"category_type"` // photographer, artist, musician, etc.
	Gender              *string           `json:"gender,omitempty" db:"gender"`
	Pronouns            *string           `json:"pronouns,omitempty" db:"pronouns"`
	LinkInBio           []LinkInBio       `json:"link_in_bio,omitempty" db:"link_in_bio"`
	Highlights          []StoryHighlight  `json:"highlights,omitempty" db:"highlights"`
	PinnedPosts         StringArray       `json:"pinned_posts,omitempty" db:"pinned_posts"`
	ProfileBadges       StringArray       `json:"profile_badges,omitempty" db:"profile_badges"` // verified, creator, etc.
	ContactOptions      *ContactOptions   `json:"contact_options,omitempty" db:"contact_options"`
	CreatorInsights     *CreatorInsights  `json:"creator_insights,omitempty" db:"creator_insights"`
	BusinessInfo        *BusinessInfo     `json:"business_info,omitempty" db:"business_info"`
	SuggestedUsers      StringArray       `json:"suggested_users,omitempty" db:"suggested_users"`
	ProfileViews        int64             `json:"profile_views" db:"profile_views"`
	ProfileViewsEnabled bool              `json:"profile_views_enabled" db:"profile_views_enabled"`
	Availability        *Availability     `json:"availability,omitempty" db:"availability"`
	CreatedAt           time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at" db:"updated_at"`
}

// LinkInBio represents a link in bio feature (similar to Linktree)
type LinkInBio struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Order int    `json:"order"`
}

// StoryHighlight represents saved story highlights
type StoryHighlight struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CoverURL  string    `json:"cover_url"`
	StoryIDs  []string  `json:"story_ids"`
	CreatedAt time.Time `json:"created_at"`
	Order     int       `json:"order"`
}

// ContactOptions for business/creator accounts
type ContactOptions struct {
	Email         *string `json:"email,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	AddressStreet *string `json:"address_street,omitempty"`
	AddressCity   *string `json:"address_city,omitempty"`
	AddressZip    *string `json:"address_zip,omitempty"`
	ShowEmail     bool    `json:"show_email"`
	ShowPhone     bool    `json:"show_phone"`
	ShowAddress   bool    `json:"show_address"`
}

// CreatorInsights for creator accounts
type CreatorInsights struct {
	IsCreatorAccount     bool      `json:"is_creator_account"`
	EnabledDate          time.Time `json:"enabled_date"`
	TotalReach           int64     `json:"total_reach"`
	TotalImpressions     int64     `json:"total_impressions"`
	TotalEngagement      int64     `json:"total_engagement"`
	AverageEngagementRate float64   `json:"average_engagement_rate"`
	TopAudiences         []string  `json:"top_audiences,omitempty"`
}

// BusinessInfo for business accounts
type BusinessInfo struct {
	IsBusinessAccount bool    `json:"is_business_account"`
	BusinessCategory  string  `json:"business_category"`
	BusinessEmail     *string `json:"business_email,omitempty"`
	BusinessPhone     *string `json:"business_phone,omitempty"`
	BusinessAddress   *string `json:"business_address,omitempty"`
	PriceRange        *string `json:"price_range,omitempty"` // $, $$, $$$, $$$$
	Hours             *string `json:"hours,omitempty"`
}

// Availability for professional accounts
type Availability struct {
	Status  string `json:"status"` // available, busy, not_available
	Message string `json:"message,omitempty"`
}

// StringArray custom type for PostgreSQL array handling
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

// JSON type handlers for complex types
func (l LinkInBio) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *LinkInBio) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, l)
}

func (s StoryHighlight) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StoryHighlight) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

func (c ContactOptions) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ContactOptions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

func (c CreatorInsights) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *CreatorInsights) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

func (b BusinessInfo) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BusinessInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, b)
}

func (a Availability) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Availability) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

// DTOs for profile management

// UpdateProfileExtendedRequest for extended profile updates
type UpdateProfileExtendedRequest struct {
	Category     *string `json:"category,omitempty" binding:"omitempty,oneof=personal creator business"`
	CategoryType *string `json:"category_type,omitempty" binding:"omitempty,max=50"`
	Gender       *string `json:"gender,omitempty" binding:"omitempty,max=50"`
	Pronouns     *string `json:"pronouns,omitempty" binding:"omitempty,max=50"`
}

// AddLinkInBioRequest for adding links
type AddLinkInBioRequest struct {
	Title string `json:"title" binding:"required,max=100"`
	URL   string `json:"url" binding:"required,url"`
}

// AddHighlightRequest for adding story highlights
type AddHighlightRequest struct {
	Title    string   `json:"title" binding:"required,max=50"`
	CoverURL string   `json:"cover_url" binding:"required,url"`
	StoryIDs []string `json:"story_ids" binding:"required"`
}

// UpdateContactOptionsRequest for business contact info
type UpdateContactOptionsRequest struct {
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Address     *string `json:"address,omitempty" binding:"omitempty,max=200"`
	ShowEmail   *bool   `json:"show_email,omitempty"`
	ShowPhone   *bool   `json:"show_phone,omitempty"`
	ShowAddress *bool   `json:"show_address,omitempty"`
}

// EnableCreatorAccountRequest for switching to creator account
type EnableCreatorAccountRequest struct {
	Category string `json:"category" binding:"required,max=50"` // photographer, artist, musician, etc.
}

// EnableBusinessAccountRequest for switching to business account
type EnableBusinessAccountRequest struct {
	BusinessCategory string  `json:"business_category" binding:"required,max=100"`
	BusinessEmail    *string `json:"business_email,omitempty" binding:"omitempty,email"`
	BusinessPhone    *string `json:"business_phone,omitempty"`
	BusinessAddress  *string `json:"business_address,omitempty"`
	PriceRange       *string `json:"price_range,omitempty" binding:"omitempty,oneof=$ $$ $$$ $$$$"`
}

// UpdateAvailabilityRequest for availability status
type UpdateAvailabilityRequest struct {
	Status  string  `json:"status" binding:"required,oneof=available busy not_available"`
	Message *string `json:"message,omitempty" binding:"omitempty,max=200"`
}

// ProfileResponse for API responses
type ProfileResponse struct {
	ID                  uuid.UUID        `json:"id"`
	UserID              uuid.UUID        `json:"user_id"`
	User                *UserResponse    `json:"user,omitempty"`
	Category            *string          `json:"category,omitempty"`
	CategoryType        *string          `json:"category_type,omitempty"`
	Gender              *string          `json:"gender,omitempty"`
	Pronouns            *string          `json:"pronouns,omitempty"`
	LinkInBio           []LinkInBio      `json:"link_in_bio,omitempty"`
	Highlights          []StoryHighlight `json:"highlights,omitempty"`
	PinnedPosts         []string         `json:"pinned_posts,omitempty"`
	ProfileBadges       []string         `json:"profile_badges,omitempty"`
	ContactOptions      *ContactOptions  `json:"contact_options,omitempty"`
	CreatorInsights     *CreatorInsights `json:"creator_insights,omitempty"`
	BusinessInfo        *BusinessInfo    `json:"business_info,omitempty"`
	ProfileViews        int64            `json:"profile_views"`
	ProfileViewsEnabled bool             `json:"profile_views_enabled"`
	Availability        *Availability    `json:"availability,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// ToProfileResponse converts Profile to ProfileResponse
func (p *Profile) ToProfileResponse() *ProfileResponse {
	return &ProfileResponse{
		ID:                  p.ID,
		UserID:              p.UserID,
		Category:            p.Category,
		CategoryType:        p.CategoryType,
		Gender:              p.Gender,
		Pronouns:            p.Pronouns,
		LinkInBio:           p.LinkInBio,
		Highlights:          p.Highlights,
		PinnedPosts:         p.PinnedPosts,
		ProfileBadges:       p.ProfileBadges,
		ContactOptions:      p.ContactOptions,
		CreatorInsights:     p.CreatorInsights,
		BusinessInfo:        p.BusinessInfo,
		ProfileViews:        p.ProfileViews,
		ProfileViewsEnabled: p.ProfileViewsEnabled,
		Availability:        p.Availability,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}
