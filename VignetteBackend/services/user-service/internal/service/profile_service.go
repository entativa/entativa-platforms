package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vignette/user-service/internal/model"
	"vignette/user-service/internal/repository"

	"github.com/google/uuid"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
	userRepo    *repository.UserRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository, userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

// GetOrCreateProfile retrieves or creates a profile for a user
func (s *ProfileService) GetOrCreateProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	// Try to get existing profile
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err == nil {
		return profile, nil
	}

	// If not found, create new profile with defaults
	if errors.Is(err, repository.ErrProfileNotFound) {
		return s.createDefaultProfile(ctx, userID)
	}

	return nil, err
}

// createDefaultProfile creates a profile with default settings
func (s *ProfileService) createDefaultProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	now := time.Now()
	personal := "personal"
	profile := &model.Profile{
		ID:                  uuid.New(),
		UserID:              userID,
		Category:            &personal,
		LinkInBio:           []model.LinkInBio{},
		Highlights:          []model.StoryHighlight{},
		PinnedPosts:         []string{},
		ProfileBadges:       []string{},
		ProfileViews:        0,
		ProfileViewsEnabled: true,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.profileRepo.CreateProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return profile, nil
}

// UpdateProfileExtended updates extended profile information
func (s *ProfileService) UpdateProfileExtended(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileExtendedRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Category != nil {
		profile.Category = req.Category
	}
	if req.CategoryType != nil {
		profile.CategoryType = req.CategoryType
	}
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.Pronouns != nil {
		profile.Pronouns = req.Pronouns
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// AddLinkInBio adds a link to bio
func (s *ProfileService) AddLinkInBio(ctx context.Context, userID uuid.UUID, req *model.AddLinkInBioRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create new link
	link := model.LinkInBio{
		ID:    uuid.New().String(),
		Title: req.Title,
		URL:   req.URL,
		Order: len(profile.LinkInBio),
	}

	// Add to profile
	if profile.LinkInBio == nil {
		profile.LinkInBio = []model.LinkInBio{}
	}
	profile.LinkInBio = append(profile.LinkInBio, link)

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// RemoveLinkInBio removes a link from bio
func (s *ProfileService) RemoveLinkInBio(ctx context.Context, userID uuid.UUID, linkID string) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter out the link
	newLinks := []model.LinkInBio{}
	for _, link := range profile.LinkInBio {
		if link.ID != linkID {
			newLinks = append(newLinks, link)
		}
	}
	profile.LinkInBio = newLinks

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// AddHighlight adds a story highlight
func (s *ProfileService) AddHighlight(ctx context.Context, userID uuid.UUID, req *model.AddHighlightRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create new highlight
	highlight := model.StoryHighlight{
		ID:        uuid.New().String(),
		Title:     req.Title,
		CoverURL:  req.CoverURL,
		StoryIDs:  req.StoryIDs,
		CreatedAt: time.Now(),
		Order:     len(profile.Highlights),
	}

	// Add to profile
	if profile.Highlights == nil {
		profile.Highlights = []model.StoryHighlight{}
	}
	profile.Highlights = append(profile.Highlights, highlight)

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// RemoveHighlight removes a story highlight
func (s *ProfileService) RemoveHighlight(ctx context.Context, userID uuid.UUID, highlightID string) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter out the highlight
	newHighlights := []model.StoryHighlight{}
	for _, highlight := range profile.Highlights {
		if highlight.ID != highlightID {
			newHighlights = append(newHighlights, highlight)
		}
	}
	profile.Highlights = newHighlights

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateContactOptions updates business contact options
func (s *ProfileService) UpdateContactOptions(ctx context.Context, userID uuid.UUID, req *model.UpdateContactOptionsRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize if nil
	if profile.ContactOptions == nil {
		profile.ContactOptions = &model.ContactOptions{}
	}

	// Update fields
	if req.Email != nil {
		profile.ContactOptions.Email = req.Email
	}
	if req.PhoneNumber != nil {
		profile.ContactOptions.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		// Split address into street/city/zip if needed
		profile.ContactOptions.AddressStreet = req.Address
	}
	if req.ShowEmail != nil {
		profile.ContactOptions.ShowEmail = *req.ShowEmail
	}
	if req.ShowPhone != nil {
		profile.ContactOptions.ShowPhone = *req.ShowPhone
	}
	if req.ShowAddress != nil {
		profile.ContactOptions.ShowAddress = *req.ShowAddress
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// EnableCreatorAccount switches account to creator account
func (s *ProfileService) EnableCreatorAccount(ctx context.Context, userID uuid.UUID, req *model.EnableCreatorAccountRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize creator insights
	if profile.CreatorInsights == nil {
		profile.CreatorInsights = &model.CreatorInsights{
			IsCreatorAccount: true,
			EnabledDate:      time.Now(),
			TopAudiences:     []string{},
		}
	} else {
		profile.CreatorInsights.IsCreatorAccount = true
		profile.CreatorInsights.EnabledDate = time.Now()
	}

	// Update category
	creator := "creator"
	profile.Category = &creator
	profile.CategoryType = &req.Category

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// EnableBusinessAccount switches account to business account
func (s *ProfileService) EnableBusinessAccount(ctx context.Context, userID uuid.UUID, req *model.EnableBusinessAccountRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize business info
	profile.BusinessInfo = &model.BusinessInfo{
		IsBusinessAccount: true,
		BusinessCategory:  req.BusinessCategory,
		BusinessEmail:     req.BusinessEmail,
		BusinessPhone:     req.BusinessPhone,
		BusinessAddress:   req.BusinessAddress,
		PriceRange:        req.PriceRange,
	}

	// Update category
	business := "business"
	profile.Category = &business

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateAvailability updates availability status
func (s *ProfileService) UpdateAvailability(ctx context.Context, userID uuid.UUID, req *model.UpdateAvailabilityRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	message := ""
	if req.Message != nil {
		message = *req.Message
	}

	profile.Availability = &model.Availability{
		Status:  req.Status,
		Message: message,
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// IncrementProfileViews increments profile view count
func (s *ProfileService) IncrementProfileViews(ctx context.Context, userID uuid.UUID) error {
	return s.profileRepo.IncrementProfileViews(ctx, userID)
}

// GetProfileWithUser gets profile along with user information
func (s *ProfileService) GetProfileWithUser(ctx context.Context, userID uuid.UUID) (*model.ProfileResponse, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := profile.ToProfileResponse()
	response.User = user.ToUserResponse()

	return response, nil
}
