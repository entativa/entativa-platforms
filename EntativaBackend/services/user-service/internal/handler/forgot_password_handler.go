package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ForgotPasswordRequest represents the forgot password request payload
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the reset password request payload
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ForgotPasswordResponse represents the forgot password response
type ForgotPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PasswordResetToken represents a password reset token in the database
type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// HandleForgotPassword handles the forgot password request
func (h *AuthHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate email format
	if !isValidEmail(req.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}
	
	// Find user by email
	user, err := h.userRepo.FindByEmail(r.Context(), req.Email)
	if err != nil {
		// Don't reveal if user exists or not for security
		respondWithJSON(w, http.StatusOK, ForgotPasswordResponse{
			Success: true,
			Message: "If an account exists with this email, you will receive a password reset link shortly.",
		})
		return
	}
	
	// Generate secure reset token
	token, err := generateSecureToken(32)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate reset token")
		return
	}
	
	// Store reset token in database (expires in 1 hour)
	resetToken := PasswordResetToken{
		ID:        generateUUID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}
	
	if err := h.tokenRepo.CreateResetToken(r.Context(), &resetToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create reset token")
		return
	}
	
	// Send password reset email
	resetLink := buildResetLink(token)
	go h.emailService.SendPasswordResetEmail(user.Email, user.FirstName, resetLink)
	
	// Log the reset attempt
	h.auditLog.LogPasswordResetRequest(user.ID, r.RemoteAddr)
	
	respondWithJSON(w, http.StatusOK, ForgotPasswordResponse{
		Success: true,
		Message: "If an account exists with this email, you will receive a password reset link shortly.",
	})
}

// HandleResetPassword handles the password reset with token
func (h *AuthHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate password strength
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	// Find and validate reset token
	resetToken, err := h.tokenRepo.FindResetToken(r.Context(), req.Token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}
	
	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		respondWithError(w, http.StatusBadRequest, "Reset token has expired")
		return
	}
	
	// Check if token was already used
	if resetToken.Used {
		respondWithError(w, http.StatusBadRequest, "Reset token has already been used")
		return
	}
	
	// Hash the new password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	
	// Update user password
	if err := h.userRepo.UpdatePassword(r.Context(), resetToken.UserID, hashedPassword); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}
	
	// Mark token as used
	if err := h.tokenRepo.MarkTokenAsUsed(r.Context(), resetToken.ID); err != nil {
		// Log error but don't fail the request
		h.logger.Error("Failed to mark token as used", err)
	}
	
	// Invalidate all existing sessions for this user (force re-login)
	if err := h.sessionRepo.InvalidateAllUserSessions(r.Context(), resetToken.UserID); err != nil {
		// Log error but don't fail the request
		h.logger.Error("Failed to invalidate sessions", err)
	}
	
	// Log the password reset
	h.auditLog.LogPasswordReset(resetToken.UserID, r.RemoteAddr)
	
	respondWithJSON(w, http.StatusOK, ForgotPasswordResponse{
		Success: true,
		Message: "Your password has been successfully reset. Please log in with your new password.",
	})
}

// HandleVerifyResetToken verifies if a reset token is valid
func (h *AuthHandler) HandleVerifyResetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}
	
	// Find reset token
	resetToken, err := h.tokenRepo.FindResetToken(r.Context(), token)
	if err != nil {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"valid": false,
			"message": "Invalid reset token",
		})
		return
	}
	
	// Check if expired
	if time.Now().After(resetToken.ExpiresAt) {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"valid": false,
			"message": "Reset token has expired",
		})
		return
	}
	
	// Check if already used
	if resetToken.Used {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"valid": false,
			"message": "Reset token has already been used",
		})
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"valid": true,
		"message": "Token is valid",
		"expires_at": resetToken.ExpiresAt,
	})
}

// Helper functions

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func buildResetLink(token string) string {
	// In production, use your actual domain
	baseURL := "https://app.entativa.com"
	return baseURL + "/reset-password?token=" + token
}

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return &ValidationError{"Password must be at least 8 characters long"}
	}
	
	hasUpper := false
	hasLower := false
	hasNumber := false
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}
	
	if !hasUpper {
		return &ValidationError{"Password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return &ValidationError{"Password must contain at least one lowercase letter"}
	}
	if !hasNumber {
		return &ValidationError{"Password must contain at least one number"}
	}
	
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
