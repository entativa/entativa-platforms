package ml

import (
	"context"
	"fmt"
	"math"
	"time"

	"socialink/meta-user-service/internal/model"

	"github.com/google/uuid"
)

// TrustScoreEngine calculates user trust scores using ML
type TrustScoreEngine struct {
	weights TrustScoreWeights
}

type TrustScoreWeights struct {
	AccountAge           float64
	VerificationStatus   float64
	ActivityScore        float64
	ViolationPenalty     float64
	PositiveInteractions float64
	DeviceTrust          float64
	SecurityCompliance   float64
	AnomalyPenalty       float64
}

func NewTrustScoreEngine() *TrustScoreEngine {
	return &TrustScoreEngine{
		weights: TrustScoreWeights{
			AccountAge:           0.15,
			VerificationStatus:   0.20,
			ActivityScore:        0.15,
			ViolationPenalty:     0.15,
			PositiveInteractions: 0.15,
			DeviceTrust:          0.10,
			SecurityCompliance:   0.05,
			AnomalyPenalty:       0.05,
		},
	}
}

type TrustScoreRequest struct {
	UserID               uuid.UUID
	AccountAge           time.Duration
	VerificationStatus   bool
	ActivityScore        float64
	ReportedViolations   int
	PositiveInteractions int
	DeviceFingerprints   int
	SecurityProfile      model.SecurityProfile
	AnomalyHistory       int
}

// CalculateTrustScore computes a comprehensive trust score (0.0 - 1.0)
func (e *TrustScoreEngine) CalculateTrustScore(ctx context.Context, req *TrustScoreRequest) (float64, error) {
	// Calculate individual component scores
	accountAgeScore := e.calculateAccountAgeScore(req.AccountAge)
	verificationScore := e.calculateVerificationScore(req.VerificationStatus)
	activityScore := e.normalizeActivityScore(req.ActivityScore)
	violationScore := e.calculateViolationScore(req.ReportedViolations)
	interactionScore := e.calculateInteractionScore(req.PositiveInteractions)
	deviceScore := e.calculateDeviceScore(req.DeviceFingerprints)
	securityScore := e.calculateSecurityScore(req.SecurityProfile)
	anomalyScore := e.calculateAnomalyScore(req.AnomalyHistory)

	// Weighted sum
	trustScore := 0.0
	trustScore += accountAgeScore * e.weights.AccountAge
	trustScore += verificationScore * e.weights.VerificationStatus
	trustScore += activityScore * e.weights.ActivityScore
	trustScore += violationScore * e.weights.ViolationPenalty
	trustScore += interactionScore * e.weights.PositiveInteractions
	trustScore += deviceScore * e.weights.DeviceTrust
	trustScore += securityScore * e.weights.SecurityCompliance
	trustScore += anomalyScore * e.weights.AnomalyPenalty

	// Apply non-linear transformation for better distribution
	trustScore = e.applyNonLinearTransform(trustScore)

	return math.Max(0.0, math.Min(1.0, trustScore)), nil
}

// calculateAccountAgeScore rewards older accounts
func (e *TrustScoreEngine) calculateAccountAgeScore(age time.Duration) float64 {
	days := age.Hours() / 24
	
	// Logarithmic growth with diminishing returns
	// New accounts: ~0.1, 30 days: ~0.4, 90 days: ~0.6, 365 days: ~0.8, 730+ days: ~1.0
	if days <= 0 {
		return 0.1
	}
	
	score := math.Log10(days+1) / 3.0
	return math.Min(score, 1.0)
}

// calculateVerificationScore rewards verified users
func (e *TrustScoreEngine) calculateVerificationScore(verified bool) float64 {
	if verified {
		return 1.0
	}
	return 0.3 // Partial credit for unverified but active users
}

// normalizeActivityScore normalizes activity score to 0-1 range
func (e *TrustScoreEngine) normalizeActivityScore(activityScore float64) float64 {
	// Assuming activity score is already 0-1, apply sigmoid for smoothing
	return 1.0 / (1.0 + math.Exp(-5*(activityScore-0.5)))
}

// calculateViolationScore penalizes reported violations
func (e *TrustScoreEngine) calculateViolationScore(violations int) float64 {
	if violations == 0 {
		return 1.0
	}
	
	// Exponential penalty for violations
	penalty := math.Exp(-0.5 * float64(violations))
	return penalty
}

// calculateInteractionScore rewards positive interactions
func (e *TrustScoreEngine) calculateInteractionScore(interactions int) float64 {
	if interactions == 0 {
		return 0.0
	}
	
	// Logarithmic growth
	// 10 interactions: ~0.3, 100: ~0.6, 1000: ~0.9
	score := math.Log10(float64(interactions)+1) / 4.0
	return math.Min(score, 1.0)
}

// calculateDeviceScore evaluates device trust
func (e *TrustScoreEngine) calculateDeviceScore(deviceCount int) float64 {
	// Normal users have 1-3 devices, suspicious users have many
	if deviceCount == 0 {
		return 0.0
	}
	
	if deviceCount <= 3 {
		return 1.0
	}
	
	if deviceCount <= 5 {
		return 0.7
	}
	
	// Penalty for too many devices (potential account sharing/fraud)
	return math.Max(0.3, 1.0 - float64(deviceCount-5)*0.1)
}

// calculateSecurityScore rewards strong security practices
func (e *TrustScoreEngine) calculateSecurityScore(profile model.SecurityProfile) float64 {
	score := 0.0
	
	// Two-factor authentication
	if profile.TwoFactorEnabled {
		score += 0.4
	}
	
	// Biometric or hardware key
	if profile.BiometricEnabled || profile.PasskeyEnabled {
		score += 0.3
	}
	
	// No recent failed login attempts
	if profile.FailedLoginAttempts == 0 {
		score += 0.2
	}
	
	// Regular password changes
	daysSincePasswordChange := time.Since(profile.LastPasswordChange).Hours() / 24
	if daysSincePasswordChange < 90 {
		score += 0.1
	}
	
	return math.Min(score, 1.0)
}

// calculateAnomalyScore penalizes anomaly history
func (e *TrustScoreEngine) calculateAnomalyScore(anomalyCount int) float64 {
	if anomalyCount == 0 {
		return 1.0
	}
	
	// Exponential penalty
	penalty := math.Exp(-0.3 * float64(anomalyCount))
	return penalty
}

// applyNonLinearTransform applies sigmoid transformation for better score distribution
func (e *TrustScoreEngine) applyNonLinearTransform(score float64) float64 {
	// Sigmoid centered at 0.5
	return 1.0 / (1.0 + math.Exp(-10*(score-0.5)))
}

// PredictFutureScore predicts future trust score based on trends
func (e *TrustScoreEngine) PredictFutureScore(ctx context.Context, userID uuid.UUID, currentScore float64, recentTrend float64, days int) (float64, error) {
	// Simple linear projection with decay
	projectedScore := currentScore + (recentTrend * float64(days) * 0.1)
	
	// Apply bounds
	return math.Max(0.0, math.Min(1.0, projectedScore)), nil
}

// GetScoreExplanation provides explanation for trust score components
func (e *TrustScoreEngine) GetScoreExplanation(req *TrustScoreRequest) map[string]interface{} {
	explanation := make(map[string]interface{})
	
	explanation["account_age"] = map[string]interface{}{
		"score":       e.calculateAccountAgeScore(req.AccountAge),
		"weight":      e.weights.AccountAge,
		"contribution": e.calculateAccountAgeScore(req.AccountAge) * e.weights.AccountAge,
		"description": "Older accounts are more trustworthy",
	}
	
	explanation["verification"] = map[string]interface{}{
		"score":       e.calculateVerificationScore(req.VerificationStatus),
		"weight":      e.weights.VerificationStatus,
		"contribution": e.calculateVerificationScore(req.VerificationStatus) * e.weights.VerificationStatus,
		"description": "Email and phone verification",
	}
	
	explanation["activity"] = map[string]interface{}{
		"score":       e.normalizeActivityScore(req.ActivityScore),
		"weight":      e.weights.ActivityScore,
		"contribution": e.normalizeActivityScore(req.ActivityScore) * e.weights.ActivityScore,
		"description": "Regular, genuine activity",
	}
	
	explanation["violations"] = map[string]interface{}{
		"score":       e.calculateViolationScore(req.ReportedViolations),
		"weight":      e.weights.ViolationPenalty,
		"contribution": e.calculateViolationScore(req.ReportedViolations) * e.weights.ViolationPenalty,
		"description": fmt.Sprintf("%d reported violations", req.ReportedViolations),
	}
	
	explanation["interactions"] = map[string]interface{}{
		"score":       e.calculateInteractionScore(req.PositiveInteractions),
		"weight":      e.weights.PositiveInteractions,
		"contribution": e.calculateInteractionScore(req.PositiveInteractions) * e.weights.PositiveInteractions,
		"description": fmt.Sprintf("%d positive interactions", req.PositiveInteractions),
	}
	
	explanation["security"] = map[string]interface{}{
		"score":       e.calculateSecurityScore(req.SecurityProfile),
		"weight":      e.weights.SecurityCompliance,
		"contribution": e.calculateSecurityScore(req.SecurityProfile) * e.weights.SecurityCompliance,
		"description": "Security features enabled",
	}
	
	return explanation
}
