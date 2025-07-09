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
		return Line{
			Number:  number,
			Type:    TypeBlock,
			Content: blockText,
		}
	}
	
	// Regular text
	return Line{
		Number:  number,
		Type:    TypeText,
		Content: trimmed,
	}
}

// RenderToHTML converts markdown text to HTML
func RenderToHTML(text string) string {
	if text == "" {
		return ""
	}
	
	html := text
	
	// Convert bold text first: **text** -> <b>text</b>
	boldPattern := regexp.MustCompile(`\*\*(.*?)\*\*`)
	html = boldPattern.ReplaceAllString(html, `<b>$1</b>`)
	
	// Convert italic text: *text* -> <i>text</i> (only single asterisks now)
	italicPattern := regexp.MustCompile(`\*([^*]+?)\*`)
	html = italicPattern.ReplaceAllString(html, `<i>$1</i>`)
	
	// Convert page links: [[page]] -> <a href="page">page</a>
	linkPattern := regexp.MustCompile(`\[\[(.*?)\]\]`)
	html = linkPattern.ReplaceAllString(html, `<a href="$1">$1</a>`)
	
	return html
}