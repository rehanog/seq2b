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
	"strings"
	"testing"
)

func TestEditPreservesLogseqMetadata(t *testing.T) {
	tests := []struct {
		name           string
		originalBlock  string
		editedContent  string
		shouldPreserve []string
	}{
		{
			name: "Edit preserves block ID",
			originalBlock: `- This is a block with some text
  id:: 550e8400-e29b-41d4-a716-446655440000`,
			editedContent: `This is a block with EDITED text
id:: 550e8400-e29b-41d4-a716-446655440000`,
			shouldPreserve: []string{"id:: 550e8400-e29b-41d4-a716-446655440000"},
		},
		{
			name: "Edit preserves properties",
			originalBlock: `- Task description
  priority:: high
  assigned:: [[John Doe]]
  due:: 2025-01-25`,
			editedContent: `Updated task description
priority:: high
assigned:: [[John Doe]]
due:: 2025-01-25`,
			shouldPreserve: []string{"priority:: high", "assigned:: [[John Doe]]", "due:: 2025-01-25"},
		},
		{
			name: "Edit preserves tags",
			originalBlock: `- Working on #logseq #parser #golang implementation`,
			editedContent: `Still working on #logseq #parser #golang implementation`,
			shouldPreserve: []string{"#logseq", "#parser", "#golang"},
		},
		{
			name: "Edit preserves TODO state and priority",
			originalBlock: `- TODO [#A] Fix critical bug`,
			editedContent: `TODO [#A] Fix critical bug in parser`,
			shouldPreserve: []string{"TODO", "[#A]"},
		},
		{
			name: "Edit preserves block references",
			originalBlock: `- See details in ((550e8400-e29b-41d4-a716-446655440000))`,
			editedContent: `Check the details in ((550e8400-e29b-41d4-a716-446655440000))`,
			shouldPreserve: []string{"((550e8400-e29b-41d4-a716-446655440000))"},
		},
		{
			name: "Edit preserves query blocks",
			originalBlock: `- {{query (todo TODO)}}`,
			editedContent: `{{query (todo TODO)}}`,
			shouldPreserve: []string{"{{query (todo TODO)}}"},
		},
		{
			name: "Edit preserves complex metadata",
			originalBlock: `- TODO [#A] Task with #urgent tag and ((ref))
  id:: abc-123
  status:: in-progress
  SCHEDULED: <2025-01-20>`,
			editedContent: `TODO [#A] Updated task with #urgent tag and ((ref))
id:: abc-123
status:: in-progress
SCHEDULED: <2025-01-20>`,
			shouldPreserve: []string{
				"TODO", "[#A]", "#urgent", "((ref))",
				"id:: abc-123", "status:: in-progress",
				"SCHEDULED: <2025-01-20>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the original block
			lines := strings.Split(tt.originalBlock, "\n")
			parsedLines := make([]Line, len(lines))
			for i, line := range lines {
				parsedLines[i] = ParseLine(i+1, line)
			}

			// Create a block
			block := &Block{
				Lines: parsedLines,
			}
			block.updateContent()

			// Store original metadata
			originalBlockID := block.BlockID
			originalProperties := make(map[string]string)
			for k, v := range block.Properties {
				originalProperties[k] = v
			}
			originalTags := make([]string, len(block.Tags))
			copy(originalTags, block.Tags)

			// Edit the block content
			block.SetContent(tt.editedContent)

			// Verify metadata is preserved
			if originalBlockID != "" && block.BlockID != originalBlockID {
				t.Errorf("Block ID changed from %s to %s", originalBlockID, block.BlockID)
			}

			// Check that all expected content is preserved
			for _, expected := range tt.shouldPreserve {
				if !strings.Contains(block.Content, expected) {
					t.Errorf("Expected content %q not found after edit", expected)
				}
			}

			// Verify the edited content is applied
			if !strings.Contains(block.Content, strings.Split(tt.editedContent, "\n")[0]) {
				t.Errorf("Edited content not applied correctly")
			}
		})
	}
}

func TestBlockReferencePreservation(t *testing.T) {
	// Test that we can edit a block without losing references to it
	block1 := &Block{
		Lines: []Line{
			ParseLine(1, "- Original block content"),
			ParseLine(2, "  id:: block-id-123"),
		},
	}
	block1.updateContent()

	block2 := &Block{
		Lines: []Line{
			ParseLine(1, "- This references ((block-id-123))"),
		},
	}
	block2.updateContent()

	// Verify block1 has the ID
	if block1.BlockID != "block-id-123" {
		t.Errorf("Block ID not parsed correctly: %s", block1.BlockID)
	}

	// Edit block1
	block1.SetContent("Edited block content\nid:: block-id-123")

	// Verify ID is preserved
	if block1.BlockID != "block-id-123" {
		t.Errorf("Block ID lost after edit: %s", block1.BlockID)
	}

	// Verify block2 still has the reference
	if !strings.Contains(block2.Content, "((block-id-123))") {
		t.Errorf("Block reference lost in referencing block")
	}
}

func TestPropertyEditPreservation(t *testing.T) {
	// Test editing a block with properties
	content := `- Meeting notes
  date:: 2025-01-19
  attendees:: [[John]], [[Jane]]
  outcome:: successful`

	lines := strings.Split(content, "\n")
	parsedLines := make([]Line, len(lines))
	for i, line := range lines {
		parsedLines[i] = ParseLine(i+1, line)
	}

	block := &Block{Lines: parsedLines}
	block.updateContent()

	// Verify properties are parsed
	if len(block.Properties) != 3 {
		t.Errorf("Expected 3 properties, got %d", len(block.Properties))
	}

	// Edit just the main text
	newContent := `Updated meeting notes
date:: 2025-01-19
attendees:: [[John]], [[Jane]]
outcome:: successful`

	block.SetContent(newContent)

	// Verify properties are preserved
	if len(block.Properties) != 3 {
		t.Errorf("Properties lost after edit: expected 3, got %d", len(block.Properties))
	}

	if block.Properties["date"] != "2025-01-19" {
		t.Errorf("Date property changed: %s", block.Properties["date"])
	}
}

func TestEditWithNewMetadata(t *testing.T) {
	// Test adding new metadata during edit
	originalContent := "- Simple block"
	
	block := &Block{
		Lines: []Line{ParseLine(1, originalContent)},
	}
	block.updateContent()

	// Edit to add metadata
	newContent := `- Simple block with metadata
  id:: new-id-456
  priority:: high`

	block.SetContent(newContent)

	// Debug output
	t.Logf("Block content: %q", block.Content)
	t.Logf("Block ID: %q", block.BlockID)
	t.Logf("Properties: %+v", block.Properties)
	t.Logf("Lines: %d", len(block.Lines))
	for i, line := range block.Lines {
		t.Logf("Line %d: Type=%v, Content=%q, BlockID=%q, Props=%+v", 
			i, line.Type, line.Content, line.BlockID, line.Properties)
	}

	// Verify new metadata is added
	if block.BlockID != "new-id-456" {
		t.Errorf("New block ID not added: %s", block.BlockID)
	}

	if block.Properties["priority"] != "high" {
		t.Errorf("New property not added: %s", block.Properties["priority"])
	}
}