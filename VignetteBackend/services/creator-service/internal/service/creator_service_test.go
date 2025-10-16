package service

import (
	"testing"

	"vignette/creator-service/internal/model"
)

func TestMonetizationRequirements(t *testing.T) {
	tests := []struct {
		name        string
		followers   int
		posts       int
		shouldMeet  bool
	}{
		{"Meets requirements", 15000, 150, true},
		{"Just meets", 10000, 100, true},
		{"Low followers", 5000, 150, false},
		{"Low posts", 15000, 50, false},
		{"Both low", 5000, 50, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meets := tt.followers >= MinFollowersForMonetization && tt.posts >= MinPostsForMonetization
			if meets != tt.shouldMeet {
				t.Errorf("Got %v, want %v for followers=%d, posts=%d", meets, tt.shouldMeet, tt.followers, tt.posts)
			}
		})
	}
}

func TestEngagementRateCalculation(t *testing.T) {
	tests := []struct {
		name     string
		likes    int
		comments int
		shares   int
		saves    int
		reach    int
		expected float64
	}{
		{"High engagement", 100, 50, 20, 30, 1000, 20.0},
		{"Low engagement", 10, 5, 2, 3, 1000, 2.0},
		{"Zero reach", 100, 50, 20, 30, 0, 0.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engagement := tt.likes + tt.comments + tt.shares + tt.saves
			var rate float64
			if tt.reach > 0 {
				rate = float64(engagement) / float64(tt.reach) * 100
			}
			
			if rate != tt.expected {
				t.Errorf("Got %v, want %v", rate, tt.expected)
			}
		})
	}
}

func TestAccountTypeValidation(t *testing.T) {
	validTypes := []model.AccountType{
		model.AccountTypePersonal,
		model.AccountTypeBusiness,
		model.AccountTypeCreator,
	}
	
	for _, accountType := range validTypes {
		if accountType == "" {
			t.Errorf("Account type should not be empty")
		}
	}
}
