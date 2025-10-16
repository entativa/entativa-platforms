package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/entativa/vignette/settings-service/internal/model"
	"github.com/google/uuid"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Create(ctx context.Context, settings *model.UserSettings) error {
	query := `
		INSERT INTO user_settings (
			id, user_id, appearance, privacy, notifications, chat, media,
			data_storage, security, accessibility, language, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		settings.ID, settings.UserID,
		settings.Appearance, settings.Privacy, settings.Notifications,
		settings.Chat, settings.Media, settings.DataStorage,
		settings.Security, settings.Accessibility, settings.Language,
		settings.CreatedAt, settings.UpdatedAt,
	)

	return err
}

func (r *SettingsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserSettings, error) {
	query := `
		SELECT id, user_id, appearance, privacy, notifications, chat, media,
		       data_storage, security, accessibility, language, created_at, updated_at
		FROM user_settings
		WHERE user_id = $1
	`

	settings := &model.UserSettings{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&settings.ID, &settings.UserID,
		&settings.Appearance, &settings.Privacy, &settings.Notifications,
		&settings.Chat, &settings.Media, &settings.DataStorage,
		&settings.Security, &settings.Accessibility, &settings.Language,
		&settings.CreatedAt, &settings.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("settings not found")
	}
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *SettingsRepository) Update(ctx context.Context, settings *model.UserSettings) error {
	query := `
		UPDATE user_settings
		SET appearance = $1, privacy = $2, notifications = $3, chat = $4,
		    media = $5, data_storage = $6, security = $7,
		    accessibility = $8, language = $9, updated_at = $10
		WHERE user_id = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		settings.Appearance, settings.Privacy, settings.Notifications,
		settings.Chat, settings.Media, settings.DataStorage,
		settings.Security, settings.Accessibility, settings.Language,
		settings.UpdatedAt, settings.UserID,
	)

	return err
}
