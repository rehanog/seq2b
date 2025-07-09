package main

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
		// List items
		{
			name:        "simple list item",
			input:       "- Item one",
			wantType:    TypeList,
			wantContent: "Item one",
		},
		{
			name:        "list with extra spaces",
			input:       "  - Indented item  ",
			wantType:    TypeList,
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
			got := parseLine(1, tt.input)
			
			if got.Type != tt.wantType {
				t.Errorf("parseLine() type = %v, want %v", got.Type, tt.wantType)
			}
			
			if got.Content != tt.wantContent {
				t.Errorf("parseLine() content = %q, want %q", got.Content, tt.wantContent)
			}
			
			if got.Type == TypeHeader && got.HeaderLevel != tt.wantLevel {
				t.Errorf("parseLine() header level = %v, want %v", got.HeaderLevel, tt.wantLevel)
			}
		})
	}
}

// Test that line numbers are correctly assigned
func TestLineNumbers(t *testing.T) {
	testCases := []int{1, 5, 10, 100}
	
	for _, lineNum := range testCases {
		result := parseLine(lineNum, "some text")
		if result.Number != lineNum {
			t.Errorf("Expected line number %d, got %d", lineNum, result.Number)
		}
	}
}

// Example of a benchmark test
func BenchmarkParseLine(b *testing.B) {
	// This runs the parseLine function b.N times
	for i := 0; i < b.N; i++ {
		parseLine(1, "## This is a header with **bold** text")
	}
}