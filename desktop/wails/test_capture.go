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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestCapture handles output capture for testing
type TestCapture struct {
	outputDir string
	enabled   bool
}

// NewTestCapture creates a new test capture instance
func NewTestCapture(enabled bool) *TestCapture {
	outputDir := filepath.Join("frontend", "test-output")
	
	if enabled {
		// Create output directories
		dirs := []string{
			filepath.Join(outputDir, "dom"),
			filepath.Join(outputDir, "api"),
			filepath.Join(outputDir, "state"),
			filepath.Join(outputDir, "logs"),
		}
		
		for _, dir := range dirs {
			os.MkdirAll(dir, 0755)
		}
	}
	
	return &TestCapture{
		outputDir: outputDir,
		enabled:   enabled,
	}
}

// CaptureDOM captures the rendered DOM for a page
func (tc *TestCapture) CaptureDOM(pageName string, html string) error {
	if !tc.enabled {
		return nil
	}
	
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("page-%s-%s.html", sanitizeFilename(pageName), timestamp)
	path := filepath.Join(tc.outputDir, "dom", filename)
	
	return os.WriteFile(path, []byte(html), 0644)
}

// LogAPICall logs an API method call with parameters and result
func (tc *TestCapture) LogAPICall(method string, params interface{}, result interface{}) error {
	if !tc.enabled {
		return nil
	}
	
	timestamp := time.Now().Format("20060102-150405")
	
	// Create API call record
	record := map[string]interface{}{
		"timestamp": timestamp,
		"method":    method,
		"params":    params,
		"result":    result,
	}
	
	// Append to daily log file
	dateStr := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("api-calls-%s.json", dateStr)
	path := filepath.Join(tc.outputDir, "api", filename)
	
	return tc.appendJSON(path, record)
}

// ExportState exports application state
func (tc *TestCapture) ExportState(state interface{}) error {
	if !tc.enabled {
		return nil
	}
	
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("app-state-%s.json", timestamp)
	path := filepath.Join(tc.outputDir, "state", filename)
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// LogUserAction logs a user interaction
func (tc *TestCapture) LogUserAction(action string) error {
	if !tc.enabled {
		return nil
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, action)
	
	// Append to daily log
	dateStr := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("user-actions-%s.log", dateStr)
	path := filepath.Join(tc.outputDir, "logs", filename)
	
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	
	_, err = f.WriteString(logEntry)
	return err
}

// CaptureNavigationHistory exports the navigation history
func (tc *TestCapture) CaptureNavigationHistory(history []string) error {
	if !tc.enabled {
		return nil
	}
	
	path := filepath.Join(tc.outputDir, "state", "navigation-history.json")
	
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// Helper to append JSON records to a file
func (tc *TestCapture) appendJSON(path string, record interface{}) error {
	// Read existing records
	var records []interface{}
	
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &records)
	}
	
	// Append new record
	records = append(records, record)
	
	// Write back
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// Helper to sanitize filenames
func sanitizeFilename(name string) string {
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
		" ", "-",
	)
	return replacer.Replace(strings.ToLower(name))
}

// LogResourceError logs failed resource loads
func (tc *TestCapture) LogResourceError(resourceInfo map[string]interface{}) error {
	if !tc.enabled {
		return nil
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Create error record
	record := map[string]interface{}{
		"timestamp": timestamp,
		"error":     resourceInfo,
	}
	
	// Append to daily resource error log
	dateStr := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("resource-errors-%s.json", dateStr)
	path := filepath.Join(tc.outputDir, "logs", filename)
	
	return tc.appendJSON(path, record)
}

// Add test capture instance to App
var testCapture *TestCapture

// InitTestCapture initializes the test capture system
func (a *App) InitTestCapture() {
	testCapture = NewTestCapture(a.TestMode)
	
	if a.TestMode {
		fmt.Println("Test mode enabled - outputs will be captured to frontend/test-output/")
	}
}

// IsTestMode returns whether test mode is enabled (exposed to frontend)
func (a *App) IsTestMode() bool {
	return a.TestMode
}

// CaptureDOM captures DOM content from frontend
func (a *App) CaptureDOM(pageName string, html string) error {
	if a.TestMode && testCapture != nil {
		return testCapture.CaptureDOM(pageName, html)
	}
	return nil
}

// LogUserAction logs a user action from frontend
func (a *App) LogUserAction(action string) error {
	if a.TestMode && testCapture != nil {
		return testCapture.LogUserAction(action)
	}
	return nil
}

// CaptureNavigationHistory captures navigation history from frontend
func (a *App) CaptureNavigationHistory(history []string) error {
	if a.TestMode && testCapture != nil {
		return testCapture.CaptureNavigationHistory(history)
	}
	return nil
}

// LogResourceError logs failed resource loads (images, scripts, etc)
func (a *App) LogResourceError(resourceInfo map[string]interface{}) error {
	if a.TestMode && testCapture != nil {
		return testCapture.LogResourceError(resourceInfo)
	}
	return nil
}