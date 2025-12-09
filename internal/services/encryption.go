package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/kodra-pay/encryption-service/internal/dto"
)

type EncryptionService struct {
	encryptionKey []byte
	tokens        map[int]encryptedToken
	tokenCounter  int
	mu            sync.RWMutex
}

type encryptedToken struct {
	EncryptedData string
	Last4         string
	Brand         string
	CreatedAt     time.Time
}

func NewEncryptionService(encryptionKey string) (*EncryptionService, error) {
	// Decode the hex encryption key
	key, err := hex.DecodeString(encryptionKey)
	if err != nil {
		// If not hex, use the raw string (must be 16, 24, or 32 bytes for AES)
		key = []byte(encryptionKey)
	}

	// Ensure key is 32 bytes for AES-256
	if len(key) != 32 {
		// Pad or truncate to 32 bytes
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	}

	return &EncryptionService{
		encryptionKey: key,
		tokens:        make(map[int]encryptedToken),
		tokenCounter:  1000, // Start from 1000 for more realistic token IDs
	}, nil
}

func (s *EncryptionService) TokenizeCard(_ context.Context, req dto.TokenizeCardRequest) dto.TokenResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encrypt the PAN
	encryptedPAN, err := s.encrypt(req.Pan)
	if err != nil {
		return dto.TokenResponse{
			TokenID: 0,
			Last4:   "",
			Brand:   "error",
		}
	}

	// Generate token ID
	s.tokenCounter++
	tokenID := s.tokenCounter

	// Store encrypted token
	s.tokens[tokenID] = encryptedToken{
		EncryptedData: encryptedPAN,
		Last4:         last4(req.Pan),
		Brand:         detectBrand(req.Pan),
		CreatedAt:     time.Now(),
	}

	return dto.TokenResponse{
		TokenID: tokenID,
		Last4:   last4(req.Pan),
		Brand:   detectBrand(req.Pan),
	}
}

func (s *EncryptionService) GetToken(_ context.Context, id int) dto.TokenResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, exists := s.tokens[id]
	if !exists {
		return dto.TokenResponse{
			TokenID: 0,
			Last4:   "0000",
			Brand:   "not_found",
		}
	}

	return dto.TokenResponse{
		TokenID: id,
		Last4:   token.Last4,
		Brand:   token.Brand,
	}
}

func (s *EncryptionService) EncryptData(_ context.Context, data string) (string, error) {
	return s.encrypt(data)
}

func (s *EncryptionService) DecryptData(_ context.Context, encryptedData string) (string, error) {
	return s.decrypt(encryptedData)
}

// encrypt encrypts plaintext using AES-256-GCM
func (s *EncryptionService) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts ciphertext using AES-256-GCM
func (s *EncryptionService) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func last4(pan string) string {
	if len(pan) < 4 {
		return pan
	}
	return pan[len(pan)-4:]
}

func detectBrand(pan string) string {
	if len(pan) < 2 {
		return "unknown"
	}

	// Simple brand detection based on BIN
	switch pan[0] {
	case '4':
		return "visa"
	case '5':
		return "mastercard"
	case '3':
		if len(pan) >= 2 && pan[1] == '7' {
			return "amex"
		}
		return "unknown"
	case '6':
		return "discover"
	default:
		return "unknown"
	}
}
