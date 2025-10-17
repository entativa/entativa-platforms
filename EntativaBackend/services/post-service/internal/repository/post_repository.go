package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"socialink/post-service/internal/model"

	"github.com/google/uuid"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]model.Post, *string, error)
	IncrementLikes(ctx context.Context, postID uuid.UUID) error
	DecrementLikes(ctx context.Context, postID uuid.UUID) error
	IncrementComments(ctx context.Context, postID uuid.UUID) error
	DecrementComments(ctx context.Context, postID uuid.UUID) error
	IncrementShares(ctx context.Context, postID uuid.UUID) error
	IncrementViews(ctx context.Context, postID uuid.UUID) error
	IncrementSaves(ctx context.Context, postID uuid.UUID) error
	DecrementSaves(ctx context.Context, postID uuid.UUID) error
	GetByHashtag(ctx context.Context, hashtag string, limit, offset int) ([]model.Post, error)
	GetReels(ctx context.Context, limit, offset int) ([]model.Post, error)
	GetTrendingPosts(ctx context.Context, limit int, timeWindow time.Duration) ([]model.Post, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (
			id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			filter_used, is_carousel, likes_count, comments_count, views_count,
			saves_count, shares_count, is_edited, is_sponsored, is_reels,
			comments_enabled, likes_visible, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING created_at, updated_at
	`

	mediaIDsJSON, _ := json.Marshal(post.MediaIDs)
	taggedUserIDsJSON, _ := json.Marshal(post.TaggedUserIDs)
	hashtagsJSON, _ := json.Marshal(post.Hashtags)
	locationJSON, _ := json.Marshal(post.Location)

	return r.db.QueryRowContext(
		ctx, query,
		post.ID, post.UserID, post.Caption, mediaIDsJSON, locationJSON, taggedUserIDsJSON,
		hashtagsJSON, post.FilterUsed, post.IsCarousel, post.LikesCount, post.CommentsCount,
		post.ViewsCount, post.SavesCount, post.SharesCount, post.IsEdited, post.IsSponsored,
		post.IsReels, post.CommentsEnabled, post.LikesVisible, post.CreatedAt, post.UpdatedAt,
	).Scan(&post.CreatedAt, &post.UpdatedAt)
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	query := `
		SELECT id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			   filter_used, is_carousel, likes_count, comments_count, views_count,
			   saves_count, shares_count, is_edited, edited_at, is_sponsored, is_reels,
			   comments_enabled, likes_visible, created_at, updated_at, deleted_at
		FROM posts
		WHERE id = $1 AND deleted_at IS NULL
	`

	post := &model.Post{}
	var mediaIDsJSON, taggedUserIDsJSON, hashtagsJSON, locationJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.Caption, &mediaIDsJSON, &locationJSON,
		&taggedUserIDsJSON, &hashtagsJSON, &post.FilterUsed, &post.IsCarousel,
		&post.LikesCount, &post.CommentsCount, &post.ViewsCount, &post.SavesCount,
		&post.SharesCount, &post.IsEdited, &post.EditedAt, &post.IsSponsored,
		&post.IsReels, &post.CommentsEnabled, &post.LikesVisible,
		&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}

	// Unmarshal JSON fields
	if len(mediaIDsJSON) > 0 {
		json.Unmarshal(mediaIDsJSON, &post.MediaIDs)
	}
	if len(taggedUserIDsJSON) > 0 {
		json.Unmarshal(taggedUserIDsJSON, &post.TaggedUserIDs)
	}
	if len(hashtagsJSON) > 0 {
		json.Unmarshal(hashtagsJSON, &post.Hashtags)
	}
	if len(locationJSON) > 0 && string(locationJSON) != "null" {
		var location model.Location
		if json.Unmarshal(locationJSON, &location) == nil {
			post.Location = &location
		}
	}

	return post, nil
}

func (r *postRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Post, error) {
	query := `
		SELECT id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			   filter_used, is_carousel, likes_count, comments_count, views_count,
			   saves_count, shares_count, is_edited, edited_at, is_sponsored, is_reels,
			   comments_enabled, likes_visible, created_at, updated_at, deleted_at
		FROM posts
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	query := `
		UPDATE posts
		SET caption = $1, location = $2, hashtags = $3, is_edited = $4,
			edited_at = $5, updated_at = $6, comments_enabled = $7, likes_visible = $8
		WHERE id = $9 AND deleted_at IS NULL
	`

	hashtagsJSON, _ := json.Marshal(post.Hashtags)
	locationJSON, _ := json.Marshal(post.Location)

	result, err := r.db.ExecContext(
		ctx, query,
		post.Caption, locationJSON, hashtagsJSON, post.IsEdited, post.EditedAt,
		post.UpdatedAt, post.CommentsEnabled, post.LikesVisible, post.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE posts SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

func (r *postRepository) GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]model.Post, *string, error) {
	query := `
		SELECT id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			   filter_used, is_carousel, likes_count, comments_count, views_count,
			   saves_count, shares_count, is_edited, edited_at, is_sponsored, is_reels,
			   comments_enabled, likes_visible, created_at, updated_at, deleted_at
		FROM posts
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	
	if cursor != "" {
		query += ` AND created_at < (SELECT created_at FROM posts WHERE id = $1)`
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

	posts, err := r.scanPosts(rows)
	if err != nil {
		return nil, nil, err
	}

	var nextCursor *string
	if len(posts) > limit {
		posts = posts[:limit]
		cursorStr := posts[limit-1].ID.String()
		nextCursor = &cursorStr
	}

	return posts, nextCursor, nil
}

func (r *postRepository) IncrementLikes(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET likes_count = likes_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) DecrementLikes(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET likes_count = GREATEST(0, likes_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) IncrementComments(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET comments_count = comments_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) DecrementComments(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET comments_count = GREATEST(0, comments_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) IncrementShares(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET shares_count = shares_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) IncrementViews(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) IncrementSaves(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET saves_count = saves_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) DecrementSaves(ctx context.Context, postID uuid.UUID) error {
	query := `UPDATE posts SET saves_count = GREATEST(0, saves_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *postRepository) GetByHashtag(ctx context.Context, hashtag string, limit, offset int) ([]model.Post, error) {
	query := `
		SELECT id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			   filter_used, is_carousel, likes_count, comments_count, views_count,
			   saves_count, shares_count, is_edited, edited_at, is_sponsored, is_reels,
			   comments_enabled, likes_visible, created_at, updated_at, deleted_at
		FROM posts
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

	return r.scanPosts(rows)
}

func (r *postRepository) GetReels(ctx context.Context, limit, offset int) ([]model.Post, error) {
	query := `
		SELECT id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			   filter_used, is_carousel, likes_count, comments_count, views_count,
			   saves_count, shares_count, is_edited, edited_at, is_sponsored, is_reels,
			   comments_enabled, likes_visible, created_at, updated_at, deleted_at
		FROM posts
		WHERE deleted_at IS NULL AND is_reels = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *postRepository) GetTrendingPosts(ctx context.Context, limit int, timeWindow time.Duration) ([]model.Post, error) {
	// Facebook explore algorithm: weighted engagement
	query := `
		SELECT id, user_id, caption, media_ids, location, tagged_user_ids, hashtags,
			   filter_used, is_carousel, likes_count, comments_count, views_count,
			   saves_count, shares_count, is_edited, edited_at, is_sponsored, is_reels,
			   comments_enabled, likes_visible, created_at, updated_at, deleted_at
		FROM posts
		WHERE deleted_at IS NULL
		AND created_at > $1
		ORDER BY (
			likes_count + 
			views_count / 10 + 
			saves_count * 2 + 
			shares_count * 3 + 
			comments_count * 2
		) DESC
		LIMIT $2
	`

	since := time.Now().Add(-timeWindow)
	rows, err := r.db.QueryContext(ctx, query, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *postRepository) scanPosts(rows *sql.Rows) ([]model.Post, error) {
	var posts []model.Post

	for rows.Next() {
		post := model.Post{}
		var mediaIDsJSON, taggedUserIDsJSON, hashtagsJSON, locationJSON []byte

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Caption, &mediaIDsJSON, &locationJSON,
			&taggedUserIDsJSON, &hashtagsJSON, &post.FilterUsed, &post.IsCarousel,
			&post.LikesCount, &post.CommentsCount, &post.ViewsCount, &post.SavesCount,
			&post.SharesCount, &post.IsEdited, &post.EditedAt, &post.IsSponsored,
			&post.IsReels, &post.CommentsEnabled, &post.LikesVisible,
			&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if len(mediaIDsJSON) > 0 {
			json.Unmarshal(mediaIDsJSON, &post.MediaIDs)
		}
		if len(taggedUserIDsJSON) > 0 {
			json.Unmarshal(taggedUserIDsJSON, &post.TaggedUserIDs)
		}
		if len(hashtagsJSON) > 0 {
			json.Unmarshal(hashtagsJSON, &post.Hashtags)
		}
		if len(locationJSON) > 0 && string(locationJSON) != "null" {
			var location model.Location
			if json.Unmarshal(locationJSON, &location) == nil {
				post.Location = &location
			}
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}
