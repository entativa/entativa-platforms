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

type LikeService struct {
	likeRepo    repository.LikeRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	redis       *redis.Client
	kafka       *kafka.Producer
}

func NewLikeService(
	likeRepo repository.LikeRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	redis *redis.Client,
	kafka *kafka.Producer,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		redis:       redis,
		kafka:       kafka,
	}
}

// LikePost creates or updates a like/reaction on a post
func (s *LikeService) LikePost(ctx context.Context, postID, userID uuid.UUID, reactionType model.ReactionType) error {
	// Verify post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found")
	}

	// Check if already liked
	existingLike, _ := s.likeRepo.GetPostLike(ctx, userID, postID)

	like := &model.Like{
		ID:           uuid.New(),
		UserID:       userID,
		PostID:       &postID,
		ReactionType: reactionType,
		CreatedAt:    time.Now(),
	}

	if existingLike != nil {
		// Update existing reaction
		like.ID = existingLike.ID
	} else {
		// New like - increment count
		if err := s.postRepo.IncrementLikes(ctx, postID); err != nil {
			fmt.Printf("Failed to increment like count: %v\n", err)
		}
	}

	if err := s.likeRepo.CreatePostLike(ctx, like); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	// Invalidate caches
	s.invalidatePostCache(ctx, postID)

	// Publish event
	s.publishPostLikedEvent(like, post.UserID)

	return nil
}

// UnlikePost removes a like from a post
func (s *LikeService) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	// Check if like exists
	existingLike, err := s.likeRepo.GetPostLike(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}

	if existingLike == nil {
		return fmt.Errorf("post not liked")
	}

	// Delete like
	if err := s.likeRepo.DeletePostLike(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	// Decrement count
	if err := s.postRepo.DecrementLikes(ctx, postID); err != nil {
		fmt.Printf("Failed to decrement like count: %v\n", err)
	}

	// Invalidate caches
	s.invalidatePostCache(ctx, postID)

	// Publish event
	s.publishPostUnlikedEvent(postID, userID)

	return nil
}

// LikeComment creates or updates a like/reaction on a comment
func (s *LikeService) LikeComment(ctx context.Context, commentID, userID uuid.UUID, reactionType model.ReactionType) error {
	// Verify comment exists
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("comment not found")
	}

	// Check if already liked
	existingLike, _ := s.likeRepo.GetCommentLike(ctx, userID, commentID)

	like := &model.Like{
		ID:           uuid.New(),
		UserID:       userID,
		CommentID:    &commentID,
		ReactionType: reactionType,
		CreatedAt:    time.Now(),
	}

	if existingLike != nil {
		// Update existing reaction
		like.ID = existingLike.ID
	} else {
		// New like - increment count
		if err := s.commentRepo.IncrementLikes(ctx, commentID); err != nil {
			fmt.Printf("Failed to increment comment like count: %v\n", err)
		}
	}

	if err := s.likeRepo.CreateCommentLike(ctx, like); err != nil {
		return fmt.Errorf("failed to like comment: %w", err)
	}

	// Invalidate caches
	s.invalidateCommentsCache(ctx, comment.PostID)

	// Publish event
	s.publishCommentLikedEvent(like, comment.UserID, comment.PostID)

	return nil
}

// UnlikeComment removes a like from a comment
func (s *LikeService) UnlikeComment(ctx context.Context, commentID, userID uuid.UUID) error {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("comment not found")
	}

	// Check if like exists
	existingLike, err := s.likeRepo.GetCommentLike(ctx, userID, commentID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}

	if existingLike == nil {
		return fmt.Errorf("comment not liked")
	}

	// Delete like
	if err := s.likeRepo.DeleteCommentLike(ctx, userID, commentID); err != nil {
		return fmt.Errorf("failed to unlike comment: %w", err)
	}

	// Decrement count
	if err := s.commentRepo.DecrementLikes(ctx, commentID); err != nil {
		fmt.Printf("Failed to decrement comment like count: %v\n", err)
	}

	// Invalidate caches
	s.invalidateCommentsCache(ctx, comment.PostID)

	// Publish event
	s.publishCommentUnlikedEvent(commentID, userID)

	return nil
}

// GetPostLikers retrieves users who liked a post
func (s *LikeService) GetPostLikers(ctx context.Context, postID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	return s.likeRepo.GetPostLikers(ctx, postID, limit, offset)
}

// GetUserPostReaction gets the user's reaction to a post
func (s *LikeService) GetUserPostReaction(ctx context.Context, userID, postID uuid.UUID) (*model.ReactionType, error) {
	return s.likeRepo.GetUserPostReaction(ctx, userID, postID)
}

// Cache methods
func (s *LikeService) invalidatePostCache(ctx context.Context, postID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("post:%s", postID.String())
	s.redis.Del(ctx, key)
}

func (s *LikeService) invalidateCommentsCache(ctx context.Context, postID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("comments:%s", postID.String())
	s.redis.Del(ctx, key)
}

// Kafka event publishing
func (s *LikeService) publishPostLikedEvent(like *model.Like, postOwnerID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":    "post.liked",
		"like_id":       like.ID.String(),
		"post_id":       like.PostID.String(),
		"user_id":       like.UserID.String(),
		"post_owner_id": postOwnerID.String(),
		"reaction_type": like.ReactionType,
		"created_at":    like.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", like.ID.String(), event)
}

func (s *LikeService) publishPostUnlikedEvent(postID, userID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "post.unliked",
		"post_id":    postID.String(),
		"user_id":    userID.String(),
		"unliked_at": time.Now(),
	}

	s.kafka.PublishEvent(context.Background(), "post-events", fmt.Sprintf("%s_%s", postID, userID), event)
}

func (s *LikeService) publishCommentLikedEvent(like *model.Like, commentOwnerID, postID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":       "comment.liked",
		"like_id":          like.ID.String(),
		"comment_id":       like.CommentID.String(),
		"user_id":          like.UserID.String(),
		"comment_owner_id": commentOwnerID.String(),
		"post_id":          postID.String(),
		"reaction_type":    like.ReactionType,
		"created_at":       like.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", like.ID.String(), event)
}

func (s *LikeService) publishCommentUnlikedEvent(commentID, userID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "comment.unliked",
		"comment_id": commentID.String(),
		"user_id":    userID.String(),
		"unliked_at": time.Now(),
	}

	s.kafka.PublishEvent(context.Background(), "post-events", fmt.Sprintf("%s_%s", commentID, userID), event)
}
