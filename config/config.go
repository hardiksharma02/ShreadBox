package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	MaxFileSize     int64 // in bytes
	StoragePath     string
	CleanupInterval time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		MaxFileSize:     getMaxFileSizeBytes(),
		StoragePath:     getEnvOrDefault("STORAGE_PATH", "./storage"),
		CleanupInterval: getCleanupInterval(),
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(config.StoragePath, 0755); err != nil {
		panic("Failed to create storage directory: " + err.Error())
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getMaxFileSizeBytes() int64 {
	maxSizeMB := getEnvOrDefault("MAX_FILE_SIZE", "10") // Default 10MB
	maxSize, err := strconv.ParseInt(maxSizeMB, 10, 64)
	if err != nil {
		maxSize = 10 // Default to 10MB if parsing fails
	}
	return maxSize * 1024 * 1024 // Convert to bytes
}

func getCleanupInterval() time.Duration {
	interval := getEnvOrDefault("CLEANUP_INTERVAL", "5m")
	duration, err := time.ParseDuration(interval)
	if err != nil {
		duration = 5 * time.Minute // Default to 5 minutes if parsing fails
	}
	return duration
}
