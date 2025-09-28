# Project Context for Claude

## Project Overview
Building a high-performance Logseq replacement in Go with focus on:
- Performance (faster than existing solutions)
- Security (signed binaries, sandboxed plugins)
- Reliable sync (Git/JJ based, no data loss)
- AI integration as first-class citizen
- Minimal feature set (no bloat)

**PRIMARY TARGET**: Enable existing Logseq users to migrate as soon as possible. Focus on core features that block migration, not nice-to-haves.

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

## Version Control Workflow
- **Primary VCS**: We use Jujutsu (jj) as our primary version control system
- **Git Integration**: Git is used as a colocated repository for GitHub interaction
- **Workflow**: 
  - Make changes and use `jj` commands for local version control
  - Use `git push` to sync with GitHub when needed
  - Avoid git commands for branching/commits - use jj instead

## Commands to Run
- Build GUI: `./scripts/build_seq2b.sh` (creates binary in `bin/`)
- Run GUI Dev: `cd desktop/wails && wails dev`
- Test: `go test ./...`
- Run CLI (testing): `go run tools/cli/main.go [file]`

## Project Structure
```
/desktop/wails     - Desktop GUI application (the main product)
/pkg/parser        - Shared parsing library
/internal
  /storage         - Persistence layer (BadgerDB cache)
/tools             - Internal development tools
  /cli             - Testing CLI (formerly cmd/seq2b)
  /cache-demo      - Cache demonstration
  /benchmark       - Performance testing
/scripts           - Build and utility scripts
/bin               - Production binaries (git-ignored)
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
1. Commit the changes: `jj commit -m "Your commit message"`
2. **CRITICAL**: Move the main bookmark to the new commit: `jj bookmark set main -r @-`
3. Push to GitHub: `jj git push` or `git push origin main`

**Alternative if you want to edit the commit message later:**
1. Create a snapshot: `jj describe -m "Your commit message"` (this updates current @ working copy)
2. To edit later: `jj describe @- -m "New message"`

**Why this happens**: After `jj commit`, you're on a new empty working copy (@), and the commit you just made is at @-. The main bookmark needs to be explicitly moved to @- before pushing, otherwise it stays pointing at the old commit.

**DO NOT**: Try to use `git push` without moving the main bookmark first - git will complain about "not on a branch"

### Current Progress
- [x] Initial project setup
- [x] Step 1.1: Basic file reader
- [x] Step 1.2: Line parser with headers (with tests)
- [ ] Step 1.3: Parse basic markdown (bold, italic, links)

## Next Tasks
See project-docs/PROJECT_PLAN.md for detailed task list and progress tracking.

## Response Formatting
- Always end responses with actual model information in backticks for subtle formatting (use the real model that generated the response, not a hardcoded string)
- Do not start responses with model information

## Writing Guidelines
- When writing blog posts or content, reference project-docs/voice.md for tone and style guidelines

## Code Overview and Explanation Guidelines
See project-docs/EXPLANATION.md for detailed guidelines on:
- Code overview format (Driver and Delegation level)
- One-page chunks to avoid scrolling
- Always including filenames above code snippets
- Visual diagrams and architecture focus

## Test-Driven Development (TDD) for Bug Fixes
When fixing bugs, ALWAYS follow this workflow:
1. **Write a failing test first** that reproduces the bug
2. **Verify the test fails** with the current code
3. **Fix the code** to make the test pass
4. **Verify the test now passes**

This ensures:
- The test actually captures the bug
- The fix actually solves the problem
- We have regression test coverage

## Testing Guidelines
**IMPORTANT**: Always use the `bug-investigator-tester-go` agent when:
- Writing new tests
- Running tests
- Debugging test failures
- Implementing TDD workflow
This agent has specialized tools for Go testing and ensures tests are properly written and executed.

## Agent Memory Boundaries
**IMPORTANT**: Each specialized agent maintains its own persistent memory in `.claude/agent/[agent-name]/`.
Memory segregation rules:
- **DO NOT** read or modify other agents' persistent memory files
- **DO NOT** access files in `.claude/agent/` directories other than your own
- Each agent should **ONLY** access its own memory at `.claude/agent/[agent-name]/persistent_memory.md`
- When using the Task tool to launch agents, they will handle their own memory automatically

Current agent directories:
- `.claude/agent/bug-investigator/` - Testing and bug investigation agent memory

## Architecture Decision Records (ADRs)
**IMPORTANT**: Before making any major architectural changes or design decisions:
1. **Create an ADR** in `/docs/adr/` following the template
2. **Name format**: `NNN-description-YYYY-MM-DD.md` (e.g., `001-testing-strategy-2025-01-18.md`)
3. **Document**:
   - Context and problem statement
   - Decision and rationale
   - Implementation approach
   - Consequences (positive, negative, neutral)
   - Alternatives considered
4. **Get approval** before implementing major changes

ADRs help maintain a clear history of architectural decisions and their reasoning.
See `/docs/adr/template.md` for the standard format.

## Design Principles

### Avoid Magic Solutions
**Prefer explicit scripts over hidden configuration.** If something happens automatically, there should be a clear, visible mechanism showing how it works.

Examples:
- ✅ Good: A build script that explicitly copies binaries to `bin/`
- ❌ Bad: Hidden configuration that magically places files somewhere
- ✅ Good: Clear file paths and operations in scripts
- ❌ Bad: Buried configuration in tool-specific files

This principle ensures maintainability and reduces confusion for new contributors.
