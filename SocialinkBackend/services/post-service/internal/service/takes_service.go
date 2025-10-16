package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/repository"
	"socialink/post-service/pkg/kafka"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TakesService struct {
	takesRepo    repository.TakesRepository
	bttRepo      repository.BTTRepository
	templateRepo repository.TemplateRepository
	trendRepo    repository.TrendRepository
	redis        *redis.Client
	kafka        *kafka.Producer
}

func NewTakesService(
	takesRepo repository.TakesRepository,
	bttRepo repository.BTTRepository,
	templateRepo repository.TemplateRepository,
	trendRepo repository.TrendRepository,
	redis *redis.Client,
	kafka *kafka.Producer,
) *TakesService {
	return &TakesService{
		takesRepo:    takesRepo,
		bttRepo:      bttRepo,
		templateRepo: templateRepo,
		trendRepo:    trendRepo,
		redis:        redis,
		kafka:        kafka,
	}
}

// CreateTake creates a new Take
func (s *TakesService) CreateTake(ctx context.Context, userID uuid.UUID, req *model.CreateTakeRequest) (*model.Take, error) {
	// Extract hashtags from caption
	hashtags := s.extractHashtags(req.Caption)
	if len(req.Hashtags) > 0 {
		hashtags = append(hashtags, req.Hashtags...)
		hashtags = s.deduplicateStrings(hashtags)
	}

	// Handle trend participation
	var trendID *uuid.UUID
	if req.TrendKeyword != nil && *req.TrendKeyword != "" {
		trend, err := s.JoinOrCreateTrend(ctx, userID, *req.TrendKeyword, req.MediaID)
		if err == nil && trend != nil {
			trendID = &trend.ID
		}
	}

	// Create Take
	take := &model.Take{
		ID:              uuid.New(),
		UserID:          userID,
		Caption:         req.Caption,
		MediaID:         req.MediaID,
		AudioTrackID:    req.AudioTrackID,
		Duration:        0, // Will be set from media metadata
		ThumbnailURL:    "", // Will be set from media service
		Hashtags:        hashtags,
		FilterUsed:      req.FilterUsed,
		Location:        req.Location,
		TaggedUserIDs:   req.TaggedUserIDs,
		TemplateID:      req.TemplateID,
		TrendID:         trendID,
		HasBTT:          false,
		ViewsCount:      0,
		LikesCount:      0,
		CommentsCount:   0,
		SharesCount:     0,
		SavesCount:      0,
		RemixCount:      0,
		CommentsEnabled: req.CommentsEnabled,
		RemixEnabled:    req.RemixEnabled,
		IsSponsored:     false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.takesRepo.Create(ctx, take); err != nil {
		return nil, fmt.Errorf("failed to create Take: %w", err)
	}

	// Increment template usage if used
	if req.TemplateID != nil {
		s.templateRepo.IncrementUsage(ctx, *req.TemplateID)
	}

	// Publish event
	s.publishTakeCreatedEvent(take)

	// Invalidate feed cache
	s.invalidateFeedCache(ctx, userID)

	return take, nil
}

// CreateBTT creates Behind-the-Takes content
func (s *TakesService) CreateBTT(ctx context.Context, takeID, userID uuid.UUID, req *model.CreateBTTRequest) (*model.BehindTheTakes, error) {
	// Verify Take exists and belongs to user
	take, err := s.takesRepo.GetByID(ctx, takeID)
	if err != nil {
		return nil, fmt.Errorf("Take not found")
	}

	if take.UserID != userID {
		return nil, fmt.Errorf("permission denied: not the Take owner")
	}

	// Create BTT
	btt := &model.BehindTheTakes{
		ID:          uuid.New(),
		TakeID:      takeID,
		UserID:      userID,
		MediaIDs:    req.MediaIDs,
		Description: req.Description,
		Steps:       req.Steps,
		Equipment:   req.Equipment,
		Software:    req.Software,
		Tips:        req.Tips,
		ViewsCount:  0,
		LikesCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.bttRepo.Create(ctx, btt); err != nil {
		return nil, fmt.Errorf("failed to create BTT: %w", err)
	}

	// Mark Take as having BTT
	take.HasBTT = true
	take.UpdatedAt = time.Now()
	s.takesRepo.Update(ctx, take)

	// Publish event
	s.publishBTTCreatedEvent(btt, takeID)

	return btt, nil
}

// CreateTemplate creates a template from a Take
func (s *TakesService) CreateTemplate(ctx context.Context, takeID, userID uuid.UUID, req *model.CreateTemplateRequest) (*model.TakeTemplate, error) {
	// Verify Take exists and user can create template
	take, err := s.takesRepo.GetByID(ctx, takeID)
	if err != nil {
		return nil, fmt.Errorf("Take not found")
	}

	// Only owner or if remix_enabled
	if take.UserID != userID && !take.RemixEnabled {
		return nil, fmt.Errorf("permission denied: remixing not allowed")
	}

	// Create template
	template := &model.TakeTemplate{
		ID:             uuid.New(),
		OriginalTakeID: takeID,
		CreatorID:      take.UserID, // Original creator
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		ThumbnailURL:   take.ThumbnailURL,
		AudioTrackID:   take.AudioTrackID,
		Effects:        req.Effects,
		Transitions:    req.Transitions,
		TimingCues:     req.TimingCues,
		UsageCount:     0,
		IsPublic:       req.IsPublic,
		IsFeatured:     false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	// Increment remix count on original Take
	s.takesRepo.IncrementRemixes(ctx, takeID)

	// Publish event
	s.publishTemplateCreatedEvent(template)

	return template, nil
}

// JoinOrCreateTrend joins an existing trend or creates a new one
func (s *TakesService) JoinOrCreateTrend(ctx context.Context, userID uuid.UUID, keyword string, takeID uuid.UUID) (*model.TakeTrend, error) {
	// Normalize keyword
	keyword = strings.ToLower(strings.TrimSpace(keyword))

	// Check if trend exists
	trend, err := s.trendRepo.GetByKeyword(ctx, keyword)
	if err != nil {
		return nil, err
	}

	if trend != nil {
		// Trend exists - increment participants
		s.trendRepo.IncrementParticipants(ctx, trend.ID)
		return trend, nil
	}

	// Create new trend - user becomes originator
	take, err := s.takesRepo.GetByID(ctx, takeID)
	if err != nil {
		return nil, fmt.Errorf("Take not found")
	}

	trend = &model.TakeTrend{
		ID:               uuid.New(),
		Keyword:          keyword,
		OriginatorID:     userID, // Deep-link to originator
		OriginTakeID:     takeID, // Deep-link to origin Take
		DisplayName:      strings.Title(keyword),
		Description:      fmt.Sprintf("Trend started by creator"),
		Category:         "Challenge",
		ThumbnailURL:     take.ThumbnailURL,
		AudioTrackID:     take.AudioTrackID,
		ParticipantCount: 1,
		ViewsCount:       0,
		IsActive:         true,
		IsFeatured:       false,
		StartedAt:        time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.trendRepo.Create(ctx, trend); err != nil {
		return nil, fmt.Errorf("failed to create trend: %w", err)
	}

	// Publish event
	s.publishTrendCreatedEvent(trend)

	return trend, nil
}

// GetTake retrieves a Take by ID
func (s *TakesService) GetTake(ctx context.Context, takeID uuid.UUID) (*model.Take, error) {
	// Try cache
	if take, err := s.getTakeFromCache(ctx, takeID); err == nil && take != nil {
		return take, nil
	}

	// Get from database
	take, err := s.takesRepo.GetByID(ctx, takeID)
	if err != nil {
		return nil, err
	}

	// Cache it
	s.cacheTake(ctx, take)

	// Increment views
	go s.takesRepo.IncrementViews(context.Background(), takeID)

	return take, nil
}

// GetBTT retrieves Behind-the-Takes content
func (s *TakesService) GetBTT(ctx context.Context, takeID uuid.UUID) (*model.BehindTheTakes, error) {
	btt, err := s.bttRepo.GetByTakeID(ctx, takeID)
	if err != nil {
		return nil, err
	}

	if btt == nil {
		return nil, fmt.Errorf("no Behind-the-Takes content for this Take")
	}

	// Increment views
	go s.bttRepo.IncrementViews(context.Background(), btt.ID)

	return btt, nil
}

// GetTrendingTakes retrieves trending Takes
func (s *TakesService) GetTrendingTakes(ctx context.Context, limit int) ([]model.Take, error) {
	// Check cache
	cacheKey := fmt.Sprintf("takes:trending:%d", limit)
	if cached, err := s.getTrendingFromCache(ctx, cacheKey); err == nil && len(cached) > 0 {
		return cached, nil
	}

	// Get from database
	takes, err := s.takesRepo.GetTrending(ctx, limit, 48*time.Hour)
	if err != nil {
		return nil, err
	}

	// Cache for 10 minutes
	s.cacheTrending(ctx, cacheKey, takes, 10*time.Minute)

	return takes, nil
}

// GetTrendingBTT retrieves trending Behind-the-Takes
func (s *TakesService) GetTrendingBTT(ctx context.Context, limit int) ([]model.BehindTheTakes, error) {
	return s.bttRepo.GetTrending(ctx, limit)
}

// GetActiveTrends retrieves active trends
func (s *TakesService) GetActiveTrends(ctx context.Context, limit int) ([]model.TakeTrend, error) {
	return s.trendRepo.GetActive(ctx, limit)
}

// GetTrendTakes retrieves Takes for a specific trend (deep-linked)
func (s *TakesService) GetTrendTakes(ctx context.Context, trendID uuid.UUID, limit, offset int) ([]model.Take, *model.TakeTrend, error) {
	// Get trend info (includes originator info)
	trend, err := s.trendRepo.GetByID(ctx, trendID)
	if err != nil {
		return nil, nil, err
	}

	// Get Takes participating in trend
	takes, err := s.takesRepo.GetByTrendID(ctx, trendID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	// Increment trend views
	go s.trendRepo.IncrementViews(context.Background(), trendID)

	return takes, trend, nil
}

// GetTemplates retrieves templates
func (s *TakesService) GetTemplates(ctx context.Context, category *string, limit, offset int) ([]model.TakeTemplate, error) {
	if category != nil {
		return s.templateRepo.GetByCategory(ctx, *category, limit, offset)
	}
	return s.templateRepo.GetTrending(ctx, limit)
}

// Helper methods

func (s *TakesService) extractHashtags(caption string) []string {
	var hashtags []string
	words := strings.Fields(caption)
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			hashtag := strings.TrimPrefix(word, "#")
			hashtag = strings.ToLower(hashtag)
			hashtag = strings.TrimRight(hashtag, ".,!?;:")
			if len(hashtag) > 0 {
				hashtags = append(hashtags, hashtag)
			}
		}
	}
	
	return s.deduplicateStrings(hashtags)
}

func (s *TakesService) deduplicateStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, str := range strs {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}
	
	return result
}

// Cache methods
func (s *TakesService) cacheTake(ctx context.Context, take *model.Take) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("take:%s", take.ID.String())
	data, _ := json.Marshal(take)
	s.redis.Set(ctx, key, data, 1*time.Hour)
}

func (s *TakesService) getTakeFromCache(ctx context.Context, takeID uuid.UUID) (*model.Take, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	key := fmt.Sprintf("take:%s", takeID.String())
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var take model.Take
	if err := json.Unmarshal(data, &take); err != nil {
		return nil, err
	}

	return &take, nil
}

func (s *TakesService) invalidateFeedCache(ctx context.Context, userID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("takes:feed:%s", userID.String())
	s.redis.Del(ctx, key)
}

func (s *TakesService) cacheTrending(ctx context.Context, key string, takes []model.Take, ttl time.Duration) {
	if s.redis == nil {
		return
	}

	data, _ := json.Marshal(takes)
	s.redis.Set(ctx, key, data, ttl)
}

func (s *TakesService) getTrendingFromCache(ctx context.Context, key string) ([]model.Take, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var takes []model.Take
	if err := json.Unmarshal(data, &takes); err != nil {
		return nil, err
	}

	return takes, nil
}

// Kafka events
func (s *TakesService) publishTakeCreatedEvent(take *model.Take) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":  "take.created",
		"take_id":     take.ID.String(),
		"user_id":     take.UserID.String(),
		"trend_id":    take.TrendID,
		"template_id": take.TemplateID,
		"hashtags":    take.Hashtags,
		"created_at":  take.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "takes-events", take.ID.String(), event)
}

func (s *TakesService) publishBTTCreatedEvent(btt *model.BehindTheTakes, takeID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "btt.created",
		"btt_id":     btt.ID.String(),
		"take_id":    takeID.String(),
		"user_id":    btt.UserID.String(),
		"created_at": btt.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "takes-events", btt.ID.String(), event)
}

func (s *TakesService) publishTemplateCreatedEvent(template *model.TakeTemplate) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":     "template.created",
		"template_id":    template.ID.String(),
		"creator_id":     template.CreatorID.String(),
		"original_take":  template.OriginalTakeID.String(),
		"category":       template.Category,
		"created_at":     template.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "takes-events", template.ID.String(), event)
}

func (s *TakesService) publishTrendCreatedEvent(trend *model.TakeTrend) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":    "trend.created",
		"trend_id":      trend.ID.String(),
		"keyword":       trend.Keyword,
		"originator_id": trend.OriginatorID.String(), // Deep-link
		"origin_take":   trend.OriginTakeID.String(),
		"created_at":    trend.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "takes-events", trend.ID.String(), event)
}
