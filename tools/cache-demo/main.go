package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	
	"github.com/rehanog/seq2b/internal/storage"
)

func main() {
	fmt.Println("=== BadgerDB Cache Demo ===\n")
	
	// Create a demo library directory
	demoLibrary := "./demo-library"
	if err := os.MkdirAll(filepath.Join(demoLibrary, "pages"), 0755); err != nil {
		log.Fatal(err)
	}
	
	// Show cache location
	cacheDir := filepath.Join(demoLibrary, "cache")
	fmt.Printf("Library: %s\n", demoLibrary)
	fmt.Printf("Cache location: %s\n", cacheDir)
	
	// List current cache files
	fmt.Println("\nCache files BEFORE:")
	listCacheFiles(cacheDir)
	
	// Create a test cache
	cache, err := storage.NewCacheManager(demoLibrary)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()
	
	// Create test data
	type DemoPage struct {
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
	}
	
	// Save some pages
	fmt.Println("\nSaving pages to cache...")
	for i := 1; i <= 5; i++ {
		page := DemoPage{
			Title:     fmt.Sprintf("Demo Page %d", i),
			Content:   fmt.Sprintf("This is content for page %d", i),
			Timestamp: time.Now(),
		}
		
		testFile := fmt.Sprintf("/tmp/demo%d.md", i)
		os.WriteFile(testFile, []byte(page.Content), 0644)
		
		err := cache.SavePage(page, page.Title, testFile, nil)
		if err != nil {
			log.Printf("Error saving page %d: %v", i, err)
		} else {
			fmt.Printf("  ✓ Saved: %s\n", page.Title)
		}
		
		// Small delay to show files appearing
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Println("\nCache files AFTER saving:")
	listCacheFiles(cacheDir)
	
	// Read back from cache
	fmt.Println("\nReading from cache:")
	for i := 1; i <= 5; i++ {
		pageName := fmt.Sprintf("Demo Page %d", i)
		testFile := fmt.Sprintf("/tmp/demo%d.md", i)
		
		cached, hit, err := cache.GetPage(pageName, testFile)
		if err != nil {
			log.Printf("Error reading %s: %v", pageName, err)
			continue
		}
		
		if !hit {
			fmt.Printf("  ✗ Cache miss: %s\n", pageName)
			continue
		}
		
		// Unmarshal and display
		if rawJSON, ok := cached.(json.RawMessage); ok {
			var page DemoPage
			if err := json.Unmarshal(rawJSON, &page); err == nil {
				fmt.Printf("  ✓ Cache hit: %s (saved at %s)\n", 
					page.Title, page.Timestamp.Format("15:04:05"))
			}
		}
	}
	
	// Show cache metadata
	fmt.Println("\nCache validation:")
	valid, err := cache.ValidateCache()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else if valid {
		fmt.Println("  ✓ Cache is valid")
	} else {
		fmt.Println("  ✗ Cache is invalid")
	}
	
	// Demonstrate file modification detection
	fmt.Println("\nModifying a file...")
	time.Sleep(10 * time.Millisecond) // Ensure timestamp changes
	os.WriteFile("/tmp/demo1.md", []byte("Modified content"), 0644)
	
	_, hit, _ := cache.GetPage("Demo Page 1", "/tmp/demo1.md")
	if !hit {
		fmt.Println("  ✓ Cache correctly detected file modification")
	} else {
		fmt.Println("  ✗ Cache failed to detect modification")
	}
	
	fmt.Println("\nDemo complete!")
}

func listCacheFiles(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("  (cache directory not found)\n")
		return
	}
	
	for _, file := range files {
		info, _ := file.Info()
		fmt.Printf("  %s (%d bytes)\n", file.Name(), info.Size())
	}
	
	// Show total size
	var totalSize int64
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	fmt.Printf("  Total: %.2f MB\n", float64(totalSize)/1024/1024)
}