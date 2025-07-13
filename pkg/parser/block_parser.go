// MIT License
//
// Copyright (c) 2025 Rehan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package parser

import (
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
	Content     string      // Combined content from all lines
	TodoInfo    TodoInfo    // TODO state and checkbox information
	HTMLContent string      // Rendered HTML (cached)
	Segments    []Segment   // Parsed markdown segments (for frontend rendering)
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



// updateContent updates the combined content from all lines
func (b *Block) updateContent() {
	contents := []string{}
	for _, line := range b.Lines {
		contents = append(contents, line.Content)
	}
	b.Content = strings.Join(contents, "\n")
	
	// Use already-parsed TODO information from the first line
	if len(b.Lines) > 0 {
		b.TodoInfo = b.Lines[0].TodoInfo
	}
	
	// Parse markdown segments for frontend rendering
	// Remove TODO prefix if present before parsing segments
	contentForSegments := b.Content
	if b.TodoInfo.TodoState != TodoStateNone || b.TodoInfo.CheckboxState != CheckboxNone {
		contentForSegments = RemoveTodoPrefix(b.Content)
	}
	b.Segments = ParseMarkdownSegments(contentForSegments)
}

// SetContent updates the block's content and reparses it
func (b *Block) SetContent(newContent string) {
	b.Content = newContent
	
	// Update lines by re-parsing them
	lines := strings.Split(newContent, "\n")
	b.Lines = make([]Line, len(lines))
	for i, line := range lines {
		// Re-parse each line to get updated TODO info and references
		b.Lines[i] = ParseLine(i+1, line)
	}
	
	// Use the parsed TODO information from first line
	if len(b.Lines) > 0 {
		b.TodoInfo = b.Lines[0].TodoInfo
	}
	
	// Parse markdown segments for frontend rendering
	// Remove TODO prefix if present before parsing segments
	contentForSegments := b.Content
	if b.TodoInfo.TodoState != TodoStateNone || b.TodoInfo.CheckboxState != CheckboxNone {
		contentForSegments = RemoveTodoPrefix(b.Content)
	}
	b.Segments = ParseMarkdownSegments(contentForSegments)
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
		content := b.GetContent()
		
		// If there's TODO info, render it specially
		if b.TodoInfo.TodoState != TodoStateNone || b.TodoInfo.CheckboxState != CheckboxNone {
			// Remove the TODO/checkbox prefix for clean rendering
			contentWithoutPrefix := RemoveTodoPrefix(content)
			b.HTMLContent = RenderToHTML(contentWithoutPrefix)
		} else {
			b.HTMLContent = RenderToHTML(content)
		}
	}
	return b.HTMLContent
}