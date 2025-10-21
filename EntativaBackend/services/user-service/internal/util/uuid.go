package util

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// ParseUUID parses a UUID string
func ParseUUID(u string) (uuid.UUID, error) {
	return uuid.Parse(u)
}
