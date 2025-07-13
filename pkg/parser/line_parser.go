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
	"regexp"
	"strings"
)

// LineType represents the type of markdown line
type LineType int

const (
	TypeEmpty LineType = iota
	TypeHeader
	TypeText
	TypeBlock // Changed from TypeList - represents a Logseq block
)

// Line represents a parsed line from the markdown file
type Line struct {
	Number      int
	Type        LineType
	Content     string
	HeaderLevel int // Only used for headers
	
	// Parsed data - populated during line parsing
	TodoInfo    TodoInfo  // TODO state and checkbox information
	References  []string  // [[page]] references found in this line
}

// Regular expressions for parsing
var (
	pageRefPattern = regexp.MustCompile(`\[\[(.*?)\]\]`)
)

// ParseLine analyzes a single line and returns its type and content
func ParseLine(number int, line string) Line {
	trimmed := strings.TrimSpace(line)
	
	// Empty line
	if trimmed == "" {
		return Line{Number: number, Type: TypeEmpty}
	}
	
	// Header (starts with #)
	if strings.HasPrefix(trimmed, "#") {
		level := 0
		for _, ch := range trimmed {
			if ch == '#' {
				level++
			} else {
				break
			}
		}
		// Extract header text (remove # and trim)
		headerText := strings.TrimSpace(trimmed[level:])
		return Line{
			Number:      number,
			Type:        TypeHeader,
			Content:     headerText,
			HeaderLevel: level,
		}
	}
	
	// Block (starts with -)
	if strings.HasPrefix(trimmed, "-") {
		blockText := strings.TrimSpace(trimmed[1:])
		
		// Parse TODO information from the block content
		todoInfo := ParseTodoInfo(blockText)
		
		// Extract page references
		references := extractPageReferences(blockText)
		
		return Line{
			Number:     number,
			Type:       TypeBlock,
			Content:    blockText,
			TodoInfo:   todoInfo,
			References: references,
		}
	}
	
	// Regular text
	// Extract page references even from regular text
	references := extractPageReferences(trimmed)
	
	return Line{
		Number:     number,
		Type:       TypeText,
		Content:    trimmed,
		References: references,
	}
}

// extractPageReferences finds all [[page]] references in text
func extractPageReferences(text string) []string {
	matches := pageRefPattern.FindAllStringSubmatch(text, -1)
	references := make([]string, 0, len(matches))
	
	for _, match := range matches {
		if len(match) > 1 {
			references = append(references, match[1])
		}
	}
	
	return references
}

// IsPageReference checks if a reference is a page (not a date)
func IsPageReference(ref string) bool {
	return !IsDatePage(ref)
}