# Running seq2b

This document describes how to run the different versions of seq2b after the mobile-ready restructure.

## Directory Structure

```
seq2b/
├── pkg/parser/              # 📦 Shared parsing library
├── cmd/seq2b/main.go   # 🖥️ CLI tool
├── desktop/wails/          # 🖥️ Desktop app (Wails)
├── mobile/                 # 📱 Future mobile apps
│   ├── ios/
│   └── android/
└── testdata/               # 🧪 Test data
```

## CLI Tool

The CLI tool is excellent for testing and debugging parser logic.

### Run CLI
```bash
# From project root
go run cmd/seq2b/main.go testdata/library_test_0/pages
```

### CLI Output
- Block structure and hierarchy
- Backlink analysis
- Orphan page detection
- Page relationship summary

### CLI Usage
```bash
# Single file
go run cmd/seq2b/main.go testdata/library_test_0/pages/page-a.md

# Directory
go run cmd/seq2b/main.go testdata/library_test_0/pages
```

## Desktop App (Wails)

The desktop app provides a native GUI with proper block indentation and navigation.

### Development Mode (Hot Reload)
```bash
# Change to desktop directory
cd desktop/wails

# Run with hot reload
$(go env GOPATH)/bin/wails dev
```

**Development Features:**
- Hot reload for frontend changes (HTML/CSS/JS)
- Go backend recompiles automatically
- Browser DevTools available (right-click → Inspect)
- Runs on `http://localhost:5173/`

### Production Build
```bash
# Change to desktop directory
cd desktop/wails

# Build the application
$(go env GOPATH)/bin/wails build

# Run the built application
./build/bin/seq2b-wails.app/Contents/MacOS/seq2b-wails
```

**Production Features:**
- Optimized bundle
- Single executable
- No DevTools
- Better performance

### Alternative Launch (macOS)
```bash
# You can also double-click the .app bundle in Finder
open ./build/bin/seq2b-wails.app
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

## Next Steps

Ready for Phase 3 development:
- Advanced parsing features (properties, tags)
- Storage layer implementation
- Git/JJ sync
- Security features
- AI integration