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

// SegmentType represents different types of text segments
type SegmentType int

const (
	SegmentText SegmentType = iota
	SegmentBold
	SegmentItalic
	SegmentLink
)

// Segment represents a parsed text segment
type Segment struct {
	Type    SegmentType
	Content string
	Target  string // For links, the target page
}

// RenderToHTML converts markdown text to HTML
// NOTE: This function is temporarily kept for compatibility
// but will be removed once frontend rendering is implemented
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

// ParseMarkdownSegments parses markdown text into structured segments
func ParseMarkdownSegments(text string) []Segment {
	if text == "" {
		return []Segment{}
	}
	
	segments := []Segment{}
	remaining := text
	
	// Combined pattern to match bold, italic, or links
	pattern := regexp.MustCompile(`(\*\*.*?\*\*|\*[^*]+?\*|\[\[.*?\]\])`)
	
	for {
		loc := pattern.FindStringIndex(remaining)
		if loc == nil {
			// No more patterns, add remaining text
			if remaining != "" {
				segments = append(segments, Segment{
					Type:    SegmentText,
					Content: remaining,
				})
			}
			break
		}
		
		// Add text before the match
		if loc[0] > 0 {
			segments = append(segments, Segment{
				Type:    SegmentText,
				Content: remaining[:loc[0]],
			})
		}
		
		// Extract and classify the match
		match := remaining[loc[0]:loc[1]]
		
		if strings.HasPrefix(match, "**") && strings.HasSuffix(match, "**") {
			// Bold text
			content := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentBold,
				Content: content,
			})
		} else if strings.HasPrefix(match, "*") && strings.HasSuffix(match, "*") {
			// Italic text
			content := match[1 : len(match)-1]
			segments = append(segments, Segment{
				Type:    SegmentItalic,
				Content: content,
			})
		} else if strings.HasPrefix(match, "[[") && strings.HasSuffix(match, "]]") {
			// Link
			target := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentLink,
				Content: target,
				Target:  target,
			})
		}
		
		// Continue with remaining text
		remaining = remaining[loc[1]:]
	}
	
	return segments
}