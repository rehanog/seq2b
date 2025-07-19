# Logseq Feature Compatibility Matrix

## Analysis Date: 2025-01-19
Based on analysis of candideu/Logseq-Demo-Graph

## Feature Support Status

### ✅ Fully Supported Features

| Feature | Logseq Syntax | seq2b Status | Notes |
|---------|---------------|--------------|-------|
| Headers | `# Header` through `##### Header` | ✅ Fully supported | All heading levels work |
| Bold text | `**bold**` | ✅ Fully supported | |
| Italic text | `*italic*` or `_italic_` | ✅ Fully supported | |
| Bold+Italic | `***text***` | ✅ Fully supported | |
| Page links | `[[Page Name]]` | ✅ Fully supported | With navigation |
| External links | `[text](url)` | ✅ Fully supported | |
| Nested blocks | Indentation with `-` | ✅ Fully supported | Proper hierarchy |
| TODO states | `TODO`, `DONE` | ✅ Fully supported | Basic states only |
| Checkboxes | `- [ ]`, `- [x]` | ✅ Fully supported | With nesting |
| Images | `![alt](path)` | ✅ Fully supported | Basic syntax |
| Backlinks | Automatic from `[[links]]` | ✅ Fully supported | At bottom of page |

### ⚠️ Partially Supported Features

| Feature | Logseq Syntax | seq2b Status | Missing |
|---------|---------------|--------------|---------|
| Date pages | `[[2025-01-19]]`, `[[Jan 19th, 2025]]` | ⚠️ Partial | Limited date format support |
| TODO states | `DOING`, `NOW`, `WAIT`, `LATER`, `CANCELLED` | ⚠️ Partial | Only TODO/DONE supported |

### ❌ Not Supported Features

| Feature | Logseq Syntax | Priority | Notes |
|---------|---------------|----------|-------|
| Block IDs | `id:: UUID` | High | Core for references |
| Block references | `((block-id))` | High | Essential feature |
| Block embeds | `{{embed ((id))}}` | High | Transclusion |
| Page embeds | `{{embed [[page]]}}` | Medium | |
| Tags | `#tag` | High | Inline tags |
| Properties | `property:: value` | High | Page/block metadata |
| Scheduled dates | `SCHEDULED: <2025-01-19>` | Medium | |
| Deadline dates | `DEADLINE: <2025-01-19>` | Medium | |
| Collapsed blocks | `collapsed:: true` | Medium | UI state |
| Strikethrough | `~~text~~` | Low | |
| Highlight | `==text==` or `^^text^^` | Low | |
| Block quotes | `> quote` | Low | |
| Inline code | `` `code` `` | Low | |
| Code blocks | ` ```language\ncode\n``` ` | Medium | |
| LaTeX | `$$formula$$` | Low | Math support |
| Tables | `| col1 | col2 |` | Medium | Markdown tables |
| Numbered lists | `1. item` | Medium | |
| Horizontal rule | `---` | Low | |
| Queries | `{{query}}` | High | Dynamic content |
| Cloze deletion | `{{cloze text}}` | Low | Flashcards |
| Flashcard tag | `#card` | Low | Study features |
| YouTube embeds | `{{video url}}` | Low | |
| YouTube timestamp | `{{youtube-timestamp 123}}` | Low | |
| PDF embeds | `![pdf](url.pdf)` | Medium | Phase 8 planned |
| PDF highlights | `((pdf-highlight-id))` | Low | |
| Image attributes | `{:height 200, :width 192}` | Low | |
| Alias links | `[alias]([[page]])` | Medium | |
| Page hierarchy | `parent/child` pages | Medium | |
| Multiple property values | `tags:: value1, value2` | Medium | |
| Unlinked references | Text mentions | Medium | Phase 5 planned |
| Numbered list type | `logseq.order-list-type:: number` | Low | |
| Org-mode blocks | `#+BEGIN_TIP...#+END_TIP` | N/A | Out of scope - MD only |
| Whiteboards | Separate feature | Low | Different paradigm |
| Slash commands | `/command` | Medium | |
| Templates | `{{template}}` | Low | |
| Macros | `{{macro}}` | Low | |
| Favorites | Page favorites | Low | |
| Contents sidebar | Right sidebar | Low | |
| Re-indexing | Manual reindex | Low | |
| Settings UI | Configuration | Low | |
| Plugins | Plugin API | Low | Phase 10 planned |

## Implementation Priority

### Phase 6: Immediate Priorities (Core Logseq Features)
1. **Block IDs & References** - Essential for Logseq compatibility
   - Generate and store block IDs
   - Parse and render block references
   - Implement block embeds
2. **Tags** - Core organizational feature
   - Parse inline #tags
   - Create tag pages
   - Tag queries
3. **Properties** - Metadata system
   - Parse key:: value syntax
   - Support page and block properties
   - Handle special properties (collapsed, scheduled, etc.)

### Phase 6+: Secondary Priorities
1. **Extended TODO States** - DOING, NOW, WAIT, LATER, CANCELLED
2. **Basic Queries** - {{query}} for TODOs and tags
3. **Markdown Extensions** - Tables, quotes, code blocks
4. **Page Features** - Hierarchies, aliases, unlinked refs

### Future Phases
- PDF support (Phase 8)
- Plugin system (Phase 10)
- Advanced features (whiteboards, templates, macros)

## Test Coverage Needed

### High Priority Test Cases
1. Block reference cycles
2. Nested embeds
3. Property inheritance
4. Query performance with large graphs
5. Tag indexing and search
6. Date format parsing edge cases

### Sample Test Files Needed
1. `block-references-test.md` - Complex ref patterns
2. `properties-test.md` - All property types
3. `queries-test.md` - Various query examples
4. `tags-test.md` - Tag hierarchies and searches
5. `edge-cases-test.md` - Malformed syntax handling