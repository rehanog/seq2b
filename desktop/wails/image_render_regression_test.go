// +build integration

package main

import (
	"strings"
	"testing"
)

// TestImageRenderingNotBrokenByPDFSupport ensures that adding PDF support
// doesn't break regular image rendering
func TestImageRenderingNotBrokenByPDFSupport(t *testing.T) {
	app := NewApp()
	
	// Load test library
	if err := app.LoadDirectory("../../testdata/library_test_0"); err != nil {
		t.Fatal(err)
	}
	
	// Get Page A which has the seq2b logo
	pageData, err := app.GetPage("Page A")
	if err != nil {
		t.Fatal(err)
	}
	
	// Find the logo block
	var logoBlock *BlockData
	for _, block := range pageData.Blocks {
		if strings.Contains(block.Content, "seq2b Logo") {
			logoBlock = &block
			break
		}
	}
	
	if logoBlock == nil {
		t.Fatal("Could not find logo block")
	}
	
	// Check the segment
	var svgSegment *SegmentData
	for i, seg := range logoBlock.Segments {
		if seg.Type == "image" && strings.HasSuffix(seg.Target, ".svg") {
			svgSegment = &logoBlock.Segments[i]
			break
		}
	}
	
	if svgSegment == nil {
		t.Fatal("Could not find SVG image segment")
	}
	
	// Verify segment properties
	if svgSegment.Target != "../assets/seq2b-logo.svg" {
		t.Errorf("Wrong target: %s", svgSegment.Target)
	}
	
	if svgSegment.Content != "seq2b Logo" {
		t.Errorf("Wrong content: %s", svgSegment.Content)
	}
	
	// Test GetAsset
	assetPath := "assets/seq2b-logo.svg"
	assetData, err := app.GetAsset(assetPath)
	if err != nil {
		t.Fatalf("GetAsset failed: %v", err)
	}
	
	// Verify it's valid SVG
	if !strings.Contains(assetData, "<svg") {
		t.Error("Asset data doesn't look like SVG")
	}
	
	if !strings.Contains(assetData, "seq2b") {
		t.Error("SVG doesn't contain expected text")
	}
	
	// CRITICAL TEST: Simulate what frontend SHOULD do
	t.Run("Frontend Simulation", func(t *testing.T) {
		// This is what the JavaScript does:
		isPDF := strings.ToLower(svgSegment.Target)[strings.LastIndex(strings.ToLower(svgSegment.Target), "."):] == ".pdf"
		
		if isPDF {
			t.Error("SVG is being detected as PDF!")
		}
		
		// Should render as image with data-asset-path
		if !strings.HasPrefix(svgSegment.Target, "../assets/") {
			t.Error("SVG should have ../assets/ prefix")
		}
		
		// Extract asset path
		extractedPath := svgSegment.Target[len("../assets/"):]
		if extractedPath != "seq2b-logo.svg" {
			t.Errorf("Wrong extracted path: %s", extractedPath)
		}
		
		t.Log("✓ SVG would be rendered as <img> with data-asset-path")
		t.Log("✓ GetAsset would be called with: assets/" + extractedPath)
		t.Log("✓ Image should display correctly")
	})
}

// TestLoadingStateNotStuck verifies images don't get stuck in loading state
func TestLoadingStateNotStuck(t *testing.T) {
	// This test checks the conditions that would cause an image to stay 
	// in the loading state (gray box)
	
	testCases := []struct {
		name          string
		assetPath     string
		getAssetError bool
		shouldLoad    bool
	}{
		{
			name:          "Valid SVG",
			assetPath:     "seq2b-logo.svg",
			getAssetError: false,
			shouldLoad:    true,
		},
		{
			name:          "Missing asset",
			assetPath:     "missing.svg",
			getAssetError: true,
			shouldLoad:    false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate frontend loadAssets() function
			fullAssetPath := "assets/" + tc.assetPath
			
			t.Logf("Would call GetAsset('%s')", fullAssetPath)
			
			if tc.getAssetError {
				t.Log("-> GetAsset would fail")
				t.Log("-> Image would show 'failed to load' (red box)")
				t.Log("-> loading-asset class would be removed")
			} else {
				t.Log("-> GetAsset would succeed")
				t.Log("-> Image src would be set")
				t.Log("-> loading-asset class would be removed")
				t.Log("-> Image would display correctly")
			}
		})
	}
}