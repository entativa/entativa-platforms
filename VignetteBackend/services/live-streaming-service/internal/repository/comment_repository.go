package repository

import (
	"context"
	"database/sql"

	"github.com/entativa/vignette/live-streaming-service/internal/model"
	"github.com/google/uuid"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *model.StreamComment) error {
	query := `
		INSERT INTO stream_comments (id, stream_id, user_id, content, is_pinned, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		comment.ID, comment.StreamID, comment.UserID, comment.Content, comment.IsPinned, comment.CreatedAt,
	)

	return err
}

func (r *CommentRepository) GetByStream(ctx context.Context, streamID uuid.UUID, limit, offset int) ([]*model.StreamComment, error) {
	query := `
		SELECT id, stream_id, user_id, content, is_pinned, created_at
		FROM stream_comments
		WHERE stream_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, streamID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.StreamComment{}
	for rows.Next() {
		comment := &model.StreamComment{}
		err := rows.Scan(
			&comment.ID, &comment.StreamID, &comment.UserID,
			&comment.Content, &comment.IsPinned, &comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) PinComment(ctx context.Context, commentID uuid.UUID) error {
	// First unpin all comments in this stream
	query1 := `
		UPDATE stream_comments SET is_pinned = FALSE
		WHERE stream_id = (SELECT stream_id FROM stream_comments WHERE id = $1)
	`
	_, err := r.db.ExecContext(ctx, query1, commentID)
	if err != nil {
		return err
	}

	// Pin this comment
	query2 := `UPDATE stream_comments SET is_pinned = TRUE WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query2, commentID)
	return err
}

func (r *CommentRepository) Delete(ctx context.Context, commentID, streamerID uuid.UUID) error {
	// Only streamer can delete comments
	query := `
		DELETE FROM stream_comments
		WHERE id = $1 AND stream_id IN (
			SELECT id FROM live_streams WHERE streamer_id = $2
		)
	`

	_, err := r.db.ExecContext(ctx, query, commentID, streamerID)
	return err
}
