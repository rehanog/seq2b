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
)

// MultiPageResult represents the result of parsing multiple pages
type MultiPageResult struct {
	Pages     map[string]*Page  // Map of page name to page
	Backlinks *BacklinkIndex    // Cross-page backlink index
	Errors    []error          // Any parsing errors
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

