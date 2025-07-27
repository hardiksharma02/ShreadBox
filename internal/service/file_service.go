package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/hardiksharma/shreadbox/internal/domain"
)

type fileService struct {
	repo      domain.FileRepository
	encryptor domain.FileEncryptor
}

// NewFileService creates a new file service instance
func NewFileService(repo domain.FileRepository, encryptor domain.FileEncryptor) domain.FileService {
	return &fileService{
		repo:      repo,
		encryptor: encryptor,
	}
}

// Upload handles the file upload process
func (s *fileService) Upload(name string, size int64, contentType string, data io.Reader, expiryDuration time.Duration, downloads int, message string) (*domain.FileResponse, error) {
	// Read file data
	fileData, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	// Encrypt file data
	encryptedData, key, err := s.encryptor.Encrypt(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Create file entity
	file := &domain.File{
		Name:          name,
		Size:          size,
		ContentType:   contentType,
		Data:          encryptedData,
		EncryptionKey: key,
		ExpiresAt:     time.Now().Add(expiryDuration),
		DownloadsLeft: downloads,
		Message:       message,
		CreatedAt:     time.Now(),
	}

	// Save file
	if err := s.repo.Save(file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate download URL
	downloadURL := fmt.Sprintf("/api/download/%s", file.ID)

	// Return response
	return &domain.FileResponse{
		Token:       file.ID,
		ExpiresAt:   file.ExpiresAt,
		FileName:    file.Name,
		DownloadURL: downloadURL,
	}, nil
}

// Download handles the file download process
func (s *fileService) Download(id string) (*domain.File, error) {
	// Get file
	file, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Decrypt file data
	decryptedData, err := s.encryptor.Decrypt(file.Data, file.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %w", err)
	}

	// Update file data with decrypted content
	file.Data = decryptedData

	return file, nil
}

// GetStatus retrieves the current status of a file
func (s *fileService) GetStatus(id string) (*domain.FileStatus, error) {
	file, err := s.repo.GetMetadata(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	return &domain.FileStatus{
		FileName:      file.Name,
		ExpiresAt:     file.ExpiresAt,
		DownloadsLeft: file.DownloadsLeft,
		Message:       file.Message,
	}, nil
}
