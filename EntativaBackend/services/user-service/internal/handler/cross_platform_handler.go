package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CrossPlatformSignInRequest represents cross-platform sign-in request
type CrossPlatformSignInRequest struct {
	Platform    string `json:"platform" validate:"required,oneof=vignette entativa"`
	AccessToken string `json:"access_token" validate:"required"`
}

// CrossPlatformSignInResponse represents cross-platform sign-in response
type CrossPlatformSignInResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    *CrossPlatformSignInData `json:"data,omitempty"`
}

// CrossPlatformSignInData contains user and token data
type CrossPlatformSignInData struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int           `json:"expires_in"`
	IsNewAccount bool          `json:"is_new_account"`
}

// VignetteUserInfo represents user info from Vignette API
type VignetteUserInfo struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	FullName          string `json:"full_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
	IsVerified        bool   `json:"is_verified"`
}

// HandleCrossPlatformSignIn handles signing in with another platform's credentials
func (h *AuthHandler) HandleCrossPlatformSignIn(w http.ResponseWriter, r *http.Request) {
	var req CrossPlatformSignInRequest
	
	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate platform
	if req.Platform != "vignette" && req.Platform != "entativa" {
		respondWithError(w, http.StatusBadRequest, "Invalid platform. Must be 'vignette' or 'entativa'")
		return
	}
	
	// Verify token with the other platform
	var userInfo *VignetteUserInfo
	var err error
	
	if req.Platform == "vignette" {
		userInfo, err = h.verifyVignetteToken(r.Context(), req.AccessToken)
	} else {
		// For Entativa tokens (when Vignette is calling)
		userInfo, err = h.verifyEntativaToken(r.Context(), req.AccessToken)
	}
	
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token from "+req.Platform)
		return
	}
	
	// Check if user already exists in our system
	existingUser, err := h.userRepo.FindByEmail(r.Context(), userInfo.Email)
	
	isNewAccount := false
	var user *User
	
	if err != nil || existingUser == nil {
		// Create new user from platform data
		isNewAccount = true
		user, err = h.createUserFromCrossPlatform(r.Context(), userInfo, req.Platform)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user account")
			return
		}
	} else {
		user = existingUser
		// Update cross-platform link
		if err := h.userRepo.LinkCrossPlatformAccount(r.Context(), user.ID, req.Platform, userInfo.ID); err != nil {
			// Log but don't fail
			h.logger.Warn("Failed to link cross-platform account", err)
		}
	}
	
	// Generate new access token for our platform
	accessToken, err := h.generateAccessToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}
	
	// Create session
	session := &Session{
		ID:           generateUUID(),
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: generateRefreshToken(),
		DeviceInfo:   r.UserAgent(),
		IPAddress:    r.RemoteAddr,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}
	
	if err := h.sessionRepo.CreateSession(r.Context(), session); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}
	
	// Log the cross-platform sign-in
	h.auditLog.LogCrossPlatformSignIn(user.ID, req.Platform, r.RemoteAddr)
	
	// Return response
	respondWithJSON(w, http.StatusOK, CrossPlatformSignInResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully signed in with %s", req.Platform),
		Data: &CrossPlatformSignInData{
			User:         mapUserToResponse(user),
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    86400, // 24 hours
			IsNewAccount: isNewAccount,
		},
	})
}

// HandleCheckCrossPlatformAccount checks if user exists on platform
func (h *AuthHandler) HandleCheckCrossPlatformAccount(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "Email parameter is required")
		return
	}
	
	// Check if user exists
	user, err := h.userRepo.FindByEmail(r.Context(), email)
	exists := err == nil && user != nil
	
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"exists": exists,
		"email":  email,
	})
}

// verifyVignetteToken verifies a token with Vignette API
func (h *AuthHandler) verifyVignetteToken(ctx context.Context, token string) (*VignetteUserInfo, error) {
	// Call Vignette API to verify token and get user info
	vignetteAPIURL := getEnvOrDefault("VIGNETTE_API_URL", "http://localhost:8002/api/v1")
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", vignetteAPIURL+"/auth/me", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token with Vignette: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}
	
	var result struct {
		Success bool              `json:"success"`
		Data    *VignetteUserInfo `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Vignette response: %w", err)
	}
	
	if !result.Success || result.Data == nil {
		return nil, fmt.Errorf("invalid response from Vignette")
	}
	
	return result.Data, nil
}

// verifyEntativaToken verifies an Entativa token (when Vignette is calling)
func (h *AuthHandler) verifyEntativaToken(ctx context.Context, token string) (*VignetteUserInfo, error) {
	// Parse JWT token
	claims, err := h.parseAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	
	// Get user from database
	user, err := h.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// Map Entativa user to VignetteUserInfo format
	return &VignetteUserInfo{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		FullName:          fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		ProfilePictureURL: user.ProfilePictureURL,
		IsVerified:        user.IsVerified,
	}, nil
}

// createUserFromCrossPlatform creates a new user from cross-platform data
func (h *AuthHandler) createUserFromCrossPlatform(
	ctx context.Context,
	userInfo *VignetteUserInfo,
	sourcePlatform string,
) (*User, error) {
	// Split full name into first and last name
	firstName, lastName := splitFullName(userInfo.FullName)
	
	// Generate a random password (user won't use it, they'll use cross-platform sign-in)
	randomPassword := generateRandomPassword(32)
	hashedPassword, err := hashPassword(randomPassword)
	if err != nil {
		return nil, err
	}
	
	user := &User{
		ID:                generateUUID(),
		FirstName:         firstName,
		LastName:          lastName,
		Email:             userInfo.Email,
		Username:          userInfo.Username,
		PasswordHash:      hashedPassword,
		ProfilePictureURL: userInfo.ProfilePictureURL,
		IsActive:          true,
		IsVerified:        userInfo.IsVerified,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	
	if err := h.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	
	// Link the cross-platform account
	if err := h.userRepo.LinkCrossPlatformAccount(ctx, user.ID, sourcePlatform, userInfo.ID); err != nil {
		h.logger.Warn("Failed to link cross-platform account", err)
	}
	
	// Send welcome email (async)
	go h.emailService.SendWelcomeEmail(user.Email, user.FirstName)
	
	return user, nil
}

// Helper functions

func splitFullName(fullName string) (string, string) {
	parts := []rune(fullName)
	spaceIndex := -1
	
	for i, char := range parts {
		if char == ' ' {
			spaceIndex = i
			break
		}
	}
	
	if spaceIndex == -1 {
		return fullName, ""
	}
	
	firstName := string(parts[:spaceIndex])
	lastName := string(parts[spaceIndex+1:])
	
	return firstName, lastName
}

func generateRandomPassword(length int) string {
	// Generate a secure random password
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
