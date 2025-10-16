package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vignette/post-service/internal/model"

	"github.com/google/uuid"
)

type TakesRepository interface {
	Create(ctx context.Context, take *model.Take) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Take, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Take, error)
	GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]model.Take, *string, error)
	GetByHashtag(ctx context.Context, hashtag string, limit, offset int) ([]model.Take, error)
	GetByTrendID(ctx context.Context, trendID uuid.UUID, limit, offset int) ([]model.Take, error)
	GetByTemplateID(ctx context.Context, templateID uuid.UUID, limit, offset int) ([]model.Take, error)
	Update(ctx context.Context, take *model.Take) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementViews(ctx context.Context, takeID uuid.UUID) error
	IncrementLikes(ctx context.Context, takeID uuid.UUID) error
	DecrementLikes(ctx context.Context, takeID uuid.UUID) error
	IncrementComments(ctx context.Context, takeID uuid.UUID) error
	DecrementComments(ctx context.Context, takeID uuid.UUID) error
	IncrementShares(ctx context.Context, takeID uuid.UUID) error
	IncrementSaves(ctx context.Context, takeID uuid.UUID) error
	DecrementSaves(ctx context.Context, takeID uuid.UUID) error
	IncrementRemixes(ctx context.Context, takeID uuid.UUID) error
	GetTrending(ctx context.Context, limit int, timeWindow time.Duration) ([]model.Take, error)
}

type takesRepository struct {
	db *sql.DB
}

func NewTakesRepository(db *sql.DB) TakesRepository {
	return &takesRepository{db: db}
}

func (r *takesRepository) Create(ctx context.Context, take *model.Take) error {
	query := `
		INSERT INTO takes (
			id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25
		)
		RETURNING created_at, updated_at
	`

	hashtagsJSON, _ := json.Marshal(take.Hashtags)
	locationJSON, _ := json.Marshal(take.Location)
	taggedJSON, _ := json.Marshal(take.TaggedUserIDs)

	return r.db.QueryRowContext(
		ctx, query,
		take.ID, take.UserID, take.Caption, take.MediaID, take.AudioTrackID, take.Duration,
		take.ThumbnailURL, hashtagsJSON, take.FilterUsed, locationJSON, taggedJSON,
		take.TemplateID, take.TrendID, take.HasBTT, take.ViewsCount, take.LikesCount,
		take.CommentsCount, take.SharesCount, take.SavesCount, take.RemixCount,
		take.CommentsEnabled, take.RemixEnabled, take.IsSponsored, take.CreatedAt, take.UpdatedAt,
	).Scan(&take.CreatedAt, &take.UpdatedAt)
}

func (r *takesRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Take, error) {
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE id = $1 AND deleted_at IS NULL
	`

	take := &model.Take{}
	var hashtagsJSON, locationJSON, taggedJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&take.ID, &take.UserID, &take.Caption, &take.MediaID, &take.AudioTrackID,
		&take.Duration, &take.ThumbnailURL, &hashtagsJSON, &take.FilterUsed, &locationJSON,
		&taggedJSON, &take.TemplateID, &take.TrendID, &take.HasBTT, &take.ViewsCount,
		&take.LikesCount, &take.CommentsCount, &take.SharesCount, &take.SavesCount,
		&take.RemixCount, &take.CommentsEnabled, &take.RemixEnabled, &take.IsSponsored,
		&take.CreatedAt, &take.UpdatedAt, &take.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("take not found")
		}
		return nil, err
	}

	// Unmarshal JSONB
	if len(hashtagsJSON) > 0 {
		json.Unmarshal(hashtagsJSON, &take.Hashtags)
	}
	if len(locationJSON) > 0 && string(locationJSON) != "null" {
		var location model.Location
		if json.Unmarshal(locationJSON, &location) == nil {
			take.Location = &location
		}
	}
	if len(taggedJSON) > 0 {
		json.Unmarshal(taggedJSON, &take.TaggedUserIDs)
	}

	return take, nil
}

func (r *takesRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Take, error) {
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTakes(rows)
}

func (r *takesRepository) GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]model.Take, *string, error) {
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	
	if cursor != "" {
		query += ` AND created_at < (SELECT created_at FROM takes WHERE id = $1)`
		var cursorID uuid.UUID
		if err := cursorID.UnmarshalText([]byte(cursor)); err == nil {
			args = append(args, cursorID)
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	takes, err := r.scanTakes(rows)
	if err != nil {
		return nil, nil, err
	}

	var nextCursor *string
	if len(takes) > limit {
		takes = takes[:limit]
		cursorStr := takes[limit-1].ID.String()
		nextCursor = &cursorStr
	}

	return takes, nextCursor, nil
}

func (r *takesRepository) GetByHashtag(ctx context.Context, hashtag string, limit, offset int) ([]model.Take, error) {
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE deleted_at IS NULL
		AND hashtags ? $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, hashtag, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTakes(rows)
}

func (r *takesRepository) GetByTrendID(ctx context.Context, trendID uuid.UUID, limit, offset int) ([]model.Take, error) {
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE trend_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, trendID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTakes(rows)
}

func (r *takesRepository) GetByTemplateID(ctx context.Context, templateID uuid.UUID, limit, offset int) ([]model.Take, error) {
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE template_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, templateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTakes(rows)
}

func (r *takesRepository) Update(ctx context.Context, take *model.Take) error {
	query := `
		UPDATE takes
		SET caption = $1, hashtags = $2, location = $3, has_btt = $4,
			comments_enabled = $5, remix_enabled = $6, updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL
	`

	hashtagsJSON, _ := json.Marshal(take.Hashtags)
	locationJSON, _ := json.Marshal(take.Location)

	result, err := r.db.ExecContext(
		ctx, query,
		take.Caption, hashtagsJSON, locationJSON, take.HasBTT,
		take.CommentsEnabled, take.RemixEnabled, take.UpdatedAt, take.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("take not found")
	}

	return nil
}

func (r *takesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE takes SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("take not found")
	}

	return nil
}

func (r *takesRepository) IncrementViews(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) IncrementLikes(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET likes_count = likes_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) DecrementLikes(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET likes_count = GREATEST(0, likes_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) IncrementComments(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET comments_count = comments_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) DecrementComments(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET comments_count = GREATEST(0, comments_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) IncrementShares(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET shares_count = shares_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) IncrementSaves(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET saves_count = saves_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) DecrementSaves(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET saves_count = GREATEST(0, saves_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) IncrementRemixes(ctx context.Context, takeID uuid.UUID) error {
	query := `UPDATE takes SET remix_count = remix_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, takeID)
	return err
}

func (r *takesRepository) GetTrending(ctx context.Context, limit int, timeWindow time.Duration) ([]model.Take, error) {
	// Algorithm: weighted engagement (views + likes*3 + shares*5 + remixes*10)
	query := `
		SELECT id, user_id, caption, media_id, audio_track_id, duration, thumbnail_url,
			   hashtags, filter_used, location, tagged_user_ids, template_id, trend_id,
			   has_btt, views_count, likes_count, comments_count, shares_count, saves_count,
			   remix_count, comments_enabled, remix_enabled, is_sponsored, created_at, updated_at, deleted_at
		FROM takes
		WHERE deleted_at IS NULL
		AND created_at > $1
		ORDER BY (
			views_count / 100 + 
			likes_count * 3 + 
			comments_count * 2 +
			shares_count * 5 + 
			saves_count * 4 +
			remix_count * 10
		) DESC
		LIMIT $2
	`

	since := time.Now().Add(-timeWindow)
	rows, err := r.db.QueryContext(ctx, query, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTakes(rows)
}

func (r *takesRepository) scanTakes(rows *sql.Rows) ([]model.Take, error) {
	var takes []model.Take

	for rows.Next() {
		take := model.Take{}
		var hashtagsJSON, locationJSON, taggedJSON []byte

		err := rows.Scan(
			&take.ID, &take.UserID, &take.Caption, &take.MediaID, &take.AudioTrackID,
			&take.Duration, &take.ThumbnailURL, &hashtagsJSON, &take.FilterUsed, &locationJSON,
			&taggedJSON, &take.TemplateID, &take.TrendID, &take.HasBTT, &take.ViewsCount,
			&take.LikesCount, &take.CommentsCount, &take.SharesCount, &take.SavesCount,
			&take.RemixCount, &take.CommentsEnabled, &take.RemixEnabled, &take.IsSponsored,
			&take.CreatedAt, &take.UpdatedAt, &take.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB
		if len(hashtagsJSON) > 0 {
			json.Unmarshal(hashtagsJSON, &take.Hashtags)
		}
		if len(locationJSON) > 0 && string(locationJSON) != "null" {
			var location model.Location
			if json.Unmarshal(locationJSON, &location) == nil {
				take.Location = &location
			}
		}
		if len(taggedJSON) > 0 {
			json.Unmarshal(taggedJSON, &take.TaggedUserIDs)
		}

		takes = append(takes, take)
	}

	return takes, rows.Err()
}
