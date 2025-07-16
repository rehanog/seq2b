package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestEditPersistenceAcrossNavigation verifies that edits are saved before navigation
func TestEditPersistenceAcrossNavigation(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "seq2b-nav-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create initial page
	pageAFile := filepath.Join(tempDir, "page-a.md")
	initialContent := `# Page A

- First block
- Second block
- Third block`

	if err := os.WriteFile(pageAFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write Page A: %v", err)
	}

	// Create app and load directory
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}

	// Get Page A (to ensure it's loaded)
	_, err = app.GetPage("Page A")
	if err != nil {
		t.Fatalf("Failed to get Page A: %v", err)
	}

	// Simulate editing the second block to add a link
	blockPath := BlockPath{1} // Second block (index 1)
	newContent := "Second block with [[Page B]] link"
	
	delta, err := app.UpdateBlockAtPath("Page A", blockPath, newContent)
	if err != nil {
		t.Fatalf("Failed to update block: %v", err)
	}

	// Verify the delta shows the update
	if delta["action"] != "update" {
		t.Errorf("Expected action 'update', got %v", delta["action"])
	}

	// Simulate navigation to Page B (which will auto-create it)
	_, err = app.GetPage("Page B")
	if err != nil {
		t.Fatalf("Failed to navigate to Page B: %v", err)
	}

	// Verify Page B was created
	pageBFile := filepath.Join(tempDir, "page-b.md")
	if _, err := os.Stat(pageBFile); os.IsNotExist(err) {
		t.Error("Page B file was not created")
	}

	// Now simulate going back to Page A
	pageDataAfterNav, err := app.GetPage("Page A")
	if err != nil {
		t.Fatalf("Failed to get Page A after navigation: %v", err)
	}

	// Check that the edit persisted
	if len(pageDataAfterNav.Blocks) < 2 {
		t.Fatalf("Expected at least 2 blocks, got %d", len(pageDataAfterNav.Blocks))
	}

	secondBlock := pageDataAfterNav.Blocks[1]
	if !strings.Contains(secondBlock.Content, "[[Page B]]") {
		t.Errorf("Edit was lost! Expected block to contain '[[Page B]]', got: %s", secondBlock.Content)
	}

	// Also verify the file on disk has the changes
	savedContent, err := os.ReadFile(pageAFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !strings.Contains(string(savedContent), "[[Page B]]") {
		t.Errorf("Edit not saved to disk! File content:\n%s", string(savedContent))
	}
}

// TestRapidNavigationWithEdits tests quick navigation while editing
func TestRapidNavigationWithEdits(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-rapid-nav-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test pages
	for _, pageName := range []string{"Alpha", "Beta", "Gamma"} {
		file := filepath.Join(tempDir, strings.ToLower(pageName)+".md")
		content := "# " + pageName + "\n\n- Content of " + pageName
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", pageName, err)
		}
	}

	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}

	// Simulate rapid edits and navigation
	edits := []struct {
		page    string
		path    BlockPath
		content string
	}{
		{"Alpha", BlockPath{0}, "Content of Alpha with [[Beta]] link"},
		{"Beta", BlockPath{0}, "Content of Beta with [[Gamma]] link"},
		{"Gamma", BlockPath{0}, "Content of Gamma with [[Alpha]] link"},
	}

	for _, edit := range edits {
		// Make edit
		if _, err := app.UpdateBlockAtPath(edit.page, edit.path, edit.content); err != nil {
			t.Errorf("Failed to update %s: %v", edit.page, err)
		}
		
		// Small delay to simulate user action
		time.Sleep(10 * time.Millisecond)
		
		// Navigate to next page
		_, err := app.GetPage(edits[(len(edits)+1)%len(edits)].page)
		if err != nil {
			t.Errorf("Failed to navigate: %v", err)
		}
	}

	// Verify all edits persisted
	for _, edit := range edits {
		pageData, err := app.GetPage(edit.page)
		if err != nil {
			t.Errorf("Failed to get %s: %v", edit.page, err)
			continue
		}

		if len(pageData.Blocks) == 0 {
			t.Errorf("No blocks found for %s", edit.page)
			continue
		}

		if pageData.Blocks[0].Content != edit.content {
			t.Errorf("Edit lost for %s! Expected: %s, Got: %s", 
				edit.page, edit.content, pageData.Blocks[0].Content)
		}
	}
}

// TestBacklinkUpdateOnNavigation verifies backlinks update correctly
func TestBacklinkUpdateOnNavigation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-backlink-nav-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create initial page
	pageFile := filepath.Join(tempDir, "journal.md")
	content := `# Journal

- Today I learned about Go
- Need to check [[Resources]]`

	if err := os.WriteFile(pageFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write journal: %v", err)
	}

	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}

	// Edit to add a new reference
	_, err = app.UpdateBlockAtPath("Journal", BlockPath{0}, "Today I learned about Go and [[Testing]]")
	if err != nil {
		t.Fatalf("Failed to update block: %v", err)
	}

	// Navigate to the new page (auto-creates it)
	_, err = app.GetPage("Testing")
	if err != nil {
		t.Fatalf("Failed to navigate to Testing: %v", err)
	}

	// Check backlinks
	backlinks := app.GetBacklinks("Testing")
	if len(backlinks) != 1 || len(backlinks["Journal"]) == 0 {
		t.Errorf("Expected backlink from Journal to Testing, got: %v", backlinks)
	}

	// Navigate back and verify
	pageData, err := app.GetPage("Journal")
	if err != nil {
		t.Fatalf("Failed to navigate back to Journal: %v", err)
	}

	// Verify the edit persisted
	if !strings.Contains(pageData.Blocks[0].Content, "[[Testing]]") {
		t.Error("Edit with new link was lost after navigation")
	}
}