package main

import (
	"os"
	"path/filepath"
	"testing"
	
)

func TestAddBlock(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "seq2b-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test markdown file
	testFile := filepath.Join(tempDir, "test-page.md")
	content := `# Test Page

- Block 1
- Block 2
  - Block 2.1
- Block 3`
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// Create app and load the directory
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Get the page to find block IDs
	pageData, err := app.GetPage("Test Page")
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}
	
	if len(pageData.Blocks) != 3 {
		t.Fatalf("Expected 3 top-level blocks, got %d", len(pageData.Blocks))
	}
	
	// Test 1: Add block after Block 1
	block1ID := pageData.Blocks[0].ID
	result1, err := app.AddBlock("Test Page", "", block1ID, "New block after 1", 0)
	if err != nil {
		t.Fatalf("Failed to add block after Block 1: %v", err)
	}
	
	newBlockID1 := result1["id"].(string)
	
	// Verify the block was added
	pageData, err = app.GetPage("Test Page")
	if err != nil {
		t.Fatalf("Failed to get page after adding block: %v", err)
	}
	
	if len(pageData.Blocks) != 4 {
		t.Fatalf("Expected 4 blocks after addition, got %d", len(pageData.Blocks))
	}
	
	// Check that the new block is in the right position
	if pageData.Blocks[1].ID != newBlockID1 {
		t.Errorf("New block not in expected position")
	}
	
	if pageData.Blocks[1].Content != "New block after 1" {
		t.Errorf("New block content mismatch: got %q", pageData.Blocks[1].Content)
	}
	
	// Test 2: Add child block to Block 2
	block2ID := pageData.Blocks[2].ID // Block 2 is now at index 2
	result2, err := app.AddBlock("Test Page", block2ID, "", "New child of Block 2", 1)
	if err != nil {
		t.Fatalf("Failed to add child block: %v", err)
	}
	
	newBlockID2 := result2["id"].(string)
	
	// Verify the child was added
	pageData, err = app.GetPage("Test Page")
	if err != nil {
		t.Fatalf("Failed to get page after adding child: %v", err)
	}
	
	// Block 2 should now have 2 children
	if len(pageData.Blocks[2].Children) != 2 {
		t.Fatalf("Expected 2 children for Block 2, got %d", len(pageData.Blocks[2].Children))
	}
	
	// The new child should be last (since we didn't specify afterBlockID)
	if pageData.Blocks[2].Children[1].ID != newBlockID2 {
		t.Errorf("New child block not in expected position")
	}
	
	// Test 3: Try to add block with non-existent parent
	_, err = app.AddBlock("Test Page", "non-existent-id", "", "Should fail", 0)
	if err == nil {
		t.Errorf("Expected error when adding block with non-existent parent")
	}
	
	// Test 4: Try to add block to non-existent page
	_, err = app.AddBlock("Non-existent Page", "", "", "Should fail", 0)
	if err == nil {
		t.Errorf("Expected error when adding block to non-existent page")
	}
}

func TestUpdateBlockAfterAddBlock(t *testing.T) {
	// This test specifically checks the bug we encountered where
	// updating a newly created block would fail because the backend
	// didn't refresh its page structure after AddBlock.
	// The fix was to ensure RefreshPages is called after adding a block.
	tempDir, err := os.MkdirTemp("", "seq2b-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test markdown file
	testFile := filepath.Join(tempDir, "test-page.md")
	content := `# Test Page

- Initial block`
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// Create app and load the directory
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Get the initial block ID
	pageData, err := app.GetPage("Test Page")
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}
	
	initialBlockID := pageData.Blocks[0].ID
	
	// Simulate the Enter key scenario:
	// 1. Update the current block with new content
	err = app.UpdateBlock("Test Page", initialBlockID, "[[link text]]")
	if err != nil {
		t.Fatalf("Failed to update initial block: %v", err)
	}
	
	// 2. Add a new block after it
	result, err := app.AddBlock("Test Page", "", initialBlockID, "", 0)
	if err != nil {
		t.Fatalf("Failed to add new block: %v", err)
	}
	
	newBlockID := result["id"].(string)
	
	// 3. Try to update the newly created block (this was failing before)
	// First, let's check if the block exists in the backend
	pageData, _ = app.GetPage("Test Page")
	foundNewBlock := false
	for _, block := range pageData.Blocks {
		if block.ID == newBlockID {
			foundNewBlock = true
			break
		}
	}
	
	if !foundNewBlock {
		t.Logf("New block %s not found in page after AddBlock", newBlockID)
		t.Logf("Page has %d blocks", len(pageData.Blocks))
		for i, b := range pageData.Blocks {
			t.Logf("Block %d: ID=%s, Content=%q", i, b.ID, b.Content)
		}
	}
	
	err = app.UpdateBlock("Test Page", newBlockID, "New content")
	if err != nil {
		t.Errorf("Failed to update newly created block: %v", err)
	}
	
	// Verify the updates
	pageData, err = app.GetPage("Test Page")
	if err != nil {
		t.Fatalf("Failed to get page after updates: %v", err)
	}
	
	if len(pageData.Blocks) != 2 {
		t.Fatalf("Expected 2 blocks, got %d", len(pageData.Blocks))
	}
	
	if pageData.Blocks[0].Content != "[[link text]]" {
		t.Errorf("First block content mismatch: got %q", pageData.Blocks[0].Content)
	}
	
	if pageData.Blocks[1].Content != "New content" {
		t.Errorf("Second block content mismatch: got %q", pageData.Blocks[1].Content)
	}
}

func TestAddBlockWithMarkdownContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test markdown file
	testFile := filepath.Join(tempDir, "test-page.md")
	content := `# Test Page

- Block 1`
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// Create app and load the directory
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Add a block with markdown content including links
	result, err := app.AddBlock("Test Page", "", "", "Text with [[Page Link]] and **bold**", 0)
	if err != nil {
		t.Fatalf("Failed to add block with markdown: %v", err)
	}
	
	// Get the block data from the result
	blockData := result["block"].(BlockData)
	
	// Check that segments were parsed correctly
	if len(blockData.Segments) == 0 {
		t.Errorf("Expected segments to be parsed, got none")
	}
	
	// Verify the segments include the link
	hasLink := false
	for _, segment := range blockData.Segments {
		if segment.Type == "link" && segment.Target == "Page Link" {
			hasLink = true
			break
		}
	}
	
	if !hasLink {
		t.Errorf("Expected to find link segment for [[Page Link]]")
	}
}