// +build integration

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEnterKeyScenario simulates the exact bug scenario:
// 1. User types [[link]]
// 2. User presses Enter (creates new block)
// 3. User types in new block
// 4. User presses Enter again (this was failing)
func TestEnterKeyScenario(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a date page
	testFile := filepath.Join(tempDir, "2025-01-13.md")
	content := `# Jan 13th, 2025

- First block`
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Get initial state
	pageData, err := app.GetPage("Jan 13th, 2025")
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}
	
	firstBlockID := pageData.Blocks[0].ID
	
	// Simulate typing [[link]] in first block
	err = app.UpdateBlock("Jan 13th, 2025", firstBlockID, "[[asdfas]]")
	if err != nil {
		t.Fatalf("Failed to update first block: %v", err)
	}
	
	// Simulate pressing Enter - creates new block
	result1, err := app.AddBlock("Jan 13th, 2025", "", firstBlockID, "", 0)
	if err != nil {
		t.Fatalf("Failed to create first new block: %v", err)
	}
	
	newBlockID1 := result1["id"].(string)
	
	// Type something in the new block
	err = app.UpdateBlock("Jan 13th, 2025", newBlockID1, "Some text")
	if err != nil {
		t.Fatalf("Failed to update new block: %v", err)
	}
	
	// Press Enter again - this was failing before the fix
	result2, err := app.AddBlock("Jan 13th, 2025", "", newBlockID1, "", 0)
	if err != nil {
		t.Fatalf("Failed to create second new block: %v. This is the bug we fixed!", err)
	}
	
	newBlockID2 := result2["id"].(string)
	
	// Verify we can update the second new block too
	err = app.UpdateBlock("Jan 13th, 2025", newBlockID2, "More text")
	if err != nil {
		t.Fatalf("Failed to update second new block: %v", err)
	}
	
	// Verify final state
	pageData, err = app.GetPage("Jan 13th, 2025")
	if err != nil {
		t.Fatalf("Failed to get final page state: %v", err)
	}
	
	if len(pageData.Blocks) != 3 {
		t.Errorf("Expected 3 blocks, got %d", len(pageData.Blocks))
	}
	
	expectedContents := []string{"[[asdfas]]", "Some text", "More text"}
	for i, expected := range expectedContents {
		if pageData.Blocks[i].Content != expected {
			t.Errorf("Block %d: expected %q, got %q", i, expected, pageData.Blocks[i].Content)
		}
	}
}