package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheMetrics_Basic(t *testing.T) {
	metrics := NewCacheMetrics()
	
	// Record some operations
	metrics.RecordSave(10*time.Millisecond, 1024, nil)
	metrics.RecordSave(5*time.Millisecond, 512, nil)
	metrics.RecordSave(15*time.Millisecond, 2048, errors.New("save error"))
	
	metrics.RecordGet(2*time.Millisecond, true, nil)  // hit
	metrics.RecordGet(3*time.Millisecond, false, nil) // miss
	metrics.RecordGet(1*time.Millisecond, false, errors.New("get error"))
	
	stats := metrics.GetStats()
	
	// Verify counters
	if stats.Saves != 3 {
		t.Errorf("Expected 3 saves, got %d", stats.Saves)
	}
	
	if stats.SaveErrors != 1 {
		t.Errorf("Expected 1 save error, got %d", stats.SaveErrors)
	}
	
	if stats.Gets != 3 {
		t.Errorf("Expected 3 gets, got %d", stats.Gets)
	}
	
	if stats.GetErrors != 1 {
		t.Errorf("Expected 1 get error, got %d", stats.GetErrors)
	}
	
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	
	if stats.HitRate != 50.0 {
		t.Errorf("Expected 50%% hit rate, got %.1f%%", stats.HitRate)
	}
	
	// Verify sizes
	expectedBytes := int64(1024 + 512) // Error save not counted
	if stats.TotalBytes != expectedBytes {
		t.Errorf("Expected %d bytes, got %d", expectedBytes, stats.TotalBytes)
	}
}

func TestCacheMetrics_Timing(t *testing.T) {
	metrics := NewCacheMetrics()
	
	// Record operations with known durations
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}
	
	for _, d := range durations {
		metrics.RecordSave(d, 100, nil)
		metrics.RecordGet(d/2, true, nil)
	}
	
	stats := metrics.GetStats()
	
	// Average save time should be 20ms
	expectedAvgSave := 20 * time.Millisecond
	if stats.AvgSaveTime < 19*time.Millisecond || stats.AvgSaveTime > 21*time.Millisecond {
		t.Errorf("Expected avg save time ~%v, got %v", expectedAvgSave, stats.AvgSaveTime)
	}
	
	// Average get time should be 10ms
	expectedAvgGet := 10 * time.Millisecond
	if stats.AvgGetTime < 9*time.Millisecond || stats.AvgGetTime > 11*time.Millisecond {
		t.Errorf("Expected avg get time ~%v, got %v", expectedAvgGet, stats.AvgGetTime)
	}
}

func TestCacheMetrics_Reset(t *testing.T) {
	metrics := NewCacheMetrics()
	
	// Add some data
	metrics.RecordSave(10*time.Millisecond, 1024, nil)
	metrics.RecordGet(5*time.Millisecond, true, nil)
	
	// Verify data exists
	stats := metrics.GetStats()
	if stats.Saves == 0 || stats.Gets == 0 {
		t.Error("Expected non-zero metrics before reset")
	}
	
	// Reset
	metrics.Reset()
	
	// Verify all cleared
	stats = metrics.GetStats()
	if stats.Saves != 0 || stats.Gets != 0 || stats.Hits != 0 {
		t.Error("Expected zero metrics after reset")
	}
}

func TestMetricsCacheManager_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewMetricsCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	
	// Create test file
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	// Perform operations
	page := &MockPage{Title: "Test", Content: "Content"}
	
	// Save
	err = cache.SavePage(page, "test", testFile, nil)
	if err != nil {
		t.Error(err)
	}
	
	// Get (should be hit)
	_, hit, err := cache.GetPage("test", testFile)
	if err != nil || !hit {
		t.Error("Expected cache hit")
	}
	
	// Get non-existent (should be miss)
	_, hit, err = cache.GetPage("nonexistent", testFile)
	if err != nil || hit {
		t.Error("Expected cache miss")
	}
	
	// Check metrics
	stats := cache.GetMetrics()
	
	if stats.Saves != 1 {
		t.Errorf("Expected 1 save, got %d", stats.Saves)
	}
	
	if stats.Gets != 2 {
		t.Errorf("Expected 2 gets, got %d", stats.Gets)
	}
	
	if stats.HitRate != 50.0 {
		t.Errorf("Expected 50%% hit rate, got %.1f%%", stats.HitRate)
	}
	
	// Rates should be positive
	if stats.SavesPerSec <= 0 || stats.GetsPerSec <= 0 {
		t.Error("Expected positive operation rates")
	}
}

func TestCacheMetrics_Concurrent(t *testing.T) {
	metrics := NewCacheMetrics()
	
	// Run concurrent operations
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				metrics.RecordSave(time.Millisecond, 100, nil)
				metrics.RecordGet(time.Millisecond, true, nil)
			}
			done <- true
		}()
	}
	
	// Wait for completion
	for i := 0; i < 10; i++ {
		<-done
	}
	
	stats := metrics.GetStats()
	
	// Should have recorded all operations
	if stats.Saves != 1000 {
		t.Errorf("Expected 1000 saves, got %d", stats.Saves)
	}
	
	if stats.Gets != 1000 {
		t.Errorf("Expected 1000 gets, got %d", stats.Gets)
	}
}

func BenchmarkMetrics_RecordSave(b *testing.B) {
	metrics := NewCacheMetrics()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordSave(time.Millisecond, 1024, nil)
	}
}

func BenchmarkMetrics_RecordGet(b *testing.B) {
	metrics := NewCacheMetrics()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordGet(time.Millisecond, true, nil)
	}
}

func BenchmarkMetrics_GetStats(b *testing.B) {
	metrics := NewCacheMetrics()
	
	// Add some data
	for i := 0; i < 1000; i++ {
		metrics.RecordSave(time.Millisecond, 1024, nil)
		metrics.RecordGet(time.Millisecond, true, nil)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metrics.GetStats()
	}
}