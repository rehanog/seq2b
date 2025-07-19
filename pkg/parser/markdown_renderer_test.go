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

func TestParseMarkdownSegments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Segment
	}{
		{
			name:  "Plain text",
			input: "Just some plain text",
			expected: []Segment{
				{Type: SegmentText, Content: "Just some plain text"},
			},
		},
		{
			name:  "Bold text",
			input: "This is **bold** text",
			expected: []Segment{
				{Type: SegmentText, Content: "This is "},
				{Type: SegmentBold, Content: "bold"},
				{Type: SegmentText, Content: " text"},
			},
		},
		{
			name:  "Italic text",
			input: "This is *italic* text",
			expected: []Segment{
				{Type: SegmentText, Content: "This is "},
				{Type: SegmentItalic, Content: "italic"},
				{Type: SegmentText, Content: " text"},
			},
		},
		{
			name:  "Link",
			input: "See [[Page A]] for details",
			expected: []Segment{
				{Type: SegmentText, Content: "See "},
				{Type: SegmentLink, Content: "Page A", Target: "Page A"},
				{Type: SegmentText, Content: " for details"},
			},
		},
		{
			name:  "Image with alt text",
			input: "Here's an image: ![Alt text](path/to/image.png)",
			expected: []Segment{
				{Type: SegmentText, Content: "Here's an image: "},
				{Type: SegmentImage, Content: "Alt text", Target: "path/to/image.png", Alt: "Alt text"},
			},
		},
		{
			name:  "Image without alt text",
			input: "Image: ![](image.png)",
			expected: []Segment{
				{Type: SegmentText, Content: "Image: "},
				{Type: SegmentImage, Content: "", Target: "image.png", Alt: ""},
			},
		},
		{
			name:  "Mixed content",
			input: "**Bold** and *italic* with [[link]] and ![image](pic.jpg)",
			expected: []Segment{
				{Type: SegmentBold, Content: "Bold"},
				{Type: SegmentText, Content: " and "},
				{Type: SegmentItalic, Content: "italic"},
				{Type: SegmentText, Content: " with "},
				{Type: SegmentLink, Content: "link", Target: "link"},
				{Type: SegmentText, Content: " and "},
				{Type: SegmentImage, Content: "image", Target: "pic.jpg", Alt: "image"},
			},
		},
		{
			name:  "Multiple images",
			input: "![First](1.png) and ![Second](2.png)",
			expected: []Segment{
				{Type: SegmentImage, Content: "First", Target: "1.png", Alt: "First"},
				{Type: SegmentText, Content: " and "},
				{Type: SegmentImage, Content: "Second", Target: "2.png", Alt: "Second"},
			},
		},
		{
			name:  "Multiple links close together",
			input: "Check [[Page A]] and [[Page B]] for details",
			expected: []Segment{
				{Type: SegmentText, Content: "Check "},
				{Type: SegmentLink, Content: "Page A", Target: "Page A"},
				{Type: SegmentText, Content: " and "},
				{Type: SegmentLink, Content: "Page B", Target: "Page B"},
				{Type: SegmentText, Content: " for details"},
			},
		},
		{
			name:  "Link at end of line",
			input: "See reference [[Page C]]",
			expected: []Segment{
				{Type: SegmentText, Content: "See reference "},
				{Type: SegmentLink, Content: "Page C", Target: "Page C"},
			},
		},
		{
			name:  "PDF link",
			input: "Check this [[../assets/document.pdf]]",
			expected: []Segment{
				{Type: SegmentText, Content: "Check this "},
				{Type: SegmentLink, Content: "../assets/document.pdf", Target: "../assets/document.pdf"},
			},
		},
		{
			name:  "External PDF link",
			input: "External PDF: [[https://example.com/file.pdf]]",
			expected: []Segment{
				{Type: SegmentText, Content: "External PDF: "},
				{Type: SegmentLink, Content: "https://example.com/file.pdf", Target: "https://example.com/file.pdf"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := ParseMarkdownSegments(tt.input)
			
			if len(segments) != len(tt.expected) {
				t.Errorf("Expected %d segments, got %d", len(tt.expected), len(segments))
				return
			}
			
			for i, seg := range segments {
				exp := tt.expected[i]
				if seg.Type != exp.Type {
					t.Errorf("Segment %d: expected type %v, got %v", i, exp.Type, seg.Type)
				}
				if seg.Content != exp.Content {
					t.Errorf("Segment %d: expected content %q, got %q", i, exp.Content, seg.Content)
				}
				if seg.Target != exp.Target {
					t.Errorf("Segment %d: expected target %q, got %q", i, exp.Target, seg.Target)
				}
				if seg.Alt != exp.Alt {
					t.Errorf("Segment %d: expected alt %q, got %q", i, exp.Alt, seg.Alt)
				}
			}
		})
	}
}

func TestRenderToHTML_Images(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple image",
			input:    "![Alt text](image.png)",
			expected: `<img src="image.png" alt="Alt text">`,
		},
		{
			name:     "Image without alt",
			input:    "![](image.png)",
			expected: `<img src="image.png" alt="">`,
		},
		{
			name:     "Image with path",
			input:    "![Diagram](../images/diagram.png)",
			expected: `<img src="../images/diagram.png" alt="Diagram">`,
		},
		{
			name:     "Mixed with other markdown",
			input:    "**Bold** and ![image](pic.jpg) and [[link]]",
			expected: `<b>Bold</b> and <img src="pic.jpg" alt="image"> and <a href="link">link</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderToHTML(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}