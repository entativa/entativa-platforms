package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vignette/post-service/internal/model"
	"vignette/post-service/internal/repository"
	"vignette/post-service/pkg/kafka"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PostService struct {
	postRepo repository.PostRepository
	likeRepo repository.LikeRepository
	commentRepo repository.CommentRepository
	saveRepo repository.SaveRepository
	redis    *redis.Client
	kafka    *kafka.Producer
}

func NewPostService(
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	saveRepo repository.SaveRepository,
	redis *redis.Client,
	kafka *kafka.Producer,
) *PostService {
	return &PostService{
		postRepo:    postRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		saveRepo:    saveRepo,
		redis:       redis,
		kafka:       kafka,
	}
}

// CreatePost creates a new Instagram-style post (media required)
func (s *PostService) CreatePost(ctx context.Context, userID uuid.UUID, req *model.CreatePostRequest) (*model.Post, error) {
	// Validate media is provided (Instagram requires media)
	if len(req.MediaIDs) == 0 {
		return nil, fmt.Errorf("at least one media attachment is required")
	}

	// Extract hashtags from caption
	hashtags := s.extractHashtags(req.Caption)
	if len(req.Hashtags) > 0 {
		hashtags = append(hashtags, req.Hashtags...)
		hashtags = s.deduplicateStrings(hashtags)
	}

	// Determine if carousel (multiple images)
	isCarousel := len(req.MediaIDs) > 1

	// Create post
	post := &model.Post{
		ID:              uuid.New(),
		UserID:          userID,
		Caption:         req.Caption,
		MediaIDs:        req.MediaIDs,
		Location:        req.Location,
		TaggedUserIDs:   req.TaggedUserIDs,
		Hashtags:        hashtags,
		FilterUsed:      req.FilterUsed,
		IsCarousel:      isCarousel,
		LikesCount:      0,
		CommentsCount:   0,
		ViewsCount:      0,
		SavesCount:      0,
		SharesCount:     0,
		IsEdited:        false,
		CommentsEnabled: req.CommentsEnabled,
		LikesVisible:    req.LikesVisible,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Publish event to Kafka
	s.publishPostCreatedEvent(post)

	// Invalidate user's feed cache
	s.invalidateFeedCache(ctx, userID)

	return post, nil
}

// GetPost retrieves a post by ID
func (s *PostService) GetPost(ctx context.Context, postID uuid.UUID) (*model.Post, error) {
	// Try cache first
	if post, err := s.getPostFromCache(ctx, postID); err == nil && post != nil {
		return post, nil
	}

	// Get from database
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Cache it
	s.cachePost(ctx, post)

	// Increment view count asynchronously
	go func() {
		s.postRepo.IncrementViews(context.Background(), postID)
	}()

	return post, nil
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(ctx context.Context, postID, userID uuid.UUID, req *model.UpdatePostRequest) (*model.Post, error) {
	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if post.UserID != userID {
		return nil, fmt.Errorf("permission denied: not the post owner")
	}

	// Update fields
	if req.Caption != nil {
		post.Caption = *req.Caption
		// Re-extract hashtags
		post.Hashtags = s.extractHashtags(post.Caption)
		if req.Hashtags != nil {
			post.Hashtags = append(post.Hashtags, *req.Hashtags...)
			post.Hashtags = s.deduplicateStrings(post.Hashtags)
		}
	}
	if req.Location != nil {
		post.Location = req.Location
	}
	if req.CommentsEnabled != nil {
		post.CommentsEnabled = *req.CommentsEnabled
	}
	if req.LikesVisible != nil {
		post.LikesVisible = *req.LikesVisible
	}

	now := time.Now()
	post.IsEdited = true
	post.EditedAt = &now
	post.UpdatedAt = now

	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Invalidate cache
	s.invalidatePostCache(ctx, postID)

	// Publish event
	s.publishPostUpdatedEvent(post)

	return post, nil
}

// DeletePost soft deletes a post
func (s *PostService) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
	// Get post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Verify ownership
	if post.UserID != userID {
		return fmt.Errorf("permission denied: not the post owner")
	}

	// Delete
	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Invalidate cache
	s.invalidatePostCache(ctx, postID)
	s.invalidateFeedCache(ctx, userID)

	// Publish event
	s.publishPostDeletedEvent(postID, userID)

	return nil
}

// GetUserPosts retrieves posts by a user
func (s *PostService) GetUserPosts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Post, error) {
	return s.postRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetFeed retrieves the user's personalized feed
func (s *PostService) GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]model.Post, *string, error) {
	// Try cache first for first page
	if cursor == "" {
		if cachedPosts, err := s.getFeedFromCache(ctx, userID, limit); err == nil && len(cachedPosts) > 0 {
			var nextCursor *string
			if len(cachedPosts) > limit {
				cachedPosts = cachedPosts[:limit]
				cursorStr := cachedPosts[limit-1].ID.String()
				nextCursor = &cursorStr
			}
			return cachedPosts, nextCursor, nil
		}
	}

	// Get from database
	posts, nextCursor, err := s.postRepo.GetFeed(ctx, userID, cursor, limit)
	if err != nil {
		return nil, nil, err
	}

	// Cache first page
	if cursor == "" && len(posts) > 0 {
		s.cacheFeed(ctx, userID, posts)
	}

	return posts, nextCursor, nil
}

// GetExplorePosts retrieves explore page posts
func (s *PostService) GetExplorePosts(ctx context.Context, limit int) ([]model.Post, error) {
	// Check cache
	cacheKey := fmt.Sprintf("explore:posts:%d", limit)
	if cachedPosts, err := s.getTrendingFromCache(ctx, cacheKey); err == nil && len(cachedPosts) > 0 {
		return cachedPosts, nil
	}

	// Get trending posts for explore
	posts, err := s.postRepo.GetTrendingPosts(ctx, limit, 48*time.Hour)
	if err != nil {
		return nil, err
	}

	// Cache for 10 minutes
	s.cacheTrending(ctx, cacheKey, posts, 10*time.Minute)

	return posts, nil
}

// SavePost saves a post to user's saved collection
func (s *PostService) SavePost(ctx context.Context, userID, postID uuid.UUID, collection *string) error {
	// Verify post exists
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		return fmt.Errorf("post not found")
	}

	save := &model.Save{
		ID:         uuid.New(),
		UserID:     userID,
		PostID:     postID,
		Collection: collection,
		CreatedAt:  time.Now(),
	}

	if err := s.saveRepo.Create(ctx, save); err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}

	// Increment saves count
	s.postRepo.IncrementSaves(ctx, postID)

	return nil
}

// UnsavePost removes a post from saved collection
func (s *PostService) UnsavePost(ctx context.Context, userID, postID uuid.UUID) error {
	if err := s.saveRepo.Delete(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unsave post: %w", err)
	}

	// Decrement saves count
	s.postRepo.DecrementSaves(ctx, postID)

	return nil
}

// GetSavedPosts retrieves user's saved posts
func (s *PostService) GetSavedPosts(ctx context.Context, userID uuid.UUID, collection *string, limit, offset int) ([]model.Save, error) {
	return s.saveRepo.GetByUserID(ctx, userID, collection, limit, offset)
}

// Helper functions

func (s *PostService) extractHashtags(caption string) []string {
	var hashtags []string
	words := strings.Fields(caption)
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			hashtag := strings.TrimPrefix(word, "#")
			hashtag = strings.ToLower(hashtag)
			// Remove punctuation from end
			hashtag = strings.TrimRight(hashtag, ".,!?;:")
			if len(hashtag) > 0 {
				hashtags = append(hashtags, hashtag)
			}
		}
	}
	
	return s.deduplicateStrings(hashtags)
}

func (s *PostService) deduplicateStrings(strs []string) []string {
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

// Cache methods (same as Socialink with minor adjustments)
func (s *PostService) cachePost(ctx context.Context, post *model.Post) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("post:%s", post.ID.String())
	data, _ := json.Marshal(post)
	s.redis.Set(ctx, key, data, 1*time.Hour)
}

func (s *PostService) getPostFromCache(ctx context.Context, postID uuid.UUID) (*model.Post, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	key := fmt.Sprintf("post:%s", postID.String())
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var post model.Post
	if err := json.Unmarshal(data, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostService) invalidatePostCache(ctx context.Context, postID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("post:%s", postID.String())
	s.redis.Del(ctx, key)
}

func (s *PostService) cacheFeed(ctx context.Context, userID uuid.UUID, posts []model.Post) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("feed:%s", userID.String())
	data, _ := json.Marshal(posts)
	s.redis.Set(ctx, key, data, 10*time.Minute)
}

func (s *PostService) getFeedFromCache(ctx context.Context, userID uuid.UUID, limit int) ([]model.Post, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	key := fmt.Sprintf("feed:%s", userID.String())
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) invalidateFeedCache(ctx context.Context, userID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("feed:%s", userID.String())
	s.redis.Del(ctx, key)
}

func (s *PostService) cacheTrending(ctx context.Context, key string, posts []model.Post, ttl time.Duration) {
	if s.redis == nil {
		return
	}

	data, _ := json.Marshal(posts)
	s.redis.Set(ctx, key, data, ttl)
}

func (s *PostService) getTrendingFromCache(ctx context.Context, key string) ([]model.Post, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// Kafka event publishing
func (s *PostService) publishPostCreatedEvent(post *model.Post) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "post.created",
		"post_id":    post.ID.String(),
		"user_id":    post.UserID.String(),
		"is_reels":   post.IsReels,
		"hashtags":   post.Hashtags,
		"created_at": post.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", post.ID.String(), event)
}

func (s *PostService) publishPostUpdatedEvent(post *model.Post) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "post.updated",
		"post_id":    post.ID.String(),
		"user_id":    post.UserID.String(),
		"updated_at": post.UpdatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", post.ID.String(), event)
}

func (s *PostService) publishPostDeletedEvent(postID, userID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "post.deleted",
		"post_id":    postID.String(),
		"user_id":    userID.String(),
		"deleted_at": time.Now(),
	}

	s.kafka.PublishEvent(context.Background(), "post-events", postID.String(), event)
}
