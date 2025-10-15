package repository

import (
	"database/sql"
	"fmt"
	"time"

	"vignette/user-service/internal/model"

	"github.com/google/uuid"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(session *model.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, access_token, refresh_token, device_info, 
			ip_address, user_agent, expires_at, created_at, last_active_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		query,
		session.ID,
		session.UserID,
		session.AccessToken,
		session.RefreshToken,
		session.DeviceInfo,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
		session.CreatedAt,
		session.LastActiveAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// FindByAccessToken finds a session by access token
func (r *SessionRepository) FindByAccessToken(token string) (*model.Session, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, device_info, 
		       ip_address, user_agent, expires_at, created_at, last_active_at
		FROM sessions
		WHERE access_token = $1
	`

	session := &model.Session{}
	err := r.db.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.AccessToken,
		&session.RefreshToken,
		&session.DeviceInfo,
		&session.IPAddress,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.LastActiveAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	return session, nil
}

// DeleteByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

// DeleteExpiredSessions deletes all expired sessions
func (r *SessionRepository) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	return err
}

// UpdateLastActive updates the last active timestamp
func (r *SessionRepository) UpdateLastActive(sessionID uuid.UUID) error {
	query := `UPDATE sessions SET last_active_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), sessionID)
	return err
}
