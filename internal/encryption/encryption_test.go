package encryption

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	// Generate a key
	key, err := GenerateKey()

	// Assert no error
	assert.NoError(t, err)

	// Assert key length
	assert.Equal(t, KeySize, len(key))

	// Generate another key and ensure it's different
	key2, err := GenerateKey()
	assert.NoError(t, err)
	assert.NotEqual(t, key, key2, "Generated keys should be unique")
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		shouldError bool
	}{
		{
			name:        "valid data",
			data:        []byte("test data"),
			shouldError: false,
		},
		{
			name:        "empty data",
			data:        make([]byte, 0),
			shouldError: false,
		},
		{
			name:        "large data",
			data:        bytes.Repeat([]byte("a"), 1000),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate a key
			key, err := GenerateKey()
			assert.NoError(t, err)

			// Encrypt data
			encrypted, err := Encrypt(tt.data, key)
			if tt.shouldError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEqual(t, tt.data, encrypted)

			// Decrypt data
			decrypted, err := Decrypt(encrypted, key)
			assert.NoError(t, err)
			if len(tt.data) == 0 {
				assert.Empty(t, decrypted)
			} else {
				assert.Equal(t, tt.data, decrypted)
			}
		})
	}
}

func TestEncryptDecrypt_InvalidKey(t *testing.T) {
	data := []byte("test data")

	// Test with invalid key size
	invalidKey := make([]byte, KeySize-1)
	_, err := Encrypt(data, invalidKey)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidKeySize, err)

	// Test decryption with invalid key
	key, _ := GenerateKey()
	encrypted, _ := Encrypt(data, key)
	_, err = Decrypt(encrypted, invalidKey)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidKeySize, err)
}

func TestEncryptFile(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		shouldError bool
	}{
		{
			name:        "valid file",
			data:        []byte("test file content"),
			shouldError: false,
		},
		{
			name:        "empty file",
			data:        make([]byte, 0),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt file
			encrypted, key, err := EncryptFile(tt.data)
			if tt.shouldError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, encrypted)
			assert.NotNil(t, key)
			assert.Len(t, key, KeySize)

			// Decrypt and verify
			decrypted, err := DecryptFile(encrypted, key)
			assert.NoError(t, err)
			if len(tt.data) == 0 {
				assert.Empty(t, decrypted)
			} else {
				assert.Equal(t, tt.data, decrypted)
			}
		})
	}
}
