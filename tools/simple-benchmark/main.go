package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"github.com/rehanog/seq2b/pkg/parser"
	"github.com/rehanog/seq2b/internal/storage"
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
	
	// Count files
	files, err := filepath.Glob(filepath.Join(pagesDir, "*.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding files: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Testing with %d markdown files\n\n", len(files))
	
	// Clear cache first
	fmt.Println("Clearing cache...")
	cache, err := storage.NewCacheManager(pagesDir)
	if err == nil {
		cache.Clear()
		cache.Close()
	}
	
	// Cold start
	fmt.Println("=== Cold Start (no cache) ===")
	start := time.Now()
	result1, err := parser.ParseDirectoryWithCache(pagesDir)
	coldTime := time.Since(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Time: %v\n", coldTime)
	fmt.Printf("Pages: %d\n", len(result1.Pages))
	fmt.Printf("Speed: %.2f pages/sec\n\n", float64(len(files))/coldTime.Seconds())
	
	// Warm start (cache should be populated now)
	fmt.Println("=== Warm Start (with cache) ===")
	start = time.Now()
	result2, err := parser.ParseDirectoryWithCache(pagesDir)
	warmTime := time.Since(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Time: %v\n", warmTime)
	fmt.Printf("Pages: %d\n", len(result2.Pages))
	fmt.Printf("Speed: %.2f pages/sec\n\n", float64(len(files))/warmTime.Seconds())
	
	// Summary
	fmt.Println("=== Summary ===")
	speedup := float64(coldTime) / float64(warmTime)
	fmt.Printf("Cold: %v\n", coldTime)
	fmt.Printf("Warm: %v\n", warmTime)
	fmt.Printf("Speedup: %.2fx faster\n", speedup)
	fmt.Printf("Time saved: %v (%.1f%% reduction)\n", 
		coldTime-warmTime, 
		(1.0-float64(warmTime)/float64(coldTime))*100)
}