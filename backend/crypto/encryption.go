package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

var encryptionKey []byte

// Init loads the encryption key from environment variable.
// Must be called before using Encrypt/Decrypt functions.
func Init() error {
	keyStr := os.Getenv("WHITELABEL_ENCRYPTION_KEY")
	if keyStr == "" {
		return fmt.Errorf("WHITELABEL_ENCRYPTION_KEY environment variable not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return fmt.Errorf("failed to decode encryption key: %w", err)
	}

	if len(key) != 32 {
		return fmt.Errorf("encryption key must be 32 bytes (256 bits), got %d bytes", len(key))
	}

	encryptionKey = key
	return nil
}

// IsInitialized returns true if the encryption key has been loaded.
func IsInitialized() bool {
	return len(encryptionKey) == 32
}

// Encrypt encrypts plaintext using AES-256-GCM and returns a base64-encoded string.
// The nonce is prepended to the ciphertext.
func Encrypt(plaintext string) (string, error) {
	if !IsInitialized() {
		return "", fmt.Errorf("encryption not initialized")
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal appends the ciphertext to nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext that was encrypted with Encrypt.
func Decrypt(encoded string) (string, error) {
	if !IsInitialized() {
		return "", fmt.Errorf("encryption not initialized")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
