package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Session represents a user session
type Session struct {
	ID            string
	UserID        string
	AccessToken   string
	RefreshToken  string
	DeviceInfo    string
	IPAddress     string
	UserAgent     string
	ExpiresAt     time.Time
	CreatedAt     time.Time
	LastActiveAt  time.Time
}

// SessionRepository handles database operations for sessions
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession creates a new session
func (r *SessionRepository) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, access_token, refresh_token,
			device_info, ip_address, user_agent,
			expires_at, created_at, last_active_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := r.db.ExecContext(
		ctx, query,
		session.ID, session.UserID, session.AccessToken, session.RefreshToken,
		session.DeviceInfo, session.IPAddress, session.UserAgent,
		session.ExpiresAt, session.CreatedAt, session.LastActiveAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	
	return nil
}

// FindSessionByToken finds a session by access token
func (r *SessionRepository) FindSessionByToken(ctx context.Context, token string) (*Session, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token,
		       device_info, ip_address, user_agent,
		       expires_at, created_at, last_active_at
		FROM sessions
		WHERE access_token = $1 AND expires_at > NOW()
	`
	
	session := &Session{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.AccessToken, &session.RefreshToken,
		&session.DeviceInfo, &session.IPAddress, &session.UserAgent,
		&session.ExpiresAt, &session.CreatedAt, &session.LastActiveAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	
	return session, nil
}

// UpdateLastActive updates the last active timestamp
func (r *SessionRepository) UpdateLastActive(ctx context.Context, sessionID string) error {
	query := `
		UPDATE sessions
		SET last_active_at = $1
		WHERE id = $2
	`
	
	_, err := r.db.ExecContext(ctx, query, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("failed to update last active: %w", err)
	}
	
	return nil
}

// DeleteSession deletes a session (logout)
func (r *SessionRepository) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE access_token = $1`
	
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	
	return nil
}

// InvalidateAllUserSessions invalidates all sessions for a user
func (r *SessionRepository) InvalidateAllUserSessions(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate sessions: %w", err)
	}
	
	return nil
}

// DeleteExpiredSessions deletes all expired sessions (cleanup job)
func (r *SessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("Deleted %d expired sessions\n", rowsAffected)
	}
	
	return nil
}

// GetUserSessions gets all active sessions for a user
func (r *SessionRepository) GetUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token,
		       device_info, ip_address, user_agent,
		       expires_at, created_at, last_active_at
		FROM sessions
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY last_active_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()
	
	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		err := rows.Scan(
			&session.ID, &session.UserID, &session.AccessToken, &session.RefreshToken,
			&session.DeviceInfo, &session.IPAddress, &session.UserAgent,
			&session.ExpiresAt, &session.CreatedAt, &session.LastActiveAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}
	
	return sessions, rows.Err()
}
