package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"socialink/live-streaming-service/internal/model"
	"socialink/live-streaming-service/internal/repository"
	"socialink/live-streaming-service/internal/grpc"
	"github.com/google/uuid"
)

// Stub types for external dependencies
type KafkaProducer struct{}
type RedisClient struct{}

func (r *RedisClient) SetLiveStream(ctx context.Context, streamID, streamerID uuid.UUID) error { return nil }
func (r *RedisClient) RemoveLiveStream(ctx context.Context, streamID uuid.UUID) error { return nil }
func (r *RedisClient) AddViewer(ctx context.Context, streamID, viewerID uuid.UUID) error { return nil }
func (r *RedisClient) RemoveViewer(ctx context.Context, streamID, viewerID uuid.UUID) error { return nil }
func (r *RedisClient) PublishComment(ctx context.Context, streamID uuid.UUID, comment *model.StreamComment) error { return nil }
func (r *RedisClient) PublishReaction(ctx context.Context, streamID uuid.UUID, reaction *model.StreamReaction) error { return nil }
func (r *RedisClient) GetFollowerCount(ctx context.Context, userID uuid.UUID) (int, error) { return 0, fmt.Errorf("not implemented") }

func (k *KafkaProducer) PublishStreamEvent(streamID, streamerID uuid.UUID, eventType string) {}

type StreamingService struct {
	streamRepo   *repository.StreamRepository
	viewerRepo   *repository.ViewerRepository
	commentRepo  *repository.CommentRepository
	mediaGRPC    *grpc.MediaServiceClient // gRPC client for media service
	kafka        *KafkaProducer
	redis        *RedisClient
}

func NewStreamingService(
	streamRepo *repository.StreamRepository,
	viewerRepo *repository.ViewerRepository,
	commentRepo *repository.CommentRepository,
	mediaGRPC *grpc.MediaServiceClient,
	kafka *KafkaProducer,
	redis *RedisClient,
) *StreamingService {
	return &StreamingService{
		streamRepo:  streamRepo,
		viewerRepo:  viewerRepo,
		commentRepo: commentRepo,
		mediaGRPC:   mediaGRPC,
		kafka:       kafka,
		redis:       redis,
	}
}

// ============================================
// STREAM MANAGEMENT
// ============================================

// CreateStream - Create a new stream (scheduled or ready to go live)
func (s *StreamingService) CreateStream(ctx context.Context, req *model.CreateStreamRequest, streamerID uuid.UUID) (*model.LiveStream, error) {
	// ELIGIBILITY CHECK: Verify follower threshold
	eligible, err := s.CheckEligibility(ctx, streamerID)
	if err != nil {
		return nil, err
	}
	if !eligible.Eligible {
		return nil, fmt.Errorf("not eligible to stream: %s", eligible.Reason)
	}

	// Generate stream key (secure random)
	streamKey, err := generateStreamKey()
	if err != nil {
		return nil, err
	}

	streamID := uuid.New()

	// Build stream URLs
	rtmpURL := fmt.Sprintf("rtmp://stream.vignette.com/live/%s", streamKey)
	hlsURL := fmt.Sprintf("https://stream.vignette.com/hls/%s.m3u8", streamID)
	webrtcURL := fmt.Sprintf("wss://stream.vignette.com/webrtc/%s", streamID)

	stream := &model.LiveStream{
		ID:            streamID,
		StreamerID:    streamerID,
		Title:         req.Title,
		Description:   req.Description,
		ThumbnailURL:  req.ThumbnailURL,
		Status:        model.StatusScheduled,
		Quality:       req.Quality,
		IsPrivate:     req.IsPrivate,
		Category:      req.Category,
		Tags:          req.Tags,
		StreamKey:     streamKey,
		RTMPUrl:       rtmpURL,
		HLSUrl:        hlsURL,
		WebRTCUrl:     webrtcURL,
		ViewerCount:   0,
		PeakViewers:   0,
		TotalViews:    0,
		RecordStream:  req.RecordStream,
		ScheduledFor:  req.ScheduledFor,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.streamRepo.Create(ctx, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// StartStream - Start a live stream
func (s *StreamingService) StartStream(ctx context.Context, streamID, streamerID uuid.UUID) error {
	// Get stream
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return fmt.Errorf("stream not found")
	}

	// Verify ownership
	if stream.StreamerID != streamerID {
		return fmt.Errorf("unauthorized")
	}

	// Check status
	if stream.Status == model.StatusLive {
		return fmt.Errorf("stream already live")
	}
	if stream.Status == model.StatusEnded || stream.Status == model.StatusCancelled {
		return fmt.Errorf("stream already %s", stream.Status)
	}

	// Re-check eligibility (in case they lost followers)
	eligible, err := s.CheckEligibility(ctx, streamerID)
	if err != nil {
		return err
	}
	if !eligible.Eligible {
		return fmt.Errorf("no longer eligible: %s", eligible.Reason)
	}

	// Update status to live
	now := time.Now()
	stream.Status = model.StatusLive
	stream.StartedAt = &now
	stream.UpdatedAt = now

	err = s.streamRepo.Update(ctx, stream)
	if err != nil {
		return err
	}

	// Cache stream as live (for quick lookups)
	s.redis.SetLiveStream(ctx, streamID, streamerID)

	// Publish Kafka event (for notifications)
	s.kafka.PublishStreamEvent(streamID, streamerID, "stream_started")

	return nil
}

// EndStream - End a live stream
func (s *StreamingService) EndStream(ctx context.Context, streamID, streamerID uuid.UUID) error {
	// Get stream
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return fmt.Errorf("stream not found")
	}

	// Verify ownership
	if stream.StreamerID != streamerID {
		return fmt.Errorf("unauthorized")
	}

	// Check if live
	if stream.Status != model.StatusLive {
		return fmt.Errorf("stream not live")
	}

	// Calculate duration
	now := time.Now()
	duration := 0
	if stream.StartedAt != nil {
		duration = int(now.Sub(*stream.StartedAt).Seconds())
	}

	// Update status
	stream.Status = model.StatusEnded
	stream.EndedAt = &now
	stream.Duration = duration
	stream.UpdatedAt = now

	err = s.streamRepo.Update(ctx, stream)
	if err != nil {
		return err
	}

	// Remove from live cache
	s.redis.RemoveLiveStream(ctx, streamID)

	// Mark all viewers as inactive
	s.viewerRepo.MarkAllInactive(ctx, streamID)

	// If recording enabled, save to media service via gRPC
	if stream.RecordStream {
		go s.saveRecording(streamID, streamerID)
	}

	// Publish event
	s.kafka.PublishStreamEvent(streamID, streamerID, "stream_ended")

	return nil
}

// CheckEligibility - Check if user is eligible to stream
func (s *StreamingService) CheckEligibility(ctx context.Context, userID uuid.UUID) (*model.StreamEligibility, error) {
	// Get follower count (from user service or cache)
	followerCount, err := s.getFollowerCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check threshold (100 followers for Vignette)
	required := model.MinFriendsSocialink
	
	if followerCount >= required {
		return &model.StreamEligibility{
			Eligible:      true,
			FollowerCount: followerCount,
			Required:      required,
		}, nil
	}

	return &model.StreamEligibility{
		Eligible:      false,
		Reason:        fmt.Sprintf("Need %d followers to go live (you have %d)", required, followerCount),
		FollowerCount: followerCount,
		Required:      required,
	}, nil
}

// ============================================
// VIEWER MANAGEMENT
// ============================================

// JoinStream - User joins a stream
func (s *StreamingService) JoinStream(ctx context.Context, streamID, viewerID uuid.UUID) error {
	// Get stream
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return fmt.Errorf("stream not found")
	}

	// Check if live
	if stream.Status != model.StatusLive {
		return fmt.Errorf("stream not live")
	}

	// Check if private
	if stream.IsPrivate {
		// In production: Check if viewer is allowed
		// For now, just continue
	}

	// Check if already viewing
	existing, _ := s.viewerRepo.GetViewer(ctx, streamID, viewerID)
	if existing != nil && existing.IsActive {
		return nil // Already viewing
	}

	// Create viewer record
	viewer := &model.StreamViewer{
		ID:       uuid.New(),
		StreamID: streamID,
		ViewerID: viewerID,
		JoinedAt: time.Now(),
		IsActive: true,
	}

	err = s.viewerRepo.Create(ctx, viewer)
	if err != nil {
		return err
	}

	// Update Redis (active viewers)
	s.redis.AddViewer(ctx, streamID, viewerID)

	return nil
}

// LeaveStream - User leaves a stream
func (s *StreamingService) LeaveStream(ctx context.Context, streamID, viewerID uuid.UUID) error {
	// Get viewer
	viewer, err := s.viewerRepo.GetViewer(ctx, streamID, viewerID)
	if err != nil {
		return err
	}

	if !viewer.IsActive {
		return nil // Already left
	}

	// Calculate watch time
	watchTime := int(time.Since(viewer.JoinedAt).Seconds())
	now := time.Now()

	viewer.IsActive = false
	viewer.LeftAt = &now
	viewer.WatchTime = watchTime

	err = s.viewerRepo.Update(ctx, viewer)
	if err != nil {
		return err
	}

	// Update Redis
	s.redis.RemoveViewer(ctx, streamID, viewerID)

	return nil
}

// ============================================
// COMMENTS & REACTIONS
// ============================================

// PostComment - Post a comment on live stream
func (s *StreamingService) PostComment(ctx context.Context, streamID, userID uuid.UUID, content string) (*model.StreamComment, error) {
	// Get stream
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return nil, fmt.Errorf("stream not found")
	}

	// Check if live
	if stream.Status != model.StatusLive {
		return nil, fmt.Errorf("stream not live")
	}

	// Create comment
	comment := &model.StreamComment{
		ID:        uuid.New(),
		StreamID:  streamID,
		UserID:    userID,
		Content:   content,
		IsPinned:  false,
		CreatedAt: time.Now(),
	}

	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	// Publish to WebSocket (real-time)
	s.redis.PublishComment(ctx, streamID, comment)

	return comment, nil
}

// AddReaction - Add reaction to stream
func (s *StreamingService) AddReaction(ctx context.Context, streamID, userID uuid.UUID, reactionType string) error {
	// Get stream
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return fmt.Errorf("stream not found")
	}

	// Check if live
	if stream.Status != model.StatusLive {
		return fmt.Errorf("stream not live")
	}

	// Create/update reaction
	reaction := &model.StreamReaction{
		ID:        uuid.New(),
		StreamID:  streamID,
		UserID:    userID,
		Type:      reactionType,
		CreatedAt: time.Now(),
	}

	// In production: Use UPSERT
	// For now, just create
	// err = s.reactionRepo.Upsert(ctx, reaction)

	// Publish to WebSocket (real-time)
	s.redis.PublishReaction(ctx, streamID, reaction)

	return nil
}

// ============================================
// HELPERS
// ============================================

func (s *StreamingService) getFollowerCount(ctx context.Context, userID uuid.UUID) (int, error) {
	// Check cache first
	count, err := s.redis.GetFollowerCount(ctx, userID)
	if err == nil && count >= 0 {
		return count, nil
	}

	// Query from user service (HTTP or gRPC)
	// For now, mock
	return 500, nil // Mock: user has 500 followers
}

func generateStreamKey() (string, error) {
	// Generate 32-byte random key
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *StreamingService) saveRecording(streamID, streamerID uuid.UUID) {
	// In production: Call media service via gRPC to save recording
	// mediaGRPC.SaveStreamRecording(streamID, streamerID, recordingPath)
	
	// This runs in goroutine (async)
	ctx := context.Background()
	
	// Simulate recording save
	time.Sleep(5 * time.Second)
	
	recordingURL := fmt.Sprintf("https://cdn.vignette.com/recordings/%s.mp4", streamID)
	
	// Update stream with recording URL
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return
	}
	
	stream.RecordingURL = &recordingURL
	s.streamRepo.Update(ctx, stream)
}

// GetStream - Get stream by ID
func (s *StreamingService) GetStream(ctx context.Context, streamID uuid.UUID) (*model.LiveStream, error) {
	return s.streamRepo.GetByID(ctx, streamID)
}

// GetLiveStreams - Get all live streams
func (s *StreamingService) GetLiveStreams(ctx context.Context, limit, offset int, category string) ([]*model.LiveStream, error) {
	return s.streamRepo.GetLive(ctx, limit, offset, category)
}

// GetStreamAnalytics - Get stream analytics
func (s *StreamingService) GetStreamAnalytics(ctx context.Context, streamID uuid.UUID) (*model.StreamAnalytics, error) {
	stream, err := s.streamRepo.GetByID(ctx, streamID)
	if err != nil {
		return nil, err
	}

	uniqueViewers, err := s.viewerRepo.GetUniqueViewerCount(ctx, streamID)
	if err != nil {
		return nil, err
	}

	avgWatchTime, err := s.viewerRepo.GetAverageWatchTime(ctx, streamID)
	if err != nil {
		return nil, err
	}

	analytics := &model.StreamAnalytics{
		StreamID:       streamID,
		TotalViews:     stream.TotalViews,
		UniqueViewers:  uniqueViewers,
		PeakViewers:    stream.PeakViewers,
		AverageViewers: float64(stream.TotalViews) / float64(stream.Duration),
		AverageWatchTime: avgWatchTime,
		CommentsCount:  stream.CommentsCount,
		ReactionsCount: stream.LikesCount,
		SharesCount:    stream.SharesCount,
	}

	return analytics, nil
}
