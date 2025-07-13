# Code Explanation Guidelines

## Format Rules

### 1. One Page Size Chunks
- Keep explanations to one screen of information
- User should not need to scroll
- Break complex topics into multiple chunks
- Allow for questions between chunks

### 2. Always Include Filename
- Add filename comment before EVERY code snippet
- Format: `// filename.go` or `// path/to/file.js`
- Even for continued discussion of same file

### 3. Code Overview Level
When asked for "code overview" or "walkthrough", provide **Code Flow Architecture** at **Driver and Delegation** level:

- Identify the "driver" code (main entry point/boss)
- Show calling hierarchy (who calls whom, in what order)
- Explain the "why" of flow
- Show delegation patterns
- Illustrate data transformation stages

### 4. Use Visual Diagrams
- ASCII flow charts for call hierarchies
- Simple box diagrams for architecture
- Tree structures for nested data

### 5. Exclude Implementation Details
- No line-by-line walkthroughs
- Skip error handling specifics
- Avoid parameter details unless crucial
- Focus on the "what" and "why", not "how"

## Example Format

```
## Topic Title

// filename.js
```javascript
function example() {
    // Key concept here
}
```

**What this does**: Brief explanation

**Why it matters**: Context

Ready for next chunk? →
```