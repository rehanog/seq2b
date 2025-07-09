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

package parser

import (
	"testing"
)

func TestCalculateIndentLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"no indent", "- Item", 0},
		{"2 spaces", "  - Item", 1},
		{"4 spaces", "    - Item", 2},
		{"6 spaces", "      - Item", 3},
		{"tab treated as spaces", "\t- Item", 0}, // tabs not counted
		{"mixed content", "  some text", 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateIndentLevel(tt.input)
			if result != tt.expected {
				t.Errorf("calculateIndentLevel(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildBlockTree(t *testing.T) {
	// Create test contexts
	contexts := []parseContext{
		{Line{1, TypeBlock, "A", 0}, 0},        // Top level A
		{Line{2, TypeBlock, "A.1", 0}, 1},      // Child of A
		{Line{3, TypeBlock, "A.2", 0}, 1},      // Another child of A
		{Line{4, TypeBlock, "A.2.1", 0}, 2},    // Child of A.2
		{Line{5, TypeBlock, "B", 0}, 0},        // Top level B
		{Line{6, TypeBlock, "B.1", 0}, 1},      // Child of B
	}
	
	blocks := BuildBlockTree(contexts)
	
	// Check top level
	if len(blocks) != 2 {
		t.Errorf("Expected 2 top-level blocks, got %d", len(blocks))
	}
	
	// Check first block hierarchy
	if blocks[0].GetContent() != "A" {
		t.Errorf("First block content = %q, want %q", blocks[0].GetContent(), "A")
	}
	if len(blocks[0].Children) != 2 {
		t.Errorf("First block children = %d, want 2", len(blocks[0].Children))
	}
	
	// Check nested child
	if blocks[0].Children[1].GetContent() != "A.2" {
		t.Errorf("A.2 content = %q, want %q", blocks[0].Children[1].GetContent(), "A.2")
	}
	if len(blocks[0].Children[1].Children) != 1 {
		t.Errorf("A.2 children = %d, want 1", len(blocks[0].Children[1].Children))
	}
	if blocks[0].Children[1].Children[0].GetContent() != "A.2.1" {
		t.Errorf("A.2.1 content = %q, want %q", blocks[0].Children[1].Children[0].GetContent(), "A.2.1")
	}
	
	// Check parent relationships
	if blocks[0].Children[0].Parent != blocks[0] {
		t.Error("A.1 parent should be A")
	}
	if blocks[0].Children[1].Children[0].Parent != blocks[0].Children[1] {
		t.Error("A.2.1 parent should be A.2")
	}
	
	// Check depths
	if blocks[0].Depth != 0 {
		t.Errorf("A depth = %d, want 0", blocks[0].Depth)
	}
	if blocks[0].Children[0].Depth != 1 {
		t.Errorf("A.1 depth = %d, want 1", blocks[0].Children[0].Depth)
	}
	if blocks[0].Children[1].Children[0].Depth != 2 {
		t.Errorf("A.2.1 depth = %d, want 2", blocks[0].Children[1].Children[0].Depth)
	}
}

func TestParseFile(t *testing.T) {
	input := `# Test Page

- Top level block
  - Nested block with **bold**
    - Double nested with [[link]]
  - Another nested
- Second top level

Regular paragraph text.`

	result, err := ParseFile(input)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	
	// Check page title
	if result.Page.Title != "Test Page" {
		t.Errorf("Page title = %q, want %q", result.Page.Title, "Test Page")
	}
	
	// Check block count
	if len(result.Page.Blocks) != 2 {
		t.Errorf("Top-level blocks = %d, want 2", len(result.Page.Blocks))
	}
	
	// Check total blocks
	if len(result.Page.AllBlocks) != 5 {
		t.Errorf("Total blocks = %d, want 5", len(result.Page.AllBlocks))
	}
	
	// Check nested structure
	firstBlock := result.Page.Blocks[0]
	if len(firstBlock.Children) != 2 {
		t.Errorf("First block children = %d, want 2", len(firstBlock.Children))
	}
	
	// Check HTML rendering in nested blocks
	nestedBlock := firstBlock.Children[0]
	html := nestedBlock.RenderHTML()
	expectedHTML := "Nested block with <b>bold</b>"
	if html != expectedHTML {
		t.Errorf("Nested block HTML = %q, want %q", html, expectedHTML)
	}
}

func TestBlockMethods(t *testing.T) {
	// Test AddChild
	parent := &Block{ID: "parent", Depth: 0}
	child := &Block{ID: "child"}
	
	parent.AddChild(child)
	
	if child.Parent != parent {
		t.Error("Child parent not set correctly")
	}
	if child.Depth != 1 {
		t.Errorf("Child depth = %d, want 1", child.Depth)
	}
	if len(parent.Children) != 1 {
		t.Errorf("Parent children = %d, want 1", len(parent.Children))
	}
}

func TestGetAllBlocks(t *testing.T) {
	// Build a simple tree
	page := &Page{
		Blocks: []*Block{
			{
				ID: "1",
				Children: []*Block{
					{ID: "1.1"},
					{
						ID: "1.2",
						Children: []*Block{
							{ID: "1.2.1"},
						},
					},
				},
			},
			{
				ID: "2",
				Children: []*Block{
					{ID: "2.1"},
				},
			},
		},
	}
	
	allBlocks := page.GetAllBlocks()
	
	// Should have 6 blocks total
	if len(allBlocks) != 6 {
		t.Errorf("GetAllBlocks returned %d blocks, want 6", len(allBlocks))
	}
	
	// Check that all IDs are present
	idMap := make(map[string]bool)
	for _, block := range allBlocks {
		idMap[block.ID] = true
	}
	
	expectedIDs := []string{"1", "1.1", "1.2", "1.2.1", "2", "2.1"}
	for _, id := range expectedIDs {
		if !idMap[id] {
			t.Errorf("Block ID %s not found in GetAllBlocks result", id)
		}
	}
}