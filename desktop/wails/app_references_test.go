package main

import (
	"testing"
	"path/filepath"
)

func TestReferencesDisplay(t *testing.T) {
	app := NewApp()
	testDir := filepath.Join("..", "..", "testdata", "library_test_0", "pages")
	absDir, _ := filepath.Abs(testDir)
	app.LoadDirectory(absDir)

	// Test that Page B has backlinks from Page A
	pageBData, err := app.GetPage("Page B")
	if err != nil {
		t.Fatalf("Failed to get Page B: %v", err)
	}

	if len(pageBData.Backlinks) == 0 {
		t.Error("Page B should have backlinks")
	}

	// Check that Page A is in the backlinks
	foundPageA := false
	for _, backlink := range pageBData.Backlinks {
		if backlink.SourcePage == "Page A" {
			foundPageA = true
			t.Logf("Found backlink from Page A with %d references", len(backlink.BlockIDs))
		}
	}

	if !foundPageA {
		t.Error("Page B should have backlinks from Page A")
	}

	// Test a page without backlinks
	pageWithoutBacklinks, err := app.GetPage("Jul 13th, 2025")
	if err != nil {
		t.Fatalf("Failed to get date page: %v", err)
	}

	if len(pageWithoutBacklinks.Backlinks) > 0 {
		t.Error("Date page should not have backlinks")
	}
}