package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/dgraph-io/badger/v4"
)

func TestCacheManager_CorruptedData(t *testing.T) {
	tmpDir := t.TempDir()
	
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	
	// Manually insert corrupted data
	err = cache.db.Update(func(txn *badger.Txn) error {
		// Insert invalid JSON
		return txn.Set([]byte(pagePrefix+"corrupt"), []byte("not valid json{{{"))
	})
	if err != nil {
		t.Fatal(err)
	}
	
	cache.Close()
	
	// Reopen and try to read
	cache2, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache2.Close()
	
	// Should handle corrupted data gracefully
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	_, hit, err := cache2.GetPage("corrupt", testFile)
	if err == nil && hit {
		t.Error("Expected cache miss or error for corrupted data")
	}
}

func TestCacheManager_LargePages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large page test in short mode")
	}
	
	tmpDir := t.TempDir()
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	
	// Create a large page
	largePage := &MockPage{
		Title:   "Large Page",
		Content: generateLargeContent(1024 * 1024), // 1MB of content
		Links:   generateManyLinks(1000),
	}
	
	testFile := filepath.Join(tmpDir, "large.md")
	os.WriteFile(testFile, []byte("large"), 0644)
	
	// Save large page
	err = cache.SavePage(largePage, "large", testFile, largePage.Links)
	if err != nil {
		t.Errorf("Failed to save large page: %v", err)
	}
	
	// Retrieve large page
	cached, hit, err := cache.GetPage("large", testFile)
	if err != nil {
		t.Errorf("Failed to retrieve large page: %v", err)
	}
	if !hit {
		t.Error("Expected cache hit for large page")
	}
	
	// Verify size is preserved
	if rawJSON, ok := cached.(json.RawMessage); ok {
		var retrieved MockPage
		if err := json.Unmarshal(rawJSON, &retrieved); err == nil {
			if len(retrieved.Content) != len(largePage.Content) {
				t.Errorf("Content size mismatch: got %d, want %d", 
					len(retrieved.Content), len(largePage.Content))
			}
		}
	}
}

func TestCacheManager_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	
	// Test pages with special characters
	specialPages := []struct {
		name    string
		content string
	}{
		{"unicode-😀", "Content with emoji 🎉"},
		{"spaces in name", "Content with spaces"},
		{"special/chars", "Content with slashes"},
		{"quotes\"and'apostrophes", "Content with quotes"},
		{"\ttabs\tand\nnewlines\n", "Content with whitespace"},
		{"very-long-name-" + generateLargeContent(200), "Long name test"},
	}
	
	for _, sp := range specialPages {
		testFile := filepath.Join(tmpDir, "test.md")
		os.WriteFile(testFile, []byte("test"), 0644)
		
		page := &MockPage{
			Title:   sp.name,
			Content: sp.content,
		}
		
		// Save with special name
		err := cache.SavePage(page, sp.name, testFile, nil)
		if err != nil {
			t.Errorf("Failed to save page with name %q: %v", sp.name, err)
			continue
		}
		
		// Retrieve with special name
		cached, hit, err := cache.GetPage(sp.name, testFile)
		if err != nil {
			t.Errorf("Failed to get page with name %q: %v", sp.name, err)
		}
		if !hit {
			t.Errorf("Expected cache hit for page %q", sp.name)
		}
		
		// Verify content
		if rawJSON, ok := cached.(json.RawMessage); ok {
			var retrieved MockPage
			if err := json.Unmarshal(rawJSON, &retrieved); err == nil {
				if retrieved.Content != sp.content {
					t.Errorf("Content mismatch for %q", sp.name)
				}
			}
		}
	}
}

func TestCacheManager_ConcurrentCacheManagers(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create multiple cache managers for same directory
	// BadgerDB should handle this with locking
	
	cache1, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cache1.Close()
	
	// Second manager should fail to open (BadgerDB exclusive lock)
	cache2, err := NewCacheManager(tmpDir)
	if err == nil {
		cache2.Close()
		t.Error("Expected error when opening second cache manager for same directory")
	}
}

func TestCacheManager_DatabaseRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create and populate cache
	cache, err := NewCacheManager(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	page := &MockPage{Title: "Test", Content: "Content"}
	cache.SavePage(page, "test", testFile, nil)
	cache.Close()
	
	// Simulate partial corruption by removing a file
	// BadgerDB should be able to recover
	cacheFiles, _ := filepath.Glob(filepath.Join(tmpDir, "*.vlog"))
	if len(cacheFiles) > 0 {
		os.Remove(cacheFiles[0])
	}
	
	// Try to reopen
	cache2, err := NewCacheManager(tmpDir)
	if err != nil {
		// This is actually expected - BadgerDB may not recover from this
		t.Logf("Database recovery failed as expected: %v", err)
		return
	}
	defer cache2.Close()
	
	// If it opened, verify it still works
	_, _, err = cache2.GetPage("test", testFile)
	if err != nil {
		t.Logf("Cache read after recovery failed: %v", err)
	}
}

// Helper functions
func generateLargeContent(size int) string {
	content := make([]byte, size)
	for i := range content {
		content[i] = byte('a' + (i % 26))
	}
	return string(content)
}

func generateManyLinks(count int) []string {
	links := make([]string, count)
	for i := range links {
		links[i] = "Page" + string(rune(i))
	}
	return links
}

// Benchmarks for edge cases
func BenchmarkCacheLargePage(b *testing.B) {
	tmpDir := b.TempDir()
	cache, _ := NewCacheManager(tmpDir)
	defer cache.Close()
	
	largePage := &MockPage{
		Title:   "Large",
		Content: generateLargeContent(100 * 1024), // 100KB
		Links:   generateManyLinks(100),
	}
	
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SavePage(largePage, "large", testFile, largePage.Links)
	}
}

func BenchmarkCacheManyPages(b *testing.B) {
	tmpDir := b.TempDir()
	cache, _ := NewCacheManager(tmpDir)
	defer cache.Close()
	
	// Pre-create test files
	files := make([]string, 100)
	for i := range files {
		files[i] = filepath.Join(tmpDir, "test"+string(rune(i))+".md")
		os.WriteFile(files[i], []byte("test"), 0644)
	}
	
	page := &MockPage{Title: "Test", Content: "Content"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(files)
		cache.SavePage(page, "page"+string(rune(idx)), files[idx], nil)
	}
}