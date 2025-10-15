package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"socialink/meta-user-service/internal/model"
	"socialink/meta-user-service/internal/repository"

	"github.com/google/uuid"
)

// CrossPlatformSyncService synchronizes user data across platforms
type CrossPlatformSyncService struct {
	repo           *repository.MetaUserRepository
	eventPublisher EventPublisher
	vignetteClient *VignetteClient // gRPC client for Vignette
	socialinkClient *SocialinkClient // gRPC client for Socialink
}

func NewCrossPlatformSyncService(
	repo *repository.MetaUserRepository,
	eventPublisher EventPublisher,
	vignetteClient *VignetteClient,
	socialinkClient *SocialinkClient,
) *CrossPlatformSyncService {
	return &CrossPlatformSyncService{
		repo:           repo,
		eventPublisher: eventPublisher,
		vignetteClient: vignetteClient,
		socialinkClient: socialinkClient,
	}
}

// SyncUserData synchronizes user data across platforms
func (s *CrossPlatformSyncService) SyncUserData(ctx context.Context, metaUserID uuid.UUID, sourcePlatform string) error {
	user, err := s.repo.GetByID(ctx, metaUserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.PlatformLinks.SyncEnabled {
		return nil // Sync disabled
	}

	// Determine target platform
	targetPlatform := "vignette"
	if sourcePlatform == "vignette" {
		targetPlatform = "socialink"
	}

	// Create sync payload
	syncData := &SyncData{
		MetaUserID:      metaUserID,
		Email:           user.Email,
		EmailVerified:   user.EmailVerified,
		PhoneNumber:     user.PhoneNumber,
		PhoneVerified:   user.PhoneVerified,
		PrivacySettings: user.PrivacySettings,
		ProfileData:     s.extractProfileData(user),
		SyncedAt:        time.Now(),
	}

	// Send to target platform
	var syncErr error
	if targetPlatform == "vignette" && user.PlatformLinks.VignetteUserID != nil {
		syncErr = s.syncToVignette(ctx, *user.PlatformLinks.VignetteUserID, syncData)
	} else if targetPlatform == "socialink" && user.PlatformLinks.SocialinkUserID != nil {
		syncErr = s.syncToSocialink(ctx, *user.PlatformLinks.SocialinkUserID, syncData)
	}

	if syncErr != nil {
		// Record sync conflict
		user.PlatformLinks.SyncConflicts++
		user.PlatformLinks.LinkStatus = model.LinkStatusConflicted
		s.repo.Update(ctx, user)
		
		log.Printf("Sync error for user %s: %v", metaUserID, syncErr)
		return syncErr
	}

	// Update sync timestamp
	now := time.Now()
	user.PlatformLinks.LastSyncedAt = &now
	user.PlatformLinks.LinkStatus = model.LinkStatusLinked
	s.repo.Update(ctx, user)

	// Publish sync event
	s.eventPublisher.PublishUserEvent(ctx, "meta.sync.completed", map[string]interface{}{
		"meta_user_id":    metaUserID,
		"source_platform": sourcePlatform,
		"target_platform": targetPlatform,
		"synced_at":       now,
	})

	return nil
}

// SyncProfileUpdates synchronizes profile updates
func (s *CrossPlatformSyncService) SyncProfileUpdates(ctx context.Context, metaUserID uuid.UUID, updates map[string]interface{}) error {
	user, err := s.repo.GetByID(ctx, metaUserID)
	if err != nil {
		return err
	}

	if !user.PlatformLinks.SyncEnabled {
		return nil
	}

	// Sync to both platforms
	errChan := make(chan error, 2)

	if user.PlatformLinks.SocialinkUserID != nil {
		go func() {
			errChan <- s.socialinkClient.UpdateProfile(ctx, *user.PlatformLinks.SocialinkUserID, updates)
		}()
	} else {
		errChan <- nil
	}

	if user.PlatformLinks.VignetteUserID != nil {
		go func() {
			errChan <- s.vignetteClient.UpdateProfile(ctx, *user.PlatformLinks.VignetteUserID, updates)
		}()
	} else {
		errChan <- nil
	}

	// Collect errors
	var errors []error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("sync errors: %v", errors)
	}

	return nil
}

// SyncPrivacySettings synchronizes privacy settings across platforms
func (s *CrossPlatformSyncService) SyncPrivacySettings(ctx context.Context, metaUserID uuid.UUID, settings model.PrivacySettings) error {
	user, err := s.repo.GetByID(ctx, metaUserID)
	if err != nil {
		return err
	}

	// Update meta user
	user.PrivacySettings = settings
	user.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	if !user.PlatformLinks.SyncEnabled {
		return nil
	}

	// Sync to platforms
	settingsJSON, _ := json.Marshal(settings)
	updates := map[string]interface{}{
		"privacy_settings": string(settingsJSON),
	}

	return s.SyncProfileUpdates(ctx, metaUserID, updates)
}

// EnableCrossPlatformSync enables synchronization between platforms
func (s *CrossPlatformSyncService) EnableCrossPlatformSync(ctx context.Context, metaUserID uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, metaUserID)
	if err != nil {
		return err
	}

	user.PlatformLinks.SyncEnabled = true
	user.PlatformLinks.LinkStatus = model.LinkStatusLinked
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	// Publish event
	s.eventPublisher.PublishUserEvent(ctx, "meta.sync.enabled", map[string]interface{}{
		"meta_user_id": metaUserID,
	})

	// Perform initial sync
	return s.SyncUserData(ctx, metaUserID, "meta")
}

// DisableCrossPlatformSync disables synchronization
func (s *CrossPlatformSyncService) DisableCrossPlatformSync(ctx context.Context, metaUserID uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, metaUserID)
	if err != nil {
		return err
	}

	user.PlatformLinks.SyncEnabled = false
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	s.eventPublisher.PublishUserEvent(ctx, "meta.sync.disabled", map[string]interface{}{
		"meta_user_id": metaUserID,
	})

	return nil
}

// ResolveSyncConflict resolves a sync conflict
func (s *CrossPlatformSyncService) ResolveSyncConflict(ctx context.Context, metaUserID uuid.UUID, resolution ConflictResolution) error {
	user, err := s.repo.GetByID(ctx, metaUserID)
	if err != nil {
		return err
	}

	// Apply resolution strategy
	switch resolution.Strategy {
	case "use_source":
		// Use source platform data
		return s.SyncUserData(ctx, metaUserID, resolution.SourcePlatform)
	case "use_target":
		// Use target platform data
		targetPlatform := "vignette"
		if resolution.SourcePlatform == "vignette" {
			targetPlatform = "socialink"
		}
		return s.SyncUserData(ctx, metaUserID, targetPlatform)
	case "merge":
		// Merge data from both platforms
		return s.mergePlatformData(ctx, user)
	default:
		return fmt.Errorf("unknown resolution strategy: %s", resolution.Strategy)
	}
}

// Helper methods

func (s *CrossPlatformSyncService) syncToVignette(ctx context.Context, vignetteUserID uuid.UUID, data *SyncData) error {
	if s.vignetteClient == nil {
		return fmt.Errorf("vignette client not configured")
	}
	return s.vignetteClient.SyncUserData(ctx, vignetteUserID, data)
}

func (s *CrossPlatformSyncService) syncToSocialink(ctx context.Context, socialinkUserID uuid.UUID, data *SyncData) error {
	if s.socialinkClient == nil {
		return fmt.Errorf("socialink client not configured")
	}
	return s.socialinkClient.SyncUserData(ctx, socialinkUserID, data)
}

func (s *CrossPlatformSyncService) extractProfileData(user *model.MetaUser) map[string]interface{} {
	return map[string]interface{}{
		"email":          user.Email,
		"email_verified": user.EmailVerified,
		"phone_number":   user.PhoneNumber,
		"phone_verified": user.PhoneVerified,
		"account_tier":   user.AccountTier,
		"trust_score":    user.TrustScore,
	}
}

func (s *CrossPlatformSyncService) mergePlatformData(ctx context.Context, user *model.MetaUser) error {
	// Complex merge logic - prioritize more recent data, handle conflicts
	// This is a placeholder for sophisticated merge algorithm
	user.PlatformLinks.SyncConflicts = 0
	user.PlatformLinks.LinkStatus = model.LinkStatusLinked
	return s.repo.Update(ctx, user)
}

// Supporting types

type SyncData struct {
	MetaUserID      uuid.UUID
	Email           string
	EmailVerified   bool
	PhoneNumber     *string
	PhoneVerified   bool
	PrivacySettings model.PrivacySettings
	ProfileData     map[string]interface{}
	SyncedAt        time.Time
}

type ConflictResolution struct {
	Strategy       string // "use_source", "use_target", "merge"
	SourcePlatform string
}

// Platform client interfaces (would be implemented with gRPC)

type VignetteClient struct {
	// gRPC client implementation
}

func (c *VignetteClient) SyncUserData(ctx context.Context, userID uuid.UUID, data *SyncData) error {
	// gRPC call to Vignette user service
	return nil
}

func (c *VignetteClient) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// gRPC call to Vignette user service
	return nil
}

type SocialinkClient struct {
	// gRPC client implementation
}

func (c *SocialinkClient) SyncUserData(ctx context.Context, userID uuid.UUID, data *SyncData) error {
	// gRPC call to Socialink user service
	return nil
}

func (c *SocialinkClient) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// gRPC call to Socialink user service
	return nil
}
