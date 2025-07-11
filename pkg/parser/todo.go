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
	"regexp"
	"strings"
)

// TodoState represents the state of a TODO item
type TodoState string

const (
	TodoStateNone     TodoState = ""
	TodoStateTodo     TodoState = "TODO"
	TodoStateDoing    TodoState = "DOING"
	TodoStateDone     TodoState = "DONE"
	TodoStateWaiting  TodoState = "WAITING"
	TodoStateCanceled TodoState = "CANCELED"
	TodoStateLater    TodoState = "LATER"
	TodoStateNow      TodoState = "NOW"
)

// CheckboxState represents the state of a checkbox
type CheckboxState string

const (
	CheckboxNone    CheckboxState = ""
	CheckboxUnchecked CheckboxState = "[ ]"
	CheckboxChecked   CheckboxState = "[x]"
	CheckboxPartial   CheckboxState = "[-]"
)

// TodoInfo contains TODO-related information for a block
type TodoInfo struct {
	TodoState     TodoState
	CheckboxState CheckboxState
	Priority      string // A, B, C, etc. from TODO [#A]
}

var (
	// Regex to match TODO states at the beginning of a block
	// Matches: TODO, DOING, DONE, WAITING, CANCELED, LATER, NOW
	todoStateRegex = regexp.MustCompile(`^(TODO|DOING|DONE|WAITING|CANCELED|LATER|NOW)\s+`)
	
	// Regex to match TODO with priority
	// Matches: TODO [#A], DONE [#B], etc.
	todoPriorityRegex = regexp.MustCompile(`^(TODO|DOING|DONE|WAITING|CANCELED|LATER|NOW)\s+\[#([A-Z])\]\s+`)
	
	// Regex to match checkboxes
	// Matches: [ ], [x], [X], [-]
	checkboxRegex = regexp.MustCompile(`^\[([ xX\-])\]\s+`)
)

// ParseTodoInfo extracts TODO information from block content
func ParseTodoInfo(content string) TodoInfo {
	info := TodoInfo{}
	trimmed := strings.TrimSpace(content)
	
	// Check for TODO with priority first
	if matches := todoPriorityRegex.FindStringSubmatch(trimmed); len(matches) > 0 {
		info.TodoState = TodoState(matches[1])
		info.Priority = matches[2]
		return info
	}
	
	// Check for TODO state without priority
	if matches := todoStateRegex.FindStringSubmatch(trimmed); len(matches) > 0 {
		info.TodoState = TodoState(matches[1])
		return info
	}
	
	// Check for checkbox
	if matches := checkboxRegex.FindStringSubmatch(trimmed); len(matches) > 0 {
		switch strings.ToLower(matches[1]) {
		case " ":
			info.CheckboxState = CheckboxUnchecked
		case "x":
			info.CheckboxState = CheckboxChecked
		case "-":
			info.CheckboxState = CheckboxPartial
		}
	}
	
	return info
}

// RemoveTodoPrefix removes TODO state and checkbox prefixes from content
func RemoveTodoPrefix(content string) string {
	trimmed := strings.TrimSpace(content)
	
	// Remove TODO with priority
	if todoPriorityRegex.MatchString(trimmed) {
		trimmed = todoPriorityRegex.ReplaceAllString(trimmed, "")
	} else if todoStateRegex.MatchString(trimmed) {
		// Remove TODO without priority
		trimmed = todoStateRegex.ReplaceAllString(trimmed, "")
	}
	
	// Remove checkbox
	if checkboxRegex.MatchString(trimmed) {
		trimmed = checkboxRegex.ReplaceAllString(trimmed, "")
	}
	
	return trimmed
}

// GetTodoBlocks returns all blocks with TODO states or checkboxes
func GetTodoBlocks(blocks []*Block) []*Block {
	var todoBlocks []*Block
	
	for _, block := range blocks {
		if block.TodoInfo.TodoState != TodoStateNone || block.TodoInfo.CheckboxState != CheckboxNone {
			todoBlocks = append(todoBlocks, block)
		}
		// Recursively check children
		todoBlocks = append(todoBlocks, GetTodoBlocks(block.Children)...)
	}
	
	return todoBlocks
}

// FilterBlocksByTodoState returns blocks matching the given TODO state
func FilterBlocksByTodoState(blocks []*Block, state TodoState) []*Block {
	var filtered []*Block
	
	for _, block := range blocks {
		if block.TodoInfo.TodoState == state {
			filtered = append(filtered, block)
		}
		// Recursively check children
		filtered = append(filtered, FilterBlocksByTodoState(block.Children, state)...)
	}
	
	return filtered
}