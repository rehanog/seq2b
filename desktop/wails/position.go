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
	"fmt"
	"github.com/rehanog/seq2b/pkg/parser"
)

// BlockPath represents a path to a block using array indices
// For example, [0, 2, 1] means: 1st top-level block -> 3rd child -> 2nd sub-child
type BlockPath []int

// FindBlockByPath traverses the block tree using the given path
func FindBlockByPath(blocks []*parser.Block, path BlockPath) (*parser.Block, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	
	// Start with top-level blocks
	currentLevel := blocks
	var currentBlock *parser.Block
	
	for i, index := range path {
		if index < 0 || index >= len(currentLevel) {
			return nil, fmt.Errorf("invalid index %d at path position %d (max: %d)", index, i, len(currentLevel)-1)
		}
		
		currentBlock = currentLevel[index]
		
		// If not the last element in path, move to children
		if i < len(path)-1 {
			currentLevel = currentBlock.Children
		}
	}
	
	return currentBlock, nil
}

// GetBlockPath finds the path to a specific block in the tree
func GetBlockPath(blocks []*parser.Block, targetBlock *parser.Block) (BlockPath, bool) {
	var path BlockPath
	
	// Helper function for recursive search
	var search func([]*parser.Block, BlockPath) (BlockPath, bool)
	search = func(blocks []*parser.Block, currentPath BlockPath) (BlockPath, bool) {
		for i, block := range blocks {
			// Create new path with current index
			newPath := append(append(BlockPath{}, currentPath...), i)
			
			if block == targetBlock {
				return newPath, true
			}
			
			// Search children
			if foundPath, found := search(block.Children, newPath); found {
				return foundPath, true
			}
		}
		return nil, false
	}
	
	return search(blocks, path)
}

// InsertBlockAtPath inserts a block at the specified position
// The last element of the path is the index where the block will be inserted
func InsertBlockAtPath(blocks []*parser.Block, newBlock *parser.Block, path BlockPath) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path")
	}
	
	// If path has only one element, insert at top level
	if len(path) == 1 {
		index := path[0]
		if index < 0 || index > len(blocks) {
			return fmt.Errorf("invalid insertion index %d (max: %d)", index, len(blocks))
		}
		// Note: This modifies the slice in place, which won't work for top-level
		// The caller needs to handle top-level insertion specially
		return fmt.Errorf("cannot modify top-level blocks slice in place")
	}
	
	// Find parent using all but last element of path
	parentPath := path[:len(path)-1]
	parent, err := FindBlockByPath(blocks, parentPath)
	if err != nil {
		return fmt.Errorf("parent not found: %w", err)
	}
	
	// Insert into parent's children
	index := path[len(path)-1]
	if index < 0 || index > len(parent.Children) {
		return fmt.Errorf("invalid insertion index %d (max: %d)", index, len(parent.Children))
	}
	
	// Insert the block
	parent.Children = append(parent.Children[:index], 
		append([]*parser.Block{newBlock}, parent.Children[index:]...)...)
	
	// Set parent reference
	newBlock.Parent = parent
	newBlock.Depth = parent.Depth + 1
	
	return nil
}

// RemoveBlockAtPath removes a block at the specified path
func RemoveBlockAtPath(blocks []*parser.Block, path BlockPath) (*parser.Block, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	
	// If path has only one element, remove from top level
	if len(path) == 1 {
		index := path[0]
		if index < 0 || index >= len(blocks) {
			return nil, fmt.Errorf("invalid index %d (max: %d)", index, len(blocks)-1)
		}
		// Note: This needs special handling at the caller level
		return nil, fmt.Errorf("cannot modify top-level blocks slice in place")
	}
	
	// Find parent using all but last element of path
	parentPath := path[:len(path)-1]
	parent, err := FindBlockByPath(blocks, parentPath)
	if err != nil {
		return nil, fmt.Errorf("parent not found: %w", err)
	}
	
	// Remove from parent's children
	index := path[len(path)-1]
	if index < 0 || index >= len(parent.Children) {
		return nil, fmt.Errorf("invalid index %d (max: %d)", index, len(parent.Children)-1)
	}
	
	removed := parent.Children[index]
	parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
	
	return removed, nil
}

// ShiftedPaths calculates which block paths change after an insertion or deletion
type PathShift struct {
	OldPath BlockPath
	NewPath BlockPath
}

// CalculatePathShifts determines how paths change after inserting at insertPath
func CalculatePathShiftsAfterInsert(blocks []*parser.Block, insertPath BlockPath) []PathShift {
	var shifts []PathShift
	
	// Helper to check if a path is affected by the insertion
	isAffected := func(path BlockPath) bool {
		// Check each level of the path
		for i := 0; i < len(path) && i < len(insertPath); i++ {
			if i < len(insertPath)-1 {
				// For parent levels, must match exactly
				if path[i] != insertPath[i] {
					return false
				}
			} else {
				// At insertion level, affected if index >= insertion index
				return path[i] >= insertPath[i]
			}
		}
		return false
	}
	
	// Traverse all blocks and find affected paths
	var traverse func([]*parser.Block, BlockPath)
	traverse = func(blocks []*parser.Block, currentPath BlockPath) {
		for i, block := range blocks {
			path := append(append(BlockPath{}, currentPath...), i)
			
			if isAffected(path) {
				// Calculate new path
				newPath := append(BlockPath{}, path...)
				// Increment the index at the insertion level
				if len(path) == len(insertPath) {
					newPath[len(newPath)-1]++
				}
				
				shifts = append(shifts, PathShift{
					OldPath: path,
					NewPath: newPath,
				})
			}
			
			// Traverse children
			traverse(block.Children, path)
		}
	}
	
	traverse(blocks, BlockPath{})
	return shifts
}