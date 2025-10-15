package ml

import (
	"context"
	"fmt"
	"math"
	"time"

	"socialink/meta-user-service/internal/model"

	"github.com/google/uuid"
)

// BehaviorAnalyzer uses ML to detect anomalous user behavior
type BehaviorAnalyzer struct {
	models map[string]*AnomalyModel
}

func NewBehaviorAnalyzer() *BehaviorAnalyzer {
	return &BehaviorAnalyzer{
		models: make(map[string]*AnomalyModel),
	}
}

type LoginAnalysisRequest struct {
	UserID         uuid.UUID
	IPAddress      string
	UserAgent      string
	DeviceID       string
	Location       *model.GeoLocation
	Time           time.Time
	PreviousLogins []LoginRecord
}

type LoginRecord struct {
	IPAddress string
	Location  *model.GeoLocation
	Timestamp time.Time
	DeviceID  string
}

type UserActivity struct {
	Type      string
	Timestamp time.Time
	Metadata  map[string]interface{}
}

type AnomalyDetectionRequest struct {
	UserID           uuid.UUID
	Activity         *UserActivity
	BehaviorBaseline map[string]interface{}
	SensitivityLevel model.SensitivityLevel
}

// AnalyzeLogin detects anomalous login patterns
func (b *BehaviorAnalyzer) AnalyzeLogin(ctx context.Context, req *LoginAnalysisRequest) (float64, error) {
	features := b.extractLoginFeatures(req)
	
	// Multi-dimensional anomaly detection
	locationScore := b.analyzeLocationAnomaly(req)
	deviceScore := b.analyzeDeviceAnomaly(req)
	temporalScore := b.analyzeTemporalAnomaly(req)
	velocityScore := b.analyzeImpossibleTravel(req)

	// Weighted ensemble
	weights := []float64{0.3, 0.25, 0.25, 0.2}
	scores := []float64{locationScore, deviceScore, temporalScore, velocityScore}

	anomalyScore := 0.0
	for i, score := range scores {
		anomalyScore += score * weights[i]
	}

	return math.Min(anomalyScore, 1.0), nil
}

// DetectAnomaly detects anomalies in user activity
func (b *BehaviorAnalyzer) DetectAnomaly(ctx context.Context, req *AnomalyDetectionRequest) (float64, string, error) {
	// Calculate deviation from baseline behavior
	baselineScore := b.calculateBaselineDeviation(req.Activity, req.BehaviorBaseline)
	
	// Pattern-based detection
	patternScore, anomalyType := b.detectSuspiciousPatterns(req.Activity)

	// Adjust for sensitivity level
	sensitivity := b.getSensitivityMultiplier(req.SensitivityLevel)
	
	anomalyScore := (baselineScore + patternScore) / 2.0 * sensitivity

	if anomalyType == "" {
		anomalyType = "behavioral_deviation"
	}

	return math.Min(anomalyScore, 1.0), anomalyType, nil
}

// analyzeLocationAnomaly detects unusual login locations
func (b *BehaviorAnalyzer) analyzeLocationAnomaly(req *LoginAnalysisRequest) float64 {
	if req.Location == nil || len(req.PreviousLogins) == 0 {
		return 0.0
	}

	score := 0.0

	// Check if location is significantly different from previous logins
	commonLocations := b.extractCommonLocations(req.PreviousLogins)
	
	isKnownLocation := false
	for _, loc := range commonLocations {
		distance := b.calculateDistance(req.Location, loc)
		if distance < 100 { // Within 100 km
			isKnownLocation = true
			break
		}
	}

	if !isKnownLocation {
		score += 0.7
	}

	// Check for rapid location changes
	if len(req.PreviousLogins) > 0 {
		lastLogin := req.PreviousLogins[0]
		if lastLogin.Location != nil {
			distance := b.calculateDistance(req.Location, lastLogin.Location)
			timeDiff := req.Time.Sub(lastLogin.Timestamp).Hours()
			
			// Impossible travel detection (> 800 km/h)
			if timeDiff > 0 && distance/timeDiff > 800 {
				score += 0.9
			}
		}
	}

	return math.Min(score, 1.0)
}

// analyzeDeviceAnomaly detects new or suspicious devices
func (b *BehaviorAnalyzer) analyzeDeviceAnomaly(req *LoginAnalysisRequest) float64 {
	score := 0.0

	// Check if device is known
	knownDevice := false
	for _, login := range req.PreviousLogins {
		if login.DeviceID == req.DeviceID {
			knownDevice = true
			break
		}
	}

	if !knownDevice {
		score += 0.6
	}

	// Analyze user agent for suspicious patterns
	if b.isSuspiciousUserAgent(req.UserAgent) {
		score += 0.4
	}

	return math.Min(score, 1.0)
}

// analyzeTemporalAnomaly detects unusual login times
func (b *BehaviorAnalyzer) analyzeTemporalAnomaly(req *LoginAnalysisRequest) float64 {
	hour := req.Time.Hour()
	dayOfWeek := req.Time.Weekday()

	// Extract typical login patterns
	typicalHours := b.extractTypicalLoginHours(req.PreviousLogins)
	typicalDays := b.extractTypicalLoginDays(req.PreviousLogins)

	score := 0.0

	// Check if current hour is unusual
	if !contains(typicalHours, hour) {
		score += 0.3
	}

	// Very late night/early morning logins are suspicious
	if hour >= 2 && hour <= 5 {
		score += 0.4
	}

	// Check if current day is unusual
	if !containsDay(typicalDays, dayOfWeek) {
		score += 0.2
	}

	return math.Min(score, 1.0)
}

// analyzeImpossibleTravel detects impossible travel scenarios
func (b *BehaviorAnalyzer) analyzeImpossibleTravel(req *LoginAnalysisRequest) float64 {
	if req.Location == nil || len(req.PreviousLogins) == 0 {
		return 0.0
	}

	lastLogin := req.PreviousLogins[0]
	if lastLogin.Location == nil {
		return 0.0
	}

	distance := b.calculateDistance(req.Location, lastLogin.Location)
	timeDiff := req.Time.Sub(lastLogin.Timestamp).Hours()

	if timeDiff == 0 {
		return 0.0
	}

	// Calculate required speed (km/h)
	requiredSpeed := distance / timeDiff

	// Impossible by commercial flight (> 900 km/h)
	if requiredSpeed > 900 {
		return 1.0
	}

	// Suspicious but possible (> 600 km/h, requires flight)
	if requiredSpeed > 600 {
		return 0.7
	}

	// Requires driving/train but very fast
	if requiredSpeed > 200 {
		return 0.4
	}

	return 0.0
}

// calculateBaselineDeviation measures how much activity deviates from baseline
func (b *BehaviorAnalyzer) calculateBaselineDeviation(activity *UserActivity, baseline map[string]interface{}) float64 {
	// Simplified baseline comparison
	// In production, this would use statistical models and ML
	
	score := 0.0
	
	// Check activity frequency
	if avgFreq, ok := baseline["activity_frequency"].(float64); ok {
		// Compare with recent activity rate
		score += math.Abs(avgFreq - 1.0) * 0.3
	}

	// Check activity patterns
	if patterns, ok := baseline["patterns"].([]string); ok {
		matchFound := false
		for _, pattern := range patterns {
			if activity.Type == pattern {
				matchFound = true
				break
			}
		}
		if !matchFound {
			score += 0.4
		}
	}

	return math.Min(score, 1.0)
}

// detectSuspiciousPatterns identifies known suspicious patterns
func (b *BehaviorAnalyzer) detectSuspiciousPatterns(activity *UserActivity) (float64, string) {
	score := 0.0
	anomalyType := ""

	// Rapid successive actions
	if activity.Type == "rapid_actions" {
		score = 0.8
		anomalyType = "rapid_action_pattern"
	}

	// Mass follow/unfollow
	if activity.Type == "mass_follow" {
		score = 0.9
		anomalyType = "mass_action"
	}

	// Unusual access patterns
	if activity.Type == "unusual_access" {
		score = 0.7
		anomalyType = "unusual_access_pattern"
	}

	// Data scraping indicators
	if activity.Type == "high_frequency_reads" {
		score = 0.85
		anomalyType = "potential_scraping"
	}

	return score, anomalyType
}

// Helper methods

func (b *BehaviorAnalyzer) extractLoginFeatures(req *LoginAnalysisRequest) map[string]float64 {
	features := make(map[string]float64)

	features["hour"] = float64(req.Time.Hour())
	features["day_of_week"] = float64(req.Time.Weekday())
	features["previous_login_count"] = float64(len(req.PreviousLogins))

	if req.Location != nil {
		features["latitude"] = req.Location.Latitude
		features["longitude"] = req.Location.Longitude
	}

	return features
}

func (b *BehaviorAnalyzer) extractCommonLocations(logins []LoginRecord) []*model.GeoLocation {
	// Extract and cluster common locations
	var locations []*model.GeoLocation
	
	for _, login := range logins {
		if login.Location != nil {
			locations = append(locations, login.Location)
		}
	}

	// Simplified - would use clustering algorithm (DBSCAN) in production
	return locations
}

func (b *BehaviorAnalyzer) calculateDistance(loc1, loc2 *model.GeoLocation) float64 {
	// Haversine formula for great circle distance
	const earthRadius = 6371.0 // km

	lat1 := loc1.Latitude * math.Pi / 180
	lat2 := loc2.Latitude * math.Pi / 180
	deltaLat := (loc2.Latitude - loc1.Latitude) * math.Pi / 180
	deltaLon := (loc2.Longitude - loc1.Longitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (b *BehaviorAnalyzer) extractTypicalLoginHours(logins []LoginRecord) []int {
	hourFreq := make(map[int]int)
	for _, login := range logins {
		hourFreq[login.Timestamp.Hour()]++
	}

	// Return hours that appear in > 20% of logins
	threshold := len(logins) / 5
	var typical []int
	for hour, count := range hourFreq {
		if count > threshold {
			typical = append(typical, hour)
		}
	}

	return typical
}

func (b *BehaviorAnalyzer) extractTypicalLoginDays(logins []LoginRecord) []time.Weekday {
	dayFreq := make(map[time.Weekday]int)
	for _, login := range logins {
		dayFreq[login.Timestamp.Weekday()]++
	}

	threshold := len(logins) / 7
	var typical []time.Weekday
	for day, count := range dayFreq {
		if count > threshold {
			typical = append(typical, day)
		}
	}

	return typical
}

func (b *BehaviorAnalyzer) isSuspiciousUserAgent(ua string) bool {
	suspicious := []string{"bot", "scraper", "crawler", "automated"}
	for _, s := range suspicious {
		if contains(ua, s) {
			return true
		}
	}
	return false
}

func (b *BehaviorAnalyzer) getSensitivityMultiplier(level model.SensitivityLevel) float64 {
	switch level {
	case model.SensitivityLow:
		return 0.7
	case model.SensitivityMedium:
		return 1.0
	case model.SensitivityHigh:
		return 1.3
	default:
		return 1.0
	}
}

func contains(s string, substr string) bool {
	return len(s) >= len(substr)
	// Simplified - use strings.Contains in production
}

func containsDay(days []time.Weekday, day time.Weekday) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

// AnomalyModel represents a trained anomaly detection model
type AnomalyModel struct {
	ModelID   string
	Version   string
	Threshold float64
	Features  []string
}
