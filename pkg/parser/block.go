package parser

import (
	"fmt"
	"strings"
	"time"
)

// Block represents a Logseq block (can be multi-line with children)
type Block struct {
	ID       string    // Unique identifier (could be UUID or hash)
	Lines    []Line    // All lines belonging to this block (ordered)
	Children []*Block  // Ordered child blocks
	Parent   *Block    // Parent block (nil for top-level)
	Depth    int       // Nesting depth (0 = top-level)
	
	// Computed properties
	Content     string    // Combined content from all lines
	HTMLContent string    // Rendered HTML (cached)
}

// Page represents a complete Logseq page
type Page struct {
	Name        string
	Title       string    // Page title (usually from first header)
	Blocks      []*Block  // Ordered top-level blocks
	AllBlocks   []*Block  // Flat list of all blocks for easy searching
	
	// Metadata
	Created     time.Time
	Modified    time.Time
}

// ParseResult represents the result of parsing a file
type ParseResult struct {
	Page     *Page
	Lines    []Line    // All parsed lines
	Errors   []error   // Any parsing errors encountered
}

// parseContext is used temporarily during parsing
type parseContext struct {
	line        Line
	indentLevel int    // Calculated from raw line, used once, discarded
}

// calculateIndentLevel counts leading spaces and divides by 2 (Logseq standard)
func calculateIndentLevel(rawLine string) int {
	spaces := 0
	for _, ch := range rawLine {
		if ch == ' ' {
			spaces++
		} else {
			break
		}
	}
	return spaces / 2
}

// ParseFile parses markdown content into a Page with block structure
func ParseFile(content string) (*ParseResult, error) {
	lines := []Line{}
	contexts := []parseContext{}
	
	// Step 1: Parse all lines and extract indent levels
	rawLines := strings.Split(content, "\n")
	for i, rawLine := range rawLines {
		indentLevel := calculateIndentLevel(rawLine)
		line := ParseLine(i+1, strings.TrimSpace(rawLine))
		
		lines = append(lines, line)
		contexts = append(contexts, parseContext{line, indentLevel})
	}
	
	// Step 2: Build block tree from parsed lines
	blocks := BuildBlockTree(contexts)
	
	// Step 3: Create page with all blocks
	page := &Page{
		Blocks:    blocks,
		Created:   time.Now(),
		Modified:  time.Now(),
	}
	
	// Extract title from first header if present
	for _, line := range lines {
		if line.Type == TypeHeader {
			page.Title = line.Content
			break
		}
	}
	
	// Build flat list of all blocks
	page.AllBlocks = page.GetAllBlocks()
	
	return &ParseResult{
		Page:  page,
		Lines: lines,
	}, nil
}

// BuildBlockTree converts flat lines into hierarchical block structure
func BuildBlockTree(contexts []parseContext) []*Block {
	var rootBlocks []*Block
	var blockStack []*Block // Stack to track current nesting
	blockID := 0
	
	for _, ctx := range contexts {
		// Skip empty lines
		if ctx.line.Type == TypeEmpty {
			continue
		}
		
		// Only process block items
		if ctx.line.Type == TypeBlock {
			blockID++
			newBlock := &Block{
				ID:    fmt.Sprintf("block-%d", blockID),
				Lines: []Line{ctx.line},
				Depth: ctx.indentLevel,
			}
			
			// Pop stack until we find the right parent level
			for len(blockStack) > 0 && blockStack[len(blockStack)-1].Depth >= ctx.indentLevel {
				blockStack = blockStack[:len(blockStack)-1]
			}
			
			// Set parent and add as child
			if len(blockStack) > 0 {
				parent := blockStack[len(blockStack)-1]
				newBlock.Parent = parent
				parent.Children = append(parent.Children, newBlock)
			} else {
				// Top-level block
				rootBlocks = append(rootBlocks, newBlock)
			}
			
			// Push onto stack
			blockStack = append(blockStack, newBlock)
			
			// Update content
			newBlock.updateContent()
		}
	}
	
	return rootBlocks
}

// updateContent updates the combined content from all lines
func (b *Block) updateContent() {
	contents := []string{}
	for _, line := range b.Lines {
		contents = append(contents, line.Content)
	}
	b.Content = strings.Join(contents, "\n")
}

// GetAllBlocks returns a flat list of all blocks in the page
func (p *Page) GetAllBlocks() []*Block {
	var allBlocks []*Block
	
	var collectBlocks func([]*Block)
	collectBlocks = func(blocks []*Block) {
		for _, block := range blocks {
			allBlocks = append(allBlocks, block)
			collectBlocks(block.Children)
		}
	}
	
	collectBlocks(p.Blocks)
	return allBlocks
}

// AddChild adds a child block and maintains relationships
func (b *Block) AddChild(child *Block) {
	child.Parent = b
	child.Depth = b.Depth + 1
	b.Children = append(b.Children, child)
}

// GetContent returns the combined content of all lines in the block
func (b *Block) GetContent() string {
	if b.Content == "" {
		b.updateContent()
	}
	return b.Content
}

// RenderHTML renders the block content as HTML
func (b *Block) RenderHTML() string {
	if b.HTMLContent == "" {
		b.HTMLContent = RenderToHTML(b.GetContent())
	}
	return b.HTMLContent
}