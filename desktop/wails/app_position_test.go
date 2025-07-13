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

func TestPositionalAddAndUpdate(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "seq2b_pos_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test page
	pageContent := `# Test Page

- First block
  - First child
  - Second child
- Second block
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
	
	// Test 1: Add block at top level (position 2)
	t.Run("AddBlockAtTopLevel", func(t *testing.T) {
		delta, err := app.AddBlockAtPath("Test Page", BlockPath{2}, "Third block")
		if err != nil {
			t.Errorf("Failed to add block: %v", err)
		}
		
		// Verify delta
		if delta["action"] != "add" {
			t.Errorf("Expected action 'add', got %v", delta["action"])
		}
		
		path := delta["path"].(BlockPath)
		if len(path) != 1 || path[0] != 2 {
			t.Errorf("Expected path [2], got %v", path)
		}
		
		// Verify block was added
		page := app.pages["Test Page"]
		if len(page.Blocks) != 3 {
			t.Errorf("Expected 3 blocks, got %d", len(page.Blocks))
		}
		if page.Blocks[2].Content != "Third block" {
			t.Errorf("Expected 'Third block', got %s", page.Blocks[2].Content)
		}
	})
	
	// Test 2: Update block using path
	t.Run("UpdateBlockAtPath", func(t *testing.T) {
		delta, err := app.UpdateBlockAtPath("Test Page", BlockPath{0, 1}, "Updated second child")
		if err != nil {
			t.Errorf("Failed to update block: %v", err)
		}
		
		// Verify delta
		if delta["action"] != "update" {
			t.Errorf("Expected action 'update', got %v", delta["action"])
		}
		
		// Verify block was updated
		page := app.pages["Test Page"]
		secondChild := page.Blocks[0].Children[1]
		if secondChild.Content != "Updated second child" {
			t.Errorf("Expected 'Updated second child', got %s", secondChild.Content)
		}
	})
	
	// Test 3: Add nested block
	t.Run("AddNestedBlock", func(t *testing.T) {
		_, err := app.AddBlockAtPath("Test Page", BlockPath{0, 2}, "Third child")
		if err != nil {
			t.Errorf("Failed to add nested block: %v", err)
		}
		
		// Verify block was added
		page := app.pages["Test Page"]
		if len(page.Blocks[0].Children) != 3 {
			t.Errorf("Expected 3 children, got %d", len(page.Blocks[0].Children))
		}
		if page.Blocks[0].Children[2].Content != "Third child" {
			t.Errorf("Expected 'Third child', got %s", page.Blocks[0].Children[2].Content)
		}
		
		// Verify parent relationship
		if page.Blocks[0].Children[2].Parent != page.Blocks[0] {
			t.Errorf("Parent relationship not set correctly")
		}
	})
	
	// Test 4: Verify path shifts
	t.Run("PathShiftsAfterInsert", func(t *testing.T) {
		// Insert at beginning
		delta, err := app.AddBlockAtPath("Test Page", BlockPath{0}, "New first block")
		if err != nil {
			t.Errorf("Failed to add block at beginning: %v", err)
		}
		
		shiftsInterface := delta["shifts"]
		if shiftsInterface != nil {
			shifts := shiftsInterface.([]PathShift)
			// Should have shifts for blocks that moved
			if len(shifts) == 0 {
				t.Errorf("Expected path shifts, got none")
			}
		}
	})
}

func TestPositionalAPIErrors(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "seq2b_pos_err_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test page
	pageContent := `# Test Page

- Single block
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
	
	// Test invalid paths
	tests := []struct {
		name string
		path BlockPath
	}{
		{"Empty path", BlockPath{}},
		{"Invalid index", BlockPath{5}},
		{"Invalid nested path", BlockPath{0, 10}},
		{"Too deep path", BlockPath{0, 0, 0}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := app.UpdateBlockAtPath("Test Page", tt.path, "content")
			if err == nil {
				t.Errorf("Expected error for %s, got none", tt.name)
			}
		})
	}
}