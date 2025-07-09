package main

import (
	"testing"
	
	"fyne.io/fyne/v2/test"
	"github.com/rehan/logseq-go/internal/parser"
)

func TestGUICreation(t *testing.T) {
	// Create test app
	app := test.NewApp()
	window := app.NewWindow("Test")
	
	gui := &LogseqGUI{
		window:    window,
		pages:     make(map[string]*parser.Page),
	}
	
	// Test UI setup
	gui.setupUI()
	if gui.blockView == nil {
		t.Error("Block view not created")
	}
	if gui.statusBar == nil {
		t.Error("Status bar not created")
	}
	
	// Test content creation
	content := gui.createContent()
	if content == nil {
		t.Error("Content not created")
	}
}

func TestFormatBlock(t *testing.T) {
	gui := &LogseqGUI{}
	
	// Create test block
	block := &parser.Block{
		Content: "Test block",
		Children: []*parser.Block{
			{Content: "Child 1"},
			{Content: "Child 2"},
		},
	}
	
	result := gui.formatBlock(block, "")
	expected := "• Test block\n  • Child 1\n  • Child 2\n"
	
	if result != expected {
		t.Errorf("formatBlock() = %q, want %q", result, expected)
	}
}