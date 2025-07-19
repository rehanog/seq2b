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
	SegmentImage
	SegmentTag          // #tag
	SegmentBlockRef     // ((block-id))
	SegmentProperty     // key:: value
	SegmentBlockID      // id:: UUID
	SegmentQuery        // {{query}}
	SegmentEmbed        // {{embed}}
	SegmentStrikethrough // ~~text~~
	SegmentHighlight    // ==text== or ^^text^^
)

// Segment represents a parsed text segment
type Segment struct {
	Type    SegmentType
	Content string
	Target  string // For links, the target page; for images, the image path
	Alt     string // For images, the alt text
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
	
	// Convert images: ![alt](src) -> <img src="src" alt="alt">
	imagePattern := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	html = imagePattern.ReplaceAllString(html, `<img src="$2" alt="$1">`)
	
	return html
}

// ParseMarkdownSegments parses markdown text into structured segments
func ParseMarkdownSegments(text string) []Segment {
	if text == "" {
		return []Segment{}
	}
	
	segments := []Segment{}
	remaining := text
	
	// Combined pattern to match all markdown and Logseq features
	// Order matters - more specific patterns first
	pattern := regexp.MustCompile(`(` +
		`\{\{query.*?\}\}|` +                    // {{query}} blocks
		`\{\{embed.*?\}\}|` +                    // {{embed}} blocks
		`\(\([a-fA-F0-9\-]+\)\)|` +              // ((block-id)) references
		`~~.*?~~|` +                             // ~~strikethrough~~
		`==.*?==|` +                             // ==highlight==
		`\^\^.*?\^\^|` +                         // ^^highlight^^
		`#[a-zA-Z0-9\-_/]+|` +                   // #tags
		`\bid::\s*[a-fA-F0-9\-]+|` +             // id:: UUID
		`[a-zA-Z][a-zA-Z0-9\-_]*::\s*[^\n]+|` + // property:: value
		`\*\*.*?\*\*|` +                         // **bold**
		`\*[^*]+?\*|` +                          // *italic*
		`\[([^\]]+)\]\(\[\[([^\]]+)\]\]\)|` +   // [text]([[page]]) - named page link
		`\[([^\]]+)\]\(([^)]+)\)|` +            // [text](url) - markdown link  
		`\[\[.*?\]\]|` +                         // [[page link]]
		`!\[.*?\]\(.*?\)` +                      // ![image](url)
		`)`)
	
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
		
		// Classify the match based on its pattern
		switch {
		case strings.HasPrefix(match, "{{query"):
			segments = append(segments, Segment{
				Type:    SegmentQuery,
				Content: match,
			})
		case strings.HasPrefix(match, "{{embed"):
			segments = append(segments, Segment{
				Type:    SegmentEmbed,
				Content: match,
			})
		case strings.HasPrefix(match, "((") && strings.HasSuffix(match, "))"):
			blockId := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentBlockRef,
				Content: blockId,
				Target:  blockId,
			})
		case strings.HasPrefix(match, "~~") && strings.HasSuffix(match, "~~"):
			content := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentStrikethrough,
				Content: content,
			})
		case strings.HasPrefix(match, "==") && strings.HasSuffix(match, "=="):
			content := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentHighlight,
				Content: content,
			})
		case strings.HasPrefix(match, "^^") && strings.HasSuffix(match, "^^"):
			content := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentHighlight,
				Content: content,
			})
		case strings.HasPrefix(match, "#"):
			tag := match[1:] // Remove the #
			segments = append(segments, Segment{
				Type:    SegmentTag,
				Content: tag,
				Target:  tag,
			})
		case strings.Contains(match, "id::"):
			// Extract the UUID part
			parts := strings.SplitN(match, "::", 2)
			if len(parts) == 2 {
				id := strings.TrimSpace(parts[1])
				segments = append(segments, Segment{
					Type:    SegmentBlockID,
					Content: match,
					Target:  id,
				})
			}
		case strings.Contains(match, "::") && !strings.Contains(match, "id::"):
			// Property
			parts := strings.SplitN(match, "::", 2)
			if len(parts) == 2 {
				segments = append(segments, Segment{
					Type:    SegmentProperty,
					Content: match,
				})
			}
		case strings.HasPrefix(match, "**") && strings.HasSuffix(match, "**"):
			// Bold text
			content := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentBold,
				Content: content,
			})
		case strings.HasPrefix(match, "*") && strings.HasSuffix(match, "*"):
			// Italic text
			content := match[1 : len(match)-1]
			segments = append(segments, Segment{
				Type:    SegmentItalic,
				Content: content,
			})
		case strings.HasPrefix(match, "[") && strings.Contains(match, "]([[") && strings.HasSuffix(match, "]])"):
			// Named page link: [text]([[page]])
			namedLinkPattern := regexp.MustCompile(`\[([^\]]+)\]\(\[\[([^\]]+)\]\]\)`)
			if matches := namedLinkPattern.FindStringSubmatch(match); len(matches) == 3 {
				segments = append(segments, Segment{
					Type:    SegmentLink,
					Content: matches[1], // Display text
					Target:  matches[2], // Page name
				})
			}
		case strings.HasPrefix(match, "[") && strings.Contains(match, "](") && strings.HasSuffix(match, ")") && !strings.HasPrefix(match, "!["):
			// Markdown link: [text](url)
			linkPattern := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
			if matches := linkPattern.FindStringSubmatch(match); len(matches) == 3 {
				segments = append(segments, Segment{
					Type:    SegmentLink,
					Content: matches[1], // Display text
					Target:  matches[2], // URL or path
				})
			}
		case strings.HasPrefix(match, "[[") && strings.HasSuffix(match, "]]"):
			// Link
			target := match[2 : len(match)-2]
			segments = append(segments, Segment{
				Type:    SegmentLink,
				Content: target,
				Target:  target,
			})
		case strings.HasPrefix(match, "!["):
			// Image: ![alt text](path/to/image.png)
			imagePattern := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
			if matches := imagePattern.FindStringSubmatch(match); len(matches) == 3 {
				segments = append(segments, Segment{
					Type:    SegmentImage,
					Alt:     matches[1],
					Target:  matches[2],
					Content: matches[1], // Use alt text as content
				})
			}
		default:
			// Fallback - treat as text
			segments = append(segments, Segment{
				Type:    SegmentText,
				Content: match,
			})
		}
		
		// Continue with remaining text
		remaining = remaining[loc[1]:]
	}
	
	return segments
}