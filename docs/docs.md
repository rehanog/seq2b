---
layout: default
title: Documentation
---

# Documentation

<div class="container" style="padding: 2rem 0;">

## Installation

### Prerequisites
- Go 1.24 or higher
- Node.js 16+ (for desktop GUI)
- Git or Jujutsu (for version control)

### Install Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone and Build
```bash
# Clone the repository
git clone https://github.com/rehanog/seq2b.git
cd seq2b

# Build CLI tool
go build -o seq2b-cli cmd/seq2b/main.go

# Build desktop app
cd desktop/wails
wails build
```

## Usage Guide

### CLI Tool

The CLI tool is perfect for automation, testing, and understanding your knowledge base structure.

#### Basic Usage
```bash
# Parse a single file
./seq2b-cli path/to/note.md

# Parse a directory
./seq2b-cli path/to/pages/
```

#### Output Includes
- **Block Structure**: Hierarchical view of all blocks
- **Backlinks**: Who links to whom
- **Statistics**: Page count, block count, orphans
- **Link Graph**: Visual representation of connections

#### Example Output
```
Directory: testdata/pages
Pages found: 4
==============

Pages:
  Page A - 5 blocks
  Page B - 5 blocks
  Page C - 3 blocks
  Page D - 3 blocks

=== Backlink Analysis ===

Page A:
  ← Referenced by:
    - Page B (in blocks: block-2)
    - Page C (in blocks: block-2)
  → References:
    - Page C (1 times)
    - Page B (2 times)
```

### Desktop GUI

The desktop application provides a native, fast interface for navigating your knowledge base.

#### Navigation
- **Click [[page links]]** to navigate between pages
- **Back button** or press `Escape` to go back
- **Backlinks sidebar** shows all pages linking to current page

#### Keyboard Shortcuts
- `Escape` - Go back to previous page
- More shortcuts coming soon!

#### Features
- **Real-time rendering** of markdown content
- **Proper block indentation** with visual hierarchy
- **Fast page switching** with no lag
- **Native OS integration** for better performance

## Project Structure

```
seq2b/
├── pkg/parser/          # Shared parsing library
│   ├── parser.go       # Main parser logic
│   ├── block.go        # Block structure
│   ├── backlink.go     # Backlink indexing
│   └── multi_file.go   # Directory parsing
├── cmd/seq2b/      # CLI application
│   └── main.go         # CLI entry point
├── desktop/wails/      # Desktop GUI
│   ├── app.go          # Backend logic
│   ├── main.go         # Wails entry point
│   └── frontend/       # Web-based UI
└── mobile/             # Future mobile apps
    ├── ios/
    └── android/
```

## Configuration

Currently, Seq2B works out of the box with sensible defaults. Configuration options are coming in future releases.

### Planned Configuration
- Custom keybindings
- Theme selection
- Parser options
- Plugin settings

## API Reference

### Parser Package

```go
import "github.com/rehanog/seq2b/pkg/parser"

// Parse a single file
result, err := parser.ParseFile(content)

// Parse a directory
result, err := parser.ParseDirectory(dirPath)

// Access parsed data
for pageName, page := range result.Pages {
    fmt.Printf("%s has %d blocks\n", 
        pageName, len(page.AllBlocks))
}

// Get backlinks
backlinks := result.Backlinks.GetBacklinks("Page Name")
```

### Desktop App API

The desktop app exposes these functions to the frontend:

```typescript
// Load a directory
await LoadDirectory(path: string)

// Get page data
await GetPage(name: string)

// Get all pages
await GetPageList()

// Get backlinks for a page
await GetBacklinks(name: string)
```

## Troubleshooting

### Common Issues

**Build fails with "wails: command not found"**
```bash
# Use full path
$(go env GOPATH)/bin/wails build
```

**Cannot find parser package**
```bash
# Ensure you're in the project root
cd /path/to/seq2b
go mod tidy
```

**GUI doesn't show indentation**
- Check browser console for errors
- Ensure CSS is loading correctly
- Try clearing browser cache

### Getting Help

- [GitHub Issues](https://github.com/rehanog/seq2b/issues)
- [Discussions](https://github.com/rehanog/seq2b/discussions)
- Check existing issues before creating new ones

</div>