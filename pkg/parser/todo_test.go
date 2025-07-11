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
)

func TestParseTodoInfo(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantTodo TodoState
		wantCheck CheckboxState
		wantPriority string
	}{
		{
			name:     "TODO state",
			content:  "TODO Write tests",
			wantTodo: TodoStateTodo,
		},
		{
			name:     "DONE state",
			content:  "DONE Implement feature",
			wantTodo: TodoStateDone,
		},
		{
			name:     "DOING state",
			content:  "DOING Working on this",
			wantTodo: TodoStateDoing,
		},
		{
			name:     "WAITING state",
			content:  "WAITING For review",
			wantTodo: TodoStateWaiting,
		},
		{
			name:     "TODO with priority A",
			content:  "TODO [#A] Important task",
			wantTodo: TodoStateTodo,
			wantPriority: "A",
		},
		{
			name:     "DONE with priority B",
			content:  "DONE [#B] Less important task",
			wantTodo: TodoStateDone,
			wantPriority: "B",
		},
		{
			name:      "Unchecked checkbox",
			content:   "[ ] Buy milk",
			wantCheck: CheckboxUnchecked,
		},
		{
			name:      "Checked checkbox",
			content:   "[x] Buy bread",
			wantCheck: CheckboxChecked,
		},
		{
			name:      "Checked checkbox uppercase",
			content:   "[X] Buy eggs",
			wantCheck: CheckboxChecked,
		},
		{
			name:      "Partial checkbox",
			content:   "[-] Shopping list",
			wantCheck: CheckboxPartial,
		},
		{
			name:     "Regular text",
			content:  "Just some regular text",
			wantTodo: TodoStateNone,
			wantCheck: CheckboxNone,
		},
		{
			name:     "TODO in middle of text",
			content:  "Need to TODO this later",
			wantTodo: TodoStateNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := ParseTodoInfo(tt.content)
			
			if info.TodoState != tt.wantTodo {
				t.Errorf("ParseTodoInfo() TodoState = %v, want %v", info.TodoState, tt.wantTodo)
			}
			
			if info.CheckboxState != tt.wantCheck {
				t.Errorf("ParseTodoInfo() CheckboxState = %v, want %v", info.CheckboxState, tt.wantCheck)
			}
			
			if info.Priority != tt.wantPriority {
				t.Errorf("ParseTodoInfo() Priority = %v, want %v", info.Priority, tt.wantPriority)
			}
		})
	}
}

func TestRemoveTodoPrefix(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Remove TODO",
			content: "TODO Write tests",
			want:    "Write tests",
		},
		{
			name:    "Remove TODO with priority",
			content: "TODO [#A] Important task",
			want:    "Important task",
		},
		{
			name:    "Remove checkbox",
			content: "[ ] Buy milk",
			want:    "Buy milk",
		},
		{
			name:    "Remove checked checkbox",
			content: "[x] Buy bread",
			want:    "Buy bread",
		},
		{
			name:    "No prefix to remove",
			content: "Regular text",
			want:    "Regular text",
		},
		{
			name:    "TODO in middle stays",
			content: "Need to TODO this",
			want:    "Need to TODO this",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveTodoPrefix(tt.content)
			if got != tt.want {
				t.Errorf("RemoveTodoPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoBlockParsing(t *testing.T) {
	content := `# My Tasks

- TODO [#A] High priority task
- TODO Normal task
- DONE Completed task
- [ ] Unchecked item
- [x] Checked item
- Regular block without TODO`

	result, err := ParseFile(content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Find TODO blocks
	todoBlocks := GetTodoBlocks(result.Page.Blocks)
	if len(todoBlocks) != 5 { // Should find 5 blocks with TODO/checkbox
		t.Errorf("Expected 5 TODO blocks, got %d", len(todoBlocks))
	}

	// Test filtering by state
	todoOnly := FilterBlocksByTodoState(result.Page.Blocks, TodoStateTodo)
	if len(todoOnly) != 2 {
		t.Errorf("Expected 2 TODO blocks, got %d", len(todoOnly))
	}

	doneOnly := FilterBlocksByTodoState(result.Page.Blocks, TodoStateDone)
	if len(doneOnly) != 1 {
		t.Errorf("Expected 1 DONE block, got %d", len(doneOnly))
	}
}

func TestNestedTodoBlocks(t *testing.T) {
	content := `- TODO Parent task
  - [ ] Subtask 1
  - [x] Subtask 2
  - TODO Nested TODO
    - [ ] Sub-subtask`

	result, err := ParseFile(content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Should find all TODO blocks including nested ones
	todoBlocks := GetTodoBlocks(result.Page.Blocks)
	if len(todoBlocks) != 5 {
		t.Errorf("Expected 5 TODO blocks (including nested), got %d", len(todoBlocks))
	}

	// Check that parent block has TODO state
	if len(result.Page.Blocks) > 0 {
		parentBlock := result.Page.Blocks[0]
		if parentBlock.TodoInfo.TodoState != TodoStateTodo {
			t.Errorf("Parent block should have TODO state")
		}
		
		// Check first child has checkbox
		if len(parentBlock.Children) > 0 {
			firstChild := parentBlock.Children[0]
			if firstChild.TodoInfo.CheckboxState != CheckboxUnchecked {
				t.Errorf("First child should have unchecked checkbox")
			}
		}
	}
}