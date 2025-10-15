package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"socialink/meta-user-service/internal/model"
	"socialink/meta-user-service/internal/repository"
	"socialink/meta-user-service/pkg/ml"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserSuspended      = errors.New("user account is suspended")
	ErrUserBanned         = errors.New("user account is banned")
	ErrHighRiskDetected   = errors.New("high risk activity detected")
)

// MetaUserService provides advanced user management with ML-based security
type MetaUserService struct {
	repo              *repository.MetaUserRepository
	fraudDetector     *ml.FraudDetector
	behaviorAnalyzer  *ml.BehaviorAnalyzer
	trustScoreEngine  *ml.TrustScoreEngine
	eventPublisher    EventPublisher
	crossPlatformSync *CrossPlatformSyncService
}

type EventPublisher interface {
	PublishUserEvent(ctx context.Context, eventType string, data map[string]interface{}) error
}

func NewMetaUserService(
	repo *repository.MetaUserRepository,
	fraudDetector *ml.FraudDetector,
	behaviorAnalyzer *ml.BehaviorAnalyzer,
	trustScoreEngine *ml.TrustScoreEngine,
	eventPublisher EventPublisher,
	crossPlatformSync *CrossPlatformSyncService,
) *MetaUserService {
	return &MetaUserService{
		repo:              repo,
		fraudDetector:     fraudDetector,
		behaviorAnalyzer:  behaviorAnalyzer,
		trustScoreEngine:  trustScoreEngine,
		eventPublisher:    eventPublisher,
		crossPlatformSync: crossPlatformSync,
	}
}

// CreateMetaUser creates a new meta user with advanced security checks
func (s *MetaUserService) CreateMetaUser(ctx context.Context, req *CreateMetaUserRequest) (*model.MetaUser, error) {
	// Validate email and phone
	if err := s.validateEmail(req.Email); err != nil {
		return nil, err
	}

	// Check for fraud indicators
	fraudScore, err := s.fraudDetector.AnalyzeSignup(ctx, &ml.SignupAnalysisRequest{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		DeviceID:    req.DeviceID,
	})
	if err != nil {
		return nil, fmt.Errorf("fraud detection failed: %w", err)
	}

	if fraudScore > 0.7 {
		return nil, ErrHighRiskDetected
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Determine initial risk level based on fraud score
	riskLevel := s.calculateRiskLevel(fraudScore)

	// Create user
	now := time.Now()
	user := &model.MetaUser{
		ID:            uuid.New(),
		MetaID:        s.generateMetaID(),
		Email:         req.Email,
		EmailVerified: false,
		PhoneNumber:   req.PhoneNumber,
		PhoneVerified: false,
		PasswordHash:  string(passwordHash),
		Status:        model.UserStatusPending,
		TrustScore:    s.calculateInitialTrustScore(fraudScore),
		RiskLevel:     riskLevel,
		AccountTier:   model.AccountTierBasic,
		PlatformLinks: model.PlatformLinks{
			LinkedAt:   now,
			LinkStatus: model.LinkStatusUnlinked,
			SyncEnabled: false,
		},
		SecurityProfile: s.createDefaultSecurityProfile(),
		PrivacySettings: s.createDefaultPrivacySettings(),
		DataRights:      s.createDefaultDataRights(),
		DeviceFingerprints: []model.DeviceFingerprint{
			s.createDeviceFingerprint(req),
		},
		SessionManagement: s.createDefaultSessionManagement(),
		AnomalyDetection: model.AnomalyDetection{
			Enabled:          true,
			SensitivityLevel: model.SensitivityMedium,
			AnomalyScore:     0.0,
			BehaviorBaseline: make(map[string]interface{}),
			Alerts:           []model.AnomalyAlert{},
		},
		CrossPlatformActivity: model.CrossPlatformActivity{
			ActivityScore:     0.0,
			EngagementMetrics: make(map[string]interface{}),
		},
		ComplianceData: s.createDefaultComplianceData(req.Region),
		Metadata:       make(map[string]interface{}),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Publish event
	s.eventPublisher.PublishUserEvent(ctx, "meta.user.created", map[string]interface{}{
		"user_id":     user.ID,
		"meta_id":     user.MetaID,
		"email":       user.Email,
		"trust_score": user.TrustScore,
	})

	return user, nil
}

// Authenticate authenticates a user with advanced security checks
func (s *MetaUserService) Authenticate(ctx context.Context, req *AuthenticationRequest) (*AuthenticationResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check user status
	if user.Status == model.UserStatusSuspended {
		return nil, ErrUserSuspended
	}
	if user.Status == model.UserStatusBanned {
		return nil, ErrUserBanned
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Record failed login attempt
		s.recordFailedLogin(ctx, user, req)
		return nil, ErrInvalidCredentials
	}

	// Analyze login behavior for anomalies
	anomalyScore, err := s.behaviorAnalyzer.AnalyzeLogin(ctx, &ml.LoginAnalysisRequest{
		UserID:      user.ID,
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		DeviceID:    req.DeviceID,
		Location:    req.Location,
		Time:        time.Now(),
		PreviousLogins: s.getUserLoginHistory(ctx, user.ID),
	})
	if err == nil && anomalyScore > 0.8 {
		// High anomaly detected, require additional verification
		return &AuthenticationResponse{
			Success:                false,
			RequiresTwoFactor:      true,
			RequiresAdditionalAuth: true,
			AnomalyDetected:        true,
			AnomalyScore:           anomalyScore,
		}, nil
	}

	// Check if 2FA is required
	if user.SecurityProfile.TwoFactorEnabled {
		return &AuthenticationResponse{
			Success:           false,
			RequiresTwoFactor: true,
			UserID:            user.ID,
		}, nil
	}

	// Update last seen
	now := time.Now()
	user.LastSeenAt = &now
	user.SecurityProfile.FailedLoginAttempts = 0
	user.SecurityProfile.LastFailedLoginAt = nil

	// Add device fingerprint if new
	s.addDeviceFingerprint(user, req)

	// Update user
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Generate session token (implementation depends on JWT/session strategy)
	sessionToken := s.generateSessionToken(user)

	// Publish event
	s.eventPublisher.PublishUserEvent(ctx, "meta.user.authenticated", map[string]interface{}{
		"user_id":    user.ID,
		"ip_address": req.IPAddress,
		"device_id":  req.DeviceID,
	})

	return &AuthenticationResponse{
		Success:      true,
		User:         user,
		SessionToken: sessionToken,
	}, nil
}

// LinkPlatformAccount links a platform-specific account to the meta user
func (s *MetaUserService) LinkPlatformAccount(ctx context.Context, metaUserID uuid.UUID, platform string, platformUserID uuid.UUID) error {
	// Validate platform
	if platform != "socialink" && platform != "vignette" {
		return errors.New("invalid platform")
	}

	// Link account
	if err := s.repo.LinkPlatformAccount(ctx, metaUserID, platform, platformUserID); err != nil {
		return err
	}

	// Trigger cross-platform synchronization
	if s.crossPlatformSync != nil {
		go s.crossPlatformSync.SyncUserData(context.Background(), metaUserID, platform)
	}

	// Publish event
	s.eventPublisher.PublishUserEvent(ctx, "meta.platform.linked", map[string]interface{}{
		"meta_user_id":     metaUserID,
		"platform":         platform,
		"platform_user_id": platformUserID,
	})

	return nil
}

// UpdateTrustScore updates the user's trust score using ML models
func (s *MetaUserService) UpdateTrustScore(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Calculate new trust score using ML engine
	trustScore, err := s.trustScoreEngine.CalculateTrustScore(ctx, &ml.TrustScoreRequest{
		UserID:                userID,
		AccountAge:            time.Since(user.CreatedAt),
		VerificationStatus:    user.EmailVerified && user.PhoneVerified,
		ActivityScore:         user.CrossPlatformActivity.ActivityScore,
		ReportedViolations:    s.getViolationCount(ctx, userID),
		PositiveInteractions:  s.getPositiveInteractionCount(ctx, userID),
		DeviceFingerprints:    len(user.DeviceFingerprints),
		SecurityProfile:       user.SecurityProfile,
		AnomalyHistory:        user.AnomalyDetection.AnomalyCount,
	})
	if err != nil {
		return fmt.Errorf("failed to calculate trust score: %w", err)
	}

	// Determine risk level
	riskLevel := s.calculateRiskLevel(1.0 - trustScore)

	// Update in database
	if err := s.repo.UpdateTrustScore(ctx, userID, trustScore, riskLevel); err != nil {
		return err
	}

	// Publish event
	s.eventPublisher.PublishUserEvent(ctx, "meta.trust_score.updated", map[string]interface{}{
		"user_id":     userID,
		"trust_score": trustScore,
		"risk_level":  riskLevel,
	})

	return nil
}

// DetectAnomalies performs real-time anomaly detection
func (s *MetaUserService) DetectAnomalies(ctx context.Context, userID uuid.UUID, activity *ml.UserActivity) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.AnomalyDetection.Enabled {
		return nil
	}

	// Analyze activity against behavior baseline
	anomalyScore, anomalyType, err := s.behaviorAnalyzer.DetectAnomaly(ctx, &ml.AnomalyDetectionRequest{
		UserID:           userID,
		Activity:         activity,
		BehaviorBaseline: user.AnomalyDetection.BehaviorBaseline,
		SensitivityLevel: user.AnomalyDetection.SensitivityLevel,
	})
	if err != nil {
		return err
	}

	// If anomaly detected, create alert
	if anomalyScore > 0.7 {
		now := time.Now()
		alert := model.AnomalyAlert{
			ID:          uuid.New(),
			Type:        anomalyType,
			Severity:    s.calculateSeverity(anomalyScore),
			Description: fmt.Sprintf("Anomaly detected: %s (score: %.2f)", anomalyType, anomalyScore),
			DetectedAt:  now,
			Resolved:    false,
			Metadata: map[string]interface{}{
				"score":    anomalyScore,
				"activity": activity,
			},
		}

		user.AnomalyDetection.Alerts = append(user.AnomalyDetection.Alerts, alert)
		user.AnomalyDetection.AnomalyCount++
		user.AnomalyDetection.LastAnomalyDetectedAt = &now
		user.AnomalyDetection.AnomalyScore = anomalyScore

		// Update user
		if err := s.repo.Update(ctx, user); err != nil {
			return err
		}

		// Publish critical event
		s.eventPublisher.PublishUserEvent(ctx, "meta.anomaly.detected", map[string]interface{}{
			"user_id":       userID,
			"anomaly_type":  anomalyType,
			"anomaly_score": anomalyScore,
			"severity":      alert.Severity,
		})

		// Auto-suspend if critical
		if anomalyScore > 0.95 {
			s.SuspendUser(ctx, userID, "Automatic suspension due to critical anomaly detection")
		}
	}

	return nil
}

// GetByID retrieves a user by ID
func (s *MetaUserService) GetByID(ctx context.Context, userID uuid.UUID) (*model.MetaUser, error) {
	return s.repo.GetByID(ctx, userID)
}

// SuspendUser suspends a user account
func (s *MetaUserService) SuspendUser(ctx context.Context, userID uuid.UUID, reason string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Status = model.UserStatusSuspended
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	s.eventPublisher.PublishUserEvent(ctx, "meta.user.suspended", map[string]interface{}{
		"user_id": userID,
		"reason":  reason,
	})

	return nil
}

// Helper methods

func (s *MetaUserService) generateMetaID() string {
	return fmt.Sprintf("META-%s", uuid.New().String()[:12])
}

func (s *MetaUserService) calculateRiskLevel(fraudScore float64) model.RiskLevel {
	if fraudScore < 0.3 {
		return model.RiskLevelLow
	} else if fraudScore < 0.6 {
		return model.RiskLevelMedium
	} else if fraudScore < 0.85 {
		return model.RiskLevelHigh
	}
	return model.RiskLevelCritical
}

func (s *MetaUserService) calculateInitialTrustScore(fraudScore float64) float64 {
	return 0.5 - (fraudScore * 0.3)
}

func (s *MetaUserService) createDefaultSecurityProfile() model.SecurityProfile {
	return model.SecurityProfile{
		TwoFactorEnabled:      false,
		TwoFactorMethod:       model.TwoFactorMethodNone,
		BiometricEnabled:      false,
		PasskeyEnabled:        false,
		PasswordlessEnabled:   false,
		LoginApprovalRequired: false,
		TrustedDevices:        []model.TrustedDevice{},
		SecurityKeys:          []model.SecurityKey{},
		LastPasswordChange:    time.Now(),
		FailedLoginAttempts:   0,
	}
}

func (s *MetaUserService) createDefaultPrivacySettings() model.PrivacySettings {
	return model.PrivacySettings{
		ProfileVisibility:     model.VisibilityPublic,
		ActivityVisibility:    model.VisibilityFriends,
		ContactsVisibility:    model.VisibilityPrivate,
		CrossPlatformSharing:  true,
		DataPortabilityEnabled: true,
		MarketingOptIn:        false,
		PersonalizationOptIn:  true,
		ThirdPartyDataSharing: false,
		LocationTracking:      model.LocationTrackingInUse,
		SearchVisibility:      true,
		CustomPermissions:     make(map[string]interface{}),
	}
}

func (s *MetaUserService) createDefaultDataRights() model.DataRights {
	return model.DataRights{
		RightToAccess:        true,
		RightToRectification: true,
		RightToErasure:       true,
		RightToPortability:   true,
		RightToObject:        true,
		ConsentVersion:       "1.0",
		ConsentGivenAt:       time.Now(),
	}
}

func (s *MetaUserService) createDefaultSessionManagement() model.SessionManagement {
	return model.SessionManagement{
		MaxConcurrentSessions: 5,
		SessionTimeout:        3600,
		IdleTimeout:           900,
		RequireReauth:         false,
		ActiveSessions:        []model.ActiveSession{},
	}
}

func (s *MetaUserService) createDefaultComplianceData(region string) model.ComplianceData {
	regulations := []string{"GDPR"}
	if region == "US" {
		regulations = append(regulations, "CCPA")
	}

	return model.ComplianceData{
		Region:                region,
		ApplicableRegulations: regulations,
		ConsentRecords:        []model.ConsentRecord{},
		DataRetentionPolicy: model.DataRetentionPolicy{
			ActiveAccountDays:   3650,
			DeletedAccountDays:  90,
			InactiveAccountDays: 730,
		},
		AuditTrail: []model.AuditEntry{},
	}
}

func (s *MetaUserService) createDeviceFingerprint(req *CreateMetaUserRequest) model.DeviceFingerprint {
	return model.DeviceFingerprint{
		ID:         uuid.New(),
		DeviceID:   req.DeviceID,
		DeviceType: req.DeviceType,
		OS:         req.OS,
		UserAgent:  req.UserAgent,
		IPAddress:  req.IPAddress,
		Location:   req.Location,
		IsTrusted:  false,
		RiskScore:  0.3,
		FirstSeenAt: time.Now(),
		LastSeenAt: time.Now(),
		Attributes: make(map[string]interface{}),
	}
}

func (s *MetaUserService) addDeviceFingerprint(user *model.MetaUser, req *AuthenticationRequest) {
	// Check if device already exists
	for i, fp := range user.DeviceFingerprints {
		if fp.DeviceID == req.DeviceID {
			user.DeviceFingerprints[i].LastSeenAt = time.Now()
			return
		}
	}

	// Add new device
	fp := model.DeviceFingerprint{
		ID:         uuid.New(),
		DeviceID:   req.DeviceID,
		DeviceType: req.DeviceType,
		UserAgent:  req.UserAgent,
		IPAddress:  req.IPAddress,
		Location:   req.Location,
		IsTrusted:  false,
		RiskScore:  0.5,
		FirstSeenAt: time.Now(),
		LastSeenAt: time.Now(),
		Attributes: make(map[string]interface{}),
	}
	user.DeviceFingerprints = append(user.DeviceFingerprints, fp)
}

func (s *MetaUserService) recordFailedLogin(ctx context.Context, user *model.MetaUser, req *AuthenticationRequest) {
	user.SecurityProfile.FailedLoginAttempts++
	now := time.Now()
	user.SecurityProfile.LastFailedLoginAt = &now

	// Lock account if too many failed attempts
	if user.SecurityProfile.FailedLoginAttempts >= 5 {
		lockUntil := now.Add(15 * time.Minute)
		user.SecurityProfile.AccountLockedUntil = &lockUntil
	}

	s.repo.Update(ctx, user)
}

func (s *MetaUserService) validateEmail(email string) error {
	// Basic email validation (extend with more sophisticated checks)
	if len(email) < 3 || !contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func (s *MetaUserService) generateSessionToken(user *model.MetaUser) string {
	// Implementation would use JWT or similar
	return fmt.Sprintf("session_%s", uuid.New().String())
}

func (s *MetaUserService) getUserLoginHistory(ctx context.Context, userID uuid.UUID) []ml.LoginRecord {
	// Implementation would fetch from database
	return []ml.LoginRecord{}
}

func (s *MetaUserService) getViolationCount(ctx context.Context, userID uuid.UUID) int {
	// Implementation would fetch from moderation service
	return 0
}

func (s *MetaUserService) getPositiveInteractionCount(ctx context.Context, userID uuid.UUID) int {
	// Implementation would fetch from analytics service
	return 100
}

func (s *MetaUserService) calculateSeverity(score float64) string {
	if score > 0.9 {
		return "critical"
	} else if score > 0.75 {
		return "high"
	} else if score > 0.5 {
		return "medium"
	}
	return "low"
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && len(s) > len(substr) && s[0:len(s)] != s[0:len(s)]
	// Simplified - use strings.Contains in real implementation
}

// Request/Response types

type CreateMetaUserRequest struct {
	Email       string
	Password    string
	PhoneNumber *string
	IPAddress   string
	UserAgent   string
	DeviceID    string
	DeviceType  string
	OS          string
	Location    *model.GeoLocation
	Region      string
}

type AuthenticationRequest struct {
	Email      string
	Password   string
	IPAddress  string
	UserAgent  string
	DeviceID   string
	DeviceType string
	Location   *model.GeoLocation
}

type AuthenticationResponse struct {
	Success                bool
	User                   *model.MetaUser
	SessionToken           string
	RequiresTwoFactor      bool
	RequiresAdditionalAuth bool
	AnomalyDetected        bool
	AnomalyScore           float64
	UserID                 uuid.UUID
}
