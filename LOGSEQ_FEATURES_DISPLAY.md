# How seq2b Displays Logseq Features

## Overview
When you open Page A or Page B in seq2b, you'll see all Logseq features preserved and displayed with visual indicators.

## What You'll See:

### 1. **Properties** (at top of page)
```
type:: documentation
tags:: example, demo, logseq-features
author:: [[John Doe]]
```
- Displayed with light blue background
- Shows as: `type:: documentation` in a styled box

### 2. **Block IDs**
```
id:: 550e8400-e29b-41d4-a716-446655440001
```
- Displayed in small gray monospace font
- Preserves the full ID for future use

### 3. **Tags**
```
#urgent #bug #critical #project/frontend
```
- Displayed with blue background and rounded corners
- Each tag is individually styled

### 4. **Block References**
```
((550e8400-e29b-41d4-a716-446655440001))
```
- Displayed in purple monospace with tooltip
- Shows as: `((550e8400-e29b-41d4-a716-446655440001))` with hover text "Block reference: [id]"

### 5. **Queries**
```
{{query (todo TODO)}}
```
- Displayed in yellow bordered box
- Shows the full query text as placeholder

### 6. **Embeds**
```
{{embed ((550e8400-e29b-41d4-a716-446655440001))}}
```
- Displayed in yellow bordered box
- Shows the full embed syntax as placeholder

### 7. **Formatting**
- `==highlighted text==` → displayed with yellow background (like HTML `<mark>`)
- `^^caret highlight^^` → displayed with yellow background
- `~~strikethrough~~` → displayed with line through text
- `**bold**` → displayed as bold
- `*italic*` → displayed as italic

### 8. **Extended TODO States**
- `TODO` → Orange badge
- `NOW` → Red badge
- `DOING` → Blue badge
- `WAIT` → Yellow badge
- `LATER` → Gray badge
- `CANCELLED` → Red badge
- `DONE` → Green badge

### 9. **Special Properties**
```
SCHEDULED: <2025-01-14 Tue>
DEADLINE: <2025-01-20 Mon>
collapsed:: true
```
- All preserved and displayed as properties

## Key Features:
1. **Nothing is lost** - All Logseq syntax is preserved
2. **Visual indicators** - Users can see what features are present
3. **Tooltips** - Hover over block references to see IDs
4. **Placeholders** - Complex features show what they are (e.g., "Query: ...")
5. **Graceful degradation** - Features we don't support yet are still visible

## Example Display:

When viewing Page A, a block like:
```
- TODO [#A] Import [[adsa]] with #urgent tag
```

Would display as:
- Bullet point
- Orange "TODO" badge with "[#A]" priority
- "Import" as plain text
- "adsa" as a blue clickable link
- "with" as plain text
- "#urgent" with blue tag styling
- "tag" as plain text

This ensures users can:
1. See all their Logseq content
2. Navigate between pages
3. Edit blocks without losing metadata
4. Understand what features are being used
5. Have confidence their data is preserved