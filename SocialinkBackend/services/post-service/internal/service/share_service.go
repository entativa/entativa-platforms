package service

import (
	"context"
	"fmt"
	"time"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/repository"
	"socialink/post-service/pkg/kafka"

	"github.com/google/uuid"
)

type ShareService struct {
	shareRepo repository.ShareRepository
	postRepo  repository.PostRepository
	kafka     *kafka.Producer
}

func NewShareService(
	shareRepo repository.ShareRepository,
	postRepo repository.PostRepository,
	kafka *kafka.Producer,
) *ShareService {
	return &ShareService{
		shareRepo: shareRepo,
		postRepo:  postRepo,
		kafka:     kafka,
	}
}

// SharePost creates a share of a post
func (s *ShareService) SharePost(ctx context.Context, originalPostID, userID uuid.UUID, req *model.SharePostRequest) (*model.Share, error) {
	// Verify original post exists
	originalPost, err := s.postRepo.GetByID(ctx, originalPostID)
	if err != nil {
		return nil, fmt.Errorf("original post not found")
	}

	// Check if post is shareable (respect privacy)
	if originalPost.Privacy == model.PrivacyOnlyMe {
		return nil, fmt.Errorf("this post cannot be shared")
	}

	// Create share
	share := &model.Share{
		ID:             uuid.New(),
		UserID:         userID,
		OriginalPostID: originalPostID,
		Caption:        req.Caption,
		Privacy:        req.Privacy,
		CreatedAt:      time.Now(),
	}

	if err := s.shareRepo.Create(ctx, share); err != nil {
		return nil, fmt.Errorf("failed to share post: %w", err)
	}

	// Increment share count on original post
	if err := s.postRepo.IncrementShares(ctx, originalPostID); err != nil {
		fmt.Printf("Failed to increment share count: %v\n", err)
	}

	// Publish event
	s.publishPostSharedEvent(share, originalPost.UserID)

	return share, nil
}

// GetUserShares retrieves shares by a user
func (s *ShareService) GetUserShares(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Share, error) {
	return s.shareRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetPostShares retrieves shares of a post
func (s *ShareService) GetPostShares(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Share, error) {
	return s.shareRepo.GetByOriginalPostID(ctx, postID, limit, offset)
}

// DeleteShare removes a share
func (s *ShareService) DeleteShare(ctx context.Context, shareID, userID uuid.UUID) error {
	// Get share
	share, err := s.shareRepo.GetByID(ctx, shareID)
	if err != nil {
		return err
	}

	// Verify ownership
	if share.UserID != userID {
		return fmt.Errorf("permission denied: not the share owner")
	}

	// Delete share
	if err := s.shareRepo.Delete(ctx, shareID); err != nil {
		return fmt.Errorf("failed to delete share: %w", err)
	}

	// Note: We don't decrement shares_count to preserve historical data
	// But you could choose to do so

	// Publish event
	s.publishPostUnsharedEvent(shareID, share.OriginalPostID, userID)

	return nil
}

// Kafka event publishing
func (s *ShareService) publishPostSharedEvent(share *model.Share, postOwnerID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":    "post.shared",
		"share_id":      share.ID.String(),
		"post_id":       share.OriginalPostID.String(),
		"user_id":       share.UserID.String(),
		"post_owner_id": postOwnerID.String(),
		"privacy":       share.Privacy,
		"created_at":    share.CreatedAt,
	}

	s.kafka.PublishEvent(context.Background(), "post-events", share.ID.String(), event)
}

func (s *ShareService) publishPostUnsharedEvent(shareID, postID, userID uuid.UUID) {
	if s.kafka == nil {
		return
	}

	event := map[string]interface{}{
		"event_type": "post.unshared",
		"share_id":   shareID.String(),
		"post_id":    postID.String(),
		"user_id":    userID.String(),
		"deleted_at": time.Now(),
	}

	s.kafka.PublishEvent(context.Background(), "post-events", shareID.String(), event)
}
