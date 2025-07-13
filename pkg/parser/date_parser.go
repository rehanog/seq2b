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
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Common date formats that Logseq supports
var dateFormats = []string{
	"2006-01-02",           // ISO format: 2025-01-15
	"Jan 2, 2006",          // Jan 15, 2025
	"January 2, 2006",      // January 15, 2025
	"Jan 2nd, 2006",        // Jan 15th, 2025 (with ordinal)
	"January 2nd, 2006",    // January 15th, 2025 (with ordinal)
	"2006/01/02",           // 2025/01/15
	"02-01-2006",           // 15-01-2025 (DD-MM-YYYY)
	"02/01/2006",           // 15/01/2025 (DD/MM/YYYY)
}

// Regular expressions for date patterns
var (
	// Matches ordinal suffixes (1st, 2nd, 3rd, 4th, etc.)
	ordinalPattern = regexp.MustCompile(`(\d+)(st|nd|rd|th)`)
	
	// Matches [[date]] references that might be dates
	dateRefPattern = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
)

// IsDatePage checks if a page title represents a date
func IsDatePage(title string) bool {
	_, err := ParseDateTitle(title)
	return err == nil
}

// ParseDateTitle attempts to parse a page title as a date
func ParseDateTitle(title string) (time.Time, error) {
	// First, try to parse as-is
	if date, err := tryParseDate(title); err == nil {
		return date, nil
	}
	
	// If that fails, try removing ordinal suffixes
	cleaned := removeOrdinalSuffixes(title)
	if cleaned != title {
		if date, err := tryParseDate(cleaned); err == nil {
			return date, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("not a valid date: %s", title)
}

// tryParseDate attempts to parse a string using various date formats
func tryParseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	
	for _, format := range dateFormats {
		if date, err := time.Parse(format, s); err == nil {
			return date, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("no matching date format")
}

// removeOrdinalSuffixes removes ordinal suffixes from dates (1st -> 1, 2nd -> 2, etc.)
func removeOrdinalSuffixes(s string) string {
	return ordinalPattern.ReplaceAllString(s, "$1")
}

// FormatDateForPage formats a date for use as a page title
// Uses the display format: Jan 15th, 2025
func FormatDateForPage(date time.Time) string {
	day := date.Day()
	suffix := getOrdinalSuffix(day)
	return fmt.Sprintf("%s %d%s, %d", date.Format("Jan"), day, suffix, date.Year())
}

// FormatDateISO formats a date in ISO format (YYYY-MM-DD)
func FormatDateISO(date time.Time) string {
	return date.Format("2006-01-02")
}

// getOrdinalSuffix returns the ordinal suffix for a day (1st, 2nd, 3rd, 4th, etc.)
func getOrdinalSuffix(day int) string {
	switch day {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}

// ExtractDateReferences finds all [[date]] references in text
func ExtractDateReferences(text string) []time.Time {
	matches := dateRefPattern.FindAllStringSubmatch(text, -1)
	dates := []time.Time{}
	
	for _, match := range matches {
		if len(match) > 1 {
			if date, err := ParseDateTitle(match[1]); err == nil {
				dates = append(dates, date)
			}
		}
	}
	
	return dates
}

// GetTodayPageTitle returns today's date formatted as a page title
func GetTodayPageTitle() string {
	return FormatDateForPage(time.Now())
}

// GetDatePageFilename returns the filename for a date page
// Uses ISO format for consistency: 2025-01-15.md
func GetDatePageFilename(date time.Time) string {
	return FormatDateISO(date) + ".md"
}

// ParseDateFromFilename extracts date from a filename like 2025-01-15.md
func ParseDateFromFilename(filename string) (time.Time, error) {
	// Remove .md extension if present
	name := strings.TrimSuffix(filename, ".md")
	
	// Try to parse as ISO date
	return time.Parse("2006-01-02", name)
}

// IsWithinDateRange checks if a date is within a range
func IsWithinDateRange(date, start, end time.Time) bool {
	return !date.Before(start) && !date.After(end)
}

// GetWeekNumber returns the ISO week number for a date
func GetWeekNumber(date time.Time) int {
	_, week := date.ISOWeek()
	return week
}

// RelativeDateString returns a human-readable relative date string
func RelativeDateString(date time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	target := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	
	days := int(target.Sub(today).Hours() / 24)
	
	switch days {
	case 0:
		return "Today"
	case 1:
		return "Tomorrow"
	case -1:
		return "Yesterday"
	default:
		if days > 0 && days <= 7 {
			return fmt.Sprintf("In %d days", days)
		} else if days < 0 && days >= -7 {
			return fmt.Sprintf("%d days ago", -days)
		}
		return FormatDateForPage(date)
	}
}