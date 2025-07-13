// MIT License
//
// Copyright (c) 2025 Rehan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBugFixedWithPositionalAPI demonstrates that the block ID bug is fixed
// when using positional addressing instead of IDs
func TestBugFixedWithPositionalAPI(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "seq2b_bug_fixed_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test page
	pageContent := `# Test Page

- [[link text]]
`
	
	pagePath := filepath.Join(tempDir, "test-page.md")
	if err := os.WriteFile(pagePath, []byte(pageContent), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	// Create app and load pages
	app := &App{}
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Step 1: Add a new block at position 1
	addDelta, err := app.AddBlockAtPath("Test Page", BlockPath{1}, "New block content")
	if err != nil {
		t.Fatalf("Failed to add block: %v", err)
	}
	
	// Verify the block was added
	page := app.pages["Test Page"]
	if len(page.Blocks) != 2 {
		t.Fatalf("Expected 2 blocks after add, got %d", len(page.Blocks))
	}
	
	// Get the path from the delta (should be [1])
	addedPath := addDelta["path"].(BlockPath)
	if len(addedPath) != 1 || addedPath[0] != 1 {
		t.Fatalf("Expected path [1], got %v", addedPath)
	}
	
	// Step 2: Update the newly added block using its path
	updateDelta, err := app.UpdateBlockAtPath("Test Page", addedPath, "Updated content")
	if err != nil {
		t.Fatalf("Failed to update block: %v", err)
	}
	
	// Verify the update succeeded
	if updateDelta["action"] != "update" {
		t.Errorf("Expected action 'update', got %v", updateDelta["action"])
	}
	
	// Verify the content was actually updated
	updatedBlock, err := FindBlockByPath(page.Blocks, addedPath)
	if err != nil {
		t.Fatalf("Failed to find block by path: %v", err)
	}
	
	if updatedBlock.Content != "Updated content" {
		t.Errorf("Expected 'Updated content', got '%s'", updatedBlock.Content)
	}
	
	// Step 3: Add another block and update the first one again
	_, err = app.AddBlockAtPath("Test Page", BlockPath{2}, "Third block")
	if err != nil {
		t.Fatalf("Failed to add third block: %v", err)
	}
	
	// The first added block is still at position [1]
	_, err = app.UpdateBlockAtPath("Test Page", BlockPath{1}, "Final content")
	if err != nil {
		t.Fatalf("Failed to update block again: %v", err)
	}
	
	// Verify final state
	finalBlock, err := FindBlockByPath(page.Blocks, BlockPath{1})
	if err != nil {
		t.Fatalf("Failed to find block: %v", err)
	}
	
	if finalBlock.Content != "Final content" {
		t.Errorf("Expected 'Final content', got '%s'", finalBlock.Content)
	}
	
	t.Log("Bug is fixed! Block updates work correctly with positional addressing")
}