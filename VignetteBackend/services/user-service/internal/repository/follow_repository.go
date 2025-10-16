package repository

import (
	"context"
	"database/sql"

	"github.com/entativa/vignette/user-service/internal/model"
	"github.com/google/uuid"
)

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Create - Create a follow relationship
func (r *FollowRepository) Create(ctx context.Context, follow *model.Follow) error {
	query := `
		INSERT INTO follows (id, follower_id, following_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		follow.ID, follow.FollowerID, follow.FollowingID, follow.Status, follow.CreatedAt,
	)

	return err
}

// Delete - Remove a follow relationship
func (r *FollowRepository) Delete(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`

	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	return err
}

// IsFollowing - Check if user A follows user B
func (r *FollowRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND following_id = $2 AND status = 'active'
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

// GetFollowers - Get list of followers
func (r *FollowRepository) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	query := `
		SELECT follower_id FROM follows
		WHERE following_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []uuid.UUID
	for rows.Next() {
		var followerID uuid.UUID
		if err := rows.Scan(&followerID); err != nil {
			return nil, err
		}
		followers = append(followers, followerID)
	}

	return followers, nil
}

// GetFollowing - Get list of users being followed
func (r *FollowRepository) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	query := `
		SELECT following_id FROM follows
		WHERE follower_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []uuid.UUID
	for rows.Next() {
		var followingID uuid.UUID
		if err := rows.Scan(&followingID); err != nil {
			return nil, err
		}
		following = append(following, followingID)
	}

	return following, nil
}

// GetFollowersCount - Count followers
func (r *FollowRepository) GetFollowersCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM follows
		WHERE following_id = $1 AND status = 'active'
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// GetFollowingCount - Count following
func (r *FollowRepository) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM follows
		WHERE follower_id = $1 AND status = 'active'
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// GetMutualFollows - Get mutual follows between two users
func (r *FollowRepository) GetMutualFollows(ctx context.Context, userID1, userID2 uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT f1.following_id
		FROM follows f1
		INNER JOIN follows f2 ON f1.following_id = f2.following_id
		WHERE f1.follower_id = $1 AND f2.follower_id = $2
		  AND f1.status = 'active' AND f2.status = 'active'
	`

	rows, err := r.db.QueryContext(ctx, query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mutuals []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		mutuals = append(mutuals, userID)
	}

	return mutuals, nil
}
