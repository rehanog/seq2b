// +build integration

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPDFLinkRendering tests that PDF links are properly detected and rendered
func TestPDFLinkRendering(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-pdf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a page with PDF links
	testFile := filepath.Join(tempDir, "pdf-test.md")
	content := `# PDF Test Page

- Check the [[../assets/report.pdf]] for details
- External PDF: [[https://example.com/document.pdf]]
- Regular page link: [[Page A]]`
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Get the page
	pageData, err := app.GetPage("PDF Test Page")
	if err != nil {
		t.Fatalf("Failed to get page: %v", err)
	}
	
	if len(pageData.Blocks) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(pageData.Blocks))
	}
	
	// Check that PDF links are parsed correctly
	block1 := pageData.Blocks[0]
	if len(block1.Segments) == 0 {
		t.Fatal("Block 1 should have segments")
	}
	
	// Find the PDF link segment
	foundPDFLink := false
	for _, seg := range block1.Segments {
		if seg.Type == "link" && seg.Target == "../assets/report.pdf" {
			foundPDFLink = true
			break
		}
	}
	if !foundPDFLink {
		t.Error("Local PDF link not found in segments")
	}
	
	// Check external PDF link
	block2 := pageData.Blocks[1]
	foundExternalPDF := false
	for _, seg := range block2.Segments {
		if seg.Type == "link" && seg.Target == "https://example.com/document.pdf" {
			foundExternalPDF = true
			break
		}
	}
	if !foundExternalPDF {
		t.Error("External PDF link not found in segments")
	}
	
	// Verify regular page links still work
	block3 := pageData.Blocks[2]
	foundPageLink := false
	for _, seg := range block3.Segments {
		if seg.Type == "link" && seg.Target == "Page A" {
			foundPageLink = true
			break
		}
	}
	if !foundPageLink {
		t.Error("Regular page link not found in segments")
	}
}

// TestPDFAssetLoading tests that PDF assets can be loaded via GetAsset
func TestPDFAssetLoading(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "seq2b-pdf-asset-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create assets directory
	assetsDir := filepath.Join(tempDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}
	
	// Create a simple PDF file (just a placeholder for testing)
	pdfPath := filepath.Join(assetsDir, "test.pdf")
	pdfContent := []byte("%PDF-1.4\n% Test PDF")
	if err := os.WriteFile(pdfPath, pdfContent, 0644); err != nil {
		t.Fatalf("Failed to write PDF file: %v", err)
	}
	
	app := NewApp()
	if err := app.LoadDirectory(tempDir); err != nil {
		t.Fatalf("Failed to load directory: %v", err)
	}
	
	// Test GetAsset for PDF
	assetData, err := app.GetAsset("assets/test.pdf")
	if err != nil {
		t.Fatalf("Failed to get PDF asset: %v", err)
	}
	
	// GetAsset returns base64 for binary files
	if assetData == "" {
		t.Error("PDF asset data should not be empty")
	}
}