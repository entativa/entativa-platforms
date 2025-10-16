package repository

import (
	"context"
	"database/sql"
	"fmt"

	"vignette/post-service/internal/model"

	"github.com/google/uuid"
)

type LikeRepository interface {
	CreatePostLike(ctx context.Context, like *model.Like) error
	CreateCommentLike(ctx context.Context, like *model.Like) error
	DeletePostLike(ctx context.Context, userID, postID uuid.UUID) error
	DeleteCommentLike(ctx context.Context, userID, commentID uuid.UUID) error
	GetPostLike(ctx context.Context, userID, postID uuid.UUID) (*model.Like, error)
	GetCommentLike(ctx context.Context, userID, commentID uuid.UUID) (*model.Like, error)
	GetPostLikers(ctx context.Context, postID uuid.UUID, limit, offset int) ([]uuid.UUID, error)
	GetCommentLikers(ctx context.Context, commentID uuid.UUID, limit, offset int) ([]uuid.UUID, error)
	GetUserPostReaction(ctx context.Context, userID, postID uuid.UUID) (*model.ReactionType, error)
}

type likeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) CreatePostLike(ctx context.Context, like *model.Like) error {
	query := `
		INSERT INTO likes (id, user_id, post_id, comment_id, reaction_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, post_id) WHERE post_id IS NOT NULL AND comment_id IS NULL
		DO UPDATE SET reaction_type = $5
		RETURNING created_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		like.ID, like.UserID, like.PostID, nil, like.ReactionType, like.CreatedAt,
	).Scan(&like.CreatedAt)
}

func (r *likeRepository) CreateCommentLike(ctx context.Context, like *model.Like) error {
	query := `
		INSERT INTO likes (id, user_id, post_id, comment_id, reaction_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, comment_id) WHERE comment_id IS NOT NULL AND post_id IS NULL
		DO UPDATE SET reaction_type = $5
		RETURNING created_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		like.ID, like.UserID, nil, like.CommentID, like.ReactionType, like.CreatedAt,
	).Scan(&like.CreatedAt)
}

func (r *likeRepository) DeletePostLike(ctx context.Context, userID, postID uuid.UUID) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2 AND comment_id IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, userID, postID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("like not found")
	}

	return nil
}

func (r *likeRepository) DeleteCommentLike(ctx context.Context, userID, commentID uuid.UUID) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND comment_id = $2 AND post_id IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, userID, commentID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("like not found")
	}

	return nil
}

func (r *likeRepository) GetPostLike(ctx context.Context, userID, postID uuid.UUID) (*model.Like, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, reaction_type, created_at
		FROM likes
		WHERE user_id = $1 AND post_id = $2 AND comment_id IS NULL
	`

	like := &model.Like{}
	err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(
		&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.ReactionType, &like.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return like, nil
}

func (r *likeRepository) GetCommentLike(ctx context.Context, userID, commentID uuid.UUID) (*model.Like, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, reaction_type, created_at
		FROM likes
		WHERE user_id = $1 AND comment_id = $2 AND post_id IS NULL
	`

	like := &model.Like{}
	err := r.db.QueryRowContext(ctx, query, userID, commentID).Scan(
		&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.ReactionType, &like.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return like, nil
}

func (r *likeRepository) GetPostLikers(ctx context.Context, postID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	query := `
		SELECT user_id FROM likes
		WHERE post_id = $1 AND comment_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}

func (r *likeRepository) GetCommentLikers(ctx context.Context, commentID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	query := `
		SELECT user_id FROM likes
		WHERE comment_id = $1 AND post_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, commentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}

func (r *likeRepository) GetUserPostReaction(ctx context.Context, userID, postID uuid.UUID) (*model.ReactionType, error) {
	query := `
		SELECT reaction_type FROM likes
		WHERE user_id = $1 AND post_id = $2 AND comment_id IS NULL
	`

	var reaction model.ReactionType
	err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&reaction)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reaction, nil
}
