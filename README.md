# Logseq-Go

A high-performance Logseq replacement written in Go.

## Features (Planned)

- **Performance**: Native Go implementation, no Electron overhead
- **Security**: Signed binaries, sandboxed plugins via WASM
- **Sync**: Git/JJ based synchronization with zero data loss
- **AI Integration**: First-class AI support with provider-agnostic interface
- **Native GUI**: True native performance and look/feel
- **Minimal**: No feature bloat, focused on core functionality

## Current Status

Early development - building the markdown parser.

## Development

```bash
# Run tests
go test ./...

# Run the parser
go run cmd/logseq-go/main.go [markdown-file]

# Build
go build -o logseq-go cmd/logseq-go/main.go
```

## Project Structure

```
logseq-go/
├── cmd/logseq-go/      # Main application entry point
├── internal/parser/    # Markdown parsing logic
├── testdata/          # Test markdown files
└── docs/              # Documentation
```

See [PROJECT_PLAN.md](PROJECT_PLAN.md) for detailed development roadmap.