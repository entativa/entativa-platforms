package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Take represents a take/reel in the database
type Take struct {
	ID            string
	UserID        string
	Username      string
	UserAvatar    string
	VideoURL      string
	ThumbnailURL  string
	Caption       string
	AudioName     string
	AudioURL      string
	Duration      int
	LikesCount    int
	CommentsCount int
	SharesCount   int
	ViewsCount    int
	IsLiked       bool
	IsSaved       bool
	Hashtags      []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TakeComment represents a comment on a take
type TakeComment struct {
	ID         string    `json:"id"`
	TakeID     string    `json:"take_id"`
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	UserAvatar string    `json:"user_avatar,omitempty"`
	Text       string    `json:"text"`
	LikesCount int       `json:"likes_count"`
	IsLiked    bool      `json:"is_liked"`
	CreatedAt  time.Time `json:"created_at"`
}

// TakesRepository handles database operations for takes
type TakesRepository struct {
	db *sql.DB
}

// NewTakesRepository creates a new takes repository
func NewTakesRepository(db *sql.DB) *TakesRepository {
	return &TakesRepository{db: db}
}

// GetFeed returns paginated takes feed
func (r *TakesRepository) GetFeed(ctx context.Context, userID string, page, limit int) ([]*Take, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT 
			t.id, t.user_id, t.video_url, t.thumbnail_url, t.caption,
			t.audio_name, t.audio_url, t.duration, t.likes_count,
			t.comments_count, t.shares_count, t.views_count,
			t.hashtags, t.created_at,
			u.username, u.profile_picture_url,
			COALESCE((SELECT true FROM take_likes WHERE take_id = t.id AND user_id = $1), false) as is_liked,
			COALESCE((SELECT true FROM take_saves WHERE take_id = t.id AND user_id = $1), false) as is_saved
		FROM takes t
		JOIN users u ON t.user_id = u.id
		WHERE t.is_active = true
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch takes feed: %w", err)
	}
	defer rows.Close()
	
	var takes []*Take
	for rows.Next() {
		take := &Take{}
		var hashtagsJSON []byte
		
		err := rows.Scan(
			&take.ID, &take.UserID, &take.VideoURL, &take.ThumbnailURL,
			&take.Caption, &take.AudioName, &take.AudioURL, &take.Duration,
			&take.LikesCount, &take.CommentsCount, &take.SharesCount,
			&take.ViewsCount, &hashtagsJSON, &take.CreatedAt,
			&take.Username, &take.UserAvatar, &take.IsLiked, &take.IsSaved,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan take: %w", err)
		}
		
		// Parse hashtags
		if len(hashtagsJSON) > 0 {
			json.Unmarshal(hashtagsJSON, &take.Hashtags)
		}
		
		takes = append(takes, take)
	}
	
	return takes, rows.Err()
}

// GetByID returns a specific take
func (r *TakesRepository) GetByID(ctx context.Context, takeID, userID string) (*Take, error) {
	query := `
		SELECT 
			t.id, t.user_id, t.video_url, t.thumbnail_url, t.caption,
			t.audio_name, t.audio_url, t.duration, t.likes_count,
			t.comments_count, t.shares_count, t.views_count,
			t.hashtags, t.created_at,
			u.username, u.profile_picture_url,
			COALESCE((SELECT true FROM take_likes WHERE take_id = t.id AND user_id = $2), false) as is_liked,
			COALESCE((SELECT true FROM take_saves WHERE take_id = t.id AND user_id = $2), false) as is_saved
		FROM takes t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1 AND t.is_active = true
	`
	
	take := &Take{}
	var hashtagsJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, takeID, userID).Scan(
		&take.ID, &take.UserID, &take.VideoURL, &take.ThumbnailURL,
		&take.Caption, &take.AudioName, &take.AudioURL, &take.Duration,
		&take.LikesCount, &take.CommentsCount, &take.SharesCount,
		&take.ViewsCount, &hashtagsJSON, &take.CreatedAt,
		&take.Username, &take.UserAvatar, &take.IsLiked, &take.IsSaved,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("take not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch take: %w", err)
	}
	
	// Parse hashtags
	if len(hashtagsJSON) > 0 {
		json.Unmarshal(hashtagsJSON, &take.Hashtags)
	}
	
	return take, nil
}

// LikeTake likes a take
func (r *TakesRepository) LikeTake(ctx context.Context, takeID, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insert like
	_, err = tx.ExecContext(ctx, `
		INSERT INTO take_likes (take_id, user_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (take_id, user_id) DO NOTHING
	`, takeID, userID, time.Now())
	if err != nil {
		return err
	}
	
	// Increment likes count
	_, err = tx.ExecContext(ctx, `
		UPDATE takes
		SET likes_count = likes_count + 1
		WHERE id = $1
	`, takeID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// UnlikeTake unlikes a take
func (r *TakesRepository) UnlikeTake(ctx context.Context, takeID, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete like
	_, err = tx.ExecContext(ctx, `
		DELETE FROM take_likes
		WHERE take_id = $1 AND user_id = $2
	`, takeID, userID)
	if err != nil {
		return err
	}
	
	// Decrement likes count
	_, err = tx.ExecContext(ctx, `
		UPDATE takes
		SET likes_count = GREATEST(likes_count - 1, 0)
		WHERE id = $1
	`, takeID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// IncrementViews increments the view count
func (r *TakesRepository) IncrementViews(ctx context.Context, takeID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE takes
		SET views_count = views_count + 1
		WHERE id = $1
	`, takeID)
	
	return err
}

// GetComments returns comments for a take
func (r *TakesRepository) GetComments(ctx context.Context, takeID string, page, limit int) ([]*TakeComment, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT 
			c.id, c.take_id, c.user_id, c.text, c.likes_count, c.created_at,
			u.username, u.profile_picture_url
		FROM take_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.take_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, takeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer rows.Close()
	
	var comments []*TakeComment
	for rows.Next() {
		comment := &TakeComment{}
		
		err := rows.Scan(
			&comment.ID, &comment.TakeID, &comment.UserID, &comment.Text,
			&comment.LikesCount, &comment.CreatedAt, &comment.Username,
			&comment.UserAvatar,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		
		comments = append(comments, comment)
	}
	
	return comments, rows.Err()
}

// AddComment adds a comment to a take
func (r *TakesRepository) AddComment(ctx context.Context, comment *TakeComment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insert comment
	_, err = tx.ExecContext(ctx, `
		INSERT INTO take_comments (id, take_id, user_id, text, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, comment.ID, comment.TakeID, comment.UserID, comment.Text, comment.CreatedAt)
	if err != nil {
		return err
	}
	
	// Increment comments count
	_, err = tx.ExecContext(ctx, `
		UPDATE takes
		SET comments_count = comments_count + 1
		WHERE id = $1
	`, comment.TakeID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}
