# Memory Leak Patterns

## Common Memory Leaks in Go Applications

### 1. Goroutine Leaks
**Symptoms**:
- Increasing memory usage over time
- Growing number of goroutines
- Application becomes unresponsive

**Common Causes**:
```go
// BAD: Goroutine leak - channel never closed
func LeakyWorker() {
    ch := make(chan int)
    go func() {
        for v := range ch {  // This goroutine never exits
            process(v)
        }
    }()
    // ch is never closed, goroutine runs forever
}

// GOOD: Proper cleanup
func ProperWorker() {
    ch := make(chan int)
    done := make(chan struct{})
    
    go func() {
        for {
            select {
            case v := <-ch:
                process(v)
            case <-done:
                return  // Goroutine exits cleanly
            }
        }
    }()
    
    // Later: close(done) to stop the worker
}
```

**Detection Test**:
```go
func TestGoroutineLeak(t *testing.T) {
    initialGoroutines := runtime.NumGoroutine()
    
    // Run operation that might leak
    for i := 0; i < 100; i++ {
        StartWorker()
    }
    
    // Give time for goroutines to finish
    time.Sleep(100 * time.Millisecond)
    
    finalGoroutines := runtime.NumGoroutine()
    leaked := finalGoroutines - initialGoroutines
    
    if leaked > 0 {
        t.Errorf("Leaked %d goroutines", leaked)
        
        // Debug: print stack traces
        buf := make([]byte, 1<<20)
        runtime.Stack(buf, true)
        t.Logf("Stack traces:\n%s", buf)
    }
}
```

### 2. Channel Buffer Leaks
**Symptoms**:
- Memory grows with buffered channels
- Channels hold references to large objects

**Example**:
```go
// BAD: Unbounded buffer growth
type EventBus struct {
    events chan Event  // Large buffer, never drained
}

func (e *EventBus) Publish(event Event) {
    select {
    case e.events <- event:
    default:
        // Silently drops events - but buffer still holds old ones
    }
}

// GOOD: Bounded buffer with cleanup
type EventBus struct {
    events chan Event
    done   chan struct{}
}

func (e *EventBus) Start() {
    go func() {
        for {
            select {
            case event := <-e.events:
                e.process(event)
            case <-e.done:
                // Drain remaining events
                for len(e.events) > 0 {
                    <-e.events
                }
                return
            }
        }
    }()
}
```

### 3. Timer and Ticker Leaks
**Symptoms**:
- Timers not stopped
- Tickers running indefinitely

**Pattern**:
```go
// BAD: Ticker leak
func LeakyPoller() {
    ticker := time.NewTicker(1 * time.Second)
    
    go func() {
        for range ticker.C {
            poll()
        }
    }()
    // ticker.Stop() never called
}

// GOOD: Proper cleanup
func ProperPoller(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()  // Always stop tickers
    
    for {
        select {
        case <-ticker.C:
            poll()
        case <-ctx.Done():
            return
        }
    }
}
```

### 4. Slice Growing Without Bounds
**Symptoms**:
- Slices keep growing, never shrink
- Historical data never cleaned up

**Example**:
```go
// BAD: Unbounded growth
type Cache struct {
    items []Item
}

func (c *Cache) Add(item Item) {
    c.items = append(c.items, item)  // Never removes old items
}

// GOOD: Bounded cache with eviction
type Cache struct {
    items    []Item
    maxSize  int
}

func (c *Cache) Add(item Item) {
    c.items = append(c.items, item)
    
    // Evict old items
    if len(c.items) > c.maxSize {
        // Remove oldest half
        copy(c.items, c.items[len(c.items)/2:])
        c.items = c.items[:len(c.items)/2]
    }
}
```

### 5. Map Memory Retention
**Symptoms**:
- Maps never shrink even after deleting keys
- Large maps hold memory after clearing

**Pattern**:
```go
// BAD: Map retains memory
var cache = make(map[string]*BigObject)

func ClearCache() {
    for k := range cache {
        delete(cache, k)
    }
    // Map still holds internal buckets
}

// GOOD: Replace map to free memory
func ClearCache() {
    cache = make(map[string]*BigObject)  // Old map can be GC'd
}

// BETTER: Use sync.Pool for temporary objects
var pool = sync.Pool{
    New: func() interface{} {
        return &BigObject{}
    },
}
```

### 6. Context Leaks
**Symptoms**:
- Context with values holds references
- Cancelled contexts not properly handled

**Example**:
```go
// BAD: Context holds large value
func ProcessWithContext(data *LargeData) {
    ctx := context.WithValue(context.Background(), "data", data)
    
    go func() {
        // This goroutine holds reference to data via context
        <-ctx.Done()  // Context never cancelled
    }()
}

// GOOD: Proper context lifecycle
func ProcessWithContext(ctx context.Context, data *LargeData) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()  // Ensure context is cancelled
    
    go func() {
        select {
        case <-time.After(5 * time.Second):
            process(data)
        case <-ctx.Done():
            return  // Clean exit
        }
    }()
}
```

## Memory Leak Detection Tools

### 1. Runtime Memory Stats
```go
func PrintMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc = %v MB\n", m.Alloc / 1024 / 1024)
    fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc / 1024 / 1024)
    fmt.Printf("Sys = %v MB\n", m.Sys / 1024 / 1024)
    fmt.Printf("NumGC = %v\n", m.NumGC)
    fmt.Printf("Goroutines = %v\n", runtime.NumGoroutine())
}
```

### 2. Heap Profiling
```bash
# Run with profiling
go test -memprofile=mem.prof -bench=.

# Analyze profile
go tool pprof mem.prof
> top
> list FunctionName
> web  # Generates visual graph
```

### 3. Leak Detection Test
```go
func TestMemoryLeak(t *testing.T) {
    // Force GC and get baseline
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    // Run operation multiple times
    for i := 0; i < 1000; i++ {
        OperationToTest()
    }
    
    // Force GC and get final stats
    runtime.GC()
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    // Check for significant growth
    growth := m2.Alloc - m1.Alloc
    if growth > 10*1024*1024 { // 10MB threshold
        t.Errorf("Memory grew by %d bytes", growth)
    }
}
```

## Prevention Strategies

### 1. Always Use Defer for Cleanup
```go
func ProcessFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // Guaranteed cleanup
    
    // Process file
    return nil
}
```

### 2. Use Context for Lifecycle Management
```go
func StartWorker(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                doWork()
            case <-ctx.Done():
                return  // Clean shutdown
            }
        }
    }()
}
```

### 3. Implement Proper Shutdown
```go
type Service struct {
    done chan struct{}
    wg   sync.WaitGroup
}

func (s *Service) Start() {
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        for {
            select {
            case <-s.done:
                return
            default:
                s.process()
            }
        }
    }()
}

func (s *Service) Stop() {
    close(s.done)
    s.wg.Wait()  // Wait for goroutines to finish
}
```

## Memory Leak Checklist

- [ ] All goroutines have exit conditions
- [ ] Channels are closed when no longer needed
- [ ] Timers and tickers are stopped
- [ ] Maps are replaced or cleared properly
- [ ] Slices don't grow unbounded
- [ ] File handles are closed
- [ ] Network connections are closed
- [ ] Contexts are cancelled
- [ ] Event listeners are unregistered
- [ ] Circular references are avoided

## Seq2b Specific Concerns

### Frontend Memory
- DOM elements not cleaned up
- Event listeners not removed
- WebSocket connections not closed
- Large images cached indefinitely

### Backend Memory
- Page cache growing without bounds
- Search index never pruned
- Undo/redo history unlimited
- File watchers not stopped

## References
- [Go Memory Management](https://golang.org/doc/effective_go.html#allocation_intro)
- [Finding Goroutine Leaks](https://github.com/uber-go/goleak)
- [pprof Profiling Guide](https://blog.golang.org/pprof)