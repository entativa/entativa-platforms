package repository

import (
	"database/sql"
	"fmt"
	"time"

	"vignette/user-service/internal/model"

	"github.com/google/uuid"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(token *model.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (
			id, user_id, token, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	return nil
}

// FindValidTokens finds all valid (not used, not expired) tokens
func (r *PasswordResetRepository) FindValidTokens() ([]*model.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE used_at IS NULL AND expires_at > $1
	`

	rows, err := r.db.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*model.PasswordResetToken
	for rows.Next() {
		token := &model.PasswordResetToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.UsedAt,
			&token.CreatedAt,
		)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// MarkAsUsed marks a token as used
func (r *PasswordResetRepository) MarkAsUsed(tokenID uuid.UUID) error {
	query := `UPDATE password_reset_tokens SET used_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), tokenID)
	return err
}

// DeleteExpiredTokens deletes expired tokens
func (r *PasswordResetRepository) DeleteExpiredTokens() error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	return err
}
