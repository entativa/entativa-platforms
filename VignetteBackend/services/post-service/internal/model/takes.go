package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Take represents a short-form video (Instagram Reels equivalent)
type Take struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	Caption         string         `json:"caption" db:"caption"`
	MediaID         uuid.UUID      `json:"media_id" db:"media_id"` // Single video
	AudioTrackID    *uuid.UUID     `json:"audio_track_id,omitempty" db:"audio_track_id"`
	Duration        float64        `json:"duration" db:"duration"`
	ThumbnailURL    string         `json:"thumbnail_url" db:"thumbnail_url"`
	Hashtags        StringList     `json:"hashtags" db:"hashtags"`
	FilterUsed      *string        `json:"filter_used,omitempty" db:"filter_used"`
	Location        *Location      `json:"location,omitempty" db:"location"`
	TaggedUserIDs   UUIDList       `json:"tagged_user_ids" db:"tagged_user_ids"`
	
	// Takes-specific features
	TemplateID      *uuid.UUID     `json:"template_id,omitempty" db:"template_id"` // If created from template
	TrendID         *uuid.UUID     `json:"trend_id,omitempty" db:"trend_id"` // If part of a trend
	HasBTT          bool           `json:"has_btt" db:"has_btt"` // Has Behind-the-Takes
	
	// Engagement
	ViewsCount      int64          `json:"views_count" db:"views_count"`
	LikesCount      int64          `json:"likes_count" db:"likes_count"`
	CommentsCount   int64          `json:"comments_count" db:"comments_count"`
	SharesCount     int64          `json:"shares_count" db:"shares_count"`
	SavesCount      int64          `json:"saves_count" db:"saves_count"`
	RemixCount      int64          `json:"remix_count" db:"remix_count"` // Times used as template
	
	// Settings
	CommentsEnabled bool           `json:"comments_enabled" db:"comments_enabled"`
	RemixEnabled    bool           `json:"remix_enabled" db:"remix_enabled"`
	
	// Metadata
	IsSponsored     bool           `json:"is_sponsored" db:"is_sponsored"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
}

// BehindTheTakes represents the behind-the-scenes content for a Take
type BehindTheTakes struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	TakeID      uuid.UUID      `json:"take_id" db:"take_id"`
	UserID      uuid.UUID      `json:"user_id" db:"user_id"`
	MediaIDs    MediaIDList    `json:"media_ids" db:"media_ids"` // Multiple BTS media
	Description string         `json:"description" db:"description"`
	Steps       StepsList      `json:"steps" db:"steps"` // Step-by-step breakdown
	Equipment   StringList     `json:"equipment" db:"equipment"` // Gear used
	Software    StringList     `json:"software" db:"software"` // Apps/tools used
	Tips        StringList     `json:"tips" db:"tips"` // Creator tips
	ViewsCount  int64          `json:"views_count" db:"views_count"`
	LikesCount  int64          `json:"likes_count" db:"likes_count"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// TakeTemplate represents a reusable template from a Take
type TakeTemplate struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	OriginalTakeID  uuid.UUID      `json:"original_take_id" db:"original_take_id"`
	CreatorID       uuid.UUID      `json:"creator_id" db:"creator_id"`
	Name            string         `json:"name" db:"name"`
	Description     string         `json:"description" db:"description"`
	Category        string         `json:"category" db:"category"` // Dance, Comedy, Tutorial, etc.
	ThumbnailURL    string         `json:"thumbnail_url" db:"thumbnail_url"`
	AudioTrackID    *uuid.UUID     `json:"audio_track_id,omitempty" db:"audio_track_id"`
	Effects         EffectsList    `json:"effects" db:"effects"` // Visual effects used
	Transitions     StringList     `json:"transitions" db:"transitions"`
	TimingCues      CuesList       `json:"timing_cues" db:"timing_cues"` // Beat markers, etc.
	UsageCount      int64          `json:"usage_count" db:"usage_count"`
	IsPublic        bool           `json:"is_public" db:"is_public"`
	IsFeatured      bool           `json:"is_featured" db:"is_featured"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// TakeTrend represents a trending topic/challenge in Takes
type TakeTrend struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	Keyword         string         `json:"keyword" db:"keyword"` // Case-insensitive unique
	OriginatorID    uuid.UUID      `json:"originator_id" db:"originator_id"` // Who started it
	OriginTakeID    uuid.UUID      `json:"origin_take_id" db:"origin_take_id"` // Original Take
	DisplayName     string         `json:"display_name" db:"display_name"` // Pretty name
	Description     string         `json:"description" db:"description"`
	Category        string         `json:"category" db:"category"`
	ThumbnailURL    string         `json:"thumbnail_url" db:"thumbnail_url"`
	AudioTrackID    *uuid.UUID     `json:"audio_track_id,omitempty" db:"audio_track_id"`
	ParticipantCount int64         `json:"participant_count" db:"participant_count"`
	ViewsCount      int64          `json:"views_count" db:"views_count"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	IsFeatured      bool           `json:"is_featured" db:"is_featured"`
	StartedAt       time.Time      `json:"started_at" db:"started_at"`
	PeakAt          *time.Time     `json:"peak_at,omitempty" db:"peak_at"`
	ExpiresAt       *time.Time     `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// Effect represents a visual effect in a template
type Effect struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"` // filter, transition, sticker, etc.
	Timestamp float64 `json:"timestamp"` // When to apply
	Duration  float64 `json:"duration"`
	Intensity float64 `json:"intensity"`
}

// TimingCue represents a timing marker in a template
type TimingCue struct {
	Timestamp   float64 `json:"timestamp"`
	Type        string  `json:"type"` // beat, cut, transition
	Description string  `json:"description"`
}

// CreationStep represents a step in BTT
type CreationStep struct {
	StepNumber  int      `json:"step_number"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MediaURL    *string  `json:"media_url,omitempty"`
	Duration    *float64 `json:"duration,omitempty"`
}

// Custom JSONB types
type EffectsList []Effect
type CuesList []TimingCue
type StepsList []CreationStep

// Scan/Value implementations for EffectsList
func (e *EffectsList) Scan(value interface{}) error {
	if value == nil {
		*e = []Effect{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		*e = []Effect{}
		return nil
	}
	return json.Unmarshal(bytes, e)
}

func (e EffectsList) Value() (driver.Value, error) {
	if len(e) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(e)
}

// Scan/Value implementations for CuesList
func (c *CuesList) Scan(value interface{}) error {
	if value == nil {
		*c = []TimingCue{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		*c = []TimingCue{}
		return nil
	}
	return json.Unmarshal(bytes, c)
}

func (c CuesList) Value() (driver.Value, error) {
	if len(c) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(c)
}

// Scan/Value implementations for StepsList
func (s *StepsList) Scan(value interface{}) error {
	if value == nil {
		*s = []CreationStep{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		*s = []CreationStep{}
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s StepsList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(s)
}

// DTOs

type CreateTakeRequest struct {
	Caption         string      `json:"caption" binding:"max=2200"`
	MediaID         uuid.UUID   `json:"media_id" binding:"required"`
	AudioTrackID    *uuid.UUID  `json:"audio_track_id,omitempty"`
	Hashtags        []string    `json:"hashtags,omitempty"`
	FilterUsed      *string     `json:"filter_used,omitempty"`
	Location        *Location   `json:"location,omitempty"`
	TaggedUserIDs   []uuid.UUID `json:"tagged_user_ids,omitempty"`
	TemplateID      *uuid.UUID  `json:"template_id,omitempty"`
	TrendKeyword    *string     `json:"trend_keyword,omitempty"` // Case-insensitive
	CommentsEnabled bool        `json:"comments_enabled"`
	RemixEnabled    bool        `json:"remix_enabled"`
}

type CreateBTTRequest struct {
	MediaIDs    []uuid.UUID    `json:"media_ids" binding:"required,min=1"`
	Description string         `json:"description" binding:"required"`
	Steps       []CreationStep `json:"steps" binding:"required,min=1"`
	Equipment   []string       `json:"equipment,omitempty"`
	Software    []string       `json:"software,omitempty"`
	Tips        []string       `json:"tips,omitempty"`
}

type CreateTemplateRequest struct {
	Name         string      `json:"name" binding:"required,max=100"`
	Description  string      `json:"description" binding:"required"`
	Category     string      `json:"category" binding:"required"`
	Effects      []Effect    `json:"effects,omitempty"`
	Transitions  []string    `json:"transitions,omitempty"`
	TimingCues   []TimingCue `json:"timing_cues,omitempty"`
	IsPublic     bool        `json:"is_public"`
}

type JoinTrendRequest struct {
	TrendKeyword string `json:"trend_keyword" binding:"required"`
}

type CreateTrendRequest struct {
	Keyword      string     `json:"keyword" binding:"required,max=50"`
	DisplayName  string     `json:"display_name" binding:"required,max=100"`
	Description  string     `json:"description" binding:"required"`
	Category     string     `json:"category" binding:"required"`
	AudioTrackID *uuid.UUID `json:"audio_track_id,omitempty"`
}

type TakeResponse struct {
	*Take
	User          *UserInfo    `json:"user,omitempty"`
	MediaURL      string       `json:"media_url"`
	ThumbnailURL  string       `json:"thumbnail_url"`
	AudioTrack    *AudioInfo   `json:"audio_track,omitempty"`
	Template      *TemplateInfo `json:"template,omitempty"`
	Trend         *TrendInfo   `json:"trend,omitempty"`
	HasBTT        bool         `json:"has_btt"`
	IsLiked       bool         `json:"is_liked"`
	IsSaved       bool         `json:"is_saved"`
}

type BTTResponse struct {
	*BehindTheTakes
	Take      *Take      `json:"take,omitempty"`
	Creator   *UserInfo  `json:"creator,omitempty"`
	MediaURLs []string   `json:"media_urls"`
}

type TemplateResponse struct {
	*TakeTemplate
	Creator       *UserInfo `json:"creator,omitempty"`
	OriginalTake  *Take     `json:"original_take,omitempty"`
	RecentTakes   []Take    `json:"recent_takes,omitempty"`
}

type TrendResponse struct {
	*TakeTrend
	Originator    *UserInfo `json:"originator,omitempty"`
	OriginTake    *Take     `json:"origin_take,omitempty"`
	TrendingTakes []Take    `json:"trending_takes,omitempty"`
	Rank          int       `json:"rank,omitempty"`
}

type AudioInfo struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Artist   string    `json:"artist"`
	Duration float64   `json:"duration"`
	URL      string    `json:"url"`
}

type TemplateInfo struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ThumbnailURL string    `json:"thumbnail_url"`
	UsageCount   int64     `json:"usage_count"`
}

type TrendInfo struct {
	ID               uuid.UUID `json:"id"`
	Keyword          string    `json:"keyword"`
	DisplayName      string    `json:"display_name"`
	OriginatorID     uuid.UUID `json:"originator_id"`
	OriginatorName   string    `json:"originator_name"`
	ParticipantCount int64     `json:"participant_count"`
	Rank             int       `json:"rank,omitempty"`
}
