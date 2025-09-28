import './style.css';
import { GetPage, GetPageList, UpdateBlock, AddBlock, UpdateBlockAtPath, AddBlockAtPath, IsTestMode, CaptureDOM, LogUserAction, CaptureNavigationHistory, GetAsset, LogResourceError } from '../wailsjs/go/main/App';

// Application state
let currentPage = getTodayPageName();
let navigationHistory = [];
let testMode = false;

// Test mode utilities
async function checkTestMode() {
    try {
        testMode = await IsTestMode();
        if (testMode) {
            console.log('Test mode enabled - capturing outputs');
            setupResourceErrorMonitoring();
        }
    } catch (err) {
        console.error('Failed to check test mode:', err);
    }
}

// Setup resource error monitoring for test mode
function setupResourceErrorMonitoring() {
    // Monitor image loading errors
    document.addEventListener('error', async function(e) {
        if (e.target.tagName === 'IMG' && testMode) {
            try {
                await LogResourceError({
                    type: 'image',
                    src: e.target.src,
                    originalSrc: e.target.getAttribute('src') || e.target.getAttribute('data-asset-path'),
                    alt: e.target.alt,
                    message: 'Failed to load image',
                    naturalWidth: e.target.naturalWidth,
                    naturalHeight: e.target.naturalHeight,
                    complete: e.target.complete
                });
            } catch (err) {
                console.error('Failed to log resource error:', err);
            }
        }
    }, true);
}

async function captureCurrentDOM() {
    if (!testMode) return;
    
    try {
        const html = document.documentElement.outerHTML;
        await CaptureDOM(currentPage, html);
    } catch (err) {
        console.error('Failed to capture DOM:', err);
    }
}

async function logAction(action) {
    if (!testMode) return;
    
    try {
        await LogUserAction(action);
    } catch (err) {
        console.error('Failed to log action:', err);
    }
}

async function captureNavHistory() {
    if (!testMode) return;
    
    try {
        await CaptureNavigationHistory([...navigationHistory, currentPage]);
    } catch (err) {
        console.error('Failed to capture navigation history:', err);
    }
}

// Path calculation utilities
function getBlockPath(blockDiv) {
    const path = [];
    let current = blockDiv;
    
    while (current && current.classList.contains('block')) {
        // Find index among siblings
        const parent = current.parentNode;
        const siblings = Array.from(parent.children).filter(child => 
            child.classList.contains('block')
        );
        const index = siblings.indexOf(current);
        
        if (index !== -1) {
            path.unshift(index); // Add to beginning
        }
        
        // Move up to parent block
        if (parent.classList.contains('block-children')) {
            current = parent.parentNode; // Get parent block
        } else {
            // Reached top level
            break;
        }
    }
    
    return path;
}

// Find block element by path
function findBlockByPath(path) {
    if (!path || path.length === 0) return null;
    
    let container = blocksContainer;
    let blockDiv = null;
    
    for (let i = 0; i < path.length; i++) {
        const index = path[i];
        const blocks = Array.from(container.children).filter(child => 
            child.classList.contains('block')
        );
        
        if (index >= 0 && index < blocks.length) {
            blockDiv = blocks[index];
            
            // If not last element, move to children container
            if (i < path.length - 1) {
                const childrenContainer = blockDiv.querySelector('.block-children');
                if (!childrenContainer) return null;
                container = childrenContainer;
            }
        } else {
            return null;
        }
    }
    
    return blockDiv;
}

// DOM elements
const pageTitle = document.getElementById('pageTitle');
const backButton = document.getElementById('backButton');
const blocksContainer = document.getElementById('blocks');
const linkedReferencesContainer = document.getElementById('linkedReferencesList');
const unlinkedReferencesContainer = document.getElementById('unlinkedReferencesList');
const loadingDiv = document.getElementById('loading');

// Initialize the application
document.addEventListener('DOMContentLoaded', async () => {
    await checkTestMode();
    loadPage(currentPage);
    setupEventListeners();
});

// Event listeners
function setupEventListeners() {
    backButton.addEventListener('click', goBack);
    
    // Add home button listener
    const homeButton = document.getElementById('homeButton');
    if (homeButton) {
        homeButton.addEventListener('click', goToToday);
    }
}

// Load and display a page
async function loadPage(pageName) {
    try {
        loadingDiv.style.display = 'block';
        blocksContainer.innerHTML = '';
        linkedReferencesContainer.innerHTML = '';
        unlinkedReferencesContainer.innerHTML = '';
        
        // Clear properties
        const propertiesContainer = document.getElementById('pageProperties');
        propertiesContainer.innerHTML = '';
        propertiesContainer.classList.remove('has-properties');
        
        // Update navigation BEFORE loading (so it works even if page is auto-created)
        if (currentPage !== pageName) {
            navigationHistory.push(currentPage);
            backButton.disabled = false;
        }
        
        const pageData = await GetPage(pageName);
        
        // Update current page BEFORE rendering (so blocks get correct page name)
        currentPage = pageName;
        
        // Update page title
        pageTitle.textContent = pageData.title;
        
        // Render page properties
        renderPageProperties(pageData.properties);
        
        // Render blocks
        renderBlocks(pageData.blocks, blocksContainer);
        
        // Load any asset images
        await loadAssets();
        
        // Render linked references
        renderLinkedReferences(pageData.backlinks);
        
        // Show/hide references sections based on content
        updateReferencesSectionVisibility();
        
        loadingDiv.style.display = 'none';
        
        // Focus on first block if it's empty (new page)
        if (pageData.blocks.length === 1 && !pageData.blocks[0].content.trim()) {
            const firstBlock = blocksContainer.querySelector('.block-text');
            if (firstBlock) {
                firstBlock.contentEditable = 'true';
                firstBlock.focus();
            }
        }
        
        // Load any assets that need the GetAsset API
        // Add small delay to ensure DOM is fully rendered
        console.log('Scheduling loadAssets call...');
        setTimeout(async () => {
            console.log('Calling loadAssets now...');
            await loadAssets();
            console.log('loadAssets completed');
        }, 10);
        
        // Capture DOM and navigation history in test mode
        await captureCurrentDOM();
        await captureNavHistory();
        
        // Validate all images loaded correctly in test mode
        if (testMode) {
            await validateAllImages();
        }
        
    } catch (error) {
        console.error('Error loading page:', error);
        loadingDiv.innerHTML = `Error loading page: ${error.message || error}`;
        // Revert navigation on error
        if (navigationHistory.length > 0 && navigationHistory[navigationHistory.length - 1] === currentPage) {
            navigationHistory.pop();
            if (navigationHistory.length === 0) {
                backButton.disabled = true;
            }
        }
    }
}

// Load assets that need the GetAsset API
async function loadAssets() {
    const assetImages = document.querySelectorAll('img[data-asset-path]');
    console.log('loadAssets: Found', assetImages.length, 'images with data-asset-path');
    
    for (const img of assetImages) {
        const assetPath = img.getAttribute('data-asset-path');
        if (assetPath) {
            // GetAsset expects path to include 'assets/' prefix
            const fullAssetPath = assetPath.startsWith('assets/') ? assetPath : 'assets/' + assetPath;
            
            try {
                console.log('Loading asset:', fullAssetPath);
                const assetData = await GetAsset(fullAssetPath);
                console.log('Asset data received:', assetData ? assetData.substring(0, 100) + '...' : 'null');
                
                // If it's SVG (text), create a data URL
                if (assetPath.endsWith('.svg')) {
                    img.src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(assetData);
                } else {
                    // For other formats, assume it's already a data URL from backend
                    img.src = assetData;
                }
                
                img.classList.remove('loading-asset');
                img.removeAttribute('data-asset-path');
            } catch (err) {
                console.error('Failed to load asset:', fullAssetPath, err);
                console.error('Error details:', err.message || err);
                img.alt = img.alt + ' (failed to load)';
                img.classList.add('failed-asset');
                img.classList.remove('loading-asset');
            }
        }
    }
}

// Validate all images loaded correctly (for test mode)
async function validateAllImages() {
    const images = document.querySelectorAll('img');
    for (const img of images) {
        // Check if image failed to load
        if (!img.complete || img.naturalWidth === 0) {
            try {
                await LogResourceError({
                    type: 'image_validation_failed',
                    src: img.src,
                    originalSrc: img.getAttribute('src') || img.getAttribute('data-asset-path'),
                    alt: img.alt,
                    message: 'Image not loaded properly during validation',
                    complete: img.complete,
                    naturalWidth: img.naturalWidth,
                    naturalHeight: img.naturalHeight,
                    classList: Array.from(img.classList)
                });
            } catch (err) {
                console.error('Failed to log image validation error:', err);
            }
        }
    }
}

// Render blocks recursively
function renderBlocks(blocks, container) {
    blocks.forEach(block => {
        const blockElement = createBlockElement(block);
        container.appendChild(blockElement);
    });
}

// Count all children recursively
function countAllChildren(block) {
    let count = 0;
    if (block.children) {
        count = block.children.length;
        block.children.forEach(child => {
            count += countAllChildren(child);
        });
    }
    return count;
}

// Toggle all children blocks recursively
function toggleAllChildren(blockElement, collapse) {
    const allChildBlocks = blockElement.querySelectorAll('.block-has-children');
    allChildBlocks.forEach(childBlock => {
        const toggle = childBlock.querySelector('.collapse-toggle');
        const blockId = childBlock.getAttribute('data-block-id');
        const collapseKey = `collapsed_${currentPage}_${blockId}`;
        
        if (collapse) {
            childBlock.classList.add('collapsed');
            if (toggle) toggle.innerHTML = '▶';
            localStorage.setItem(collapseKey, 'true');
        } else {
            childBlock.classList.remove('collapsed');
            if (toggle) toggle.innerHTML = '▼';
            localStorage.removeItem(collapseKey);
        }
    });
    
    // Also toggle the main block
    const mainToggle = blockElement.querySelector('.collapse-toggle');
    const mainBlockId = blockElement.getAttribute('data-block-id');
    const mainCollapseKey = `collapsed_${currentPage}_${mainBlockId}`;
    
    if (collapse) {
        blockElement.classList.add('collapsed');
        if (mainToggle) mainToggle.innerHTML = '▶';
        localStorage.setItem(mainCollapseKey, 'true');
    } else {
        blockElement.classList.remove('collapsed');
        if (mainToggle) mainToggle.innerHTML = '▼';
        localStorage.removeItem(mainCollapseKey);
    }
}

// Create a block element
function createBlockElement(block) {
    const blockDiv = document.createElement('div');
    blockDiv.className = 'block';
    blockDiv.setAttribute('data-depth', block.depth);
    blockDiv.setAttribute('data-block-id', block.id);
    blockDiv.setAttribute('data-page', currentPage); // Store page name for saving
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'block-content';
    
    // Create a consistent prefix element with bullet always visible
    const prefixDiv = document.createElement('div');
    prefixDiv.className = 'block-prefix';
    
    // Add collapse/expand toggle for blocks with children
    if (block.children && block.children.length > 0) {
        const collapseToggle = document.createElement('div');
        collapseToggle.className = 'collapse-toggle';
        collapseToggle.innerHTML = '▼'; // Down arrow when expanded
        collapseToggle.title = 'Click to collapse/expand';
        
        // Check if this block should be collapsed from localStorage
        const collapseKey = `collapsed_${currentPage}_${block.id}`;
        const isCollapsed = localStorage.getItem(collapseKey) === 'true';
        
        if (isCollapsed) {
            blockDiv.classList.add('collapsed');
            collapseToggle.innerHTML = '▶'; // Right arrow when collapsed
        }
        
        collapseToggle.addEventListener('click', function(e) {
            e.stopPropagation();
            const isCurrentlyCollapsed = blockDiv.classList.contains('collapsed');
            
            // Alt+Click to toggle all children recursively
            if (e.altKey) {
                toggleAllChildren(blockDiv, !isCurrentlyCollapsed);
            } else {
                if (isCurrentlyCollapsed) {
                    blockDiv.classList.remove('collapsed');
                    collapseToggle.innerHTML = '▼';
                    localStorage.removeItem(collapseKey);
                } else {
                    blockDiv.classList.add('collapsed');
                    collapseToggle.innerHTML = '▶';
                    localStorage.setItem(collapseKey, 'true');
                }
            }
        });
        
        prefixDiv.appendChild(collapseToggle);
    }
    
    // Always show bullet as a circle
    const bullet = document.createElement('div');
    bullet.className = 'block-bullet';
    prefixDiv.appendChild(bullet);
    
    // Add TODO state or checkbox after bullet if present
    if (block.todoState) {
        const todoSpan = document.createElement('span');
        todoSpan.className = `todo-state todo-${block.todoState.toLowerCase()}`;
        todoSpan.textContent = block.todoState;
        if (block.priority) {
            const prioritySpan = document.createElement('span');
            prioritySpan.className = 'todo-priority';
            prioritySpan.textContent = `[#${block.priority}]`;
            todoSpan.appendChild(prioritySpan);
        }
        prefixDiv.appendChild(todoSpan);
    } else if (block.checkboxState) {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.className = 'block-checkbox';
        checkbox.checked = block.checkboxState === '[x]';
        checkbox.disabled = true; // For now, make it read-only
        prefixDiv.appendChild(checkbox);
    }
    
    contentDiv.appendChild(prefixDiv);
    
    // Create editable text div
    const textDiv = document.createElement('div');
    textDiv.className = 'block-text';
    textDiv.contentEditable = 'false'; // Start with editing disabled
    
    // Use segments if available, otherwise fall back to htmlContent
    if (block.segments && block.segments.length > 0) {
        // Debug log for image blocks
        if (block.content && block.content.includes('![')) {
            console.log('Block with image:', block.content);
            console.log('Segments:', block.segments);
        }
        textDiv.innerHTML = renderSegmentsToHTML(block.segments);
    } else if (block.htmlContent) {
        textDiv.innerHTML = processLinksInHTML(block.htmlContent);
    } else if (block.content) {
        // If no segments or HTML, parse the content on the fly
        const segments = parseMarkdownToSegments(block.content);
        textDiv.innerHTML = renderSegmentsToHTML(segments);
    } else {
        // For empty blocks, add a zero-width space to ensure editability
        textDiv.innerHTML = '\u200B';
    }
    
    // Store original content for editing
    textDiv.setAttribute('data-raw-content', block.content);
    
    // Handle focus - show raw markdown
    textDiv.addEventListener('focus', function() {
        this.contentEditable = 'true'; // Enable editing
        this.classList.add('editing');
        const rawContent = this.getAttribute('data-raw-content') || '';
        this.textContent = rawContent;
        
        // Place cursor at end
        const range = document.createRange();
        const sel = window.getSelection();
        if (this.childNodes.length > 0) {
            range.selectNodeContents(this);
            range.collapse(false);
        } else {
            // For empty content, set cursor at beginning
            range.setStart(this, 0);
            range.collapse(true);
        }
        sel.removeAllRanges();
        sel.addRange(range);
    });
    
    // Handle blur - save and show rendered HTML
    textDiv.addEventListener('blur', async function(e) {
        // Prevent re-triggering if already processing
        if (this.dataset.saving === 'true') return;
        
        this.dataset.saving = 'true';
        this.contentEditable = 'false'; // Disable editing
        this.classList.remove('editing');
        const newContent = this.textContent;
        const blockId = blockDiv.getAttribute('data-block-id');
        const pageName = blockDiv.getAttribute('data-page');
        
        // Only save if content actually changed
        if (newContent !== this.getAttribute('data-raw-content')) {
            await saveBlockEdit(pageName, blockId, newContent, this);
        } else {
            // Just restore the HTML if no changes
            if (block.segments && block.segments.length > 0) {
                this.innerHTML = renderSegmentsToHTML(block.segments);
            } else if (block.htmlContent) {
                this.innerHTML = processLinksInHTML(block.htmlContent);
            } else {
                // Parse content on the fly for temporary blocks
                const segments = parseMarkdownToSegments(this.getAttribute('data-raw-content'));
                this.innerHTML = renderSegmentsToHTML(segments);
            }
            // Load any asset images
            await loadAssets();
        }
        
        this.dataset.saving = 'false';
    });
    
    // Handle clicks - prevent edit mode when clicking links
    textDiv.addEventListener('click', function(e) {
        // If already editing, prevent link navigation
        if (this.classList.contains('editing')) {
            e.preventDefault();
            e.stopPropagation();
            return;
        }
        
        // Check if click was on a link (either page link or PDF link)
        const link = e.target.closest('.page-link, .pdf-link');
        if (link) {
            // Let the link handle the click, don't enter edit mode
            e.stopPropagation();
            return;
        }
        
        // Otherwise, enter edit mode by temporarily enabling contentEditable and focusing
        this.contentEditable = 'true';
        this.focus();
    });
    
    // Handle keyboard navigation
    textDiv.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleEnterKey(this, blockDiv);
        } else if (e.key === 'Escape') {
            e.preventDefault();
            this.blur();
        } else if (e.key === 'Tab') {
            e.preventDefault();
            if (e.shiftKey) {
                handleOutdent(blockDiv);
            } else {
                handleIndent(blockDiv);
            }
        } else if (e.key === 'Backspace' || e.key === 'Delete') {
            const selection = window.getSelection();
            const range = selection.getRangeAt(0);
            const cursorPos = range.startOffset;
            
            // If at beginning of block and pressing backspace, merge with previous
            if (e.key === 'Backspace' && cursorPos === 0 && this.textContent.length > 0) {
                e.preventDefault();
                handleMergeWithPrevious(this, blockDiv);
            }
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            handleArrowUp(this, blockDiv);
        } else if (e.key === 'ArrowDown') {
            e.preventDefault();
            handleArrowDown(this, blockDiv);
        }
    });
    
    contentDiv.appendChild(textDiv);
    blockDiv.appendChild(contentDiv);
    
    // Add children if any
    if (block.children && block.children.length > 0) {
        blockDiv.classList.add('block-has-children');
        const childrenContainer = document.createElement('div');
        childrenContainer.className = 'block-children';
        renderBlocks(block.children, childrenContainer);
        blockDiv.appendChild(childrenContainer);
    }
    
    return blockDiv;
}

// Save block edit
async function saveBlockEdit(pageName, blockId, newContent, textDiv) {
    try {
        // Log edit action
        await logAction(`Edit block in ${pageName}: "${newContent.substring(0, 50)}..."`);
        
        // Store the new raw content
        textDiv.setAttribute('data-raw-content', newContent);
        
        // Get block path
        const blockDiv = textDiv.closest('.block');
        const path = getBlockPath(blockDiv);
        
        // Update the backend using path-based API
        const delta = await UpdateBlockAtPath(pageName, path, newContent);
        
        // Apply the delta update
        if (delta && delta.action === 'update') {
            const updatedBlock = delta.block;
            
            // Update the HTML content without reloading the whole page
            if (updatedBlock.segments && updatedBlock.segments.length > 0) {
                textDiv.innerHTML = renderSegmentsToHTML(updatedBlock.segments);
            } else if (updatedBlock.htmlContent) {
                textDiv.innerHTML = processLinksInHTML(updatedBlock.htmlContent);
            } else if (updatedBlock.content) {
                // Fallback: parse content on the fly if no segments or HTML
                const segments = parseMarkdownToSegments(updatedBlock.content);
                textDiv.innerHTML = renderSegmentsToHTML(segments);
            } else {
                // Last resort: use the raw content
                textDiv.innerHTML = escapeHtml(newContent);
            }
            
            // Load any new asset images
            await loadAssets();
            
            // Update TODO/checkbox state if changed
            if (updatedBlock.todoState !== undefined || updatedBlock.checkboxState !== undefined) {
                updateBlockPrefix(blockDiv, updatedBlock);
            }
        } else {
            // Fallback to old API if new one fails
            await UpdateBlock(pageName, blockId, newContent);
            // Reload page to sync
            await loadPage(pageName);
        }
        
    } catch (error) {
        console.error('Failed to save block:', error);
        // Revert to original content on error
        textDiv.textContent = textDiv.getAttribute('data-raw-content');
    }
}

// Update block prefix (TODO state, checkbox, etc)
function updateBlockPrefix(blockDiv, blockData) {
    const prefixDiv = blockDiv.querySelector('.block-prefix');
    if (!prefixDiv) return;

    // Preserve collapse toggle and bullet
    const collapseToggle = prefixDiv.querySelector('.collapse-toggle');
    const bullet = prefixDiv.querySelector('.block-bullet');

    // Clear prefix content
    prefixDiv.innerHTML = '';

    // Re-add collapse toggle first (if it exists)
    if (collapseToggle) prefixDiv.appendChild(collapseToggle);

    // Then add bullet
    if (bullet) prefixDiv.appendChild(bullet);
    
    // Add TODO state or checkbox
    if (blockData.todoState) {
        const todoSpan = document.createElement('span');
        todoSpan.className = `todo-state todo-${blockData.todoState.toLowerCase()}`;
        todoSpan.textContent = blockData.todoState;
        if (blockData.priority) {
            const prioritySpan = document.createElement('span');
            prioritySpan.className = 'todo-priority';
            prioritySpan.textContent = `[#${blockData.priority}]`;
            todoSpan.appendChild(prioritySpan);
        }
        prefixDiv.appendChild(todoSpan);
    } else if (blockData.checkboxState) {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.className = 'block-checkbox';
        checkbox.checked = blockData.checkboxState === '[x]';
        checkbox.disabled = true;
        prefixDiv.appendChild(checkbox);
    }
}

// Process HTML content to make links clickable
function processLinksInHTML(html) {
    // Convert <a href="page">text</a> to clickable links
    return html.replace(/<a href="([^"]+)">([^<]+)<\/a>/g, (match, href, text) => {
        return `<a href="#" class="page-link" onclick="navigateToPage('${text}')">${text}</a>`;
    });
}

// Render segments to HTML
function renderSegmentsToHTML(segments) {
    if (!segments || segments.length === 0) {
        return '';
    }
    
    return segments.map(segment => {
        switch (segment.type) {
            case 'bold':
                return `<b>${escapeHtml(segment.content)}</b>`;
            case 'italic':
                return `<i>${escapeHtml(segment.content)}</i>`;
            case 'link':
                // Check if this is a PDF link
                if (segment.target && (segment.target.toLowerCase().endsWith('.pdf') || segment.content.toLowerCase().endsWith('.pdf'))) {
                    console.log('Rendering PDF link:', segment.target);
                    return `<a href="#" class="pdf-link" onclick="openPDF('${escapeHtml(segment.target)}'); return false;">${escapeHtml(segment.content)}</a>`;
                } else {
                    return `<a href="#" class="page-link" onclick="navigateToPage('${escapeHtml(segment.target)}'); return false;">${escapeHtml(segment.content)}</a>`;
                }
            case 'image':
                console.log('Processing image segment:', segment.target, 'content:', segment.content);
                // Check if this is actually a PDF (Logseq uses image syntax for PDFs)
                if (segment.target && segment.target.toLowerCase().endsWith('.pdf')) {
                    console.log('Rendering PDF as image syntax:', segment.target);
                    return `<a href="#" class="pdf-link" onclick="openPDF('${escapeHtml(segment.target)}'); return false;">${escapeHtml(segment.content || segment.alt || 'PDF')}</a>`;
                }
                // Check if this is a relative asset path
                else if (segment.target && segment.target.startsWith('../assets/')) {
                    // Transform to use asset loading
                    const assetPath = segment.target.substring('../assets/'.length);
                    console.log('Rendering as asset image:', assetPath);
                    return `<img data-asset-path="${escapeHtml(assetPath)}" alt="${escapeHtml(segment.alt || segment.content)}" class="embedded-image loading-asset">`;
                } else {
                    // Regular image (absolute URL or other path)
                    console.log('Rendering as regular image:', segment.target);
                    return `<img src="${escapeHtml(segment.target)}" alt="${escapeHtml(segment.alt || segment.content)}" class="embedded-image">`;
                }
            case 'tag':
                return `<span class="tag">#${escapeHtml(segment.content)}</span>`;
            case 'blockRef':
                return `<span class="block-reference" title="Block reference: ${escapeHtml(segment.target)}">((${escapeHtml(segment.content)}))</span>`;
            case 'property':
                return `<span class="property">${escapeHtml(segment.content)}</span>`;
            case 'blockId':
                return `<span class="block-id" title="Block ID: ${escapeHtml(segment.target)}">${escapeHtml(segment.content)}</span>`;
            case 'query':
                return `<span class="query-block">${escapeHtml(segment.content)}</span>`;
            case 'embed':
                return `<span class="embed-block">${escapeHtml(segment.content)}</span>`;
            case 'strikethrough':
                return `<del>${escapeHtml(segment.content)}</del>`;
            case 'highlight':
                return `<mark>${escapeHtml(segment.content)}</mark>`;
            case 'text':
            default:
                return escapeHtml(segment.content);
        }
    }).join('');
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Simple client-side markdown parser for temporary blocks
function parseMarkdownToSegments(text) {
    if (!text) return [];
    
    const segments = [];
    
    // Combined pattern for links and images
    const pattern = /(\[\[([^\]]+)\]\]|!\[([^\]]*)\]\(([^)]+)\))/g;
    let lastIndex = 0;
    let match;
    
    while ((match = pattern.exec(text)) !== null) {
        // Add text before the match
        if (match.index > lastIndex) {
            segments.push({
                type: 'text',
                content: text.substring(lastIndex, match.index)
            });
        }
        
        // Check if it's an image or a link
        if (match[0].startsWith('![')) {
            // Image: match[3] is alt text, match[4] is src
            segments.push({
                type: 'image',
                content: match[3] || '', // alt text
                alt: match[3] || '',
                target: match[4]
            });
        } else {
            // Link: match[2] is the link target
            segments.push({
                type: 'link',
                content: match[2],
                target: match[2]
            });
        }
        
        lastIndex = match.index + match[0].length;
    }
    
    // Add remaining text
    if (lastIndex < text.length) {
        segments.push({
            type: 'text',
            content: text.substring(lastIndex)
        });
    }
    
    return segments.length > 0 ? segments : [{type: 'text', content: text}];
}

// PDF Viewer functionality
let currentPDF = null;
let sidebarResizer = null;

window.openPDF = async function(pdfPath) {
    console.log('Opening PDF:', pdfPath);
    
    const sidebar = document.getElementById('rightSidebar');
    const resizeHandle = document.getElementById('resizeHandle');
    const iframe = document.getElementById('pdfFrame');
    const titleSpan = document.querySelector('.sidebar-title');
    const toggleBtn = document.getElementById('toggleSidebar');
    
    // Store current PDF path
    currentPDF = pdfPath;
    
    // Update title
    titleSpan.textContent = pdfPath.split('/').pop();
    
    try {
        // Check if it's a relative path (../assets/ or similar)
        if (pdfPath.startsWith('../')) {
            // Use GetAsset to load the PDF
            const assetPath = pdfPath.substring(3); // Remove ../
            const base64Data = await GetAsset(assetPath);
            
            // The backend now returns complete data URLs for PDFs
            iframe.src = base64Data;
        } else {
            // For absolute paths or URLs, load directly
            iframe.src = pdfPath;
        }
        
        // Show the sidebar and resize handle
        sidebar.style.display = 'flex';
        resizeHandle.style.display = 'block';
        toggleBtn.style.display = 'block';
        toggleBtn.textContent = '⬅ Hide Sidebar';
        
        // Log PDF open action
        await logAction(`Opened PDF: ${pdfPath}`);
    } catch (error) {
        console.error('Error loading PDF:', error);
        alert(`Failed to load PDF: ${error.message}`);
    }
}

// Setup PDF sidebar controls
document.addEventListener('DOMContentLoaded', function() {
    const sidebar = document.getElementById('rightSidebar');
    const resizeHandle = document.getElementById('resizeHandle');
    const closeBtn = document.getElementById('sidebarClose');
    const toggleBtn = document.getElementById('toggleSidebar');
    const mainContainer = document.getElementById('mainContainer');
    
    // Close button
    closeBtn.addEventListener('click', function() {
        closeSidebar();
    });
    
    // Toggle button
    toggleBtn.addEventListener('click', function() {
        if (sidebar.style.display === 'none') {
            // Show sidebar
            sidebar.style.display = 'flex';
            resizeHandle.style.display = 'block';
            toggleBtn.textContent = '⬅ Hide Sidebar';
        } else {
            // Hide sidebar temporarily
            sidebar.style.display = 'none';
            resizeHandle.style.display = 'none';
            toggleBtn.textContent = '➡ Show Sidebar';
        }
    });
    
    // Escape key to close
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape' && sidebar.style.display === 'flex') {
            // Only close if not editing
            if (!document.querySelector('.block-text.editing')) {
                closeSidebar();
            }
        }
    });
    
    // Resize functionality
    let isResizing = false;
    let startX = 0;
    let startWidth = 0;
    
    resizeHandle.addEventListener('mousedown', function(e) {
        isResizing = true;
        startX = e.pageX;
        startWidth = sidebar.offsetWidth;
        document.body.style.cursor = 'ew-resize';
        e.preventDefault();
    });
    
    document.addEventListener('mousemove', function(e) {
        if (!isResizing) return;
        
        const currentX = e.pageX;
        const diff = startX - currentX;
        const newWidth = startWidth + diff;
        
        // Limit sidebar width between 200px and 80% of container
        const minWidth = 200;
        const maxWidth = mainContainer.offsetWidth * 0.8;
        
        if (newWidth >= minWidth && newWidth <= maxWidth) {
            sidebar.style.width = newWidth + 'px';
        }
    });
    
    document.addEventListener('mouseup', function() {
        if (isResizing) {
            isResizing = false;
            document.body.style.cursor = '';
        }
    });
});

// Close sidebar helper
function closeSidebar() {
    const sidebar = document.getElementById('rightSidebar');
    const resizeHandle = document.getElementById('resizeHandle');
    const iframe = document.getElementById('pdfFrame');
    const toggleBtn = document.getElementById('toggleSidebar');
    
    sidebar.style.display = 'none';
    resizeHandle.style.display = 'none';
    toggleBtn.style.display = 'none';
    iframe.src = ''; // Clear the iframe
    currentPDF = null;
}

// Navigate to a page
window.navigateToPage = async function(pageName) {
    // Log navigation action
    await logAction(`Navigate to page: ${pageName}`);
    
    // Save any active edits before navigating
    const editingBlock = document.querySelector('.block-text.editing');
    if (editingBlock) {
        // Trigger blur to save
        editingBlock.blur();
        // Wait a bit for save to complete
        await new Promise(resolve => setTimeout(resolve, 100));
    }
    
    loadPage(pageName);
};

// Go back to previous page
async function goBack() {
    if (navigationHistory.length > 0) {
        // Save any active edits before navigating
        const editingBlock = document.querySelector('.block-text.editing');
        if (editingBlock) {
            // Trigger blur to save
            editingBlock.blur();
            // Wait a bit for save to complete
            await new Promise(resolve => setTimeout(resolve, 100));
        }
        
        const previousPage = navigationHistory.pop();
        
        // Don't add to history when going back
        const tempCurrent = currentPage;
        currentPage = previousPage;
        
        loadPage(previousPage);
        
        // Disable back button if no more history
        if (navigationHistory.length === 0) {
            backButton.disabled = true;
        }
    }
}

// Get today's page name in Logseq format
function getTodayPageName() {
    const today = new Date();
    const monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    const month = monthNames[today.getMonth()];
    const day = today.getDate();
    const year = today.getFullYear();
    
    // Add ordinal suffix
    let suffix = "th";
    if (day === 1 || day === 21 || day === 31) suffix = "st";
    else if (day === 2 || day === 22) suffix = "nd";
    else if (day === 3 || day === 23) suffix = "rd";
    
    return `${month} ${day}${suffix}, ${year}`;
}

// Go to today's journal page
async function goToToday() {
    // Save any active edits before navigating
    const editingBlock = document.querySelector('.block-text.editing');
    if (editingBlock) {
        // Trigger blur to save
        editingBlock.blur();
        // Wait a bit for save to complete
        await new Promise(resolve => setTimeout(resolve, 100));
    }
    
    const todayPageName = getTodayPageName();
    // Navigate to today's page
    loadPage(todayPageName);
}

// Render page properties
function renderPageProperties(properties) {
    const container = document.getElementById('pageProperties');
    
    if (!properties || Object.keys(properties).length === 0) {
        return;
    }
    
    // Show the container
    container.classList.add('has-properties');
    
    // Render each property
    for (const [key, value] of Object.entries(properties)) {
        const propertyItem = document.createElement('div');
        propertyItem.className = 'property-item';
        
        const keyElement = document.createElement('span');
        keyElement.className = 'property-key';
        keyElement.textContent = key;
        
        const valueElement = document.createElement('span');
        valueElement.className = 'property-value';
        valueElement.textContent = value;
        
        propertyItem.appendChild(keyElement);
        propertyItem.appendChild(valueElement);
        container.appendChild(propertyItem);
    }
}

// Render linked references
function renderLinkedReferences(backlinks) {
    const container = document.getElementById('linkedReferences');
    
    if (!backlinks || backlinks.length === 0) {
        linkedReferencesContainer.innerHTML = '<div class="empty-state">No linked references</div>';
        container.style.display = 'none';
        return;
    }
    
    container.style.display = 'block';
    
    backlinks.forEach(backlink => {
        const referenceDiv = document.createElement('div');
        referenceDiv.className = 'reference-item';
        
        const sourceDiv = document.createElement('div');
        sourceDiv.className = 'reference-source';
        sourceDiv.textContent = backlink.sourcePage;
        sourceDiv.onclick = () => navigateToPage(backlink.sourcePage);
        
        const contentDiv = document.createElement('div');
        contentDiv.className = 'reference-content';
        // TODO: In future, show actual block content with context
        contentDiv.textContent = `${backlink.blockIds.length} reference${backlink.blockIds.length > 1 ? 's' : ''}`;
        
        referenceDiv.appendChild(sourceDiv);
        referenceDiv.appendChild(contentDiv);
        linkedReferencesContainer.appendChild(referenceDiv);
    });
}

// Update references section visibility
function updateReferencesSectionVisibility() {
    const referencesSection = document.getElementById('references');
    
    // Show references section only if there are any references
    const hasLinkedRefs = linkedReferencesContainer.children.length > 0 && 
                         !linkedReferencesContainer.querySelector('.empty-state');
    const hasUnlinkedRefs = unlinkedReferencesContainer.children.length > 0;
    
    if (hasLinkedRefs || hasUnlinkedRefs) {
        referencesSection.style.display = 'block';
        setupReferenceCollapsible();
    } else {
        referencesSection.style.display = 'none';
    }
}

// Setup collapsible references sections
function setupReferenceCollapsible() {
    const headers = document.querySelectorAll('.references-header');
    headers.forEach(header => {
        // Remove any existing click handlers
        const newHeader = header.cloneNode(true);
        header.parentNode.replaceChild(newHeader, header);
        
        newHeader.style.cursor = 'pointer';
        newHeader.onclick = function() {
            const section = this.parentElement;
            section.classList.toggle('collapsed');
            const content = section.querySelector('[id$="List"]');
            if (section.classList.contains('collapsed')) {
                content.style.display = 'none';
            } else {
                content.style.display = 'block';
            }
        };
    });
}

// Error handling
window.addEventListener('error', (event) => {
    console.error('Application error:', event.error);
});

// Handle Enter key - split block at cursor
async function handleEnterKey(textDiv, blockDiv) {
    const selection = window.getSelection();
    const range = selection.getRangeAt(0);
    const textContent = textDiv.textContent;
    
    // Get the actual cursor position
    let cursorPos = 0;
    if (range.startContainer.nodeType === Node.TEXT_NODE) {
        cursorPos = range.startOffset;
    } else {
        // If not in a text node, we're probably at the end
        cursorPos = textContent.length;
    }
    
    // Split content at cursor position
    const beforeCursor = textContent.substring(0, cursorPos);
    const afterCursor = textContent.substring(cursorPos);
    
    // Update current block with text before cursor
    textDiv.textContent = beforeCursor;
    textDiv.blur(); // This will save the current block
    
    try {
        // Get current block path
        const currentPath = getBlockPath(blockDiv);
        const currentDepth = parseInt(blockDiv.getAttribute('data-depth'));
        const pageName = currentPage;
        
        // Determine where to insert the new block
        let insertPath;
        
        // Check if current block already has children
        const childrenContainer = blockDiv.querySelector('.block-children');
        
        if (childrenContainer && afterCursor) {
            // If has children AND we're moving text, insert as first child
            insertPath = [...currentPath, 0]; // Insert at beginning of children
        } else {
            // Otherwise insert after current block at same level
            if (currentPath.length === 1) {
                // Top level block - insert after it
                insertPath = [currentPath[0] + 1];
            } else {
                // Nested block - insert after it in parent's children
                const parentPath = currentPath.slice(0, -1);
                const indexInParent = currentPath[currentPath.length - 1];
                insertPath = [...parentPath, indexInParent + 1];
            }
        }
        
        // Call backend to create the new block using positional API
        const delta = await AddBlockAtPath(pageName, insertPath, afterCursor || '');
        
        // Apply incremental update instead of reloading
        if (delta && delta.action === 'add') {
            // Create the new block element
            const newBlockElement = createBlockElement(delta.block);
            
            // Insert it at the correct position
            if (insertPath.length === 1) {
                // Top level insertion
                const index = insertPath[0];
                const blocks = Array.from(blocksContainer.children).filter(child => 
                    child.classList.contains('block')
                );
                
                if (index >= blocks.length) {
                    blocksContainer.appendChild(newBlockElement);
                } else {
                    blocksContainer.insertBefore(newBlockElement, blocks[index]);
                }
            } else {
                // Nested insertion
                const parentPath = insertPath.slice(0, -1);
                const parentBlock = findBlockByPath(parentPath);
                
                if (parentBlock) {
                    let childrenContainer = parentBlock.querySelector('.block-children');
                    if (!childrenContainer) {
                        parentBlock.classList.add('block-has-children');
                        childrenContainer = document.createElement('div');
                        childrenContainer.className = 'block-children';
                        parentBlock.appendChild(childrenContainer);
                    }
                    
                    const index = insertPath[insertPath.length - 1];
                    const children = Array.from(childrenContainer.children).filter(child => 
                        child.classList.contains('block')
                    );
                    
                    if (index >= children.length) {
                        childrenContainer.appendChild(newBlockElement);
                    } else {
                        childrenContainer.insertBefore(newBlockElement, children[index]);
                    }
                }
            }
            
            // Apply path shifts if any
            if (delta.shifts && delta.shifts.length > 0) {
                applyPathShifts(delta.shifts);
            }
            
            // Focus the new block
            const newTextDiv = newBlockElement.querySelector('.block-text');
            if (newTextDiv) {
                // Store whether we should position cursor at start
                const shouldPositionAtStart = afterCursor.length > 0;
                
                // Enable contentEditable before focusing
                newTextDiv.contentEditable = 'true';
                newTextDiv.focus();
                
                // After focus handler converts to markdown, position cursor
                if (shouldPositionAtStart && afterCursor) {
                    setTimeout(() => {
                        const range = document.createRange();
                        const sel = window.getSelection();
                        if (newTextDiv.childNodes.length > 0 && newTextDiv.childNodes[0]) {
                            range.setStart(newTextDiv.childNodes[0], 0);
                        } else {
                            range.setStart(newTextDiv, 0);
                        }
                        range.collapse(true);
                        sel.removeAllRanges();
                        sel.addRange(range);
                    }, 10);
                }
            }
        } else {
            // Fallback: reload the page
            await loadPage(currentPage);
        }
    } catch (error) {
        console.error('Failed to create new block:', error);
        // Fallback: reload the page
        await loadPage(currentPage);
    }
}

// Apply path shifts after insertion/deletion
function applyPathShifts(shifts) {
    // Path shifts are used to update any stored references
    // For now, we don't store paths, so this is a no-op
    // In the future, if we cache paths, we'd update them here
}

// Handle indent (Tab)
function handleIndent(blockDiv) {
    const currentDepth = parseInt(blockDiv.getAttribute('data-depth'));
    const prevSibling = blockDiv.previousSibling;
    
    // Can only indent if there's a previous sibling to become parent
    if (prevSibling && prevSibling.classList.contains('block')) {
        // Remove from current position
        blockDiv.remove();
        
        // Find or create children container in previous sibling
        let childrenContainer = prevSibling.querySelector('.block-children');
        if (!childrenContainer) {
            prevSibling.classList.add('block-has-children');
            childrenContainer = document.createElement('div');
            childrenContainer.className = 'block-children';
            prevSibling.appendChild(childrenContainer);
        }
        
        // Update depth and add to children
        blockDiv.setAttribute('data-depth', currentDepth + 1);
        childrenContainer.appendChild(blockDiv);
        
        // Refocus the text
        const textDiv = blockDiv.querySelector('.block-text');
        if (textDiv) textDiv.focus();
    }
}

// Handle outdent (Shift+Tab)
function handleOutdent(blockDiv) {
    const currentDepth = parseInt(blockDiv.getAttribute('data-depth'));
    
    // Can only outdent if not at top level
    if (currentDepth > 0) {
        const parentChildren = blockDiv.parentNode;
        const parentBlock = parentChildren.parentNode;
        
        // Remove from current position
        blockDiv.remove();
        
        // Update depth
        blockDiv.setAttribute('data-depth', currentDepth - 1);
        
        // Insert after parent block
        if (parentBlock.nextSibling) {
            parentBlock.parentNode.insertBefore(blockDiv, parentBlock.nextSibling);
        } else {
            parentBlock.parentNode.appendChild(blockDiv);
        }
        
        // If parent has no more children, remove children container
        if (parentChildren.children.length === 0) {
            parentBlock.classList.remove('block-has-children');
            parentChildren.remove();
        }
        
        // Refocus the text
        const textDiv = blockDiv.querySelector('.block-text');
        if (textDiv) textDiv.focus();
    }
}

// Handle merge with previous block (Backspace at start)
function handleMergeWithPrevious(textDiv, blockDiv) {
    const prevBlock = findPreviousEditableBlock(blockDiv);
    if (!prevBlock) return;
    
    const prevTextDiv = prevBlock.querySelector('.block-text');
    if (!prevTextDiv) return;
    
    // Get the content to merge
    const contentToMerge = textDiv.textContent;
    const prevContent = prevTextDiv.getAttribute('data-raw-content') || '';
    
    // Save cursor position (will be at end of previous content)
    const cursorPos = prevContent.length;
    
    // Update previous block with merged content
    prevTextDiv.setAttribute('data-raw-content', prevContent + contentToMerge);
    
    // Remove current block
    blockDiv.remove();
    
    // Focus previous block and position cursor
    prevTextDiv.contentEditable = 'true';
    prevTextDiv.focus();
    
    // Position cursor after the focus handler runs
    setTimeout(() => {
        const range = document.createRange();
        const sel = window.getSelection();
        if (prevTextDiv.childNodes[0]) {
            range.setStart(prevTextDiv.childNodes[0], cursorPos);
            range.collapse(true);
            sel.removeAllRanges();
            sel.addRange(range);
        }
    }, 10);
}

// Handle arrow up navigation
function handleArrowUp(textDiv, blockDiv) {
    const prevBlock = findPreviousEditableBlock(blockDiv);
    if (!prevBlock) return;
    
    const prevTextDiv = prevBlock.querySelector('.block-text');
    if (!prevTextDiv) return;
    
    // Get current cursor position
    const selection = window.getSelection();
    const range = selection.getRangeAt(0);
    const cursorPos = range.startOffset;
    
    // Save current content before blur
    textDiv.blur();
    
    // Focus previous block
    prevTextDiv.contentEditable = 'true';
    prevTextDiv.focus();
    
    // Try to maintain horizontal cursor position
    setTimeout(() => {
        const range = document.createRange();
        const sel = window.getSelection();
        const textLength = prevTextDiv.textContent.length;
        const newPos = Math.min(cursorPos, textLength);
        
        if (prevTextDiv.childNodes.length > 0 && prevTextDiv.childNodes[0]) {
            range.setStart(prevTextDiv.childNodes[0], newPos);
            range.collapse(true);
        } else {
            // For empty content, set cursor at beginning
            range.setStart(prevTextDiv, 0);
            range.collapse(true);
        }
        sel.removeAllRanges();
        sel.addRange(range);
    }, 10);
}

// Handle arrow down navigation
function handleArrowDown(textDiv, blockDiv) {
    const nextBlock = findNextEditableBlock(blockDiv);
    if (!nextBlock) return;
    
    const nextTextDiv = nextBlock.querySelector('.block-text');
    if (!nextTextDiv) return;
    
    // Get current cursor position
    const selection = window.getSelection();
    const range = selection.getRangeAt(0);
    const cursorPos = range.startOffset;
    
    // Save current content before blur
    textDiv.blur();
    
    // Focus next block
    nextTextDiv.contentEditable = 'true';
    nextTextDiv.focus();
    
    // Try to maintain horizontal cursor position
    setTimeout(() => {
        const range = document.createRange();
        const sel = window.getSelection();
        const textLength = nextTextDiv.textContent.length;
        const newPos = Math.min(cursorPos, textLength);
        
        if (nextTextDiv.childNodes.length > 0 && nextTextDiv.childNodes[0]) {
            range.setStart(nextTextDiv.childNodes[0], newPos);
            range.collapse(true);
        } else {
            // For empty content, set cursor at beginning
            range.setStart(nextTextDiv, 0);
            range.collapse(true);
        }
        sel.removeAllRanges();
        sel.addRange(range);
    }, 10);
}

// Find previous editable block (handling nesting)
function findPreviousEditableBlock(blockDiv) {
    // First check previous sibling
    let prev = blockDiv.previousSibling;
    if (prev && prev.classList && prev.classList.contains('block')) {
        // If previous has children, get the last descendant
        let lastChild = prev;
        while (lastChild.querySelector('.block-children > .block:last-child')) {
            lastChild = lastChild.querySelector('.block-children > .block:last-child');
        }
        return lastChild;
    }
    
    // No previous sibling, check parent
    const parentContainer = blockDiv.parentNode;
    if (parentContainer && parentContainer.classList.contains('block-children')) {
        const parentBlock = parentContainer.parentNode;
        if (parentBlock && parentBlock.classList.contains('block')) {
            return parentBlock;
        }
    }
    
    return null;
}

// Find next editable block (handling nesting)
function findNextEditableBlock(blockDiv) {
    // First check if current block has children
    const childContainer = blockDiv.querySelector('.block-children');
    if (childContainer) {
        const firstChild = childContainer.querySelector('.block');
        if (firstChild) return firstChild;
    }
    
    // No children, check next sibling
    let next = blockDiv.nextSibling;
    if (next && next.classList && next.classList.contains('block')) {
        return next;
    }
    
    // No next sibling, traverse up and find parent's next sibling
    let current = blockDiv;
    while (current) {
        const parentContainer = current.parentNode;
        if (parentContainer && parentContainer.classList.contains('block-children')) {
            const parentBlock = parentContainer.parentNode;
            if (parentBlock && parentBlock.classList.contains('block')) {
                const parentNext = parentBlock.nextSibling;
                if (parentNext && parentNext.classList && parentNext.classList.contains('block')) {
                    return parentNext;
                }
                current = parentBlock;
            } else {
                break;
            }
        } else {
            break;
        }
    }
    
    return null;
}

// Keyboard shortcuts
document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
        // Only go back if not editing
        if (!document.querySelector('.block-text.editing')) {
            goBack();
        }
    }
    
    // Cmd/Ctrl + Shift + Left Arrow: Collapse all blocks
    if ((event.metaKey || event.ctrlKey) && event.shiftKey && event.key === 'ArrowLeft') {
        event.preventDefault();
        const allBlocks = document.querySelectorAll('.block-has-children');
        allBlocks.forEach(block => {
            const toggle = block.querySelector('.collapse-toggle');
            const blockId = block.getAttribute('data-block-id');
            const collapseKey = `collapsed_${currentPage}_${blockId}`;
            
            block.classList.add('collapsed');
            if (toggle) toggle.innerHTML = '▶';
            localStorage.setItem(collapseKey, 'true');
        });
    }
    
    // Cmd/Ctrl + Shift + Right Arrow: Expand all blocks
    if ((event.metaKey || event.ctrlKey) && event.shiftKey && event.key === 'ArrowRight') {
        event.preventDefault();
        const allBlocks = document.querySelectorAll('.block-has-children');
        allBlocks.forEach(block => {
            const toggle = block.querySelector('.collapse-toggle');
            const blockId = block.getAttribute('data-block-id');
            const collapseKey = `collapsed_${currentPage}_${blockId}`;

            block.classList.remove('collapsed');
            if (toggle) toggle.innerHTML = '▼';
            localStorage.removeItem(collapseKey);
        });
    }

    // Tab: Toggle collapse/expand for focused block with children
    if (event.key === 'Tab' && !document.querySelector('.block-text.editing')) {
        const focusedElement = document.activeElement;
        const blockDiv = focusedElement?.closest('.block');

        if (blockDiv && blockDiv.classList.contains('block-has-children')) {
            event.preventDefault();
            const toggle = blockDiv.querySelector('.collapse-toggle');
            if (toggle) {
                toggle.click();
            }
        }
    }
});
