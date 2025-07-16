import { vi } from 'vitest';

// Mock Wails runtime APIs
window.go = {
  main: {
    App: {
      GetPage: vi.fn(),
      UpdateBlockAtPath: vi.fn(),
      AddBlockAtPath: vi.fn(),
      GetPageList: vi.fn(),
      GetBacklinks: vi.fn(),
      UpdateBlock: vi.fn(),
      AddBlock: vi.fn(),
      LoadDirectory: vi.fn(),
      RefreshPages: vi.fn()
    }
  }
};

// Mock DOM elements that main.js expects
document.body.innerHTML = `
  <div id="pageTitle"></div>
  <button id="backButton" disabled></button>
  <button id="homeButton"></button>
  <div id="blocks"></div>
  <div id="backlinksList"></div>
  <div id="loading" style="display: none;"></div>
`;

// Note: Using real timers for async tests