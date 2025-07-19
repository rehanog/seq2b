package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Mock page structure for testing
type MockPage struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Links   []string `json:"links"`
}

func TestCacheManager_SaveAndGetPage(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	
	// Create cache manager
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cache.Close()
	
	// Create test page
	page := &MockPage{
		Title:   "Test Page",
		Content: "This is test content",
		Links:   []string{"Page A", "Page B"},
	}
	
	// Create a temporary file to simulate page file
	testFile := filepath.Join(tmpDir, "test-page.md")
	if err := os.WriteFile(testFile, []byte("# Test Page"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Test saving page
	deps := []string{"Page A", "Page B"}
	if err := cache.SavePage(page, "test-page", testFile, deps); err != nil {
		t.Errorf("SavePage failed: %v", err)
	}
	
	// Test retrieving page
	cachedData, hit, err := cache.GetPage("test-page", testFile)
	if err != nil {
		t.Errorf("GetPage failed: %v", err)
	}
	if !hit {
		t.Error("Expected cache hit, got miss")
	}
	
	// Unmarshal and verify
	if rawJSON, ok := cachedData.(json.RawMessage); ok {
		var retrieved MockPage
		if err := json.Unmarshal(rawJSON, &retrieved); err != nil {
			t.Errorf("Failed to unmarshal cached page: %v", err)
		}
		
		if retrieved.Title != page.Title {
			t.Errorf("Title mismatch: got %s, want %s", retrieved.Title, page.Title)
		}
		if retrieved.Content != page.Content {
			t.Errorf("Content mismatch: got %s, want %s", retrieved.Content, page.Content)
		}
		if len(retrieved.Links) != len(page.Links) {
			t.Errorf("Links count mismatch: got %d, want %d", len(retrieved.Links), len(page.Links))
		}
	} else {
		t.Errorf("Unexpected type returned: %T", cachedData)
	}
}

func TestCacheManager_FileModificationCheck(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cache.Close()
	
	// Create test file
	testFile := filepath.Join(tmpDir, "test-page.md")
	if err := os.WriteFile(testFile, []byte("# Original"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Save page to cache
	page := &MockPage{Title: "Test", Content: "Original content"}
	if err := cache.SavePage(page, "test-page", testFile, nil); err != nil {
		t.Fatalf("SavePage failed: %v", err)
	}
	
	// Should get cache hit
	_, hit, err := cache.GetPage("test-page", testFile)
	if err != nil {
		t.Errorf("GetPage failed: %v", err)
	}
	if !hit {
		t.Error("Expected cache hit for unmodified file")
	}
	
	// Modify file (need to ensure timestamp changes)
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(testFile, []byte("# Modified"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	
	// Should get cache miss due to modification
	_, hit, err = cache.GetPage("test-page", testFile)
	if err != nil {
		t.Errorf("GetPage failed: %v", err)
	}
	if hit {
		t.Error("Expected cache miss for modified file")
	}
}

func TestCacheManager_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cache.Close()
	
	// Try to get page for non-existent file
	nonExistentFile := filepath.Join(tmpDir, "does-not-exist.md")
	_, hit, err := cache.GetPage("test-page", nonExistentFile)
	
	if err != nil {
		t.Errorf("GetPage should not error for non-existent file: %v", err)
	}
	if hit {
		t.Error("Expected cache miss for non-existent file")
	}
}

func TestCacheManager_ValidateCache(t *testing.T) {
	tmpDir := t.TempDir()
	
	// First cache manager
	cache1, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	
	// Save metadata
	if err := cache1.SaveMetadata(); err != nil {
		t.Fatalf("SaveMetadata failed: %v", err)
	}
	
	// Validate should succeed
	valid, err := cache1.ValidateCache()
	if err != nil {
		t.Errorf("ValidateCache failed: %v", err)
	}
	if !valid {
		t.Error("Expected cache to be valid")
	}
	cache1.Close()
	
	// Create new cache manager with different path
	cache2, err := NewCacheManager(tmpDir + "/different")
	if err != nil {
		t.Fatalf("Failed to create second cache manager: %v", err)
	}
	defer cache2.Close()
	
	// Validation should fail due to path mismatch
	valid, err = cache2.ValidateCache()
	if err != nil {
		t.Errorf("ValidateCache failed: %v", err)
	}
	if valid {
		t.Error("Expected cache to be invalid due to path mismatch")
	}
}

func TestCacheManager_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cache.Close()
	
	// Save some data
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	page := &MockPage{Title: "Test"}
	if err := cache.SavePage(page, "test", testFile, nil); err != nil {
		t.Fatalf("SavePage failed: %v", err)
	}
	
	// Verify data exists
	_, hit, _ := cache.GetPage("test", testFile)
	if !hit {
		t.Error("Expected cache hit before clear")
	}
	
	// Clear cache
	if err := cache.Clear(); err != nil {
		t.Errorf("Clear failed: %v", err)
	}
	
	// Verify data is gone
	_, hit, _ = cache.GetPage("test", testFile)
	if hit {
		t.Error("Expected cache miss after clear")
	}
}

func TestCacheManager_Backlinks(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cache.Close()
	
	// Mock backlinks data
	type MockBlockRef struct {
		PageName string `json:"page_name"`
		BlockID  string `json:"block_id"`
		Position int    `json:"position"`
	}
	
	backlinks := map[string][]MockBlockRef{
		"Page A": {
			{PageName: "Page B", BlockID: "123", Position: 10},
			{PageName: "Page C", BlockID: "456", Position: 20},
		},
	}
	
	// Save backlinks
	if err := cache.SaveBacklinks("Test Page", backlinks); err != nil {
		t.Errorf("SaveBacklinks failed: %v", err)
	}
	
	// Retrieve backlinks
	cached, hit, err := cache.GetBacklinks("Test Page")
	if err != nil {
		t.Errorf("GetBacklinks failed: %v", err)
	}
	if !hit {
		t.Error("Expected backlinks cache hit")
	}
	
	// Verify data
	if cached != nil {
		// The actual verification would depend on the structure
		t.Log("Successfully retrieved backlinks")
	}
}

func TestCacheManager_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cache.Close()
	
	// Create test file
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Run concurrent operations
	done := make(chan bool, 3)
	
	// Writer 1
	go func() {
		for i := 0; i < 10; i++ {
			page := &MockPage{Title: "Test", Content: "Content " + string(rune(i))}
			cache.SavePage(page, "test", testFile, nil)
		}
		done <- true
	}()
	
	// Writer 2
	go func() {
		for i := 0; i < 10; i++ {
			cache.SaveBacklinks("test", map[string]interface{}{"data": i})
		}
		done <- true
	}()
	
	// Reader
	go func() {
		for i := 0; i < 20; i++ {
			cache.GetPage("test", testFile)
			cache.GetBacklinks("test")
		}
		done <- true
	}()
	
	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// If we get here without panic/deadlock, concurrent access is working
	t.Log("Concurrent access test passed")
}

func TestGetCacheDir(t *testing.T) {
	// Test cache directory creation
	dir, err := getCacheDir()
	if err != nil {
		t.Errorf("getCacheDir failed: %v", err)
	}
	
	// Verify directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("Cache directory was not created")
	}
	
	// Should contain "seq2b" in path
	if filepath.Base(dir) != "seq2b" {
		t.Errorf("Cache directory should end with 'seq2b', got: %s", dir)
	}
}

// Benchmark cache operations
func BenchmarkCacheSave(b *testing.B) {
	tmpDir := b.TempDir()
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		b.Fatal(err)
	}
	defer cache.Close()
	
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	page := &MockPage{
		Title:   "Test Page",
		Content: "This is a longer content string that simulates real page content with multiple sentences and paragraphs.",
		Links:   []string{"Page A", "Page B", "Page C", "Page D", "Page E"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SavePage(page, "test", testFile, page.Links)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	tmpDir := b.TempDir()
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		b.Fatal(err)
	}
	defer cache.Close()
	
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	page := &MockPage{Title: "Test Page", Content: "Content"}
	cache.SavePage(page, "test", testFile, nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetPage("test", testFile)
	}
}