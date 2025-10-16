package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/entativa/vignette/creator-service/internal/model"
	"github.com/google/uuid"
)

type CreatorRepository struct {
	db *sql.DB
}

func NewCreatorRepository(db *sql.DB) *CreatorRepository {
	return &CreatorRepository{db: db}
}

func (r *CreatorRepository) Create(ctx context.Context, profile *model.CreatorProfile) error {
	query := `
		INSERT INTO creator_profiles (
			id, user_id, account_type, display_name, bio, category,
			badges, email, phone, website,
			monetization_enabled, monetization_status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.ExecContext(ctx, query,
		profile.ID, profile.UserID, profile.AccountType, profile.DisplayName, profile.Bio, profile.Category,
		model.BadgeArray(profile.Badges), profile.Email, profile.Phone, profile.Website,
		profile.MonetizationEnabled, profile.MonetizationStatus,
		profile.CreatedAt, profile.UpdatedAt,
	)

	return err
}

func (r *CreatorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.CreatorProfile, error) {
	query := `
		SELECT id, user_id, account_type, display_name, bio, category,
		       badges, email, phone, website,
		       monetization_enabled, monetization_status,
		       created_at, updated_at
		FROM creator_profiles
		WHERE user_id = $1
	`

	profile := &model.CreatorProfile{}
	var badges model.BadgeArray
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.AccountType, &profile.DisplayName, &profile.Bio, &profile.Category,
		&badges, &profile.Email, &profile.Phone, &profile.Website,
		&profile.MonetizationEnabled, &profile.MonetizationStatus,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("creator profile not found")
	}
	if err != nil {
		return nil, err
	}

	profile.Badges = badges
	return profile, nil
}

func (r *CreatorRepository) Update(ctx context.Context, profile *model.CreatorProfile) error {
	query := `
		UPDATE creator_profiles
		SET display_name = $1, bio = $2, category = $3,
		    email = $4, website = $5, updated_at = $6
		WHERE user_id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		profile.DisplayName, profile.Bio, profile.Category,
		profile.Email, profile.Website, profile.UpdatedAt,
		profile.UserID,
	)

	return err
}

func (r *CreatorRepository) UpdateMonetizationStatus(ctx context.Context, userID uuid.UUID, status model.MonetizationStatus, enabled bool) error {
	query := `
		UPDATE creator_profiles
		SET monetization_status = $1, monetization_enabled = $2, updated_at = NOW()
		WHERE user_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, enabled, userID)
	return err
}

func (r *CreatorRepository) AddBadge(ctx context.Context, userID uuid.UUID, badge model.CreatorBadge) error {
	// Add badge to array
	query := `
		UPDATE creator_profiles
		SET badges = badges || $1::jsonb, updated_at = NOW()
		WHERE user_id = $2
	`

	badgeJSON, _ := json.Marshal([]model.CreatorBadge{badge})
	_, err := r.db.ExecContext(ctx, query, badgeJSON, userID)
	return err
}

func (r *CreatorRepository) RecordAnalytics(ctx context.Context, analytics *model.CreatorAnalytics) error {
	query := `
		INSERT INTO creator_analytics (
			id, user_id, date,
			followers_count, followers_gained, followers_lost,
			total_likes, total_comments, total_shares, total_views, engagement_rate,
			posts_count, takes_count, stories_count,
			accounts_reached, accounts_engaged,
			age_gender_breakdown, top_locations,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (user_id, date) DO UPDATE
		SET followers_count = $4, followers_gained = $5, followers_lost = $6,
		    total_likes = $7, total_comments = $8, total_shares = $9, total_views = $10, engagement_rate = $11,
		    posts_count = $12, takes_count = $13, stories_count = $14,
		    accounts_reached = $15, accounts_engaged = $16,
		    age_gender_breakdown = $17, top_locations = $18
	`

	_, err := r.db.ExecContext(ctx, query,
		analytics.ID, analytics.UserID, analytics.Date,
		analytics.FollowersCount, analytics.FollowersGained, analytics.FollowersLost,
		analytics.TotalLikes, analytics.TotalComments, analytics.TotalShares, analytics.TotalViews, analytics.EngagementRate,
		analytics.PostsCount, analytics.TakesCount, analytics.StoriesCount,
		analytics.AccountsReached, analytics.AccountsEngaged,
		analytics.AgeGenderBreakdown, analytics.TopLocations,
		analytics.CreatedAt,
	)

	return err
}

func (r *CreatorRepository) GetAnalytics(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*model.CreatorAnalytics, error) {
	query := `
		SELECT id, user_id, date,
		       followers_count, followers_gained, followers_lost,
		       total_likes, total_comments, total_shares, total_views, engagement_rate,
		       posts_count, takes_count, stories_count,
		       accounts_reached, accounts_engaged,
		       age_gender_breakdown, top_locations,
		       created_at
		FROM creator_analytics
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	analytics := []*model.CreatorAnalytics{}
	for rows.Next() {
		a := &model.CreatorAnalytics{}
		err := rows.Scan(
			&a.ID, &a.UserID, &a.Date,
			&a.FollowersCount, &a.FollowersGained, &a.FollowersLost,
			&a.TotalLikes, &a.TotalComments, &a.TotalShares, &a.TotalViews, &a.EngagementRate,
			&a.PostsCount, &a.TakesCount, &a.StoriesCount,
			&a.AccountsReached, &a.AccountsEngaged,
			&a.AgeGenderBreakdown, &a.TopLocations,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, nil
}

func (r *CreatorRepository) RecordContentInsight(ctx context.Context, insight *model.ContentInsights) error {
	query := `
		INSERT INTO content_insights (
			id, content_id, content_type,
			impressions, reach, likes, comments, shares, saves, engagement,
			from_home, from_explore, from_profile, from_hashtags, from_other,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (content_id) DO UPDATE
		SET impressions = $4, reach = $5, likes = $6, comments = $7, shares = $8, saves = $9, engagement = $10,
		    from_home = $11, from_explore = $12, from_profile = $13, from_hashtags = $14, from_other = $15,
		    updated_at = $17
	`

	_, err := r.db.ExecContext(ctx, query,
		insight.ID, insight.ContentID, insight.ContentType,
		insight.Impressions, insight.Reach, insight.Likes, insight.Comments, insight.Shares, insight.Saves, insight.Engagement,
		insight.FromHome, insight.FromExplore, insight.FromProfile, insight.FromHashtags, insight.FromOther,
		insight.CreatedAt, insight.UpdatedAt,
	)

	return err
}

func (r *CreatorRepository) GetTopContent(ctx context.Context, userID uuid.UUID, contentType string, limit int) ([]*model.ContentInsights, error) {
	query := `
		SELECT ci.id, ci.content_id, ci.content_type,
		       ci.impressions, ci.reach, ci.likes, ci.comments, ci.shares, ci.saves,
		       ci.engagement, ci.engagement_rate, ci.created_at
		FROM content_insights ci
		WHERE ci.content_type = $1
		ORDER BY ci.engagement_rate DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, contentType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	insights := []*model.ContentInsights{}
	for rows.Next() {
		i := &model.ContentInsights{}
		err := rows.Scan(
			&i.ID, &i.ContentID, &i.ContentType,
			&i.Impressions, &i.Reach, &i.Likes, &i.Comments, &i.Shares, &i.Saves,
			&i.Engagement, &i.EngagementRate, &i.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		insights = append(insights, i)
	}

	return insights, nil
}
