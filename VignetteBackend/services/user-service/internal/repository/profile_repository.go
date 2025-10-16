package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"vignette/user-service/internal/model"
	"vignette/user-service/pkg/cache"

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
			id, user_id, category, category_type, gender, pronouns,
			link_in_bio, highlights, pinned_posts, profile_badges,
			contact_options, creator_insights, business_info,
			profile_views, profile_views_enabled, availability,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	linkInBioJSON, _ := json.Marshal(profile.LinkInBio)
	highlightsJSON, _ := json.Marshal(profile.Highlights)
	pinnedPostsJSON, _ := json.Marshal(profile.PinnedPosts)
	profileBadgesJSON, _ := json.Marshal(profile.ProfileBadges)
	contactOptionsJSON, _ := json.Marshal(profile.ContactOptions)
	creatorInsightsJSON, _ := json.Marshal(profile.CreatorInsights)
	businessInfoJSON, _ := json.Marshal(profile.BusinessInfo)
	availabilityJSON, _ := json.Marshal(profile.Availability)

	_, err := r.db.ExecContext(ctx, query,
		profile.ID, profile.UserID, profile.Category, profile.CategoryType,
		profile.Gender, profile.Pronouns, linkInBioJSON, highlightsJSON,
		pinnedPostsJSON, profileBadgesJSON, contactOptionsJSON, creatorInsightsJSON,
		businessInfoJSON, profile.ProfileViews, profile.ProfileViewsEnabled,
		availabilityJSON, profile.CreatedAt, profile.UpdatedAt,
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
		SELECT id, user_id, category, category_type, gender, pronouns,
		       link_in_bio, highlights, pinned_posts, profile_badges,
		       contact_options, creator_insights, business_info,
		       profile_views, profile_views_enabled, availability,
		       created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	var linkInBioJSON, highlightsJSON, pinnedPostsJSON, profileBadgesJSON []byte
	var contactOptionsJSON, creatorInsightsJSON, businessInfoJSON, availabilityJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Category, &profile.CategoryType,
		&profile.Gender, &profile.Pronouns, &linkInBioJSON, &highlightsJSON,
		&pinnedPostsJSON, &profileBadgesJSON, &contactOptionsJSON, &creatorInsightsJSON,
		&businessInfoJSON, &profile.ProfileViews, &profile.ProfileViewsEnabled,
		&availabilityJSON, &profile.CreatedAt, &profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Unmarshal JSON fields
	if len(linkInBioJSON) > 0 {
		json.Unmarshal(linkInBioJSON, &profile.LinkInBio)
	}
	if len(highlightsJSON) > 0 {
		json.Unmarshal(highlightsJSON, &profile.Highlights)
	}
	if len(pinnedPostsJSON) > 0 {
		json.Unmarshal(pinnedPostsJSON, &profile.PinnedPosts)
	}
	if len(profileBadgesJSON) > 0 {
		json.Unmarshal(profileBadgesJSON, &profile.ProfileBadges)
	}
	if len(contactOptionsJSON) > 0 {
		json.Unmarshal(contactOptionsJSON, &profile.ContactOptions)
	}
	if len(creatorInsightsJSON) > 0 {
		json.Unmarshal(creatorInsightsJSON, &profile.CreatorInsights)
	}
	if len(businessInfoJSON) > 0 {
		json.Unmarshal(businessInfoJSON, &profile.BusinessInfo)
	}
	if len(availabilityJSON) > 0 {
		json.Unmarshal(availabilityJSON, &profile.Availability)
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
			category = $2, category_type = $3, gender = $4, pronouns = $5,
			link_in_bio = $6, highlights = $7, pinned_posts = $8, profile_badges = $9,
			contact_options = $10, creator_insights = $11, business_info = $12,
			profile_views = $13, profile_views_enabled = $14, availability = $15,
			updated_at = $16
		WHERE user_id = $1
	`

	linkInBioJSON, _ := json.Marshal(profile.LinkInBio)
	highlightsJSON, _ := json.Marshal(profile.Highlights)
	pinnedPostsJSON, _ := json.Marshal(profile.PinnedPosts)
	profileBadgesJSON, _ := json.Marshal(profile.ProfileBadges)
	contactOptionsJSON, _ := json.Marshal(profile.ContactOptions)
	creatorInsightsJSON, _ := json.Marshal(profile.CreatorInsights)
	businessInfoJSON, _ := json.Marshal(profile.BusinessInfo)
	availabilityJSON, _ := json.Marshal(profile.Availability)

	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.Category, profile.CategoryType, profile.Gender,
		profile.Pronouns, linkInBioJSON, highlightsJSON, pinnedPostsJSON,
		profileBadgesJSON, contactOptionsJSON, creatorInsightsJSON, businessInfoJSON,
		profile.ProfileViews, profile.ProfileViewsEnabled, availabilityJSON,
		profile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("profile:user:%s", profile.UserID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// IncrementProfileViews increments the profile view count
func (r *ProfileRepository) IncrementProfileViews(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE profiles SET profile_views = profile_views + 1 WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment profile views: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("profile:user:%s", userID)
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
