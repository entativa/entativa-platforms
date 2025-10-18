package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"socialink/user-service/internal/model"
	"socialink/user-service/internal/repository"
	"socialink/user-service/internal/util"

	"github.com/google/uuid"
)

type PasswordResetService struct {
	userRepo              *repository.UserRepository
	passwordResetRepo     *repository.PasswordResetRepository
	kafkaProducer         *KafkaProducer
}

func NewPasswordResetService(
	userRepo *repository.UserRepository,
	passwordResetRepo *repository.PasswordResetRepository,
	kafkaProducer *KafkaProducer,
) *PasswordResetService {
	return &PasswordResetService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		kafkaProducer:     kafkaProducer,
	}
}

// RequestPasswordReset creates a password reset token and sends email
func (s *PasswordResetService) RequestPasswordReset(email string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if user exists or not (security)
		return nil
	}

	// Generate secure random token
	token := generateSecureToken()
	hashedToken, err := util.HashPassword(token)
	if err != nil {
		return fmt.Errorf("failed to hash token: %w", err)
	}

	// Create password reset token (expires in 1 hour)
	resetToken := &model.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	err = s.passwordResetRepo.Create(resetToken)
	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send password reset email via Kafka event
	event := map[string]interface{}{
		"event_type": "password_reset_requested",
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"token":      token, // Send unhashed token in email
		"expires_at": resetToken.ExpiresAt,
		"timestamp":  time.Now(),
	}

	if s.kafkaProducer != nil {
		_ = s.kafkaProducer.PublishEvent("user-events", event)
	}

	return nil
}

// ResetPassword resets password using token
func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	// Validate password
	if !util.ValidatePassword(newPassword) {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Find valid reset token
	resetTokens, err := s.passwordResetRepo.FindValidTokens()
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	var validToken *model.PasswordResetToken
	for _, rt := range resetTokens {
		if util.CheckPassword(token, rt.Token) {
			validToken = rt
			break
		}
	}

	if validToken == nil {
		return fmt.Errorf("invalid or expired token")
	}

	// Get user
	user, err := s.userRepo.FindByID(validToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user.Password = hashedPassword
	err = s.userRepo.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	now := time.Now()
	validToken.UsedAt = &now
	_ = s.passwordResetRepo.MarkAsUsed(validToken.ID)

	// Send password changed event
	if s.kafkaProducer != nil {
		event := map[string]interface{}{
			"event_type": "password_changed",
			"user_id":    user.ID.String(),
			"email":      user.Email,
			"timestamp":  time.Now(),
		}
		_ = s.kafkaProducer.PublishEvent("user-events", event)
	}

	return nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
