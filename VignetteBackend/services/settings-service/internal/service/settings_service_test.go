package service

import (
	"testing"

	"vignette/settings-service/internal/crypto"
)

func TestEncryptionRoundtrip(t *testing.T) {
	encService := crypto.NewEncryptionService()
	
	// Test data
	originalKeys := []byte("test encryption keys data")
	passphrase := "TestPass123"
	
	// Generate salt
	salt, err := encService.GenerateSalt()
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}
	
	// Encrypt
	encrypted, err := encService.EncryptKeys(originalKeys, passphrase, salt)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	
	// Decrypt
	decrypted, err := encService.DecryptKeys(encrypted, passphrase, salt)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	
	// Verify
	if string(decrypted) != string(originalKeys) {
		t.Errorf("Decrypted data doesn't match original. Got: %s, Want: %s", string(decrypted), string(originalKeys))
	}
}

func TestPINValidation(t *testing.T) {
	encService := crypto.NewEncryptionService()
	
	tests := []struct {
		name    string
		pin     string
		wantErr bool
	}{
		{"Valid PIN", "123456", false},
		{"Too short", "12345", true},
		{"Too long", "1234567", true},
		{"Non-numeric", "12345a", true},
		{"Empty", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := encService.ValidatePIN(tt.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassphraseValidation(t *testing.T) {
	encService := crypto.NewEncryptionService()
	
	tests := []struct {
		name    string
		pass    string
		wantErr bool
	}{
		{"Valid passphrase", "MyPass123", false},
		{"Strong passphrase", "SuperSecure999!", false},
		{"Too short", "Pass1", true},
		{"No number", "Password", true},
		{"No letter", "12345678", true},
		{"Empty", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := encService.ValidatePassphrase(tt.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassphrase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPINHashing(t *testing.T) {
	encService := crypto.NewEncryptionService()
	
	pin := "123456"
	
	// Hash
	hash, err := encService.HashPINOrPassphrase(pin)
	if err != nil {
		t.Fatalf("Failed to hash PIN: %v", err)
	}
	
	// Verify correct PIN
	if !encService.VerifyPINOrPassphrase(pin, hash) {
		t.Error("Failed to verify correct PIN")
	}
	
	// Verify incorrect PIN
	if encService.VerifyPINOrPassphrase("654321", hash) {
		t.Error("Incorrectly verified wrong PIN")
	}
}
