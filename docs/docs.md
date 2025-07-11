---
layout: default
title: Documentation
---

# Documentation

<div class="blog-post-meta">
  Getting Started Guide
</div>

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

# Build desktop GUI
cd desktop/wails
wails build
```

## Usage

### CLI Tool
The command-line interface provides analysis and conversion tools:

```bash
# Analyze a single file
./seq2b-cli path/to/file.md

# Process a directory
./seq2b-cli path/to/directory

# Generate HTML output
./seq2b-cli --html path/to/file.md > output.html
```

### Desktop Application
Launch the native desktop app for interactive editing:

```bash
cd desktop/wails
wails dev  # Development mode
# or
./build/seq2b  # Production build
```

## Features

### Block-Based Structure
seq2b uses an outliner approach where content is organized in hierarchical blocks. Each block can contain text, links, and nested sub-blocks with proper indentation.

### Bidirectional Linking
Create connections between pages using `[[page name]]` syntax. The system automatically tracks forward and backward links, building a comprehensive knowledge graph.

### File Format
All notes are stored as standard Markdown files on your local file system. No proprietary formats, no vendor lock-in. Your data remains accessible with any text editor.

### Performance
Built in Go for exceptional speed. Parse thousands of pages in seconds, not minutes. Real-time link detection and HTML generation.

## Architecture

seq2b follows a modular design:

- **Parser Library** (`pkg/parser/`): Core Markdown processing shared across platforms
- **CLI Tool** (`cmd/seq2b/`): Command-line interface for batch operations
- **Desktop GUI** (`desktop/wails/`): Native desktop application using Wails framework
- **Mobile Apps** (`mobile/`): Future iOS and Android applications

## Development

### Running Tests
```bash
go test ./...
```

### Development Server
```bash
cd desktop/wails
wails dev
```

### Building for Production
```bash
# CLI
go build -o seq2b cmd/seq2b/main.go

# Desktop
cd desktop/wails
wails build
```

## API Reference

### Parser Functions
- `ParseFile(filename string)`: Parse a single Markdown file
- `ParseDirectory(path string)`: Process all `.md` files in a directory
- `BuildBacklinks()`: Generate bidirectional link index
- `ExportHTML()`: Convert to HTML with proper formatting

### Configuration
The parser can be configured for different Logseq compatibility modes and custom block indentation preferences.

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