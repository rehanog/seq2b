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
	"reflect"
	"testing"
	"github.com/rehanog/seq2b/pkg/parser"
)

func TestFindBlockByPath(t *testing.T) {
	// Create test block structure
	block1 := &parser.Block{Content: "Block 1"}
	block2 := &parser.Block{Content: "Block 2"}
	block1Child1 := &parser.Block{Content: "Block 1 Child 1", Parent: block1, Depth: 1}
	block1Child2 := &parser.Block{Content: "Block 1 Child 2", Parent: block1, Depth: 1}
	block1Child1SubChild := &parser.Block{Content: "Block 1 Child 1 SubChild", Parent: block1Child1, Depth: 2}
	
	block1.Children = []*parser.Block{block1Child1, block1Child2}
	block1Child1.Children = []*parser.Block{block1Child1SubChild}
	
	blocks := []*parser.Block{block1, block2}
	
	tests := []struct {
		name    string
		path    BlockPath
		want    *parser.Block
		wantErr bool
	}{
		{
			name:    "Find top-level block",
			path:    BlockPath{0},
			want:    block1,
			wantErr: false,
		},
		{
			name:    "Find nested block",
			path:    BlockPath{0, 1},
			want:    block1Child2,
			wantErr: false,
		},
		{
			name:    "Find deeply nested block",
			path:    BlockPath{0, 0, 0},
			want:    block1Child1SubChild,
			wantErr: false,
		},
		{
			name:    "Invalid index",
			path:    BlockPath{2},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty path",
			path:    BlockPath{},
			want:    nil,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindBlockByPath(blocks, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBlockByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindBlockByPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBlockPath(t *testing.T) {
	// Create test block structure
	block1 := &parser.Block{Content: "Block 1"}
	block2 := &parser.Block{Content: "Block 2"}
	block1Child1 := &parser.Block{Content: "Block 1 Child 1", Parent: block1}
	block1Child2 := &parser.Block{Content: "Block 1 Child 2", Parent: block1}
	block1Child1SubChild := &parser.Block{Content: "Block 1 Child 1 SubChild", Parent: block1Child1}
	
	block1.Children = []*parser.Block{block1Child1, block1Child2}
	block1Child1.Children = []*parser.Block{block1Child1SubChild}
	
	blocks := []*parser.Block{block1, block2}
	
	// Create a block not in the tree
	orphanBlock := &parser.Block{Content: "Orphan"}
	
	tests := []struct {
		name       string
		target     *parser.Block
		wantPath   BlockPath
		wantFound  bool
	}{
		{
			name:      "Find top-level block",
			target:    block1,
			wantPath:  BlockPath{0},
			wantFound: true,
		},
		{
			name:      "Find nested block",
			target:    block1Child2,
			wantPath:  BlockPath{0, 1},
			wantFound: true,
		},
		{
			name:      "Find deeply nested block",
			target:    block1Child1SubChild,
			wantPath:  BlockPath{0, 0, 0},
			wantFound: true,
		},
		{
			name:      "Block not in tree",
			target:    orphanBlock,
			wantPath:  nil,
			wantFound: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotFound := GetBlockPath(blocks, tt.target)
			if gotFound != tt.wantFound {
				t.Errorf("GetBlockPath() found = %v, want %v", gotFound, tt.wantFound)
				return
			}
			if !reflect.DeepEqual(gotPath, tt.wantPath) {
				t.Errorf("GetBlockPath() path = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestCalculatePathShiftsAfterInsert(t *testing.T) {
	// Create test block structure
	block1 := &parser.Block{Content: "Block 1"}
	block2 := &parser.Block{Content: "Block 2"}
	block3 := &parser.Block{Content: "Block 3"}
	
	blocks := []*parser.Block{block1, block2, block3}
	
	tests := []struct {
		name       string
		insertPath BlockPath
		wantShifts []PathShift
	}{
		{
			name:       "Insert at beginning",
			insertPath: BlockPath{0},
			wantShifts: []PathShift{
				{OldPath: BlockPath{0}, NewPath: BlockPath{1}},
				{OldPath: BlockPath{1}, NewPath: BlockPath{2}},
				{OldPath: BlockPath{2}, NewPath: BlockPath{3}},
			},
		},
		{
			name:       "Insert in middle",
			insertPath: BlockPath{1},
			wantShifts: []PathShift{
				{OldPath: BlockPath{1}, NewPath: BlockPath{2}},
				{OldPath: BlockPath{2}, NewPath: BlockPath{3}},
			},
		},
		{
			name:       "Insert at end",
			insertPath: BlockPath{3},
			wantShifts: []PathShift{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShifts := CalculatePathShiftsAfterInsert(blocks, tt.insertPath)
			if !reflect.DeepEqual(gotShifts, tt.wantShifts) {
				t.Errorf("CalculatePathShiftsAfterInsert() = %v, want %v", gotShifts, tt.wantShifts)
			}
		})
	}
}