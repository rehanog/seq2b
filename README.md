# seq2b

A high-performance, cross-platform knowledge management system built in Go, inspired by Logseq. Features a native desktop GUI with proper block indentation, bidirectional linking, and a mobile-ready architecture.

![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.24%2B-blue.svg)
![Platform](https://img.shields.io/badge/platform-macOS%20|%20Linux%20|%20Windows-lightgrey.svg)

## ✨ Features

- **🚀 High Performance**: Written in Go for speed and efficiency
- **🔗 Bidirectional Linking**: Automatic backlink detection and navigation
- **📱 Mobile-Ready Architecture**: Structured for future iOS/Android apps
- **🖥️ Native Desktop GUI**: Built with Wails for native feel
- **🎯 Clean Block Hierarchy**: Proper visual indentation for nested blocks
- **💾 Persistent Cache**: BadgerDB-powered cache for instant startup

## 🎯 Project Goals

- **Performance**: Faster than existing solutions
- **Security**: Signed binaries and sandboxed plugins (coming soon)
- **Reliable Sync**: Git/JJ based with no data loss (coming soon)
- **AI Integration**: First-class AI capabilities (coming soon)
- **Minimal Design**: No feature bloat, focused functionality

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or higher
- Node.js 16+ (for Wails GUI)
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Installation

```bash
# Clone the repository
git clone https://github.com/rehanog/seq2b.git
cd seq2b

# Build the desktop app
./scripts/build_seq2b.sh

# Run the desktop app
# macOS: open bin/seq2b.app
# Linux: ./bin/seq2b
# Windows: bin\seq2b.exe
```

### Development

```bash
# Run in development mode with hot reload
cd desktop/wails
wails dev

# Run tests
go test ./...

# Run the testing CLI
go run tools/cli/main.go testdata/library_test_0/pages
```

## 📖 Usage

### Desktop Application

1. Build the app with `./scripts/build_seq2b.sh`
2. Launch from the `bin/` directory
3. Navigate between pages by clicking [[page links]]
4. Use the back button or press Escape to go back
5. View backlinks in the sidebar
6. Enjoy proper block indentation and instant loading!

## 🏗️ Architecture

```
seq2b/
├── desktop/wails/      # Desktop GUI application
├── pkg/parser/         # Shared parsing library
├── internal/storage/   # Cache and persistence
├── tools/              # Development tools
│   ├── cli/           # Testing CLI
│   └── benchmark/     # Performance tests
├── scripts/           # Build scripts
└── bin/               # Production binaries
```

### Key Components

- **Parser**: Logseq-compatible markdown parser with block support
- **Cache**: BadgerDB-powered persistent cache for instant startup
- **GUI**: Web-based UI in native window (Wails)
- **Backlinks**: Automatic bidirectional link detection

## 🛠️ Development

### Running Tests

```bash
go test ./...
```

### Code Structure

- `pkg/parser/`: Core parsing logic (shared across platforms)
- `cmd/seq2b/`: CLI application
- `desktop/wails/`: Desktop GUI application

### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📋 Roadmap

- [x] **Phase 1**: Core markdown parser with block support
- [x] **Phase 2**: Desktop GUI with Wails
- [ ] **Phase 3**: Advanced parsing (properties, tags, TODOs)
- [ ] **Phase 4**: Persistent storage layer
- [ ] **Phase 5**: Git/JJ sync system
- [ ] **Phase 6**: Security and plugin system
- [ ] **Phase 7**: AI integration
- [ ] **Phase 8**: API and web interface

## 🔗 Links

- [Documentation](https://github.com/rehanog/seq2b/wiki)
- [Issue Tracker](https://github.com/rehanog/seq2b/issues)
- [Discussions](https://github.com/rehanog/seq2b/discussions)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Logseq](https://logseq.com/)
- Built with [Wails](https://wails.io/)
- Written in [Go](https://golang.org/)

---

**Note**: This is an early-stage project. APIs and features may change.