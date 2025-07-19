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

package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/rehanog/seq2b/internal/storage"
)

// ParseResult represents the result of parsing a file
type ParseResult struct {
	Page     *Page
	Lines    []Line    // All parsed lines
	Errors   []error   // Any parsing errors encountered
}

// MultiPageResult represents the result of parsing multiple pages
type MultiPageResult struct {
	Pages     map[string]*Page  // Map of page name to page
	Backlinks *BacklinkIndex    // Cross-page backlink index
	Errors    []error          // Any parsing errors
}

// parseContext is used temporarily during parsing
type parseContext struct {
	line        Line
	indentLevel int    // Calculated from raw line, used once, discarded
}

// calculateIndentLevel counts leading spaces and divides by 2 (Logseq standard)
func calculateIndentLevel(rawLine string) int {
	spaces := 0
	for _, ch := range rawLine {
		if ch == ' ' {
			spaces++
		} else {
			break
		}
	}
	return spaces / 2
}

// ParseFile parses markdown content into a Page with block structure
func ParseFile(content string) (*ParseResult, error) {
	lines := []Line{}
	contexts := []parseContext{}
	
	// Step 1: Parse all lines and extract indent levels
	rawLines := strings.Split(content, "\n")
	for i, rawLine := range rawLines {
		indentLevel := calculateIndentLevel(rawLine)
		line := ParseLine(i+1, strings.TrimSpace(rawLine))
		
		lines = append(lines, line)
		contexts = append(contexts, parseContext{line, indentLevel})
	}
	
	// Step 2: Build block tree from parsed lines
	blocks := BuildBlockTree(contexts)
	
	// Step 3: Create page with all blocks
	page := &Page{
		Blocks:    blocks,
		Created:   time.Now(),
		Modified:  time.Now(),
	}
	
	// Extract title from first header if present
	for _, line := range lines {
		if line.Type == TypeHeader {
			page.Title = line.Content
			break
		}
	}
	
	// Build flat list of all blocks
	page.AllBlocks = page.GetAllBlocks()
	
	return &ParseResult{
		Page:  page,
		Lines: lines,
	}, nil
}

// ParseDirectory parses all markdown files in a directory
func ParseDirectory(dirPath string) (*MultiPageResult, error) {
	result := &MultiPageResult{
		Pages:     make(map[string]*Page),
		Backlinks: NewBacklinkIndex(),
		Errors:    []error{},
	}
	
	// Find all markdown files
	files, err := filepath.Glob(filepath.Join(dirPath, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("error finding files: %w", err)
	}
	
	// Parse each file
	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			result.Errors = append(result.Errors, 
				fmt.Errorf("error reading %s: %w", filePath, err))
			continue
		}
		
		// Parse the file
		parseResult, err := ParseFile(string(content))
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Errorf("error parsing %s: %w", filePath, err))
			continue
		}
		
		// Store the page
		page := parseResult.Page
		result.Pages[page.Title] = page
		
		// Add to backlink index
		result.Backlinks.AddPage(page)
	}
	
	return result, nil
}

// ParseDirectoryWithCache parses a directory using cache for unchanged files
func ParseDirectoryWithCache(dirPath string) (*MultiPageResult, error) {
	result := &MultiPageResult{
		Pages:     make(map[string]*Page),
		Backlinks: NewBacklinkIndex(),
		Errors:    []error{},
	}
	
	// Initialize cache
	cache, err := storage.NewCacheManager(dirPath)
	if err != nil {
		// Fall back to regular parsing if cache fails
		return ParseDirectory(dirPath)
	}
	defer cache.Close()
	
	// Validate cache
	valid, err := cache.ValidateCache()
	if err != nil {
		fmt.Printf("Cache validation error: %v\n", err)
		valid = false
	}
	if !valid {
		fmt.Println("Cache invalid, clearing and rebuilding...")
		cache.Clear()
		cache.SaveMetadata()
	} else {
		fmt.Println("Cache is valid, using cached data...")
	}
	
	// Find all markdown files
	files, err := filepath.Glob(filepath.Join(dirPath, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("error finding files: %w", err)
	}
	
	cacheHits := 0
	cacheMisses := 0
	startTime := time.Now()
	
	// Parse each file
	for _, filePath := range files {
		// Extract page name from filename
		baseName := filepath.Base(filePath)
		pageName := baseName[:len(baseName)-3] // Remove .md extension
		
		// Try cache first
		if valid {
			cachedPage, hit, err := cache.GetPage(pageName, filePath)
			if err != nil && err.Error() != "Key not found" {
				fmt.Printf("Cache get error for %s: %v\n", pageName, err)
			}
			if hit {
				// Unmarshal the raw JSON into a Page
				if rawJSON, ok := cachedPage.(json.RawMessage); ok {
					var page Page
					if err := json.Unmarshal(rawJSON, &page); err == nil {
						cacheHits++
						result.Pages[page.Title] = &page
						result.Backlinks.AddPage(&page)
						continue
					} else {
						fmt.Printf("Cache unmarshal error for %s: %v\n", pageName, err)
					}
				} else {
					fmt.Printf("Cache type assertion failed for %s (got %T)\n", pageName, cachedPage)
				}
			}
		}
		
		// Cache miss - parse the file
		cacheMisses++
		content, err := os.ReadFile(filePath)
		if err != nil {
			result.Errors = append(result.Errors, 
				fmt.Errorf("error reading %s: %w", filePath, err))
			continue
		}
		
		// Parse the file
		parseResult, err := ParseFile(string(content))
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Errorf("error parsing %s: %w", filePath, err))
			continue
		}
		
		// Store the page
		page := parseResult.Page
		result.Pages[page.Title] = page
		
		// Extract dependencies
		var dependencies []string
		for _, block := range page.Blocks {
			extractBlockDependencies(block, &dependencies)
		}
		
		// Save to cache
		if err := cache.SavePage(page, pageName, filePath, dependencies); err != nil {
			// Log but don't fail
			result.Errors = append(result.Errors,
				fmt.Errorf("warning: failed to cache %s: %w", pageName, err))
		}
		
		// Add to backlink index
		result.Backlinks.AddPage(page)
	}
	
	// Save backlinks to cache
	for pageName := range result.Pages {
		backlinks := result.Backlinks.GetBacklinks(pageName)
		if backlinks != nil && len(backlinks) > 0 {
			if err := cache.SaveBacklinks(pageName, backlinks); err != nil {
				// Log but don't fail
				result.Errors = append(result.Errors,
					fmt.Errorf("warning: failed to cache backlinks for %s: %w", pageName, err))
			}
		}
	}
	
	elapsed := time.Since(startTime)
	fmt.Printf("Parsed %d files in %v (cache hits: %d, misses: %d)\n", 
		len(files), elapsed, cacheHits, cacheMisses)
	
	return result, nil
}

// extractBlockDependencies recursively extracts all links from a block
func extractBlockDependencies(block *Block, dependencies *[]string) {
	for _, segment := range block.Segments {
		if segment.Type == SegmentLink && segment.Target != "" {
			*dependencies = append(*dependencies, segment.Target)
		}
	}
	for _, child := range block.Children {
		extractBlockDependencies(child, dependencies)
	}
}

// ParseFiles parses specific markdown files
func ParseFiles(filePaths []string) (*MultiPageResult, error) {
	result := &MultiPageResult{
		Pages:     make(map[string]*Page),
		Backlinks: NewBacklinkIndex(),
		Errors:    []error{},
	}
	
	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Errorf("error reading %s: %w", filePath, err))
			continue
		}
		
		// Parse the file
		parseResult, err := ParseFile(string(content))
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Errorf("error parsing %s: %w", filePath, err))
			continue
		}
		
		// Store the page
		page := parseResult.Page
		result.Pages[page.Title] = page
		
		// Add to backlink index
		result.Backlinks.AddPage(page)
	}
	
	return result, nil
}

// TitleToFilename converts a page title to a filename
func TitleToFilename(title string) string {
	// Simple conversion: lowercase, replace spaces with hyphens, add .md
	filename := strings.ToLower(title)
	filename = strings.ReplaceAll(filename, " ", "-")
	
	// Remove any special characters that might be problematic
	filename = strings.ReplaceAll(filename, "/", "-")
	filename = strings.ReplaceAll(filename, "\\", "-")
	filename = strings.ReplaceAll(filename, ":", "-")
	
	return filename + ".md"
}

// BuildBlockTree converts flat lines into hierarchical block structure
func BuildBlockTree(contexts []parseContext) []*Block {
	var rootBlocks []*Block
	var blockStack []*Block // Stack to track current nesting
	blockID := 0
	
	for _, ctx := range contexts {
		// Skip empty lines
		if ctx.line.Type == TypeEmpty {
			continue
		}
		
		// Only process block items
		if ctx.line.Type == TypeBlock {
			blockID++
			newBlock := &Block{
				ID:    fmt.Sprintf("block-%d", blockID),
				Lines: []Line{ctx.line},
				Depth: ctx.indentLevel,
			}
			
			// Pop stack until we find the right parent level
			for len(blockStack) > 0 && blockStack[len(blockStack)-1].Depth >= ctx.indentLevel {
				blockStack = blockStack[:len(blockStack)-1]
			}
			
			// Set parent and add as child
			if len(blockStack) > 0 {
				parent := blockStack[len(blockStack)-1]
				newBlock.Parent = parent
				parent.Children = append(parent.Children, newBlock)
			} else {
				// Top-level block
				rootBlocks = append(rootBlocks, newBlock)
			}
			
			// Push onto stack
			blockStack = append(blockStack, newBlock)
			
			// Update content
			newBlock.updateContent()
		}
	}
	
	return rootBlocks
}