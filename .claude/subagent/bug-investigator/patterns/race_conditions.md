# Race Condition Patterns

## Common Race Conditions in Seq2b

### 1. Concurrent Block Editing
**Symptoms**: 
- Lost edits when multiple blocks updated simultaneously
- Inconsistent state after rapid navigation

**Detection**:
```bash
go test -race ./desktop/wails
```

**Example Test**:
```go
func TestConcurrentBlockUpdate(t *testing.T) {
    app := NewApp()
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            app.UpdateBlock(Block{
                UUID: "same-uuid",
                Content: fmt.Sprintf("Content %d", id),
            })
        }(i)
    }
    
    wg.Wait()
    // Check final state is consistent
}
```

### 2. File System Access
**Symptoms**:
- File corruption during simultaneous reads/writes
- Incomplete saves when switching pages quickly

**Prevention**:
- Use file locks or mutex per file
- Implement write-through cache with proper synchronization

### 3. Cache Invalidation
**Symptoms**:
- Stale data served after updates
- Inconsistent backlinks after page edits

**Common Locations**:
- `internal/storage/cache.go`
- Page navigation handlers
- Search index updates

### 4. WebSocket/Event Handling
**Symptoms**:
- Duplicate event processing
- Out-of-order message handling
- Memory leaks from unclosed channels

**Test Pattern**:
```go
func TestEventOrdering(t *testing.T) {
    events := make(chan Event, 100)
    results := make([]int, 0)
    var mu sync.Mutex
    
    // Send events
    for i := 0; i < 100; i++ {
        events <- Event{ID: i}
    }
    
    // Process concurrently
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for e := range events {
                mu.Lock()
                results = append(results, e.ID)
                mu.Unlock()
            }
        }()
    }
    
    close(events)
    wg.Wait()
    
    // Verify all events processed
    if len(results) != 100 {
        t.Errorf("Lost events: got %d, want 100", len(results))
    }
}
```

## Detection Tools

### Go Race Detector
```bash
# Run all tests with race detection
go test -race ./...

# Run specific test with verbose output
go test -race -v -run TestConcurrent ./pkg/parser
```

### Stress Testing
```bash
# Run test multiple times to catch intermittent races
for i in {1..100}; do
    go test -race -run TestBlockUpdate ./desktop/wails || break
done
```

### Memory Sanitizers
```bash
# Check for data races and memory issues
go test -race -msan ./...
```

## Prevention Strategies

### 1. Proper Synchronization
- Always use `sync.Mutex` or `sync.RWMutex` for shared state
- Prefer channels for communication between goroutines
- Use `sync.Once` for one-time initialization

### 2. Immutable Data Structures
- Return copies instead of references
- Use functional update patterns
- Implement copy-on-write semantics

### 3. Actor Model
- Encapsulate state within goroutines
- Communicate via messages only
- No shared memory access

## Known Issues in Codebase

### Frontend-Backend Communication
- **Issue**: Rapid clicks can trigger multiple backend calls
- **Solution**: Implement debouncing and request cancellation

### Page Loading
- **Issue**: Concurrent page loads can corrupt display
- **Solution**: Cancel previous load before starting new one

### Search Index Updates
- **Issue**: Index updates during search can cause panics
- **Solution**: Use RWMutex with proper read/write separation

## Testing Checklist

- [ ] Run with `-race` flag
- [ ] Test with high concurrency (100+ goroutines)
- [ ] Check for deadlocks with timeout tests
- [ ] Verify no goroutine leaks
- [ ] Test rapid user interactions
- [ ] Simulate network delays and failures
- [ ] Check cleanup in error paths

## References
- [Go Data Race Detector](https://golang.org/doc/articles/race_detector.html)
- [Effective Go - Concurrency](https://golang.org/doc/effective_go.html#concurrency)
- [Common Go Concurrency Bugs](https://github.com/golang/go/wiki/CommonMistakes)