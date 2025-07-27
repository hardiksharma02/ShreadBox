package domain

import (
	"io"
	"time"
)

// File represents the core file entity
type File struct {
	ID            string
	Name          string
	Size          int64
	ContentType   string
	Data          []byte
	EncryptionKey []byte
	ExpiresAt     time.Time
	DownloadsLeft int
	Message       string
	CreatedAt     time.Time
}

// FileRepository defines the interface for file storage operations
type FileRepository interface {
	Save(file *File) error
	Get(id string) (*File, error)
	Delete(id string) error
	GetMetadata(id string) (*File, error)
	CleanupExpired() error
}

// FileEncryptor defines the interface for file encryption operations
type FileEncryptor interface {
	Encrypt(data []byte) ([]byte, []byte, error)
	Decrypt(data []byte, key []byte) ([]byte, error)
}

// FileService defines the interface for file business logic
type FileService interface {
	Upload(name string, size int64, contentType string, data io.Reader, expiryDuration time.Duration, downloads int, message string) (*FileResponse, error)
	Download(id string) (*File, error)
	GetStatus(id string) (*FileStatus, error)
}

// FileResponse represents the response after successful file upload
type FileResponse struct {
	Token       string    `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
	FileName    string    `json:"file_name"`
	DownloadURL string    `json:"download_url"`
}

// FileStatus represents the current status of a file
type FileStatus struct {
	FileName      string    `json:"file_name"`
	ExpiresAt     time.Time `json:"expires_at"`
	DownloadsLeft int       `json:"downloads_left"`
	Message       string    `json:"message,omitempty"`
}
