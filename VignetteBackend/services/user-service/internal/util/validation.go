package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Email regex pattern
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	
	// Username regex (Instagram-style for Vignette)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
)

// IsValidEmail validates an email address
func IsValidEmail(email string) bool {
	if email == "" || len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidateEmail validates email and returns error if invalid
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	if len(email) > 255 {
		return fmt.Errorf("email is too long")
	}
	
	if !IsValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// IsValidUsername validates a username (Instagram-style rules)
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	
	// Check if matches pattern
	if !usernameRegex.MatchString(username) {
		return false
	}
	
	// Cannot start or end with period
	if strings.HasPrefix(username, ".") || strings.HasSuffix(username, ".") {
		return false
	}
	
	// Cannot have consecutive periods
	if strings.Contains(username, "..") {
		return false
	}
	
	return true
}

// ValidateUsername validates username and returns error if invalid
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	
	if username == "" {
		return fmt.Errorf("username is required")
	}
	
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	
	if len(username) > 30 {
		return fmt.Errorf("username must be 30 characters or less")
	}
	
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, periods, and underscores")
	}
	
	if strings.HasPrefix(username, ".") || strings.HasSuffix(username, ".") {
		return fmt.Errorf("username cannot start or end with a period")
	}
	
	if strings.Contains(username, "..") {
		return fmt.Errorf("username cannot have consecutive periods")
	}
	
	return nil
}

// ValidateName validates first/last name
func ValidateName(name, fieldName string) error {
	name = strings.TrimSpace(name)
	
	if name == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	
	if len(name) < 2 {
		return fmt.Errorf("%s must be at least 2 characters", fieldName)
	}
	
	if len(name) > 50 {
		return fmt.Errorf("%s must be 50 characters or less", fieldName)
	}
	
	// Check if contains only letters, spaces, hyphens, apostrophes
	for _, char := range name {
		if !unicode.IsLetter(char) && char != ' ' && char != '-' && char != '\'' {
			return fmt.Errorf("%s can only contain letters, spaces, hyphens, and apostrophes", fieldName)
		}
	}
	
	return nil
}

// SanitizeInput removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	return input
}

// ValidateGender validates gender field
func ValidateGender(gender string) error {
	validGenders := map[string]bool{
		"male":                true,
		"female":              true,
		"non_binary":          true,
		"prefer_not_to_say":   true,
		"custom":              true,
	}
	
	if !validGenders[gender] {
		return fmt.Errorf("invalid gender value")
	}
	
	return nil
}

// SplitFullName splits a full name into first and last name
func SplitFullName(fullName string) (firstName, lastName string) {
	fullName = strings.TrimSpace(fullName)
	
	parts := strings.SplitN(fullName, " ", 2)
	
	if len(parts) == 1 {
		return parts[0], ""
	}
	
	return parts[0], parts[1]
}
