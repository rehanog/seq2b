package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPageNameCaseSensitivity verifies that page lookups work regardless of case
func TestPageNameCaseSensitivity(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "seq2b-case-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an existing page with uppercase title
	pageFile := filepath.Join(tempDir, "page-a.md")
	content := `# Page A

- This is Page A with capital letters
- It should be findable by [[page a]] link`

	if err := os.WriteFile(pageFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write page: %v", err)
	}

	// Create app and load directory
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}

	// Test 1: Should find "Page A" when looking for "page a"
	pageData, err := app.GetPage("page a")
	if err != nil {
		t.Errorf("Failed to get page with lowercase name: %v", err)
	} else if pageData.Title != "Page A" {
		t.Errorf("Expected title 'Page A', got '%s'", pageData.Title)
	}

	// Test 2: Should find "Page A" when looking for "PAGE A"
	pageData, err = app.GetPage("PAGE A")
	if err != nil {
		t.Errorf("Failed to get page with uppercase name: %v", err)
	} else if pageData.Title != "Page A" {
		t.Errorf("Expected title 'Page A', got '%s'", pageData.Title)
	}

	// Test 3: Should find "Page A" when looking for "Page A" (exact match)
	pageData, err = app.GetPage("Page A")
	if err != nil {
		t.Errorf("Failed to get page with exact name: %v", err)
	} else if pageData.Title != "Page A" {
		t.Errorf("Expected title 'Page A', got '%s'", pageData.Title)
	}
}

// TestCreatePageWithExistingFile tests creating a page when file already exists with different case
func TestCreatePageWithExistingFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-existing-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an existing file with uppercase title
	pageFile := filepath.Join(tempDir, "test-page.md")
	existingContent := `# Test Page

- Existing content with uppercase title`

	if err := os.WriteFile(pageFile, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing page: %v", err)
	}

	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}

	// Try to get "test page" (lowercase) - should find existing "Test Page"
	pageData, err := app.GetPage("test page")
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}

	// Should preserve the original title case
	if pageData.Title != "Test Page" {
		t.Errorf("Expected original title 'Test Page', got '%s'", pageData.Title)
	}

	// Content should be unchanged
	if len(pageData.Blocks) == 0 || !strings.Contains(pageData.Blocks[0].Content, "Existing content") {
		t.Error("Existing content was not preserved")
	}

	// File should not be overwritten
	savedContent, err := os.ReadFile(pageFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if !strings.Contains(string(savedContent), "# Test Page") {
		t.Error("File was overwritten with different case")
	}
}

// TestBacklinksCaseInsensitive verifies backlinks work with different cases
func TestBacklinksCaseInsensitive(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-backlink-case-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source page with mixed case links
	sourceFile := filepath.Join(tempDir, "source.md")
	sourceContent := `# Source

- Link to [[Page Target]] with mixed case
- Link to [[page target]] with lowercase
- Link to [[PAGE TARGET]] with uppercase`

	targetFile := filepath.Join(tempDir, "page-target.md")
	targetContent := `# Page Target

- This is the target page`

	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("Failed to write source: %v", err)
	}
	if err := os.WriteFile(targetFile, []byte(targetContent), 0644); err != nil {
		t.Fatalf("Failed to write target: %v", err)
	}

	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}

	// Should find backlinks using any case
	backlinks := app.GetBacklinks("Page Target")
	if len(backlinks) == 0 || len(backlinks["Source"]) == 0 {
		t.Errorf("Expected backlinks from Source page, got: %v", backlinks)
	}

	// Should also find backlinks with different case queries
	backlinks = app.GetBacklinks("page target")
	if len(backlinks) == 0 || len(backlinks["Source"]) == 0 {
		t.Error("Backlinks lookup should be case-insensitive")
	}
	
	// Should work with uppercase too
	backlinks = app.GetBacklinks("PAGE TARGET")
	if len(backlinks) == 0 || len(backlinks["Source"]) == 0 {
		t.Error("Backlinks lookup should work with uppercase")
	}
}