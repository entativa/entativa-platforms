package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"user-service/internal/config"
	"user-service/internal/logger"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/internal/util"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userRepo     *repository.UserRepository
	sessionRepo  *repository.SessionRepository
	tokenRepo    *repository.TokenRepository
	emailService *service.EmailService
	auditLog     *service.AuditLog
	logger       *logger.Logger
	config       *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	tokenRepo *repository.TokenRepository,
	emailService *service.EmailService,
	auditLog *service.AuditLog,
	logger *logger.Logger,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenRepo:    tokenRepo,
		emailService: emailService,
		auditLog:     auditLog,
		logger:       logger,
		config:       cfg,
	}
}

// SignUpRequest represents sign up request payload
type SignUpRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Birthday  string `json:"birthday"`
	Gender    string `json:"gender"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID                string  `json:"id"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Email             string  `json:"email"`
	Username          string  `json:"username"`
	Birthday          *string `json:"birthday,omitempty"`
	Gender            *string `json:"gender,omitempty"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	CoverPhotoURL     *string `json:"cover_photo_url,omitempty"`
	IsActive          bool    `json:"is_active"`
	CreatedAt         string  `json:"created_at"`
}

// AuthResponse represents authentication response with token
type AuthResponse struct {
	User        *UserResponse `json:"user"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int           `json:"expires_in"`
}

// HandleSignUp handles user sign up
func (h *AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate input
	if err := util.ValidateName(req.FirstName, "First name"); err != nil {
		util.RespondWithValidationError(w, "first_name", err.Error())
		return
	}
	
	if err := util.ValidateName(req.LastName, "Last name"); err != nil {
		util.RespondWithValidationError(w, "last_name", err.Error())
		return
	}
	
	if err := util.ValidateEmail(req.Email); err != nil {
		util.RespondWithValidationError(w, "email", err.Error())
		return
	}
	
	if err := util.ValidatePasswordStrength(req.Password); err != nil {
		util.RespondWithValidationError(w, "password", err.Error())
		return
	}
	
	// Check if email already exists
	exists, _ := h.userRepo.CheckEmailExists(r.Context(), req.Email)
	if exists {
		util.RespondWithError(w, http.StatusConflict, "Email already registered")
		return
	}
	
	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", err)
		util.RespondWithInternalError(w, "Failed to create account")
		return
	}
	
	// Generate username from email
	username := generateUsernameFromEmail(req.Email)
	
	// Ensure username is unique
	username, err = h.ensureUniqueUsername(r.Context(), username)
	if err != nil {
		h.logger.Error("Failed to generate unique username", err)
		util.RespondWithInternalError(w, "Failed to create account")
		return
	}
	
	// Parse birthday
	var birthday *time.Time
	if req.Birthday != "" {
		parsedBirthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthday = &parsedBirthday
		}
	}
	
	// Create user
	user := &repository.User{
		ID:           util.GenerateUUID(),
		FirstName:    util.SanitizeInput(req.FirstName),
		LastName:     util.SanitizeInput(req.LastName),
		Email:        req.Email,
		Username:     username,
		PasswordHash: hashedPassword,
		Birthday:     birthday,
		Gender:       &req.Gender,
		IsActive:     true,
		IsDeleted:    false,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	if err := h.userRepo.CreateUser(r.Context(), user); err != nil {
		h.logger.Error("Failed to create user", err)
		util.RespondWithInternalError(w, "Failed to create account")
		return
	}
	
	// Generate access token
	accessToken, err := util.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate access token", err)
		util.RespondWithInternalError(w, "Failed to create session")
		return
	}
	
	// Generate refresh token
	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", err)
		util.RespondWithInternalError(w, "Failed to create session")
		return
	}
	
	// Create session
	session := &repository.Session{
		ID:           util.GenerateUUID(),
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceInfo:   r.UserAgent(),
		IPAddress:    getIPAddress(r),
		UserAgent:    r.UserAgent(),
		ExpiresAt:    time.Now().Add(h.config.Security.SessionExpiry),
		CreatedAt:    time.Now(),
		LastActiveAt: time.Now(),
	}
	
	if err := h.sessionRepo.CreateSession(r.Context(), session); err != nil {
		h.logger.Error("Failed to create session", err)
		util.RespondWithInternalError(w, "Failed to create session")
		return
	}
	
	// Send welcome email (async)
	go h.emailService.SendWelcomeEmail(user.Email, user.FirstName)
	
	// Log sign up
	h.auditLog.LogSignUp(user.ID, getIPAddress(r), r.UserAgent())
	
	// Return response
	util.RespondWithCreated(w, "Account created successfully! Welcome to Entativa!", AuthResponse{
		User:        mapUserToResponse(user),
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 24 hours
	})
}

// HandleLogin handles user login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Find user by email or username
	user, err := h.userRepo.FindByEmailOrUsername(r.Context(), req.EmailOrUsername)
	if err != nil {
		h.auditLog.LogFailedLogin(req.EmailOrUsername, getIPAddress(r), "user_not_found")
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	
	// Check if user is active
	if !user.IsActive {
		util.RespondWithError(w, http.StatusForbidden, "Account is deactivated")
		return
	}
	
	// Verify password
	if !util.ComparePassword(user.PasswordHash, req.Password) {
		h.auditLog.LogFailedLogin(req.EmailOrUsername, getIPAddress(r), "invalid_password")
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	
	// Generate tokens
	accessToken, err := util.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate access token", err)
		util.RespondWithInternalError(w, "Failed to create session")
		return
	}
	
	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", err)
		util.RespondWithInternalError(w, "Failed to create session")
		return
	}
	
	// Create session
	session := &repository.Session{
		ID:           util.GenerateUUID(),
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceInfo:   r.UserAgent(),
		IPAddress:    getIPAddress(r),
		UserAgent:    r.UserAgent(),
		ExpiresAt:    time.Now().Add(h.config.Security.SessionExpiry),
		CreatedAt:    time.Now(),
		LastActiveAt: time.Now(),
	}
	
	if err := h.sessionRepo.CreateSession(r.Context(), session); err != nil {
		h.logger.Error("Failed to create session", err)
		util.RespondWithInternalError(w, "Failed to create session")
		return
	}
	
	// Update last login
	if err := h.userRepo.UpdateLastLogin(r.Context(), user.ID); err != nil {
		h.logger.Warn("Failed to update last login", err)
	}
	
	// Log login
	h.auditLog.LogLogin(user.ID, getIPAddress(r), r.UserAgent())
	
	// Return response
	util.RespondWithSuccess(w, "Login successful! Welcome back!", AuthResponse{
		User:        mapUserToResponse(user),
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   86400,
	})
}

// HandleGetCurrentUser returns the current authenticated user
func (h *AuthHandler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}
	
	util.RespondWithSuccess(w, "", mapUserToResponse(user))
}

// HandleLogout handles user logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get token from header
	authHeader := r.Header.Get("Authorization")
	token, err := util.ExtractTokenFromHeader(authHeader)
	if err != nil {
		util.RespondWithUnauthorized(w, "Invalid authorization header")
		return
	}
	
	// Get user from context
	user, ok := r.Context().Value("user").(*repository.User)
	if ok {
		h.auditLog.LogLogout(user.ID, getIPAddress(r))
	}
	
	// Delete session
	if err := h.sessionRepo.DeleteSession(r.Context(), token); err != nil {
		h.logger.Warn("Failed to delete session", err)
	}
	
	util.RespondWithSuccess(w, "Logged out successfully", nil)
}

// HandleRefreshToken refreshes an access token
func (h *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate refresh token
	claims, err := util.ParseAccessToken(req.RefreshToken)
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	
	// Get user
	user, err := h.userRepo.FindByID(r.Context(), claims.UserID)
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}
	
	// Generate new access token
	accessToken, err := util.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate access token", err)
		util.RespondWithInternalError(w, "Failed to refresh token")
		return
	}
	
	util.RespondWithSuccess(w, "Token refreshed successfully", map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   86400,
	})
}

// HandleGetUser gets a user by ID
func (h *AuthHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// Implementation placeholder
	util.RespondWithError(w, http.StatusNotImplemented, "Not implemented")
}

// HandleUpdateUser updates a user's profile
func (h *AuthHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation placeholder
	util.RespondWithError(w, http.StatusNotImplemented, "Not implemented")
}

// HandleDeleteUser deletes a user account
func (h *AuthHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Implementation placeholder
	util.RespondWithError(w, http.StatusNotImplemented, "Not implemented")
}

// Helper methods

func (h *AuthHandler) generateAccessToken(userID string) (string, error) {
	user, err := h.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return "", err
	}
	return util.GenerateAccessToken(user.ID, user.Username, user.Email)
}

func (h *AuthHandler) parseAccessToken(token string) (*util.TokenClaims, error) {
	return util.ParseAccessToken(token)
}

func (h *AuthHandler) ensureUniqueUsername(ctx context.Context, baseUsername string) (string, error) {
	username := baseUsername
	counter := 1
	
	for {
		exists, err := h.userRepo.CheckUsernameExists(ctx, username)
		if err != nil {
			return "", err
		}
		
		if !exists {
			return username, nil
		}
		
		// Add number suffix
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
		
		// Safety check
		if counter > 9999 {
			return "", fmt.Errorf("failed to generate unique username")
		}
	}
}

// Utility functions

func mapUserToResponse(user *repository.User) *UserResponse {
	var birthday *string
	if user.Birthday != nil {
		birthdayStr := user.Birthday.Format("2006-01-02")
		birthday = &birthdayStr
	}
	
	return &UserResponse{
		ID:                user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		Username:          user.Username,
		Birthday:          birthday,
		Gender:            user.Gender,
		ProfilePictureURL: user.ProfilePictureURL,
		CoverPhotoURL:     user.CoverPhotoURL,
		IsActive:          user.IsActive,
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
	}
}

func generateUsernameFromEmail(email string) string {
	// Extract part before @
	atIndex := 0
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}
	
	if atIndex == 0 {
		return "user"
	}
	
	username := email[:atIndex]
	
	// Remove special characters except dots and underscores
	cleanUsername := ""
	for _, char := range username {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
		   (char >= '0' && char <= '9') || char == '.' || char == '_' {
			cleanUsername += string(char)
		}
	}
	
	if cleanUsername == "" {
		return "user"
	}
	
	return cleanUsername
}

func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header (proxy/load balancer)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	
	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	
	// Fallback to RemoteAddr
	return r.RemoteAddr
}

func generateRefreshToken() string {
	token, _ := util.GenerateSecureToken(32)
	return token
}
