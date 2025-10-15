package util

import (
	"regexp"
	"strings"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._]{3,30}$`)
)

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidUsername checks if a username is valid (Instagram-style)
// Username must be 3-30 characters, can contain letters, numbers, periods, and underscores
func IsValidUsername(username string) bool {
	// Check length
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	
	// Check format
	if !usernameRegex.MatchString(username) {
		return false
	}
	
	// Username cannot start or end with a period
	if strings.HasPrefix(username, ".") || strings.HasSuffix(username, ".") {
		return false
	}
	
	// Username cannot have consecutive periods
	if strings.Contains(username, "..") {
		return false
	}
	
	return true
}

// SanitizeUsername sanitizes a username by removing invalid characters
func SanitizeUsername(username string) string {
	// Convert to lowercase
	username = strings.ToLower(username)
	
	// Remove spaces and special characters except . and _
	re := regexp.MustCompile(`[^a-z0-9._]`)
	username = re.ReplaceAllString(username, "")
	
	// Remove leading/trailing periods
	username = strings.Trim(username, ".")
	
	// Replace consecutive periods
	for strings.Contains(username, "..") {
		username = strings.ReplaceAll(username, "..", ".")
	}
	
	return username
}

// ValidatePassword checks password strength
func ValidatePassword(password string) bool {
	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}
	
	// Maximum 128 characters
	if len(password) > 128 {
		return false
	}
	
	return true
}

// ValidateBio checks if bio is valid (max 150 characters for Instagram-like)
func ValidateBio(bio string) bool {
	return len(bio) <= 150
}
