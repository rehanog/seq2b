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

// +build ignore

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Test result types
type TestResult struct {
	Name     string
	Passed   bool
	Messages []string
}

type APICall struct {
	Timestamp string                 `json:"timestamp"`
	Method    string                 `json:"method"`
	Params    map[string]interface{} `json:"params"`
	Result    map[string]interface{} `json:"result"`
}

// Main verification function
func main() {
	fmt.Println("=== Seq2b Test Output Verification ===")
	fmt.Println()

	outputDir := filepath.Join("frontend", "test-output")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Println("Error: test-output directory not found. Please run tests with -test-mode flag first.")
		os.Exit(1)
	}

	results := []TestResult{}

	// Verify navigation test
	results = append(results, verifyNavigationTest(outputDir))

	// Verify edit persistence test
	results = append(results, verifyEditPersistenceTest(outputDir))

	// Verify backlinks test
	results = append(results, verifyBacklinksTest(outputDir))
	
	// Check for resource loading errors
	results = append(results, verifyResourceLoading(outputDir))

	// Print results
	fmt.Println("\n=== Test Results ===")
	passed := 0
	failed := 0
	for _, result := range results {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
			failed++
		} else {
			passed++
		}
		
		fmt.Printf("\n[%s] %s\n", status, result.Name)
		for _, msg := range result.Messages {
			fmt.Printf("  %s\n", msg)
		}
	}

	fmt.Printf("\n=== Summary: %d passed, %d failed ===\n", passed, failed)
	
	if failed > 0 {
		os.Exit(1)
	}
}

// Verify navigation with case sensitivity
func verifyNavigationTest(outputDir string) TestResult {
	result := TestResult{
		Name:     "Page Navigation with Case Sensitivity",
		Passed:   true,
		Messages: []string{},
	}

	// Check user actions log
	actionLog := filepath.Join(outputDir, "logs", fmt.Sprintf("user-actions-%s.log", time.Now().Format("2006-01-02")))
	if actions, err := readUserActions(actionLog); err == nil {
		foundNavigation := false
		for _, action := range actions {
			if strings.Contains(action, "Navigate to page: page a") ||
			   strings.Contains(action, "Navigate to page: Page a") ||
			   strings.Contains(action, "Navigate to page: Page A") {
				foundNavigation = true
				result.Messages = append(result.Messages, "✓ Found navigation action to page a")
				break
			}
		}
		if !foundNavigation {
			result.Passed = false
			result.Messages = append(result.Messages, "✗ No navigation to 'page a' found in user actions")
		}
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("⚠ Could not read user actions log: %v", err))
	}

	// Check API calls
	apiLog := filepath.Join(outputDir, "api", fmt.Sprintf("api-calls-%s.json", time.Now().Format("2006-01-02")))
	if calls, err := readAPICalls(apiLog); err == nil {
		foundGetPage := false
		for _, call := range calls {
			if call.Method == "GetPage" {
				if params, ok := call.Params["pageName"].(string); ok {
					if strings.EqualFold(params, "page a") {
						foundGetPage = true
						result.Messages = append(result.Messages, fmt.Sprintf("✓ Found GetPage API call for '%s'", params))
						
						// Check if result has no error
						if call.Result != nil {
							if errVal, hasErr := call.Result["error"]; hasErr && errVal == nil {
								result.Messages = append(result.Messages, "✓ GetPage call succeeded (no error)")
							} else if hasErr && errVal != nil {
								result.Passed = false
								result.Messages = append(result.Messages, fmt.Sprintf("✗ GetPage returned error: %v", errVal))
							}
						}
						break
					}
				}
			}
		}
		if !foundGetPage {
			result.Passed = false
			result.Messages = append(result.Messages, "✗ No GetPage API call found for 'page a'")
		}
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("⚠ Could not read API calls log: %v", err))
	}

	// Check DOM capture
	domFiles := findDOMFiles(outputDir, "page-a", "page a", "Page a", "Page A")
	if len(domFiles) > 0 {
		result.Messages = append(result.Messages, fmt.Sprintf("✓ Found %d DOM capture(s) for page a", len(domFiles)))
		
		// Verify DOM contains expected content
		if content, err := os.ReadFile(domFiles[0]); err == nil {
			htmlContent := string(content)
			if strings.Contains(htmlContent, "Page A") || strings.Contains(htmlContent, "Page a") || strings.Contains(htmlContent, "page a") {
				result.Messages = append(result.Messages, "✓ DOM contains page title")
			}
		}
	} else {
		result.Messages = append(result.Messages, "⚠ No DOM captures found for page a")
	}

	// Check navigation history
	historyFile := filepath.Join(outputDir, "state", "navigation-history.json")
	if history, err := readNavigationHistory(historyFile); err == nil {
		for _, page := range history {
			if strings.EqualFold(page, "page a") {
				result.Messages = append(result.Messages, "✓ Navigation history includes page a")
				break
			}
		}
	}

	return result
}

// Verify edit persistence test
func verifyEditPersistenceTest(outputDir string) TestResult {
	result := TestResult{
		Name:     "Edit Persistence",
		Passed:   true,
		Messages: []string{},
	}

	// Check for edit action before navigation
	actionLog := filepath.Join(outputDir, "logs", fmt.Sprintf("user-actions-%s.log", time.Now().Format("2006-01-02")))
	if actions, err := readUserActions(actionLog); err == nil {
		foundEdit := false
		foundNavAfterEdit := false
		editIndex := -1
		
		for i, action := range actions {
			if strings.Contains(action, "Edit block") && strings.Contains(action, "[[new page]]") {
				foundEdit = true
				editIndex = i
				result.Messages = append(result.Messages, "✓ Found edit action with [[new page]] link")
			} else if foundEdit && i > editIndex && strings.Contains(action, "Navigate to page: new page") {
				foundNavAfterEdit = true
				result.Messages = append(result.Messages, "✓ Found navigation to 'new page' after edit")
			}
		}
		
		if !foundEdit {
			result.Messages = append(result.Messages, "⚠ No edit action with [[new page]] found")
		}
		if foundEdit && !foundNavAfterEdit {
			result.Messages = append(result.Messages, "⚠ No navigation to 'new page' found after edit")
		}
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("⚠ Could not read user actions log: %v", err))
	}

	// Check for UpdateBlockAtPath API call
	apiLog := filepath.Join(outputDir, "api", fmt.Sprintf("api-calls-%s.json", time.Now().Format("2006-01-02")))
	if calls, err := readAPICalls(apiLog); err == nil {
		foundUpdate := false
		for _, call := range calls {
			if call.Method == "UpdateBlockAtPath" || call.Method == "UpdateBlock" {
				foundUpdate = true
				result.Messages = append(result.Messages, fmt.Sprintf("✓ Found %s API call", call.Method))
				
				// Check if the update contains the expected content
				if params, ok := call.Params["newContent"].(string); ok {
					if strings.Contains(params, "[[new page]]") {
						result.Messages = append(result.Messages, "✓ Update contains [[new page]] link")
					}
				}
				break
			}
		}
		if !foundUpdate {
			result.Passed = false
			result.Messages = append(result.Messages, "✗ No UpdateBlockAtPath or UpdateBlock API call found")
		}
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("⚠ Could not read API calls log: %v", err))
	}

	return result
}

// Verify backlinks test
func verifyBacklinksTest(outputDir string) TestResult {
	result := TestResult{
		Name:     "Backlink Updates",
		Passed:   true,
		Messages: []string{},
	}

	// Look for DOM captures of Page B
	domFiles := findDOMFiles(outputDir, "page-b", "Page B", "Page b")
	if len(domFiles) > 0 {
		result.Messages = append(result.Messages, fmt.Sprintf("✓ Found %d DOM capture(s) for Page B", len(domFiles)))
		
		// Check if DOM contains backlinks section
		backlinkFound := false
		for _, file := range domFiles {
			if content, err := os.ReadFile(file); err == nil {
				htmlContent := string(content)
				if strings.Contains(htmlContent, "backlink") || strings.Contains(htmlContent, "Backlink") {
					backlinkFound = true
					result.Messages = append(result.Messages, "✓ DOM contains backlinks section")
					
					// Try to find specific backlink references
					if strings.Contains(htmlContent, "backlink-item") || strings.Contains(htmlContent, "backlink-source") {
						result.Messages = append(result.Messages, "✓ Backlinks are rendered in DOM")
					}
					break
				}
			}
		}
		
		if !backlinkFound {
			result.Messages = append(result.Messages, "⚠ No backlinks section found in DOM captures")
		}
	} else {
		result.Messages = append(result.Messages, "⚠ No DOM captures found for Page B")
	}

	// Check API calls for GetPage on Page B
	apiLog := filepath.Join(outputDir, "api", fmt.Sprintf("api-calls-%s.json", time.Now().Format("2006-01-02")))
	if calls, err := readAPICalls(apiLog); err == nil {
		for _, call := range calls {
			if call.Method == "GetPage" {
				if params, ok := call.Params["pageName"].(string); ok {
					if strings.EqualFold(params, "page b") {
						// Check if result contains backlinks
						if data, ok := call.Result["data"].(map[string]interface{}); ok {
							if backlinks, ok := data["backlinks"].([]interface{}); ok && len(backlinks) > 0 {
								result.Messages = append(result.Messages, fmt.Sprintf("✓ GetPage result contains %d backlink(s)", len(backlinks)))
							} else {
								result.Messages = append(result.Messages, "⚠ GetPage result has no backlinks")
							}
						}
						break
					}
				}
			}
		}
	}

	return result
}

// Helper functions

func readUserActions(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var actions []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		actions = append(actions, scanner.Text())
	}
	return actions, scanner.Err()
}

func readAPICalls(filename string) ([]APICall, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var calls []APICall
	err = json.Unmarshal(data, &calls)
	return calls, err
}

func readNavigationHistory(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var history []string
	err = json.Unmarshal(data, &history)
	return history, err
}

func findDOMFiles(outputDir string, patterns ...string) []string {
	domDir := filepath.Join(outputDir, "dom")
	files, err := os.ReadDir(domDir)
	if err != nil {
		return []string{}
	}

	var matches []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		filename := file.Name()
		for _, pattern := range patterns {
			// Convert pattern to lowercase and replace spaces with hyphens
			normalizedPattern := strings.ToLower(strings.ReplaceAll(pattern, " ", "-"))
			if strings.Contains(filename, normalizedPattern) {
				matches = append(matches, filepath.Join(domDir, filename))
				break
			}
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(matches, func(i, j int) bool {
		fi, _ := os.Stat(matches[i])
		fj, _ := os.Stat(matches[j])
		return fi.ModTime().After(fj.ModTime())
	})

	return matches
}

// Verify resource loading
func verifyResourceLoading(outputDir string) TestResult {
	result := TestResult{
		Name:     "Resource Loading Validation",
		Passed:   true,
		Messages: []string{},
	}

	// Check for resource error logs
	resourceErrorLog := filepath.Join(outputDir, "logs", fmt.Sprintf("resource-errors-%s.json", time.Now().Format("2006-01-02")))
	if data, err := os.ReadFile(resourceErrorLog); err == nil {
		var errors []interface{}
		if err := json.Unmarshal(data, &errors); err == nil {
			if len(errors) > 0 {
				result.Passed = false
				result.Messages = append(result.Messages, fmt.Sprintf("❌ Found %d resource loading error(s)", len(errors)))
				
				// Analyze the errors
				for _, errItem := range errors {
					if errMap, ok := errItem.(map[string]interface{}); ok {
						if errData, ok := errMap["error"].(map[string]interface{}); ok {
							errType := errData["type"]
							src := errData["originalSrc"]
							message := errData["message"]
							result.Messages = append(result.Messages, 
								fmt.Sprintf("   - %s: %s (%s)", errType, src, message))
						}
					}
				}
			} else {
				result.Messages = append(result.Messages, "✓ No resource loading errors found")
			}
		}
	} else {
		result.Messages = append(result.Messages, "✓ No resource error log found (good if no errors expected)")
	}

	return result
}