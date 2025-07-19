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
	TodoInfo    TodoInfo            // TODO state and checkbox information
	References  []string            // [[page]] references found in this line
	BlockID     string              // id:: UUID if present
	Properties  map[string]string   // key:: value properties
	Tags        []string            // #tag references
}

// Regular expressions for parsing
var (
	pageRefPattern     = regexp.MustCompile(`\[\[(.*?)\]\]`)
	blockIDPattern     = regexp.MustCompile(`id::\s*([a-zA-Z0-9\-]+)`)
	propertyPattern    = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9\-_]*)::\s*(.*)$`)
	tagPattern         = regexp.MustCompile(`(?:^|\s)#([a-zA-Z0-9\-_/]+)`)
	blockRefPattern    = regexp.MustCompile(`\(\(([a-fA-F0-9\-]+)\)\)`)
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
		
		// Extract block ID if present
		blockID := extractBlockID(blockText)
		
		// Check if this line is a property
		properties := extractProperties(blockText)
		
		// Extract tags
		tags := extractTags(blockText)
		
		return Line{
			Number:     number,
			Type:       TypeBlock,
			Content:    blockText,
			TodoInfo:   todoInfo,
			References: references,
			BlockID:    blockID,
			Properties: properties,
			Tags:       tags,
		}
	}
	
	// Regular text
	// Extract features even from regular text
	references := extractPageReferences(trimmed)
	blockID := extractBlockID(trimmed)
	properties := extractProperties(trimmed)
	tags := extractTags(trimmed)
	
	return Line{
		Number:     number,
		Type:       TypeText,
		Content:    trimmed,
		References: references,
		BlockID:    blockID,
		Properties: properties,
		Tags:       tags,
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

// extractBlockID finds id:: UUID in text
func extractBlockID(text string) string {
	matches := blockIDPattern.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractProperties finds key:: value properties in text
func extractProperties(text string) map[string]string {
	properties := make(map[string]string)
	
	// Check if entire line is a property
	if matches := propertyPattern.FindStringSubmatch(text); len(matches) > 2 {
		key := matches[1]
		value := strings.TrimSpace(matches[2])
		
		// Don't treat block ID as a regular property
		if key != "id" {
			properties[key] = value
		}
	}
	
	return properties
}

// extractTags finds all #tag references in text
func extractTags(text string) []string {
	matches := tagPattern.FindAllStringSubmatch(text, -1)
	tags := make([]string, 0, len(matches))
	
	for _, match := range matches {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}
	
	return tags
}