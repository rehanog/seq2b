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
- Build: `go build -o seq2b cmd/seq2b/main.go`
- Test: `go test ./...`
- Run: `go run cmd/seq2b/main.go [file]`

## Project Structure
```
/cmd/seq2b     - Main application entry point
/internal
  /parser          - Markdown parsing logic
  /storage         - Persistence layer (future)
  /sync            - Git/JJ integration (future)
  /security        - Code signing, plugin verification (future)
  /ai              - AI provider interfaces (future)
/testdata          - Test markdown files
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

**IMPORTANT**: Do NOT add "Generated with Claude Code" or "Co-Authored-By: Claude" to commit messages. Keep them clean and professional.

### JJ (Jujutsu) Workflow - IMPORTANT!
When committing and pushing with jj:
1. Create a snapshot: `jj describe -m "Your commit message"`
2. Commit the changes: `jj commit -m "Your commit message"`
3. **CRITICAL**: Move the main bookmark to the new commit: `jj bookmark set main -r @-`
4. Push to GitHub: `jj git push --branch main`

**Why this happens**: After `jj commit`, you're on a new empty working copy (@), and the commit you just made is at @-. The main bookmark needs to be explicitly moved to @- before pushing, otherwise it stays pointing at the old commit.

### Current Progress
- [x] Initial project setup
- [x] Step 1.1: Basic file reader
- [x] Step 1.2: Line parser with headers (with tests)
- [ ] Step 1.3: Parse basic markdown (bold, italic, links)

## Next Tasks
See PROJECT_PLAN.md for detailed task list and progress tracking.

## Response Formatting
- Always end responses with actual model information in backticks for subtle formatting (use the real model that generated the response, not a hardcoded string)
- Do not start responses with model information

## Writing Guidelines
- When writing blog posts or content, reference voice.md for tone and style guidelines