package cleanup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of StorageInterface
type MockStorage struct {
	mock.Mock
	cleanupCalled bool
}

func (m *MockStorage) CleanupExpired() error {
	m.cleanupCalled = true
	args := m.Called()
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	storage := new(MockStorage)
	interval := 5 * time.Minute

	service := NewService(storage, interval)
	assert.NotNil(t, service)
	assert.Equal(t, interval, service.interval)
	assert.False(t, service.isRunning)
}

func TestService_StartStop(t *testing.T) {
	storage := new(MockStorage)
	interval := 100 * time.Millisecond
	service := NewService(storage, interval)

	// Setup mock expectations
	storage.On("CleanupExpired").Return(nil)

	// Test Start
	service.Start()
	assert.True(t, service.isRunning)

	// Wait for at least one cleanup cycle
	time.Sleep(interval * 2)

	// Test Stop
	service.Stop()
	assert.False(t, service.isRunning)

	// Verify cleanup was called
	assert.True(t, storage.cleanupCalled)
	storage.AssertExpectations(t)
}

func TestService_MultipleStartStop(t *testing.T) {
	storage := new(MockStorage)
	service := NewService(storage, time.Minute)

	// Setup mock expectations
	storage.On("CleanupExpired").Return(nil)

	// Test multiple starts
	service.Start()
	assert.True(t, service.isRunning)
	service.Start() // Should not affect isRunning state
	assert.True(t, service.isRunning)

	// Test multiple stops
	service.Stop()
	assert.False(t, service.isRunning)
	service.Stop() // Should not affect isRunning state
	assert.False(t, service.isRunning)
}

func TestService_CleanupError(t *testing.T) {
	storage := new(MockStorage)
	interval := 100 * time.Millisecond
	service := NewService(storage, interval)

	// Setup mock to return an error
	storage.On("CleanupExpired").Return(error(nil))

	// Start service
	service.Start()

	// Wait for at least one cleanup cycle
	time.Sleep(interval * 2)

	// Stop service
	service.Stop()

	// Verify cleanup was still called despite error
	assert.True(t, storage.cleanupCalled)
	storage.AssertExpectations(t)
}
