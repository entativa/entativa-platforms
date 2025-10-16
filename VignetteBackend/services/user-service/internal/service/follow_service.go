package service

import (
	"context"
	"fmt"
	"time"

	"github.com/entativa/vignette/user-service/internal/model"
	"github.com/entativa/vignette/user-service/internal/repository"
	"github.com/google/uuid"
)

type FollowService struct {
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
	kafka      *KafkaProducer
}

func NewFollowService(
	followRepo *repository.FollowRepository,
	userRepo *repository.UserRepository,
	kafka *KafkaProducer,
) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
		kafka:      kafka,
	}
}

// Follow - Follow a user (Instagram-style, instant)
func (s *FollowService) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	// Validate: Can't follow yourself
	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if already following
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return fmt.Errorf("already following this user")
	}

	// Create follow relationship
	follow := &model.Follow{
		ID:          uuid.New(),
		FollowerID:  followerID,
		FollowingID: followingID,
		Status:      "active",
		CreatedAt:   time.Now(),
	}

	err = s.followRepo.Create(ctx, follow)
	if err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishFollowEvent(followerID, followingID, "followed")

	return nil
}

// Unfollow - Unfollow a user
func (s *FollowService) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	// Check if following
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return fmt.Errorf("not following this user")
	}

	// Remove follow
	err = s.followRepo.Delete(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishFollowEvent(followerID, followingID, "unfollowed")

	return nil
}

// GetFollowers - Get list of followers
func (s *FollowService) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	return s.followRepo.GetFollowers(ctx, userID, limit, offset)
}

// GetFollowing - Get list of users being followed
func (s *FollowService) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	return s.followRepo.GetFollowing(ctx, userID, limit, offset)
}

// GetFollowStats - Get follower/following counts
func (s *FollowService) GetFollowStats(ctx context.Context, userID uuid.UUID) (*model.FollowStats, error) {
	followersCount, err := s.followRepo.GetFollowersCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	followingCount, err := s.followRepo.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.FollowStats{
		UserID:         userID,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
	}, nil
}

// IsFollowing - Check if user A follows user B
func (s *FollowService) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return s.followRepo.IsFollowing(ctx, followerID, followingID)
}

// GetMutualFollows - Get mutual follows between two users
func (s *FollowService) GetMutualFollows(ctx context.Context, userID1, userID2 uuid.UUID) ([]uuid.UUID, error) {
	return s.followRepo.GetMutualFollows(ctx, userID1, userID2)
}

// RemoveFollower - Remove a follower (block them from following you)
func (s *FollowService) RemoveFollower(ctx context.Context, userID, followerID uuid.UUID) error {
	// Remove the follow relationship
	return s.followRepo.Delete(ctx, followerID, userID)
}
