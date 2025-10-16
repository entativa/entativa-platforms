package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"socialink/user-service/internal/model"
	"socialink/user-service/internal/repository"

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
	profile := &model.Profile{
		ID:        uuid.New(),
		UserID:    userID,
		Languages: []string{},
		Work:      []model.WorkExperience{},
		Education: []model.EducationEntry{},
		Visibility: &model.ProfileVisibility{
			Bio:              "public",
			Work:             "public",
			Education:        "public",
			ContactInfo:      "friends",
			RelationshipInfo: "friends",
			Hometown:         "public",
			Birthday:         "friends",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.profileRepo.CreateProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return profile, nil
}

// UpdateProfileInfo updates basic profile information
func (s *ProfileService) UpdateProfileInfo(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileInfoRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Hometown != nil {
		profile.Hometown = req.Hometown
	}
	if req.CurrentCity != nil {
		profile.CurrentCity = req.CurrentCity
	}
	if req.RelationshipStatus != nil {
		profile.RelationshipStatus = req.RelationshipStatus
	}
	if req.About != nil {
		profile.About = req.About
	}
	if req.FavoriteQuotes != nil {
		profile.FavoriteQuotes = req.FavoriteQuotes
	}
	if req.Website != nil {
		profile.Website = req.Website
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// AddWorkExperience adds work experience to profile
func (s *ProfileService) AddWorkExperience(ctx context.Context, userID uuid.UUID, req *model.AddWorkExperienceRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		endDate = &parsedDate
	}

	// Create new work entry
	workEntry := model.WorkExperience{
		ID:          uuid.New().String(),
		Company:     req.Company,
		Position:    req.Position,
		City:        req.City,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		IsCurrent:   req.IsCurrent,
	}

	// Add to profile
	if profile.Work == nil {
		profile.Work = []model.WorkExperience{}
	}
	profile.Work = append(profile.Work, workEntry)

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// RemoveWorkExperience removes work experience from profile
func (s *ProfileService) RemoveWorkExperience(ctx context.Context, userID uuid.UUID, workID string) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter out the work entry
	newWork := []model.WorkExperience{}
	for _, work := range profile.Work {
		if work.ID != workID {
			newWork = append(newWork, work)
		}
	}
	profile.Work = newWork

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// AddEducation adds education to profile
func (s *ProfileService) AddEducation(ctx context.Context, userID uuid.UUID, req *model.AddEducationRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create new education entry
	educationEntry := model.EducationEntry{
		ID:              uuid.New().String(),
		School:          req.School,
		Degree:          req.Degree,
		FieldOfStudy:    req.FieldOfStudy,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		Description:     req.Description,
		IsCurrentlyHere: req.IsCurrentlyHere,
	}

	// Add to profile
	if profile.Education == nil {
		profile.Education = []model.EducationEntry{}
	}
	profile.Education = append(profile.Education, educationEntry)

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// RemoveEducation removes education from profile
func (s *ProfileService) RemoveEducation(ctx context.Context, userID uuid.UUID, educationID string) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter out the education entry
	newEducation := []model.EducationEntry{}
	for _, edu := range profile.Education {
		if edu.ID != educationID {
			newEducation = append(newEducation, edu)
		}
	}
	profile.Education = newEducation

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateContactInfo updates contact information
func (s *ProfileService) UpdateContactInfo(ctx context.Context, userID uuid.UUID, req *model.UpdateContactInfoRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize if nil
	if profile.ContactInfo == nil {
		profile.ContactInfo = &model.ContactInfo{}
	}

	// Update fields
	if req.Email != nil {
		profile.ContactInfo.Email = req.Email
	}
	if req.PhoneNumber != nil {
		profile.ContactInfo.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		profile.ContactInfo.Address = req.Address
	}
	if req.City != nil {
		profile.ContactInfo.City = req.City
	}
	if req.State != nil {
		profile.ContactInfo.State = req.State
	}
	if req.ZipCode != nil {
		profile.ContactInfo.ZipCode = req.ZipCode
	}
	if req.Country != nil {
		profile.ContactInfo.Country = req.Country
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateSocialLinks updates social media links
func (s *ProfileService) UpdateSocialLinks(ctx context.Context, userID uuid.UUID, req *model.UpdateSocialLinksRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize if nil
	if profile.SocialLinks == nil {
		profile.SocialLinks = &model.SocialLinks{}
	}

	// Update fields
	if req.Instagram != nil {
		profile.SocialLinks.Instagram = req.Instagram
	}
	if req.Twitter != nil {
		profile.SocialLinks.Twitter = req.Twitter
	}
	if req.LinkedIn != nil {
		profile.SocialLinks.LinkedIn = req.LinkedIn
	}
	if req.YouTube != nil {
		profile.SocialLinks.YouTube = req.YouTube
	}
	if req.GitHub != nil {
		profile.SocialLinks.GitHub = req.GitHub
	}
	if req.Website != nil {
		profile.SocialLinks.Website = req.Website
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateVisibility updates profile visibility settings
func (s *ProfileService) UpdateVisibility(ctx context.Context, userID uuid.UUID, req *model.UpdateVisibilityRequest) (*model.Profile, error) {
	profile, err := s.GetOrCreateProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize if nil
	if profile.Visibility == nil {
		profile.Visibility = &model.ProfileVisibility{
			Bio:              "public",
			Work:             "public",
			Education:        "public",
			ContactInfo:      "friends",
			RelationshipInfo: "friends",
			Hometown:         "public",
			Birthday:         "friends",
		}
	}

	// Update fields
	if req.Bio != nil {
		profile.Visibility.Bio = *req.Bio
	}
	if req.Work != nil {
		profile.Visibility.Work = *req.Work
	}
	if req.Education != nil {
		profile.Visibility.Education = *req.Education
	}
	if req.ContactInfo != nil {
		profile.Visibility.ContactInfo = *req.ContactInfo
	}
	if req.RelationshipInfo != nil {
		profile.Visibility.RelationshipInfo = *req.RelationshipInfo
	}
	if req.Hometown != nil {
		profile.Visibility.Hometown = *req.Hometown
	}
	if req.Birthday != nil {
		profile.Visibility.Birthday = *req.Birthday
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
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
