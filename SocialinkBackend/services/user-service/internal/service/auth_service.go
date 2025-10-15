package service

import (
	"errors"
	"fmt"
	"time"

	"socialink/user-service/internal/config"
	"socialink/user-service/internal/model"
	"socialink/user-service/internal/repository"
	"socialink/user-service/internal/util"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email/username or password")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrInvalidBirthday    = errors.New("invalid birthday - must be at least 13 years old")
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

	// Check if email exists
	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, repository.ErrEmailExists
	}

	// Parse birthday
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return nil, fmt.Errorf("invalid birthday format, use YYYY-MM-DD")
	}

	// Validate birthday (must be at least 13 years old)
	if !util.IsValidBirthday(birthday) {
		return nil, ErrInvalidBirthday
	}

	// Validate password
	if !util.ValidatePassword(req.Password) {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate unique username
	username := util.GenerateUsername(req.FirstName, req.LastName)
	
	// Ensure username is unique
	for {
		exists, err := s.userRepo.UsernameExists(username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if !exists {
			break
		}
		// Add random suffix if username exists
		username = util.GenerateUsername(req.FirstName, req.LastName)
	}

	// Create user
	user := &model.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Username:  username,
		Password:  hashedPassword,
		Birthday:  birthday,
		Gender:    req.Gender,
		IsActive:  true, // Instantly active - Meta-level approach
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		user.FirstName,
		user.LastName,
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

// Login authenticates a user with email/username and password
func (s *AuthService) Login(req *model.LoginRequest, ipAddress, userAgent string) (*model.AuthResponse, error) {
	// Find user by email or username
	var user *model.User
	var err error

	if util.IsValidEmail(req.EmailOrUsername) {
		user, err = s.userRepo.FindByEmail(req.EmailOrUsername)
	} else {
		user, err = s.userRepo.FindByUsername(req.EmailOrUsername)
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
		user.FirstName,
		user.LastName,
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
