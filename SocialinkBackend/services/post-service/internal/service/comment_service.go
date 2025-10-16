package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/repository"
	"socialink/post-service/pkg/kafka"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	redis       *redis.Client
	kafka       *kafka.Producer
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	redis *redis.Client,
	kafka *kafka.Producer,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		redis:       redis,
		kafka:       kafka,
	}
}

// CreateComment creates a new comment on a post
func (s *CommentService) CreateComment(ctx context.Context, postID, userID uuid.UUID, req *model.CreateCommentRequest) (*model.Comment, error) {
	// Verify post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// If parent comment specified, verify it exists and belongs to same post
	if req.ParentID != nil {
		parentComment, err := s.commentRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent comment not found")
		}
		if parentComment.PostID != postID {
			return nil, fmt.Errorf("parent comment does not belong to this post")
		}
	}

	// Create comment
	comment := &model.Comment{
		ID:         uuid.New(),
		PostID:     postID,
		UserID:     userID,
		ParentID:   req.ParentID,
		Content:    req.Content,
		MediaID:    req.MediaID,
		LikesCount: 0,
		IsEdited:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Increment post comment count
	if err := s.postRepo.IncrementComments(ctx, postID); err != nil {
		// Log error but don't fail the comment creation
		fmt.Printf("Failed to increment comment count: %v\n", err)
	}

	// Invalidate caches
	s.invalidatePostCache(ctx, postID)
	s.invalidateCommentsCache(ctx, postID)

	// Publish event
	s.publishCommentCreatedEvent(comment, post.UserID)

	return comment, nil
}

// GetComments retrieves comments for a post
func (s *CommentService) GetComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	// Try cache for first page
	if offset == 0 {
		if cachedComments, err := s.getCommentsFromCache(ctx, postID, limit); err == nil && len(cachedComments) > 0 {
			return cachedComments, nil
		}
	}

	// Get from database
	comments, err := s.commentRepo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache first page
	if offset == 0 {
		s.cacheComments(ctx, postID, comments)
	}

	return comments, nil
}

// GetReplies retrieves replies to a comment
func (s *CommentService) GetReplies(ctx context.Context, commentID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	return s.commentRepo.GetReplies(ctx, commentID, limit, offset)
}

// UpdateComment updates a comment
func (s *CommentService) UpdateComment(ctx context.Context, commentID, userID uuid.UUID, req *model.UpdateCommentRequest) (*model.Comment, error) {
	// Get existing comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if comment.UserID != userID {
		return nil, fmt.Errorf("permission denied: not the comment owner")
	}

	// Update
	comment.Content = req.Content
	now := time.Now()
	comment.IsEdited = true
	comment.EditedAt = &now
	comment.UpdatedAt = now

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// Invalidate caches
	s.invalidateCommentsCache(ctx, comment.PostID)

	// Publish event
	s.publishCommentUpdatedEvent(comment)

	return comment, nil
}

// DeleteComment soft deletes a comment
func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Verify ownership
	if comment.UserID != userID {
		return fmt.Errorf("permission denied: not the comment owner")
	}

	// Delete
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Decrement post comment count
	if err := s.postRepo.DecrementComments(ctx, comment.PostID); err != nil {
		fmt.Printf("Failed to decrement comment count: %v\n", err)
	}

	// Invalidate caches
	s.invalidatePostCache(ctx, comment.PostID)
	s.invalidateCommentsCache(ctx, comment.PostID)

	// Publish event
	s.publishCommentDeletedEvent(commentID, comment.PostID, userID)

	return nil
}

// Cache methods
func (s *CommentService) cacheComments(ctx context.Context, postID uuid.UUID, comments []model.Comment) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("comments:%s", postID.String())
	data, _ := json.Marshal(comments)
	s.redis.Set(ctx, key, data, 30*time.Minute)
}

func (s *CommentService) getCommentsFromCache(ctx context.Context, postID uuid.UUID, limit int) ([]model.Comment, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not available")
	}

	key := fmt.Sprintf("comments:%s", postID.String())
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var comments []model.Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, err
	}

	if len(comments) > limit {
		comments = comments[:limit]
	}

	return comments, nil
}

func (s *CommentService) invalidateCommentsCache(ctx context.Context, postID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("comments:%s", postID.String())
	s.redis.Del(ctx, key)
}

func (s *CommentService) invalidatePostCache(ctx context.Context, postID uuid.UUID) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("post:%s", postID.String())
	s.redis.Del(ctx, key)
}

// Kafka event publishing
func (s *CommentService) publishCommentCreatedEvent(comment *model.Comment, postOwnerID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":    "comment.created",
		"comment_id":    comment.ID.String(),
		"post_id":       comment.PostID.String(),
		"user_id":       comment.UserID.String(),
		"post_owner_id": postOwnerID.String(),
		"parent_id":     comment.ParentID,
		"created_at":    comment.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", comment.ID.String(), event)
}

func (s *CommentService) publishCommentUpdatedEvent(comment *model.Comment) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "comment.updated",
		"comment_id": comment.ID.String(),
		"post_id":    comment.PostID.String(),
		"user_id":    comment.UserID.String(),
		"updated_at": comment.UpdatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", comment.ID.String(), event)
}

func (s *CommentService) publishCommentDeletedEvent(commentID, postID, userID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "comment.deleted",
		"comment_id": commentID.String(),
		"post_id":    postID.String(),
		"user_id":    userID.String(),
		"deleted_at": time.Now(),
	}

	s.kafka.PublishEvent(context.Background(), "post-events", commentID.String(), event)
}
