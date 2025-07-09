---
layout: default
title: Home
---

<section class="hero">
    <div class="container">
        <h1>seq2b</h1>
        <p class="tagline">High-performance knowledge management built with Go</p>
        <div class="hero-buttons">
            <a href="{{ '/download' | relative_url }}" class="btn btn-primary">Download Now</a>
            <a href="https://github.com/rehanog/seq2b" class="btn btn-secondary">View on GitHub</a>
        </div>
    </div>
</section>

<section class="features">
    <div class="container">
        <h2>Why seq2b?</h2>
        <div class="feature-grid">
            <div class="feature-card">
                <div class="feature-icon">🚀</div>
                <h3>Blazing Fast</h3>
                <p>Built with Go for exceptional performance. Parse thousands of pages in seconds, not minutes.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🔗</div>
                <h3>Smart Linking</h3>
                <p>Automatic bidirectional linking with real-time backlink detection. Never lose track of connections.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">📱</div>
                <h3>Cross-Platform</h3>
                <p>Native desktop app today, mobile apps tomorrow. One codebase, all your devices.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🎯</div>
                <h3>Clean Interface</h3>
                <p>Proper block indentation and visual hierarchy. Focus on your thoughts, not the tool.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🛡️</div>
                <h3>Secure by Design</h3>
                <p>Coming soon: Signed binaries and sandboxed plugins. Your data stays yours.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🤖</div>
                <h3>AI Ready</h3>
                <p>Future-proof architecture for AI integration. Local LLMs and semantic search on the roadmap.</p>
            </div>
        </div>
    </div>
</section>

<section class="demo">
    <div class="container">
        <h2>See It In Action</h2>
        <div style="text-align: center; margin: 2rem 0;">
            <img src="{{ '/assets/images/screenshot.png' | relative_url }}" alt="Seq2B Screenshot" style="max-width: 100%; border-radius: 12px; box-shadow: 0 10px 30px rgba(0,0,0,0.1);">
        </div>
    </div>
</section>

<section class="getting-started">
    <div class="container">
        <h2>Quick Start</h2>
        <div style="background-color: #f8f9fa; padding: 2rem; border-radius: 12px;">
            <h3>CLI Tool</h3>
            <pre><code># Parse a directory of markdown files
go run cmd/seq2b/main.go path/to/your/pages/

# See block structure, backlinks, and more
./seq2b-cli your-notes.md</code></pre>
            
            <h3>Desktop GUI</h3>
            <pre><code># Launch the desktop application
cd desktop/wails
wails dev

# Or build for production
wails build</code></pre>
        </div>
    </div>
</section>

<section class="roadmap">
    <div class="container">
        <h2>Roadmap</h2>
        <div style="max-width: 800px; margin: 0 auto;">
            <div style="display: flex; align-items: center; margin: 1rem 0;">
                <span style="color: #28a745; font-size: 1.5rem; margin-right: 1rem;">✓</span>
                <div>
                    <strong>Phase 1 & 2: Core Parser & Desktop GUI</strong>
                    <p style="margin: 0; color: #666;">Fully functional desktop app with backlinks</p>
                </div>
            </div>
            <div style="display: flex; align-items: center; margin: 1rem 0;">
                <span style="color: #ffc107; font-size: 1.5rem; margin-right: 1rem;">◐</span>
                <div>
                    <strong>Phase 3: Advanced Parsing</strong>
                    <p style="margin: 0; color: #666;">Properties, tags, TODO states</p>
                </div>
            </div>
            <div style="display: flex; align-items: center; margin: 1rem 0;">
                <span style="color: #6c757d; font-size: 1.5rem; margin-right: 1rem;">○</span>
                <div>
                    <strong>Phase 4-8: Storage, Sync, AI & More</strong>
                    <p style="margin: 0; color: #666;">Persistent storage, Git/JJ sync, AI integration</p>
                </div>
            </div>
        </div>
    </div>
</section>