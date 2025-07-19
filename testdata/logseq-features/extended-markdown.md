# Extended Markdown Features

## Text Formatting
- ~~Strikethrough text~~
- ==Highlighted text with equals==
- ^^Highlighted text with carets^^
- ~~==Combined strikethrough and highlight==~~

## Block Quotes
> This is a block quote
> It can span multiple lines
> > And can be nested
> > > Multiple levels deep

- Regular text
> Single line quote

## Code
- Inline code: `const x = 42;`
- Inline with backticks: `npm install`

## Code Blocks
```javascript
function hello(name) {
  console.log(`Hello, ${name}!`);
}
```

```python
def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)
```

```
Plain code block without language
No syntax highlighting
```

## LaTeX Math
- Inline math: $E = mc^2$
- Display math: $$\frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$

Complex equation:
$$
\begin{align}
\nabla \times \vec{E} &= -\frac{\partial \vec{B}}{\partial t} \\
\nabla \times \vec{B} &= \mu_0 \vec{J} + \mu_0 \epsilon_0 \frac{\partial \vec{E}}{\partial t}
\end{align}
$$

## Tables
| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |

| Left aligned | Center aligned | Right aligned |
|:-------------|:--------------:|--------------:|
| Left         | Center         | Right         |
| 123          | 456            | 789           |

| Feature | Supported | Priority |
|---------|-----------|----------|
| Tables  | ❌        | Medium   |
| LaTeX   | ❌        | Low      |
| Code    | ❌        | Medium   |

## Lists

### Ordered Lists
1. First item
2. Second item
3. Third item
   1. Nested item
   2. Another nested
      1. Deep nesting

### Mixed Lists
1. Ordered item
   - Unordered sub-item
   - Another sub-item
2. Second ordered
   - [ ] Checkbox in ordered list
   - [x] Completed checkbox

### Definition Lists (if supported)
Term 1
: Definition 1

Term 2
: Definition 2a
: Definition 2b

## Horizontal Rules
---
***
___

## Footnotes
This text has a footnote[^1].

[^1]: This is the footnote content.

## HTML (if allowed)
<details>
<summary>Click to expand</summary>
This is hidden content that can be revealed.
</details>

<mark>Highlighted with HTML</mark>

## Org-Mode Blocks (Out of Scope)
Note: seq2b focuses on Markdown files only. Org-mode syntax is not supported.

## Special Logseq Syntax

### Cloze Deletion
- This is a {{cloze cloze}} for spaced repetition

### YouTube Timestamps
- {{youtube-timestamp 125}} - Jump to 2:05 in video

### Draws/Excalidraw
- [[draws/2025-01-19-12-34-56.excalidraw]]

### Templates
- {{template my-template}}

### Macros
- {{macro-name arg1 arg2}}