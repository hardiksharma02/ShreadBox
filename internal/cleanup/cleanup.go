package cleanup

import (
	"log"
	"time"
)

// Service represents the cleanup service
type Service struct {
	storage   StorageInterface
	interval  time.Duration
	stopChan  chan struct{}
	isRunning bool
}

// StorageInterface defines the methods required from the storage service
type StorageInterface interface {
	CleanupExpired() error
}

// NewService creates a new cleanup service
func NewService(storage StorageInterface, interval time.Duration) *Service {
	return &Service{
		storage:  storage,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start begins the cleanup routine
func (s *Service) Start() {
	if s.isRunning {
		return
	}

	s.isRunning = true
	go s.cleanupRoutine()
}

// Stop stops the cleanup routine
func (s *Service) Stop() {
	if !s.isRunning {
		return
	}

	s.stopChan <- struct{}{}
	s.isRunning = false
}

// cleanupRoutine runs the cleanup process at regular intervals
func (s *Service) cleanupRoutine() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.runCleanup(); err != nil {
				log.Printf("Cleanup error: %v", err)
			}
		case <-s.stopChan:
			log.Println("Cleanup service stopped")
			return
		}
	}
}

// runCleanup executes a single cleanup operation
func (s *Service) runCleanup() error {
	log.Println("Running cleanup operation...")
	return s.storage.CleanupExpired()
}
