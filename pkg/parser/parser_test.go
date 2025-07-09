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

func TestParseLine(t *testing.T) {
	// Table-driven tests - a Go idiom for testing multiple cases
	tests := []struct {
		name        string
		input       string
		wantType    LineType
		wantContent string
		wantLevel   int
	}{
		// Empty lines
		{
			name:        "empty line",
			input:       "",
			wantType:    TypeEmpty,
			wantContent: "",
		},
		{
			name:        "whitespace only",
			input:       "   \t  ",
			wantType:    TypeEmpty,
			wantContent: "",
		},
		// Headers
		{
			name:        "h1 header",
			input:       "# My Header",
			wantType:    TypeHeader,
			wantContent: "My Header",
			wantLevel:   1,
		},
		{
			name:        "h2 header",
			input:       "## Section Two",
			wantType:    TypeHeader,
			wantContent: "Section Two",
			wantLevel:   2,
		},
		{
			name:        "h3 with extra spaces",
			input:       "###   Spaced Header  ",
			wantType:    TypeHeader,
			wantContent: "Spaced Header",
			wantLevel:   3,
		},
		// Block items
		{
			name:        "simple block",
			input:       "- Item one",
			wantType:    TypeBlock,
			wantContent: "Item one",
		},
		{
			name:        "block with extra spaces",
			input:       "  - Indented item  ",
			wantType:    TypeBlock,
			wantContent: "Indented item",
		},
		// Regular text
		{
			name:        "plain text",
			input:       "This is some text",
			wantType:    TypeText,
			wantContent: "This is some text",
		},
		{
			name:        "text with markdown",
			input:       "Text with **bold** and *italic*",
			wantType:    TypeText,
			wantContent: "Text with **bold** and *italic*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLine(1, tt.input)
			
			if got.Type != tt.wantType {
				t.Errorf("ParseLine() type = %v, want %v", got.Type, tt.wantType)
			}
			
			if got.Content != tt.wantContent {
				t.Errorf("ParseLine() content = %q, want %q", got.Content, tt.wantContent)
			}
			
			if got.Type == TypeHeader && got.HeaderLevel != tt.wantLevel {
				t.Errorf("ParseLine() header level = %v, want %v", got.HeaderLevel, tt.wantLevel)
			}
		})
	}
}

// Test that line numbers are correctly assigned
func TestLineNumbers(t *testing.T) {
	testCases := []int{1, 5, 10, 100}
	
	for _, lineNum := range testCases {
		result := ParseLine(lineNum, "some text")
		if result.Number != lineNum {
			t.Errorf("Expected line number %d, got %d", lineNum, result.Number)
		}
	}
}

// Test HTML rendering functionality
func TestRenderToHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "bold text",
			input:    "This has **bold** text",
			expected: "This has <b>bold</b> text",
		},
		{
			name:     "italic text",
			input:    "This has *italic* text",
			expected: "This has <i>italic</i> text",
		},
		{
			name:     "page link",
			input:    "This has a [[page link]]",
			expected: `This has a <a href="page link">page link</a>`,
		},
		{
			name:     "bold and italic",
			input:    "This has **bold** and *italic*",
			expected: "This has <b>bold</b> and <i>italic</i>",
		},
		{
			name:     "multiple links",
			input:    "Link to [[Page A]] and [[Page B]]",
			expected: `Link to <a href="Page A">Page A</a> and <a href="Page B">Page B</a>`,
		},
		{
			name:     "multiple bold",
			input:    "**First** and **second** bold",
			expected: "<b>First</b> and <b>second</b> bold",
		},
		{
			name:     "mixed formatting",
			input:    "**Bold** text with *italic* and [[link]]",
			expected: `<b>Bold</b> text with <i>italic</i> and <a href="link">link</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderToHTML(tt.input)
			if result != tt.expected {
				t.Errorf("RenderToHTML() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Example of a benchmark test
func BenchmarkParseLine(b *testing.B) {
	// This runs the ParseLine function b.N times
	for i := 0; i < b.N; i++ {
		ParseLine(1, "## This is a header with **bold** text")
	}
}

// Benchmark HTML rendering
func BenchmarkRenderToHTML(b *testing.B) {
	text := "This has **bold** and *italic* text with [[links]] to [[other pages]]"
	for i := 0; i < b.N; i++ {
		RenderToHTML(text)
	}
}