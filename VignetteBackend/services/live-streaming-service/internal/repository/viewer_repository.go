package repository

import (
	"context"
	"database/sql"

	"github.com/entativa/vignette/live-streaming-service/internal/model"
	"github.com/google/uuid"
)

type ViewerRepository struct {
	db *sql.DB
}

func NewViewerRepository(db *sql.DB) *ViewerRepository {
	return &ViewerRepository{db: db}
}

func (r *ViewerRepository) Create(ctx context.Context, viewer *model.StreamViewer) error {
	query := `
		INSERT INTO stream_viewers (id, stream_id, viewer_id, joined_at, is_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (stream_id, viewer_id) DO UPDATE
		SET is_active = TRUE, joined_at = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		viewer.ID, viewer.StreamID, viewer.ViewerID, viewer.JoinedAt, viewer.IsActive,
	)

	return err
}

func (r *ViewerRepository) GetViewer(ctx context.Context, streamID, viewerID uuid.UUID) (*model.StreamViewer, error) {
	query := `
		SELECT id, stream_id, viewer_id, joined_at, left_at, watch_time, is_active
		FROM stream_viewers
		WHERE stream_id = $1 AND viewer_id = $2
	`

	viewer := &model.StreamViewer{}
	err := r.db.QueryRowContext(ctx, query, streamID, viewerID).Scan(
		&viewer.ID, &viewer.StreamID, &viewer.ViewerID, &viewer.JoinedAt,
		&viewer.LeftAt, &viewer.WatchTime, &viewer.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return viewer, nil
}

func (r *ViewerRepository) Update(ctx context.Context, viewer *model.StreamViewer) error {
	query := `
		UPDATE stream_viewers
		SET left_at = $1, watch_time = $2, is_active = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		viewer.LeftAt, viewer.WatchTime, viewer.IsActive, viewer.ID,
	)

	return err
}

func (r *ViewerRepository) MarkAllInactive(ctx context.Context, streamID uuid.UUID) error {
	query := `UPDATE stream_viewers SET is_active = FALSE WHERE stream_id = $1 AND is_active = TRUE`
	_, err := r.db.ExecContext(ctx, query, streamID)
	return err
}

func (r *ViewerRepository) GetUniqueViewerCount(ctx context.Context, streamID uuid.UUID) (int, error) {
	query := `SELECT COUNT(DISTINCT viewer_id) FROM stream_viewers WHERE stream_id = $1`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, streamID).Scan(&count)
	return count, err
}

func (r *ViewerRepository) GetAverageWatchTime(ctx context.Context, streamID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(AVG(watch_time), 0) FROM stream_viewers WHERE stream_id = $1`
	
	var avg float64
	err := r.db.QueryRowContext(ctx, query, streamID).Scan(&avg)
	return int(avg), err
}
