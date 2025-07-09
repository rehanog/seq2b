package parser

import (
	"fmt"
	"os"
	"path/filepath"
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