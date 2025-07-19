// +build integration

package main

import (
	"os"
	"strings"
	"testing"
)

// SimulateFrontendRendering simulates what the frontend JavaScript does
// This helps us catch rendering issues before they happen
func TestSimulateFrontendRendering(t *testing.T) {
	app := NewApp()
	
	// Load the actual test library
	if err := app.LoadDirectory("../../testdata/library_test_0"); err != nil {
		t.Fatal(err)
	}
	
	pageData, err := app.GetPage("Page A")
	if err != nil {
		t.Fatal(err)
	}
	
	// Find blocks with images and PDFs
	var imageCount, pdfCount int
	
	for _, block := range pageData.Blocks {
		for _, seg := range block.Segments {
			if seg.Type == "image" {
				// Simulate frontend logic
				if strings.ToLower(seg.Target)[strings.LastIndex(strings.ToLower(seg.Target), "."):] == ".pdf" {
					pdfCount++
					t.Logf("PDF (as image): %s -> would render as link", seg.Target)
				} else {
					imageCount++
					t.Logf("Image: %s -> would render as img tag", seg.Target)
					
					// Check if it needs asset loading
					if strings.HasPrefix(seg.Target, "../assets/") {
						assetPath := seg.Target[len("../assets/"):]
						t.Logf("  -> Would load via GetAsset: assets/%s", assetPath)
						
						// Actually test GetAsset
						fullPath := "assets/" + assetPath
						data, err := app.GetAsset(fullPath)
						if err != nil {
							t.Errorf("  -> GetAsset would FAIL: %v", err)
						} else {
							t.Logf("  -> GetAsset OK, data length: %d", len(data))
						}
					}
				}
			}
		}
	}
	
	if imageCount == 0 {
		t.Error("No images found - this would result in no images being displayed")
	}
	
	if pdfCount == 0 {
		t.Error("No PDFs found - PDF detection might be broken")
	}
	
	t.Logf("Summary: %d images, %d PDFs", imageCount, pdfCount)
}

// TestActualRenderingScenario tests the exact scenario from Page A
func TestActualRenderingScenario(t *testing.T) {
	app := NewApp()
	
	// Create a page exactly like Page A
	testDir := t.TempDir()
	testFile := testDir + "/page-a-copy.md"
	content := `# Page A Copy
- Project overview diagram: ![seq2b Logo](../assets/seq2b-logo.svg)
- Project documentation: ![test-sample](../assets/test-sample.pdf)`
	
	if err := createTestFile(testFile, content); err != nil {
		t.Fatal(err)
	}
	
	// Create a mock assets directory structure
	assetsDir := testDir + "/../assets"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		// Try alternative approach - load from pages dir
		if err := app.LoadDirectory(testDir); err != nil {
			t.Fatal(err)
		}
	}
	
	pageData, err := app.GetPage("Page A Copy")
	if err != nil {
		t.Fatal(err)
	}
	
	// Test each block
	for i, block := range pageData.Blocks {
		if i == 0 {
			continue // Skip header
		}
		
		t.Logf("\nBlock %d: %s", i, block.Content)
		
		for _, seg := range block.Segments {
			if seg.Type == "image" {
				isPDF := strings.HasSuffix(strings.ToLower(seg.Target), ".pdf")
				
				t.Logf("  Segment: type=%s, target=%s, content=%s", seg.Type, seg.Target, seg.Content)
				t.Logf("  -> Is PDF? %v", isPDF)
				
				if isPDF {
					t.Log("  -> Frontend would render as: <a class='pdf-link'>...</a>")
				} else {
					t.Log("  -> Frontend would render as: <img ...>")
					
					if strings.HasPrefix(seg.Target, "../assets/") {
						t.Log("  -> Would use data-asset-path for loading")
					}
				}
			}
		}
	}
}

// TestJavaScriptLogic tests the exact JavaScript logic we use
func TestJavaScriptLogic(t *testing.T) {
	testCases := []struct {
		segmentType   string
		segmentTarget string
		expectedHTML  string
	}{
		{
			segmentType:   "image",
			segmentTarget: "../assets/logo.svg",
			expectedHTML:  "img with data-asset-path",
		},
		{
			segmentType:   "image", 
			segmentTarget: "../assets/report.pdf",
			expectedHTML:  "a with pdf-link class",
		},
		{
			segmentType:   "image",
			segmentTarget: "../assets/report.PDF", // Test case sensitivity
			expectedHTML:  "a with pdf-link class",
		},
		{
			segmentType:   "link",
			segmentTarget: "../assets/report.pdf",
			expectedHTML:  "a with pdf-link class", // Links to PDFs also become PDF links
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.segmentTarget, func(t *testing.T) {
			// Simulate the JavaScript switch statement
			var result string
			
			switch tc.segmentType {
			case "link":
				if strings.ToLower(tc.segmentTarget)[strings.LastIndex(strings.ToLower(tc.segmentTarget), "."):] == ".pdf" {
					result = "a with pdf-link class"
				} else {
					result = "a with page-link class"
				}
			case "image":
				if strings.ToLower(tc.segmentTarget)[strings.LastIndex(strings.ToLower(tc.segmentTarget), "."):] == ".pdf" {
					result = "a with pdf-link class"
				} else if strings.HasPrefix(tc.segmentTarget, "../assets/") {
					result = "img with data-asset-path"
				} else {
					result = "img with src"
				}
			}
			
			if result != tc.expectedHTML {
				t.Errorf("Expected %s, got %s", tc.expectedHTML, result)
			}
		})
	}
}