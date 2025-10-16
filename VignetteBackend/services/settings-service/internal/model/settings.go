package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Storage location options for encrypted keys
const (
	StorageEntativaServer StorageLocation = "entativa_server"
	StorageLocalDevice    StorageLocation = "local_device"
	StorageICloud         StorageLocation = "icloud"
	StorageGoogleDrive    StorageLocation = "google_drive"
)

// Encryption methods
const (
	EncryptionPIN        EncryptionMethod = "pin"        // 6-digit PIN
	EncryptionPassphrase EncryptionMethod = "passphrase" // Strong passphrase
)

type StorageLocation string
type EncryptionMethod string

// UserSettings represents all user settings
type UserSettings struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	
	// App Settings
	Appearance        AppearanceSettings        `json:"appearance" db:"appearance"`
	Privacy           PrivacySettings           `json:"privacy" db:"privacy"`
	Notifications     NotificationSettings      `json:"notifications" db:"notifications"`
	Chat              ChatSettings              `json:"chat" db:"chat"`
	Media             MediaSettings             `json:"media" db:"media"`
	DataStorage       DataStorageSettings       `json:"data_storage" db:"data_storage"`
	Security          SecuritySettings          `json:"security" db:"security"`
	Accessibility     AccessibilitySettings     `json:"accessibility" db:"accessibility"`
	Language          LanguageSettings          `json:"language" db:"language"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AppearanceSettings for UI customization
type AppearanceSettings struct {
	Theme              string `json:"theme"`                // light, dark, auto
	AccentColor        string `json:"accent_color"`         // #hex color
	FontSize           string `json:"font_size"`            // small, medium, large
	HighContrast       bool   `json:"high_contrast"`
	ReduceMotion       bool   `json:"reduce_motion"`
	CompactMode        bool   `json:"compact_mode"`
}

// PrivacySettings for privacy controls
type PrivacySettings struct {
	ProfileVisibility   string   `json:"profile_visibility"`    // public, friends, private
	LastSeen            string   `json:"last_seen"`             // everyone, friends, nobody
	ReadReceipts        bool     `json:"read_receipts"`
	TypingIndicator     bool     `json:"typing_indicator"`
	OnlineStatus        bool     `json:"online_status"`
	BlockedUsers        []string `json:"blocked_users"`
	AllowTagging        string   `json:"allow_tagging"`         // everyone, friends, nobody
	AllowMentions       string   `json:"allow_mentions"`        // everyone, friends, nobody
	SearchableByEmail   bool     `json:"searchable_by_email"`
	SearchableByPhone   bool     `json:"searchable_by_phone"`
	ShowActivity        bool     `json:"show_activity"`
}

// NotificationSettings for notification preferences
type NotificationSettings struct {
	PushEnabled          bool     `json:"push_enabled"`
	EmailEnabled         bool     `json:"email_enabled"`
	SMSEnabled           bool     `json:"sms_enabled"`
	
	// Notification types
	Likes                bool     `json:"likes"`
	Comments             bool     `json:"comments"`
	Mentions             bool     `json:"mentions"`
	Follows              bool     `json:"follows"`
	Messages             bool     `json:"messages"`
	GroupMessages        bool     `json:"group_messages"`
	EventInvites         bool     `json:"event_invites"`
	EventReminders       bool     `json:"event_reminders"`
	LiveStreams          bool     `json:"live_streams"`
	
	// Notification schedule
	QuietHoursEnabled    bool     `json:"quiet_hours_enabled"`
	QuietHoursStart      string   `json:"quiet_hours_start"`    // 22:00
	QuietHoursEnd        string   `json:"quiet_hours_end"`      // 08:00
	
	// Sound & vibration
	NotificationSound    string   `json:"notification_sound"`
	Vibrate              bool     `json:"vibrate"`
}

// ChatSettings for messaging preferences
type ChatSettings struct {
	// Key storage
	KeyStorageLocation   StorageLocation  `json:"key_storage_location"`
	EncryptionMethod     EncryptionMethod `json:"encryption_method"`
	BackupKeysToServer   bool             `json:"backup_keys_to_server"`
	
	// Chat behavior
	EnterToSend          bool   `json:"enter_to_send"`
	AutoDownloadMedia    bool   `json:"auto_download_media"`
	AutoPlayVideos       bool   `json:"auto_play_videos"`
	AutoPlayGifs         bool   `json:"auto_play_gifs"`
	SaveToGallery        bool   `json:"save_to_gallery"`
	
	// Message retention
	AutoDeleteMessages   bool   `json:"auto_delete_messages"`
	AutoDeleteAfterDays  int    `json:"auto_delete_after_days"`
	
	// Security
	ScreenSecurity       bool   `json:"screen_security"`      // Prevent screenshots
	IncognitoKeyboard    bool   `json:"incognito_keyboard"`   // No keyboard learning
}

// MediaSettings for media handling
type MediaSettings struct {
	AutoDownloadPhotos   bool   `json:"auto_download_photos"`
	AutoDownloadVideos   bool   `json:"auto_download_videos"`
	AutoDownloadFiles    bool   `json:"auto_download_files"`
	
	// Quality
	UploadQuality        string `json:"upload_quality"`       // original, high, medium, low
	VideoQuality         string `json:"video_quality"`        // auto, 4k, 1080p, 720p, 480p
	
	// Storage
	MediaStorageLocation string `json:"media_storage_location"` // internal, external, cloud
	AutoDeleteMedia      bool   `json:"auto_delete_media"`
	AutoDeleteAfterDays  int    `json:"auto_delete_after_days"`
}

// DataStorageSettings for data and storage
type DataStorageSettings struct {
	DataSaverMode        bool   `json:"data_saver_mode"`
	LowDataMode          bool   `json:"low_data_mode"`
	WiFiOnly             bool   `json:"wifi_only"`
	
	// Cache
	CacheSize            int    `json:"cache_size"`           // MB
	AutoClearCache       bool   `json:"auto_clear_cache"`
	AutoClearAfterDays   int    `json:"auto_clear_after_days"`
}

// SecuritySettings for account security
type SecuritySettings struct {
	TwoFactorEnabled     bool   `json:"two_factor_enabled"`
	BiometricEnabled     bool   `json:"biometric_enabled"`
	
	// App lock
	AppLockEnabled       bool   `json:"app_lock_enabled"`
	AppLockTimeout       int    `json:"app_lock_timeout"`     // seconds
	
	// Sessions
	ActiveSessions       int    `json:"active_sessions"`
	ShowLoginAlerts      bool   `json:"show_login_alerts"`
	
	// Recovery
	RecoveryEmail        string `json:"recovery_email"`
	RecoveryPhone        string `json:"recovery_phone"`
}

// AccessibilitySettings for accessibility features
type AccessibilitySettings struct {
	ScreenReader         bool   `json:"screen_reader"`
	ClosedCaptions       bool   `json:"closed_captions"`
	ColorBlindMode       string `json:"color_blind_mode"`     // none, protanopia, deuteranopia, tritanopia
	HighContrastText     bool   `json:"high_contrast_text"`
	LargeText            bool   `json:"large_text"`
	ReduceTransparency   bool   `json:"reduce_transparency"`
	VoiceControl         bool   `json:"voice_control"`
}

// LanguageSettings for language preferences
type LanguageSettings struct {
	AppLanguage          string   `json:"app_language"`         // en, es, fr, etc
	ContentLanguages     []string `json:"content_languages"`    // Preferred content languages
	TranslationEnabled   bool     `json:"translation_enabled"`
	AutoTranslate        bool     `json:"auto_translate"`
}

// EncryptedKeyBackup for storing encrypted chat keys
type EncryptedKeyBackup struct {
	ID                uuid.UUID        `json:"id" db:"id"`
	UserID            uuid.UUID        `json:"user_id" db:"user_id"`
	
	// Storage
	StorageLocation   StorageLocation  `json:"storage_location" db:"storage_location"`
	EncryptionMethod  EncryptionMethod `json:"encryption_method" db:"encryption_method"`
	
	// Encrypted data (double-encrypted: first by Signal, then by PIN/Passphrase)
	EncryptedKeys     []byte           `json:"-" db:"encrypted_keys"`           // Never sent in API
	KeysHash          string           `json:"keys_hash" db:"keys_hash"`        // SHA256 for verification
	
	// PIN/Passphrase protection (hashed)
	PINHash           string           `json:"-" db:"pin_hash"`                 // bcrypt of PIN
	Salt              string           `json:"-" db:"salt"`                     // Random salt
	Iterations        int              `json:"iterations" db:"iterations"`       // PBKDF2 iterations
	
	// Metadata (only this is visible to authorities)
	DeviceID          string           `json:"device_id" db:"device_id"`
	DeviceName        string           `json:"device_name" db:"device_name"`
	BackupVersion     int              `json:"backup_version" db:"backup_version"`
	LastBackupAt      time.Time        `json:"last_backup_at" db:"last_backup_at"`
	
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
}

// StorageLocationInfo provides info about storage options
type StorageLocationInfo struct {
	Location    StorageLocation `json:"location"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Security    string          `json:"security"`
	Warning     string          `json:"warning,omitempty"`
	Recommended bool            `json:"recommended"`
}

// API Request/Response Models

type UpdateSettingsRequest struct {
	Appearance    *AppearanceSettings    `json:"appearance,omitempty"`
	Privacy       *PrivacySettings       `json:"privacy,omitempty"`
	Notifications *NotificationSettings  `json:"notifications,omitempty"`
	Chat          *ChatSettings          `json:"chat,omitempty"`
	Media         *MediaSettings         `json:"media,omitempty"`
	DataStorage   *DataStorageSettings   `json:"data_storage,omitempty"`
	Security      *SecuritySettings      `json:"security,omitempty"`
	Accessibility *AccessibilitySettings `json:"accessibility,omitempty"`
	Language      *LanguageSettings      `json:"language,omitempty"`
}

type CreateKeyBackupRequest struct {
	StorageLocation  StorageLocation  `json:"storage_location" binding:"required"`
	EncryptionMethod EncryptionMethod `json:"encryption_method" binding:"required"`
	PIN              *string          `json:"pin,omitempty"`              // 6-digit PIN (if method is PIN)
	Passphrase       *string          `json:"passphrase,omitempty"`       // Strong passphrase (if method is passphrase)
	EncryptedKeys    string           `json:"encrypted_keys" binding:"required"` // Base64 encoded
	DeviceID         string           `json:"device_id" binding:"required"`
	DeviceName       string           `json:"device_name" binding:"required"`
}

type RestoreKeyBackupRequest struct {
	PIN              *string `json:"pin,omitempty"`
	Passphrase       *string `json:"passphrase,omitempty"`
}

type KeyBackupResponse struct {
	HasBackup        bool             `json:"has_backup"`
	StorageLocation  StorageLocation  `json:"storage_location,omitempty"`
	EncryptionMethod EncryptionMethod `json:"encryption_method,omitempty"`
	LastBackupAt     *time.Time       `json:"last_backup_at,omitempty"`
	BackupVersion    int              `json:"backup_version"`
}

type RestoreKeyBackupResponse struct {
	EncryptedKeys string    `json:"encrypted_keys"` // Base64 encoded
	BackupVersion int       `json:"backup_version"`
	BackupDate    time.Time `json:"backup_date"`
}

// Helper methods for JSONB scanning

func (s *AppearanceSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s AppearanceSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *PrivacySettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s PrivacySettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *NotificationSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s NotificationSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ChatSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s ChatSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *MediaSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s MediaSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *DataStorageSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s DataStorageSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SecuritySettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s SecuritySettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *AccessibilitySettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s AccessibilitySettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *LanguageSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s LanguageSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}
