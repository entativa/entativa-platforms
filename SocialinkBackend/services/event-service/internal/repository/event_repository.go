package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/entativa/socialink/event-service/internal/model"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *model.Event) error {
	query := `
		INSERT INTO events (
			id, creator_id, title, description, cover_photo,
			type, category, privacy,
			location_name, address, city, country, latitude, longitude, online_link,
			start_time, end_time, timezone,
			is_recurring, recurrence_rule, recurrence_end_date,
			allow_guest_invites, require_approval, max_attendees,
			co_hosts, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
	`

	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.CreatorID, event.Title, event.Description, event.CoverPhoto,
		event.Type, event.Category, event.Privacy,
		event.LocationName, event.Address, event.City, event.Country, event.Latitude, event.Longitude, event.OnlineLink,
		event.StartTime, event.EndTime, event.Timezone,
		event.IsRecurring, event.RecurrenceRule, event.RecurrenceEndDate,
		event.AllowGuestInvites, event.RequireApproval, event.MaxAttendees,
		model.UUIDArray(event.CoHosts), event.CreatedAt, event.UpdatedAt,
	)

	return err
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	query := `
		SELECT id, creator_id, title, description, cover_photo,
		       type, category, privacy,
		       location_name, address, city, country, latitude, longitude, online_link,
		       start_time, end_time, timezone,
		       is_recurring, recurrence_rule, recurrence_end_date,
		       allow_guest_invites, require_approval, max_attendees,
		       co_hosts, going_count, interested_count, view_count,
		       is_cancelled, cancelled_at, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	event := &model.Event{}
	var coHosts model.UUIDArray
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.CreatorID, &event.Title, &event.Description, &event.CoverPhoto,
		&event.Type, &event.Category, &event.Privacy,
		&event.LocationName, &event.Address, &event.City, &event.Country, &event.Latitude, &event.Longitude, &event.OnlineLink,
		&event.StartTime, &event.EndTime, &event.Timezone,
		&event.IsRecurring, &event.RecurrenceRule, &event.RecurrenceEndDate,
		&event.AllowGuestInvites, &event.RequireApproval, &event.MaxAttendees,
		&coHosts, &event.GoingCount, &event.InterestedCount, &event.ViewCount,
		&event.IsCancelled, &event.CancelledAt, &event.CreatedAt, &event.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found")
	}
	if err != nil {
		return nil, err
	}

	event.CoHosts = coHosts
	return event, nil
}

func (r *EventRepository) Update(ctx context.Context, event *model.Event) error {
	query := `
		UPDATE events
		SET title = $1, description = $2, cover_photo = $3,
		    start_time = $4, end_time = $5, online_link = $6,
		    co_hosts = $7, updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(ctx, query,
		event.Title, event.Description, event.CoverPhoto,
		event.StartTime, event.EndTime, event.OnlineLink,
		model.UUIDArray(event.CoHosts), event.UpdatedAt, event.ID,
	)

	return err
}

func (r *EventRepository) Cancel(ctx context.Context, eventID uuid.UUID) error {
	query := `UPDATE events SET is_cancelled = TRUE, cancelled_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, eventID)
	return err
}

func (r *EventRepository) IncrementViewCount(ctx context.Context, eventID uuid.UUID) error {
	query := `UPDATE events SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, eventID)
	return err
}

func (r *EventRepository) GetUpcoming(ctx context.Context, limit, offset int, category string) ([]*model.Event, error) {
	query := `
		SELECT id, creator_id, title, description, cover_photo,
		       type, category, start_time, end_time, timezone,
		       location_name, city, country, going_count, interested_count,
		       created_at
		FROM events
		WHERE start_time > NOW() AND NOT is_cancelled AND privacy = 'public'
	`
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY start_time ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*model.Event{}
	for rows.Next() {
		event := &model.Event{}
		err := rows.Scan(
			&event.ID, &event.CreatorID, &event.Title, &event.Description, &event.CoverPhoto,
			&event.Type, &event.Category, &event.StartTime, &event.EndTime, &event.Timezone,
			&event.LocationName, &event.City, &event.Country, &event.GoingCount, &event.InterestedCount,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) Search(ctx context.Context, query string, limit, offset int) ([]*model.Event, error) {
	sqlQuery := `
		SELECT id, creator_id, title, description, cover_photo,
		       type, category, start_time, going_count, interested_count
		FROM events
		WHERE to_tsvector('english', title || ' ' || description) @@ plainto_tsquery('english', $1)
		  AND NOT is_cancelled
		ORDER BY start_time ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*model.Event{}
	for rows.Next() {
		event := &model.Event{}
		err := rows.Scan(
			&event.ID, &event.CreatorID, &event.Title, &event.Description, &event.CoverPhoto,
			&event.Type, &event.Category, &event.StartTime, &event.GoingCount, &event.InterestedCount,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetNearby(ctx context.Context, lat, lng float64, radiusKm int, limit int) ([]*model.Event, error) {
	query := `
		SELECT id, creator_id, title, start_time, location_name, city, country,
		       going_count, interested_count,
		       ST_Distance(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) / 1000 as distance_km
		FROM events
		WHERE location IS NOT NULL
		  AND NOT is_cancelled
		  AND start_time > NOW()
		  AND ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3 * 1000)
		ORDER BY distance_km ASC
		LIMIT $4
	`

	rows, err := r.db.QueryContext(ctx, query, lng, lat, radiusKm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*model.Event{}
	for rows.Next() {
		event := &model.Event{}
		var distanceKm float64
		err := rows.Scan(
			&event.ID, &event.CreatorID, &event.Title, &event.StartTime,
			&event.LocationName, &event.City, &event.Country,
			&event.GoingCount, &event.InterestedCount, &distanceKm,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetByCreator(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]*model.Event, error) {
	query := `
		SELECT id, title, type, category, start_time, end_time,
		       going_count, interested_count, is_cancelled, created_at
		FROM events
		WHERE creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*model.Event{}
	for rows.Next() {
		event := &model.Event{}
		err := rows.Scan(
			&event.ID, &event.Title, &event.Type, &event.Category,
			&event.StartTime, &event.EndTime, &event.GoingCount, &event.InterestedCount,
			&event.IsCancelled, &event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) IsCoHost(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	query := `SELECT co_hosts FROM events WHERE id = $1`
	
	var coHosts model.UUIDArray
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(&coHosts)
	if err != nil {
		return false, err
	}

	for _, coHost := range coHosts {
		if coHost == userID {
			return true, nil
		}
	}

	return false, nil
}
