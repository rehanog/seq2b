package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/dgraph-io/badger/v4"
	"github.com/rehanog/seq2b/internal/storage"
)

func main() {
	libraryPath := flag.String("library", "", "Path to library directory")
	flag.Parse()
	
	if *libraryPath == "" {
		fmt.Println("Usage: cache-inspect -library /path/to/library")
		os.Exit(1)
	}
	
	// Open cache in read-only mode
	cache, err := storage.NewCacheManager(*libraryPath)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()
	
	fmt.Printf("=== Cache Contents for %s ===\n\n", *libraryPath)
	
	// Access the underlying BadgerDB (you'd need to export db field for this)
	// For demo, we'll open it directly
	cacheDir := filepath.Join(*libraryPath, "cache")
	opts := badger.DefaultOptions(cacheDir)
	opts.Logger = nil
	opts.ReadOnly = true
	
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// Iterate through all keys
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		
		it := txn.NewIterator(opts)
		defer it.Close()
		
		count := 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			
			// Get value size
			valueSize := item.ValueSize()
			
			// Show key info
			fmt.Printf("Key: %s (value size: %d bytes)\n", key, valueSize)
			
			// For small values, show content
			if valueSize < 500 && strings.HasPrefix(key, "page:") {
				err := item.Value(func(val []byte) error {
					// Try to parse as CachedPage
					var cached struct {
						Page         json.RawMessage `json:"page"`
						FileModTime  string          `json:"file_mod_time"`
						Dependencies []string        `json:"dependencies"`
					}
					
					if err := json.Unmarshal(val, &cached); err == nil {
						fmt.Printf("  Modified: %s\n", cached.FileModTime)
						if len(cached.Dependencies) > 0 {
							fmt.Printf("  Dependencies: %v\n", cached.Dependencies)
						}
						
						// Show page content preview
						var pagePreview map[string]interface{}
						if err := json.Unmarshal(cached.Page, &pagePreview); err == nil {
							if title, ok := pagePreview["title"]; ok {
								fmt.Printf("  Title: %v\n", title)
							}
						}
					}
					return nil
				})
				
				if err != nil {
					fmt.Printf("  (Error reading value: %v)\n", err)
				}
			}
			
			fmt.Println()
			count++
			
			if count > 20 {
				fmt.Println("... (showing first 20 entries)")
				break
			}
		}
		
		// Count total entries
		total := 0
		it.Rewind()
		for it.Valid() {
			total++
			it.Next()
		}
		
		fmt.Printf("\nTotal entries: %d\n", total)
		
		return nil
	})
	
	if err != nil {
		log.Fatal(err)
	}
	
	// Show database stats
	lsm, vlog := db.Size()
	fmt.Printf("\nDatabase size:\n")
	fmt.Printf("  LSM size: %.2f MB\n", float64(lsm)/1024/1024)
	fmt.Printf("  Value log size: %.2f MB\n", float64(vlog)/1024/1024)
}