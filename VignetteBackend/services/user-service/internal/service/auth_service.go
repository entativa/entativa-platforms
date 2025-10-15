package service

import (
	"errors"
	"fmt"
	"time"

	"vignette/user-service/internal/config"
	"vignette/user-service/internal/model"
	"vignette/user-service/internal/repository"
	"vignette/user-service/internal/util"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid username/email or password")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrInvalidUsername    = errors.New("invalid username format")
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	config      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		config:      cfg,
	}
}

// Signup registers a new user (Meta-level: instant access, no verification)
func (s *AuthService) Signup(req *model.SignupRequest, ipAddress, userAgent string) (*model.AuthResponse, error) {
	// Validate email
	if !util.IsValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate username
	if !util.IsValidUsername(req.Username) {
		return nil, ErrInvalidUsername
	}

	// Check if email exists
	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, repository.ErrEmailExists
	}

	// Check if username exists
	usernameExists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if usernameExists {
		return nil, repository.ErrUsernameExists
	}

	// Validate password
	if !util.ValidatePassword(req.Password) {
		return nil, fmt.Errorf("password must be 8-128 characters")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		ID:             uuid.New(),
		Username:       req.Username,
		Email:          req.Email,
		FullName:       req.FullName,
		Password:       hashedPassword,
		IsPrivate:      false, // Default to public account
		IsVerified:     false,
		IsActive:       true, // Instantly active - Meta-level approach
		IsDeleted:      false,
		FollowersCount: 0,
		FollowingCount: 0,
		PostsCount:     0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	accessToken, err := util.GenerateToken(
		user.ID,
		user.Email,
		user.Username,
		user.FullName,
		s.config.JWT.SecretKey,
		s.config.JWT.AccessTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session
	session := &model.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(), // Simple refresh token for now
		DeviceInfo:   "Web",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenTTL),
		CreatedAt:    time.Now(),
		LastActiveAt: time.Now(),
	}

	err = s.sessionRepo.Create(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return &model.AuthResponse{
		User:        user.ToUserResponse(),
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

// Login authenticates a user with username/email and password
func (s *AuthService) Login(req *model.LoginRequest, ipAddress, userAgent string) (*model.AuthResponse, error) {
	// Find user by email or username
	var user *model.User
	var err error

	if util.IsValidEmail(req.UsernameOrEmail) {
		user, err = s.userRepo.FindByEmail(req.UsernameOrEmail)
	} else {
		user, err = s.userRepo.FindByUsername(req.UsernameOrEmail)
	}

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if !util.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	accessToken, err := util.GenerateToken(
		user.ID,
		user.Email,
		user.Username,
		user.FullName,
		s.config.JWT.SecretKey,
		s.config.JWT.AccessTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session
	session := &model.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		DeviceInfo:   "Web",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessTokenTTL),
		CreatedAt:    time.Now(),
		LastActiveAt: time.Now(),
	}

	err = s.sessionRepo.Create(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return &model.AuthResponse{
		User:        user.ToUserResponse(),
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

// Logout logs out a user by deleting their sessions
func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

// ValidateToken validates a JWT token and returns user claims
func (s *AuthService) ValidateToken(token string) (*util.JWTClaims, error) {
	return util.ValidateToken(token, s.config.JWT.SecretKey)
}
