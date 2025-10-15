package util

import (
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

// GenerateUsername generates a username from first and last name
func GenerateUsername(firstName, lastName string) string {
	// Create base username
	base := strings.ToLower(firstName + "." + lastName)
	// Remove any special characters
	base = regexp.MustCompile(`[^a-z0-9._]`).ReplaceAllString(base, "")
	// Add timestamp suffix to ensure uniqueness
	suffix := time.Now().Unix() % 10000
	return base + string(rune('0'+suffix/1000)) + string(rune('0'+(suffix/100)%10)) + string(rune('0'+(suffix/10)%10)) + string(rune('0'+suffix%10))
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
