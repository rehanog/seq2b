# BadgerDB Implementation Walkthrough

## Overview

BadgerDB is an embeddable, persistent key-value database written in Go. In seq2b, we use it to cache parsed Logseq pages to dramatically improve startup times.

## How BadgerDB Works

### 1. **Architecture**
```
┌─────────────────┐
│   Application   │
├─────────────────┤
│  CacheManager   │ ← Our abstraction layer
├─────────────────┤
│    BadgerDB     │ ← Key-value store
├─────────────────┤
│  LSM Tree Files │ ← On-disk storage
└─────────────────┘
```

### 2. **Key Components**

- **SST Files** (`.sst`): Sorted String Table files containing the actual key-value data
- **Value Log** (`.vlog`): Stores the actual values (BadgerDB separates keys and values)
- **MANIFEST**: Metadata about the database state
- **DISCARD**: Tracks deleted entries for garbage collection

## Code Walkthrough

### 1. **Cache Initialization** (`cache.go:41-61`)

```go
func NewCacheManager(libraryPath string) (*CacheManager, error) {
    // Get platform-specific cache directory
    cacheDir, err := getCacheDir(libraryPath)  // <library>/cache/
    
    // Configure BadgerDB
    opts := badger.DefaultOptions(cacheDir)
    opts.Logger = nil  // Disable verbose logging
    
    // Open database (creates files if needed)
    db, err := badger.Open(opts)
}
```

### 2. **Saving Data** (`cache.go:69-103`)

```go
func SavePage(page interface{}, pageName string, filePath string, dependencies []string) error {
    // 1. Get file modification time for cache invalidation
    fileInfo, err := os.Stat(filePath)
    
    // 2. Marshal page data to JSON
    pageData, err := json.Marshal(page)
    
    // 3. Create cache entry with metadata
    cached := CachedPage{
        Page:         pageData,
        FileModTime:  fileInfo.ModTime(),
        Dependencies: dependencies,
    }
    
    // 4. Save to BadgerDB in a transaction
    err = cm.db.Update(func(txn *badger.Txn) error {
        key := []byte("page:" + pageName)  // e.g., "page:Daily Note"
        return txn.Set(key, data)
    })
}
```

### 3. **Retrieving Data** (`cache.go:106-141`)

```go
func GetPage(pageName string, filePath string) (interface{}, bool, error) {
    // 1. Check if file was modified since caching
    fileInfo, err := os.Stat(filePath)
    
    // 2. Read from BadgerDB
    err = cm.db.View(func(txn *badger.Txn) error {
        key := []byte("page:" + pageName)
        item, err := txn.Get(key)
        
        // 3. Decode value
        return item.Value(func(val []byte) error {
            return json.Unmarshal(val, &cached)
        })
    })
    
    // 4. Validate cache freshness
    if fileInfo.ModTime().After(cached.FileModTime) {
        return nil, false, nil  // Cache miss - file was modified
    }
    
    return cached.Page, true, nil  // Cache hit
}
```

## On-Disk Storage

### Location
- **Production**: `<library-path>/cache/`
- **Testing**: `<repo-root>/tmp/cache/`

### File Structure
```bash
/path/to/library/cache/
├── 000007.sst      # Small SST file with keys
├── 000127.sst      # Larger SST file (9.5KB)
├── 000130.sst      # Main data file (69KB)
├── 000024.vlog     # Value log files
├── 000025.vlog     
├── MANIFEST        # Database metadata
├── DISCARD         # Garbage collection info
└── KEYREGISTRY     # Key conflict detection
```

### What's in These Files?

1. **SST Files**: Contain the actual key-value pairs in sorted order
   - Keys like `page:2025-07-19`, `page:Daily Note 45`
   - Values are JSON-encoded page data

2. **VLOG Files**: Store the actual values
   - BadgerDB uses value log to optimize for SSDs
   - Reduces write amplification

3. **MANIFEST**: Tracks database state
   - Which SST/VLOG files are active
   - Database version information

## Monitoring Cache Activity

### 1. **Watch Files Change in Real-Time**
```bash
# Watch cache directory for changes
watch -n 1 'ls -la /path/to/library/cache/'

# See file sizes change as cache grows
du -sh /path/to/library/cache/*
```

### 2. **View Cache Contents** (Debug Tool)
```go
// Add this debug function to cache.go
func (cm *CacheManager) DebugDump() {
    cm.db.View(func(txn *badger.Txn) error {
        opts := badger.DefaultIteratorOptions
        it := txn.NewIterator(opts)
        defer it.Close()
        
        for it.Rewind(); it.Valid(); it.Next() {
            item := it.Item()
            key := string(item.Key())
            fmt.Printf("Key: %s\n", key)
        }
        return nil
    })
}
```

### 3. **Cache Metrics**
```bash
# Run the app and watch cache metrics
./seq2b --debug-cache

# Output:
# Cache stats: Hit rate=85.3%, Saves/s=127, Gets/s=1847
```

## Performance Characteristics

### Write Path
1. Application calls `SavePage()`
2. Data is serialized to JSON
3. BadgerDB writes to WAL (Write-Ahead Log)
4. Data is written to MemTable in memory
5. Periodically flushed to SST files on disk

### Read Path
1. Application calls `GetPage()`
2. BadgerDB checks MemTable first (fast)
3. If not found, checks SST files (still fast due to indexing)
4. Returns data or cache miss

### Why It's Fast
- **LSM Tree**: Optimized for write-heavy workloads
- **Separation of Keys/Values**: Reduces I/O for key lookups
- **Memory Mapping**: OS page cache speeds up reads
- **Compression**: Reduces disk usage (Snappy compression)

## Debugging Cache Issues

### 1. **Check Cache is Working**
```bash
# First run - builds cache
time ./seq2b /path/to/vault
# Parsed 1000 files in 1.5s (cache hits: 0, misses: 1000)

# Second run - uses cache
time ./seq2b /path/to/vault
# Parsed 1000 files in 0.4s (cache hits: 1000, misses: 0)
```

### 2. **Clear Cache**
```bash
rm -rf /path/to/library/cache/
```

### 3. **Verify Cache Integrity**
```go
// BadgerDB has built-in verification
badger info --dir /path/to/library/cache/
```

## Advanced Topics

### Garbage Collection
BadgerDB periodically runs GC to reclaim space:
```go
// Run GC manually
db.RunValueLogGC(0.5)  // Reclaim if >50% space is unused
```

### Transactions
All operations are ACID compliant:
- **Atomicity**: Updates succeed or fail completely
- **Consistency**: Data integrity maintained
- **Isolation**: Concurrent reads/writes handled correctly
- **Durability**: Survives crashes (WAL ensures this)

### Configuration Tuning
```go
opts := badger.DefaultOptions(dir)
opts.ValueLogFileSize = 1 << 20      // 1MB files for testing
opts.MemTableSize = 1 << 20          // 1MB memory tables
opts.NumMemtables = 2                // Number of memtables
opts.NumLevelZeroTables = 2          // L0 tables before compaction
```

## Common Issues and Solutions

1. **"Database locked" error**: Another process has the DB open
   - Solution: Ensure only one instance runs

2. **Slow performance**: Cache might be corrupted
   - Solution: Clear cache and rebuild

3. **Disk space**: Cache grows over time
   - Solution: Implement size limits or periodic cleanup

4. **Permission errors**: Cache directory not writable
   - Solution: Check directory permissions