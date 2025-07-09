package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	
	"github.com/rehan/logseq-go/internal/parser"
)

type LogseqGUI struct {
	window      fyne.Window
	blockView   *widget.Label
	statusBar   *widget.Label
	
	// Data
	pages       map[string]*parser.Page
	backlinks   *parser.BacklinkIndex
	currentPage string
}

func main() {
	// Parse command line arguments
	var startPage string
	flag.StringVar(&startPage, "page", "", "Starting page name (default: Page A)")
	flag.Parse()
	
	// Set defaults
	if startPage == "" {
		startPage = "Page A"
	}
	
	// Default directory - can be overridden by argument
	directory := "testdata/pages"
	if flag.NArg() > 0 {
		directory = flag.Arg(0)
	}
	
	// Make directory absolute
	absDir, err := filepath.Abs(directory)
	if err != nil {
		log.Fatalf("Error resolving directory: %v", err)
	}
	
	myApp := app.New()
	myWindow := myApp.NewWindow("Logseq Go")
	
	gui := &LogseqGUI{
		window:    myWindow,
		pages:     make(map[string]*parser.Page),
	}
	
	gui.setupUI()
	
	// Load the directory
	if err := gui.loadDirectory(absDir); err != nil {
		log.Fatalf("Error loading directory: %v", err)
	}
	
	// Display the starting page
	gui.displayPage(startPage)
	
	myWindow.SetContent(gui.createContent())
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func (g *LogseqGUI) setupUI() {
	// Create block view
	g.blockView = widget.NewLabel("Loading...")
	g.blockView.Wrapping = fyne.TextWrapWord
	
	// Create status bar
	g.statusBar = widget.NewLabel("Loading pages...")
}

func (g *LogseqGUI) createContent() fyne.CanvasObject {
	// Main layout - just the page content and status bar
	return container.NewBorder(
		nil,
		g.statusBar,
		nil, nil,
		container.NewScroll(g.blockView),
	)
}

func (g *LogseqGUI) loadDirectory(dirPath string) error {
	g.statusBar.SetText(fmt.Sprintf("Loading %s...", dirPath))
	
	// Parse directory
	result, err := parser.ParseDirectory(dirPath)
	if err != nil {
		return fmt.Errorf("error parsing directory: %w", err)
	}
	
	// Store results
	g.pages = result.Pages
	g.backlinks = result.Backlinks
	
	g.statusBar.SetText(fmt.Sprintf("Loaded %d pages", len(g.pages)))
	return nil
}

func (g *LogseqGUI) displayPage(pageName string) {
	page, exists := g.pages[pageName]
	if !exists {
		g.blockView.SetText(fmt.Sprintf("Page '%s' not found", pageName))
		g.statusBar.SetText(fmt.Sprintf("Page '%s' not found", pageName))
		return
	}
	
	g.currentPage = pageName
	
	// Format block content
	content := fmt.Sprintf("# %s\n\n", pageName)
	
	// Show blocks
	for _, block := range page.Blocks {
		content += g.formatBlock(block, "")
	}
	
	// Show backlinks
	backlinks := g.backlinks.GetBacklinks(pageName)
	if len(backlinks) > 0 {
		content += "\n\n## Backlinks\n"
		for sourcePage, refs := range backlinks {
			content += fmt.Sprintf("← %s (%d references)\n", sourcePage, len(refs))
		}
	}
	
	g.blockView.SetText(content)
	g.statusBar.SetText(fmt.Sprintf("Viewing: %s", pageName))
}

func (g *LogseqGUI) formatBlock(block *parser.Block, indent string) string {
	result := fmt.Sprintf("%s• %s\n", indent, block.Content)
	
	// Add children
	for _, child := range block.Children {
		result += g.formatBlock(child, indent+"  ")
	}
	
	return result
}