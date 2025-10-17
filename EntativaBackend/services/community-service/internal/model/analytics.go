package model

import (
	"time"

	"github.com/google/uuid"
)

// CommunityAnalytics represents analytics for a community
type CommunityAnalytics struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	Date        time.Time `json:"date" db:"date"`
	
	// Growth metrics
	NewMembers      int `json:"new_members" db:"new_members"`
	LeftMembers     int `json:"left_members" db:"left_members"`
	TotalMembers    int `json:"total_members" db:"total_members"`
	
	// Engagement metrics
	NewPosts        int `json:"new_posts" db:"new_posts"`
	TotalComments   int `json:"total_comments" db:"total_comments"`
	TotalLikes      int `json:"total_likes" db:"total_likes"`
	TotalShares     int `json:"total_shares" db:"total_shares"`
	
	// Activity metrics
	ActiveMembers   int     `json:"active_members" db:"active_members"` // Posted or commented
	EngagementRate  float64 `json:"engagement_rate" db:"engagement_rate"` // Active / Total
	AvgPostsPerMember float64 `json:"avg_posts_per_member" db:"avg_posts_per_member"`
	
	// Moderation metrics
	ReportsReceived int `json:"reports_received" db:"reports_received"`
	ActionsToken    int `json:"actions_taken" db:"actions_taken"`
	MembersBanned   int `json:"members_banned" db:"members_banned"`
	
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// AnalyticsSummary for date range
type AnalyticsSummary struct {
	CommunityID uuid.UUID `json:"community_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	
	// Aggregated metrics
	TotalNewMembers   int     `json:"total_new_members"`
	TotalLeftMembers  int     `json:"total_left_members"`
	NetGrowth         int     `json:"net_growth"`
	GrowthRate        float64 `json:"growth_rate"`
	
	TotalPosts        int     `json:"total_posts"`
	TotalComments     int     `json:"total_comments"`
	TotalEngagements  int     `json:"total_engagements"`
	
	AvgEngagementRate float64 `json:"avg_engagement_rate"`
	PeakActivityDate  *time.Time `json:"peak_activity_date,omitempty"`
	
	TotalReports      int     `json:"total_reports"`
	TotalModerationActions int `json:"total_moderation_actions"`
}

// MemberActivity for tracking individual member activity
type MemberActivity struct {
	CommunityID   uuid.UUID `json:"community_id"`
	UserID        uuid.UUID `json:"user_id"`
	PostCount     int       `json:"post_count"`
	CommentCount  int       `json:"comment_count"`
	LikeCount     int       `json:"like_count"`
	LastActive    time.Time `json:"last_active"`
	JoinDate      time.Time `json:"join_date"`
	DaysSinceJoin int       `json:"days_since_join"`
}

// TopContributor represents top contributors
type TopContributor struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	PostCount    int       `json:"post_count"`
	CommentCount int       `json:"comment_count"`
	LikeCount    int       `json:"like_count"`
	TotalScore   int       `json:"total_score"` // Weighted score
}

// CommunityInsights for overview
type CommunityInsights struct {
	CommunityID   uuid.UUID `json:"community_id"`
	
	// Current stats
	TotalMembers  int     `json:"total_members"`
	ActiveMembers int     `json:"active_members"`
	TotalPosts    int     `json:"total_posts"`
	
	// Growth (last 30 days)
	MemberGrowth30d int     `json:"member_growth_30d"`
	GrowthRate30d   float64 `json:"growth_rate_30d"`
	
	// Engagement (last 7 days)
	PostsLast7d     int     `json:"posts_last_7d"`
	CommentsLast7d  int     `json:"comments_last_7d"`
	EngagementRate  float64 `json:"engagement_rate"`
	
	// Top contributors (last 30 days)
	TopContributors []TopContributor `json:"top_contributors"`
	
	// Moderation (last 30 days)
	ReportsLast30d  int `json:"reports_last_30d"`
	ActionsLast30d  int `json:"actions_last_30d"`
}
