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
	"testing"
	"time"
)

func TestParseDateTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "ISO format",
			input:    "2025-01-15",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Short month with ordinal",
			input:    "Jan 15th, 2025",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Full month with ordinal",
			input:    "January 15th, 2025",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Short month without ordinal",
			input:    "Jan 15, 2025",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Different ordinals",
			input:    "Jan 1st, 2025",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "2nd ordinal",
			input:    "Jan 2nd, 2025",
			expected: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "3rd ordinal",
			input:    "Jan 3rd, 2025",
			expected: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Slash format",
			input:    "2025/01/15",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Not a date",
			input:    "Page A",
			expected: time.Time{},
			wantErr:  true,
		},
		{
			name:     "Invalid date",
			input:    "Jan 32nd, 2025",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateTitle(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("ParseDateTitle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatDateForPage(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "1st of month",
			date:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "Jan 1st, 2025",
		},
		{
			name:     "2nd of month",
			date:     time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: "Jan 2nd, 2025",
		},
		{
			name:     "3rd of month",
			date:     time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
			expected: "Jan 3rd, 2025",
		},
		{
			name:     "15th of month",
			date:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "Jan 15th, 2025",
		},
		{
			name:     "21st of month",
			date:     time.Date(2025, 1, 21, 0, 0, 0, 0, time.UTC),
			expected: "Jan 21st, 2025",
		},
		{
			name:     "22nd of month",
			date:     time.Date(2025, 1, 22, 0, 0, 0, 0, time.UTC),
			expected: "Jan 22nd, 2025",
		},
		{
			name:     "23rd of month",
			date:     time.Date(2025, 1, 23, 0, 0, 0, 0, time.UTC),
			expected: "Jan 23rd, 2025",
		},
		{
			name:     "31st of month",
			date:     time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: "Jan 31st, 2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateForPage(tt.date); got != tt.expected {
				t.Errorf("FormatDateForPage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExtractDateReferences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // number of dates found
	}{
		{
			name:     "Single date reference",
			input:    "Meeting on [[Jan 15th, 2025]]",
			expected: 1,
		},
		{
			name:     "Multiple date references",
			input:    "From [[Jan 1st, 2025]] to [[Jan 31st, 2025]]",
			expected: 2,
		},
		{
			name:     "Mixed references",
			input:    "See [[Page A]] and meeting on [[Jan 15th, 2025]]",
			expected: 1,
		},
		{
			name:     "No date references",
			input:    "Just some [[Page A]] and [[Page B]]",
			expected: 0,
		},
		{
			name:     "ISO date reference",
			input:    "Deadline: [[2025-01-15]]",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dates := ExtractDateReferences(tt.input)
			if len(dates) != tt.expected {
				t.Errorf("ExtractDateReferences() found %d dates, want %d", len(dates), tt.expected)
			}
		})
	}
}

func TestGetDatePageFilename(t *testing.T) {
	date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	expected := "2025-01-15.md"
	
	if got := GetDatePageFilename(date); got != expected {
		t.Errorf("GetDatePageFilename() = %v, want %v", got, expected)
	}
}

func TestParseFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "Valid date filename with extension",
			filename: "2025-01-15.md",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Valid date filename without extension",
			filename: "2025-01-15",
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Invalid filename",
			filename: "page-a.md",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateFromFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateFromFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("ParseDateFromFilename() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelativeDateString(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "Today",
			date:     today,
			expected: "Today",
		},
		{
			name:     "Tomorrow",
			date:     today.AddDate(0, 0, 1),
			expected: "Tomorrow",
		},
		{
			name:     "Yesterday",
			date:     today.AddDate(0, 0, -1),
			expected: "Yesterday",
		},
		{
			name:     "In 3 days",
			date:     today.AddDate(0, 0, 3),
			expected: "In 3 days",
		},
		{
			name:     "5 days ago",
			date:     today.AddDate(0, 0, -5),
			expected: "5 days ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RelativeDateString(tt.date); got != tt.expected {
				t.Errorf("RelativeDateString() = %v, want %v", got, tt.expected)
			}
		})
	}
}