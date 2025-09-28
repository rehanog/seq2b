package templates

import (
	"fmt"
	"testing"
)

// BenchmarkTemplate provides patterns for performance testing
// Use this template when testing:
// - Parser performance
// - Search operations
// - Large file handling
// - Memory allocations
// - CPU-intensive operations

// Basic benchmark
func BenchmarkOperation(b *testing.B) {
	// Setup - not included in timing
	data := setupTestData()
	
	// Reset timer after setup
	b.ResetTimer()
	
	// Run the operation b.N times
	for i := 0; i < b.N; i++ {
		result := performOperation(data)
		// Prevent compiler optimization
		_ = result
	}
}

// Benchmark with different input sizes
func BenchmarkOperationSizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// Setup for this size
			data := generateData(size)
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				_ = performOperation(data)
			}
		})
	}
}

// Benchmark parallel operations
func BenchmarkParallel(b *testing.B) {
	// Setup shared resources
	resource := setupResource()
	
	b.ResetTimer()
	
	// Run in parallel
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine gets its own pb
		for pb.Next() {
			_ = resource.Process()
		}
	})
}

// Benchmark with memory allocation tracking
func BenchmarkWithAllocs(b *testing.B) {
	b.ReportAllocs() // Enable allocation reporting
	
	for i := 0; i < b.N; i++ {
		// Operation that allocates memory
		result := make([]byte, 1024)
		processData(result)
	}
}

// Comparative benchmarks
func BenchmarkComparison(b *testing.B) {
	b.Run("OldImplementation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = oldImplementation()
		}
	})
	
	b.Run("NewImplementation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = newImplementation()
		}
	})
}

// Benchmark with custom metrics
func BenchmarkCustomMetrics(b *testing.B) {
	var totalBytes int64
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		bytes := processLargeFile()
		totalBytes += bytes
	}
	
	// Report custom metric
	b.SetBytes(totalBytes / int64(b.N))
	b.ReportMetric(float64(totalBytes)/float64(b.N), "bytes/op")
}

// Memory-focused benchmark
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("WithPooling", func(b *testing.B) {
		pool := setupPool()
		b.ResetTimer()
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			buf := pool.Get()
			process(buf)
			pool.Put(buf)
		}
	})
	
	b.Run("WithoutPooling", func(b *testing.B) {
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			buf := make([]byte, 1024)
			process(buf)
		}
	})
}

// Benchmark table pattern
func BenchmarkTable(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
		setup func()
	}{
		{"Small", "small input", setupSmall},
		{"Medium", "medium sized input data", setupMedium},
		{"Large", "very large input data with lots of content", setupLarge},
	}
	
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			bm.setup()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				_ = processString(bm.input)
			}
		})
	}
}

// Placeholder functions - replace with actual implementation
func setupTestData() interface{} { return nil }
func performOperation(data interface{}) interface{} { return nil }
func generateData(size int) interface{} { return nil }
func setupResource() *Resource { return &Resource{} }
func processData([]byte) {}
func oldImplementation() interface{} { return nil }
func newImplementation() interface{} { return nil }
func processLargeFile() int64 { return 0 }
func setupPool() *Pool { return &Pool{} }
func process(interface{}) {}
func setupSmall() {}
func setupMedium() {}
func setupLarge() {}
func processString(string) interface{} { return nil }

type Resource struct{}
func (r *Resource) Process() interface{} { return nil }

type Pool struct{}
func (p *Pool) Get() interface{} { return nil }
func (p *Pool) Put(interface{}) {}