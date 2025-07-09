package parser

import (
	"regexp"
	"strings"
)

// BacklinkIndex tracks page references across multiple pages
type BacklinkIndex struct {
	// Forward links: source page -> target pages it references
	ForwardLinks map[string]map[string][]BlockReference
	
	// Backward links: target page -> source pages that reference it
	BackwardLinks map[string]map[string][]BlockReference
}

// BlockReference records where a page reference appears
type BlockReference struct {
	PageName string // The page containing this reference
	BlockID  string
	Position int // character position in block content
}

// NewBacklinkIndex creates a new empty backlink index
func NewBacklinkIndex() *BacklinkIndex {
	return &BacklinkIndex{
		ForwardLinks:  make(map[string]map[string][]BlockReference),
		BackwardLinks: make(map[string]map[string][]BlockReference),
	}
}

// ExtractPageLinks finds all [[page]] references in text
func ExtractPageLinks(text string) []string {
	linkPattern := regexp.MustCompile(`\[\[(.*?)\]\]`)
	matches := linkPattern.FindAllStringSubmatch(text, -1)
	
	links := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			links = append(links, match[1])
		}
	}
	
	return links
}

// AddPage adds a single page to the backlink index
func (idx *BacklinkIndex) AddPage(page *Page) {
	pageName := page.Title
	
	// Initialize forward links for this page if needed
	if idx.ForwardLinks[pageName] == nil {
		idx.ForwardLinks[pageName] = make(map[string][]BlockReference)
	}
	
	// Scan all blocks for references
	for _, block := range page.AllBlocks {
		links := ExtractPageLinks(block.Content)
		
		for _, targetPage := range links {
			// Skip self-references
			if targetPage == pageName {
				continue
			}
			
			// Find position of this link in the content
			pos := strings.Index(block.Content, "[["+targetPage+"]]")
			
			ref := BlockReference{
				PageName: pageName,
				BlockID:  block.ID,
				Position: pos,
			}
			
			// Add forward link
			idx.ForwardLinks[pageName][targetPage] = append(
				idx.ForwardLinks[pageName][targetPage], ref)
			
			// Add backward link
			if idx.BackwardLinks[targetPage] == nil {
				idx.BackwardLinks[targetPage] = make(map[string][]BlockReference)
			}
			idx.BackwardLinks[targetPage][pageName] = append(
				idx.BackwardLinks[targetPage][pageName], ref)
		}
	}
}

// GetBacklinks returns all pages that link TO the given page
func (idx *BacklinkIndex) GetBacklinks(pageName string) map[string][]BlockReference {
	return idx.BackwardLinks[pageName]
}

// GetForwardLinks returns all pages that this page links TO
func (idx *BacklinkIndex) GetForwardLinks(pageName string) map[string][]BlockReference {
	return idx.ForwardLinks[pageName]
}

// GetAllPages returns all pages in the index (both sources and targets)
func (idx *BacklinkIndex) GetAllPages() []string {
	pageSet := make(map[string]bool)
	
	// Add all source pages
	for page := range idx.ForwardLinks {
		pageSet[page] = true
	}
	
	// Add all target pages
	for page := range idx.BackwardLinks {
		pageSet[page] = true
	}
	
	pages := make([]string, 0, len(pageSet))
	for page := range pageSet {
		pages = append(pages, page)
	}
	return pages
}

// IsOrphanPage returns true if a page has no incoming or outgoing links
func (idx *BacklinkIndex) IsOrphanPage(pageName string) bool {
	hasOutgoing := len(idx.ForwardLinks[pageName]) > 0
	hasIncoming := len(idx.BackwardLinks[pageName]) > 0
	return !hasOutgoing && !hasIncoming
}

// FindOrphanBlocks returns blocks that have no outgoing references
func FindOrphanBlocks(blocks []*Block) []*Block {
	orphans := make([]*Block, 0)
	
	for _, block := range blocks {
		links := ExtractPageLinks(block.Content)
		if len(links) == 0 && strings.TrimSpace(block.Content) != "" {
			orphans = append(orphans, block)
		}
	}
	
	return orphans
}