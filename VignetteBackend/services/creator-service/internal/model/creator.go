package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Creator account types
const (
	AccountTypePersonal    AccountType = "personal"
	AccountTypeBusiness    AccountType = "business"
	AccountTypeCreator     AccountType = "creator"
)

// Creator badges
const (
	BadgeVerified     CreatorBadge = "verified"      // Verified account
	BadgePartner      CreatorBadge = "partner"       // Official partner
	BadgeTopCreator   CreatorBadge = "top_creator"   // Top creator program
	BadgeTrendSetter  CreatorBadge = "trendsetter"   // Trendsetter
	BadgeRising       CreatorBadge = "rising"        // Rising star
)

// Monetization status
const (
	MonetizationPending  MonetizationStatus = "pending"
	MonetizationApproved MonetizationStatus = "approved"
	MonetizationRejected MonetizationStatus = "rejected"
	MonetizationSuspended MonetizationStatus = "suspended"
)

type AccountType string
type CreatorBadge string
type MonetizationStatus string

// CreatorProfile represents a creator's professional account
type CreatorProfile struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	UserID          uuid.UUID     `json:"user_id" db:"user_id"`
	AccountType     AccountType   `json:"account_type" db:"account_type"`
	
	// Creator info
	DisplayName     string        `json:"display_name" db:"display_name"`
	Bio             string        `json:"bio" db:"bio"`
	Category        string        `json:"category" db:"category"` // Fashion, Food, Tech, etc.
	Badges          BadgeArray    `json:"badges" db:"badges"`
	
	// Contact
	Email           *string       `json:"email,omitempty" db:"email"`
	Phone           *string       `json:"phone,omitempty" db:"phone"`
	Website         *string       `json:"website,omitempty" db:"website"`
	
	// Monetization
	MonetizationEnabled bool                `json:"monetization_enabled" db:"monetization_enabled"`
	MonetizationStatus  MonetizationStatus  `json:"monetization_status" db:"monetization_status"`
	
	// Timestamps
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

// CreatorAnalytics represents analytics for a creator
type CreatorAnalytics struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Date          time.Time `json:"date" db:"date"`
	
	// Followers
	FollowersCount     int `json:"followers_count" db:"followers_count"`
	FollowersGained    int `json:"followers_gained" db:"followers_gained"`
	FollowersLost      int `json:"followers_lost" db:"followers_lost"`
	
	// Engagement
	TotalLikes         int `json:"total_likes" db:"total_likes"`
	TotalComments      int `json:"total_comments" db:"total_comments"`
	TotalShares        int `json:"total_shares" db:"total_shares"`
	TotalViews         int `json:"total_views" db:"total_views"`
	EngagementRate     float64 `json:"engagement_rate" db:"engagement_rate"`
	
	// Content
	PostsCount         int `json:"posts_count" db:"posts_count"`
	TakesCount         int `json:"takes_count" db:"takes_count"`
	StoriesCount       int `json:"stories_count" db:"stories_count"`
	
	// Reach
	AccountsReached    int `json:"accounts_reached" db:"accounts_reached"`
	AccountsEngaged    int `json:"accounts_engaged" db:"accounts_engaged"`
	
	// Demographics (stored as JSONB)
	AgeGenderBreakdown JSONBData `json:"age_gender_breakdown" db:"age_gender_breakdown"`
	TopLocations       JSONBData `json:"top_locations" db:"top_locations"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ContentInsights represents performance of individual content
type ContentInsights struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ContentID     uuid.UUID `json:"content_id" db:"content_id"`
	ContentType   string    `json:"content_type" db:"content_type"` // post, take, story
	
	// Performance
	Impressions   int     `json:"impressions" db:"impressions"`
	Reach         int     `json:"reach" db:"reach"`
	Likes         int     `json:"likes" db:"likes"`
	Comments      int     `json:"comments" db:"comments"`
	Shares        int     `json:"shares" db:"shares"`
	Saves         int     `json:"saves" db:"saves"`
	Engagement    int     `json:"engagement" db:"engagement"`
	EngagementRate float64 `json:"engagement_rate" db:"engagement_rate"`
	
	// Audience
	FromHome      int `json:"from_home" db:"from_home"`
	FromExplore   int `json:"from_explore" db:"from_explore"`
	FromProfile   int `json:"from_profile" db:"from_profile"`
	FromHashtags  int `json:"from_hashtags" db:"from_hashtags"`
	FromOther     int `json:"from_other" db:"from_other"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// MonetizationApplication for creator monetization
type MonetizationApplication struct {
	ID             uuid.UUID          `json:"id" db:"id"`
	UserID         uuid.UUID          `json:"user_id" db:"user_id"`
	Status         MonetizationStatus `json:"status" db:"status"`
	
	// Requirements
	FollowersCount int  `json:"followers_count" db:"followers_count"`
	PostsCount     int  `json:"posts_count" db:"posts_count"`
	MeetsRequirements bool `json:"meets_requirements" db:"meets_requirements"`
	
	// Tax info
	TaxID          *string `json:"tax_id,omitempty" db:"tax_id"`
	PayoutMethod   *string `json:"payout_method,omitempty" db:"payout_method"`
	
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewedBy     *uuid.UUID `json:"reviewed_by,omitempty" db:"reviewed_by"`
	RejectionReason *string   `json:"rejection_reason,omitempty" db:"rejection_reason"`
	
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CreatorEarnings tracks creator earnings
type CreatorEarnings struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Month        time.Time `json:"month" db:"month"`
	
	// Revenue streams
	AdsRevenue         float64 `json:"ads_revenue" db:"ads_revenue"`
	TipsRevenue        float64 `json:"tips_revenue" db:"tips_revenue"`
	BrandDealsRevenue  float64 `json:"brand_deals_revenue" db:"brand_deals_revenue"`
	OtherRevenue       float64 `json:"other_revenue" db:"other_revenue"`
	TotalRevenue       float64 `json:"total_revenue" db:"total_revenue"`
	
	// Status
	IsPaid         bool       `json:"is_paid" db:"is_paid"`
	PaidAt         *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// AudienceInsight for understanding audience
type AudienceInsight struct {
	UserID           uuid.UUID `json:"user_id"`
	TopAgeGroups     JSONBData `json:"top_age_groups"`
	GenderBreakdown  JSONBData `json:"gender_breakdown"`
	TopCities        JSONBData `json:"top_cities"`
	TopCountries     JSONBData `json:"top_countries"`
	FollowerGrowth   []DailyGrowth `json:"follower_growth"`
	PeakActivityHours JSONBData `json:"peak_activity_hours"`
}

type DailyGrowth struct {
	Date     string `json:"date"`
	Gained   int    `json:"gained"`
	Lost     int    `json:"lost"`
	NetGrowth int   `json:"net_growth"`
}

// API Request/Response Models

type CreateCreatorProfileRequest struct {
	AccountType AccountType `json:"account_type" binding:"required,oneof=personal business creator"`
	DisplayName string      `json:"display_name" binding:"required"`
	Bio         string      `json:"bio" binding:"max=500"`
	Category    string      `json:"category" binding:"required"`
	Email       *string     `json:"email"`
	Website     *string     `json:"website"`
}

type UpdateCreatorProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Category    *string `json:"category,omitempty"`
	Email       *string `json:"email,omitempty"`
	Website     *string `json:"website,omitempty"`
}

type MonetizationApplicationRequest struct {
	TaxID        string `json:"tax_id" binding:"required"`
	PayoutMethod string `json:"payout_method" binding:"required,oneof=bank paypal"`
}

type AnalyticsOverview struct {
	Period           string  `json:"period"` // 7d, 30d, 90d
	FollowersCount   int     `json:"followers_count"`
	FollowersChange  int     `json:"followers_change"`
	EngagementRate   float64 `json:"engagement_rate"`
	AccountsReached  int     `json:"accounts_reached"`
	TotalPosts       int     `json:"total_posts"`
	TotalTakes       int     `json:"total_takes"`
	AverageViews     int     `json:"average_views"`
}

type TopContent struct {
	ContentID      uuid.UUID `json:"content_id"`
	ContentType    string    `json:"content_type"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	Likes          int       `json:"likes"`
	Comments       int       `json:"comments"`
	Shares         int       `json:"shares"`
	Views          int       `json:"views"`
	EngagementRate float64   `json:"engagement_rate"`
}

// Helper types

type BadgeArray []CreatorBadge

func (a *BadgeArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

func (a BadgeArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return json.Marshal(a)
}

type JSONBData map[string]interface{}

func (j *JSONBData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONBData) Value() (driver.Value, error) {
	return json.Marshal(j)
}
