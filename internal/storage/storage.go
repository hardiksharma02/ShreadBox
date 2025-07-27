package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Storage represents the file storage service
type Storage struct {
	basePath string
	files    map[string]*FileMetadata
	mu       sync.RWMutex
}

// NewStorage creates a new storage service
func NewStorage(basePath string) (*Storage, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Storage{
		basePath: basePath,
		files:    make(map[string]*FileMetadata),
	}, nil
}

// SaveFile saves an encrypted file and its metadata
func (s *Storage) SaveFile(data []byte, metadata *FileMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate unique ID if not provided
	if metadata.ID == "" {
		metadata.ID = uuid.New().String()
	}

	// Set file path
	metadata.FilePath = filepath.Join(s.basePath, metadata.ID)
	metadata.CreatedAt = time.Now()

	// Save encrypted file
	if err := os.WriteFile(metadata.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Store metadata
	s.files[metadata.ID] = metadata
	return nil
}

// GetFile retrieves a file's metadata and decrements the download counter
func (s *Storage) GetFile(id string) (*FileMetadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, exists := s.files[id]
	if !exists {
		return nil, fmt.Errorf("file not found")
	}

	// Check if file has expired
	if time.Now().After(metadata.ExpiresAt) {
		s.deleteFile(id)
		return nil, fmt.Errorf("file has expired")
	}

	// Check if downloads are exhausted
	if metadata.DownloadsLeft <= 0 {
		s.deleteFile(id)
		return nil, fmt.Errorf("download limit reached")
	}

	// Decrement download counter
	metadata.DownloadsLeft--

	// If no downloads left, schedule deletion
	if metadata.DownloadsLeft == 0 {
		go func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			s.deleteFile(id)
		}()
	}

	return metadata, nil
}

// ReadFile reads the encrypted file content
func (s *Storage) ReadFile(metadata *FileMetadata) ([]byte, error) {
	return os.ReadFile(metadata.FilePath)
}

// deleteFile removes a file and its metadata
func (s *Storage) deleteFile(id string) error {
	metadata, exists := s.files[id]
	if !exists {
		return nil
	}

	// Remove file from disk
	if err := os.Remove(metadata.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Remove metadata from memory
	delete(s.files, id)
	return nil
}

// CleanupExpired removes expired files
func (s *Storage) CleanupExpired() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var lastErr error

	for id, metadata := range s.files {
		if now.After(metadata.ExpiresAt) || metadata.DownloadsLeft <= 0 {
			if err := s.deleteFile(id); err != nil {
				lastErr = err
			}
		}
	}

	return lastErr
}

// GetFileMetadata retrieves file metadata without modifying the download counter
func (s *Storage) GetFileMetadata(id string) (*FileMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadata, exists := s.files[id]
	if !exists {
		return nil, fmt.Errorf("file not found")
	}

	return metadata, nil
}
