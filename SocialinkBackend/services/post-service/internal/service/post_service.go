package service

import (
	"context"
	"fmt"
	"time"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/repository"
	"socialink/post-service/pkg/kafka"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PostService struct {
	postRepo    repository.PostRepository
	likeRepo    repository.LikeRepository
	commentRepo repository.CommentRepository
	shareRepo   repository.ShareRepository
	redis       *redis.Client
	kafka       *kafka.Producer
}

func NewPostService(
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	shareRepo repository.ShareRepository,
	redis *redis.Client,
	kafka *kafka.Producer,
) *PostService {
	return &PostService{
		postRepo:    postRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		shareRepo:   shareRepo,
		redis:       redis,
		kafka:       kafka,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(ctx context.Context, userID uuid.UUID, req *model.CreatePostRequest) (*model.Post, error) {
	// Validate privacy setting
	if !s.isValidPrivacy(req.Privacy) {
		return nil, fmt.Errorf("invalid privacy setting")
	}

	// Create post
	post := &model.Post{
		ID:            uuid.New(),
		UserID:        userID,
		Content:       req.Content,
		MediaIDs:      req.MediaIDs,
		Privacy:       req.Privacy,
		Location:      req.Location,
		TaggedUserIDs: req.TaggedUserIDs,
		Feeling:       req.Feeling,
		Activity:      req.Activity,
		LikesCount:    0,
		CommentsCount: 0,
		SharesCount:   0,
		IsEdited:      false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
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

// GetPost retrieves a post by ID with permission check
func (s *PostService) GetPost(ctx context.Context, postID, requestingUserID uuid.UUID) (*model.Post, error) {
	// Try cache first
	if post, err := s.getPostFromCache(ctx, postID); err == nil && post != nil {
		// Check permissions
		if s.canViewPost(ctx, post, requestingUserID) {
			return post, nil
		}
		return nil, fmt.Errorf("permission denied")
	}

	// Get from database
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !s.canViewPost(ctx, post, requestingUserID) {
		return nil, fmt.Errorf("permission denied")
	}

	// Cache it
	s.cachePost(ctx, post)

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
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Privacy != nil {
		post.Privacy = *req.Privacy
	}
	if req.Location != nil {
		post.Location = req.Location
	}
	if req.Feeling != nil {
		post.Feeling = req.Feeling
	}
	if req.Activity != nil {
		post.Activity = req.Activity
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
func (s *PostService) GetUserPosts(ctx context.Context, userID, requestingUserID uuid.UUID, limit, offset int) ([]model.Post, error) {
	posts, err := s.postRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Filter based on permissions
	var visiblePosts []model.Post
	for _, post := range posts {
		if s.canViewPost(ctx, &post, requestingUserID) {
			visiblePosts = append(visiblePosts, post)
		}
	}

	return visiblePosts, nil
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

// GetTrendingPosts retrieves trending posts
func (s *PostService) GetTrendingPosts(ctx context.Context, limit int) ([]model.Post, error) {
	// Check cache
	cacheKey := fmt.Sprintf("trending:posts:%d", limit)
	if cachedPosts, err := s.getTrendingFromCache(ctx, cacheKey); err == nil && len(cachedPosts) > 0 {
		return cachedPosts, nil
	}

	// Get from database (last 24 hours)
	posts, err := s.postRepo.GetTrendingPosts(ctx, limit, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	s.cacheTrending(ctx, cacheKey, posts, 5*time.Minute)

	return posts, nil
}

// Permission checking
func (s *PostService) canViewPost(ctx context.Context, post *model.Post, viewerID uuid.UUID) bool {
	// Post owner can always view
	if post.UserID == viewerID {
		return true
	}

	// Check privacy
	switch post.Privacy {
	case model.PrivacyPublic:
		return true
	case model.PrivacyOnlyMe:
		return false
	case model.PrivacyFriends:
		// In production, check if viewer is friend
		// For now, allow if authenticated
		return viewerID != uuid.Nil
	default:
		return false
	}
}

func (s *PostService) isValidPrivacy(privacy model.Privacy) bool {
	switch privacy {
	case model.PrivacyPublic, model.PrivacyFriends, model.PrivacyFriendsExcept,
		model.PrivacySpecificFriends, model.PrivacyOnlyMe, model.PrivacyCustom:
		return true
	default:
		return false
	}
}

// Cache methods
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
		"privacy":    post.Privacy,
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

// Add missing import
import "encoding/json"
