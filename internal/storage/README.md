# Storage Package - Cache System Documentation

## Overview

The storage package provides a high-performance, persistent caching layer for parsed Logseq pages using BadgerDB. It significantly improves startup times by avoiding re-parsing of unchanged files.

## Features

- **Persistent Storage**: Uses BadgerDB for on-disk key-value storage
- **Automatic Invalidation**: Detects file modifications via timestamps
- **Concurrent Safe**: Handles multiple readers/writers gracefully
- **Metrics Collection**: Built-in performance monitoring
- **Failure Resilient**: Gracefully handles corruption, disk issues, permissions

## Performance

Based on benchmarks with production-like workloads:

- **Cache Save**: ~7.5μs per operation
- **Cache Get**: ~3μs per operation
- **Speedup**: 2.5-3x faster startup with warm cache
- **Throughput**: 150,000+ operations/second

## Usage

### Basic Usage

```go
// Create cache manager
cache, err := storage.NewCacheManager("/path/to/library")
if err != nil {
    // Falls back to non-cached parsing
}
defer cache.Close()

// Save a page
err = cache.SavePage(page, "PageName", "/path/to/file.md", dependencies)

// Retrieve a page
cachedData, hit, err := cache.GetPage("PageName", "/path/to/file.md")
if hit {
    // Use cached data
}
```

### With Metrics

```go
// Create cache with metrics
cache, err := storage.NewMetricsCacheManager("/path/to/library")
defer cache.Close()

// Use normally...

// Get metrics
stats := cache.GetMetrics()
fmt.Printf("Hit rate: %.1f%%\n", stats.HitRate)
fmt.Printf("Saves/sec: %.0f\n", stats.SavesPerSec)
```

## Cache Location

The cache is stored alongside your library data in a visible directory:

- **Production**: `<library-path>/cache/`
- **Testing**: `<repo-root>/tmp/cache/`

This design ensures:
- Each library has its own isolated cache
- Cache moves with the library when relocated
- No conflicts between different libraries
- Easy cache management (just delete the cache folder)
- Transparency - users can see what's being cached

## Failure Modes

The cache is designed to fail gracefully:

1. **Disk Full**: Returns error, falls back to non-cached parsing
2. **Corrupted DB**: Detected on open, cache is cleared and rebuilt
3. **Permission Denied**: Falls back to non-cached parsing
4. **File Modified**: Automatically detected, entry is refreshed
5. **Network Drive**: Handles slow I/O gracefully

## Monitoring

Key metrics to monitor in production:

- **Hit Rate**: Should be >80% after warm-up
- **Save/Get Errors**: Should be near zero
- **Average Latency**: Get <5ms, Save <10ms
- **Cache Size**: Monitor disk usage

## Testing

The cache system has comprehensive test coverage:

- Unit tests for all operations
- Concurrent access tests
- Failure scenario tests (disk full, corruption, etc.)
- Performance benchmarks
- Production load simulation

Run tests:
```bash
go test ./internal/storage -v
go test -bench=. ./internal/storage
```

## Limitations

- Cache keys are based on filename, not content
- No automatic size limits (relies on OS disk management)
- BadgerDB requires exclusive access (one writer)
- Initial cache build still requires full parse

## Future Improvements

- [ ] Implement cache size limits with LRU eviction
- [ ] Add cache warming on startup
- [ ] Support for distributed caching
- [ ] Compression for large pages
- [ ] Export metrics to Prometheus/OpenTelemetry