package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Stream status constants
const (
	StatusScheduled StreamStatus = "scheduled"
	StatusLive      StreamStatus = "live"
	StatusEnded     StreamStatus = "ended"
	StatusCancelled StreamStatus = "cancelled"
)

// Stream quality levels (YouTube-quality!)
const (
	Quality144p  StreamQuality = "144p"
	Quality240p  StreamQuality = "240p"
	Quality360p  StreamQuality = "360p"
	Quality480p  StreamQuality = "480p"
	Quality720p  StreamQuality = "720p"   // HD
	Quality1080p StreamQuality = "1080p"  // Full HD
	Quality1440p StreamQuality = "1440p"  // 2K
	Quality2160p StreamQuality = "2160p"  // 4K
)

// Follower thresholds to go live
const (
	MinFollowersVignette = 100  // Vignette: 100 followers to go live
	MinFriendsEntativa  = 50   // Entativa: 50 friends to go live
)

type StreamStatus string
type StreamQuality string

// LiveStream represents a live streaming session
type LiveStream struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	StreamerID  uuid.UUID     `json:"streamer_id" db:"streamer_id"`
	Title       string        `json:"title" db:"title"`
	Description string        `json:"description" db:"description"`
	ThumbnailURL *string      `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	
	// Stream configuration
	Status      StreamStatus  `json:"status" db:"status"`
	Quality     StreamQuality `json:"quality" db:"quality"`
	IsPrivate   bool          `json:"is_private" db:"is_private"`
	Category    string        `json:"category" db:"category"`
	Tags        StringArray   `json:"tags" db:"tags"`
	
	// Technical details
	StreamKey   string        `json:"-" db:"stream_key"` // Secret, not exposed in JSON
	RTMPUrl     string        `json:"rtmp_url" db:"rtmp_url"`
	HLSUrl      string        `json:"hls_url" db:"hls_url"`
	WebRTCUrl   string        `json:"webrtc_url" db:"webrtc_url"`
	
	// Analytics
	ViewerCount     int       `json:"viewer_count" db:"viewer_count"`
	PeakViewers     int       `json:"peak_viewers" db:"peak_viewers"`
	TotalViews      int       `json:"total_views" db:"total_views"`
	LikesCount      int       `json:"likes_count" db:"likes_count"`
	CommentsCount   int       `json:"comments_count" db:"comments_count"`
	SharesCount     int       `json:"shares_count" db:"shares_count"`
	
	// Recording
	RecordStream    bool      `json:"record_stream" db:"record_stream"`
	RecordingURL    *string   `json:"recording_url,omitempty" db:"recording_url"`
	
	// Timestamps
	ScheduledFor    *time.Time `json:"scheduled_for,omitempty" db:"scheduled_for"`
	StartedAt       *time.Time `json:"started_at,omitempty" db:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	Duration        int        `json:"duration" db:"duration"` // Seconds
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// StreamViewer represents a viewer watching a stream
type StreamViewer struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	StreamID   uuid.UUID  `json:"stream_id" db:"stream_id"`
	ViewerID   uuid.UUID  `json:"viewer_id" db:"viewer_id"`
	JoinedAt   time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt     *time.Time `json:"left_at,omitempty" db:"left_at"`
	WatchTime  int        `json:"watch_time" db:"watch_time"` // Seconds
	IsActive   bool       `json:"is_active" db:"is_active"`
}

// StreamComment represents a comment during live stream
type StreamComment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StreamID  uuid.UUID `json:"stream_id" db:"stream_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	IsPinned  bool      `json:"is_pinned" db:"is_pinned"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// StreamReaction represents real-time reactions
type StreamReaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StreamID  uuid.UUID `json:"stream_id" db:"stream_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"` // like, love, fire, etc.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// StreamAnalytics for detailed stats
type StreamAnalytics struct {
	StreamID            uuid.UUID `json:"stream_id"`
	TotalViews          int       `json:"total_views"`
	UniqueViewers       int       `json:"unique_viewers"`
	PeakViewers         int       `json:"peak_viewers"`
	AverageViewers      float64   `json:"average_viewers"`
	AverageWatchTime    int       `json:"average_watch_time"`
	CommentsCount       int       `json:"comments_count"`
	ReactionsCount      int       `json:"reactions_count"`
	SharesCount         int       `json:"shares_count"`
	ChatMessagesPerMin  float64   `json:"chat_messages_per_min"`
}

// API Request/Response Models

type CreateStreamRequest struct {
	Title        string         `json:"title" binding:"required,min=5,max=200"`
	Description  string         `json:"description" binding:"max=1000"`
	ThumbnailURL *string        `json:"thumbnail_url"`
	Quality      StreamQuality  `json:"quality" binding:"required"`
	IsPrivate    bool           `json:"is_private"`
	Category     string         `json:"category" binding:"required"`
	Tags         []string       `json:"tags"`
	RecordStream bool           `json:"record_stream"`
	ScheduledFor *time.Time     `json:"scheduled_for"`
}

type UpdateStreamRequest struct {
	Title        *string   `json:"title,omitempty" binding:"omitempty,min=5,max=200"`
	Description  *string   `json:"description,omitempty" binding:"omitempty,max=1000"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Category     *string   `json:"category,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

type StartStreamResponse struct {
	StreamID    uuid.UUID `json:"stream_id"`
	StreamKey   string    `json:"stream_key"`
	RTMPUrl     string    `json:"rtmp_url"`
	HLSUrl      string    `json:"hls_url"`
	WebRTCUrl   string    `json:"webrtc_url"`
	Message     string    `json:"message"`
}

type StreamCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

type StreamReactionRequest struct {
	Type string `json:"type" binding:"required,oneof=like love fire clap wow"`
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

// StreamEligibility check result
type StreamEligibility struct {
	Eligible     bool   `json:"eligible"`
	Reason       string `json:"reason,omitempty"`
	FollowerCount int   `json:"follower_count"`
	Required     int    `json:"required"`
}
