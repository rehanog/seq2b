package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"github.com/rehanog/seq2b/pkg/parser"
	"github.com/rehanog/seq2b/internal/storage"
)

func main() {
	var (
		vaultPath = flag.String("vault", "", "Path to test vault")
		runs      = flag.Int("runs", 3, "Number of benchmark runs")
		clearCache = flag.Bool("clear-cache", false, "Clear cache before benchmarking")
	)
	flag.Parse()
	
	if *vaultPath == "" {
		fmt.Fprintf(os.Stderr, "Error: -vault flag is required\n")
		flag.Usage()
		os.Exit(1)
	}
	
	// Check if vault exists
	pagesDir := filepath.Join(*vaultPath, "pages")
	if _, err := os.Stat(pagesDir); err != nil {
		if _, err := os.Stat(*vaultPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: vault path does not exist: %s\n", *vaultPath)
			os.Exit(1)
		}
		pagesDir = *vaultPath
	}
	
	// Count files
	files, err := filepath.Glob(filepath.Join(pagesDir, "*.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding markdown files: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Benchmarking vault: %s\n", *vaultPath)
	fmt.Printf("Found %d markdown files\n\n", len(files))
	
	// Clear cache if requested
	if *clearCache {
		fmt.Println("Clearing cache...")
		cache, err := storage.NewCacheManager(pagesDir)
		if err == nil {
			cache.Clear()
			cache.Close()
		}
	}
	
	// Benchmark cold start (no cache)
	fmt.Println("=== Cold Start Benchmark (no cache) ===")
	coldTimes := make([]time.Duration, *runs)
	
	for i := 0; i < *runs; i++ {
		// Clear cache before each cold run
		cache, err := storage.NewCacheManager(pagesDir)
		if err == nil {
			cache.Clear()
			cache.Close()
		}
		
		start := time.Now()
		result, err := parser.ParseDirectoryWithCache(pagesDir)
		elapsed := time.Since(start)
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			continue
		}
		
		coldTimes[i] = elapsed
		fmt.Printf("Run %d: %v (%d pages, %d backlinks)\n", 
			i+1, elapsed, len(result.Pages), countBacklinks(result.Backlinks))
		
		// Small delay between runs
		time.Sleep(100 * time.Millisecond)
	}
	
	// Calculate cold start average
	var coldTotal time.Duration
	for _, t := range coldTimes {
		coldTotal += t
	}
	coldAvg := coldTotal / time.Duration(*runs)
	
	fmt.Printf("\nCold start average: %v\n", coldAvg)
	fmt.Printf("Pages per second: %.2f\n\n", float64(len(files))/coldAvg.Seconds())
	
	// Ensure cache is populated
	fmt.Println("Populating cache...")
	_, err = parser.ParseDirectoryWithCache(pagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error populating cache: %v\n", err)
		os.Exit(1)
	}
	
	// Small delay
	time.Sleep(500 * time.Millisecond)
	
	// Benchmark warm start (with cache)
	fmt.Println("=== Warm Start Benchmark (with cache) ===")
	warmTimes := make([]time.Duration, *runs)
	
	for i := 0; i < *runs; i++ {
		start := time.Now()
		result, err := parser.ParseDirectoryWithCache(pagesDir)
		elapsed := time.Since(start)
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			continue
		}
		
		warmTimes[i] = elapsed
		fmt.Printf("Run %d: %v (%d pages, %d backlinks)\n", 
			i+1, elapsed, len(result.Pages), countBacklinks(result.Backlinks))
		
		// Small delay between runs
		time.Sleep(100 * time.Millisecond)
	}
	
	// Calculate warm start average
	var warmTotal time.Duration
	for _, t := range warmTimes {
		warmTotal += t
	}
	warmAvg := warmTotal / time.Duration(*runs)
	
	fmt.Printf("\nWarm start average: %v\n", warmAvg)
	fmt.Printf("Pages per second: %.2f\n\n", float64(len(files))/warmAvg.Seconds())
	
	// Summary
	fmt.Println("=== Summary ===")
	fmt.Printf("Cold start: %v\n", coldAvg)
	fmt.Printf("Warm start: %v\n", warmAvg)
	speedup := float64(coldAvg) / float64(warmAvg)
	fmt.Printf("Speedup: %.2fx faster with cache\n", speedup)
	fmt.Printf("Time saved: %v\n", coldAvg-warmAvg)
	
	// Memory usage estimate
	cacheDir, _ := os.UserCacheDir()
	seq2bCache := filepath.Join(cacheDir, "seq2b")
	if info, err := dirSize(seq2bCache); err == nil {
		fmt.Printf("\nCache size: %.2f MB\n", float64(info)/1024/1024)
	}
}

func countBacklinks(index *parser.BacklinkIndex) int {
	count := 0
	for _, backlinks := range index.BackwardLinks {
		for _, refs := range backlinks {
			count += len(refs)
		}
	}
	return count
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}