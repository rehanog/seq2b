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

func TestParseLogseqFeatures(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected Line
	}{
		{
			name: "Block with ID",
			line: "- This is a block with an ID",
			expected: Line{
				Type:    TypeBlock,
				Content: "This is a block with an ID",
				BlockID: "",
			},
		},
		{
			name: "Block with ID property",
			line: "- This block has id:: 550e8400-e29b-41d4-a716-446655440000",
			expected: Line{
				Type:    TypeBlock,
				Content: "This block has id:: 550e8400-e29b-41d4-a716-446655440000",
				BlockID: "550e8400-e29b-41d4-a716-446655440000",
			},
		},
		{
			name: "Block with property",
			line: "- type:: note",
			expected: Line{
				Type:       TypeBlock,
				Content:    "type:: note",
				Properties: map[string]string{"type": "note"},
			},
		},
		{
			name: "Block with tag",
			line: "- This is a #important note",
			expected: Line{
				Type:    TypeBlock,
				Content: "This is a #important note",
				Tags:    []string{"important"},
			},
		},
		{
			name: "Block with multiple tags",
			line: "- Working on #logseq #parser #golang",
			expected: Line{
				Type:    TypeBlock,
				Content: "Working on #logseq #parser #golang",
				Tags:    []string{"logseq", "parser", "golang"},
			},
		},
		{
			name: "Block reference",
			line: "- See ((550e8400-e29b-41d4-a716-446655440000))",
			expected: Line{
				Type:    TypeBlock,
				Content: "See ((550e8400-e29b-41d4-a716-446655440000))",
			},
		},
		{
			name: "NOW state",
			line: "- NOW Working on this task",
			expected: Line{
				Type:    TypeBlock,
				Content: "NOW Working on this task",
				TodoInfo: TodoInfo{
					TodoState: TodoStateNow,
				},
			},
		},
		{
			name: "WAIT state",
			line: "- WAIT Waiting for feedback",
			expected: Line{
				Type:    TypeBlock,
				Content: "WAIT Waiting for feedback",
				TodoInfo: TodoInfo{
					TodoState: TodoStateWait,
				},
			},
		},
		{
			name: "CANCELLED state",
			line: "- CANCELLED This was cancelled",
			expected: Line{
				Type:    TypeBlock,
				Content: "CANCELLED This was cancelled",
				TodoInfo: TodoInfo{
					TodoState: TodoStateCancelled,
				},
			},
		},
		{
			name: "LATER state",
			line: "- LATER Will do this eventually",
			expected: Line{
				Type:    TypeBlock,
				Content: "LATER Will do this eventually",
				TodoInfo: TodoInfo{
					TodoState: TodoStateLater,
				},
			},
		},
		{
			name: "Complex line with multiple features",
			line: "- TODO [#A] Fix #bug in [[parser]] id:: abc-123",
			expected: Line{
				Type:    TypeBlock,
				Content: "TODO [#A] Fix #bug in [[parser]] id:: abc-123",
				TodoInfo: TodoInfo{
					TodoState: TodoStateTodo,
					Priority:  "A",
				},
				References: []string{"parser"},
				Tags:       []string{"bug"},
				BlockID:    "abc-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLine(1, tt.line)
			
			// Check basic fields
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}
			if result.Content != tt.expected.Content {
				t.Errorf("Content = %v, want %v", result.Content, tt.expected.Content)
			}
			
			// Check BlockID
			if result.BlockID != tt.expected.BlockID {
				t.Errorf("BlockID = %v, want %v", result.BlockID, tt.expected.BlockID)
			}
			
			// Check TodoInfo
			if result.TodoInfo.TodoState != tt.expected.TodoInfo.TodoState {
				t.Errorf("TodoState = %v, want %v", result.TodoInfo.TodoState, tt.expected.TodoInfo.TodoState)
			}
			if result.TodoInfo.Priority != tt.expected.TodoInfo.Priority {
				t.Errorf("Priority = %v, want %v", result.TodoInfo.Priority, tt.expected.TodoInfo.Priority)
			}
			
			// Check Properties
			if len(result.Properties) != len(tt.expected.Properties) {
				t.Errorf("Properties length = %v, want %v", len(result.Properties), len(tt.expected.Properties))
			}
			for k, v := range tt.expected.Properties {
				if result.Properties[k] != v {
					t.Errorf("Properties[%s] = %v, want %v", k, result.Properties[k], v)
				}
			}
			
			// Check Tags
			if len(result.Tags) != len(tt.expected.Tags) {
				t.Errorf("Tags = %v, want %v", result.Tags, tt.expected.Tags)
			}
			
			// Check References
			if len(tt.expected.References) > 0 && len(result.References) != len(tt.expected.References) {
				t.Errorf("References = %v, want %v", result.References, tt.expected.References)
			}
		})
	}
}

func TestParseLogseqFile(t *testing.T) {
	// Test that we can parse a file with Logseq features without crashing
	content := `# Test Page
property:: value
tags:: test, import

- Normal block
- Block with id:: 123-456
- TODO Task
- NOW Current task
- WAIT Waiting task
- LATER Future task
- CANCELLED Cancelled task
- Block with #tag
- Block with [[link]]
- Block with ((123-456)) reference
- collapsed:: true
  - Hidden content
- {{query (todo TODO)}}
- {{embed ((123-456))}}
`

	lines := []string{}
	for i, line := range splitLines(content) {
		parsed := ParseLine(i+1, line)
		// The parser should not crash
		lines = append(lines, parsed.Content)
	}
	
	// Verify we parsed all lines
	if len(lines) == 0 {
		t.Error("Failed to parse any lines")
	}
}

// Helper function to split content into lines
func splitLines(content string) []string {
	var lines []string
	current := ""
	for _, r := range content {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}