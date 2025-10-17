package model

import (
	"time"

	"github.com/google/uuid"
)

// TwoFactorAuth represents 2FA configuration for a user
type TwoFactorAuth struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Secret       string     `json:"-" db:"secret"` // TOTP secret, never expose
	IsEnabled    bool       `json:"is_enabled" db:"is_enabled"`
	BackupCodes  []string   `json:"-" db:"backup_codes"` // Encrypted backup codes
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	EnabledAt    *time.Time `json:"enabled_at,omitempty" db:"enabled_at"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
}

// Enable2FARequest represents request to enable 2FA
type Enable2FARequest struct {
	Code string `json:"code" binding:"required,len=6"` // TOTP code for verification
}

// Verify2FARequest represents request to verify 2FA code
type Verify2FARequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// TwoFactorSetupResponse contains QR code and backup codes
type TwoFactorSetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"-" db:"token"` // Hashed token
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RequestPasswordResetRequest represents password reset request
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents password reset with token
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

// AccountRecoveryMethod represents recovery method
type AccountRecoveryMethod struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"` // email, phone, security_question
	Value     string    `json:"value" db:"value"`
	IsVerified bool     `json:"is_verified" db:"is_verified"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
