package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"socialink/post-service/internal/model"

	"github.com/google/uuid"
)

type TrendRepository interface {
	Create(ctx context.Context, trend *model.TakeTrend) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.TakeTrend, error)
	GetByKeyword(ctx context.Context, keyword string) (*model.TakeTrend, error)
	GetActive(ctx context.Context, limit int) ([]model.TakeTrend, error)
	GetFeatured(ctx context.Context, limit int) ([]model.TakeTrend, error)
	Update(ctx context.Context, trend *model.TakeTrend) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementParticipants(ctx context.Context, trendID uuid.UUID) error
	IncrementViews(ctx context.Context, trendID uuid.UUID) error
	SearchByKeyword(ctx context.Context, keyword string, limit int) ([]model.TakeTrend, error)
}

type trendRepository struct {
	db *sql.DB
}

func NewTrendRepository(db *sql.DB) TrendRepository {
	return &trendRepository{db: db}
}

func (r *trendRepository) Create(ctx context.Context, trend *model.TakeTrend) error {
	query := `
		INSERT INTO takes_trends (
			id, keyword, originator_id, origin_take_id, display_name, description,
			category, thumbnail_url, audio_track_id, participant_count, views_count,
			is_active, is_featured, started_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at, updated_at
	`

	// Normalize keyword to lowercase for case-insensitive matching
	keyword := strings.ToLower(strings.TrimSpace(trend.Keyword))

	return r.db.QueryRowContext(
		ctx, query,
		trend.ID, keyword, trend.OriginatorID, trend.OriginTakeID, trend.DisplayName,
		trend.Description, trend.Category, trend.ThumbnailURL, trend.AudioTrackID,
		trend.ParticipantCount, trend.ViewsCount, trend.IsActive, trend.IsFeatured,
		trend.StartedAt, trend.CreatedAt, trend.UpdatedAt,
	).Scan(&trend.CreatedAt, &trend.UpdatedAt)
}

func (r *trendRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TakeTrend, error) {
	query := `
		SELECT id, keyword, originator_id, origin_take_id, display_name, description,
			   category, thumbnail_url, audio_track_id, participant_count, views_count,
			   is_active, is_featured, started_at, peak_at, expires_at, created_at, updated_at
		FROM takes_trends
		WHERE id = $1
	`

	trend := &model.TakeTrend{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trend.ID, &trend.Keyword, &trend.OriginatorID, &trend.OriginTakeID, &trend.DisplayName,
		&trend.Description, &trend.Category, &trend.ThumbnailURL, &trend.AudioTrackID,
		&trend.ParticipantCount, &trend.ViewsCount, &trend.IsActive, &trend.IsFeatured,
		&trend.StartedAt, &trend.PeakAt, &trend.ExpiresAt, &trend.CreatedAt, &trend.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trend not found")
		}
		return nil, err
	}

	return trend, nil
}

func (r *trendRepository) GetByKeyword(ctx context.Context, keyword string) (*model.TakeTrend, error) {
	// Case-insensitive search
	query := `
		SELECT id, keyword, originator_id, origin_take_id, display_name, description,
			   category, thumbnail_url, audio_track_id, participant_count, views_count,
			   is_active, is_featured, started_at, peak_at, expires_at, created_at, updated_at
		FROM takes_trends
		WHERE LOWER(keyword) = LOWER($1) AND is_active = TRUE
	`

	trend := &model.TakeTrend{}
	err := r.db.QueryRowContext(ctx, query, keyword).Scan(
		&trend.ID, &trend.Keyword, &trend.OriginatorID, &trend.OriginTakeID, &trend.DisplayName,
		&trend.Description, &trend.Category, &trend.ThumbnailURL, &trend.AudioTrackID,
		&trend.ParticipantCount, &trend.ViewsCount, &trend.IsActive, &trend.IsFeatured,
		&trend.StartedAt, &trend.PeakAt, &trend.ExpiresAt, &trend.CreatedAt, &trend.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Trend doesn't exist yet
		}
		return nil, err
	}

	return trend, nil
}

func (r *trendRepository) GetActive(ctx context.Context, limit int) ([]model.TakeTrend, error) {
	query := `
		SELECT id, keyword, originator_id, origin_take_id, display_name, description,
			   category, thumbnail_url, audio_track_id, participant_count, views_count,
			   is_active, is_featured, started_at, peak_at, expires_at, created_at, updated_at
		FROM takes_trends
		WHERE is_active = TRUE
		AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY participant_count DESC, views_count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTrends(rows)
}

func (r *trendRepository) GetFeatured(ctx context.Context, limit int) ([]model.TakeTrend, error) {
	query := `
		SELECT id, keyword, originator_id, origin_take_id, display_name, description,
			   category, thumbnail_url, audio_track_id, participant_count, views_count,
			   is_active, is_featured, started_at, peak_at, expires_at, created_at, updated_at
		FROM takes_trends
		WHERE is_featured = TRUE AND is_active = TRUE
		ORDER BY participant_count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTrends(rows)
}

func (r *trendRepository) Update(ctx context.Context, trend *model.TakeTrend) error {
	query := `
		UPDATE takes_trends
		SET display_name = $1, description = $2, is_active = $3, is_featured = $4,
			peak_at = $5, expires_at = $6, updated_at = $7
		WHERE id = $8
	`

	result, err := r.db.ExecContext(
		ctx, query,
		trend.DisplayName, trend.Description, trend.IsActive, trend.IsFeatured,
		trend.PeakAt, trend.ExpiresAt, trend.UpdatedAt, trend.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("trend not found")
	}

	return nil
}

func (r *trendRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM takes_trends WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("trend not found")
	}

	return nil
}

func (r *trendRepository) IncrementParticipants(ctx context.Context, trendID uuid.UUID) error {
	query := `UPDATE takes_trends SET participant_count = participant_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, trendID)
	return err
}

func (r *trendRepository) IncrementViews(ctx context.Context, trendID uuid.UUID) error {
	query := `UPDATE takes_trends SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, trendID)
	return err
}

func (r *trendRepository) SearchByKeyword(ctx context.Context, keyword string, limit int) ([]model.TakeTrend, error) {
	query := `
		SELECT id, keyword, originator_id, origin_take_id, display_name, description,
			   category, thumbnail_url, audio_track_id, participant_count, views_count,
			   is_active, is_featured, started_at, peak_at, expires_at, created_at, updated_at
		FROM takes_trends
		WHERE is_active = TRUE
		AND (
			LOWER(keyword) LIKE LOWER($1) OR
			LOWER(display_name) LIKE LOWER($1)
		)
		ORDER BY participant_count DESC
		LIMIT $2
	`

	searchPattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTrends(rows)
}

func (r *trendRepository) scanTrends(rows *sql.Rows) ([]model.TakeTrend, error) {
	var trends []model.TakeTrend

	for rows.Next() {
		trend := model.TakeTrend{}
		err := rows.Scan(
			&trend.ID, &trend.Keyword, &trend.OriginatorID, &trend.OriginTakeID, &trend.DisplayName,
			&trend.Description, &trend.Category, &trend.ThumbnailURL, &trend.AudioTrackID,
			&trend.ParticipantCount, &trend.ViewsCount, &trend.IsActive, &trend.IsFeatured,
			&trend.StartedAt, &trend.PeakAt, &trend.ExpiresAt, &trend.CreatedAt, &trend.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		trends = append(trends, trend)
	}

	return trends, rows.Err()
}
