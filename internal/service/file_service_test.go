package service

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/hardiksharma/shreadbox/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFileRepository struct {
	mock.Mock
}

func (m *mockFileRepository) Save(file *domain.File) error {
	args := m.Called(file)
	if args.Error(0) == nil {
		file.ID = "test-id" // Set ID for successful saves
	}
	return args.Error(0)
}

func (m *mockFileRepository) Get(id string) (*domain.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockFileRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockFileRepository) GetMetadata(id string) (*domain.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockFileRepository) CleanupExpired() error {
	args := m.Called()
	return args.Error(0)
}

type mockFileEncryptor struct {
	mock.Mock
}

func (m *mockFileEncryptor) Encrypt(data []byte) ([]byte, []byte, error) {
	args := m.Called(data)
	if args.Error(2) != nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]byte), args.Get(1).([]byte), nil
}

func (m *mockFileEncryptor) Decrypt(data []byte, key []byte) ([]byte, error) {
	args := m.Called(data, key)
	return args.Get(0).([]byte), args.Error(1)
}

func TestFileService_Upload(t *testing.T) {
	repo := new(mockFileRepository)
	encryptor := new(mockFileEncryptor)
	service := NewFileService(repo, encryptor)

	tests := []struct {
		name          string
		fileData      []byte
		setupMocks    func()
		expectedError bool
	}{
		{
			name:     "successful upload",
			fileData: []byte("test data"),
			setupMocks: func() {
				encryptedData := []byte("encrypted")
				key := []byte("key")
				encryptor.On("Encrypt", []byte("test data")).Return(encryptedData, key, nil)
				repo.On("Save", mock.AnythingOfType("*domain.File")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "encryption fails",
			fileData: []byte("test data"),
			setupMocks: func() {
				encryptor.On("Encrypt", []byte("test data")).Return([]byte{}, []byte{}, errors.New("encryption failed"))
				// Don't set up repo.Save since encryption should fail
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			repo = new(mockFileRepository)
			encryptor = new(mockFileEncryptor)
			service = NewFileService(repo, encryptor)

			// Setup mocks
			tt.setupMocks()

			// Create test data
			reader := bytes.NewReader(tt.fileData)
			result, err := service.Upload("test.txt", int64(len(tt.fileData)), "text/plain", reader, 24*time.Hour, 1, "test message")

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "test-id", result.Token)
				assert.Equal(t, "test.txt", result.FileName)
				assert.Equal(t, "/api/download/test-id", result.DownloadURL)
			}

			// Verify mock expectations
			encryptor.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}

func TestFileService_Download(t *testing.T) {
	repo := new(mockFileRepository)
	encryptor := new(mockFileEncryptor)
	service := NewFileService(repo, encryptor)

	tests := []struct {
		name          string
		fileID        string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:   "successful download",
			fileID: "test-id",
			setupMocks: func() {
				file := &domain.File{
					ID:            "test-id",
					Data:          []byte("encrypted"),
					EncryptionKey: []byte("key"),
				}
				repo.On("Get", "test-id").Return(file, nil)
				encryptor.On("Decrypt", []byte("encrypted"), []byte("key")).Return([]byte("decrypted"), nil)
			},
			expectedError: false,
		},
		{
			name:   "file not found",
			fileID: "not-found",
			setupMocks: func() {
				repo.On("Get", "not-found").Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			repo = new(mockFileRepository)
			encryptor = new(mockFileEncryptor)
			service = NewFileService(repo, encryptor)

			// Setup mocks
			tt.setupMocks()

			// Attempt download
			result, err := service.Download(tt.fileID)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, []byte("decrypted"), result.Data)
			}

			// Verify mock expectations
			encryptor.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}

func TestFileService_GetStatus(t *testing.T) {
	repo := new(mockFileRepository)
	encryptor := new(mockFileEncryptor)
	service := NewFileService(repo, encryptor)

	tests := []struct {
		name          string
		fileID        string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:   "successful status check",
			fileID: "test-id",
			setupMocks: func() {
				file := &domain.File{
					ID:            "test-id",
					Name:          "test.txt",
					DownloadsLeft: 1,
					ExpiresAt:     time.Now().Add(time.Hour),
					Message:       "test message",
				}
				repo.On("GetMetadata", "test-id").Return(file, nil)
			},
			expectedError: false,
		},
		{
			name:   "file not found",
			fileID: "not-found",
			setupMocks: func() {
				repo.On("GetMetadata", "not-found").Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			repo = new(mockFileRepository)
			encryptor = new(mockFileEncryptor)
			service = NewFileService(repo, encryptor)

			// Setup mocks
			tt.setupMocks()

			// Check status
			result, err := service.GetStatus(tt.fileID)

			// Verify results
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "test.txt", result.FileName)
				assert.Equal(t, 1, result.DownloadsLeft)
				assert.Equal(t, "test message", result.Message)
			}

			// Verify mock expectations
			repo.AssertExpectations(t)
		})
	}
}
