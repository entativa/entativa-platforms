package repository

import (
	"context"
	"database/sql"
	"fmt"

	"socialink/event-service/internal/model"
	"github.com/google/uuid"
)

type RSVPRepository struct {
	db *sql.DB
}

func NewRSVPRepository(db *sql.DB) *RSVPRepository {
	return &RSVPRepository{db: db}
}

func (r *RSVPRepository) Upsert(ctx context.Context, rsvp *model.EventRSVP) error {
	query := `
		INSERT INTO event_rsvps (id, event_id, user_id, status, guest_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (event_id, user_id) DO UPDATE
		SET status = $4, guest_count = $5, updated_at = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		rsvp.ID, rsvp.EventID, rsvp.UserID, rsvp.Status, rsvp.GuestCount,
		rsvp.CreatedAt, rsvp.UpdatedAt,
	)

	return err
}

func (r *RSVPRepository) Get(ctx context.Context, eventID, userID uuid.UUID) (*model.EventRSVP, error) {
	query := `
		SELECT id, event_id, user_id, status, guest_count, checked_in, checked_in_at, created_at, updated_at
		FROM event_rsvps
		WHERE event_id = $1 AND user_id = $2
	`

	rsvp := &model.EventRSVP{}
	err := r.db.QueryRowContext(ctx, query, eventID, userID).Scan(
		&rsvp.ID, &rsvp.EventID, &rsvp.UserID, &rsvp.Status, &rsvp.GuestCount,
		&rsvp.CheckedIn, &rsvp.CheckedInAt, &rsvp.CreatedAt, &rsvp.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return rsvp, nil
}

func (r *RSVPRepository) GetByEvent(ctx context.Context, eventID uuid.UUID, status model.RSVPStatus, limit, offset int) ([]*model.EventRSVP, error) {
	query := `
		SELECT id, event_id, user_id, status, guest_count, checked_in, created_at
		FROM event_rsvps
		WHERE event_id = $1
	`
	args := []interface{}{eventID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT $3 OFFSET $4"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rsvps := []*model.EventRSVP{}
	for rows.Next() {
		rsvp := &model.EventRSVP{}
		err := rows.Scan(
			&rsvp.ID, &rsvp.EventID, &rsvp.UserID, &rsvp.Status,
			&rsvp.GuestCount, &rsvp.CheckedIn, &rsvp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rsvps = append(rsvps, rsvp)
	}

	return rsvps, nil
}

func (r *RSVPRepository) CheckIn(ctx context.Context, eventID, userID uuid.UUID) error {
	query := `
		UPDATE event_rsvps
		SET checked_in = TRUE, checked_in_at = NOW()
		WHERE event_id = $1 AND user_id = $2 AND status = 'going'
	`

	_, err := r.db.ExecContext(ctx, query, eventID, userID)
	return err
}

func (r *RSVPRepository) Delete(ctx context.Context, eventID, userID uuid.UUID) error {
	query := `DELETE FROM event_rsvps WHERE event_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, eventID, userID)
	return err
}

func (r *RSVPRepository) GetUserEvents(ctx context.Context, userID uuid.UUID, status model.RSVPStatus, upcoming bool, limit, offset int) ([]*uuid.UUID, error) {
	query := `
		SELECT event_id
		FROM event_rsvps
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	if status != "" {
		query += " AND status = $" + fmt.Sprint(argIdx)
		args = append(args, status)
		argIdx++
	}

	if upcoming {
		query += " AND event_id IN (SELECT id FROM events WHERE start_time > NOW() AND NOT is_cancelled)"
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(argIdx) + " OFFSET $" + fmt.Sprint(argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventIDs := []*uuid.UUID{}
	for rows.Next() {
		var eventID uuid.UUID
		if err := rows.Scan(&eventID); err != nil {
			return nil, err
		}
		eventIDs = append(eventIDs, &eventID)
	}

	return eventIDs, nil
}
