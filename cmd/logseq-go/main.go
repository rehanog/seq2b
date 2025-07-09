package main

import (
	"fmt"
	"os"
	"strings"
	
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
	
	// Split into lines
	lines := strings.Split(string(content), "\n")
	
	// Parse each line
	parsedLines := []parser.Line{}
	for i, line := range lines {
		parsed := parser.ParseLine(i+1, line)
		parsedLines = append(parsedLines, parsed)
	}
	
	// Print results
	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Lines: %d\n", len(lines))
	fmt.Println("==============")
	
	for _, line := range parsedLines {
		printLine(line)
	}
}

// printLine formats and prints a parsed line with HTML rendering
func printLine(line parser.Line) {
	switch line.Type {
	case parser.TypeEmpty:
		fmt.Printf("Line %d: Empty\n", line.Number)
	case parser.TypeHeader:
		htmlContent := parser.RenderToHTML(line.Content)
		fmt.Printf("Line %d: Header (level %d): %s\n", line.Number, line.HeaderLevel, line.Content)
		fmt.Printf("  HTML: <h%d>%s</h%d>\n", line.HeaderLevel, htmlContent, line.HeaderLevel)
	case parser.TypeList:
		htmlContent := parser.RenderToHTML(line.Content)
		fmt.Printf("Line %d: List item: %s\n", line.Number, line.Content)
		fmt.Printf("  HTML: <li>%s</li>\n", htmlContent)
	case parser.TypeText:
		htmlContent := parser.RenderToHTML(line.Content)
		fmt.Printf("Line %d: Text: %s\n", line.Number, line.Content)
		fmt.Printf("  HTML: <p>%s</p>\n", htmlContent)
	}
}