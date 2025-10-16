package service

import (
	"context"
	"fmt"
	"time"

	"github.com/entativa/vignette/creator-service/internal/model"
	"github.com/entativa/vignette/creator-service/internal/repository"
	"github.com/google/uuid"
)

// Monetization requirements
const (
	MinFollowersForMonetization = 10000 // 10K followers
	MinPostsForMonetization     = 100   // 100 posts
)

type CreatorService struct {
	creatorRepo *repository.CreatorRepository
	kafka       *KafkaProducer
	redis       *RedisClient
}

func NewCreatorService(
	creatorRepo *repository.CreatorRepository,
	kafka *KafkaProducer,
	redis *RedisClient,
) *CreatorService {
	return &CreatorService{
		creatorRepo: creatorRepo,
		kafka:       kafka,
		redis:       redis,
	}
}

// CreateCreatorProfile creates a new creator profile
func (s *CreatorService) CreateCreatorProfile(ctx context.Context, req *model.CreateCreatorProfileRequest, userID uuid.UUID) (*model.CreatorProfile, error) {
	// Check if profile already exists
	existing, _ := s.creatorRepo.GetByUserID(ctx, userID)
	if existing != nil {
		return nil, fmt.Errorf("creator profile already exists")
	}

	profile := &model.CreatorProfile{
		ID:                  uuid.New(),
		UserID:              userID,
		AccountType:         req.AccountType,
		DisplayName:         req.DisplayName,
		Bio:                 req.Bio,
		Category:            req.Category,
		Badges:              []model.CreatorBadge{},
		Email:               req.Email,
		Website:             req.Website,
		MonetizationEnabled: false,
		MonetizationStatus:  model.MonetizationPending,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	err := s.creatorRepo.Create(ctx, profile)
	if err != nil {
		return nil, err
	}

	// Publish event
	s.kafka.PublishCreatorProfileCreated(userID)

	return profile, nil
}

// GetCreatorProfile gets creator profile
func (s *CreatorService) GetCreatorProfile(ctx context.Context, userID uuid.UUID) (*model.CreatorProfile, error) {
	return s.creatorRepo.GetByUserID(ctx, userID)
}

// UpdateCreatorProfile updates profile
func (s *CreatorService) UpdateCreatorProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateCreatorProfileRequest) error {
	profile, err := s.creatorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Update fields
	if req.DisplayName != nil {
		profile.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		profile.Bio = *req.Bio
	}
	if req.Category != nil {
		profile.Category = *req.Category
	}
	if req.Email != nil {
		profile.Email = req.Email
	}
	if req.Website != nil {
		profile.Website = req.Website
	}

	profile.UpdatedAt = time.Now()

	return s.creatorRepo.Update(ctx, profile)
}

// GetAnalyticsOverview gets analytics overview
func (s *CreatorService) GetAnalyticsOverview(ctx context.Context, userID uuid.UUID, period string) (*model.AnalyticsOverview, error) {
	// Calculate date range
	var startDate time.Time
	now := time.Now()

	switch period {
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "90d":
		startDate = now.AddDate(0, 0, -90)
	default:
		startDate = now.AddDate(0, 0, -30) // Default to 30 days
	}

	// Get analytics
	analytics, err := s.creatorRepo.GetAnalytics(ctx, userID, startDate, now)
	if err != nil {
		return nil, err
	}

	if len(analytics) == 0 {
		return &model.AnalyticsOverview{
			Period: period,
		}, nil
	}

	// Calculate aggregates
	latest := analytics[0]
	oldest := analytics[len(analytics)-1]

	totalLikes := 0
	totalComments := 0
	totalShares := 0
	totalViews := 0
	totalPosts := 0
	totalTakes := 0
	totalReached := 0

	for _, a := range analytics {
		totalLikes += a.TotalLikes
		totalComments += a.TotalComments
		totalShares += a.TotalShares
		totalViews += a.TotalViews
		totalPosts += a.PostsCount
		totalTakes += a.TakesCount
		totalReached += a.AccountsReached
	}

	engagementRate := 0.0
	if totalReached > 0 {
		engagementRate = float64(totalLikes+totalComments+totalShares) / float64(totalReached) * 100
	}

	avgViews := 0
	if totalPosts+totalTakes > 0 {
		avgViews = totalViews / (totalPosts + totalTakes)
	}

	overview := &model.AnalyticsOverview{
		Period:          period,
		FollowersCount:  latest.FollowersCount,
		FollowersChange: latest.FollowersCount - oldest.FollowersCount,
		EngagementRate:  engagementRate,
		AccountsReached: totalReached,
		TotalPosts:      totalPosts,
		TotalTakes:      totalTakes,
		AverageViews:    avgViews,
	}

	return overview, nil
}

// GetAudienceInsights gets audience insights
func (s *CreatorService) GetAudienceInsights(ctx context.Context, userID uuid.UUID, days int) (*model.AudienceInsight, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	analytics, err := s.creatorRepo.GetAnalytics(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(analytics) == 0 {
		return &model.AudienceInsight{UserID: userID}, nil
	}

	// Aggregate demographics
	latest := analytics[0]

	// Build follower growth
	growth := []model.DailyGrowth{}
	for _, a := range analytics {
		growth = append(growth, model.DailyGrowth{
			Date:      a.Date.Format("2006-01-02"),
			Gained:    a.FollowersGained,
			Lost:      a.FollowersLost,
			NetGrowth: a.FollowersGained - a.FollowersLost,
		})
	}

	insight := &model.AudienceInsight{
		UserID:           userID,
		TopAgeGroups:     latest.AgeGenderBreakdown,
		GenderBreakdown:  latest.AgeGenderBreakdown,
		TopCities:        latest.TopLocations,
		TopCountries:     latest.TopLocations,
		FollowerGrowth:   growth,
		PeakActivityHours: model.JSONBData{}, // Placeholder
	}

	return insight, nil
}

// GetTopContent gets top performing content
func (s *CreatorService) GetTopContent(ctx context.Context, userID uuid.UUID, contentType string, limit int) ([]*model.TopContent, error) {
	insights, err := s.creatorRepo.GetTopContent(ctx, userID, contentType, limit)
	if err != nil {
		return nil, err
	}

	topContent := []*model.TopContent{}
	for _, insight := range insights {
		topContent = append(topContent, &model.TopContent{
			ContentID:      insight.ContentID,
			ContentType:    insight.ContentType,
			ThumbnailURL:   "", // To be fetched from content service
			Likes:          insight.Likes,
			Comments:       insight.Comments,
			Shares:         insight.Shares,
			Views:          insight.Reach,
			EngagementRate: insight.EngagementRate,
		})
	}

	return topContent, nil
}

// ApplyForMonetization applies for monetization
func (s *CreatorService) ApplyForMonetization(ctx context.Context, userID uuid.UUID, req *model.MonetizationApplicationRequest) error {
	// Get current stats (from user service or cache)
	followersCount := 15000 // Mock
	postsCount := 150       // Mock

	// Check requirements
	meetsRequirements := followersCount >= MinFollowersForMonetization && postsCount >= MinPostsForMonetization

	if !meetsRequirements {
		return fmt.Errorf("does not meet monetization requirements: need %d followers (have %d) and %d posts (have %d)",
			MinFollowersForMonetization, followersCount, MinPostsForMonetization, postsCount)
	}

	// Create application
	// In production: Store in monetization_applications table
	// For now, just update profile status
	err := s.creatorRepo.UpdateMonetizationStatus(ctx, userID, model.MonetizationPending, false)
	if err != nil {
		return err
	}

	// Publish event
	s.kafka.PublishMonetizationApplicationSubmitted(userID)

	return nil
}

// AwardBadge awards a badge to creator
func (s *CreatorService) AwardBadge(ctx context.Context, userID uuid.UUID, badge model.CreatorBadge) error {
	err := s.creatorRepo.AddBadge(ctx, userID, badge)
	if err != nil {
		return err
	}

	// Publish event
	s.kafka.PublishBadgeAwarded(userID, string(badge))

	return nil
}

// RecordContentInsight records insight for content
func (s *CreatorService) RecordContentInsight(ctx context.Context, insight *model.ContentInsights) error {
	return s.creatorRepo.RecordContentInsight(ctx, insight)
}

// Stub types
type KafkaProducer struct{}

func (k *KafkaProducer) PublishCreatorProfileCreated(userID uuid.UUID)                  {}
func (k *KafkaProducer) PublishMonetizationApplicationSubmitted(userID uuid.UUID)       {}
func (k *KafkaProducer) PublishBadgeAwarded(userID uuid.UUID, badge string)             {}

type RedisClient struct{}
