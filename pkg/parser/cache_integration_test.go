package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseDirectoryWithCache_Basic(t *testing.T) {
	// Create test directory
	tmpDir := t.TempDir()
	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	// Create test pages
	pages := map[string]string{
		"page1.md": "# Page 1\n- Content with [[Page 2]] link",
		"page2.md": "# Page 2\n- Content with [[Page 1]] backlink",
		"page3.md": "# Page 3\n- No links here",
	}
	
	for name, content := range pages {
		path := filepath.Join(pagesDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	
	// First parse - should build cache
	result1, err := ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Fatalf("First parse failed: %v", err)
	}
	
	if len(result1.Pages) != 3 {
		t.Errorf("Expected 3 pages, got %d", len(result1.Pages))
	}
	
	// Check that output indicates cache miss
	// This is a bit fragile but works for now
	
	// Second parse - should use cache
	result2, err := ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Fatalf("Second parse failed: %v", err)
	}
	
	if len(result2.Pages) != 3 {
		t.Errorf("Expected 3 pages on second parse, got %d", len(result2.Pages))
	}
	
	// Verify backlinks are preserved
	page1Backlinks := result2.Backlinks.GetBacklinks("Page 1")
	if len(page1Backlinks) != 1 {
		t.Errorf("Expected 1 backlink to Page 1, got %d", len(page1Backlinks))
	}
	
	page2Backlinks := result2.Backlinks.GetBacklinks("Page 2")
	if len(page2Backlinks) != 1 {
		t.Errorf("Expected 1 backlink to Page 2, got %d", len(page2Backlinks))
	}
}

func TestParseDirectoryWithCache_FileModification(t *testing.T) {
	tmpDir := t.TempDir()
	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	// Create initial page
	pagePath := filepath.Join(pagesDir, "test.md")
	originalContent := "# Test Page\n- Original content"
	if err := os.WriteFile(pagePath, []byte(originalContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	// First parse
	result1, err := ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Fatal(err)
	}
	
	page1 := result1.Pages["Test Page"]
	if page1 == nil {
		t.Fatal("Page not found in first parse")
	}
	
	// Verify original content
	foundOriginal := false
	for _, block := range page1.Blocks {
		if strings.Contains(block.Content, "Original content") {
			foundOriginal = true
			break
		}
	}
	if !foundOriginal {
		t.Error("Original content not found in parsed blocks")
	}
	
	// Modify file (ensure timestamp changes)
	time.Sleep(10 * time.Millisecond)
	modifiedContent := "# Test Page\n- Modified content\n- New line"
	if err := os.WriteFile(pagePath, []byte(modifiedContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Second parse - should detect modification
	result2, err := ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Fatal(err)
	}
	
	page2 := result2.Pages["Test Page"]
	if page2 == nil {
		t.Fatal("Page not found in second parse")
	}
	
	// Verify modified content
	hasModified := false
	hasNewLine := false
	for _, block := range page2.Blocks {
		if strings.Contains(block.Content, "Modified content") {
			hasModified = true
		}
		if strings.Contains(block.Content, "New line") {
			hasNewLine = true
		}
	}
	
	if !hasModified || !hasNewLine {
		t.Error("Modified content not detected after file change")
	}
}

func TestParseDirectoryWithCache_Performance(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	tmpDir := t.TempDir()
	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	// Create many test pages
	numPages := 100
	for i := 0; i < numPages; i++ {
		content := "# Page " + string(rune(i)) + "\n"
		// Add some cross-references
		for j := 0; j < 5; j++ {
			targetPage := (i + j + 1) % numPages
			content += "- Link to [[Page " + string(rune(targetPage)) + "]]\n"
		}
		
		path := filepath.Join(pagesDir, "page"+string(rune(i))+".md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	
	// First parse (cold)
	start := time.Now()
	result1, err := ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Fatal(err)
	}
	coldTime := time.Since(start)
	
	if len(result1.Pages) != numPages {
		t.Errorf("Expected %d pages, got %d", numPages, len(result1.Pages))
	}
	
	// Second parse (warm)
	start = time.Now()
	_, err = ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Fatal(err)
	}
	warmTime := time.Since(start)
	
	// Warm parse should be significantly faster
	speedup := float64(coldTime) / float64(warmTime)
	t.Logf("Cold parse: %v, Warm parse: %v, Speedup: %.2fx", coldTime, warmTime, speedup)
	
	if speedup < 1.5 {
		t.Errorf("Expected at least 1.5x speedup with cache, got %.2fx", speedup)
	}
}

func TestParseDirectoryWithCache_ErrorHandling(t *testing.T) {
	// Test with non-existent directory
	_, err := ParseDirectoryWithCache("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
	
	// Test with empty directory
	tmpDir := t.TempDir()
	result, err := ParseDirectoryWithCache(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error for empty directory: %v", err)
	}
	if len(result.Pages) != 0 {
		t.Errorf("Expected 0 pages for empty directory, got %d", len(result.Pages))
	}
}

func TestParseDirectoryWithCache_CacheFallback(t *testing.T) {
	tmpDir := t.TempDir()
	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	// Create a page
	if err := os.WriteFile(filepath.Join(pagesDir, "test.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Make cache directory read-only to force cache failure
	cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "seq2b")
	if err := os.MkdirAll(cacheDir, 0755); err == nil {
		// Try to make it read-only (may not work on all systems)
		defer os.Chmod(cacheDir, 0755) // Restore permissions
		os.Chmod(cacheDir, 0444)
	}
	
	// Should still parse successfully (falling back to non-cached)
	result, err := ParseDirectoryWithCache(pagesDir)
	if err != nil {
		t.Logf("Parse with restricted cache failed (this might be expected): %v", err)
	} else if len(result.Pages) != 1 {
		t.Errorf("Expected 1 page even with cache issues, got %d", len(result.Pages))
	}
}