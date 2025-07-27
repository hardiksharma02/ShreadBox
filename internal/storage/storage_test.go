package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	// Create temporary directory for testing
	tempDir := filepath.Join(os.TempDir(), "shreadbox-test")
	defer os.RemoveAll(tempDir)

	// Test successful creation
	storage, err := NewStorage(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
	assert.DirExists(t, tempDir)

	// Test with invalid path
	invalidPath := string([]byte{0})
	storage, err = NewStorage(invalidPath)
	assert.Error(t, err)
	assert.Nil(t, storage)
}

func TestStorage_SaveFile(t *testing.T) {
	// Setup
	tempDir := filepath.Join(os.TempDir(), "shreadbox-test")
	defer os.RemoveAll(tempDir)

	storage, err := NewStorage(tempDir)
	assert.NoError(t, err)

	// Test data
	testData := []byte("test file content")
	metadata := &FileMetadata{
		FileName:      "test.txt",
		EncryptionKey: []byte("test-key"),
		ExpiresAt:     time.Now().Add(time.Hour),
		DownloadsLeft: 1,
		Message:       "test message",
		ContentType:   "text/plain",
		FileSize:      int64(len(testData)),
	}

	// Test saving file
	err = storage.SaveFile(testData, metadata)
	assert.NoError(t, err)
	assert.NotEmpty(t, metadata.ID)
	assert.FileExists(t, metadata.FilePath)

	// Verify metadata was stored
	stored, err := storage.GetFileMetadata(metadata.ID)
	assert.NoError(t, err)
	assert.Equal(t, metadata.FileName, stored.FileName)
	assert.Equal(t, metadata.DownloadsLeft, stored.DownloadsLeft)
}

func TestStorage_GetFile(t *testing.T) {
	// Setup
	tempDir := filepath.Join(os.TempDir(), "shreadbox-test")
	defer os.RemoveAll(tempDir)

	storage, err := NewStorage(tempDir)
	assert.NoError(t, err)

	// Save a test file
	testData := []byte("test file content")
	metadata := &FileMetadata{
		FileName:      "test.txt",
		EncryptionKey: []byte("test-key"),
		ExpiresAt:     time.Now().Add(time.Hour),
		DownloadsLeft: 1,
		ContentType:   "text/plain",
		FileSize:      int64(len(testData)),
	}

	err = storage.SaveFile(testData, metadata)
	assert.NoError(t, err)

	// Test getting file
	retrieved, err := storage.GetFile(metadata.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, metadata.FileName, retrieved.FileName)
	assert.Equal(t, 0, retrieved.DownloadsLeft) // Should be decremented

	// Test getting non-existent file
	retrieved, err = storage.GetFile("non-existent")
	assert.Error(t, err)
	assert.Nil(t, retrieved)
}

func TestStorage_CleanupExpired(t *testing.T) {
	// Setup
	tempDir := filepath.Join(os.TempDir(), "shreadbox-test")
	defer os.RemoveAll(tempDir)

	storage, err := NewStorage(tempDir)
	assert.NoError(t, err)

	// Save files with different expiry times
	testCases := []struct {
		name      string
		expiresAt time.Time
		shouldBe  string
	}{
		{
			name:      "expired.txt",
			expiresAt: time.Now().Add(-time.Hour),
			shouldBe:  "deleted",
		},
		{
			name:      "valid.txt",
			expiresAt: time.Now().Add(time.Hour),
			shouldBe:  "present",
		},
	}

	savedFiles := make(map[string]string) // map[filename]id

	for _, tc := range testCases {
		metadata := &FileMetadata{
			FileName:      tc.name,
			EncryptionKey: []byte("test-key"),
			ExpiresAt:     tc.expiresAt,
			DownloadsLeft: 1,
			ContentType:   "text/plain",
			FileSize:      10,
		}

		err = storage.SaveFile([]byte("test"), metadata)
		assert.NoError(t, err)
		savedFiles[tc.name] = metadata.ID
	}

	// Run cleanup
	storage.CleanupExpired()

	// Verify results
	for _, tc := range testCases {
		_, err := storage.GetFileMetadata(savedFiles[tc.name])
		if tc.shouldBe == "deleted" {
			assert.Error(t, err, "File %s should have been deleted", tc.name)
		} else {
			assert.NoError(t, err, "File %s should still exist", tc.name)
		}
	}
}
