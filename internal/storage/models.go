package storage

import (
	"time"
)

// FileMetadata represents the metadata for a stored file
type FileMetadata struct {
	ID            string    `json:"id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	EncryptionKey []byte    `json:"-"` // Not exposed in JSON
	ExpiresAt     time.Time `json:"expires_at"`
	DownloadsLeft int       `json:"downloads_left"`
	Message       string    `json:"message,omitempty"`
	ContentType   string    `json:"content_type"`
	FileSize      int64     `json:"file_size"`
	CreatedAt     time.Time `json:"created_at"`
}

// FileUploadResponse represents the response sent back to the client after a successful upload
type FileUploadResponse struct {
	Token       string    `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
	FileName    string    `json:"file_name"`
	DownloadURL string    `json:"download_url"`
}
