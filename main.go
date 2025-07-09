package main

import (
	"fmt"
	"os"
)

func main() {
	// Step 1.1: Read a file and print its contents
	
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
	
	// Print file info
	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Size: %d bytes\n", len(content))
	fmt.Println("==============")
	fmt.Println(string(content))
}