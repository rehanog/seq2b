package parser

import (
	"testing"
)

func TestNamedLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Segment
	}{
		{
			name:  "Named page link",
			input: "Check [our documentation]([[Page A]]) for details",
			expected: []Segment{
				{Type: SegmentText, Content: "Check "},
				{Type: SegmentLink, Content: "our documentation", Target: "Page A"},
				{Type: SegmentText, Content: " for details"},
			},
		},
		{
			name:  "Named PDF link", 
			input: "Download [the report]([[../assets/report.pdf]])",
			expected: []Segment{
				{Type: SegmentText, Content: "Download "},
				{Type: SegmentLink, Content: "the report", Target: "../assets/report.pdf"},
			},
		},
		{
			name:  "Markdown style link",
			input: "Visit [our website](https://example.com)",
			expected: []Segment{
				{Type: SegmentText, Content: "Visit "},
				{Type: SegmentLink, Content: "our website", Target: "https://example.com"},
			},
		},
		{
			name:  "Mix of link styles",
			input: "See [[Page B]] or [click here]([[Page B]]) or [external](https://example.com)",
			expected: []Segment{
				{Type: SegmentText, Content: "See "},
				{Type: SegmentLink, Content: "Page B", Target: "Page B"},
				{Type: SegmentText, Content: " or "},
				{Type: SegmentLink, Content: "click here", Target: "Page B"},
				{Type: SegmentText, Content: " or "},
				{Type: SegmentLink, Content: "external", Target: "https://example.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := ParseMarkdownSegments(tt.input)
			
			if len(segments) != len(tt.expected) {
				t.Errorf("Expected %d segments, got %d", len(tt.expected), len(segments))
				for i, seg := range segments {
					t.Logf("Segment %d: type=%v, content=%q, target=%q", i, seg.Type, seg.Content, seg.Target)
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
			}
		})
	}
}