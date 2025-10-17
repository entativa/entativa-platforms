package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/entativa/socialink/settings-service/internal/model"
	"github.com/google/uuid"
)

type KeyBackupRepository struct {
	db *sql.DB
}

func NewKeyBackupRepository(db *sql.DB) *KeyBackupRepository {
	return &KeyBackupRepository{db: db}
}

func (r *KeyBackupRepository) Create(ctx context.Context, backup *model.EncryptedKeyBackup) error {
	query := `
		INSERT INTO encrypted_key_backups (
			id, user_id, storage_location, encryption_method,
			encrypted_keys, keys_hash, pin_hash, salt, iterations,
			device_id, device_name, backup_version, last_backup_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (user_id, device_id) DO UPDATE
		SET encrypted_keys = $5, keys_hash = $6, pin_hash = $7, salt = $8,
		    backup_version = encrypted_key_backups.backup_version + 1,
		    last_backup_at = $13, updated_at = $15
	`

	_, err := r.db.ExecContext(ctx, query,
		backup.ID, backup.UserID, backup.StorageLocation, backup.EncryptionMethod,
		backup.EncryptedKeys, backup.KeysHash, backup.PINHash, backup.Salt, backup.Iterations,
		backup.DeviceID, backup.DeviceName, backup.BackupVersion, backup.LastBackupAt,
		backup.CreatedAt, backup.UpdatedAt,
	)

	return err
}

func (r *KeyBackupRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.EncryptedKeyBackup, error) {
	query := `
		SELECT id, user_id, storage_location, encryption_method,
		       encrypted_keys, keys_hash, pin_hash, salt, iterations,
		       device_id, device_name, backup_version, last_backup_at,
		       created_at, updated_at
		FROM encrypted_key_backups
		WHERE user_id = $1
		ORDER BY last_backup_at DESC
		LIMIT 1
	`

	backup := &model.EncryptedKeyBackup{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&backup.ID, &backup.UserID, &backup.StorageLocation, &backup.EncryptionMethod,
		&backup.EncryptedKeys, &backup.KeysHash, &backup.PINHash, &backup.Salt, &backup.Iterations,
		&backup.DeviceID, &backup.DeviceName, &backup.BackupVersion, &backup.LastBackupAt,
		&backup.CreatedAt, &backup.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return backup, nil
}

func (r *KeyBackupRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM encrypted_key_backups WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *KeyBackupRepository) LogAccess(ctx context.Context, userID, backupID uuid.UUID, action string, deviceID, ipAddress string, success bool, failureReason string) error {
	query := `
		INSERT INTO key_backup_access_log (
			id, user_id, backup_id, action, device_id, ip_address,
			success, failure_reason, accessed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		uuid.New(), userID, backupID, action, deviceID, ipAddress,
		success, failureReason,
	)

	return err
}
