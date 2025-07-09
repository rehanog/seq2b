# Project Context for Claude

## Project Overview
Building a high-performance Logseq replacement in Go with focus on:
- Performance (faster than existing solutions)
- Security (signed binaries, sandboxed plugins)
- Reliable sync (Git/JJ based, no data loss)
- AI integration as first-class citizen
- Minimal feature set (no bloat)

## Current Status
- Go environment set up
- Basic project structure created
- Working on Phase 1: Core Markdown Parser & Data Model

## Key Technical Decisions
- Storage: BadgerDB for embedded key-value store
- Parser: Custom Logseq-compatible markdown parser
- Graph: In-memory with automatic backlink generation
- Sync: Git/JJ for version control
- Plugins: WASM for sandboxing
- AI: Provider-agnostic interface

## Commands to Run
- Build: `go build`
- Test: `go test ./...`
- Run: `go run main.go`

## Project Structure
```
/src
  /parser    - Markdown parsing logic
  /storage   - Persistence layer
  /sync      - Git/JJ integration
  /security  - Code signing, plugin verification
  /ai        - AI provider interfaces
  /ui        - CLI and API
```

## Next Tasks
See PROJECT_PLAN.md for detailed task list and progress tracking.