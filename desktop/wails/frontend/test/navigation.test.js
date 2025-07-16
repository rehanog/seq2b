import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import '../src/main.js';

describe('Navigation', () => {
  beforeEach(() => {
    // Reset mocks
    vi.clearAllMocks();
    
    // Reset DOM
    document.getElementById('blocks').innerHTML = '';
    document.getElementById('backlinksList').innerHTML = '';
    document.getElementById('loading').style.display = 'none';
    
    // Trigger DOMContentLoaded to set up event listeners
    document.dispatchEvent(new Event('DOMContentLoaded'));
  });

  afterEach(() => {
    // Clean up any active timers
    vi.clearAllTimers();
  });

  describe('navigateToPage', () => {
    it('should save active edits before navigating', async () => {
      // Mock GetPage responses
      window.go.main.App.GetPage
        .mockResolvedValueOnce({
          title: 'Page A',
          blocks: [{
            id: 'block-1',
            content: 'Original content',
            segments: [{type: 'text', content: 'Original content'}],
            depth: 0,
            children: []
          }],
          backlinks: []
        })
        .mockResolvedValueOnce({
          title: 'Page B',
          blocks: [{
            id: 'block-2',
            content: 'Page B content',
            segments: [{type: 'text', content: 'Page B content'}],
            depth: 0,
            children: []
          }],
          backlinks: []
        });

      // Load initial page
      await window.navigateToPage('Page A');
      
      // Find the block and simulate editing
      const blockText = document.querySelector('.block-text');
      expect(blockText).toBeTruthy();
      
      // Simulate editing mode
      blockText.contentEditable = 'true';
      blockText.classList.add('editing');
      blockText.textContent = 'Edited content with [[Page B]]';
      
      // Mock the save operation
      window.go.main.App.UpdateBlockAtPath.mockResolvedValueOnce({
        action: 'update',
        path: [0],
        block: {
          content: 'Edited content with [[Page B]]',
          segments: [
            {type: 'text', content: 'Edited content with '},
            {type: 'link', content: 'Page B', target: 'Page B'}
          ]
        }
      });

      // Spy on blur to ensure it's called
      const blurSpy = vi.spyOn(blockText, 'blur');
      
      // Navigate to Page B
      await window.navigateToPage('Page B');
      
      // Verify blur was called
      expect(blurSpy).toHaveBeenCalled();
      
      // Wait for blur handler to complete
      await new Promise(resolve => setTimeout(resolve, 150));
      
      // Verify save was called
      expect(window.go.main.App.UpdateBlockAtPath).toHaveBeenCalled();
    });

    it('should handle navigation when no block is being edited', async () => {
      window.go.main.App.GetPage.mockResolvedValueOnce({
        title: 'Test Page',
        blocks: [{
          id: 'block-1',
          content: 'Test content',
          segments: [{type: 'text', content: 'Test content'}],
          depth: 0,
          children: []
        }],
        backlinks: []
      });

      await window.navigateToPage('Test Page');
      
      // Verify GetPage was called
      expect(window.go.main.App.GetPage).toHaveBeenCalledWith('Test Page');
      
      // Verify no save operations were triggered
      expect(window.go.main.App.UpdateBlockAtPath).not.toHaveBeenCalled();
    });
  });

  describe('Block editing', () => {
    it('should save block content on blur', async () => {
      // Set up initial page
      window.go.main.App.GetPage.mockResolvedValueOnce({
        title: 'Edit Test',
        blocks: [{
          id: 'block-1',
          content: 'Initial content',
          segments: [{type: 'text', content: 'Initial content'}],
          depth: 0,
          children: []
        }],
        backlinks: []
      });

      await window.navigateToPage('Edit Test');
      
      const blockText = document.querySelector('.block-text');
      const blockDiv = document.querySelector('.block');
      
      // Simulate editing
      blockText.contentEditable = 'true';
      blockText.classList.add('editing');
      blockText.textContent = 'New content';
      
      // Mock save response
      window.go.main.App.UpdateBlockAtPath.mockResolvedValueOnce({
        action: 'update',
        path: [0],
        block: {
          content: 'New content',
          segments: [{type: 'text', content: 'New content'}]
        }
      });

      // Trigger blur
      blockText.dispatchEvent(new Event('blur'));
      
      // Wait for async save
      await new Promise(resolve => setTimeout(resolve, 50));
      
      // Verify save was called with correct parameters
      expect(window.go.main.App.UpdateBlockAtPath).toHaveBeenCalledWith(
        'Edit Test',
        [0],
        'New content'
      );
    });

    it('should not save if content has not changed', async () => {
      window.go.main.App.GetPage.mockResolvedValueOnce({
        title: 'No Change Test',
        blocks: [{
          id: 'block-1',
          content: 'Same content',
          segments: [{type: 'text', content: 'Same content'}],
          depth: 0,
          children: []
        }],
        backlinks: []
      });

      await window.navigateToPage('No Change Test');
      
      const blockText = document.querySelector('.block-text');
      
      // Focus and blur without changing content
      blockText.dispatchEvent(new Event('focus'));
      blockText.dispatchEvent(new Event('blur'));
      
      await new Promise(resolve => setTimeout(resolve, 50));
      
      // Verify no save was triggered
      expect(window.go.main.App.UpdateBlockAtPath).not.toHaveBeenCalled();
    });
  });

  describe('Back button', () => {
    it('should save edits before going back', async () => {
      // Mock pages
      const pageA = {
        title: 'Page A',
        blocks: [{
          id: 'block-a',
          content: 'Page A content',
          segments: [{type: 'text', content: 'Page A content'}],
          depth: 0,
          children: []
        }],
        backlinks: []
      };
      
      const pageB = {
        title: 'Page B', 
        blocks: [{
          id: 'block-b',
          content: 'Page B content',
          segments: [{type: 'text', content: 'Page B content'}],
          depth: 0,
          children: []
        }],
        backlinks: []
      };

      // Navigate A -> B
      window.go.main.App.GetPage.mockResolvedValueOnce(pageA);
      await window.navigateToPage('Page A');
      
      window.go.main.App.GetPage.mockResolvedValueOnce(pageB);
      await window.navigateToPage('Page B');
      
      // Start editing on Page B
      const blockText = document.querySelector('.block-text');
      blockText.contentEditable = 'true';
      blockText.classList.add('editing');
      blockText.textContent = 'Edited Page B content';
      
      // Mock save and page load for going back
      window.go.main.App.UpdateBlockAtPath.mockResolvedValueOnce({
        action: 'update',
        path: [0],
        block: {
          content: 'Edited Page B content',
          segments: [{type: 'text', content: 'Edited Page B content'}]
        }
      });
      
      window.go.main.App.GetPage.mockResolvedValueOnce(pageA);
      
      // Click back button
      const backButton = document.getElementById('backButton');
      backButton.disabled = false;
      backButton.click();
      
      await new Promise(resolve => setTimeout(resolve, 150));
      
      // Verify edit was saved before navigation
      expect(window.go.main.App.UpdateBlockAtPath).toHaveBeenCalledWith(
        'Page B',
        [0],
        'Edited Page B content'
      );
      
      // Verify we're back on Page A
      expect(document.getElementById('pageTitle').textContent).toBe('Page A');
    });
  });
});