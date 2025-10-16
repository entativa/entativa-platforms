package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"socialink/post-service/internal/model"

	"github.com/google/uuid"
)

type TemplateRepository interface {
	Create(ctx context.Context, template *model.TakeTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.TakeTemplate, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]model.TakeTemplate, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]model.TakeTemplate, error)
	GetFeatured(ctx context.Context, limit int) ([]model.TakeTemplate, error)
	GetTrending(ctx context.Context, limit int) ([]model.TakeTemplate, error)
	Update(ctx context.Context, template *model.TakeTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementUsage(ctx context.Context, templateID uuid.UUID) error
	Search(ctx context.Context, query string, limit, offset int) ([]model.TakeTemplate, error)
}

type templateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, template *model.TakeTemplate) error {
	query := `
		INSERT INTO takes_templates (
			id, original_take_id, creator_id, name, description, category,
			thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			usage_count, is_public, is_featured, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at, updated_at
	`

	effectsJSON, _ := json.Marshal(template.Effects)
	transitionsJSON, _ := json.Marshal(template.Transitions)
	cuesJSON, _ := json.Marshal(template.TimingCues)

	return r.db.QueryRowContext(
		ctx, query,
		template.ID, template.OriginalTakeID, template.CreatorID, template.Name,
		template.Description, template.Category, template.ThumbnailURL, template.AudioTrackID,
		effectsJSON, transitionsJSON, cuesJSON, template.UsageCount, template.IsPublic,
		template.IsFeatured, template.CreatedAt, template.UpdatedAt,
	).Scan(&template.CreatedAt, &template.UpdatedAt)
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TakeTemplate, error) {
	query := `
		SELECT id, original_take_id, creator_id, name, description, category,
			   thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			   usage_count, is_public, is_featured, created_at, updated_at
		FROM takes_templates
		WHERE id = $1
	`

	template := &model.TakeTemplate{}
	var effectsJSON, transitionsJSON, cuesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID, &template.OriginalTakeID, &template.CreatorID, &template.Name,
		&template.Description, &template.Category, &template.ThumbnailURL, &template.AudioTrackID,
		&effectsJSON, &transitionsJSON, &cuesJSON, &template.UsageCount, &template.IsPublic,
		&template.IsFeatured, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, err
	}

	// Unmarshal JSONB
	if len(effectsJSON) > 0 {
		json.Unmarshal(effectsJSON, &template.Effects)
	}
	if len(transitionsJSON) > 0 {
		json.Unmarshal(transitionsJSON, &template.Transitions)
	}
	if len(cuesJSON) > 0 {
		json.Unmarshal(cuesJSON, &template.TimingCues)
	}

	return template, nil
}

func (r *templateRepository) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]model.TakeTemplate, error) {
	query := `
		SELECT id, original_take_id, creator_id, name, description, category,
			   thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			   usage_count, is_public, is_featured, created_at, updated_at
		FROM takes_templates
		WHERE creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *templateRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]model.TakeTemplate, error) {
	query := `
		SELECT id, original_take_id, creator_id, name, description, category,
			   thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			   usage_count, is_public, is_featured, created_at, updated_at
		FROM takes_templates
		WHERE category = $1 AND is_public = TRUE
		ORDER BY usage_count DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *templateRepository) GetFeatured(ctx context.Context, limit int) ([]model.TakeTemplate, error) {
	query := `
		SELECT id, original_take_id, creator_id, name, description, category,
			   thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			   usage_count, is_public, is_featured, created_at, updated_at
		FROM takes_templates
		WHERE is_featured = TRUE AND is_public = TRUE
		ORDER BY usage_count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *templateRepository) GetTrending(ctx context.Context, limit int) ([]model.TakeTemplate, error) {
	query := `
		SELECT id, original_take_id, creator_id, name, description, category,
			   thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			   usage_count, is_public, is_featured, created_at, updated_at
		FROM takes_templates
		WHERE is_public = TRUE
		ORDER BY usage_count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *templateRepository) Update(ctx context.Context, template *model.TakeTemplate) error {
	query := `
		UPDATE takes_templates
		SET name = $1, description = $2, category = $3, is_public = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx, query,
		template.Name, template.Description, template.Category, template.IsPublic,
		template.UpdatedAt, template.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func (r *templateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM takes_templates WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func (r *templateRepository) IncrementUsage(ctx context.Context, templateID uuid.UUID) error {
	query := `UPDATE takes_templates SET usage_count = usage_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, templateID)
	return err
}

func (r *templateRepository) Search(ctx context.Context, searchQuery string, limit, offset int) ([]model.TakeTemplate, error) {
	query := `
		SELECT id, original_take_id, creator_id, name, description, category,
			   thumbnail_url, audio_track_id, effects, transitions, timing_cues,
			   usage_count, is_public, is_featured, created_at, updated_at
		FROM takes_templates
		WHERE is_public = TRUE
		AND (
			name ILIKE $1 OR
			description ILIKE $1 OR
			category ILIKE $1
		)
		ORDER BY usage_count DESC
		LIMIT $2 OFFSET $3
	`

	searchPattern := "%" + searchQuery + "%"
	rows, err := r.db.QueryContext(ctx, query, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTemplates(rows)
}

func (r *templateRepository) scanTemplates(rows *sql.Rows) ([]model.TakeTemplate, error) {
	var templates []model.TakeTemplate

	for rows.Next() {
		template := model.TakeTemplate{}
		var effectsJSON, transitionsJSON, cuesJSON []byte

		err := rows.Scan(
			&template.ID, &template.OriginalTakeID, &template.CreatorID, &template.Name,
			&template.Description, &template.Category, &template.ThumbnailURL, &template.AudioTrackID,
			&effectsJSON, &transitionsJSON, &cuesJSON, &template.UsageCount, &template.IsPublic,
			&template.IsFeatured, &template.CreatedAt, &template.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB
		if len(effectsJSON) > 0 {
			json.Unmarshal(effectsJSON, &template.Effects)
		}
		if len(transitionsJSON) > 0 {
			json.Unmarshal(transitionsJSON, &template.Transitions)
		}
		if len(cuesJSON) > 0 {
			json.Unmarshal(cuesJSON, &template.TimingCues)
		}

		templates = append(templates, template)
	}

	return templates, rows.Err()
}
