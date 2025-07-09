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
- [ ] Read a plain text file into memory
- [ ] Print contents to console
- [ ] Basic error handling

#### Step 1.2: Simple Line Parser
- [ ] Parse file line by line
- [ ] Identify headers (# Header)
- [ ] Store in simple struct

#### Step 1.3: Basic Markdown Elements
- [ ] Parse bold/italic text
- [ ] Parse links [[page references]]
- [ ] Parse bullet points

#### Step 1.4: Block Structure
- [ ] Understand Logseq's block concept
- [ ] Parse nested blocks (indentation)
- [ ] Create parent-child relationships

#### Step 1.5: Page Model
- [ ] Design Page and Block structs
- [ ] Parse multiple pages
- [ ] Create in-memory storage

#### Step 1.6: Backlinks
- [ ] Detect [[page references]]
- [ ] Build backlink index
- [ ] Create bidirectional graph

#### Step 1.7: Advanced Features
- [ ] Parse properties
- [ ] Parse tags
- [ ] Parse TODO states
- [ ] Implement file watcher

### Phase 2: Persistent Storage Layer
- [ ] Design efficient storage format (BadgerDB/BoltDB)
- [ ] Implement indexing for fast queries
- [ ] Create caching layer for performance
- [ ] Add write-ahead logging for data integrity

### Phase 3: Git/JJ Sync System
- [ ] Implement git integration with go-git
- [ ] Add jujutsu (jj) support
- [ ] Create conflict resolution system
- [ ] Implement mobile sync protocol

### Phase 4: Security Implementation
- [ ] Set up code signing for binaries
- [ ] Design plugin verification system
- [ ] Implement plugin sandboxing with WASM
- [ ] Create capability-based permissions

### Phase 5: AI Integration
- [ ] Design AI provider interface
- [ ] Implement local LLM support (Ollama/llama.cpp)
- [ ] Add semantic search with embeddings
- [ ] Create AI-powered linking suggestions

### Phase 6: Native GUI Development
- [ ] Evaluate native GUI frameworks (Wails, Fyne, Gio)
- [ ] Design native UI/UX following platform guidelines
- [ ] Implement core editor view with native performance
- [ ] Add native file browser and search
- [ ] Create platform-specific installers

### Phase 7: API & CLI
- [ ] Create REST/gRPC API for extensions
- [ ] Implement comprehensive CLI
- [ ] Optional: Minimal web UI for remote access

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

## Next Steps
1. Set up Go module structure
2. Implement basic markdown parser
3. Create test suite with Logseq sample data
4. Design storage schema