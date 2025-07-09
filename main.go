package main

import (
	"fmt"
	"os"
	"strings"
)

// LineType represents the type of markdown line
type LineType int

const (
	TypeEmpty LineType = iota
	TypeHeader
	TypeText
	TypeList
)

// Line represents a parsed line from the markdown file
type Line struct {
	Number      int
	Type        LineType
	Content     string
	HeaderLevel int // Only used for headers
}

func main() {
	// Step 1.2: Parse file line by line, identify headers
	
	// Get filename from command line args, or use default
	filename := "test.md"
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
	parsedLines := []Line{}
	for i, line := range lines {
		parsed := parseLine(i+1, line)
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

// parseLine analyzes a single line and returns its type and content
func parseLine(number int, line string) Line {
	trimmed := strings.TrimSpace(line)
	
	// Empty line
	if trimmed == "" {
		return Line{Number: number, Type: TypeEmpty}
	}
	
	// Header (starts with #)
	if strings.HasPrefix(trimmed, "#") {
		level := 0
		for _, ch := range trimmed {
			if ch == '#' {
				level++
			} else {
				break
			}
		}
		// Extract header text (remove # and trim)
		headerText := strings.TrimSpace(trimmed[level:])
		return Line{
			Number:      number,
			Type:        TypeHeader,
			Content:     headerText,
			HeaderLevel: level,
		}
	}
	
	// List item (starts with -)
	if strings.HasPrefix(trimmed, "-") {
		listText := strings.TrimSpace(trimmed[1:])
		return Line{
			Number:  number,
			Type:    TypeList,
			Content: listText,
		}
	}
	
	// Regular text
	return Line{
		Number:  number,
		Type:    TypeText,
		Content: trimmed,
	}
}

// printLine formats and prints a parsed line
func printLine(line Line) {
	switch line.Type {
	case TypeEmpty:
		fmt.Printf("Line %d: Empty\n", line.Number)
	case TypeHeader:
		fmt.Printf("Line %d: Header (level %d): %s\n", line.Number, line.HeaderLevel, line.Content)
	case TypeList:
		fmt.Printf("Line %d: List item: %s\n", line.Number, line.Content)
	case TypeText:
		fmt.Printf("Line %d: Text: %s\n", line.Number, line.Content)
	}
}