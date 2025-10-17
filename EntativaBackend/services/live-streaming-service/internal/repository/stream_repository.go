package repository

import (
	"context"
	"database/sql"
	"fmt"

	"socialink/live-streaming-service/internal/model"
	"github.com/google/uuid"
)

type StreamRepository struct {
	db *sql.DB
}

func NewStreamRepository(db *sql.DB) *StreamRepository {
	return &StreamRepository{db: db}
}

func (r *StreamRepository) Create(ctx context.Context, stream *model.LiveStream) error {
	query := `
		INSERT INTO live_streams (
			id, streamer_id, title, description, thumbnail_url,
			status, quality, is_private, category, tags,
			stream_key, rtmp_url, hls_url, webrtc_url,
			viewer_count, peak_viewers, total_views, record_stream,
			scheduled_for, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	_, err := r.db.ExecContext(ctx, query,
		stream.ID, stream.StreamerID, stream.Title, stream.Description, stream.ThumbnailURL,
		stream.Status, stream.Quality, stream.IsPrivate, stream.Category, model.StringArray(stream.Tags),
		stream.StreamKey, stream.RTMPUrl, stream.HLSUrl, stream.WebRTCUrl,
		stream.ViewerCount, stream.PeakViewers, stream.TotalViews, stream.RecordStream,
		stream.ScheduledFor, stream.CreatedAt, stream.UpdatedAt,
	)

	return err
}

func (r *StreamRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.LiveStream, error) {
	query := `
		SELECT id, streamer_id, title, description, thumbnail_url,
		       status, quality, is_private, category, tags,
		       stream_key, rtmp_url, hls_url, webrtc_url,
		       viewer_count, peak_viewers, total_views, likes_count, comments_count, shares_count,
		       record_stream, recording_url, scheduled_for, started_at, ended_at, duration,
		       created_at, updated_at
		FROM live_streams
		WHERE id = $1
	`

	stream := &model.LiveStream{}
	var tags model.StringArray
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stream.ID, &stream.StreamerID, &stream.Title, &stream.Description, &stream.ThumbnailURL,
		&stream.Status, &stream.Quality, &stream.IsPrivate, &stream.Category, &tags,
		&stream.StreamKey, &stream.RTMPUrl, &stream.HLSUrl, &stream.WebRTCUrl,
		&stream.ViewerCount, &stream.PeakViewers, &stream.TotalViews, &stream.LikesCount,
		&stream.CommentsCount, &stream.SharesCount, &stream.RecordStream, &stream.RecordingURL,
		&stream.ScheduledFor, &stream.StartedAt, &stream.EndedAt, &stream.Duration,
		&stream.CreatedAt, &stream.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("stream not found")
	}
	if err != nil {
		return nil, err
	}

	stream.Tags = tags
	return stream, nil
}

func (r *StreamRepository) Update(ctx context.Context, stream *model.LiveStream) error {
	query := `
		UPDATE live_streams
		SET title = $1, description = $2, status = $3, viewer_count = $4,
		    peak_viewers = $5, total_views = $6, recording_url = $7,
		    started_at = $8, ended_at = $9, duration = $10, updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.ExecContext(ctx, query,
		stream.Title, stream.Description, stream.Status, stream.ViewerCount,
		stream.PeakViewers, stream.TotalViews, stream.RecordingURL,
		stream.StartedAt, stream.EndedAt, stream.Duration, stream.UpdatedAt,
		stream.ID,
	)

	return err
}

func (r *StreamRepository) GetLive(ctx context.Context, limit, offset int, category string) ([]*model.LiveStream, error) {
	query := `
		SELECT id, streamer_id, title, description, thumbnail_url,
		       status, quality, is_private, category, tags,
		       rtmp_url, hls_url, webrtc_url,
		       viewer_count, peak_viewers, total_views, likes_count, comments_count,
		       started_at, created_at
		FROM live_streams
		WHERE status = 'live'
	`
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY viewer_count DESC, started_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	streams := []*model.LiveStream{}
	for rows.Next() {
		stream := &model.LiveStream{}
		var tags model.StringArray
		err := rows.Scan(
			&stream.ID, &stream.StreamerID, &stream.Title, &stream.Description, &stream.ThumbnailURL,
			&stream.Status, &stream.Quality, &stream.IsPrivate, &stream.Category, &tags,
			&stream.RTMPUrl, &stream.HLSUrl, &stream.WebRTCUrl,
			&stream.ViewerCount, &stream.PeakViewers, &stream.TotalViews, &stream.LikesCount, &stream.CommentsCount,
			&stream.StartedAt, &stream.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		stream.Tags = tags
		streams = append(streams, stream)
	}

	return streams, nil
}

func (r *StreamRepository) GetByStreamer(ctx context.Context, streamerID uuid.UUID, limit, offset int) ([]*model.LiveStream, error) {
	query := `
		SELECT id, streamer_id, title, status, viewer_count, total_views,
		       scheduled_for, started_at, ended_at, created_at
		FROM live_streams
		WHERE streamer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, streamerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	streams := []*model.LiveStream{}
	for rows.Next() {
		stream := &model.LiveStream{}
		err := rows.Scan(
			&stream.ID, &stream.StreamerID, &stream.Title, &stream.Status, &stream.ViewerCount,
			&stream.TotalViews, &stream.ScheduledFor, &stream.StartedAt, &stream.EndedAt, &stream.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}

	return streams, nil
}
