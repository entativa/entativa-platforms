package ml

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/stat"
)

// FraudDetector uses ML models to detect fraudulent signup attempts
type FraudDetector struct {
	threshold       float64
	emailBlacklist  map[string]bool
	ipRiskScores    map[string]float64
	deviceRiskCache map[string]float64
}

func NewFraudDetector() *FraudDetector {
	return &FraudDetector{
		threshold:       0.7,
		emailBlacklist:  make(map[string]bool),
		ipRiskScores:    make(map[string]float64),
		deviceRiskCache: make(map[string]float64),
	}
}

type SignupAnalysisRequest struct {
	Email       string
	PhoneNumber *string
	IPAddress   string
	UserAgent   string
	DeviceID    string
	Timestamp   time.Time
}

// AnalyzeSignup performs multi-factor fraud analysis on signup
func (f *FraudDetector) AnalyzeSignup(ctx context.Context, req *SignupAnalysisRequest) (float64, error) {
	features := f.extractSignupFeatures(req)
	
	// Multi-factor scoring
	emailScore := f.analyzeEmail(req.Email)
	ipScore := f.analyzeIP(req.IPAddress)
	deviceScore := f.analyzeDevice(req.DeviceID, req.UserAgent)
	velocityScore := f.analyzeSignupVelocity(req.IPAddress, req.DeviceID)
	patternScore := f.analyzePatterns(features)

	// Weighted ensemble scoring
	weights := []float64{0.25, 0.20, 0.20, 0.20, 0.15}
	scores := []float64{emailScore, ipScore, deviceScore, velocityScore, patternScore}
	
	fraudScore := stat.Mean(scores, weights)

	return fraudScore, nil
}

// analyzeEmail checks email for suspicious patterns
func (f *FraudDetector) analyzeEmail(email string) float64 {
	score := 0.0

	// Check blacklist
	if f.emailBlacklist[email] {
		return 1.0
	}

	// Disposable email detection
	disposableDomains := []string{"tempmail.com", "guerrillamail.com", "10minutemail.com"}
	for _, domain := range disposableDomains {
		if contains(email, domain) {
			score += 0.7
		}
	}

	// Pattern analysis: random character sequences
	if f.hasRandomPattern(email) {
		score += 0.3
	}

	// Numeric-heavy emails are suspicious
	numCount := 0
	for _, ch := range email {
		if ch >= '0' && ch <= '9' {
			numCount++
		}
	}
	if float64(numCount)/float64(len(email)) > 0.5 {
		score += 0.2
	}

	return math.Min(score, 1.0)
}

// analyzeIP evaluates IP address reputation
func (f *FraudDetector) analyzeIP(ip string) float64 {
	// Check cache
	if cachedScore, ok := f.ipRiskScores[ip]; ok {
		return cachedScore
	}

	score := 0.0

	// VPN/Proxy detection (simplified)
	if f.isVPNIP(ip) {
		score += 0.4
	}

	// Geographic risk assessment
	if f.isHighRiskRegion(ip) {
		score += 0.3
	}

	// Rate limiting check
	if f.isRateLimited(ip) {
		score += 0.5
	}

	// Cache result
	f.ipRiskScores[ip] = score

	return math.Min(score, 1.0)
}

// analyzeDevice examines device fingerprint
func (f *FraudDetector) analyzeDevice(deviceID, userAgent string) float64 {
	// Check cache
	if cachedScore, ok := f.deviceRiskCache[deviceID]; ok {
		return cachedScore
	}

	score := 0.0

	// Unusual user agent
	if f.isUnusualUserAgent(userAgent) {
		score += 0.3
	}

	// Device seen in multiple rapid signups
	if f.isReusedDevice(deviceID) {
		score += 0.5
	}

	// Automation detection
	if f.detectsAutomation(userAgent) {
		score += 0.6
	}

	f.deviceRiskCache[deviceID] = score

	return math.Min(score, 1.0)
}

// analyzeSignupVelocity detects velocity attacks
func (f *FraudDetector) analyzeSignupVelocity(ip, deviceID string) float64 {
	score := 0.0

	// Signups per hour from same IP
	ipVelocity := f.getIPSignupVelocity(ip)
	if ipVelocity > 10 {
		score += 0.7
	} else if ipVelocity > 5 {
		score += 0.4
	}

	// Signups from same device
	deviceVelocity := f.getDeviceSignupVelocity(deviceID)
	if deviceVelocity > 5 {
		score += 0.6
	}

	return math.Min(score, 1.0)
}

// analyzePatterns uses pattern recognition for fraud detection
func (f *FraudDetector) analyzePatterns(features map[string]float64) float64 {
	// Simplified pattern analysis using feature vector
	// In production, this would use a trained ML model (Random Forest, Neural Network)
	
	suspiciousPatterns := 0
	totalPatterns := 0

	for key, value := range features {
		totalPatterns++
		
		// Example pattern checks
		if key == "email_entropy" && value < 2.0 {
			suspiciousPatterns++
		}
		if key == "signup_hour" && (value < 6 || value > 23) {
			suspiciousPatterns++ // Unusual hours
		}
	}

	if totalPatterns == 0 {
		return 0.0
	}

	return float64(suspiciousPatterns) / float64(totalPatterns)
}

// extractSignupFeatures extracts feature vector for ML analysis
func (f *FraudDetector) extractSignupFeatures(req *SignupAnalysisRequest) map[string]float64 {
	features := make(map[string]float64)

	// Temporal features
	features["signup_hour"] = float64(req.Timestamp.Hour())
	features["signup_day_of_week"] = float64(req.Timestamp.Weekday())

	// Email features
	features["email_length"] = float64(len(req.Email))
	features["email_entropy"] = f.calculateEntropy(req.Email)
	features["email_has_numbers"] = boolToFloat(f.hasNumbers(req.Email))

	// Device features
	features["user_agent_length"] = float64(len(req.UserAgent))

	return features
}

// Helper methods

func (f *FraudDetector) hasRandomPattern(s string) bool {
	// Simplified random pattern detection
	entropy := f.calculateEntropy(s)
	return entropy > 3.5 // High entropy suggests random
}

func (f *FraudDetector) calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}

	entropy := 0.0
	length := float64(len(s))
	for _, count := range freq {
		p := float64(count) / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

func (f *FraudDetector) isVPNIP(ip string) bool {
	// Simplified VPN detection - would use external API in production
	vpnRanges := []string{"10.", "192.168."}
	for _, prefix := range vpnRanges {
		if len(ip) >= len(prefix) && ip[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func (f *FraudDetector) isHighRiskRegion(ip string) bool {
	// Would use GeoIP database in production
	return false
}

func (f *FraudDetector) isRateLimited(ip string) bool {
	// Would check against rate limiting service
	return false
}

func (f *FraudDetector) isUnusualUserAgent(ua string) bool {
	// Check for bot-like user agents
	botIndicators := []string{"bot", "crawler", "spider", "scraper"}
	for _, indicator := range botIndicators {
		if contains(ua, indicator) {
			return true
		}
	}
	return false
}

func (f *FraudDetector) isReusedDevice(deviceID string) bool {
	// Would check database for device reuse
	return false
}

func (f *FraudDetector) detectsAutomation(ua string) bool {
	// Detect automation tools
	automationTools := []string{"selenium", "puppeteer", "playwright", "phantom"}
	for _, tool := range automationTools {
		if contains(ua, tool) {
			return true
		}
	}
	return false
}

func (f *FraudDetector) getIPSignupVelocity(ip string) int {
	// Would query database for recent signups from IP
	return 0
}

func (f *FraudDetector) getDeviceSignupVelocity(deviceID string) int {
	// Would query database for recent signups from device
	return 0
}

func (f *FraudDetector) hasNumbers(s string) bool {
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			return true
		}
	}
	return false
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != substr
	// Simplified - use strings.Contains in production
}
