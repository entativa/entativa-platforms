package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"vignette/post-service/internal/model"

	"github.com/google/uuid"
)

type BTTRepository interface {
	Create(ctx context.Context, btt *model.BehindTheTakes) error
	GetByTakeID(ctx context.Context, takeID uuid.UUID) (*model.BehindTheTakes, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.BehindTheTakes, error)
	Update(ctx context.Context, btt *model.BehindTheTakes) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementViews(ctx context.Context, bttID uuid.UUID) error
	IncrementLikes(ctx context.Context, bttID uuid.UUID) error
	GetTrending(ctx context.Context, limit int) ([]model.BehindTheTakes, error)
}

type bttRepository struct {
	db *sql.DB
}

func NewBTTRepository(db *sql.DB) BTTRepository {
	return &bttRepository{db: db}
}

func (r *bttRepository) Create(ctx context.Context, btt *model.BehindTheTakes) error {
	query := `
		INSERT INTO behind_the_takes (
			id, take_id, user_id, media_ids, description, steps, equipment,
			software, tips, views_count, likes_count, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING created_at, updated_at
	`

	mediaIDsJSON, _ := json.Marshal(btt.MediaIDs)
	stepsJSON, _ := json.Marshal(btt.Steps)
	equipmentJSON, _ := json.Marshal(btt.Equipment)
	softwareJSON, _ := json.Marshal(btt.Software)
	tipsJSON, _ := json.Marshal(btt.Tips)

	return r.db.QueryRowContext(
		ctx, query,
		btt.ID, btt.TakeID, btt.UserID, mediaIDsJSON, btt.Description, stepsJSON,
		equipmentJSON, softwareJSON, tipsJSON, btt.ViewsCount, btt.LikesCount,
		btt.CreatedAt, btt.UpdatedAt,
	).Scan(&btt.CreatedAt, &btt.UpdatedAt)
}

func (r *bttRepository) GetByTakeID(ctx context.Context, takeID uuid.UUID) (*model.BehindTheTakes, error) {
	query := `
		SELECT id, take_id, user_id, media_ids, description, steps, equipment,
			   software, tips, views_count, likes_count, created_at, updated_at
		FROM behind_the_takes
		WHERE take_id = $1
	`

	btt := &model.BehindTheTakes{}
	var mediaIDsJSON, stepsJSON, equipmentJSON, softwareJSON, tipsJSON []byte

	err := r.db.QueryRowContext(ctx, query, takeID).Scan(
		&btt.ID, &btt.TakeID, &btt.UserID, &mediaIDsJSON, &btt.Description, &stepsJSON,
		&equipmentJSON, &softwareJSON, &tipsJSON, &btt.ViewsCount, &btt.LikesCount,
		&btt.CreatedAt, &btt.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No BTT for this Take
		}
		return nil, err
	}

	// Unmarshal JSONB
	if len(mediaIDsJSON) > 0 {
		json.Unmarshal(mediaIDsJSON, &btt.MediaIDs)
	}
	if len(stepsJSON) > 0 {
		json.Unmarshal(stepsJSON, &btt.Steps)
	}
	if len(equipmentJSON) > 0 {
		json.Unmarshal(equipmentJSON, &btt.Equipment)
	}
	if len(softwareJSON) > 0 {
		json.Unmarshal(softwareJSON, &btt.Software)
	}
	if len(tipsJSON) > 0 {
		json.Unmarshal(tipsJSON, &btt.Tips)
	}

	return btt, nil
}

func (r *bttRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.BehindTheTakes, error) {
	query := `
		SELECT id, take_id, user_id, media_ids, description, steps, equipment,
			   software, tips, views_count, likes_count, created_at, updated_at
		FROM behind_the_takes
		WHERE id = $1
	`

	btt := &model.BehindTheTakes{}
	var mediaIDsJSON, stepsJSON, equipmentJSON, softwareJSON, tipsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&btt.ID, &btt.TakeID, &btt.UserID, &mediaIDsJSON, &btt.Description, &stepsJSON,
		&equipmentJSON, &softwareJSON, &tipsJSON, &btt.ViewsCount, &btt.LikesCount,
		&btt.CreatedAt, &btt.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("BTT not found")
		}
		return nil, err
	}

	// Unmarshal JSONB
	if len(mediaIDsJSON) > 0 {
		json.Unmarshal(mediaIDsJSON, &btt.MediaIDs)
	}
	if len(stepsJSON) > 0 {
		json.Unmarshal(stepsJSON, &btt.Steps)
	}
	if len(equipmentJSON) > 0 {
		json.Unmarshal(equipmentJSON, &btt.Equipment)
	}
	if len(softwareJSON) > 0 {
		json.Unmarshal(softwareJSON, &btt.Software)
	}
	if len(tipsJSON) > 0 {
		json.Unmarshal(tipsJSON, &btt.Tips)
	}

	return btt, nil
}

func (r *bttRepository) Update(ctx context.Context, btt *model.BehindTheTakes) error {
	query := `
		UPDATE behind_the_takes
		SET description = $1, steps = $2, equipment = $3, software = $4, tips = $5, updated_at = $6
		WHERE id = $7
	`

	stepsJSON, _ := json.Marshal(btt.Steps)
	equipmentJSON, _ := json.Marshal(btt.Equipment)
	softwareJSON, _ := json.Marshal(btt.Software)
	tipsJSON, _ := json.Marshal(btt.Tips)

	result, err := r.db.ExecContext(
		ctx, query,
		btt.Description, stepsJSON, equipmentJSON, softwareJSON, tipsJSON, btt.UpdatedAt, btt.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("BTT not found")
	}

	return nil
}

func (r *bttRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM behind_the_takes WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("BTT not found")
	}

	return nil
}

func (r *bttRepository) IncrementViews(ctx context.Context, bttID uuid.UUID) error {
	query := `UPDATE behind_the_takes SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, bttID)
	return err
}

func (r *bttRepository) IncrementLikes(ctx context.Context, bttID uuid.UUID) error {
	query := `UPDATE behind_the_takes SET likes_count = likes_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, bttID)
	return err
}

func (r *bttRepository) GetTrending(ctx context.Context, limit int) ([]model.BehindTheTakes, error) {
	query := `
		SELECT id, take_id, user_id, media_ids, description, steps, equipment,
			   software, tips, views_count, likes_count, created_at, updated_at
		FROM behind_the_takes
		ORDER BY (views_count + likes_count * 5) DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bttList []model.BehindTheTakes
	for rows.Next() {
		btt := model.BehindTheTakes{}
		var mediaIDsJSON, stepsJSON, equipmentJSON, softwareJSON, tipsJSON []byte

		err := rows.Scan(
			&btt.ID, &btt.TakeID, &btt.UserID, &mediaIDsJSON, &btt.Description, &stepsJSON,
			&equipmentJSON, &softwareJSON, &tipsJSON, &btt.ViewsCount, &btt.LikesCount,
			&btt.CreatedAt, &btt.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB
		if len(mediaIDsJSON) > 0 {
			json.Unmarshal(mediaIDsJSON, &btt.MediaIDs)
		}
		if len(stepsJSON) > 0 {
			json.Unmarshal(stepsJSON, &btt.Steps)
		}
		if len(equipmentJSON) > 0 {
			json.Unmarshal(equipmentJSON, &btt.Equipment)
		}
		if len(softwareJSON) > 0 {
			json.Unmarshal(softwareJSON, &btt.Software)
		}
		if len(tipsJSON) > 0 {
			json.Unmarshal(tipsJSON, &btt.Tips)
		}

		bttList = append(bttList, btt)
	}

	return bttList, rows.Err()
}
