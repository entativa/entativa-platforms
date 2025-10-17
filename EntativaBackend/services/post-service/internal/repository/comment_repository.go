package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"socialink/post-service/internal/model"

	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Comment, error)
	GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetRepliesCount(ctx context.Context, commentID uuid.UUID) (int64, error)
	IncrementLikes(ctx context.Context, commentID uuid.UUID) error
	DecrementLikes(ctx context.Context, commentID uuid.UUID) error
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (
			id, post_id, user_id, parent_id, content, media_id,
			likes_count, is_edited, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		comment.ID, comment.PostID, comment.UserID, comment.ParentID,
		comment.Content, comment.MediaID, comment.LikesCount, comment.IsEdited,
		comment.CreatedAt, comment.UpdatedAt,
	).Scan(&comment.CreatedAt, &comment.UpdatedAt)
}

func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, content, media_id,
			   likes_count, is_edited, edited_at, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL
	`

	comment := &model.Comment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID,
		&comment.Content, &comment.MediaID, &comment.LikesCount, &comment.IsEdited,
		&comment.EditedAt, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, err
	}

	return comment, nil
}

func (r *commentRepository) GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, content, media_id,
			   likes_count, is_edited, edited_at, created_at, updated_at, deleted_at
		FROM comments
		WHERE post_id = $1 AND parent_id IS NULL AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComments(rows)
}

func (r *commentRepository) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, content, media_id,
			   likes_count, is_edited, edited_at, created_at, updated_at, deleted_at
		FROM comments
		WHERE parent_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComments(rows)
}

func (r *commentRepository) Update(ctx context.Context, comment *model.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, is_edited = $2, edited_at = $3, updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx, query,
		comment.Content, comment.IsEdited, comment.EditedAt, comment.UpdatedAt, comment.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE comments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

func (r *commentRepository) GetRepliesCount(ctx context.Context, commentID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM comments WHERE parent_id = $1 AND deleted_at IS NULL`
	
	var count int64
	err := r.db.QueryRowContext(ctx, query, commentID).Scan(&count)
	return count, err
}

func (r *commentRepository) IncrementLikes(ctx context.Context, commentID uuid.UUID) error {
	query := `UPDATE comments SET likes_count = likes_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, commentID)
	return err
}

func (r *commentRepository) DecrementLikes(ctx context.Context, commentID uuid.UUID) error {
	query := `UPDATE comments SET likes_count = GREATEST(0, likes_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, commentID)
	return err
}

func (r *commentRepository) scanComments(rows *sql.Rows) ([]model.Comment, error) {
	var comments []model.Comment

	for rows.Next() {
		comment := model.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID,
			&comment.Content, &comment.MediaID, &comment.LikesCount, &comment.IsEdited,
			&comment.EditedAt, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, rows.Err()
}
