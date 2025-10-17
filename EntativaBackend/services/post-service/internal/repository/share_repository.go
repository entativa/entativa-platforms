package repository

import (
	"context"
	"database/sql"
	"fmt"

	"socialink/post-service/internal/model"

	"github.com/google/uuid"
)

type ShareRepository interface {
	Create(ctx context.Context, share *model.Share) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Share, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Share, error)
	GetByOriginalPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Share, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type shareRepository struct {
	db *sql.DB
}

func NewShareRepository(db *sql.DB) ShareRepository {
	return &shareRepository{db: db}
}

func (r *shareRepository) Create(ctx context.Context, share *model.Share) error {
	query := `
		INSERT INTO shares (id, user_id, original_post_id, caption, privacy, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		share.ID, share.UserID, share.OriginalPostID, share.Caption, share.Privacy, share.CreatedAt,
	).Scan(&share.CreatedAt)
}

func (r *shareRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Share, error) {
	query := `
		SELECT id, user_id, original_post_id, caption, privacy, created_at
		FROM shares
		WHERE id = $1
	`

	share := &model.Share{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&share.ID, &share.UserID, &share.OriginalPostID, &share.Caption, &share.Privacy, &share.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("share not found")
		}
		return nil, err
	}

	return share, nil
}

func (r *shareRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Share, error) {
	query := `
		SELECT id, user_id, original_post_id, caption, privacy, created_at
		FROM shares
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShares(rows)
}

func (r *shareRepository) GetByOriginalPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Share, error) {
	query := `
		SELECT id, user_id, original_post_id, caption, privacy, created_at
		FROM shares
		WHERE original_post_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShares(rows)
}

func (r *shareRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shares WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("share not found")
	}

	return nil
}

func (r *shareRepository) scanShares(rows *sql.Rows) ([]model.Share, error) {
	var shares []model.Share

	for rows.Next() {
		share := model.Share{}
		err := rows.Scan(
			&share.ID, &share.UserID, &share.OriginalPostID, &share.Caption, &share.Privacy, &share.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}

	return shares, rows.Err()
}
