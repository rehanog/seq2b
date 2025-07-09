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

## Commit Strategy

### When to Commit
1. After completing each numbered step (1.1, 1.2, etc.)
2. Before major refactoring
3. When switching between different parts of the codebase
4. At natural stopping points in a work session

### Commit Message Format
- Step commits: "Step X.Y: Brief description"
- Feature commits: "Add/Update/Fix: Feature description"
- Refactor commits: "Refactor: What was changed"

### Current Progress
- [x] Initial project setup
- [x] Step 1.1: Basic file reader
- [ ] Step 1.2: Line parser with headers

## Next Tasks
See PROJECT_PLAN.md for detailed task list and progress tracking.