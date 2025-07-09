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
	"fmt"
	"os"
	"strings"
	
	"github.com/rehan/logseq-go/pkg/parser"
)

func main() {
	// Check if we have a directory or file argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: logseq-go <file.md> or logseq-go <directory>")
		return
	}
	
	path := os.Args[1]
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing path: %v\n", err)
		return
	}
	
	// Handle directory vs single file
	if fileInfo.IsDir() {
		handleDirectory(path)
	} else {
		handleSingleFile(path)
	}
}

func handleSingleFile(filename string) {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	// Parse the file with block structure
	result, err := parser.ParseFile(string(content))
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}
	
	// Print results
	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Title: %s\n", result.Page.Title)
	fmt.Printf("Lines: %d\n", len(result.Lines))
	fmt.Printf("Top-level blocks: %d\n", len(result.Page.Blocks))
	fmt.Printf("Total blocks: %d\n", len(result.Page.AllBlocks))
	fmt.Println("==============")
	
	// Print block tree
	fmt.Println("\nBlock Structure:")
	printBlockTree(result.Page.Blocks, "")
	
	// Print block relationships summary
	fmt.Println("\nBlock Relationships:")
	for _, block := range result.Page.AllBlocks {
		parentInfo := "top-level"
		if block.Parent != nil {
			parentInfo = fmt.Sprintf("child of %s", block.Parent.ID)
		}
		childCount := len(block.Children)
		childInfo := "no children"
		if childCount > 0 {
			childInfo = fmt.Sprintf("%d children", childCount)
		}
		fmt.Printf("  %s: %s, %s\n", block.ID, parentInfo, childInfo)
	}
	
	// For single file, just show what links it contains
	fmt.Println("\n=== Page References ===")
	
	// Extract all links from this page
	allLinks := make(map[string]int)
	for _, block := range result.Page.AllBlocks {
		links := parser.ExtractPageLinks(block.Content)
		for _, link := range links {
			allLinks[link]++
		}
	}
	
	if len(allLinks) > 0 {
		fmt.Println("This page references:")
		for link, count := range allLinks {
			fmt.Printf("  [[%s]] - %d time(s)\n", link, count)
		}
	} else {
		fmt.Println("No outgoing page references")
	}
	
	// Show orphan blocks
	orphans := parser.FindOrphanBlocks(result.Page.AllBlocks)
	if len(orphans) > 0 {
		fmt.Println("\nBlocks with no outgoing references:")
		for _, block := range orphans {
			fmt.Printf("  %s: %q\n", block.ID, block.Content)
		}
	}
}

func handleDirectory(dirPath string) {
	// Parse all files in directory
	result, err := parser.ParseDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error parsing directory: %v\n", err)
		return
	}
	
	// Report any errors
	if len(result.Errors) > 0 {
		fmt.Println("Parsing errors:")
		for _, err := range result.Errors {
			fmt.Printf("  - %v\n", err)
		}
		fmt.Println()
	}
	
	fmt.Printf("Directory: %s\n", dirPath)
	fmt.Printf("Pages found: %d\n", len(result.Pages))
	fmt.Println("==============")
	
	// Show page summaries
	fmt.Println("\nPages:")
	for pageName, page := range result.Pages {
		fmt.Printf("  %s - %d blocks\n", pageName, len(page.AllBlocks))
	}
	
	// Show backlinks for each page
	fmt.Println("\n=== Backlink Analysis ===")
	for _, pageName := range result.Backlinks.GetAllPages() {
		backlinks := result.Backlinks.GetBacklinks(pageName)
		forwardLinks := result.Backlinks.GetForwardLinks(pageName)
		
		fmt.Printf("\n%s:\n", pageName)
		
		// Show backlinks (who references this page)
		if len(backlinks) > 0 {
			fmt.Println("  ← Referenced by:")
			for sourcePage, refs := range backlinks {
				blockIDs := make([]string, len(refs))
				for i, ref := range refs {
					blockIDs[i] = ref.BlockID
				}
				fmt.Printf("    - %s (in blocks: %s)\n", sourcePage, strings.Join(blockIDs, ", "))
			}
		} else {
			fmt.Println("  ← No incoming references")
		}
		
		// Show forward links (what this page references)
		if len(forwardLinks) > 0 {
			fmt.Println("  → References:")
			for targetPage, refs := range forwardLinks {
				fmt.Printf("    - %s (%d times)\n", targetPage, len(refs))
			}
		} else {
			fmt.Println("  → No outgoing references")
		}
		
		// Mark orphans
		if result.Backlinks.IsOrphanPage(pageName) {
			fmt.Println("  [ORPHAN PAGE]")
		}
	}
	
	// Summary
	fmt.Println("\n=== Summary ===")
	orphanCount := 0
	for _, page := range result.Backlinks.GetAllPages() {
		if result.Backlinks.IsOrphanPage(page) {
			orphanCount++
		}
	}
	fmt.Printf("Total pages: %d\n", len(result.Pages))
	fmt.Printf("Orphan pages: %d\n", orphanCount)
}

// printBlockTree recursively prints the block hierarchy
func printBlockTree(blocks []*parser.Block, indent string) {
	for _, block := range blocks {
		parentID := "none"
		if block.Parent != nil {
			parentID = block.Parent.ID
		}
		fmt.Printf("%s├── [%s] (parent: %s, depth: %d): %s\n", 
			indent, block.ID, parentID, block.Depth, block.GetContent())
		fmt.Printf("%s    HTML: %s\n", indent, block.RenderHTML())
		if len(block.Children) > 0 {
			printBlockTree(block.Children, indent+"│   ")
		}
	}
}