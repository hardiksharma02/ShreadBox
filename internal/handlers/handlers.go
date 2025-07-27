package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardiksharma/shreadbox/internal/encryption"
	"github.com/hardiksharma/shreadbox/internal/storage"
)

// Handler represents the HTTP handler
type Handler struct {
	storage *storage.Storage
}

// NewHandler creates a new handler instance
func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

// Upload handles file upload requests
func (h *Handler) Upload(c *gin.Context) {
	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Read file data
	fileData := make([]byte, header.Size)
	if _, err := file.Read(fileData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Parse form parameters
	expiryTime := c.PostForm("expiry_time")
	downloadsStr := c.PostForm("downloads_allowed")
	message := c.PostForm("message")

	// Parse expiry time
	duration, err := time.ParseDuration(expiryTime)
	if err != nil {
		duration = 24 * time.Hour // Default to 24 hours
	}

	// Parse downloads allowed
	downloads, err := strconv.Atoi(downloadsStr)
	if err != nil || downloads < 1 {
		downloads = 1 // Default to 1 download
	}

	// Encrypt file
	encryptedData, key, err := encryption.EncryptFile(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	// Create metadata
	metadata := &storage.FileMetadata{
		FileName:      header.Filename,
		EncryptionKey: key,
		ExpiresAt:     time.Now().Add(duration),
		DownloadsLeft: downloads,
		Message:       message,
		ContentType:   header.Header.Get("Content-Type"),
		FileSize:      header.Size,
	}

	// Save file
	if err := h.storage.SaveFile(encryptedData, metadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate download URL
	downloadURL := fmt.Sprintf("/api/download/%s", metadata.ID)

	// Return response
	c.JSON(http.StatusOK, storage.FileUploadResponse{
		Token:       metadata.ID,
		ExpiresAt:   metadata.ExpiresAt,
		FileName:    metadata.FileName,
		DownloadURL: downloadURL,
	})
}

// Download handles file download requests
func (h *Handler) Download(c *gin.Context) {
	// Get file ID from URL
	fileID := c.Param("token")

	// Get file metadata
	metadata, err := h.storage.GetFile(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or expired"})
		return
	}

	// Read encrypted file
	encryptedData, err := h.storage.ReadFile(metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Decrypt file
	decryptedData, err := encryption.DecryptFile(encryptedData, metadata.EncryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		return
	}

	// Set response headers
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(metadata.FileName)))
	c.Header("Content-Type", metadata.ContentType)
	c.Header("Content-Length", strconv.FormatInt(metadata.FileSize, 10))

	// Send file
	c.Data(http.StatusOK, metadata.ContentType, decryptedData)
}

// Status handles file status requests
func (h *Handler) Status(c *gin.Context) {
	fileID := c.Param("token")

	metadata, err := h.storage.GetFileMetadata(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_name":      metadata.FileName,
		"expires_at":     metadata.ExpiresAt,
		"downloads_left": metadata.DownloadsLeft,
		"message":        metadata.Message,
	})
}
