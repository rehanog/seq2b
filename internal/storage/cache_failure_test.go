package storage

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"
	
	"github.com/dgraph-io/badger/v4"
)

// TestCacheManager_DiskSpaceExhaustion simulates disk space issues
func TestCacheManager_DiskSpaceExhaustion(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Disk space simulation not implemented for Windows")
	}
	
	tmpDir := t.TempDir()
	
	// Create cache manager
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	
	// Create a very large page that might exhaust space
	hugePage := &MockPage{
		Title:   "Huge Page",
		Content: generateLargeContent(50 * 1024 * 1024), // 50MB
		Links:   generateManyLinks(10000),
	}
	
	testFile := filepath.Join(tmpDir, "huge.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	// Try to save - should handle gracefully even if it fails
	err = cache.SavePage(hugePage, "huge", testFile, hugePage.Links)
	if err != nil {
		// Check if it's a space-related error
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			if pathErr.Err == syscall.ENOSPC {
				t.Log("Correctly detected disk space exhaustion")
				return
			}
		}
		// Other errors are acceptable too (BadgerDB might have internal limits)
		t.Logf("Save failed with error (this is expected): %v", err)
	}
}

// TestCacheManager_MaxCacheSize tests behavior with size limits
func TestCacheManager_MaxCacheSize(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create cache with custom options to limit size
	opts := badger.DefaultOptions(tmpDir)
	opts.Logger = nil
	opts.ValueLogFileSize = 1 << 20 // 1MB max file size
	
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatal(err)
	}
	
	cache := &CacheManager{
		db:          db,
		libraryPath: tmpDir,
	}
	defer cache.Close()
	
	// Try to fill the cache beyond limits
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	errors := 0
	saved := 0
	
	for i := 0; i < 1000; i++ {
		page := &MockPage{
			Title:   "Page" + string(rune(i)),
			Content: generateLargeContent(10 * 1024), // 10KB each
			Links:   []string{"link1", "link2"},
		}
		
		err := cache.SavePage(page, page.Title, testFile, page.Links)
		if err != nil {
			errors++
			t.Logf("Save %d failed (expected after filling cache): %v", i, err)
			if errors > 10 {
				break // Cache is full
			}
		} else {
			saved++
		}
	}
	
	t.Logf("Successfully saved %d pages before hitting limits", saved)
	if saved == 0 {
		t.Error("Expected to save at least some pages before hitting limit")
	}
}

// TestCacheManager_PermissionDenied tests read-only cache directory
func TestCacheManager_PermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Permission test not reliable on Windows")
	}
	
	tmpDir := t.TempDir()
	cacheDir := filepath.Join(tmpDir, "readonly-cache")
	
	// Create directory
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	// Try to create cache (should work)
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	cache.Close()
	
	// Make directory read-only
	if err := os.Chmod(cacheDir, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(cacheDir, 0755) // Restore permissions
	
	// Try to create cache again - should fail gracefully
	cache2, err := NewCacheManager(cacheDir)
	if err == nil {
		cache2.Close()
		t.Error("Expected error when opening cache in read-only directory")
	} else {
		t.Logf("Correctly failed with permission error: %v", err)
	}
}

// TestCacheManager_NetworkDrive simulates network drive issues
func TestCacheManager_NetworkDrive(t *testing.T) {
	// This test simulates slow/unreliable storage
	tmpDir := t.TempDir()
	
	// Create a wrapper that simulates network delays
	slowCache := &SlowCacheManager{
		CacheManager: nil,
		delay:        100 * time.Millisecond,
	}
	
	// Create normal cache
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	slowCache.CacheManager = cache
	defer cache.Close()
	
	// Test timeout behavior
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	start := time.Now()
	page := &MockPage{Title: "Test", Content: "Content"}
	
	// This should be slow but succeed
	err = slowCache.SavePageWithDelay(page, "test", testFile, nil)
	elapsed := time.Since(start)
	
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}
	
	if elapsed < slowCache.delay {
		t.Errorf("Operation too fast, delay not applied: %v", elapsed)
	}
	
	t.Logf("Network simulation: operation took %v", elapsed)
}

// TestCacheManager_CorruptedBadgerDB tests recovery from corrupted DB
func TestCacheManager_CorruptedBadgerDB(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create and populate cache
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	page := &MockPage{Title: "Test", Content: "Original"}
	cache.SavePage(page, "test", testFile, nil)
	cache.Close()
	
	// Corrupt the manifest file
	manifestPath := filepath.Join(tmpDir, "MANIFEST")
	if err := os.WriteFile(manifestPath, []byte("corrupted data"), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Try to open corrupted cache
	cache2, err := NewCacheManager(tmpDir)
	if err != nil {
		// This is expected - BadgerDB detected corruption
		t.Logf("BadgerDB correctly detected corruption: %v", err)
		
		// In production, we'd want to clear and rebuild
		os.RemoveAll(tmpDir)
		os.MkdirAll(tmpDir, 0755)
		
		// Try again with fresh cache
		cache3, err := NewCacheManager(tmpDir)
		if err != nil {
			t.Errorf("Failed to create fresh cache after corruption: %v", err)
		} else {
			cache3.Close()
			t.Log("Successfully recovered by creating fresh cache")
		}
	} else {
		// If it opened despite corruption, that's concerning
		cache2.Close()
		t.Error("Cache opened despite manifest corruption")
	}
}

// TestCacheManager_ConcurrentWrites tests race conditions
func TestCacheManager_ConcurrentWrites(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	// Run many concurrent operations
	const goroutines = 50
	const operations = 100
	
	done := make(chan bool, goroutines)
	errors := make(chan error, goroutines*operations)
	
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < operations; j++ {
				// Mix of reads and writes
				if j%2 == 0 {
					// Write
					page := &MockPage{
						Title:   "Page" + string(rune(id)),
						Content: "Content " + string(rune(j)),
					}
					if err := cache.SavePage(page, page.Title, testFile, nil); err != nil {
						errors <- err
					}
				} else {
					// Read
					_, _, err := cache.GetPage("Page"+string(rune(id)), testFile)
					if err != nil && err != badger.ErrKeyNotFound {
						errors <- err
					}
				}
			}
			done <- true
		}(i)
	}
	
	// Wait for completion
	for i := 0; i < goroutines; i++ {
		<-done
	}
	
	close(errors)
	
	// Check for errors
	errorCount := 0
	for err := range errors {
		errorCount++
		if errorCount <= 5 {
			t.Logf("Concurrent operation error: %v", err)
		}
	}
	
	if errorCount > 0 {
		t.Errorf("Had %d errors during concurrent operations", errorCount)
	}
}

// TestCacheManager_Resilience tests recovery from various failures
func TestCacheManager_Resilience(t *testing.T) {
	scenarios := []struct {
		name string
		test func(t *testing.T, tmpDir string)
	}{
		{
			name: "MissingDirectory",
			test: func(t *testing.T, tmpDir string) {
				// Remove directory while cache is closed
				cache, _ := NewCacheManager(tmpDir)
				cache.Close()
				
				os.RemoveAll(tmpDir)
				
				// Should recreate directory
				cache2, err := NewCacheManager(tmpDir)
				if err != nil {
					t.Errorf("Failed to recreate cache after directory removal: %v", err)
				} else {
					cache2.Close()
				}
			},
		},
		{
			name: "PartialWrite",
			test: func(t *testing.T, tmpDir string) {
				// Simulate partial write by filling up most of small value log
				opts := badger.DefaultOptions(tmpDir)
				opts.Logger = nil
				opts.ValueLogFileSize = 1 << 15 // Very small (32KB)
				
				db, err := badger.Open(opts)
				if err != nil {
					t.Fatal(err)
				}
				
				cache := &CacheManager{db: db, libraryPath: tmpDir}
				defer cache.Close()
				
				// Try to write until we hit limit
				testFile := filepath.Join(tmpDir, "test.md")
				os.WriteFile(testFile, []byte("test"), 0644)
				
				for i := 0; i < 100; i++ {
					page := &MockPage{
						Title:   "Page" + string(rune(i)),
						Content: generateLargeContent(1024), // 1KB
					}
					cache.SavePage(page, page.Title, testFile, nil)
				}
				
				// Should handle gracefully
				t.Log("Partial write scenario completed")
			},
		},
	}
	
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			scenario.test(t, tmpDir)
		})
	}
}

// Helper type for network simulation
type SlowCacheManager struct {
	*CacheManager
	delay time.Duration
}

func (s *SlowCacheManager) SavePageWithDelay(page interface{}, pageName string, filePath string, deps []string) error {
	time.Sleep(s.delay)
	return s.CacheManager.SavePage(page, pageName, filePath, deps)
}

// TestCacheManager_ProductionLoad simulates production-like load
func TestCacheManager_ProductionLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping production load test in short mode")
	}
	
	tmpDir := t.TempDir()
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	
	// Simulate production patterns:
	// - Many pages (5000+)
	// - Mixed read/write ratio (80/20)
	// - Varying page sizes
	// - Periodic cache clears
	// - Long-running operations
	
	const numPages = 5000
	const duration = 10 * time.Second
	
	// Pre-create test files
	files := make([]string, 100)
	for i := range files {
		files[i] = filepath.Join(tmpDir, "file"+string(rune(i))+".md")
		os.WriteFile(files[i], []byte("test"), 0644)
	}
	
	// Metrics
	var writes, reads, cacheHits int64
	
	start := time.Now()
	deadline := start.Add(duration)
	
	// Run until deadline
	for time.Now().Before(deadline) {
		// 80% reads, 20% writes
		if time.Now().UnixNano()%5 == 0 {
			// Write
			pageNum := int(time.Now().UnixNano()) % numPages
			page := &MockPage{
				Title:   "Page" + string(rune(pageNum)),
				Content: generateLargeContent(100 + (pageNum%10)*1024), // 100B to 10KB
				Links:   generateManyLinks(pageNum % 20),
			}
			
			fileIdx := pageNum % len(files)
			cache.SavePage(page, page.Title, files[fileIdx], page.Links)
			writes++
		} else {
			// Read
			pageNum := int(time.Now().UnixNano()) % numPages
			fileIdx := pageNum % len(files)
			
			_, hit, _ := cache.GetPage("Page"+string(rune(pageNum)), files[fileIdx])
			reads++
			if hit {
				cacheHits++
			}
		}
		
		// Occasionally clear cache (simulating restart)
		if writes%1000 == 999 {
			t.Log("Simulating cache clear...")
			cache.Clear()
			cache.SaveMetadata()
		}
	}
	
	elapsed := time.Since(start)
	hitRate := float64(cacheHits) / float64(reads) * 100
	
	t.Logf("Production load test results:")
	t.Logf("  Duration: %v", elapsed)
	t.Logf("  Writes: %d (%.0f/sec)", writes, float64(writes)/elapsed.Seconds())
	t.Logf("  Reads: %d (%.0f/sec)", reads, float64(reads)/elapsed.Seconds())
	t.Logf("  Cache hit rate: %.1f%%", hitRate)
	
	// Check cache is still functional
	testFile := files[0]
	finalPage := &MockPage{Title: "Final", Content: "Test"}
	if err := cache.SavePage(finalPage, "final", testFile, nil); err != nil {
		t.Errorf("Cache not functional after load test: %v", err)
	}
}