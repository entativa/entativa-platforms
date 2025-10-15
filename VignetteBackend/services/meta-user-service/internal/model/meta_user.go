package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// MetaUser represents a unified user identity across Meta platforms (Socialink & Vignette)
type MetaUser struct {
	ID                    uuid.UUID              `json:"id" db:"id"`
	MetaID                string                 `json:"meta_id" db:"meta_id"` // Unique Meta-level identifier
	Email                 string                 `json:"email" db:"email"`
	EmailVerified         bool                   `json:"email_verified" db:"email_verified"`
	PhoneNumber           *string                `json:"phone_number,omitempty" db:"phone_number"`
	PhoneVerified         bool                   `json:"phone_verified" db:"phone_verified"`
	PasswordHash          string                 `json:"-" db:"password_hash"`
	Status                UserStatus             `json:"status" db:"status"`
	TrustScore            float64                `json:"trust_score" db:"trust_score"` // ML-computed trust score
	RiskLevel             RiskLevel              `json:"risk_level" db:"risk_level"`
	AccountTier           AccountTier            `json:"account_tier" db:"account_tier"`
	PlatformLinks         PlatformLinks          `json:"platform_links" db:"platform_links"`
	SecurityProfile       SecurityProfile        `json:"security_profile" db:"security_profile"`
	PrivacySettings       PrivacySettings        `json:"privacy_settings" db:"privacy_settings"`
	DataRights            DataRights             `json:"data_rights" db:"data_rights"`
	DeviceFingerprints    []DeviceFingerprint    `json:"device_fingerprints" db:"device_fingerprints"`
	BiometricTokens       []BiometricToken       `json:"biometric_tokens,omitempty" db:"biometric_tokens"`
	FederatedIdentities   []FederatedIdentity    `json:"federated_identities" db:"federated_identities"`
	SessionManagement     SessionManagement      `json:"session_management" db:"session_management"`
	AnomalyDetection      AnomalyDetection       `json:"anomaly_detection" db:"anomaly_detection"`
	CrossPlatformActivity CrossPlatformActivity  `json:"cross_platform_activity" db:"cross_platform_activity"`
	ComplianceData        ComplianceData         `json:"compliance_data" db:"compliance_data"`
	Metadata              map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
	LastSeenAt            *time.Time             `json:"last_seen_at,omitempty" db:"last_seen_at"`
	DeletedAt             *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
	UserStatusPending   UserStatus = "pending"
	UserStatusDeleted   UserStatus = "deleted"
)

type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

type AccountTier string

const (
	AccountTierBasic      AccountTier = "basic"
	AccountTierPremium    AccountTier = "premium"
	AccountTierBusiness   AccountTier = "business"
	AccountTierEnterprise AccountTier = "enterprise"
)

// PlatformLinks maintains links to platform-specific user accounts
type PlatformLinks struct {
	SocialinkUserID   *uuid.UUID `json:"socialink_user_id,omitempty"`
	VignetteUserID    *uuid.UUID `json:"vignette_user_id,omitempty"`
	LinkedAt          time.Time  `json:"linked_at"`
	LinkStatus        LinkStatus `json:"link_status"`
	SyncEnabled       bool       `json:"sync_enabled"`
	LastSyncedAt      *time.Time `json:"last_synced_at,omitempty"`
	SyncConflicts     int        `json:"sync_conflicts"`
}

type LinkStatus string

const (
	LinkStatusLinked     LinkStatus = "linked"
	LinkStatusUnlinked   LinkStatus = "unlinked"
	LinkStatusPending    LinkStatus = "pending"
	LinkStatusConflicted LinkStatus = "conflicted"
)

// SecurityProfile contains advanced security settings
type SecurityProfile struct {
	TwoFactorEnabled      bool                   `json:"two_factor_enabled"`
	TwoFactorMethod       TwoFactorMethod        `json:"two_factor_method"`
	BackupCodes           []string               `json:"backup_codes,omitempty"`
	BiometricEnabled      bool                   `json:"biometric_enabled"`
	PasskeyEnabled        bool                   `json:"passkey_enabled"`
	PasswordlessEnabled   bool                   `json:"passwordless_enabled"`
	SecurityKeys          []SecurityKey          `json:"security_keys,omitempty"`
	LoginApprovalRequired bool                   `json:"login_approval_required"`
	TrustedDevices        []TrustedDevice        `json:"trusted_devices"`
	IPWhitelist           []string               `json:"ip_whitelist,omitempty"`
	GeoRestrictions       []string               `json:"geo_restrictions,omitempty"`
	LastPasswordChange    time.Time              `json:"last_password_change"`
	PasswordExpiresAt     *time.Time             `json:"password_expires_at,omitempty"`
	FailedLoginAttempts   int                    `json:"failed_login_attempts"`
	LastFailedLoginAt     *time.Time             `json:"last_failed_login_at,omitempty"`
	AccountLockedUntil    *time.Time             `json:"account_locked_until,omitempty"`
}

type TwoFactorMethod string

const (
	TwoFactorMethodNone         TwoFactorMethod = "none"
	TwoFactorMethodTOTP         TwoFactorMethod = "totp"
	TwoFactorMethodSMS          TwoFactorMethod = "sms"
	TwoFactorMethodEmail        TwoFactorMethod = "email"
	TwoFactorMethodBiometric    TwoFactorMethod = "biometric"
	TwoFactorMethodHardwareKey  TwoFactorMethod = "hardware_key"
)

// PrivacySettings for cross-platform privacy management
type PrivacySettings struct {
	ProfileVisibility     VisibilityLevel        `json:"profile_visibility"`
	ActivityVisibility    VisibilityLevel        `json:"activity_visibility"`
	ContactsVisibility    VisibilityLevel        `json:"contacts_visibility"`
	CrossPlatformSharing  bool                   `json:"cross_platform_sharing"`
	DataPortabilityEnabled bool                  `json:"data_portability_enabled"`
	MarketingOptIn        bool                   `json:"marketing_opt_in"`
	PersonalizationOptIn  bool                   `json:"personalization_opt_in"`
	ThirdPartyDataSharing bool                   `json:"third_party_data_sharing"`
	LocationTracking      LocationTrackingLevel  `json:"location_tracking"`
	SearchVisibility      bool                   `json:"search_visibility"`
	CustomPermissions     map[string]interface{} `json:"custom_permissions"`
}

type VisibilityLevel string

const (
	VisibilityPublic     VisibilityLevel = "public"
	VisibilityFriends    VisibilityLevel = "friends"
	VisibilityFollowers  VisibilityLevel = "followers"
	VisibilityPrivate    VisibilityLevel = "private"
	VisibilityCustom     VisibilityLevel = "custom"
)

type LocationTrackingLevel string

const (
	LocationTrackingAlways LocationTrackingLevel = "always"
	LocationTrackingInUse  LocationTrackingLevel = "in_use"
	LocationTrackingNever  LocationTrackingLevel = "never"
)

// DataRights for GDPR/CCPA compliance
type DataRights struct {
	RightToAccess         bool       `json:"right_to_access"`
	RightToRectification  bool       `json:"right_to_rectification"`
	RightToErasure        bool       `json:"right_to_erasure"`
	RightToPortability    bool       `json:"right_to_portability"`
	RightToObject         bool       `json:"right_to_object"`
	DataExportRequested   bool       `json:"data_export_requested"`
	DataExportRequestedAt *time.Time `json:"data_export_requested_at,omitempty"`
	DataExportCompletedAt *time.Time `json:"data_export_completed_at,omitempty"`
	DeletionRequested     bool       `json:"deletion_requested"`
	DeletionRequestedAt   *time.Time `json:"deletion_requested_at,omitempty"`
	DeletionScheduledFor  *time.Time `json:"deletion_scheduled_for,omitempty"`
	ConsentVersion        string     `json:"consent_version"`
	ConsentGivenAt        time.Time  `json:"consent_given_at"`
}

// DeviceFingerprint for device identification and security
type DeviceFingerprint struct {
	ID              uuid.UUID              `json:"id"`
	DeviceID        string                 `json:"device_id"`
	DeviceType      string                 `json:"device_type"`
	OS              string                 `json:"os"`
	OSVersion       string                 `json:"os_version"`
	Browser         string                 `json:"browser,omitempty"`
	BrowserVersion  string                 `json:"browser_version,omitempty"`
	UserAgent       string                 `json:"user_agent"`
	IPAddress       string                 `json:"ip_address"`
	Location        *GeoLocation           `json:"location,omitempty"`
	IsTrusted       bool                   `json:"is_trusted"`
	RiskScore       float64                `json:"risk_score"`
	FirstSeenAt     time.Time              `json:"first_seen_at"`
	LastSeenAt      time.Time              `json:"last_seen_at"`
	Attributes      map[string]interface{} `json:"attributes"`
}

type GeoLocation struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
}

// BiometricToken for biometric authentication
type BiometricToken struct {
	ID              uuid.UUID `json:"id"`
	TokenHash       string    `json:"token_hash"`
	BiometricType   string    `json:"biometric_type"` // face, fingerprint, iris
	DeviceID        string    `json:"device_id"`
	PublicKey       string    `json:"public_key"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	UsageCount      int       `json:"usage_count"`
}

// FederatedIdentity for SSO and third-party authentication
type FederatedIdentity struct {
	ID          uuid.UUID `json:"id"`
	Provider    string    `json:"provider"` // google, apple, facebook, etc.
	ProviderID  string    `json:"provider_id"`
	Email       string    `json:"email"`
	IsVerified  bool      `json:"is_verified"`
	AccessToken string    `json:"access_token,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	TokenExpiry *time.Time `json:"token_expiry,omitempty"`
	LinkedAt    time.Time `json:"linked_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

// SecurityKey for hardware security keys (WebAuthn/FIDO2)
type SecurityKey struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	CredentialID    string    `json:"credential_id"`
	PublicKey       string    `json:"public_key"`
	AttestationType string    `json:"attestation_type"`
	Counter         uint32    `json:"counter"`
	CreatedAt       time.Time `json:"created_at"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
}

// TrustedDevice represents a device that has been marked as trusted
type TrustedDevice struct {
	ID              uuid.UUID `json:"id"`
	DeviceID        string    `json:"device_id"`
	DeviceName      string    `json:"device_name"`
	DeviceType      string    `json:"device_type"`
	TrustedAt       time.Time `json:"trusted_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
}

// SessionManagement for advanced session control
type SessionManagement struct {
	MaxConcurrentSessions int             `json:"max_concurrent_sessions"`
	SessionTimeout        int             `json:"session_timeout"` // in seconds
	IdleTimeout           int             `json:"idle_timeout"`
	RequireReauth         bool            `json:"require_reauth"`
	ActiveSessions        []ActiveSession `json:"active_sessions"`
}

type ActiveSession struct {
	SessionID     string     `json:"session_id"`
	DeviceID      string     `json:"device_id"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	Location      *GeoLocation `json:"location,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	ExpiresAt     time.Time  `json:"expires_at"`
}

// AnomalyDetection for ML-based fraud detection
type AnomalyDetection struct {
	Enabled               bool                   `json:"enabled"`
	SensitivityLevel      SensitivityLevel       `json:"sensitivity_level"`
	AnomalyScore          float64                `json:"anomaly_score"`
	LastAnomalyDetectedAt *time.Time             `json:"last_anomaly_detected_at,omitempty"`
	AnomalyCount          int                    `json:"anomaly_count"`
	BehaviorBaseline      map[string]interface{} `json:"behavior_baseline"`
	Alerts                []AnomalyAlert         `json:"alerts"`
}

type SensitivityLevel string

const (
	SensitivityLow    SensitivityLevel = "low"
	SensitivityMedium SensitivityLevel = "medium"
	SensitivityHigh   SensitivityLevel = "high"
)

type AnomalyAlert struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	DetectedAt  time.Time              `json:"detected_at"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// CrossPlatformActivity tracks user activity across platforms
type CrossPlatformActivity struct {
	SocialinkLastActive   *time.Time             `json:"socialink_last_active,omitempty"`
	VignetteLastActive    *time.Time             `json:"vignette_last_active,omitempty"`
	PreferredPlatform     *string                `json:"preferred_platform,omitempty"`
	ActivityScore         float64                `json:"activity_score"`
	EngagementMetrics     map[string]interface{} `json:"engagement_metrics"`
	CrossPlatformPosts    int                    `json:"cross_platform_posts"`
	CrossPlatformMessages int                    `json:"cross_platform_messages"`
}

// ComplianceData for regulatory compliance
type ComplianceData struct {
	Region              string                 `json:"region"`
	ApplicableRegulations []string             `json:"applicable_regulations"`
	ConsentRecords      []ConsentRecord        `json:"consent_records"`
	DataRetentionPolicy DataRetentionPolicy    `json:"data_retention_policy"`
	AuditTrail          []AuditEntry           `json:"audit_trail"`
}

type ConsentRecord struct {
	ID          uuid.UUID `json:"id"`
	ConsentType string    `json:"consent_type"`
	Version     string    `json:"version"`
	Granted     bool      `json:"granted"`
	GrantedAt   time.Time `json:"granted_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
}

type DataRetentionPolicy struct {
	ActiveAccountDays   int `json:"active_account_days"`
	DeletedAccountDays  int `json:"deleted_account_days"`
	InactiveAccountDays int `json:"inactive_account_days"`
}

type AuditEntry struct {
	ID        uuid.UUID              `json:"id"`
	Action    string                 `json:"action"`
	Actor     string                 `json:"actor"`
	Timestamp time.Time              `json:"timestamp"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Details   map[string]interface{} `json:"details"`
}

// Scan implementations for PostgreSQL JSON/JSONB
func (p *PlatformLinks) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &p)
}

func (p PlatformLinks) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (s *SecurityProfile) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &s)
}

func (s SecurityProfile) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (p *PrivacySettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &p)
}

func (p PrivacySettings) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (d *DataRights) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &d)
}

func (d DataRights) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DeviceFingerprint) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &d)
}

func (d DeviceFingerprint) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (s *SessionManagement) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &s)
}

func (s SessionManagement) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (a *AnomalyDetection) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &a)
}

func (a AnomalyDetection) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (c *CrossPlatformActivity) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &c)
}

func (c CrossPlatformActivity) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ComplianceData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &c)
}

func (c ComplianceData) Value() (driver.Value, error) {
	return json.Marshal(c)
}
