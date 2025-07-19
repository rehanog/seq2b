# ADR-004: Use BadgerDB for Persistent Cache

Date: 2025-01-19

## Status

Accepted

## Context

As seq2b users work with larger Logseq vaults (1000+ pages), the startup time becomes a significant usability issue. Every time the application starts, it must:

1. Read all markdown files from disk
2. Parse markdown syntax into structured data
3. Build the backlink index
4. Extract metadata (tags, TODOs, etc.)

For a 5000-page vault, this can take 3-5 seconds on every startup, which degrades the user experience, especially for users who open and close the app frequently.

### Requirements
- Dramatically reduce startup time for large vaults
- Maintain data integrity (never show stale data)
- Handle failures gracefully (corruption, disk full, etc.)
- Minimal impact on binary size
- Cross-platform compatibility

### Options Considered

1. **No Cache** - Parse everything on each startup
   - ✅ Simple, no complexity
   - ❌ Slow startup for large vaults
   - ❌ Poor user experience

2. **SQLite** - Embedded SQL database
   - ✅ Well-known, battle-tested
   - ✅ Good query capabilities
   - ❌ Requires SQL schema design
   - ❌ Overkill for key-value storage
   - ❌ Larger binary size

3. **BoltDB** - Pure Go key-value store
   - ✅ Simple API
   - ✅ Pure Go (no CGO)
   - ❌ No longer actively maintained
   - ❌ Single writer limitation

4. **BadgerDB** - Modern Go key-value store
   - ✅ High performance (LSM tree)
   - ✅ Pure Go (no CGO)
   - ✅ Actively maintained by Dgraph
   - ✅ Built-in compression
   - ✅ ACID transactions
   - ❌ Larger than BoltDB
   - ❌ More complex than BoltDB

## Decision

We will use **BadgerDB v4** for persistent caching of parsed pages.

### Rationale

1. **Performance**: BadgerDB's LSM tree architecture is optimized for our write pattern (bulk writes during cache build, mostly reads afterward)

2. **Reliability**: ACID transactions ensure cache consistency even if the app crashes during writes

3. **Maintenance**: Actively developed and used in production by Dgraph and others

4. **Features**: Built-in compression reduces cache size, garbage collection prevents unbounded growth

5. **Pure Go**: No CGO dependencies means easy cross-platform builds

## Implementation Details

### Cache Key Structure
```
page:<pagename>      → Cached page data
backlinks:<pagename> → Backlink index for page
cache_metadata       → Version and validation info
```

### Cache Invalidation Strategy
- Store file modification time with each cached entry
- On read, compare current file mtime with cached mtime
- If file is newer, invalidate cache entry

### Cache Location
- Production: `<library-path>/cache/`
- Testing: `<repo-root>/tmp/cache/`

The cache is stored alongside the library data in a visible directory to ensure:
- Each library has its own isolated cache
- Cache moves with the library when relocated
- No conflicts between different libraries
- Easy cache management (just delete the cache folder)
- Transparency - users can see what's being cached

### Failure Handling
- If cache operations fail, fall back to direct parsing
- Log warnings but don't crash
- Provide option to clear cache if corrupted

## Consequences

### Positive
- **2.5-3x faster startup** for large vaults (measured)
- **Transparent to users** - works automatically
- **Graceful degradation** - failures don't break the app
- **Future-proof** - can add more cached data (search index, etc.)

### Negative
- **Increased complexity** - new dependency and subsystem
- **Disk usage** - cache can grow to ~10-20% of vault size
- **Binary size** - adds ~5MB to binary size
- **Potential for bugs** - cache invalidation is notoriously hard

### Mitigations
- Comprehensive test suite including failure scenarios
- Metrics collection to monitor cache health
- Clear documentation for troubleshooting
- Easy cache clearing mechanism

## Metrics

After implementation, we measure:
- Cold start (no cache): 3.7s for 5000 pages
- Warm start (with cache): 1.4s for 5000 pages
- **Speedup: 2.6x faster**
- Cache operations: <10μs per operation
- Cache size: ~50MB for 5000 pages

## References
- [BadgerDB Documentation](https://dgraph.io/docs/badger/)
- [LSM Tree Explanation](https://en.wikipedia.org/wiki/Log-structured_merge-tree)
- [Our BadgerDB Walkthrough](./BADGERDB_WALKTHROUGH.md)