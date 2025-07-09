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
	"reflect"
	"testing"
)

func TestExtractPageLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "no links",
			input:    "Just plain text",
			expected: []string{},
		},
		{
			name:     "single link",
			input:    "Reference to [[Page A]]",
			expected: []string{"Page A"},
		},
		{
			name:     "multiple links",
			input:    "Links to [[Page A]] and [[Page B]]",
			expected: []string{"Page A", "Page B"},
		},
		{
			name:     "duplicate links",
			input:    "First [[Page A]] and another [[Page A]]",
			expected: []string{"Page A", "Page A"}, // Keep duplicates for counting
		},
		{
			name:     "link with spaces",
			input:    "Link to [[My Complex Page Name]]",
			expected: []string{"My Complex Page Name"},
		},
		{
			name:     "mixed formatting",
			input:    "Text with **bold** and [[Page Link]] and *italic*",
			expected: []string{"Page Link"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPageLinks(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractPageLinks(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBacklinkIndex(t *testing.T) {
	// Create test pages
	pageA := &Page{
		Title: "Page A",
		Blocks: []*Block{
			{
				ID:      "block-1",
				Content: "References [[Page B]] and [[Page C]]",
			},
			{
				ID:      "block-2",
				Content: "Another [[Page B]] reference",
			},
			{
				ID:      "block-3",
				Content: "No references here",
			},
		},
	}
	pageA.AllBlocks = pageA.Blocks

	pageB := &Page{
		Title: "Page B",
		Blocks: []*Block{
			{
				ID:      "block-1",
				Content: "Back to [[Page A]]",
			},
		},
	}
	pageB.AllBlocks = pageB.Blocks

	// Build index
	idx := NewBacklinkIndex()
	idx.AddPage(pageA)
	idx.AddPage(pageB)

	// Test forward links
	t.Run("forward links", func(t *testing.T) {
		forwardLinks := idx.GetForwardLinks("Page A")
		if len(forwardLinks) != 2 {
			t.Errorf("Page A forward links = %d, want 2", len(forwardLinks))
		}
		
		// Check Page B references from Page A
		pageBRefs := forwardLinks["Page B"]
		if len(pageBRefs) != 2 {
			t.Errorf("Page A -> Page B refs = %d, want 2", len(pageBRefs))
		}
	})

	// Test backward links
	t.Run("backward links", func(t *testing.T) {
		backlinks := idx.GetBacklinks("Page A")
		if len(backlinks) != 1 {
			t.Errorf("Page A backlinks = %d, want 1", len(backlinks))
		}
		
		// Check that Page B links to Page A
		if refs, ok := backlinks["Page B"]; !ok || len(refs) != 1 {
			t.Errorf("Expected Page B to link to Page A once")
		}
	})

	// Test orphan detection
	t.Run("orphan detection", func(t *testing.T) {
		// Add an orphan page
		orphanPage := &Page{
			Title:  "Orphan Page",
			Blocks: []*Block{{ID: "block-1", Content: "No links"}},
		}
		orphanPage.AllBlocks = orphanPage.Blocks
		idx.AddPage(orphanPage)
		
		if !idx.IsOrphanPage("Orphan Page") {
			t.Error("Expected Orphan Page to be detected as orphan")
		}
		
		if idx.IsOrphanPage("Page A") {
			t.Error("Page A should not be an orphan")
		}
	})
}

func TestFindOrphanBlocks(t *testing.T) {
	blocks := []*Block{
		{ID: "block-1", Content: "Has [[link]]"},
		{ID: "block-2", Content: "No links here"},
		{ID: "block-3", Content: ""},
		{ID: "block-4", Content: "Multiple [[Page A]] and [[Page B]]"},
		{ID: "block-5", Content: "Just text"},
	}

	orphans := FindOrphanBlocks(blocks)
	
	if len(orphans) != 2 {
		t.Errorf("FindOrphanBlocks returned %d blocks, want 2", len(orphans))
	}
	
	// Check that the correct blocks are identified as orphans
	orphanIDs := make(map[string]bool)
	for _, block := range orphans {
		orphanIDs[block.ID] = true
	}
	
	if !orphanIDs["block-2"] || !orphanIDs["block-5"] {
		t.Errorf("Expected block-2 and block-5 to be orphans")
	}
}

