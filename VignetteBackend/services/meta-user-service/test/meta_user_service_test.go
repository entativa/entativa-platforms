package test

import (
	"context"
	"testing"
	"time"

	"vignette/meta-user-service/internal/model"
	"vignette/meta-user-service/internal/service"
	"vignette/meta-user-service/pkg/ml"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockMetaUserRepository struct {
	mock.Mock
}

func (m *MockMetaUserRepository) Create(ctx context.Context, user *model.MetaUser) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockMetaUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MetaUser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MetaUser), args.Error(1)
}

func (m *MockMetaUserRepository) GetByEmail(ctx context.Context, email string) (*model.MetaUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MetaUser), args.Error(1)
}

func (m *MockMetaUserRepository) Update(ctx context.Context, user *model.MetaUser) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Mock Event Publisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) PublishUserEvent(ctx context.Context, eventType string, data map[string]interface{}) error {
	args := m.Called(ctx, eventType, data)
	return args.Error(0)
}

func TestCreateMetaUser_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockMetaUserRepository)
	mockPublisher := new(MockEventPublisher)
	fraudDetector := ml.NewFraudDetector()
	behaviorAnalyzer := ml.NewBehaviorAnalyzer()
	trustScoreEngine := ml.NewTrustScoreEngine()

	metaUserService := service.NewMetaUserService(
		mockRepo,
		fraudDetector,
		behaviorAnalyzer,
		trustScoreEngine,
		mockPublisher,
		nil,
	)

	// Mock expectations
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.MetaUser")).Return(nil)
	mockPublisher.On("PublishUserEvent", mock.Anything, "meta.user.created", mock.Anything).Return(nil)

	// Test
	ctx := context.Background()
	req := &service.CreateMetaUserRequest{
		Email:      "test@example.com",
		Password:   "securepassword123",
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		DeviceID:   "device-123",
		DeviceType: "desktop",
		OS:         "Windows",
		Region:     "US",
	}

	user, err := metaUserService.CreateMetaUser(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.MetaID)
	assert.Equal(t, model.UserStatusPending, user.Status)
	assert.GreaterOrEqual(t, user.TrustScore, 0.0)
	assert.LessOrEqual(t, user.TrustScore, 1.0)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestFraudDetection_HighRiskEmail(t *testing.T) {
	// Setup
	fraudDetector := ml.NewFraudDetector()

	// Test with disposable email
	ctx := context.Background()
	req := &ml.SignupAnalysisRequest{
		Email:       "test@tempmail.com",
		IPAddress:   "192.168.1.1",
		UserAgent:   "Mozilla/5.0",
		DeviceID:    "device-123",
		Timestamp:   time.Now(),
	}

	score, err := fraudDetector.AnalyzeSignup(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.Greater(t, score, 0.5, "Disposable email should have high fraud score")
}

func TestTrustScoreCalculation(t *testing.T) {
	// Setup
	trustScoreEngine := ml.NewTrustScoreEngine()

	// Test with good user
	ctx := context.Background()
	req := &ml.TrustScoreRequest{
		UserID:               uuid.New(),
		AccountAge:           time.Hour * 24 * 365, // 1 year old
		VerificationStatus:   true,
		ActivityScore:        0.8,
		ReportedViolations:   0,
		PositiveInteractions: 1000,
		DeviceFingerprints:   2,
		SecurityProfile: model.SecurityProfile{
			TwoFactorEnabled:    true,
			BiometricEnabled:    true,
			FailedLoginAttempts: 0,
			LastPasswordChange:  time.Now().Add(-time.Hour * 24 * 30),
		},
		AnomalyHistory: 0,
	}

	score, err := trustScoreEngine.CalculateTrustScore(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.Greater(t, score, 0.7, "Good user should have high trust score")
	assert.LessOrEqual(t, score, 1.0)
}

func TestBehaviorAnalyzer_ImpossibleTravel(t *testing.T) {
	// Setup
	analyzer := ml.NewBehaviorAnalyzer()

	// Create scenario with impossible travel
	now := time.Now()
	previousLogin := ml.LoginRecord{
		IPAddress: "1.2.3.4",
		Location: &model.GeoLocation{
			Country:   "USA",
			City:      "New York",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		Timestamp: now.Add(-time.Hour * 1), // 1 hour ago
	}

	req := &ml.LoginAnalysisRequest{
		UserID:    uuid.New(),
		IPAddress: "5.6.7.8",
		UserAgent: "Mozilla/5.0",
		DeviceID:  "device-123",
		Location: &model.GeoLocation{
			Country:   "Japan",
			City:      "Tokyo",
			Latitude:  35.6762,
			Longitude: 139.6503,
		},
		Time:           now,
		PreviousLogins: []ml.LoginRecord{previousLogin},
	}

	score, err := analyzer.AnalyzeLogin(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Greater(t, score, 0.5, "Impossible travel should trigger high anomaly score")
}

func TestAuthenticate_ValidCredentials(t *testing.T) {
	// Setup
	mockRepo := new(MockMetaUserRepository)
	mockPublisher := new(MockEventPublisher)
	fraudDetector := ml.NewFraudDetector()
	behaviorAnalyzer := ml.NewBehaviorAnalyzer()
	trustScoreEngine := ml.NewTrustScoreEngine()

	metaUserService := service.NewMetaUserService(
		mockRepo,
		fraudDetector,
		behaviorAnalyzer,
		trustScoreEngine,
		mockPublisher,
		nil,
	)

	// Create mock user with hashed password
	// Note: In real test, use bcrypt to hash the password
	mockUser := &model.MetaUser{
		ID:           uuid.New(),
		MetaID:       "META-123",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$...", // Mock hash
		Status:       model.UserStatusActive,
		SecurityProfile: model.SecurityProfile{
			TwoFactorEnabled: false,
		},
	}

	// Mock expectations
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.MetaUser")).Return(nil)
	mockPublisher.On("PublishUserEvent", mock.Anything, "meta.user.authenticated", mock.Anything).Return(nil)

	// Note: This test would need actual password hashing to work completely
	// For demonstration purposes, we're showing the structure
}

func TestCrossPlatformSync_EnableSync(t *testing.T) {
	// Setup
	mockRepo := new(MockMetaUserRepository)
	mockPublisher := new(MockEventPublisher)

	syncService := service.NewCrossPlatformSyncService(
		mockRepo,
		mockPublisher,
		nil,
		nil,
	)

	userID := uuid.New()
	mockUser := &model.MetaUser{
		ID:     userID,
		MetaID: "META-123",
		PlatformLinks: model.PlatformLinks{
			SyncEnabled: false,
			LinkStatus:  model.LinkStatusUnlinked,
		},
	}

	// Mock expectations
	mockRepo.On("GetByID", mock.Anything, userID).Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.MetaUser")).Return(nil)
	mockPublisher.On("PublishUserEvent", mock.Anything, "meta.sync.enabled", mock.Anything).Return(nil)

	// Test
	err := syncService.EnableCrossPlatformSync(context.Background(), userID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

// Benchmark tests

func BenchmarkTrustScoreCalculation(b *testing.B) {
	engine := ml.NewTrustScoreEngine()
	ctx := context.Background()
	req := &ml.TrustScoreRequest{
		UserID:               uuid.New(),
		AccountAge:           time.Hour * 24 * 365,
		VerificationStatus:   true,
		ActivityScore:        0.8,
		ReportedViolations:   0,
		PositiveInteractions: 1000,
		DeviceFingerprints:   2,
		SecurityProfile: model.SecurityProfile{
			TwoFactorEnabled: true,
		},
		AnomalyHistory: 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.CalculateTrustScore(ctx, req)
	}
}

func BenchmarkFraudDetection(b *testing.B) {
	detector := ml.NewFraudDetector()
	ctx := context.Background()
	req := &ml.SignupAnalysisRequest{
		Email:      "test@example.com",
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		DeviceID:   "device-123",
		Timestamp:  time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.AnalyzeSignup(ctx, req)
	}
}
