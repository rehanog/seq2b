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
	"strings"
	"testing"
)

func TestImageDisplayAfterEditing(t *testing.T) {
	tests := []struct {
		name           string
		initialContent string
		editedContent  string
		expectedType   SegmentType
		expectedTarget string
		expectedAlt    string
	}{
		{
			name:           "Edit image maintains image segment",
			initialContent: "Project overview diagram: ![seq2b Logo](../assets/seq2b-logo.svg)",
			editedContent:  "Project overview diagram: ![seq2b Logo](../assets/seq2b-logo.svg)",
			expectedType:   SegmentImage,
			expectedTarget: "../assets/seq2b-logo.svg",
			expectedAlt:    "seq2b Logo",
		},
		{
			name:           "Edit image path",
			initialContent: "Logo: ![alt text](old-path.png)",
			editedContent:  "Logo: ![alt text](new-path.png)",
			expectedType:   SegmentImage,
			expectedTarget: "new-path.png",
			expectedAlt:    "alt text",
		},
		{
			name:           "Edit alt text",
			initialContent: "Image: ![old alt](path.png)",
			editedContent:  "Image: ![new alt](path.png)",
			expectedType:   SegmentImage,
			expectedTarget: "path.png",
			expectedAlt:    "new alt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create initial block
			line := ParseLine(1, tt.initialContent)
			block := &Block{
				Lines:   []Line{line},
				Content: tt.initialContent,
				Depth:   0,
			}
			block.updateContent()

			// Verify initial parsing
			initialSegments := block.Segments
			if len(initialSegments) < 2 {
				t.Fatalf("Expected at least 2 segments, got %d", len(initialSegments))
			}

			// Find the image segment
			var foundInitialImage bool
			for _, seg := range initialSegments {
				if seg.Type == SegmentImage {
					foundInitialImage = true
					break
				}
			}
			if !foundInitialImage {
				t.Fatal("Initial content should have an image segment")
			}

			// Simulate editing by using SetContent
			block.SetContent(tt.editedContent)

			// Verify segments after editing
			editedSegments := block.Segments
			if len(editedSegments) < 2 {
				t.Fatalf("After editing: expected at least 2 segments, got %d", len(editedSegments))
			}

			// Find and verify the image segment after editing
			var foundEditedImage bool
			var imageSegment *Segment
			for _, seg := range editedSegments {
				if seg.Type == SegmentImage {
					foundEditedImage = true
					imageSegment = &seg
					break
				}
			}

			if !foundEditedImage {
				t.Fatalf("After editing: no image segment found. Segments: %+v", editedSegments)
			}

			// Verify image properties
			if imageSegment.Type != tt.expectedType {
				t.Errorf("Expected segment type %v, got %v", tt.expectedType, imageSegment.Type)
			}
			if imageSegment.Target != tt.expectedTarget {
				t.Errorf("Expected target %q, got %q", tt.expectedTarget, imageSegment.Target)
			}
			if imageSegment.Alt != tt.expectedAlt {
				t.Errorf("Expected alt %q, got %q", tt.expectedAlt, imageSegment.Alt)
			}
		})
	}
}

func TestComplexBlockWithImageEditing(t *testing.T) {
	// Test a more complex case similar to the bug report
	content := `Meeting scheduled for [[Jan 14th, 2025]]`
	
	// Create block
	line := ParseLine(1, content)
	block := &Block{
		Lines:   []Line{line},
		Content: content,
		Depth:   0,
	}
	block.updateContent()

	// Add the image line
	newContent := content + "\n" + "Project overview diagram: ![seq2b Logo](../assets/seq2b-logo.svg)"
	block.SetContent(newContent)

	// Verify the segments include both link and image
	var hasLink, hasImage bool
	var imageSegment *Segment
	
	for _, seg := range block.Segments {
		if seg.Type == SegmentLink {
			hasLink = true
		}
		if seg.Type == SegmentImage {
			hasImage = true
			imageSegment = &seg
		}
	}

	if !hasLink {
		t.Error("Should have a link segment for [[Jan 14th, 2025]]")
	}
	if !hasImage {
		t.Error("Should have an image segment")
	}
	if imageSegment != nil && imageSegment.Target != "../assets/seq2b-logo.svg" {
		t.Errorf("Image target incorrect: got %q", imageSegment.Target)
	}
}

func TestMultiLineBlockParsing(t *testing.T) {
	// Test that multi-line blocks correctly parse all segments
	tests := []struct {
		name           string
		content        string
		expectedSegmentTypes []SegmentType
	}{
		{
			name:    "Two lines with link and image",
			content: "Meeting scheduled for [[Jan 14th, 2025]]\nProject overview diagram: ![seq2b Logo](../assets/seq2b-logo.svg)",
			expectedSegmentTypes: []SegmentType{
				SegmentText,  // "Meeting scheduled for "
				SegmentLink,  // "[[Jan 14th, 2025]]"
				SegmentText,  // "\nProject overview diagram: "
				SegmentImage, // "![seq2b Logo](../assets/seq2b-logo.svg)"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create block by parsing lines
			lines := strings.Split(tt.content, "\n")
			parsedLines := make([]Line, len(lines))
			for i, line := range lines {
				parsedLines[i] = ParseLine(i+1, line)
			}
			
			block := &Block{
				Lines:   parsedLines,
				Content: tt.content,
				Depth:   0,
			}
			block.updateContent()

			// Log segments for debugging
			t.Logf("Content: %q", tt.content)
			t.Logf("Segments: %+v", block.Segments)

			// Check segment count
			if len(block.Segments) != len(tt.expectedSegmentTypes) {
				t.Fatalf("Expected %d segments, got %d", len(tt.expectedSegmentTypes), len(block.Segments))
			}

			// Check segment types
			for i, expectedType := range tt.expectedSegmentTypes {
				if block.Segments[i].Type != expectedType {
					t.Errorf("Segment %d: expected type %v, got %v", i, expectedType, block.Segments[i].Type)
				}
			}
		})
	}
}