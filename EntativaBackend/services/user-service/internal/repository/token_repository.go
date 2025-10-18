package repository

import (
	"context"
	"database/sql"
	"time"
)

// TokenRepository handles database operations for password reset tokens
type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository creates a new token repository instance
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// CreateResetToken creates a new password reset token
func (r *TokenRepository) CreateResetToken(ctx context.Context, token *PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens 
		(id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.Used,
		token.CreatedAt,
	)
	
	return err
}

// FindResetToken finds a reset token by token string
func (r *TokenRepository) FindResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1
		LIMIT 1
	`
	
	var resetToken PasswordResetToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.Used,
		&resetToken.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, &NotFoundError{"Reset token not found"}
	}
	if err != nil {
		return nil, err
	}
	
	return &resetToken, nil
}

// MarkTokenAsUsed marks a reset token as used
func (r *TokenRepository) MarkTokenAsUsed(ctx context.Context, tokenID string) error {
	query := `
		UPDATE password_reset_tokens
		SET used = true
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

// DeleteExpiredTokens deletes all expired reset tokens (cleanup job)
func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE expires_at < NOW()
	`
	
	_, err := r.db.ExecContext(ctx, query)
	return err
}

// GetUserResetTokens gets all reset tokens for a user
func (r *TokenRepository) GetUserResetTokens(ctx context.Context, userID string) ([]*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tokens []*PasswordResetToken
	for rows.Next() {
		var token PasswordResetToken
		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.Used,
			&token.CreatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}
	
	return tokens, rows.Err()
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
