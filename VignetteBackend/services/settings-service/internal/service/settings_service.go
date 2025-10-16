package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"vignette/settings-service/internal/crypto"
	"vignette/settings-service/internal/model"
	"vignette/settings-service/internal/repository"
	"github.com/google/uuid"
)

type SettingsService struct {
	settingsRepo *repository.SettingsRepository
	keyBackupRepo *repository.KeyBackupRepository
	encryptionSvc *crypto.EncryptionService
}

func NewSettingsService(
	settingsRepo *repository.SettingsRepository,
	keyBackupRepo *repository.KeyBackupRepository,
	encryptionSvc *crypto.EncryptionService,
) *SettingsService {
	return &SettingsService{
		settingsRepo:  settingsRepo,
		keyBackupRepo: keyBackupRepo,
		encryptionSvc: encryptionSvc,
	}
}

// ============================================
// SETTINGS MANAGEMENT
// ============================================

func (s *SettingsService) GetOrCreateSettings(ctx context.Context, userID uuid.UUID) (*model.UserSettings, error) {
	// Try to get existing settings
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err == nil {
		return settings, nil
	}

	// Create default settings
	settings = &model.UserSettings{
		ID:        uuid.New(),
		UserID:    userID,
		Appearance: model.AppearanceSettings{
			Theme:        "auto",
			AccentColor:  "#007AFF",
			FontSize:     "medium",
			HighContrast: false,
			ReduceMotion: false,
			CompactMode:  false,
		},
		Privacy: model.PrivacySettings{
			ProfileVisibility: "public",
			LastSeen:          "everyone",
			ReadReceipts:      true,
			TypingIndicator:   true,
			OnlineStatus:      true,
			BlockedUsers:      []string{},
			AllowTagging:      "everyone",
			AllowMentions:     "everyone",
			SearchableByEmail: true,
			SearchableByPhone: true,
			ShowActivity:      true,
		},
		Notifications: model.NotificationSettings{
			PushEnabled:        true,
			EmailEnabled:       true,
			SMSEnabled:         false,
			Likes:              true,
			Comments:           true,
			Mentions:           true,
			Follows:            true,
			Messages:           true,
			GroupMessages:      true,
			EventInvites:       true,
			EventReminders:     true,
			LiveStreams:        true,
			QuietHoursEnabled:  false,
			QuietHoursStart:    "22:00",
			QuietHoursEnd:      "08:00",
			NotificationSound:  "default",
			Vibrate:            true,
		},
		Chat: model.ChatSettings{
			KeyStorageLocation: model.StorageEntativaServer,
			EncryptionMethod:   model.EncryptionPassphrase,
			BackupKeysToServer: true,
			EnterToSend:        false,
			AutoDownloadMedia:  true,
			AutoPlayVideos:     true,
			AutoPlayGifs:       true,
			SaveToGallery:      false,
			AutoDeleteMessages: false,
			AutoDeleteAfterDays: 0,
			ScreenSecurity:     false,
			IncognitoKeyboard:  false,
		},
		Media: model.MediaSettings{
			AutoDownloadPhotos:   true,
			AutoDownloadVideos:   false,
			AutoDownloadFiles:    false,
			UploadQuality:        "high",
			VideoQuality:         "auto",
			MediaStorageLocation: "internal",
			AutoDeleteMedia:      false,
			AutoDeleteAfterDays:  0,
		},
		DataStorage: model.DataStorageSettings{
			DataSaverMode:      false,
			LowDataMode:        false,
			WiFiOnly:           false,
			CacheSize:          500,
			AutoClearCache:     false,
			AutoClearAfterDays: 30,
		},
		Security: model.SecuritySettings{
			TwoFactorEnabled: false,
			BiometricEnabled: false,
			AppLockEnabled:   false,
			AppLockTimeout:   300,
			ActiveSessions:   1,
			ShowLoginAlerts:  true,
			RecoveryEmail:    "",
			RecoveryPhone:    "",
		},
		Accessibility: model.AccessibilitySettings{
			ScreenReader:       false,
			ClosedCaptions:     false,
			ColorBlindMode:     "none",
			HighContrastText:   false,
			LargeText:          false,
			ReduceTransparency: false,
			VoiceControl:       false,
		},
		Language: model.LanguageSettings{
			AppLanguage:        "en",
			ContentLanguages:   []string{"en"},
			TranslationEnabled: true,
			AutoTranslate:      false,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.settingsRepo.Create(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *SettingsService) UpdateSettings(ctx context.Context, userID uuid.UUID, req *model.UpdateSettingsRequest) error {
	settings, err := s.GetOrCreateSettings(ctx, userID)
	if err != nil {
		return err
	}

	// Update fields if provided
	if req.Appearance != nil {
		settings.Appearance = *req.Appearance
	}
	if req.Privacy != nil {
		settings.Privacy = *req.Privacy
	}
	if req.Notifications != nil {
		settings.Notifications = *req.Notifications
	}
	if req.Chat != nil {
		settings.Chat = *req.Chat
	}
	if req.Media != nil {
		settings.Media = *req.Media
	}
	if req.DataStorage != nil {
		settings.DataStorage = *req.DataStorage
	}
	if req.Security != nil {
		settings.Security = *req.Security
	}
	if req.Accessibility != nil {
		settings.Accessibility = *req.Accessibility
	}
	if req.Language != nil {
		settings.Language = *req.Language
	}

	settings.UpdatedAt = time.Now()

	return s.settingsRepo.Update(ctx, settings)
}

// ============================================
// ENCRYPTED KEY BACKUP
// ============================================

func (s *SettingsService) CreateKeyBackup(ctx context.Context, userID uuid.UUID, req *model.CreateKeyBackupRequest) error {
	// Validate PIN or Passphrase
	var pinOrPassphrase string
	if req.EncryptionMethod == model.EncryptionPIN {
		if req.PIN == nil {
			return fmt.Errorf("PIN required for PIN encryption method")
		}
		if err := s.encryptionSvc.ValidatePIN(*req.PIN); err != nil {
			return err
		}
		pinOrPassphrase = *req.PIN
	} else {
		if req.Passphrase == nil {
			return fmt.Errorf("passphrase required for passphrase encryption method")
		}
		if err := s.encryptionSvc.ValidatePassphrase(*req.Passphrase); err != nil {
			return err
		}
		pinOrPassphrase = *req.Passphrase
	}

	// Decode the already-encrypted keys from client
	encryptedKeysFromClient, err := base64.StdEncoding.DecodeString(req.EncryptedKeys)
	if err != nil {
		return fmt.Errorf("invalid encrypted keys format")
	}

	// Generate salt
	salt, err := s.encryptionSvc.GenerateSalt()
	if err != nil {
		return err
	}

	// Hash PIN/Passphrase
	pinHash, err := s.encryptionSvc.HashPINOrPassphrase(pinOrPassphrase)
	if err != nil {
		return err
	}

	// Double-encrypt the keys (second layer of encryption)
	doubleEncryptedKeys, err := s.encryptionSvc.EncryptKeys(encryptedKeysFromClient, pinOrPassphrase, salt)
	if err != nil {
		return err
	}

	// Compute hash for integrity
	keysHash := s.encryptionSvc.HashData(doubleEncryptedKeys)

	// Create backup record
	backup := &model.EncryptedKeyBackup{
		ID:               uuid.New(),
		UserID:           userID,
		StorageLocation:  req.StorageLocation,
		EncryptionMethod: req.EncryptionMethod,
		EncryptedKeys:    doubleEncryptedKeys,
		KeysHash:         keysHash,
		PINHash:          pinHash,
		Salt:             crypto.EncodeBase64(salt),
		Iterations:       crypto.PBKDF2Iterations,
		DeviceID:         req.DeviceID,
		DeviceName:       req.DeviceName,
		BackupVersion:    1,
		LastBackupAt:     time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.keyBackupRepo.Create(ctx, backup)
	if err != nil {
		// Log failed attempt
		s.keyBackupRepo.LogAccess(ctx, userID, backup.ID, "create", req.DeviceID, "", false, err.Error())
		return err
	}

	// Log successful creation
	s.keyBackupRepo.LogAccess(ctx, userID, backup.ID, "create", req.DeviceID, "", true, "")

	return nil
}

func (s *SettingsService) RestoreKeyBackup(ctx context.Context, userID uuid.UUID, req *model.RestoreKeyBackupRequest) (*model.RestoreKeyBackupResponse, error) {
	// Get backup
	backup, err := s.keyBackupRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if backup == nil {
		return nil, fmt.Errorf("no backup found")
	}

	// Get PIN or Passphrase
	var pinOrPassphrase string
	if backup.EncryptionMethod == model.EncryptionPIN {
		if req.PIN == nil {
			return nil, fmt.Errorf("PIN required")
		}
		pinOrPassphrase = *req.PIN
	} else {
		if req.Passphrase == nil {
			return nil, fmt.Errorf("passphrase required")
		}
		pinOrPassphrase = *req.Passphrase
	}

	// Verify PIN/Passphrase
	if !s.encryptionSvc.VerifyPINOrPassphrase(pinOrPassphrase, backup.PINHash) {
		// Log failed attempt
		s.keyBackupRepo.LogAccess(ctx, userID, backup.ID, "failed_restore", backup.DeviceID, "", false, "incorrect PIN/passphrase")
		return nil, fmt.Errorf("incorrect PIN/passphrase")
	}

	// Decode salt
	salt, err := crypto.DecodeBase64(backup.Salt)
	if err != nil {
		return nil, err
	}

	// Decrypt keys (remove second layer of encryption)
	decryptedKeys, err := s.encryptionSvc.DecryptKeys(backup.EncryptedKeys, pinOrPassphrase, salt)
	if err != nil {
		// Log failed attempt
		s.keyBackupRepo.LogAccess(ctx, userID, backup.ID, "failed_restore", backup.DeviceID, "", false, "decryption failed")
		return nil, fmt.Errorf("decryption failed")
	}

	// Verify integrity
	computedHash := s.encryptionSvc.HashData(backup.EncryptedKeys)
	if computedHash != backup.KeysHash {
		return nil, fmt.Errorf("integrity check failed - backup may be corrupted")
	}

	// Log successful restore
	s.keyBackupRepo.LogAccess(ctx, userID, backup.ID, "restore", backup.DeviceID, "", true, "")

	// Return the decrypted keys (still encrypted by Signal, but not by PIN)
	return &model.RestoreKeyBackupResponse{
		EncryptedKeys: crypto.EncodeBase64(decryptedKeys),
		BackupVersion: backup.BackupVersion,
		BackupDate:    backup.LastBackupAt,
	}, nil
}

func (s *SettingsService) GetKeyBackupInfo(ctx context.Context, userID uuid.UUID) (*model.KeyBackupResponse, error) {
	backup, err := s.keyBackupRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if backup == nil {
		return &model.KeyBackupResponse{
			HasBackup: false,
		}, nil
	}

	return &model.KeyBackupResponse{
		HasBackup:        true,
		StorageLocation:  backup.StorageLocation,
		EncryptionMethod: backup.EncryptionMethod,
		LastBackupAt:     &backup.LastBackupAt,
		BackupVersion:    backup.BackupVersion,
	}, nil
}

func (s *SettingsService) DeleteKeyBackup(ctx context.Context, userID uuid.UUID) error {
	backup, err := s.keyBackupRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if backup == nil {
		return fmt.Errorf("no backup found")
	}

	err = s.keyBackupRepo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	// Log deletion
	s.keyBackupRepo.LogAccess(ctx, userID, backup.ID, "delete", backup.DeviceID, "", true, "")

	return nil
}

// GetStorageLocationInfo provides info about storage options
func (s *SettingsService) GetStorageLocationInfo() []*model.StorageLocationInfo {
	return []*model.StorageLocationInfo{
		{
			Location:    model.StorageEntativaServer,
			Name:        "Entativa Servers",
			Description: "Your keys are encrypted with your PIN/Passphrase and stored on Entativa servers. Only you can decrypt them.",
			Security:    "🔐 Double-encrypted (Signal + Your PIN/Passphrase). Server cannot decrypt. Authorities only get metadata.",
			Recommended: true,
		},
		{
			Location:    model.StorageLocalDevice,
			Name:        "Local Device",
			Description: "Keys stored only on this device. Lost if device is lost/reset.",
			Security:    "⚠️ Unreliable - you'll lose access if device is lost, broken, or reset.",
			Warning:     "Not recommended: Risk of data loss",
			Recommended: false,
		},
		{
			Location:    model.StorageICloud,
			Name:        "Apple iCloud",
			Description: "Keys backed up to your iCloud account.",
			Security:    "⚠️ Apple has your decryption keys. Can be accessed by Apple and shared with authorities if warranted.",
			Warning:     "Apple can access your keys",
			Recommended: false,
		},
		{
			Location:    model.StorageGoogleDrive,
			Name:        "Google Drive",
			Description: "Keys backed up to your Google Drive.",
			Security:    "⚠️ Google has your decryption keys. Can be accessed by Google and shared with authorities if warranted.",
			Warning:     "Google can access your keys",
			Recommended: false,
		},
	}
}
