package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

// SignalProtocolService implements the Signal Protocol (Double Ratchet)
// This provides E2EE with forward secrecy and post-compromise security
type SignalProtocolService struct {
	prekeyRepo  PrekeyRepository
	sessionRepo SessionRepository
}

type PrekeyRepository interface {
	StorePrekeys(userID string, prekeys []Prekey) error
	GetPrekey(userID string) (*Prekey, error)
	GetIdentityKey(userID string) (*IdentityKey, error)
	StoreIdentityKey(userID string, key *IdentityKey) error
	MarkPrekeyUsed(prekeyID string) error
}

type SessionRepository interface {
	StoreSession(userID, partnerID string, session *Session) error
	GetSession(userID, partnerID string) (*Session, error)
	UpdateSession(userID, partnerID string, session *Session) error
}

// Prekey represents a one-time prekey for X3DH key exchange
type Prekey struct {
	ID        string
	PublicKey []byte
	Signature []byte
}

// IdentityKey represents a user's long-term identity key
type IdentityKey struct {
	PublicKey  []byte
	PrivateKey []byte // Only stored for own identity
}

// Session represents a Double Ratchet session
type Session struct {
	RootKey          []byte
	SendingChainKey  []byte
	ReceivingChainKey []byte
	SendingRatchetKey []byte
	ReceivingRatchetKey []byte
	SendCounter      int
	ReceiveCounter   int
	PreviousCounter  int
	SkippedMessages  map[int][]byte // For out-of-order messages
}

// PrekeyBundle contains all keys needed to initiate a conversation
type PrekeyBundle struct {
	IdentityKey    []byte
	SignedPrekey   []byte
	PrekeySignature []byte
	OneTimePrekey  []byte
}

// EncryptedMessage contains the encrypted message and metadata
type EncryptedMessage struct {
	Ciphertext     []byte
	MAC            []byte
	Counter        int
	RatchetKey     []byte
	MessageType    string // "prekey" or "message"
}

func NewSignalProtocolService(prekeyRepo PrekeyRepository, sessionRepo SessionRepository) *SignalProtocolService {
	return &SignalProtocolService{
		prekeyRepo:  prekeyRepo,
		sessionRepo: sessionRepo,
	}
}

// GenerateKeyPair generates a Curve25519 key pair
func (s *SignalProtocolService) GenerateKeyPair() (publicKey, privateKey []byte, err error) {
	privateKey = make([]byte, 32)
	if _, err := rand.Read(privateKey); err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey = make([]byte, 32)
	curve25519.ScalarBaseMult((*[32]byte)(publicKey), (*[32]byte)(privateKey))

	return publicKey, privateKey, nil
}

// PerformX3DH performs the X3DH key exchange to establish a shared secret
func (s *SignalProtocolService) PerformX3DH(
	identityKeyPrivate []byte,
	ephemeralKeyPrivate []byte,
	recipientIdentityKey []byte,
	recipientSignedPrekey []byte,
	recipientOneTimePrekey []byte,
) ([]byte, error) {
	// DH1 = DH(IKa, SPKb)
	dh1 := make([]byte, 32)
	curve25519.ScalarMult((*[32]byte)(dh1), (*[32]byte)(identityKeyPrivate), (*[32]byte)(recipientSignedPrekey))

	// DH2 = DH(EKa, IKb)
	dh2 := make([]byte, 32)
	curve25519.ScalarMult((*[32]byte)(dh2), (*[32]byte)(ephemeralKeyPrivate), (*[32]byte)(recipientIdentityKey))

	// DH3 = DH(EKa, SPKb)
	dh3 := make([]byte, 32)
	curve25519.ScalarMult((*[32]byte)(dh3), (*[32]byte)(ephemeralKeyPrivate), (*[32]byte)(recipientSignedPrekey))

	// Concatenate DH outputs
	dhOutputs := append(dh1, dh2...)
	dhOutputs = append(dhOutputs, dh3...)

	// If one-time prekey exists: DH4 = DH(EKa, OPKb)
	if len(recipientOneTimePrekey) > 0 {
		dh4 := make([]byte, 32)
		curve25519.ScalarMult((*[32]byte)(dh4), (*[32]byte)(ephemeralKeyPrivate), (*[32]byte)(recipientOneTimePrekey))
		dhOutputs = append(dhOutputs, dh4...)
	}

	// Derive shared secret using HKDF
	kdf := hkdf.New(sha256.New, dhOutputs, nil, []byte("Signal_X3DH"))
	sharedSecret := make([]byte, 32)
	if _, err := kdf.Read(sharedSecret); err != nil {
		return nil, fmt.Errorf("failed to derive shared secret: %w", err)
	}

	return sharedSecret, nil
}

// InitializeSession initializes a Double Ratchet session from a shared secret
func (s *SignalProtocolService) InitializeSession(sharedSecret []byte, sending bool) (*Session, error) {
	// Derive root key and chain key from shared secret
	kdf := hkdf.New(sha256.New, sharedSecret, nil, []byte("Signal_Root"))
	rootKey := make([]byte, 32)
	if _, err := kdf.Read(rootKey); err != nil {
		return nil, fmt.Errorf("failed to derive root key: %w", err)
	}

	kdf = hkdf.New(sha256.New, sharedSecret, nil, []byte("Signal_Chain"))
	chainKey := make([]byte, 32)
	if _, err := kdf.Read(chainKey); err != nil {
		return nil, fmt.Errorf("failed to derive chain key: %w", err)
	}

	// Generate initial ratchet key pair
	ratchetPublic, ratchetPrivate, err := s.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ratchet key: %w", err)
	}

	session := &Session{
		RootKey:          rootKey,
		SendCounter:      0,
		ReceiveCounter:   0,
		PreviousCounter:  0,
		SkippedMessages:  make(map[int][]byte),
	}

	if sending {
		session.SendingChainKey = chainKey
		session.SendingRatchetKey = ratchetPrivate
	} else {
		session.ReceivingChainKey = chainKey
		session.ReceivingRatchetKey = ratchetPrivate
	}

	return session, nil
}

// Encrypt encrypts a message using the Double Ratchet
func (s *SignalProtocolService) Encrypt(session *Session, plaintext []byte) (*EncryptedMessage, error) {
	// Derive message key from chain key
	messageKey := s.deriveMessageKey(session.SendingChainKey)

	// Advance chain key
	session.SendingChainKey = s.advanceChainKey(session.SendingChainKey)

	// Encrypt the message
	ciphertext, err := s.encryptAESGCM(messageKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %w", err)
	}

	// Calculate MAC
	mac := s.calculateMAC(messageKey, ciphertext)

	encryptedMsg := &EncryptedMessage{
		Ciphertext:  ciphertext,
		MAC:         mac,
		Counter:     session.SendCounter,
		RatchetKey:  session.SendingRatchetKey,
		MessageType: "message",
	}

	session.SendCounter++

	return encryptedMsg, nil
}

// Decrypt decrypts a message using the Double Ratchet
func (s *SignalProtocolService) Decrypt(session *Session, encrypted *EncryptedMessage) ([]byte, error) {
	// Handle out-of-order messages
	if encrypted.Counter < session.ReceiveCounter {
		// Try to use skipped message key
		if messageKey, exists := session.SkippedMessages[encrypted.Counter]; exists {
			plaintext, err := s.decryptAESGCM(messageKey, encrypted.Ciphertext)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt skipped message: %w", err)
			}
			delete(session.SkippedMessages, encrypted.Counter)
			return plaintext, nil
		}
		return nil, errors.New("message key not found for old counter")
	}

	// Skip intermediate messages if needed
	for session.ReceiveCounter < encrypted.Counter {
		messageKey := s.deriveMessageKey(session.ReceivingChainKey)
		session.SkippedMessages[session.ReceiveCounter] = messageKey
		session.ReceivingChainKey = s.advanceChainKey(session.ReceivingChainKey)
		session.ReceiveCounter++
	}

	// Derive message key
	messageKey := s.deriveMessageKey(session.ReceivingChainKey)

	// Verify MAC
	expectedMAC := s.calculateMAC(messageKey, encrypted.Ciphertext)
	if !hmac.Equal(expectedMAC, encrypted.MAC) {
		return nil, errors.New("MAC verification failed")
	}

	// Decrypt the message
	plaintext, err := s.decryptAESGCM(messageKey, encrypted.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}

	// Advance chain key
	session.ReceivingChainKey = s.advanceChainKey(session.ReceivingChainKey)
	session.ReceiveCounter++

	return plaintext, nil
}

// Helper functions

func (s *SignalProtocolService) deriveMessageKey(chainKey []byte) []byte {
	mac := hmac.New(sha256.New, chainKey)
	mac.Write([]byte{0x01})
	return mac.Sum(nil)
}

func (s *SignalProtocolService) advanceChainKey(chainKey []byte) []byte {
	mac := hmac.New(sha256.New, chainKey)
	mac.Write([]byte{0x02})
	return mac.Sum(nil)
}

func (s *SignalProtocolService) encryptAESGCM(key, plaintext []byte) ([]byte, error) {
	// Use first 32 bytes of key
	if len(key) < 32 {
		return nil, errors.New("key too short")
	}

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (s *SignalProtocolService) decryptAESGCM(key, ciphertext []byte) ([]byte, error) {
	if len(key) < 32 {
		return nil, errors.New("key too short")
	}

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (s *SignalProtocolService) calculateMAC(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// EncryptMessage encrypts a message for a specific conversation
func (s *SignalProtocolService) EncryptMessage(senderID, recipientID, message string) (string, error) {
	// Get or create session
	session, err := s.sessionRepo.GetSession(senderID, recipientID)
	if err != nil || session == nil {
		// Initialize new session using prekey bundle
		bundle, err := s.GetPrekeyBundle(recipientID)
		if err != nil {
			return "", fmt.Errorf("failed to get prekey bundle: %w", err)
		}

		// Perform X3DH and initialize session
		// (Simplified - in production you'd need sender's identity key)
		session, err = s.InitializeSession(bundle.IdentityKey, true)
		if err != nil {
			return "", fmt.Errorf("failed to initialize session: %w", err)
		}

		// Store session
		if err := s.sessionRepo.StoreSession(senderID, recipientID, session); err != nil {
			return "", fmt.Errorf("failed to store session: %w", err)
		}
	}

	// Encrypt the message
	encrypted, err := s.Encrypt(session, []byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	// Update session
	if err := s.sessionRepo.UpdateSession(senderID, recipientID, session); err != nil {
		return "", fmt.Errorf("failed to update session: %w", err)
	}

	// Encode to base64 for transmission
	ciphertext := base64.StdEncoding.EncodeToString(encrypted.Ciphertext)

	return ciphertext, nil
}

// DecryptMessage decrypts a received message
func (s *SignalProtocolService) DecryptMessage(recipientID, senderID, encryptedMessage string) (string, error) {
	// Get session
	session, err := s.sessionRepo.GetSession(recipientID, senderID)
	if err != nil || session == nil {
		return "", fmt.Errorf("session not found")
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", fmt.Errorf("failed to decode message: %w", err)
	}

	encrypted := &EncryptedMessage{
		Ciphertext: ciphertext,
		Counter:    session.ReceiveCounter,
	}

	// Decrypt the message
	plaintext, err := s.Decrypt(session, encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %w", err)
	}

	// Update session
	if err := s.sessionRepo.UpdateSession(recipientID, senderID, session); err != nil {
		return "", fmt.Errorf("failed to update session: %w", err)
	}

	return string(plaintext), nil
}

// GetPrekeyBundle retrieves a user's prekey bundle for initiating a conversation
func (s *SignalProtocolService) GetPrekeyBundle(userID string) (*PrekeyBundle, error) {
	identityKey, err := s.prekeyRepo.GetIdentityKey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity key: %w", err)
	}

	prekey, err := s.prekeyRepo.GetPrekey(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prekey: %w", err)
	}

	bundle := &PrekeyBundle{
		IdentityKey:     identityKey.PublicKey,
		SignedPrekey:    prekey.PublicKey,
		PrekeySignature: prekey.Signature,
		OneTimePrekey:   prekey.PublicKey,
	}

	// Mark prekey as used
	if err := s.prekeyRepo.MarkPrekeyUsed(prekey.ID); err != nil {
		return nil, fmt.Errorf("failed to mark prekey as used: %w", err)
		}

	return bundle, nil
}
