package main

import (
	"fmt"
	"github.com/rehanog/seq2b/pkg/parser"
)

func main() {
	input := `# Test Page

- Block with property
  property:: block-level-value`
	
	result, err := parser.ParseFile(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Number of blocks: %d\n", len(result.Page.Blocks))
	if len(result.Page.Blocks) > 0 {
		block := result.Page.Blocks[0]
		fmt.Printf("Block ID: %s\n", block.ID)
		fmt.Printf("Block Properties: %v\n", block.Properties)
		fmt.Printf("Number of lines: %d\n", len(block.Lines))
		for i, line := range block.Lines {
			fmt.Printf("  Line %d: Type=%d, Content='%s', Properties=%v\n", i, line.Type, line.Content, line.Properties)
		}
	}
}
