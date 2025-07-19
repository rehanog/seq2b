package main

import (
	"testing"
	"path/filepath"
)

func TestImageSegmentConversion(t *testing.T) {
	app := NewApp()
	testDir := filepath.Join("..", "..", "testdata", "library_test_0", "pages")
	absDir, _ := filepath.Abs(testDir)
	app.LoadDirectory(absDir)

	// Get Page A which contains an image
	pageData, err := app.GetPage("Page A")
	if err != nil {
		t.Fatalf("Failed to get Page A: %v", err)
	}

	// Find the block with the image
	var imageBlock *BlockData
	for _, block := range pageData.Blocks {
		if containsImage(block) {
			imageBlock = &block
			break
		}
	}

	if imageBlock == nil {
		t.Fatal("No block with image found in Page A")
	}

	t.Logf("Found block with content: %s", imageBlock.Content)
	t.Logf("Block has %d segments", len(imageBlock.Segments))

	// Check if segments contain an image
	hasImageSegment := false
	for i, seg := range imageBlock.Segments {
		t.Logf("Segment %d: type=%s, content=%q, target=%q", 
			i, seg.Type, seg.Content, seg.Target)
		if seg.Type == "image" {
			hasImageSegment = true
			if seg.Target == "" {
				t.Error("Image segment has empty target")
			}
			if seg.Content == "" && seg.Alt == "" {
				t.Error("Image segment has no alt text")
			}
		}
	}

	if !hasImageSegment {
		t.Error("No image segment found in block containing image markdown")
	}

	// Also verify the deprecated HTML rendering
	if imageBlock.HTMLContent != "" {
		t.Logf("HTML content: %s", imageBlock.HTMLContent)
		if !containsString(imageBlock.HTMLContent, "<img") {
			t.Error("HTML content does not contain img tag")
		}
	}
}

func containsImage(block BlockData) bool {
	return containsString(block.Content, "![") && containsString(block.Content, "](")
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Helper function to check children recursively
func checkBlockAndChildren(block BlockData, check func(BlockData) bool) bool {
	if check(block) {
		return true
	}
	for _, child := range block.Children {
		if checkBlockAndChildren(child, check) {
			return true
		}
	}
	return false
}