# Critical Features Before Public Release

## Page Rename Functionality
**Priority: CRITICAL - Must have before public release**

### Requirements
1. When renaming a page, must update ALL references across the entire vault:
   - `[[Old Page Name]]` → `[[New Page Name]]`
   - Date pages: `[[Jan 1st, 2025]]` → `[[2025-01-01]]` (if standardizing)
   - Aliased links: `[text]([[Old Page]])` → `[text]([[New Page]])`

2. Implementation considerations:
   - Must handle case sensitivity changes
   - Must update backlink index
   - Should create redirect/alias for compatibility
   - Must be atomic (all-or-nothing operation)
   - Should handle conflicts (if new name already exists)

3. UI/UX requirements:
   - Rename option in page menu or keyboard shortcut
   - Confirmation dialog showing affected pages
   - Progress indicator for large vaults
   - Undo capability

### Technical Approach
- Parse all pages to find references
- Use existing backlink index for efficiency
- Update both in-memory and on-disk representations
- Consider using Git to track the rename operation

### Testing Requirements
- Test with 10,000+ page vaults
- Test with various link formats
- Test with special characters in page names
- Test rollback scenarios