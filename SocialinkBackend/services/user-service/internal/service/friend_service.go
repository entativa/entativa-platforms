package service

import (
	"context"
	"fmt"
	"time"

	"github.com/entativa/socialink/user-service/internal/model"
	"github.com/entativa/socialink/user-service/internal/repository"
	"github.com/google/uuid"
)

type FriendService struct {
	friendRepo *repository.FriendRepository
	userRepo   *repository.UserRepository
	kafka      *KafkaProducer
}

func NewFriendService(
	friendRepo *repository.FriendRepository,
	userRepo *repository.UserRepository,
	kafka *KafkaProducer,
) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
		kafka:      kafka,
	}
}

// ==========================
// FRIEND REQUEST WORKFLOW
// ==========================

// SendFriendRequest - Send a friend request
func (s *FriendService) SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID, message *string) error {
	// Validation: Can't send to yourself
	if senderID == receiverID {
		return fmt.Errorf("cannot send friend request to yourself")
	}

	// Check if receiver exists
	_, err := s.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if already friends
	isFriend, err := s.friendRepo.AreFriends(ctx, senderID, receiverID)
	if err != nil {
		return err
	}
	if isFriend {
		return fmt.Errorf("already friends with this user")
	}

	// Check for existing pending request (either direction)
	existingRequest, err := s.friendRepo.GetPendingRequest(ctx, senderID, receiverID)
	if err == nil && existingRequest != nil {
		if existingRequest.SenderID == senderID {
			return fmt.Errorf("friend request already sent")
		} else {
			// They sent you a request! Just accept it automatically
			return s.AcceptFriendRequest(ctx, existingRequest.ID, receiverID)
		}
	}

	// SPAM PREVENTION: Check daily limit
	todayCount, err := s.friendRepo.CountRequestsSentToday(ctx, senderID)
	if err != nil {
		return err
	}
	if todayCount >= model.MaxDailyRequests {
		return fmt.Errorf("daily friend request limit reached (%d)", model.MaxDailyRequests)
	}

	// SPAM PREVENTION: Check pending requests limit
	pendingCount, err := s.friendRepo.CountPendingRequestsSent(ctx, senderID)
	if err != nil {
		return err
	}
	if pendingCount >= model.MaxPendingRequests {
		return fmt.Errorf("too many pending requests (%d). cancel some first", model.MaxPendingRequests)
	}

	// LIMIT CHECK: Check if sender has room for more friends
	senderFriendsCount, err := s.friendRepo.GetFriendsCount(ctx, senderID)
	if err != nil {
		return err
	}
	if senderFriendsCount >= model.MaxFriendsLimit {
		return fmt.Errorf("you've reached the friend limit (%d)", model.MaxFriendsLimit)
	}

	// LIMIT CHECK: Check if receiver has room
	receiverFriendsCount, err := s.friendRepo.GetFriendsCount(ctx, receiverID)
	if err != nil {
		return err
	}
	if receiverFriendsCount >= model.MaxFriendsLimit {
		return fmt.Errorf("user has reached their friend limit")
	}

	// Create friend request
	request := &model.FriendRequest{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
		Message:    message,
		CreatedAt:  time.Now(),
	}

	err = s.friendRepo.CreateRequest(ctx, request)
	if err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishFriendRequestEvent(senderID, receiverID, "sent")

	return nil
}

// AcceptFriendRequest - Accept a friend request
func (s *FriendService) AcceptFriendRequest(ctx context.Context, requestID, receiverID uuid.UUID) error {
	// Get request
	request, err := s.friendRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("request not found")
	}

	// Verify receiver
	if request.ReceiverID != receiverID {
		return fmt.Errorf("unauthorized")
	}

	// Check if already accepted
	if request.Status != "pending" {
		return fmt.Errorf("request already %s", request.Status)
	}

	// LIMIT CHECK: Check if receiver (now accepting) has room
	receiverFriendsCount, err := s.friendRepo.GetFriendsCount(ctx, receiverID)
	if err != nil {
		return err
	}
	if receiverFriendsCount >= model.MaxFriendsLimit {
		return fmt.Errorf("you've reached the friend limit (%d)", model.MaxFriendsLimit)
	}

	// LIMIT CHECK: Check if sender still has room
	senderFriendsCount, err := s.friendRepo.GetFriendsCount(ctx, request.SenderID)
	if err != nil {
		return err
	}
	if senderFriendsCount >= model.MaxFriendsLimit {
		return fmt.Errorf("sender has reached their friend limit")
	}

	// Update request status
	err = s.friendRepo.UpdateRequestStatus(ctx, requestID, "accepted")
	if err != nil {
		return err
	}

	// Create friendship (bi-directional)
	userID1, userID2 := request.SenderID, request.ReceiverID
	if userID1.String() > userID2.String() {
		userID1, userID2 = userID2, userID1
	}

	friendship := &model.Friend{
		ID:        uuid.New(),
		UserID1:   userID1,
		UserID2:   userID2,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	err = s.friendRepo.CreateFriendship(ctx, friendship)
	if err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishFriendRequestEvent(request.SenderID, request.ReceiverID, "accepted")

	return nil
}

// RejectFriendRequest - Reject a friend request
func (s *FriendService) RejectFriendRequest(ctx context.Context, requestID, receiverID uuid.UUID) error {
	// Get request
	request, err := s.friendRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("request not found")
	}

	// Verify receiver
	if request.ReceiverID != receiverID {
		return fmt.Errorf("unauthorized")
	}

	// Check if pending
	if request.Status != "pending" {
		return fmt.Errorf("request already %s", request.Status)
	}

	// Update status
	err = s.friendRepo.UpdateRequestStatus(ctx, requestID, "rejected")
	if err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishFriendRequestEvent(request.SenderID, request.ReceiverID, "rejected")

	return nil
}

// CancelFriendRequest - Cancel a sent friend request
func (s *FriendService) CancelFriendRequest(ctx context.Context, requestID, senderID uuid.UUID) error {
	// Get request
	request, err := s.friendRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("request not found")
	}

	// Verify sender
	if request.SenderID != senderID {
		return fmt.Errorf("unauthorized")
	}

	// Check if pending
	if request.Status != "pending" {
		return fmt.Errorf("request already %s", request.Status)
	}

	// Update status
	err = s.friendRepo.UpdateRequestStatus(ctx, requestID, "cancelled")
	if err != nil {
		return err
	}

	return nil
}

// ==========================
// FRIEND MANAGEMENT
// ==========================

// Unfriend - Remove a friend
func (s *FriendService) Unfriend(ctx context.Context, userID, friendID uuid.UUID) error {
	// Check if friends
	isFriend, err := s.friendRepo.AreFriends(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if !isFriend {
		return fmt.Errorf("not friends with this user")
	}

	// Remove friendship
	err = s.friendRepo.DeleteFriendship(ctx, userID, friendID)
	if err != nil {
		return err
	}

	// Publish Kafka event
	s.kafka.PublishFriendRequestEvent(userID, friendID, "unfriended")

	return nil
}

// GetFriends - Get list of friends
func (s *FriendService) GetFriends(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	return s.friendRepo.GetFriends(ctx, userID, limit, offset)
}

// GetFriendRequests - Get incoming friend requests
func (s *FriendService) GetFriendRequests(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.FriendRequest, error) {
	return s.friendRepo.GetIncomingRequests(ctx, userID, limit, offset)
}

// GetSentRequests - Get outgoing friend requests
func (s *FriendService) GetSentRequests(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.FriendRequest, error) {
	return s.friendRepo.GetOutgoingRequests(ctx, userID, limit, offset)
}

// GetFriendStats - Get friend statistics
func (s *FriendService) GetFriendStats(ctx context.Context, userID uuid.UUID) (*model.FriendStats, error) {
	friendsCount, err := s.friendRepo.GetFriendsCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	pendingCount, err := s.friendRepo.CountPendingRequestsReceived(ctx, userID)
	if err != nil {
		return nil, err
	}

	sentCount, err := s.friendRepo.CountPendingRequestsSent(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.FriendStats{
		UserID:                 userID,
		FriendsCount:           friendsCount,
		PendingRequestsCount:   pendingCount,
		SentRequestsCount:      sentCount,
		MutualFriendsAvailable: model.MaxFriendsLimit - friendsCount,
	}, nil
}

// AreFriends - Check if two users are friends
func (s *FriendService) AreFriends(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error) {
	return s.friendRepo.AreFriends(ctx, userID1, userID2)
}

// GetMutualFriends - Get mutual friends between two users
func (s *FriendService) GetMutualFriends(ctx context.Context, userID1, userID2 uuid.UUID) ([]uuid.UUID, error) {
	return s.friendRepo.GetMutualFriends(ctx, userID1, userID2)
}
