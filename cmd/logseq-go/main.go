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

// printLine formats and prints a parsed line
func printLine(line parser.Line) {
	switch line.Type {
	case parser.TypeEmpty:
		fmt.Printf("Line %d: Empty\n", line.Number)
	case parser.TypeHeader:
		fmt.Printf("Line %d: Header (level %d): %s\n", line.Number, line.HeaderLevel, line.Content)
		printElements(line.Elements)
	case parser.TypeList:
		fmt.Printf("Line %d: List item: %s\n", line.Number, line.Content)
		printElements(line.Elements)
	case parser.TypeText:
		fmt.Printf("Line %d: Text: %s\n", line.Number, line.Content)
		printElements(line.Elements)
	}
}

// printElements shows the parsed markdown elements
func printElements(elements []parser.MarkdownElement) {
	if len(elements) == 0 {
		return
	}
	
	for i, elem := range elements {
		fmt.Printf("  Element %d: ", i+1)
		
		features := []string{}
		if elem.Bold {
			features = append(features, "bold")
		}
		if elem.Italic {
			features = append(features, "italic")
		}
		if elem.Link != "" {
			features = append(features, fmt.Sprintf("link->%s", elem.Link))
		}
		
		if len(features) > 0 {
			fmt.Printf("[%s] ", strings.Join(features, ", "))
		}
		fmt.Printf("'%s'\n", elem.Text)
	}
}