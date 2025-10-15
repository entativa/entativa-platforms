package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._]{3,30}$`)
)

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidUsername checks if a username is valid
func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

// GenerateUsername generates a clean username from first and last name
// Format: firstname.lastname (for URLs like Socialink.com/john.doe)
func GenerateUsername(firstName, lastName string) string {
	// Create base username: firstname.lastname
	first := strings.ToLower(firstName)
	last := strings.ToLower(lastName)
	
	// Remove any special characters and spaces
	first = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(first, "")
	last = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(last, "")
	
	// Base format: firstname.lastname
	base := first + "." + last
	
	return base
}

// GenerateUniqueUsername generates a unique username with suffix if needed
func GenerateUniqueUsername(firstName, lastName string) string {
	base := GenerateUsername(firstName, lastName)
	
	// Add random suffix to ensure uniqueness
	// This will be checked against database and regenerated if exists
	suffix := time.Now().UnixNano() % 9999
	if suffix == 0 {
		return base
	}
	return fmt.Sprintf("%s%d", base, suffix)
}

// ValidateDisplayName validates first/last names (relaxed policy)
// Unlike Facebook, we don't require legal names but recommend them
func ValidateDisplayName(name string) (bool, string) {
	// Remove leading/trailing whitespace
	name = strings.TrimSpace(name)
	
	// Minimum length check
	if len(name) < 1 {
		return false, "Name must be at least 1 character"
	}
	
	// Maximum length check
	if len(name) > 50 {
		return false, "Name must be less than 50 characters"
	}
	
	// Allow letters, spaces, hyphens, apostrophes, and common international characters
	// This is more relaxed than Facebook's policy
	validPattern := regexp.MustCompile(`^[\p{L}\p{M}\s\-'\.]+$`)
	if !validPattern.MatchString(name) {
		return false, "Name contains invalid characters. Use letters, spaces, hyphens, or apostrophes"
	}
	
	// Check for excessive special characters (anti-spam)
	specialChars := regexp.MustCompile(`[\-'\.]`)
	if len(specialChars.FindAllString(name, -1)) > len(name)/2 {
		return false, "Name contains too many special characters"
	}
	
	return true, ""
}

// IsLikelyRealName provides a hint if the name looks like a real name
// This is informational only, not enforced
func IsLikelyRealName(firstName, lastName string) (bool, string) {
	// Check for common patterns that suggest fake names
	full := strings.ToLower(firstName + " " + lastName)
	
	// Check for numbers
	if regexp.MustCompile(`\d`).MatchString(full) {
		return false, "Consider using your real name for better connections"
	}
	
	// Check for excessive repetition
	if regexp.MustCompile(`(.)\1{3,}`).MatchString(full) {
		return false, "Consider using your real name for authenticity"
	}
	
	// Check for common fake patterns
	fakePatterns := []string{"test", "fake", "admin", "user", "account", "temp"}
	for _, pattern := range fakePatterns {
		if strings.Contains(full, pattern) {
			return false, "We recommend using your real name to connect with friends and family"
		}
	}
	
	return true, ""
}

// IsValidBirthday checks if a birthday is valid (user must be at least 13 years old)
func IsValidBirthday(birthday time.Time) bool {
	age := time.Now().Year() - birthday.Year()
	if time.Now().YearDay() < birthday.YearDay() {
		age--
	}
	return age >= 13 && age <= 120
}

// ValidatePassword checks password strength
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	// Could add more complex validation here
	return true
}
