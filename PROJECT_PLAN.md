# Logseq Replacement Project Plan

## Project Goals
- **Performance**: Faster than existing solutions using Go
- **Security**: Signed binaries and trusted plugin system
- **Sync**: Git/JJ based with no data loss guarantee
- **AI-First**: Built-in AI capabilities as core feature
- **Minimal**: No feature bloat, focused functionality

## Architecture Overview
- **Language**: Go
- **Storage**: BadgerDB/BoltDB embedded key-value store
- **Sync**: Git/JJ for version control
- **Plugins**: WASM sandboxed environment
- **AI**: Provider-agnostic (local and cloud)
- **UI**: Native GUI (Wails/Fyne/Gio), CLI, optional web UI

## Development Phases

### Phase 1: Core Markdown Parser & Data Model (Step-by-Step)
#### Step 1.1: Basic File Reading
- [x] Read a plain text file into memory
- [x] Print contents to console
- [x] Basic error handling
- [x] Write unit tests for file reading

#### Step 1.2: Simple Line Parser
- [x] Parse file line by line
- [x] Identify headers (# Header)
- [x] Store in simple struct
- [x] Write unit tests for line parser

#### Step 1.3: Basic Markdown Elements
- [x] Parse bold/italic text (**bold**, *italic*)
- [x] Parse links [[page references]]
- [x] Parse bullet points with proper nesting
- [x] Write unit tests for markdown parsing

#### Step 1.4: Block Structure
- [x] Understand Logseq's block concept
- [x] Parse nested blocks (indentation)
- [x] Create parent-child relationships
- [x] Write unit tests for block structure

#### Step 1.5: Backlinks
- [x] Detect [[page references]]
- [x] Build backlink index
- [x] Create bidirectional graph
- [x] Write unit tests for backlink system

### Phase 2: Wails GUI Development (MVP Priority) ✅ COMPLETED
- [x] Evaluate native GUI frameworks (Wails, Fyne, Gio) - **Decision: Wails**
- [x] **2.1: Setup Wails Environment**
  - [x] Verify Wails CLI installation
  - [x] Check system dependencies (Node.js, etc.)
  - [x] Initialize Wails project structure
- [x] **2.2: Backend Migration**
  - [x] Create Wails app context (app.go)
  - [x] Migrate existing parser logic to Wails backend
  - [x] Expose Go functions to frontend:
    - [x] LoadDirectory(path string)
    - [x] GetPage(name string)
    - [x] GetPageList()
    - [x] GetBacklinks(name string)
- [x] **2.3: Frontend Development**
  - [x] Create HTML template for page display
  - [x] Add CSS for block indentation and styling (solved original Fyne issue!)
  - [x] Implement JavaScript for page navigation
  - [x] Add clickable [[page]] link handlers
- [x] **2.4: Core Features**
  - [x] Page navigation between links
  - [x] Back button functionality with history
  - [x] Backlinks sidebar with navigation
  - [x] Native macOS application feel
- [x] **2.5: Polish & Testing**
  - [x] Improve styling and layout
  - [x] Test compilation and functionality
  - [x] Performance optimization (fast builds)
  - [x] Add keyboard shortcuts (Escape = back)

### Phase 2.6: Mobile-Ready Architecture ✅ COMPLETED
- [x] Restructure for multi-platform development:
  - [x] Move parser to `pkg/parser/` for shared library
  - [x] Move Wails app to `desktop/wails/`
  - [x] Create `mobile/ios/` and `mobile/android/` directories
  - [x] Update all import paths to use shared parser
  - [x] Remove dead-end Fyne GUI version
- [x] Add comprehensive usage documentation (RUNNING.md)
- [x] Verify CLI and desktop functionality

### Phase 2.7: Open Source Publishing & Website
- [x] Add MIT license to all source files
- [x] Create main LICENSE file
- [ ] Prepare repository for GitHub
  - [ ] Add comprehensive README.md
  - [ ] Create CONTRIBUTING.md
  - [ ] Add issue templates
  - [ ] Configure GitHub Actions for CI/CD
- [ ] Create GitHub repository
  - [ ] Push code to GitHub
  - [ ] Set up branch protection
  - [ ] Enable GitHub Issues
  - [ ] Configure security alerts
- [ ] Setup GitHub Pages for website
  - [ ] Create docs/ directory
  - [ ] Design landing page
  - [ ] Add documentation
  - [ ] Configure custom domain (optional)
- [ ] Create project blog
  - [ ] Setup Jekyll or Hugo
  - [ ] Write announcement post
  - [ ] Document journey from CLI to GUI
  - [ ] Share architecture decisions
- [ ] Community Setup
  - [ ] Create Discord/Discussions
  - [ ] Add Code of Conduct
  - [ ] Set up sponsorship (GitHub Sponsors)

### Phase 3: Advanced Parsing Features (Deferred)
- [ ] Parse properties (key:: value)
- [ ] Parse tags (#tag)
- [ ] Parse TODO states (TODO, DONE)
- [ ] Implement file watcher
- [ ] Write unit tests for advanced features

### Phase 4: Persistent Storage Layer
- [ ] Design efficient storage format (BadgerDB/BoltDB)
- [ ] Write unit tests for storage interface
- [ ] Implement indexing for fast queries
- [ ] Write unit tests for indexing
- [ ] Create caching layer for performance
- [ ] Write unit tests for caching
- [ ] Add write-ahead logging for data integrity
- [ ] Write integration tests for storage pipeline

### Phase 5: Git/JJ Sync System
- [ ] Implement git integration with go-git
- [ ] Write unit tests for git operations
- [ ] Add jujutsu (jj) support
- [ ] Write unit tests for jj operations
- [ ] Create conflict resolution system
- [ ] Write unit tests for conflict resolution
- [ ] Implement mobile sync protocol
- [ ] Write integration tests for sync pipeline

### Phase 6: Security Implementation
- [ ] Set up code signing for binaries
- [ ] Write tests for signature verification
- [ ] Design plugin verification system
- [ ] Write unit tests for plugin verification
- [ ] Implement plugin sandboxing with WASM
- [ ] Write security tests for WASM sandbox
- [ ] Create capability-based permissions
- [ ] Write unit tests for permission system

### Phase 7: AI Integration
- [ ] Design AI provider interface
- [ ] Write unit tests for AI interface
- [ ] Implement local LLM support (Ollama/llama.cpp)
- [ ] Write unit tests for LLM integration
- [ ] Add semantic search with embeddings
- [ ] Write unit tests for embedding system
- [ ] Create AI-powered linking suggestions
- [ ] Write integration tests for AI features

### Phase 8: API & CLI
- [ ] Create REST/gRPC API for extensions
- [ ] Write API unit tests
- [ ] Write API integration tests
- [ ] Implement comprehensive CLI
- [ ] Write CLI unit tests
- [ ] Write CLI integration tests
- [ ] Optional: Minimal web UI for remote access
- [ ] Write web UI tests

## Technical Decisions

### Data Model
```go
type Page struct {
    ID        string
    Title     string
    Blocks    []Block
    Backlinks []Reference
    Tags      []string
    Modified  time.Time
}

type Block struct {
    ID         string
    Content    string
    Children   []Block
    Properties map[string]string
    References []Reference
}

type Reference struct {
    FromID string
    ToID   string
    Type   RefType // inline, tag, property
}
```

### Performance Targets
- Parse 10,000 pages < 1 second
- Query response < 10ms
- Memory usage < 500MB for 100k blocks
- Instant sync for changes < 1MB

### Security Requirements
- All binaries signed with developer certificate
- Plugins must be cryptographically verified
- WASM sandbox with restricted capabilities
- No network access without explicit permission

### Native GUI Framework Options

#### Option 1: Wails
- **Pros**: Native look/feel, uses system webview, small binary size
- **Cons**: Still uses web technologies (HTML/CSS/JS) for UI
- **Best for**: Rapid development with web skills

#### Option 2: Fyne
- **Pros**: Pure Go, truly native, cross-platform, material design
- **Cons**: Custom look (not system native), larger binaries
- **Best for**: Consistent cross-platform experience

#### Option 3: Gio
- **Pros**: Pure Go, immediate mode, high performance, small
- **Cons**: Lower level, more work for complex UIs
- **Best for**: Maximum performance and control

**Recommendation**: Start with Fyne for faster development, consider Gio later for performance optimization

## Testing Strategy
- Unit tests for each parsing function
- Test files with various markdown edge cases
- Benchmark tests for performance goals
- Integration tests for full pipeline

## Go Testing Conventions
- Test files: `*_test.go` in same package
- Test functions: `TestXxx(t *testing.T)`
- Table-driven tests for multiple cases
- Benchmark functions: `BenchmarkXxx(b *testing.B)`

## Next Steps
1. Set up Go module structure
2. Implement basic markdown parser with tests
3. Create test suite with Logseq sample data
4. Design storage schema