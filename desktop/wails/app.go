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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	
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
	defaultDir := "../../testdata/library_test_0/pages"
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
		// Auto-create the page
		if parser.IsDatePage(pageName) {
			// Create date page with special handling
			if err := a.createDatePage(pageName); err != nil {
				return nil, fmt.Errorf("failed to create date page: %w", err)
			}
		} else {
			// Create regular page
			if err := a.createPage(pageName); err != nil {
				return nil, fmt.Errorf("failed to create page: %w", err)
			}
		}
		
		// Refresh pages to load the new page
		if err := a.RefreshPages(); err != nil {
			return nil, fmt.Errorf("failed to refresh after creating page: %w", err)
		}
		
		// Try to get the page again
		page, exists = a.pages[pageName]
		if !exists {
			return nil, fmt.Errorf("page not found after creation: %s", pageName)
		}
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

// SegmentData represents a text segment for frontend
type SegmentData struct {
	Type    string `json:"type"`    // "text", "bold", "italic", "link", "image"
	Content string `json:"content"`
	Target  string `json:"target,omitempty"` // For links and images
	Alt     string `json:"alt,omitempty"`    // For images
}

// BlockData represents block data for frontend
type BlockData struct {
	ID string `json:"id"`
	Content string `json:"content"`
	HTMLContent string `json:"htmlContent"` // Deprecated - will be removed
	Segments []SegmentData `json:"segments"`
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
			HTMLContent: block.RenderHTML(), // Deprecated - kept for compatibility
			Segments: convertSegments(block.Segments),
			Depth: block.Depth,
			Children: convertBlocks(block.Children),
			TodoState: string(block.TodoInfo.TodoState),
			CheckboxState: string(block.TodoInfo.CheckboxState),
			Priority: block.TodoInfo.Priority,
		}
	}
	return result
}

// convertSegments converts parser segments to frontend segments
func convertSegments(segments []parser.Segment) []SegmentData {
	result := make([]SegmentData, len(segments))
	for i, seg := range segments {
		segType := "text"
		switch seg.Type {
		case parser.SegmentBold:
			segType = "bold"
		case parser.SegmentItalic:
			segType = "italic"
		case parser.SegmentLink:
			segType = "link"
		case parser.SegmentImage:
			segType = "image"
		}
		
		result[i] = SegmentData{
			Type:    segType,
			Content: seg.Content,
			Target:  seg.Target,
			Alt:     seg.Alt,
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

// UpdateBlockAtPath updates a block's content using positional addressing
func (a *App) UpdateBlockAtPath(pageName string, path BlockPath, newContent string) (map[string]interface{}, error) {
	// Work with current state for incremental updates
	page, exists := a.pages[pageName]
	if !exists {
		return nil, fmt.Errorf("page '%s' not found", pageName)
	}
	
	// Find the block by path
	block, err := FindBlockByPath(page.Blocks, path)
	if err != nil {
		return nil, fmt.Errorf("failed to find block: %w", err)
	}
	
	// Update the block content
	oldContent := block.Content
	block.SetContent(newContent)
	
	// Save the page
	if err := a.savePage(page); err != nil {
		return nil, fmt.Errorf("failed to save page: %w", err)
	}
	
	
	// Return delta for incremental update
	blockData := BlockData{
		Content:       block.Content,
		HTMLContent:   block.RenderHTML(),
		Segments:      convertSegments(block.Segments),
		Depth:         block.Depth,
		TodoState:     string(block.TodoInfo.TodoState),
		CheckboxState: string(block.TodoInfo.CheckboxState),
		Priority:      block.TodoInfo.Priority,
		Children:      []BlockData{}, // Children don't change
	}
	
	delta := map[string]interface{}{
		"action": "update",
		"path":   path,
		"block":  blockData,
		"oldContent": oldContent,
	}
	
	// Update backlinks incrementally if references changed
	oldRefs := extractPageReferences(oldContent)
	newRefs := extractPageReferences(newContent)
	
	// Update backlinks for changed references
	for _, ref := range oldRefs {
		if !contains(newRefs, ref) {
			// Reference removed
			a.removeBacklink(ref, pageName)
		}
	}
	for _, ref := range newRefs {
		if !contains(oldRefs, ref) {
			// Reference added
			a.addBacklink(ref, pageName)
		}
	}
	
	return delta, nil
}

// Helper functions for incremental backlink updates
func (a *App) addBacklink(targetPage, sourcePage string) {
	if a.backlinks == nil {
		a.backlinks = &parser.BacklinkIndex{
			ForwardLinks:  make(map[string]map[string][]parser.BlockReference),
			BackwardLinks: make(map[string]map[string][]parser.BlockReference),
		}
	}
	
	// Add to backward links
	if _, exists := a.backlinks.BackwardLinks[targetPage]; !exists {
		a.backlinks.BackwardLinks[targetPage] = make(map[string][]parser.BlockReference)
	}
	if _, exists := a.backlinks.BackwardLinks[targetPage][sourcePage]; !exists {
		a.backlinks.BackwardLinks[targetPage][sourcePage] = []parser.BlockReference{}
	}
	
	// Add to forward links
	if _, exists := a.backlinks.ForwardLinks[sourcePage]; !exists {
		a.backlinks.ForwardLinks[sourcePage] = make(map[string][]parser.BlockReference)
	}
	if _, exists := a.backlinks.ForwardLinks[sourcePage][targetPage]; !exists {
		a.backlinks.ForwardLinks[sourcePage][targetPage] = []parser.BlockReference{}
	}
}

func (a *App) removeBacklink(targetPage, sourcePage string) {
	if a.backlinks == nil {
		return
	}
	
	// Remove from backward links
	if sources, exists := a.backlinks.BackwardLinks[targetPage]; exists {
		delete(sources, sourcePage)
		if len(sources) == 0 {
			delete(a.backlinks.BackwardLinks, targetPage)
		}
	}
	
	// Remove from forward links
	if targets, exists := a.backlinks.ForwardLinks[sourcePage]; exists {
		delete(targets, targetPage)
		if len(targets) == 0 {
			delete(a.backlinks.ForwardLinks, sourcePage)
		}
	}
}

// extractPageReferences finds all [[page]] references in text
func extractPageReferences(text string) []string {
	pattern := regexp.MustCompile(`\[\[(.*?)\]\]`)
	matches := pattern.FindAllStringSubmatch(text, -1)
	references := make([]string, 0, len(matches))
	
	for _, match := range matches {
		if len(match) > 1 {
			references = append(references, match[1])
		}
	}
	
	return references
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// savePage writes a page back to disk
func (a *App) savePage(page *parser.Page) error {
	// Reconstruct the markdown content
	content := a.pageToMarkdown(page)
	
	// Determine the filename
	var filename string
	if parser.IsDatePage(page.Title) {
		// For date pages, parse the title and use ISO format
		date, err := parser.ParseDateTitle(page.Title)
		if err != nil {
			// Fall back to regular filename conversion
			filename = parser.TitleToFilename(page.Title)
		} else {
			filename = parser.GetDatePageFilename(date)
		}
	} else {
		// Regular pages use title-based filenames
		filename = parser.TitleToFilename(page.Title)
	}
	
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
	
	// Join lines and ensure trailing newline
	content := strings.Join(lines, "\n")
	if len(content) > 0 && content[len(content)-1] != '\n' {
		content += "\n"
	}
	return content
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

// createPage creates a new regular page with default content
func (a *App) createPage(pageTitle string) error {
	// Generate the filename
	filename := parser.TitleToFilename(pageTitle)
	filePath := filepath.Join(a.currentDir, filename)
	
	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File already exists, no need to create
		return nil
	}
	
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	// Write default content
	content := fmt.Sprintf(`# %s

- 
`, pageTitle)
	
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	
	return nil
}

// createDatePage creates a new date page with default content
func (a *App) createDatePage(dateTitle string) error {
	// Parse the date from the title
	date, err := parser.ParseDateTitle(dateTitle)
	if err != nil {
		return fmt.Errorf("invalid date title: %w", err)
	}
	
	// Generate the filename
	filename := parser.GetDatePageFilename(date)
	filePath := filepath.Join(a.currentDir, filename)
	
	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File already exists, no need to create
		return nil
	}
	
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	// Write default content for a journal page
	content := fmt.Sprintf(`# %s

- 

`, dateTitle)
	
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	
	return nil
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

// generateBlockID generates a unique block ID
func generateBlockID() string {
	// Generate 16 random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("block-%d", time.Now().UnixNano())
	}
	// Convert to hex string
	return hex.EncodeToString(bytes)
}

// AddBlock adds a new block to a page
func (a *App) AddBlock(pageName string, parentBlockID string, afterBlockID string, content string, depth int) (map[string]interface{}, error) {
	// Refresh to ensure we have latest content
	if err := a.RefreshPages(); err != nil {
		return nil, fmt.Errorf("failed to refresh pages: %w", err)
	}
	
	page, exists := a.pages[pageName]
	if !exists {
		return nil, fmt.Errorf("page '%s' not found", pageName)
	}
	
	// Create new block
	newBlock := &parser.Block{
		ID:       generateBlockID(),
		Content:  content,
		Depth:    depth,
		Children: []*parser.Block{},
	}
	
	// Parse the content into lines
	lines := strings.Split(content, "\n")
	newBlock.Lines = make([]parser.Line, len(lines))
	for i, line := range lines {
		newBlock.Lines[i] = parser.ParseLine(i+1, line)
	}
	
	// Update block content to parse TODO info and segments
	newBlock.SetContent(content)
	
	// Find insertion point
	if parentBlockID != "" {
		// Add as child of specific parent
		parent := findBlockByID(page.Blocks, parentBlockID)
		if parent == nil {
			return nil, fmt.Errorf("parent block '%s' not found", parentBlockID)
		}
		
		// Set parent relationship
		newBlock.Parent = parent
		newBlock.Depth = parent.Depth + 1
		
		if afterBlockID != "" {
			// Insert after specific sibling
			inserted := false
			for i, child := range parent.Children {
				if child.ID == afterBlockID {
					// Insert after this child
					parent.Children = append(parent.Children[:i+1], append([]*parser.Block{newBlock}, parent.Children[i+1:]...)...)
					inserted = true
					break
				}
			}
			if !inserted {
				return nil, fmt.Errorf("after block '%s' not found in parent's children", afterBlockID)
			}
		} else {
			// Add to end of parent's children
			parent.Children = append(parent.Children, newBlock)
		}
	} else {
		// Add as top-level block
		if afterBlockID != "" {
			// Insert after specific block
			inserted := false
			for i, block := range page.Blocks {
				if block.ID == afterBlockID {
					// Insert after this block
					page.Blocks = append(page.Blocks[:i+1], append([]*parser.Block{newBlock}, page.Blocks[i+1:]...)...)
					inserted = true
					break
				}
			}
			if !inserted {
				return nil, fmt.Errorf("after block '%s' not found", afterBlockID)
			}
		} else {
			// Add to end of page
			page.Blocks = append(page.Blocks, newBlock)
		}
	}
	
	// Save the page
	if err := a.savePage(page); err != nil {
		return nil, fmt.Errorf("failed to save page: %w", err)
	}
	
	// Refresh to update backlinks
	if err := a.RefreshPages(); err != nil {
		return nil, fmt.Errorf("failed to refresh after save: %w", err)
	}
	
	// Return the new block data
	blockData := BlockData{
		ID:            newBlock.ID,
		Content:       newBlock.Content,
		HTMLContent:   newBlock.RenderHTML(),
		Segments:      convertSegments(newBlock.Segments),
		Depth:         newBlock.Depth,
		Children:      []BlockData{}, // New block has no children yet
		TodoState:     string(newBlock.TodoInfo.TodoState),
		CheckboxState: string(newBlock.TodoInfo.CheckboxState),
		Priority:      newBlock.TodoInfo.Priority,
	}
	
	return map[string]interface{}{
		"id":    newBlock.ID,
		"block": blockData,
	}, nil
}

// AddBlockAtPath adds a new block using positional addressing
func (a *App) AddBlockAtPath(pageName string, insertPath BlockPath, content string) (map[string]interface{}, error) {
	// Work with current state for incremental updates
	page, exists := a.pages[pageName]
	if !exists {
		return nil, fmt.Errorf("page '%s' not found", pageName)
	}
	
	// Create new block
	newBlock := &parser.Block{
		Content:  content,
		Children: []*parser.Block{},
	}
	
	// Parse the content into lines
	lines := strings.Split(content, "\n")
	newBlock.Lines = make([]parser.Line, len(lines))
	for i, line := range lines {
		newBlock.Lines[i] = parser.ParseLine(i+1, line)
	}
	
	// Update block content to parse TODO info and segments
	newBlock.SetContent(content)
	
	// Handle insertion based on path
	if len(insertPath) == 0 {
		return nil, fmt.Errorf("empty insertion path")
	}
	
	if len(insertPath) == 1 {
		// Top-level insertion
		index := insertPath[0]
		if index < 0 || index > len(page.Blocks) {
			return nil, fmt.Errorf("invalid insertion index %d (max: %d)", index, len(page.Blocks))
		}
		
		// Set depth for top-level
		newBlock.Depth = 0
		
		// Insert at position
		if index == len(page.Blocks) {
			page.Blocks = append(page.Blocks, newBlock)
		} else {
			page.Blocks = append(page.Blocks[:index], 
				append([]*parser.Block{newBlock}, page.Blocks[index:]...)...)
		}
	} else {
		// Nested insertion - find parent
		parentPath := insertPath[:len(insertPath)-1]
		parent, err := FindBlockByPath(page.Blocks, parentPath)
		if err != nil {
			return nil, fmt.Errorf("parent not found: %w", err)
		}
		
		// Set parent relationship
		newBlock.Parent = parent
		newBlock.Depth = parent.Depth + 1
		
		// Insert into parent's children
		index := insertPath[len(insertPath)-1]
		if index < 0 || index > len(parent.Children) {
			return nil, fmt.Errorf("invalid insertion index %d (max: %d)", index, len(parent.Children))
		}
		
		if index == len(parent.Children) {
			parent.Children = append(parent.Children, newBlock)
		} else {
			parent.Children = append(parent.Children[:index], 
				append([]*parser.Block{newBlock}, parent.Children[index:]...)...)
		}
	}
	
	// Save the page
	if err := a.savePage(page); err != nil {
		return nil, fmt.Errorf("failed to save page: %w", err)
	}
	
	// Calculate path shifts for other blocks
	shifts := CalculatePathShiftsAfterInsert(page.Blocks, insertPath)
	
	// Return delta for incremental update
	blockData := BlockData{
		Content:       newBlock.Content,
		HTMLContent:   newBlock.RenderHTML(),
		Segments:      convertSegments(newBlock.Segments),
		Depth:         newBlock.Depth,
		TodoState:     string(newBlock.TodoInfo.TodoState),
		CheckboxState: string(newBlock.TodoInfo.CheckboxState),
		Priority:      newBlock.TodoInfo.Priority,
		Children:      []BlockData{}, // New block has no children
	}
	
	delta := map[string]interface{}{
		"action": "add",
		"path":   insertPath,
		"block":  blockData,
		"shifts": shifts,
	}
	
	// Update backlinks for any references in the new block
	refs := extractPageReferences(content)
	for _, ref := range refs {
		a.addBacklink(ref, pageName)
	}
	
	return delta, nil
}

// findBlockByID recursively searches for a block by ID
func findBlockByID(blocks []*parser.Block, targetID string) *parser.Block {
	for _, block := range blocks {
		if block.ID == targetID {
			return block
		}
		// Check children
		if found := findBlockByID(block.Children, targetID); found != nil {
			return found
		}
	}
	return nil
}
