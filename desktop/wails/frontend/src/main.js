import './style.css';
import { GetPage, GetPageList } from '../wailsjs/go/main/App';

// Application state
let currentPage = 'Page A';
let navigationHistory = [];

// DOM elements
const pageTitle = document.getElementById('pageTitle');
const backButton = document.getElementById('backButton');
const blocksContainer = document.getElementById('blocks');
const backlinksContainer = document.getElementById('backlinksList');
const loadingDiv = document.getElementById('loading');

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    loadPage(currentPage);
    setupEventListeners();
});

// Event listeners
function setupEventListeners() {
    backButton.addEventListener('click', goBack);
}

// Load and display a page
async function loadPage(pageName) {
    try {
        loadingDiv.style.display = 'block';
        blocksContainer.innerHTML = '';
        backlinksContainer.innerHTML = '';
        
        const pageData = await GetPage(pageName);
        
        // Update page title
        pageTitle.textContent = pageData.title;
        
        // Render blocks
        renderBlocks(pageData.blocks, blocksContainer);
        
        // Render backlinks
        renderBacklinks(pageData.backlinks);
        
        // Update navigation
        if (currentPage !== pageName) {
            navigationHistory.push(currentPage);
            backButton.disabled = false;
        }
        
        currentPage = pageName;
        loadingDiv.style.display = 'none';
        
    } catch (error) {
        console.error('Error loading page:', error);
        loadingDiv.innerHTML = `Error loading page: ${error.message || error}`;
    }
}

// Render blocks recursively
function renderBlocks(blocks, container) {
    blocks.forEach(block => {
        const blockElement = createBlockElement(block);
        container.appendChild(blockElement);
    });
}

// Create a block element
function createBlockElement(block) {
    const blockDiv = document.createElement('div');
    blockDiv.className = 'block';
    blockDiv.setAttribute('data-depth', block.depth);
    blockDiv.setAttribute('data-block-id', block.id);
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'block-content';
    
    const bullet = document.createElement('span');
    bullet.className = 'block-bullet';
    bullet.textContent = '•';
    
    const textDiv = document.createElement('div');
    textDiv.className = 'block-text';
    textDiv.innerHTML = processLinksInHTML(block.htmlContent);
    
    contentDiv.appendChild(bullet);
    contentDiv.appendChild(textDiv);
    blockDiv.appendChild(contentDiv);
    
    // Add children if any
    if (block.children && block.children.length > 0) {
        renderBlocks(block.children, blockDiv);
    }
    
    return blockDiv;
}

// Process HTML content to make links clickable
function processLinksInHTML(html) {
    // Convert <a href="page">text</a> to clickable links
    return html.replace(/<a href="([^"]+)">([^<]+)<\/a>/g, (match, href, text) => {
        return `<a href="#" class="page-link" onclick="navigateToPage('${text}')">${text}</a>`;
    });
}

// Navigate to a page
window.navigateToPage = function(pageName) {
    loadPage(pageName);
};

// Go back to previous page
function goBack() {
    if (navigationHistory.length > 0) {
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

// Render backlinks
function renderBacklinks(backlinks) {
    if (!backlinks || backlinks.length === 0) {
        backlinksContainer.innerHTML = '<div class="empty-state">No backlinks found</div>';
        return;
    }
    
    backlinks.forEach(backlink => {
        const backlinkDiv = document.createElement('div');
        backlinkDiv.className = 'backlink-item';
        
        const sourceDiv = document.createElement('div');
        sourceDiv.className = 'backlink-source';
        sourceDiv.textContent = backlink.sourcePage;
        sourceDiv.onclick = () => navigateToPage(backlink.sourcePage);
        
        const blocksDiv = document.createElement('div');
        blocksDiv.className = 'backlink-blocks';
        blocksDiv.textContent = `Referenced in: ${backlink.blockIds.join(', ')}`;
        
        backlinkDiv.appendChild(sourceDiv);
        backlinkDiv.appendChild(blocksDiv);
        backlinksContainer.appendChild(backlinkDiv);
    });
}

// Error handling
window.addEventListener('error', (event) => {
    console.error('Application error:', event.error);
});

// Keyboard shortcuts
document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
        goBack();
    }
});
