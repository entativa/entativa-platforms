package repository

import (
	"context"
	"database/sql"
	"fmt"

	"socialink/post-service/internal/model"

	"github.com/google/uuid"
)

type SaveRepository interface {
	Create(ctx context.Context, save *model.Save) error
	Delete(ctx context.Context, userID, postID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID, collection *string, limit, offset int) ([]model.Save, error)
	IsPostSaved(ctx context.Context, userID, postID uuid.UUID) (bool, error)
	GetCollections(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type saveRepository struct {
	db *sql.DB
}

func NewSaveRepository(db *sql.DB) SaveRepository {
	return &saveRepository{db: db}
}

func (r *saveRepository) Create(ctx context.Context, save *model.Save) error {
	query := `
		INSERT INTO saves (id, user_id, post_id, collection, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, post_id) DO NOTHING
		RETURNING created_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		save.ID, save.UserID, save.PostID, save.Collection, save.CreatedAt,
	).Scan(&save.CreatedAt)
}

func (r *saveRepository) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	query := `DELETE FROM saves WHERE user_id = $1 AND post_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, userID, postID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("save not found")
	}

	return nil
}

func (r *saveRepository) GetByUserID(ctx context.Context, userID uuid.UUID, collection *string, limit, offset int) ([]model.Save, error) {
	var query string
	var args []interface{}

	if collection != nil {
		query = `
			SELECT id, user_id, post_id, collection, created_at
			FROM saves
			WHERE user_id = $1 AND collection = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{userID, *collection, limit, offset}
	} else {
		query = `
			SELECT id, user_id, post_id, collection, created_at
			FROM saves
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{userID, limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saves []model.Save
	for rows.Next() {
		save := model.Save{}
		err := rows.Scan(&save.ID, &save.UserID, &save.PostID, &save.Collection, &save.CreatedAt)
		if err != nil {
			return nil, err
		}
		saves = append(saves, save)
	}

	return saves, rows.Err()
}

func (r *saveRepository) IsPostSaved(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM saves WHERE user_id = $1 AND post_id = $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&exists)
	return exists, err
}

func (r *saveRepository) GetCollections(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT DISTINCT collection FROM saves
		WHERE user_id = $1 AND collection IS NOT NULL
		ORDER BY collection ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []string
	for rows.Next() {
		var collection string
		if err := rows.Scan(&collection); err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, rows.Err()
}
