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

// Test markdown element parsing
func TestMarkdownElements(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantBold     bool
		wantItalic   bool
		wantLink     string
	}{
		{
			name:       "plain text",
			input:      "Just plain text",
			wantBold:   false,
			wantItalic: false,
			wantLink:   "",
		},
		{
			name:       "bold text",
			input:      "This has **bold** text",
			wantBold:   true,
			wantItalic: false,
			wantLink:   "",
		},
		{
			name:       "italic text",
			input:      "This has *italic* text",
			wantBold:   false,
			wantItalic: true,
			wantLink:   "",
		},
		{
			name:       "page link",
			input:      "This has a [[page link]]",
			wantBold:   false,
			wantItalic: false,
			wantLink:   "page link",
		},
		{
			name:       "bold and italic",
			input:      "This has **bold** and *italic*",
			wantBold:   true,
			wantItalic: true,
			wantLink:   "",
		},
		{
			name:       "link with spaces",
			input:      "Link to [[My Important Page]]",
			wantBold:   false,
			wantItalic: false,
			wantLink:   "My Important Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := ParseLine(1, tt.input)
			
			if len(line.Elements) == 0 {
				t.Fatalf("Expected elements to be parsed, got none")
			}
			
			elem := line.Elements[0] // We know our simple parser returns one element
			
			if elem.Bold != tt.wantBold {
				t.Errorf("Expected bold=%v, got %v", tt.wantBold, elem.Bold)
			}
			
			if elem.Italic != tt.wantItalic {
				t.Errorf("Expected italic=%v, got %v", tt.wantItalic, elem.Italic)
			}
			
			if elem.Link != tt.wantLink {
				t.Errorf("Expected link=%q, got %q", tt.wantLink, elem.Link)
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