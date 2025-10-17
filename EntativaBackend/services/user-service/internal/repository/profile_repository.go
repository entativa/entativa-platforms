package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"socialink/user-service/internal/model"
	"socialink/user-service/pkg/cache"

	"github.com/google/uuid"
)

var (
	ErrProfileNotFound = fmt.Errorf("profile not found")
)

type ProfileRepository struct {
	db    *sql.DB
	cache cache.Cache
}

func NewProfileRepository(db *sql.DB, cache cache.Cache) *ProfileRepository {
	return &ProfileRepository{
		db:    db,
		cache: cache,
	}
}

// CreateProfile creates a new profile for a user
func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *model.Profile) error {
	query := `
		INSERT INTO profiles (
			id, user_id, hometown, current_city, relationship_status,
			languages, interested_in, work, education, contact_info,
			about, favorite_quotes, hobbies, website, social_links,
			featured_photos, visibility, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	workJSON, _ := json.Marshal(profile.Work)
	educationJSON, _ := json.Marshal(profile.Education)
	contactInfoJSON, _ := json.Marshal(profile.ContactInfo)
	socialLinksJSON, _ := json.Marshal(profile.SocialLinks)
	visibilityJSON, _ := json.Marshal(profile.Visibility)
	languagesJSON, _ := json.Marshal(profile.Languages)
	interestedInJSON, _ := json.Marshal(profile.InterestedIn)
	hobbiesJSON, _ := json.Marshal(profile.Hobbies)
	featuredPhotosJSON, _ := json.Marshal(profile.FeaturedPhotos)

	_, err := r.db.ExecContext(ctx, query,
		profile.ID, profile.UserID, profile.Hometown, profile.CurrentCity,
		profile.RelationshipStatus, languagesJSON, interestedInJSON, workJSON,
		educationJSON, contactInfoJSON, profile.About, profile.FavoriteQuotes,
		hobbiesJSON, profile.Website, socialLinksJSON, featuredPhotosJSON,
		visibilityJSON, profile.CreatedAt, profile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	// Cache the new profile
	cacheKey := fmt.Sprintf("profile:user:%s", profile.UserID)
	r.cache.Set(ctx, cacheKey, profile, 10*time.Minute)

	return nil
}

// GetProfileByUserID retrieves a profile by user ID
func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("profile:user:%s", userID)
	var profile model.Profile
	if err := r.cache.Get(ctx, cacheKey, &profile); err == nil {
		return &profile, nil
	}

	query := `
		SELECT id, user_id, hometown, current_city, relationship_status,
		       languages, interested_in, work, education, contact_info,
		       about, favorite_quotes, hobbies, website, social_links,
		       featured_photos, visibility, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	var workJSON, educationJSON, contactInfoJSON, socialLinksJSON, visibilityJSON []byte
	var languagesJSON, interestedInJSON, hobbiesJSON, featuredPhotosJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Hometown, &profile.CurrentCity,
		&profile.RelationshipStatus, &languagesJSON, &interestedInJSON, &workJSON,
		&educationJSON, &contactInfoJSON, &profile.About, &profile.FavoriteQuotes,
		&hobbiesJSON, &profile.Website, &socialLinksJSON, &featuredPhotosJSON,
		&visibilityJSON, &profile.CreatedAt, &profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Unmarshal JSON fields
	if len(workJSON) > 0 {
		json.Unmarshal(workJSON, &profile.Work)
	}
	if len(educationJSON) > 0 {
		json.Unmarshal(educationJSON, &profile.Education)
	}
	if len(contactInfoJSON) > 0 {
		json.Unmarshal(contactInfoJSON, &profile.ContactInfo)
	}
	if len(socialLinksJSON) > 0 {
		json.Unmarshal(socialLinksJSON, &profile.SocialLinks)
	}
	if len(visibilityJSON) > 0 {
		json.Unmarshal(visibilityJSON, &profile.Visibility)
	}
	if len(languagesJSON) > 0 {
		json.Unmarshal(languagesJSON, &profile.Languages)
	}
	if len(interestedInJSON) > 0 {
		json.Unmarshal(interestedInJSON, &profile.InterestedIn)
	}
	if len(hobbiesJSON) > 0 {
		json.Unmarshal(hobbiesJSON, &profile.Hobbies)
	}
	if len(featuredPhotosJSON) > 0 {
		json.Unmarshal(featuredPhotosJSON, &profile.FeaturedPhotos)
	}

	// Cache the profile
	r.cache.Set(ctx, cacheKey, &profile, 10*time.Minute)

	return &profile, nil
}

// UpdateProfile updates an existing profile
func (r *ProfileRepository) UpdateProfile(ctx context.Context, profile *model.Profile) error {
	profile.UpdatedAt = time.Now()

	query := `
		UPDATE profiles SET
			hometown = $2, current_city = $3, relationship_status = $4,
			languages = $5, interested_in = $6, work = $7, education = $8,
			contact_info = $9, about = $10, favorite_quotes = $11, hobbies = $12,
			website = $13, social_links = $14, featured_photos = $15,
			visibility = $16, updated_at = $17
		WHERE user_id = $1
	`

	workJSON, _ := json.Marshal(profile.Work)
	educationJSON, _ := json.Marshal(profile.Education)
	contactInfoJSON, _ := json.Marshal(profile.ContactInfo)
	socialLinksJSON, _ := json.Marshal(profile.SocialLinks)
	visibilityJSON, _ := json.Marshal(profile.Visibility)
	languagesJSON, _ := json.Marshal(profile.Languages)
	interestedInJSON, _ := json.Marshal(profile.InterestedIn)
	hobbiesJSON, _ := json.Marshal(profile.Hobbies)
	featuredPhotosJSON, _ := json.Marshal(profile.FeaturedPhotos)

	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.Hometown, profile.CurrentCity, profile.RelationshipStatus,
		languagesJSON, interestedInJSON, workJSON, educationJSON, contactInfoJSON,
		profile.About, profile.FavoriteQuotes, hobbiesJSON, profile.Website,
		socialLinksJSON, featuredPhotosJSON, visibilityJSON, profile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("profile:user:%s", profile.UserID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// DeleteProfile soft deletes a profile
func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM profiles WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("profile:user:%s", userID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}
