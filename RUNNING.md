# Running seq2b

This document describes how to run the different versions of seq2b after the mobile-ready restructure.

## Directory Structure

```
seq2b/
├── desktop/wails/          # 🖥️ Desktop GUI app (the main product)
├── pkg/parser/             # 📦 Shared parsing library
├── internal/storage/       # 💾 Cache and persistence layer
├── tools/                  # 🔧 Internal development tools
│   ├── cli/               # Testing CLI
│   ├── cache-demo/        # Cache demonstration
│   └── benchmark/         # Performance testing
├── scripts/               # 📜 Build and utility scripts
├── bin/                   # 📦 Production binaries (git-ignored)
└── testdata/              # 🧪 Test data
```

## Desktop App (Main Product)

The desktop app is the primary seq2b application that end users will use.

### Quick Start
```bash
# Build the app
./scripts/build_seq2b.sh

# Run with test library
./scripts/run_seq2b.sh

# Run with your own library
./bin/seq2b.app/Contents/MacOS/seq2b -library /path/to/your/library
```

### Production Build
```bash
# Build the desktop app and copy to bin/
./scripts/build_seq2b.sh

# The binary will be available at:
# macOS: bin/seq2b.app (run with: ./bin/seq2b.app/Contents/MacOS/seq2b)
# Linux: bin/seq2b
# Windows: bin/seq2b.exe

# Run with a library (required parameter)
# macOS:
./bin/seq2b.app/Contents/MacOS/seq2b -library /path/to/library
# Linux/Windows:
./bin/seq2b -library /path/to/library
```

### Development Mode (Hot Reload)
```bash
# Easy way - uses test library
./scripts/run_seq2b.sh -dev

# Manual way - if you need custom settings
cd desktop/wails
export SEQ2B_LIBRARY_PATH=/path/to/library
wails dev
```

**Development Features:**
- Hot reload for frontend changes (HTML/CSS/JS)
- Go backend recompiles automatically
- Browser DevTools available (right-click → Inspect)
- Runs on `http://localhost:5173/`

## CLI Tool (Testing/Development)

The CLI tool is for testing and debugging parser logic, not for end users.

### Run CLI
```bash
# From project root
go run tools/cli/main.go testdata/library_test_0/pages
```

### CLI Output
- Block structure and hierarchy
- Backlink analysis
- Orphan page detection
- Page relationship summary

### CLI Usage
```bash
# Single file
go run tools/cli/main.go testdata/library_test_0/pages/page-a.md

# Directory
go run tools/cli/main.go testdata/library_test_0/pages
```


## Desktop App Features

- ✅ Block indentation using CSS custom properties
- ✅ Clickable page links with navigation
- ✅ Back button with history
- ✅ Backlinks sidebar
- ✅ Native macOS feel
- ✅ Keyboard shortcuts (Escape = back)

## Mobile Apps (Future)

The mobile directories are prepared for future development:

- `mobile/ios/` - iOS app (React Native/Flutter/native)
- `mobile/android/` - Android app (React Native/Flutter/native)

Both will import the shared `pkg/parser/` library for consistent parsing logic.

## Troubleshooting

### Wails CLI Not Found
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Or use full path
$(go env GOPATH)/bin/wails dev
```

### Import Path Issues
All code now uses the shared parser library:
```go
import "github.com/rehanog/seq2b/pkg/parser"
```

### Test Data Location
Test data is organized in Logseq-style structure:
```
testdata/
└── library_test_0/        # Test library
    ├── pages/            # Markdown pages
    └── assets/           # Images and attachments
```

## Cache System & Performance

seq2b uses BadgerDB for persistent caching to dramatically improve startup times for large vaults.

### Cache Location
- **Production**: `<library-path>/cache/`
- **Testing**: `<repo-root>/tmp/cache/`

The cache is stored alongside your library data in a visible directory. This ensures:
- Each library has its own isolated cache
- Cache moves with the library if you relocate it
- No conflicts between different libraries
- Easy to see and manage cache files

### Viewing Cache Activity

You can watch the cache in action using the demo tool:

```bash
# Build and run the cache demo
go run cmd/cache-demo/main.go

# Watch cache files appear in real-time (in another terminal)
# The cache will be in ./demo-library/cache/
watch -n 1 'ls -la ./demo-library/cache/'
```

The demo will:
1. Show the cache directory location
2. List existing cache files
3. Save some test pages to the cache
4. Show cache files growing in size
5. Demonstrate cache hit/miss behavior
6. Show file modification detection

### Cache Management

```bash
# Clear the cache for a specific library (force rebuild on next run)
rm -rf /path/to/your/library/cache/

# Check cache size for a library
du -sh /path/to/your/library/cache/

# See what's using the cache
lsof | grep seq2b
```

### Performance Monitoring

When running the desktop app, you'll see cache statistics in the console:
```
Cache is valid, using cached data...
Parsed 1000 files in 380ms (cache hits: 1000, misses: 0)
```

First run (cold cache):
- Parses all files
- Builds cache
- Slower startup

Subsequent runs (warm cache):
- Uses cached data
- 2.5-3x faster startup
- Only re-parses modified files

## Benchmarking Tools

```bash
# Generate test vaults of various sizes
go run cmd/generate-test-vault/main.go -pages 1000 -output ./test-vault

# Run performance benchmarks
go run cmd/simple-benchmark/main.go -vault ./test-vault
```

## Next Steps

Ready for Phase 3 development:
- Advanced parsing features (properties, tags)
- ~~Storage layer implementation~~ ✅ Completed with BadgerDB
- Git/JJ sync
- Security features
- AI integration