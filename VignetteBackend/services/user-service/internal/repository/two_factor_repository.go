package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"vignette/user-service/internal/model"

	"github.com/google/uuid"
)

type TwoFactorRepository struct {
	db *sql.DB
}

func NewTwoFactorRepository(db *sql.DB) *TwoFactorRepository {
	return &TwoFactorRepository{db: db}
}

// Create creates a new 2FA configuration
func (r *TwoFactorRepository) Create(twoFactor *model.TwoFactorAuth) error {
	backupCodesJSON, _ := json.Marshal(twoFactor.BackupCodes)
	
	query := `
		INSERT INTO two_factor_auth (
			id, user_id, secret, is_enabled, backup_codes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			secret = EXCLUDED.secret,
			backup_codes = EXCLUDED.backup_codes,
			created_at = EXCLUDED.created_at
	`

	_, err := r.db.Exec(
		query,
		twoFactor.ID,
		twoFactor.UserID,
		twoFactor.Secret,
		twoFactor.IsEnabled,
		backupCodesJSON,
		twoFactor.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create 2FA config: %w", err)
	}

	return nil
}

// FindByUserID finds 2FA configuration by user ID
func (r *TwoFactorRepository) FindByUserID(userID uuid.UUID) (*model.TwoFactorAuth, error) {
	query := `
		SELECT id, user_id, secret, is_enabled, backup_codes, 
		       created_at, enabled_at, last_used_at
		FROM two_factor_auth
		WHERE user_id = $1
	`

	twoFactor := &model.TwoFactorAuth{}
	var backupCodesJSON []byte

	err := r.db.QueryRow(query, userID).Scan(
		&twoFactor.ID,
		&twoFactor.UserID,
		&twoFactor.Secret,
		&twoFactor.IsEnabled,
		&backupCodesJSON,
		&twoFactor.CreatedAt,
		&twoFactor.EnabledAt,
		&twoFactor.LastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("2FA not configured")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find 2FA config: %w", err)
	}

	// Unmarshal backup codes
	if len(backupCodesJSON) > 0 {
		json.Unmarshal(backupCodesJSON, &twoFactor.BackupCodes)
	}

	return twoFactor, nil
}

// Update updates 2FA configuration
func (r *TwoFactorRepository) Update(twoFactor *model.TwoFactorAuth) error {
	backupCodesJSON, _ := json.Marshal(twoFactor.BackupCodes)
	
	query := `
		UPDATE two_factor_auth
		SET is_enabled = $1, backup_codes = $2, enabled_at = $3, last_used_at = $4
		WHERE user_id = $5
	`

	_, err := r.db.Exec(
		query,
		twoFactor.IsEnabled,
		backupCodesJSON,
		twoFactor.EnabledAt,
		twoFactor.LastUsedAt,
		twoFactor.UserID,
	)

	return err
}

// DeleteByUserID deletes 2FA configuration
func (r *TwoFactorRepository) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM two_factor_auth WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
