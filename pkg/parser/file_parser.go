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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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