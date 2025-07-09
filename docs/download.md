---
layout: default
title: Download
---

# Download

<div class="container" style="padding: 2rem 0;">

## Desktop Application

<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 2rem; margin: 2rem 0;">
  <div class="feature-card" style="text-align: center;">
    <h3>🍎 macOS</h3>
    <p>Universal binary for Intel and Apple Silicon</p>
    <a href="#" class="btn btn-primary" style="display: inline-block; margin-top: 1rem;">Download for macOS</a>
    <p style="font-size: 0.9rem; color: #666; margin-top: 0.5rem;">Requires macOS 10.15+</p>
  </div>
  
  <div class="feature-card" style="text-align: center;">
    <h3>🪟 Windows</h3>
    <p>64-bit installer for Windows</p>
    <a href="#" class="btn btn-primary" style="display: inline-block; margin-top: 1rem;">Download for Windows</a>
    <p style="font-size: 0.9rem; color: #666; margin-top: 0.5rem;">Requires Windows 10+</p>
  </div>
  
  <div class="feature-card" style="text-align: center;">
    <h3>🐧 Linux</h3>
    <p>AppImage for all distributions</p>
    <a href="#" class="btn btn-primary" style="display: inline-block; margin-top: 1rem;">Download for Linux</a>
    <p style="font-size: 0.9rem; color: #666; margin-top: 0.5rem;">64-bit only</p>
  </div>
</div>

<div style="background-color: #fff3cd; padding: 1rem; border-radius: 8px; margin: 2rem 0;">
  <strong>⚠️ Note:</strong> Pre-built binaries coming soon! For now, please build from source.
</div>

## Build from Source

### Quick Build

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

### Platform-Specific Instructions

#### macOS
```bash
# Install Xcode Command Line Tools if needed
xcode-select --install

# Build universal binary
cd desktop/wails
wails build -platform darwin/universal
```

#### Windows
```bash
# Requires MinGW or MSVC
cd desktop/wails
wails build -platform windows/amd64
```

#### Linux
```bash
# Install required dependencies
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev

# Build
cd desktop/wails
wails build -platform linux/amd64
```

## CLI Tool Only

If you only need the command-line tool:

```bash
# Install directly with go
go install github.com/rehanog/seq2b/cmd/seq2b@latest

# Or download and build
git clone https://github.com/rehanog/seq2b.git
cd seq2b
go build -o seq2b-cli cmd/seq2b/main.go
```

## System Requirements

### Minimum Requirements
- **Operating System**: macOS 10.15+, Windows 10+, or Linux (64-bit)
- **RAM**: 512MB
- **Disk Space**: 50MB

### Recommended
- **RAM**: 2GB or more for large knowledge bases
- **SSD**: For best performance with thousands of pages

## Verify Installation

After installation, verify everything is working:

```bash
# Check CLI version
./seq2b-cli --version

# Test with sample data
./seq2b-cli testdata/pages
```

## Updates

Seq2B checks for updates automatically (coming soon). You can also:
- Watch the [GitHub repository](https://github.com/rehanog/seq2b) for releases
- Subscribe to the [RSS feed](/feed.xml)
- Follow development in [Discussions](https://github.com/rehanog/seq2b/discussions)

</div>