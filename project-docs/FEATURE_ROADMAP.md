# seq2b Feature Roadmap

## 🚀 Minimum Viable Release (v1.0)
Features required before announcing seq2b to the world:

### ✅ Already Implemented
- [x] Basic markdown parsing (headers, bold, italic, links)
- [x] Block structure with proper nesting
- [x] Page links `[[Page Name]]` with navigation
- [x] Backlinks (automatic from page references)
- [x] TODO/DONE states with checkboxes
- [x] Date page support (basic)
- [x] Image embedding `![alt](path)`
- [x] Desktop GUI (Wails)
- [x] Inline editing of blocks
- [x] Daily journal pages

### 🔨 Required for v1.0
- [ ] **Import Logseq files without breaking**
  - [ ] Parse and preserve block IDs `id:: UUID`
  - [ ] Handle block references `((block-id))` without crashing
  - [ ] Parse properties `property:: value` without breaking
  - [ ] Handle tags `#tag` gracefully
  - [ ] Support extended TODO states (NOW, DOING, LATER, etc.)
  - [ ] Parse but ignore unsupported syntax (queries, embeds, etc.)
- [ ] **Basic PDF support**
  - [ ] View PDFs via `[[file.pdf]]` links
  - [ ] Basic navigation (page up/down, zoom)
  - [ ] Remember last viewed page
- [ ] **Data integrity**
  - [ ] Never lose user data on import
  - [ ] Preserve original file structure
  - [ ] Handle malformed syntax gracefully

## 🔜 Coming Soon (Post-Launch)
Features we plan to add after initial release:

### High Priority
- [ ] Unlinked references (text mentions)
- [ ] Block references `((block-id))` - full implementation
- [ ] Tags `#tag` - full implementation with tag pages
- [ ] Properties - full implementation with queries
- [ ] Search functionality
- [ ] Block embeds `{{embed ((id))}}`

### Medium Priority
- [ ] Extended markdown (tables, code blocks, quotes)
- [ ] Page properties display
- [ ] Scheduled/deadline dates
- [ ] Git sync integration
- [ ] Page aliases
- [ ] Basic queries `{{query}}`
- [ ] Strikethrough `~~text~~`
- [ ] Highlights `==text==`

### Low Priority
- [ ] PDF annotations and highlights
- [ ] Flashcards and spaced repetition
- [ ] Templates
- [ ] Custom themes
- [ ] Export functionality

## ❌ No Plans to Implement
Features that don't align with seq2b's philosophy:

- **Whiteboards** - Different paradigm, adds complexity
- **Org-mode syntax** - Focusing on Markdown only
- **Plugin system** - At least not initially, to maintain performance
- **Excalidraw integration** - Use external tools instead
- **Video/audio embeds** - Link to external players
- **Real-time collaboration** - Single-user focus
- **Cloud sync service** - Git/JJ based sync only

## 📋 Implementation Notes

### Import Strategy
For v1.0, we need to ensure Logseq files can be imported without breaking seq2b:
1. **Parse but don't render** unsupported features initially
2. **Preserve original content** even if we can't display it properly
3. **Show placeholders** for complex features (e.g., "Query: {{query}}")
4. **Log warnings** for unsupported syntax but don't crash

### PDF Strategy
Start with basic PDF viewing:
1. Use PDF.js or native viewer
2. Simple navigation controls
3. Link from markdown: `[[document.pdf]]`
4. Open in panel or modal
5. Later: annotations, highlights, text extraction

### Success Metrics for v1.0
- [ ] Can import any Logseq graph without data loss
- [ ] Can view and navigate between all pages
- [ ] Can edit blocks without breaking references
- [ ] Can view PDFs linked from pages
- [ ] Performance: <1s to load 1000 pages
- [ ] Zero crashes on malformed input