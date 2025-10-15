package service

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"vignette/user-service/internal/model"
	"vignette/user-service/internal/repository"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TwoFactorService struct {
	twoFactorRepo *repository.TwoFactorRepository
	userRepo      *repository.UserRepository
}

func NewTwoFactorService(twoFactorRepo *repository.TwoFactorRepository, userRepo *repository.UserRepository) *TwoFactorService {
	return &TwoFactorService{
		twoFactorRepo: twoFactorRepo,
		userRepo:      userRepo,
	}
}

// SetupTwoFactor generates a new 2FA secret and backup codes
func (s *TwoFactorService) SetupTwoFactor(userID uuid.UUID) (*model.TwoFactorSetupResponse, error) {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Vignette",
		AccountName: user.Email,
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Generate backup codes
	backupCodes := generateBackupCodes(10)

	// Store 2FA configuration (not enabled yet)
	twoFactor := &model.TwoFactorAuth{
		ID:          uuid.New(),
		UserID:      userID,
		Secret:      key.Secret(),
		IsEnabled:   false,
		BackupCodes: backupCodes,
		CreatedAt:   time.Now(),
	}

	err = s.twoFactorRepo.Create(twoFactor)
	if err != nil {
		return nil, fmt.Errorf("failed to create 2FA config: %w", err)
	}

	return &model.TwoFactorSetupResponse{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

// EnableTwoFactor enables 2FA after code verification
func (s *TwoFactorService) EnableTwoFactor(userID uuid.UUID, code string) error {
	// Get 2FA config
	twoFactor, err := s.twoFactorRepo.FindByUserID(userID)
	if err != nil {
		return fmt.Errorf("2FA not set up")
	}

	// Verify code
	valid := totp.Validate(code, twoFactor.Secret)
	if !valid {
		return fmt.Errorf("invalid verification code")
	}

	// Enable 2FA
	now := time.Now()
	twoFactor.IsEnabled = true
	twoFactor.EnabledAt = &now

	err = s.twoFactorRepo.Update(twoFactor)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return nil
}

// VerifyTwoFactorCode verifies a 2FA code or backup code
func (s *TwoFactorService) VerifyTwoFactorCode(userID uuid.UUID, code string) (bool, error) {
	twoFactor, err := s.twoFactorRepo.FindByUserID(userID)
	if err != nil {
		return false, err
	}

	if !twoFactor.IsEnabled {
		return false, fmt.Errorf("2FA not enabled")
	}

	// Try TOTP code first
	valid := totp.Validate(code, twoFactor.Secret)
	if valid {
		// Update last used
		now := time.Now()
		twoFactor.LastUsedAt = &now
		_ = s.twoFactorRepo.Update(twoFactor)
		return true, nil
	}

	// Try backup codes
	for i, backupCode := range twoFactor.BackupCodes {
		if backupCode == code {
			// Remove used backup code
			twoFactor.BackupCodes = append(twoFactor.BackupCodes[:i], twoFactor.BackupCodes[i+1:]...)
			_ = s.twoFactorRepo.Update(twoFactor)
			return true, nil
		}
	}

	return false, fmt.Errorf("invalid 2FA code")
}

// DisableTwoFactor disables 2FA for a user
func (s *TwoFactorService) DisableTwoFactor(userID uuid.UUID) error {
	return s.twoFactorRepo.DeleteByUserID(userID)
}

// IsTwoFactorEnabled checks if 2FA is enabled for a user
func (s *TwoFactorService) IsTwoFactorEnabled(userID uuid.UUID) (bool, error) {
	twoFactor, err := s.twoFactorRepo.FindByUserID(userID)
	if err != nil {
		return false, nil // 2FA not set up
	}
	return twoFactor.IsEnabled, nil
}

// generateBackupCodes generates random backup codes
func generateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		codes[i] = generateBackupCode()
	}
	return codes
}

// generateBackupCode generates a single random backup code
func generateBackupCode() string {
	b := make([]byte, 10)
	rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)[:16]
}
