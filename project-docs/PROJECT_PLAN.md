# seq2b Project Plan

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

### Phase 2.7: Open Source Publishing & Website ✅ MOSTLY COMPLETED
- [x] Add MIT license to all source files
- [x] Create main LICENSE file
- [x] Prepare repository for GitHub
  - [x] Add comprehensive README.md
  - [ ] Create CONTRIBUTING.md
  - [ ] Add issue templates
  - [ ] Configure GitHub Actions for CI/CD
- [x] Create GitHub repository
  - [x] Push code to GitHub
  - [ ] Set up branch protection
  - [x] Enable GitHub Issues
  - [ ] Configure security alerts
- [x] Setup GitHub Pages for website
  - [x] Create docs/ directory
  - [x] Design landing page
  - [x] Add documentation
  - [ ] Configure custom domain (optional)
- [x] Create project blog
  - [x] Setup Jekyll
  - [x] Write announcement post
  - [x] Document journey from CLI to GUI
  - [x] Share architecture decisions
- [ ] Community Setup
  - [ ] Create Discord/Discussions
  - [ ] Add Code of Conduct
  - [ ] Set up sponsorship (GitHub Sponsors)

### Phase 3: Basic Editing Capability

#### Step 3.1: TODO Lists ✅ COMPLETED
- [x] Parse TODO/DONE/WAITING states in blocks
- [x] Parse checkboxes [ ]/[x] with proper nesting
- [x] Track completion status and inheritance
- [x] Add TODO filtering and views
- [x] Write unit tests for TODO parsing

#### Step 3.2: Basic Block Editing ✅ COMPLETED
- [x] Add edit mode for individual blocks
- [x] Create simple text editor component
- [x] Handle save/cancel operations
- [x] Update parser to write changes back
- [x] Maintain block IDs during edits
- [x] Write unit tests for editing

#### Phase 3 Goal
Enable basic editing of blocks so users can modify their notes without external editors.

### Phase 3.5: Technical Debt Refactoring ✅ COMPLETED

#### Parser Architecture Refactoring ✅ COMPLETED
- [x] Restructure parser files for single responsibility:
  - [x] Create file_parser.go - Contains ParseFile() and ParseDirectory() orchestration
  - [x] Create line_parser.go - Full parsing of line-level features (TODO, tags, etc.)
  - [x] Rename block.go to block_parser.go - Pure structural organization of pre-parsed lines
  - [x] Create markdown_renderer.go - Move RenderToHTML() and formatting functions
- [x] Update parsing flow:
  - [x] Line parser fully parses line-level features in initial pass
  - [x] Line struct carries parsed data (TodoInfo, etc.)
  - [x] Block parser reorganizes pre-parsed lines into hierarchical structure
  - [x] Blocks reference already-parsed data (no re-parsing needed)
- [x] Benefits: Clean separation of concerns, single-pass parsing, better performance

#### Additional Refactoring ✅ COMPLETED
- [x] Separate parsing from rendering:
  - [x] Change parser to output structured segments (text, bold, link, etc.) instead of HTML
  - [x] Remove RenderToHTML() from Go parser
  - [x] Move HTML generation to frontend JavaScript
  - [x] Benefits: Clean separation of concerns, flexible rendering, testable parsing
- [x] Implement incremental updates for editing:
  - [x] Return edit deltas from backend (added/removed/updated blocks)
  - [x] Update only affected DOM elements instead of full page reload
  - [x] Handle structural changes (block splits, merges, indentation)
  - [x] Update backlinks incrementally
  - [x] Benefits: Much faster editing, no UI flicker, maintains scroll position

#### Phase 3.5 Goal ✅ ACHIEVED
Clean up technical debt and improve code efficiency before adding new features.

### Phase 4: Daily Driver Features

#### Step 4.1: Date Pages ✅ COMPLETED
- [x] Parse dates in standard format (YYYY-MM-DD, [[Jan 1st, 2025]], etc.)
- [x] Auto-create date pages like Logseq
- [x] Handle date page navigation and linking
- [x] Support journal-style daily notes
- [x] Write unit tests for date parsing

#### Step 4.2: Home Page with Today's Date ✅ COMPLETED
- [x] Default to today's date page on startup
- [x] Add "Home" button in GUI to return to today
- [x] Auto-create today's page if it doesn't exist
- [x] Handle date page formatting and structure
- [x] Write unit tests for home page logic (reused date parser tests)

#### Step 4.3: Embedded Images ✅ COMPLETED
- [x] Parse image markdown syntax ![alt](path/to/image.png)
- [x] Handle relative and absolute image paths
- [x] Add image rendering in GUI
- [x] Support common image formats (PNG, JPG, GIF, SVG)
- [x] Write unit tests for image parsing

#### Phase 4 Goal
Complete the minimum viable daily driver with:
- Daily journaling with date pages
- Image embedding for visual notes
- Combined with existing features (blocks, backlinks, TODOs, editing)

### Phase 5: Logseq-like Page Structure 🚧 IN PROGRESS
**Goal**: Make the UI more closely match Logseq's page layout and interaction patterns

#### Step 5.1: Move Backlinks to Page Bottom ✅ COMPLETED
- [x] Remove separate backlinks sidebar
- [x] Add "Linked References" section at bottom of page
- [x] Style with subtle separator (not bordered box)
- [x] Show source page and block context
- [x] Make references clickable for navigation

#### Step 5.2: Page Properties Display
- [ ] Parse page-level properties (tags::, alias::, etc.)
- [ ] Display properties at top of page
- [ ] Support both YAML frontmatter and key:: value syntax
- [ ] Make properties editable

#### Step 5.3: Enhanced Block Rendering (Critical for Migration)
- [ ] Add collapse/expand for blocks with children
- [ ] Show block reference count indicators
- [ ] Improve visual hierarchy for deep nesting
- [ ] Add block actions menu (on hover)

### Phase 6: Logseq Markdown Compatibility - Safe Import
**Goal**: Import Logseq files without breaking - parse but don't necessarily render all features

#### Key Strategy
- **Parse but don't render** unsupported features initially
- **Preserve original content** even if we can't display it properly
- **Show placeholders** for complex features (e.g., "Query: {{query}}")
- **Log warnings** for unsupported syntax but don't crash

#### Step 6.1: Compatibility Assessment ✅ COMPLETED
- [x] Create comprehensive feature matrix (LOGSEQ_COMPATIBILITY.md)
- [x] Test against real Logseq markdown files
- [x] Document differences and limitations
- [x] Prioritize missing features by usage

#### Step 6.2: Safe Import - Block IDs & References
- [ ] Parse `id:: UUID` without breaking
- [ ] Store block IDs but don't generate new ones
- [ ] Parse `((block-id))` as plain text or placeholder
- [ ] Show placeholder: "[Block reference: block-id]"
- [ ] Don't crash on invalid references

#### Step 6.3: Safe Import - Properties & Tags
- [ ] Parse `property:: value` syntax
- [ ] Store properties in block metadata
- [ ] Parse `#tag` syntax without breaking
- [ ] Display tags as plain text initially
- [ ] Handle SCHEDULED: and DEADLINE: dates
- [ ] Parse extended TODO states (NOW, DOING, LATER, etc.)

#### Step 6.4: Safe Import - Complex Features
- [ ] Parse `{{query}}` blocks → show "Query: [query text]"
- [ ] Parse `{{embed}}` blocks → show "Embed: [content]"
- [ ] Parse tables → preserve as code blocks initially
- [ ] Parse `~~strikethrough~~` → show as plain text
- [ ] Parse `==highlight==` → show as plain text
- [ ] Log all unsupported features for user awareness

### Phase 7: Persistence & Performance
**Goal**: Quick startup without full reparse + performance benchmarks

#### Step 7.1: Parsed Data Cache
- [ ] Design cache schema for parsed pages
- [ ] Implement with BadgerDB/BoltDB
- [ ] Compare file timestamps on startup
- [ ] Only reparse modified files
- [ ] Handle cache invalidation

#### Step 7.2: Performance Benchmarking
- [ ] Create test vault generator (1000+ pages)
- [ ] Benchmark cold startup time
- [ ] Benchmark warm startup (with cache)
- [ ] Measure page navigation speed
- [ ] Add regression test suite

#### Step 7.3: Lazy Loading Implementation
- [ ] Parse pages on first access
- [ ] Background parsing queue
- [ ] Progress indicator for parsing
- [ ] Prioritize recent/pinned pages
- [ ] Memory usage optimization

### Phase 8: PDF Integration
**Goal**: View and annotate PDFs within seq2b

#### Step 8.1: PDF Viewing
- [ ] Integrate PDF.js or native viewer
- [ ] Support [[pdf-file.pdf]] links
- [ ] Page number references [[pdf.pdf#page=5]]
- [ ] Zoom and navigation controls

#### Step 8.2: PDF Annotations
- [ ] Highlight text in PDFs
- [ ] Add margin notes
- [ ] Export annotations as blocks
- [ ] Link blocks to PDF locations
- [ ] Search PDF text content

### Phase 9: Git/JJ Sync System
- [ ] Implement git integration with go-git
- [ ] Add jujutsu (jj) support
- [ ] Create conflict resolution UI
- [ ] Auto-commit on changes
- [ ] Sync status indicators
- [ ] Handle merge conflicts gracefully

### Phase 10: Plugin System & Security
- [ ] Design plugin API surface
- [ ] WASM sandbox implementation
- [ ] Plugin marketplace/registry
- [ ] Code signing for plugins
- [ ] Capability-based permissions
- [ ] Plugin settings management

### Phase 11: AI Integration
- [ ] Provider-agnostic AI interface
- [ ] Local LLM support (Ollama)
- [ ] Cloud provider integration
- [ ] Smart linking suggestions
- [ ] Semantic search with embeddings
- [ ] AI-powered block completion

### Phase 12: Mobile & API
- [ ] iOS app with shared Go core
- [ ] Android app development
- [ ] REST API for third-party apps
- [ ] Sync protocol for mobile
- [ ] Offline-first architecture
- [ ] WebDAV support

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

## Critical Features for Logseq Migration

These are the absolute minimum features required for existing Logseq users to migrate:

1. **Phase 6: Logseq Markdown Compatibility** - Must be able to import existing files without data loss
2. **Phase 7: Persistence & Performance** - Must handle large vaults (10,000+ pages) without stuttering
3. **Local Git Versioning** - Auto-commit changes locally (without sync) for version history
4. **Block collapse/expand** (Phase 5.3) - Core navigation feature users expect

Non-critical but expected:
- Page properties (Phase 5.2) - Nice for organization but not blocking
- PDF annotations (Phase 8.2) - Important for academic users but not universal

## Future Release Features

### Future Release 1: Sync
- **Git/JJ Sync System** (was Phase 9)
  - Remote sync with conflict resolution
  - Multi-device access

### Future Release 2: Extensibility
- **Plugin System** (was Phase 10)
  - WASM sandbox
  - Plugin marketplace
  - Custom themes

### Future Release 3: Enhanced Discovery
- **Unlinked References**
  - Search for text mentions
  - Convert to linked references
- **Graph View**
  - Visual page connections
  - Interactive navigation
- **Advanced Search**
  - Full-text with highlighting
  - Date range filtering

## Next Steps
1. Complete Phase 5 (skip unlinked references)
2. Focus on Phase 6 (Logseq compatibility) 
3. Implement local Git versioning (simplified from Phase 9)
4. Complete Phase 7 (Performance for 10k+ pages)