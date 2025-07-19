package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Page templates for variety
var pageTemplates = []string{
	"Daily Note",
	"Project Page",
	"Meeting Notes",
	"Research Topic",
	"Book Notes",
	"Code Documentation",
	"Personal Journal",
	"Task List",
}

// Generate random words for content
var words = []string{
	"system", "design", "architecture", "implementation", "testing",
	"development", "process", "workflow", "documentation", "analysis",
	"research", "project", "meeting", "discussion", "decision",
	"requirement", "specification", "feature", "bug", "fix",
	"improvement", "optimization", "performance", "security", "reliability",
	"scalability", "maintainability", "usability", "accessibility", "compatibility",
}

var tags = []string{
	"#important", "#urgent", "#todo", "#done", "#inprogress",
	"#review", "#question", "#idea", "#note", "#reference",
}

func main() {
	var (
		numPages    = flag.Int("pages", 1000, "Number of pages to generate")
		outputDir   = flag.String("output", "./test-vault", "Output directory")
		linkDensity = flag.Float64("links", 0.1, "Link density (0-1)")
		seed        = flag.Int64("seed", 0, "Random seed (0 for time-based)")
	)
	flag.Parse()

	// Set random seed
	if *seed == 0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(*seed)
	}

	// Create output directory
	pagesDir := filepath.Join(*outputDir, "pages")
	assetsDir := filepath.Join(*outputDir, "assets")
	
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pages directory: %v\n", err)
		os.Exit(1)
	}
	
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating assets directory: %v\n", err)
		os.Exit(1)
	}

	// Generate page names
	pageNames := generatePageNames(*numPages)
	
	// Generate pages
	fmt.Printf("Generating %d pages in %s...\n", *numPages, *outputDir)
	startTime := time.Now()
	
	for i, pageName := range pageNames {
		content := generatePageContent(pageName, pageNames, *linkDensity)
		filePath := filepath.Join(pagesDir, fmt.Sprintf("%s.md", sanitizeFileName(pageName)))
		
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", filePath, err)
			continue
		}
		
		if (i+1) % 100 == 0 {
			fmt.Printf("Generated %d/%d pages...\n", i+1, *numPages)
		}
	}
	
	// Create a simple README
	readme := fmt.Sprintf(`# Test Vault

This is a generated test vault with %d pages for performance testing.

Generated on: %s
Link density: %.2f
Random seed: %d

## Statistics
- Total pages: %d
- Assets directory: created (empty)
- Page types: %v
`, *numPages, time.Now().Format("2006-01-02 15:04:05"), *linkDensity, *seed, *numPages, pageTemplates)

	if err := os.WriteFile(filepath.Join(*outputDir, "README.md"), []byte(readme), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing README: %v\n", err)
	}
	
	elapsed := time.Since(startTime)
	fmt.Printf("Generated %d pages in %v (%.2f pages/second)\n", 
		*numPages, elapsed, float64(*numPages)/elapsed.Seconds())
}

func generatePageNames(count int) []string {
	names := make([]string, 0, count)
	used := make(map[string]bool)
	
	// Add some date pages
	dateCount := count / 10 // 10% date pages
	for i := 0; i < dateCount; i++ {
		date := time.Now().AddDate(0, 0, -i)
		name := date.Format("2006-01-02")
		if !used[name] {
			names = append(names, name)
			used[name] = true
		}
	}
	
	// Generate remaining pages
	for len(names) < count {
		template := pageTemplates[rand.Intn(len(pageTemplates))]
		num := rand.Intn(1000)
		name := fmt.Sprintf("%s %d", template, num)
		
		if !used[name] {
			names = append(names, name)
			used[name] = true
		}
	}
	
	return names
}

func generatePageContent(pageName string, allPages []string, linkDensity float64) string {
	var content strings.Builder
	
	// Title
	content.WriteString(fmt.Sprintf("# %s\n\n", pageName))
	
	// Add some properties for certain page types
	if strings.Contains(pageName, "Project") {
		content.WriteString(fmt.Sprintf("status:: %s\n", randomChoice([]string{"active", "completed", "on-hold", "planning"})))
		content.WriteString(fmt.Sprintf("priority:: %s\n", randomChoice([]string{"high", "medium", "low"})))
		content.WriteString("\n")
	}
	
	// Generate 5-20 blocks
	numBlocks := 5 + rand.Intn(16)
	for i := 0; i < numBlocks; i++ {
		indent := strings.Repeat("\t", rand.Intn(3)) // 0-2 levels of nesting
		
		// Randomly choose block type
		blockType := rand.Float64()
		switch {
		case blockType < 0.2: // 20% TODO items
			status := randomChoice([]string{"TODO", "DONE", "DOING", "LATER", "NOW"})
			content.WriteString(fmt.Sprintf("%s- %s %s\n", indent, status, generateSentence(allPages, linkDensity)))
			
		case blockType < 0.3: // 10% tags
			tag := randomChoice(tags)
			content.WriteString(fmt.Sprintf("%s- %s %s\n", indent, generateSentence(allPages, linkDensity), tag))
			
		case blockType < 0.35: // 5% images
			content.WriteString(fmt.Sprintf("%s- ![diagram](../assets/diagram-%d.png)\n", indent, rand.Intn(100)))
			
		default: // 65% regular text
			content.WriteString(fmt.Sprintf("%s- %s\n", indent, generateSentence(allPages, linkDensity)))
		}
		
		// Add sub-blocks sometimes
		if rand.Float64() < 0.3 && indent == "" {
			subBlocks := 1 + rand.Intn(3)
			for j := 0; j < subBlocks; j++ {
				content.WriteString(fmt.Sprintf("\t- %s\n", generateSentence(allPages, linkDensity)))
			}
		}
	}
	
	return content.String()
}

func generateSentence(allPages []string, linkDensity float64) string {
	length := 5 + rand.Intn(10)
	sentenceWords := make([]string, length)
	
	for i := 0; i < length; i++ {
		if rand.Float64() < linkDensity && len(allPages) > 0 {
			// Insert a page link
			page := allPages[rand.Intn(len(allPages))]
			sentenceWords[i] = fmt.Sprintf("[[%s]]", page)
		} else {
			// Regular word
			sentenceWords[i] = randomChoice(words)
		}
	}
	
	sentence := strings.Join(sentenceWords, " ")
	return strings.ToUpper(sentence[:1]) + sentence[1:] + "."
}

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func sanitizeFileName(name string) string {
	// Replace problematic characters
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return replacer.Replace(name)
}