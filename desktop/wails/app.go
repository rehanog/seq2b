package main

import (
	"context"
	"fmt"
	"path/filepath"
	
	"github.com/rehan/logseq-go/pkg/parser"
)

// App struct
type App struct {
	ctx context.Context
	pages map[string]*parser.Page
	backlinks *parser.BacklinkIndex
	currentDir string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		pages: make(map[string]*parser.Page),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Load default directory
	defaultDir := "testdata/pages"
	if absDir, err := filepath.Abs(defaultDir); err == nil {
		a.LoadDirectory(absDir)
	}
}

// LoadDirectory loads all markdown files from a directory
func (a *App) LoadDirectory(dirPath string) error {
	a.currentDir = dirPath
	
	result, err := parser.ParseDirectory(dirPath)
	if err != nil {
		return fmt.Errorf("error parsing directory: %w", err)
	}
	
	a.pages = result.Pages
	a.backlinks = result.Backlinks
	
	return nil
}

// GetPage returns page data for display
func (a *App) GetPage(pageName string) (*PageData, error) {
	page, exists := a.pages[pageName]
	if !exists {
		return nil, fmt.Errorf("page '%s' not found", pageName)
	}
	
	// Get backlinks for this page
	backlinks := a.backlinks.GetBacklinks(pageName)
	
	return &PageData{
		Name: pageName,
		Title: page.Title,
		Blocks: convertBlocks(page.Blocks),
		Backlinks: convertBacklinks(backlinks),
	}, nil
}

// GetPageList returns all available pages
func (a *App) GetPageList() []string {
	pages := make([]string, 0, len(a.pages))
	for name := range a.pages {
		pages = append(pages, name)
	}
	return pages
}

// GetBacklinks returns backlinks for a page
func (a *App) GetBacklinks(pageName string) map[string][]string {
	backlinks := a.backlinks.GetBacklinks(pageName)
	result := make(map[string][]string)
	
	for sourcePage, refs := range backlinks {
		blockIDs := make([]string, len(refs))
		for i, ref := range refs {
			blockIDs[i] = ref.BlockID
		}
		result[sourcePage] = blockIDs
	}
	
	return result
}

// PageData represents page data for frontend
type PageData struct {
	Name string `json:"name"`
	Title string `json:"title"`
	Blocks []BlockData `json:"blocks"`
	Backlinks []BacklinkData `json:"backlinks"`
}

// BlockData represents block data for frontend
type BlockData struct {
	ID string `json:"id"`
	Content string `json:"content"`
	HTMLContent string `json:"htmlContent"`
	Depth int `json:"depth"`
	Children []BlockData `json:"children"`
}

// BacklinkData represents backlink data for frontend
type BacklinkData struct {
	SourcePage string `json:"sourcePage"`
	BlockIDs []string `json:"blockIds"`
	Count int `json:"count"`
}

// Helper functions to convert internal types to frontend types
func convertBlocks(blocks []*parser.Block) []BlockData {
	result := make([]BlockData, len(blocks))
	for i, block := range blocks {
		result[i] = BlockData{
			ID: block.ID,
			Content: block.Content,
			HTMLContent: block.RenderHTML(),
			Depth: block.Depth,
			Children: convertBlocks(block.Children),
		}
	}
	return result
}

func convertBacklinks(backlinks map[string][]parser.BlockReference) []BacklinkData {
	result := make([]BacklinkData, 0, len(backlinks))
	
	for sourcePage, refs := range backlinks {
		blockIDs := make([]string, len(refs))
		for i, ref := range refs {
			blockIDs[i] = ref.BlockID
		}
		
		result = append(result, BacklinkData{
			SourcePage: sourcePage,
			BlockIDs: blockIDs,
			Count: len(refs),
		})
	}
	
	return result
}
