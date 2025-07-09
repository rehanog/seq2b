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
	TypeList
)

// MarkdownElement represents a formatted piece of text
type MarkdownElement struct {
	Text   string
	Bold   bool
	Italic bool
	Link   string // Empty if not a link, contains page name if it's a [[link]]
}

// Line represents a parsed line from the markdown file
type Line struct {
	Number      int
	Type        LineType
	Content     string
	HeaderLevel int // Only used for headers
	Elements    []MarkdownElement // Parsed markdown elements
}

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
		parsed := Line{
			Number:      number,
			Type:        TypeHeader,
			Content:     headerText,
			HeaderLevel: level,
		}
		parsed.Elements = parseMarkdownElements(headerText)
		return parsed
	}
	
	// List item (starts with -)
	if strings.HasPrefix(trimmed, "-") {
		listText := strings.TrimSpace(trimmed[1:])
		parsed := Line{
			Number:  number,
			Type:    TypeList,
			Content: listText,
		}
		parsed.Elements = parseMarkdownElements(listText)
		return parsed
	}
	
	// Regular text
	parsed := Line{
		Number:  number,
		Type:    TypeText,
		Content: trimmed,
	}
	parsed.Elements = parseMarkdownElements(trimmed)
	return parsed
}

// parseMarkdownElements parses a string and extracts markdown formatting
func parseMarkdownElements(text string) []MarkdownElement {
	if text == "" {
		return nil
	}
	
	elements := []MarkdownElement{}
	
	// For now, let's implement a simple approach that handles one type at a time
	// We'll use regex to find patterns and build elements
	
	// Patterns for different markdown elements
	boldPattern := regexp.MustCompile(`\*\*(.*?)\*\*`)
	linkPattern := regexp.MustCompile(`\[\[(.*?)\]\]`)
	
	// Start with the full text as a single element
	current := MarkdownElement{Text: text}
	
	// Check for bold text first
	current.Bold = boldPattern.MatchString(text)
	
	// Check for italic text (single asterisk, but not adjacent to another asterisk)
	// Use negative lookbehind and lookahead to avoid matching bold markers
	italicPattern := regexp.MustCompile(`(?:^|[^*])\*([^*]+?)\*(?:[^*]|$)`)
	current.Italic = italicPattern.MatchString(text)
	
	// Check for links
	if linkPattern.MatchString(text) {
		matches := linkPattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			current.Link = matches[1]
		}
	}
	
	elements = append(elements, current)
	return elements
}