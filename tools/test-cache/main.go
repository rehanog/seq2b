package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/rehanog/seq2b/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <vault-path>\n", os.Args[0])
		os.Exit(1)
	}
	
	vaultPath := os.Args[1]
	pagesDir := filepath.Join(vaultPath, "pages")
	if _, err := os.Stat(pagesDir); err != nil {
		pagesDir = vaultPath
	}
	
	fmt.Println("First parse (should miss cache):")
	result1, err := parser.ParseDirectoryWithCache(pagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Parsed %d pages\n\n", len(result1.Pages))
	
	fmt.Println("Second parse (should hit cache):")
	result2, err := parser.ParseDirectoryWithCache(pagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Parsed %d pages\n", len(result2.Pages))
}