package parser

import (
	"testing"
)

func TestLogseqPDFLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Segment
	}{
		{
			name:  "Simple PDF link",
			input: "Download ![report](../assets/report.pdf)",
			expected: []Segment{
				{Type: SegmentText, Content: "Download "},
				{Type: SegmentImage, Content: "report", Target: "../assets/report.pdf", Alt: "report"},
			},
		},
		{
			name:  "PDF with descriptive name",
			input: "See the ![Annual Report 2024](../assets/annual-report-2024.pdf) for details",
			expected: []Segment{
				{Type: SegmentText, Content: "See the "},
				{Type: SegmentImage, Content: "Annual Report 2024", Target: "../assets/annual-report-2024.pdf", Alt: "Annual Report 2024"},
				{Type: SegmentText, Content: " for details"},
			},
		},
		{
			name:  "External PDF",
			input: "Reference: ![W3C Spec](https://www.w3.org/spec.pdf)",
			expected: []Segment{
				{Type: SegmentText, Content: "Reference: "},
				{Type: SegmentImage, Content: "W3C Spec", Target: "https://www.w3.org/spec.pdf", Alt: "W3C Spec"},
			},
		},
		{
			name:  "Mixed images and PDFs",
			input: "Logo: ![logo](../assets/logo.png) and docs: ![manual](../assets/manual.pdf)",
			expected: []Segment{
				{Type: SegmentText, Content: "Logo: "},
				{Type: SegmentImage, Content: "logo", Target: "../assets/logo.png", Alt: "logo"},
				{Type: SegmentText, Content: " and docs: "},
				{Type: SegmentImage, Content: "manual", Target: "../assets/manual.pdf", Alt: "manual"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := ParseMarkdownSegments(tt.input)
			
			if len(segments) != len(tt.expected) {
				t.Errorf("Expected %d segments, got %d", len(tt.expected), len(segments))
				for i, seg := range segments {
					t.Logf("Segment %d: type=%v, content=%q, target=%q, alt=%q", 
						i, seg.Type, seg.Content, seg.Target, seg.Alt)
				}
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