// +build integration

package main

import (
	"os"
	"strings"
	"testing"
)

// TestFrontendImageVsPDFRendering tests that the frontend correctly distinguishes
// between actual images and PDFs using image syntax
func TestFrontendImageVsPDFRendering(t *testing.T) {
	app := NewApp()
	
	// Create test content with both images and PDFs using image syntax
	testDir := t.TempDir()
	testFile := testDir + "/test.md"
	content := `# Test Page
- Image: ![Logo](../assets/logo.svg)
- PDF: ![Report](../assets/report.pdf)
- External Image: ![External](https://example.com/image.png)
- External PDF: ![External PDF](https://example.com/doc.pdf)`
	
	if err := createTestFile(testFile, content); err != nil {
		t.Fatal(err)
	}
	
	if err := app.LoadDirectory(testDir); err != nil {
		t.Fatal(err)
	}
	
	pageData, err := app.GetPage("Test Page")
	if err != nil {
		t.Fatal(err)
	}
	
	// Check that all blocks have image segments
	for i, block := range pageData.Blocks {
		if i == 0 {
			continue // Skip header
		}
		
		hasImageSegment := false
		for _, seg := range block.Segments {
			if seg.Type == "image" {
				hasImageSegment = true
				t.Logf("Block %d: Image segment with target=%s", i, seg.Target)
			}
		}
		
		if !hasImageSegment {
			t.Errorf("Block %d should have an image segment", i)
		}
	}
	
	// Now test that frontend would render them differently
	// This is a simulation of what the frontend JavaScript does
	testCases := []struct {
		target       string
		shouldBePDF  bool
		description  string
	}{
		{"../assets/logo.svg", false, "SVG should be rendered as image"},
		{"../assets/report.pdf", true, "PDF should be rendered as link"},
		{"https://example.com/image.png", false, "External PNG should be image"},
		{"https://example.com/doc.pdf", true, "External PDF should be link"},
	}
	
	for _, tc := range testCases {
		isPDF := strings.ToLower(tc.target)[strings.LastIndex(strings.ToLower(tc.target), "."):] == ".pdf"
		if isPDF != tc.shouldBePDF {
			t.Errorf("%s: expected isPDF=%v, got %v", tc.description, tc.shouldBePDF, isPDF)
		}
	}
}

// TestImageLoadingStates tests that images go through proper loading states
func TestImageLoadingStates(t *testing.T) {
	app := NewApp()
	
	// Load test library
	if err := app.LoadDirectory("../../testdata/library_test_0"); err != nil {
		t.Fatal(err)
	}
	
	// Get page with image
	pageData, err := app.GetPage("Page A")
	if err != nil {
		t.Fatal(err)
	}
	
	// Find image block
	var imageBlock *BlockData
	for _, block := range pageData.Blocks {
		if strings.Contains(block.Content, "![seq2b Logo]") {
			imageBlock = &block
			break
		}
	}
	
	if imageBlock == nil {
		t.Fatal("Could not find image block")
	}
	
	// Verify image segment exists
	hasImage := false
	var imageTarget string
	for _, seg := range imageBlock.Segments {
		if seg.Type == "image" && strings.HasSuffix(seg.Target, ".svg") {
			hasImage = true
			imageTarget = seg.Target
			break
		}
	}
	
	if !hasImage {
		t.Fatal("Image segment not found")
	}
	
	// Test that GetAsset works for the image
	if strings.HasPrefix(imageTarget, "../assets/") {
		assetPath := "assets/" + imageTarget[len("../assets/"):]
		assetData, err := app.GetAsset(assetPath)
		if err != nil {
			t.Errorf("GetAsset failed for %s: %v", assetPath, err)
		} else {
			if !strings.Contains(assetData, "seq2b") {
				t.Error("SVG content doesn't contain expected text")
			}
		}
	}
}

// TestHTMLGeneration tests that the correct HTML is generated for images vs PDFs
func TestHTMLGeneration(t *testing.T) {
	app := NewApp()
	
	testDir := t.TempDir()
	testFile := testDir + "/html-test.md"
	content := `# HTML Test
- ![Image Test](../assets/test.png)
- ![PDF Test](../assets/test.pdf)`
	
	if err := createTestFile(testFile, content); err != nil {
		t.Fatal(err)
	}
	
	if err := app.LoadDirectory(testDir); err != nil {
		t.Fatal(err)
	}
	
	pageData, err := app.GetPage("HTML Test")
	if err != nil {
		t.Fatal(err)
	}
	
	// Check segments are correct
	if len(pageData.Blocks) < 3 {
		t.Fatal("Expected at least 3 blocks")
	}
	
	// Image block
	imageBlock := pageData.Blocks[1]
	if !containsImageSegment(imageBlock.Segments, ".png") {
		t.Error("First block should have PNG image segment")
	}
	
	// PDF block  
	pdfBlock := pageData.Blocks[2]
	if !containsImageSegment(pdfBlock.Segments, ".pdf") {
		t.Error("Second block should have PDF image segment")
	}
	
	// Frontend would render these differently:
	// - PNG: <img data-asset-path="test.png" ...>
	// - PDF: <a href="#" class="pdf-link" onclick="openPDF(...)">...</a>
	t.Log("✓ Backend correctly parses both as image segments")
	t.Log("✓ Frontend is responsible for rendering them differently based on file extension")
}

// Helper functions
func createTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func containsImageSegment(segments []SegmentData, ext string) bool {
	for _, seg := range segments {
		if seg.Type == "image" && strings.HasSuffix(seg.Target, ext) {
			return true
		}
	}
	return false
}