// MIT License
//
// Copyright (c) 2025 Rehan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/rehanog/seq2b/pkg/parser"
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
	
	// Load default directory - look for testdata in parent directories
	defaultDir := "../../testdata/pages"
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

// RefreshPages reloads all pages from the current directory
func (a *App) RefreshPages() error {
	if a.currentDir == "" {
		return fmt.Errorf("no directory loaded")
	}
	return a.LoadDirectory(a.currentDir)
}

// GetPage returns page data for display
func (a *App) GetPage(pageName string) (*PageData, error) {
	// Always refresh before getting a page to ensure we have latest content
	if err := a.RefreshPages(); err != nil {
		// Log error but continue with cached version
		fmt.Printf("Warning: failed to refresh pages: %v\n", err)
	}
	
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
	// Refresh pages to get latest list
	if err := a.RefreshPages(); err != nil {
		fmt.Printf("Warning: failed to refresh pages: %v\n", err)
	}
	
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
	TodoState string `json:"todoState"`
	CheckboxState string `json:"checkboxState"`
	Priority string `json:"priority"`
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
			TodoState: string(block.TodoInfo.TodoState),
			CheckboxState: string(block.TodoInfo.CheckboxState),
			Priority: block.TodoInfo.Priority,
		}
	}
	return result
}

// UpdateBlock updates a block's content in a page
func (a *App) UpdateBlock(pageName string, blockID string, newContent string) error {
	page, exists := a.pages[pageName]
	if !exists {
		return fmt.Errorf("page '%s' not found", pageName)
	}
	
	// Find and update the block
	if updateBlockContent(page.Blocks, blockID, newContent) {
		// Save the page back to disk
		if err := a.savePage(page); err != nil {
			return fmt.Errorf("failed to save page: %w", err)
		}
		
		// Reparse the entire directory to update backlinks
		return a.RefreshPages()
	}
	
	return fmt.Errorf("block '%s' not found in page '%s'", blockID, pageName)
}

// savePage writes a page back to disk
func (a *App) savePage(page *parser.Page) error {
	// Reconstruct the markdown content
	content := a.pageToMarkdown(page)
	
	// Determine the filename (convert title to filename)
	filename := parser.TitleToFilename(page.Title)
	filePath := filepath.Join(a.currentDir, filename)
	
	// Write to file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// pageToMarkdown converts a page back to markdown format
func (a *App) pageToMarkdown(page *parser.Page) string {
	var lines []string
	
	// Add title as header
	lines = append(lines, "# " + page.Title)
	lines = append(lines, "")
	
	// Convert blocks to markdown
	a.blocksToMarkdown(page.Blocks, &lines, 0)
	
	return strings.Join(lines, "\n")
}

// blocksToMarkdown recursively converts blocks to markdown lines
func (a *App) blocksToMarkdown(blocks []*parser.Block, lines *[]string, depth int) {
	for _, block := range blocks {
		// Create indentation
		indent := strings.Repeat("  ", depth)
		
		// Add the block content with proper indentation
		blockLines := strings.Split(block.Content, "\n")
		for i, line := range blockLines {
			if i == 0 {
				*lines = append(*lines, indent + "- " + line)
			} else {
				*lines = append(*lines, indent + "  " + line)
			}
		}
		
		// Process children
		if len(block.Children) > 0 {
			a.blocksToMarkdown(block.Children, lines, depth+1)
		}
	}
}

// updateBlockContent recursively searches for and updates a block
func updateBlockContent(blocks []*parser.Block, targetID string, newContent string) bool {
	for _, block := range blocks {
		if block.ID == targetID {
			block.SetContent(newContent)
			return true
		}
		// Check children
		if updateBlockContent(block.Children, targetID, newContent) {
			return true
		}
	}
	return false
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
