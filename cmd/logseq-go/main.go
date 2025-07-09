package main

import (
	"fmt"
	"os"
	
	"github.com/rehan/logseq-go/internal/parser"
)

func main() {
	// Get filename from command line args, or use default
	filename := "testdata/basic-markdown.md"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	
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