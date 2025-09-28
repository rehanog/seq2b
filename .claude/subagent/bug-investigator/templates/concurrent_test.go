package templates

import (
	"sync"
	"testing"
)

// ConcurrentTestTemplate demonstrates testing concurrent operations
// Use this template when testing:
// - Race conditions
// - Simultaneous edits
// - Thread-safe operations
// - Deadlock scenarios

func TestConcurrentOperation(t *testing.T) {
	// Setup
	const numGoroutines = 100
	const numOperations = 1000
	
	// Shared resource (replace with your actual resource)
	resource := &YourResource{}
	
	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	
	// Channel to collect errors
	errChan := make(chan error, numGoroutines)
	
	// Launch concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()
			
			for j := 0; j < numOperations; j++ {
				// Perform operation (replace with actual operation)
				if err := resource.Operation(workerID, j); err != nil {
					errChan <- err
					return
				}
			}
		}(i)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)
	
	// Check for errors
	for err := range errChan {
		t.Errorf("Concurrent operation failed: %v", err)
	}
	
	// Verify final state
	if !resource.IsValid() {
		t.Error("Resource in invalid state after concurrent operations")
	}
}

// TestRaceCondition uses Go's race detector
// Run with: go test -race
func TestRaceCondition(t *testing.T) {
	shared := 0
	var wg sync.WaitGroup
	
	// This will be caught by race detector if not properly synchronized
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Intentional race - should use atomic or mutex
			shared++
		}()
	}
	
	wg.Wait()
	
	if shared != 100 {
		t.Errorf("Expected 100, got %d (race condition detected)", shared)
	}
}

// TestDeadlockTimeout ensures operations don't deadlock
func TestDeadlockTimeout(t *testing.T) {
	done := make(chan bool)
	
	go func() {
		// Your potentially deadlocking operation
		performOperation()
		done <- true
	}()
	
	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("Operation timed out - possible deadlock")
	}
}

// Placeholder types - replace with actual implementation
type YourResource struct {
	mu    sync.Mutex
	data  map[string]interface{}
}

func (r *YourResource) Operation(workerID, operationID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Implement operation
	return nil
}

func (r *YourResource) IsValid() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Implement validation
	return true
}

func performOperation() {
	// Implement potentially blocking operation
}