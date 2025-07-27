package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const (
	// KeySize is the size of the encryption key in bytes (32 bytes = 256 bits)
	KeySize = 32
)

var (
	ErrInvalidKeySize = errors.New("invalid key size: key must be 32 bytes")
	ErrEncryption     = errors.New("encryption failed")
	ErrDecryption     = errors.New("decryption failed")
)

// GenerateKey generates a new random 32-byte key
func GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt encrypts data using AES-GCM
func Encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeySize
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrEncryption
	}

	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrEncryption
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, ErrEncryption
	}

	// Encrypt and seal data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-GCM
func Decrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeySize
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrDecryption
	}

	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrDecryption
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, ErrDecryption
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryption
	}

	return plaintext, nil
}

// EncryptFile encrypts a file's contents
func EncryptFile(fileData []byte) ([]byte, []byte, error) {
	// Generate a new key for this file
	key, err := GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the file data
	encryptedData, err := Encrypt(fileData, key)
	if err != nil {
		return nil, nil, err
	}

	return encryptedData, key, nil
}

// DecryptFile decrypts a file's contents using the provided key
func DecryptFile(encryptedData []byte, key []byte) ([]byte, error) {
	return Decrypt(encryptedData, key)
}
