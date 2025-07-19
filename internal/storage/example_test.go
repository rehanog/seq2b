package storage_test

import (
	"fmt"
	"log"
	"time"
	
	"github.com/rehanog/seq2b/internal/storage"
)

// Example_productionUsage demonstrates production-ready cache usage
func Example_productionUsage() {
	// Initialize cache with metrics
	cache, err := storage.NewMetricsCacheManager("/path/to/library")
	if err != nil {
		// Log but don't fail - fall back to non-cached
		log.Printf("Warning: Cache initialization failed: %v", err)
		// Continue without cache
		return
	}
	defer func() {
		// Always close cache
		if err := cache.Close(); err != nil {
			log.Printf("Warning: Cache close error: %v", err)
		}
	}()
	
	// Monitor cache performance
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			stats := cache.GetMetrics()
			log.Printf("Cache stats: Hit rate=%.1f%%, Saves/s=%.0f, Gets/s=%.0f, Errors=%d",
				stats.HitRate, stats.SavesPerSec, stats.GetsPerSec, 
				stats.SaveErrors+stats.GetErrors)
			
			// Alert if performance degrades
			if stats.HitRate < 50 && stats.Gets > 100 {
				log.Printf("Warning: Low cache hit rate: %.1f%%", stats.HitRate)
			}
			
			if stats.AvgGetTime > 10*time.Millisecond {
				log.Printf("Warning: Slow cache gets: %v", stats.AvgGetTime)
			}
		}
	}()
	
	// Example page save with error handling
	page := struct {
		Title   string
		Content string
	}{
		Title:   "Example Page",
		Content: "Page content here",
	}
	
	err = cache.SavePage(page, "example", "/path/to/example.md", []string{"dep1", "dep2"})
	if err != nil {
		// Log but don't fail
		log.Printf("Warning: Failed to cache page: %v", err)
		// Continue - caching is optional
	}
	
	// Example page retrieval with fallback
	cachedData, hit, err := cache.GetPage("example", "/path/to/example.md")
	if err != nil {
		log.Printf("Warning: Cache get error: %v", err)
		// Fall back to parsing file
	} else if hit {
		fmt.Println("Cache hit! Using cached data")
		// Use cachedData
		_ = cachedData
	} else {
		fmt.Println("Cache miss - parsing file")
		// Parse file normally
	}
	
	// Graceful degradation example
	if stats := cache.GetMetrics(); stats.SaveErrors > 100 {
		log.Printf("Warning: High cache error rate, clearing cache")
		if err := cache.Clear(); err != nil {
			log.Printf("Failed to clear cache: %v", err)
		}
	}
}

// Example_disasterRecovery shows how to handle cache corruption
func Example_disasterRecovery() {
	libraryPath := "/path/to/library"
	
	// Attempt to use cache
	cache, err := storage.NewCacheManager(libraryPath)
	if err != nil {
		log.Printf("Cache error: %v", err)
		
		// Try to recover by clearing cache
		log.Println("Attempting cache recovery...")
		
		// Create temporary cache to clear it
		tempCache, err := storage.NewCacheManager("")
		if err == nil {
			tempCache.Clear()
			tempCache.Close()
		}
		
		// Try again
		cache, err = storage.NewCacheManager(libraryPath)
		if err != nil {
			log.Printf("Cache unrecoverable, proceeding without cache: %v", err)
			// Continue without cache
			return
		}
	}
	defer cache.Close()
	
	fmt.Println("Cache operational")
}

// Example_monitoring shows how to export metrics
func Example_monitoring() {
	cache, _ := storage.NewMetricsCacheManager("/path/to/library")
	defer cache.Close()
	
	// Simulate some operations
	for i := 0; i < 100; i++ {
		cache.SavePage(struct{}{}, fmt.Sprintf("page%d", i), "file.md", nil)
		cache.GetPage(fmt.Sprintf("page%d", i), "file.md")
	}
	
	// Export metrics
	stats := cache.GetMetrics()
	
	// Format for logging/monitoring system
	fmt.Printf("cache_hit_rate{library=\"mylib\"} %.2f\n", stats.HitRate)
	fmt.Printf("cache_saves_total{library=\"mylib\"} %d\n", stats.Saves)
	fmt.Printf("cache_gets_total{library=\"mylib\"} %d\n", stats.Gets)
	fmt.Printf("cache_errors_total{library=\"mylib\",op=\"save\"} %d\n", stats.SaveErrors)
	fmt.Printf("cache_errors_total{library=\"mylib\",op=\"get\"} %d\n", stats.GetErrors)
	fmt.Printf("cache_operation_duration_seconds{library=\"mylib\",op=\"save\"} %.6f\n", 
		stats.AvgSaveTime.Seconds())
	fmt.Printf("cache_operation_duration_seconds{library=\"mylib\",op=\"get\"} %.6f\n", 
		stats.AvgGetTime.Seconds())
}