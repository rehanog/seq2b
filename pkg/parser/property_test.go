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
	"os"
	"strings"
	"testing"
)

// BUG ANALYSIS: Page Properties Not Supported
// 
// CURRENT STATE:
// - extractProperties() function exists and works for individual lines ✓
// - Line struct has Properties field ✓  
// - Block struct has Properties field ✓
// - Page struct has Properties field ✓ [FIXED]
// - extractPageLevelProperties() function implemented ✓ [FIXED]
//
// EXPECTED BEHAVIOR:
// - Properties at top of page (before first header/block) should be page-level ✓
// - Page struct should have Properties map[string]string field ✓
// - Nested properties should remain as block-level only ✓
//
// ROOT CAUSE WAS:
// 1. Page struct missing Properties field in block_parser.go [FIXED]
// 2. No extractPageLevelProperties() function to identify page-level properties [FIXED]
// 3. ParseFile() doesn't call page-level property extraction logic [FIXED]

// TestBug001_ExtractPropertiesFunction - Test individual property line parsing
// This test verifies the existing extractProperties function works correctly
// ✓ PASSES - Basic property extraction works
func TestBug001_ExtractPropertiesFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     map[string]string
	}{
		{
			name:  "simple property",
			input: "tags:: test, properties, demo",
			want:  map[string]string{"tags": "test, properties, demo"},
		},
		{
			name:  "property with spaces",
			input: "alias:: property-test",
			want:  map[string]string{"alias": "property-test"},
		},
		{
			name:  "property with empty value",
			input: "status::",
			want:  map[string]string{"status": ""},
		},
		{
			name:  "property with extra spaces",
			input: "type::   example   ",
			want:  map[string]string{"type": "example"},
		},
		{
			name:  "not a property - no double colon",
			input: "just normal text",
			want:  map[string]string{},
		},
		{
			name:  "not a property - text with single colon",
			input: "time: 2:30 PM",
			want:  map[string]string{},
		},
		{
			name:  "block ID should not be treated as property",
			input: "id:: 64a9c123-456e-789f-abc0-def123456789",
			want:  map[string]string{}, // Block IDs are handled separately
		},
		{
			name:  "property with special characters",
			input: "author:: John Doe <john@example.com>",
			want:  map[string]string{"author": "John Doe <john@example.com>"},
		},
		{
			name:  "property with numbers",
			input: "version:: 1.2.3",
			want:  map[string]string{"version": "1.2.3"},
		},
		{
			name:  "property with underscores and hyphens",
			input: "custom_key:: some-value",
			want:  map[string]string{"custom_key": "some-value"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractProperties(tt.input)
			
			if len(got) != len(tt.want) {
				t.Errorf("extractProperties() returned %d properties, want %d", len(got), len(tt.want))
				t.Errorf("got: %v", got)
				t.Errorf("want: %v", tt.want)
				return
			}
			
			for key, wantValue := range tt.want {
				if gotValue, exists := got[key]; !exists {
					t.Errorf("extractProperties() missing key %q", key)
				} else if gotValue != wantValue {
					t.Errorf("extractProperties() key %q = %q, want %q", key, gotValue, wantValue)
				}
			}
		})
	}
}

// TestBug001_PagePropertiesFieldMissing - Main bug reproduction test
// This test should now PASS with our implementation
func TestBug001_PagePropertiesFieldMissing(t *testing.T) {
	// Read the test properties file
	content, err := os.ReadFile("/Users/rehan/sandbox/seq2b/testdata/library_test_0/pages/Test Properties.md")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}
	
	// Parse the file
	result, err := ParseFile(string(content))
	if err != nil {
		t.Errorf("unexpected parsing error: %v", err)
		return
	}
	
	if result.Page == nil {
		t.Fatal("page is nil")
	}
	
	// Expected page-level properties from the test file
	expectedProps := map[string]string{
		"tags":   "test, properties, demo",
		"alias":  "property-test", 
		"type":   "example",
		"status": "active",
	}
	
	// Test that page has Properties field and it's populated
	if result.Page.Properties == nil {
		t.Error("page.Properties is nil - Properties field missing from Page struct")
		return
	}
	
	// Verify each expected property
	for key, expectedValue := range expectedProps {
		if actualValue, exists := result.Page.Properties[key]; !exists {
			t.Errorf("page property %q not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("page property %q = %q, want %q", key, actualValue, expectedValue)
		}
	}
	
	t.Logf("Successfully parsed %d page properties: %v", len(result.Page.Properties), result.Page.Properties)
}

// TestBug001_Regression - Ensure block-level properties still work
// ✓ SHOULD PASS - Tests existing functionality that should not be broken
func TestBug001_Regression(t *testing.T) {
	// Test that block-level properties still work correctly
	input := `# Test Page

- Block with property
  property:: block-level-value
  - Child block
    child-prop:: child-value`
	
	result, err := ParseFile(input)
	if err != nil {
		t.Errorf("unexpected parsing error: %v", err)
		return
	}
	
	if result.Page == nil {
		t.Fatal("page is nil")
	}
	
	// Verify blocks were parsed correctly
	if len(result.Page.Blocks) == 0 {
		t.Fatal("no blocks parsed")
	}
	
	block := result.Page.Blocks[0]
	
	// Check that block properties are still working
	if block.Properties == nil {
		t.Error("block.Properties is nil")
		return
	}
	
	if propValue, exists := block.Properties["property"]; !exists {
		t.Error("block property 'property' not found")
	} else if propValue != "block-level-value" {
		t.Errorf("block property 'property' = %q, want 'block-level-value'", propValue)
	}
	
	// Check child block properties
	if len(block.Children) == 0 {
		t.Fatal("no child blocks found")
	}
	
	childBlock := block.Children[0]
	if childBlock.Properties == nil {
		t.Error("childBlock.Properties is nil")
		return
	}
	
	if childValue, exists := childBlock.Properties["child-prop"]; !exists {
		t.Error("child block property 'child-prop' not found")
	} else if childValue != "child-value" {
		t.Errorf("child block property 'child-prop' = %q, want 'child-value'", childValue)
	}
}

// TestBug001_PropertyLineIdentification - Test helper functions
// ✓ PASSES - Tests utility functions for identifying page-level properties
func TestBug001_PropertyLineIdentification(t *testing.T) {
	input := `tags:: page-level
alias:: also-page-level

# Header starts content

type:: not-page-level
- Block content`

	// Parse the lines to understand the structure
	lines := ParseFileLines(input)
	
	// Expected: first two lines should be identified as page-level properties
	pageProps := extractPageLevelProperties(lines)
	
	expectedPageProps := map[string]string{
		"tags":  "page-level",
		"alias": "also-page-level",
	}
	
	for key, expectedValue := range expectedPageProps {
		if actualValue, exists := pageProps[key]; !exists {
			t.Errorf("page property %q not found in extracted properties", key)
		} else if actualValue != expectedValue {
			t.Errorf("page property %q = %q, want %q", key, actualValue, expectedValue)
		}
	}
	
	// Should not include the "type" property as it comes after the header
	if typeValue, exists := pageProps["type"]; exists {
		t.Errorf("property 'type' should not be page-level, got: %q", typeValue)
	}
}

// Helper function for testing
func ParseFileLines(content string) []Line {
	lines := []Line{}
	rawLines := strings.Split(content, "\n")
	for i, rawLine := range rawLines {
		line := ParseLine(i+1, strings.TrimSpace(rawLine))
		lines = append(lines, line)
	}
	return lines
}

// BenchmarkBug001_PropertyParsing - Performance test for property parsing
func BenchmarkBug001_PropertyParsing(b *testing.B) {
	content := `tags:: test, properties, demo, performance
alias:: benchmark-test
type:: performance-test
status:: active
priority:: high
category:: testing
author:: test-suite
version:: 1.0.8

# Content starts here

- Some block content
- More blocks`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseFile(content)
		if err != nil {
			b.Fatalf("parsing error: %v", err)
		}
	}
}