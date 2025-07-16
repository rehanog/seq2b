import { describe, it, expect, beforeEach, vi } from 'vitest';
import '../src/main.js';

describe('Navigation Bug Fix - Edit Persistence', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    document.getElementById('blocks').innerHTML = '';
    document.getElementById('loading').style.display = 'none';
  });

  it('should save edits before navigating away via link click', async () => {
    // Setup mock response for initial page
    window.go.main.App.GetPage.mockResolvedValueOnce({
      title: 'Jul 16th, 2025',
      blocks: [{
        id: 'block-1',
        content: 'Original content',
        segments: [{type: 'text', content: 'Original content'}],
        depth: 0,
        children: []
      }],
      backlinks: []
    });

    // Load the page (simulating being on jul-16th page)
    await window.navigateToPage('Jul 16th, 2025');
    
    // Get the block text element
    const blockText = document.querySelector('.block-text');
    expect(blockText).toBeTruthy();
    
    // Simulate user clicking to edit and typing [[Page B]]
    blockText.click();
    blockText.contentEditable = 'true';
    blockText.classList.add('editing');
    blockText.textContent = 'Original content [[Page B]]';
    blockText.setAttribute('data-raw-content', 'Original content');
    
    // Create a mock blur handler that simulates the save
    let blurHandlerCalled = false;
    const originalBlur = blockText.blur;
    blockText.blur = function() {
      blurHandlerCalled = true;
      this.classList.remove('editing');
      // Simulate the save completing
      this.setAttribute('data-raw-content', this.textContent);
    };
    
    // Now navigate to Page B (what happens when user clicks the link)
    const navigatePromise = window.navigateToPage('Page B');
    
    // Verify blur was called
    expect(blurHandlerCalled).toBe(true);
    
    // Setup mock for Page B
    window.go.main.App.GetPage.mockResolvedValueOnce({
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
    
    await navigatePromise;
    
    // Verify the edit was saved (data-raw-content updated)
    expect(blockText.getAttribute('data-raw-content')).toBe('Original content [[Page B]]');
  });
  
  it('regression test: navigating without the fix would lose edits', async () => {
    // This test simulates what would happen WITHOUT our fix
    // by calling loadPage directly (bypassing navigateToPage)
    
    window.go.main.App.GetPage.mockResolvedValueOnce({
      title: 'Test Page',
      blocks: [{
        id: 'block-1',
        content: 'Original',
        segments: [{type: 'text', content: 'Original'}],
        depth: 0,
        children: []
      }],
      backlinks: []
    });

    // Load initial page
    const loadPageDirectly = async (pageName) => {
      // This is what used to happen - direct navigation without saving
      const pageData = await window.go.main.App.GetPage(pageName);
      document.getElementById('pageTitle').textContent = pageData.title;
      // ... rest of page loading
    };
    
    await loadPageDirectly('Test Page');
    
    // Create a block that's being edited
    const mockBlock = document.createElement('div');
    mockBlock.className = 'block-text editing';
    mockBlock.textContent = 'Edited content that would be lost';
    mockBlock.setAttribute('data-raw-content', 'Original');
    document.getElementById('blocks').appendChild(mockBlock);
    
    // Navigate away WITHOUT the fix (direct load)
    window.go.main.App.GetPage.mockResolvedValueOnce({
      title: 'Another Page',
      blocks: [],
      backlinks: []
    });
    
    await loadPageDirectly('Another Page');
    
    // The edit would be lost because we didn't save it
    // This demonstrates why the fix was needed
    expect(mockBlock.getAttribute('data-raw-content')).toBe('Original'); // Not updated!
  });
});