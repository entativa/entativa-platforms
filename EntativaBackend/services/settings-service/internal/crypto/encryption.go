package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// PBKDF2 iterations for key derivation (100,000 iterations)
	PBKDF2Iterations = 100000
	
	// Key size for AES-256
	KeySize = 32
	
	// Salt size
	SaltSize = 32
	
	// bcrypt cost for PIN/Passphrase hashing
	BcryptCost = 12
)

// EncryptionService handles encryption/decryption of chat keys
type EncryptionService struct{}

func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// EncryptKeys encrypts chat keys with user's PIN/Passphrase
// This is the SECOND layer of encryption (first is Signal protocol)
func (e *EncryptionService) EncryptKeys(keys []byte, pinOrPassphrase string, salt []byte) ([]byte, error) {
	// Derive encryption key from PIN/Passphrase using PBKDF2
	key := pbkdf2.Key([]byte(pinOrPassphrase), salt, PBKDF2Iterations, KeySize, sha256.New)
	
	// Create AES-256-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	
	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, keys, nil)
	
	return ciphertext, nil
}

// DecryptKeys decrypts chat keys with user's PIN/Passphrase
func (e *EncryptionService) DecryptKeys(encryptedKeys []byte, pinOrPassphrase string, salt []byte) ([]byte, error) {
	// Derive decryption key from PIN/Passphrase using PBKDF2
	key := pbkdf2.Key([]byte(pinOrPassphrase), salt, PBKDF2Iterations, KeySize, sha256.New)
	
	// Create AES-256-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(encryptedKeys) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext := encryptedKeys[:nonceSize], encryptedKeys[nonceSize:]
	
	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	
	return plaintext, nil
}

// HashPINOrPassphrase hashes PIN/Passphrase using bcrypt
func (e *EncryptionService) HashPINOrPassphrase(pinOrPassphrase string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pinOrPassphrase), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPINOrPassphrase verifies PIN/Passphrase against hash
func (e *EncryptionService) VerifyPINOrPassphrase(pinOrPassphrase, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pinOrPassphrase))
	return err == nil
}

// GenerateSalt generates a random salt
func (e *EncryptionService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// HashData computes SHA256 hash of data (for integrity verification)
func (e *EncryptionService) HashData(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// ValidatePIN validates PIN format (6 digits)
func (e *EncryptionService) ValidatePIN(pin string) error {
	if len(pin) != 6 {
		return fmt.Errorf("PIN must be exactly 6 digits")
	}
	
	for _, c := range pin {
		if c < '0' || c > '9' {
			return fmt.Errorf("PIN must contain only digits")
		}
	}
	
	return nil
}

// ValidatePassphrase validates passphrase strength
func (e *EncryptionService) ValidatePassphrase(passphrase string) error {
	if len(passphrase) < 8 {
		return fmt.Errorf("passphrase must be at least 8 characters")
	}
	
	if len(passphrase) > 128 {
		return fmt.Errorf("passphrase must be at most 128 characters")
	}
	
	// Check for complexity (at least one letter and one number)
	hasLetter := false
	hasNumber := false
	
	for _, c := range passphrase {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}
	
	if !hasLetter {
		return fmt.Errorf("passphrase must contain at least one letter")
	}
	if !hasNumber {
		return fmt.Errorf("passphrase must contain at least one number")
	}
	
	return nil
}

// EncodeBase64 encodes bytes to base64
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes base64 to bytes
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
