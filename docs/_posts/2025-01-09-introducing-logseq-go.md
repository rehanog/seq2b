---
layout: default
title: "Introducing Logseq Go: A High-Performance Knowledge Management System"
date: 2025-01-09
author: Rehan
categories: announcement
---

# Introducing Logseq Go: A High-Performance Knowledge Management System

<div class="blog-post-meta">
  January 9, 2025 • Rehan
</div>

<div class="container" style="padding: 2rem 0;">

I'm excited to announce **Logseq Go**, a new high-performance knowledge management system built from the ground up in Go. This project represents months of work solving real problems I encountered with existing tools.

## The Journey

It started with a simple need: I wanted a fast, native application for managing my knowledge base. Existing solutions were either too slow (Electron-based), too complex, or lacked the features I needed. So I decided to build my own.

### The Indentation Problem

The journey wasn't smooth. My first attempt using Fyne for the GUI hit a major roadblock - proper text indentation for nested blocks simply wasn't possible with their RichText widget. After much frustration and trying various workarounds, I realized I needed a different approach.

### Enter Wails

That's when I discovered [Wails](https://wails.io/), which combines the power of Go backends with web-based frontends in a native window. This gave me the best of both worlds:
- Native performance from Go
- Full control over text rendering with HTML/CSS
- Cross-platform compatibility
- Small binary size

## What Makes Logseq Go Special?

### 🚀 Performance First
Written in Go, Logseq Go can parse thousands of markdown files in seconds. No more waiting for your knowledge base to load.

### 🔗 Smart Bidirectional Linking
Automatic backlink detection means you never lose track of connections between your ideas. The relationship graph builds itself as you write.

### 📱 Mobile-Ready Architecture
While we're starting with desktop, the architecture is designed for mobile from day one. The parser is a shared library that can be used across platforms.

### 🎯 Clean, Focused Interface
Proper block indentation with visual hierarchy helps you see the structure of your thoughts at a glance. No clutter, just your content.

## Current Status

Today, Logseq Go includes:
- ✅ Full markdown parser with Logseq-compatible block structure
- ✅ Native desktop GUI for macOS, Windows, and Linux
- ✅ CLI tool for automation and scripting
- ✅ Automatic backlink detection and navigation
- ✅ Clean visual hierarchy with proper indentation

## What's Next?

The roadmap is ambitious but achievable:

**Phase 3** (Next): Advanced parsing features like properties, tags, and TODO states

**Phase 4**: Persistent storage layer with BadgerDB for lightning-fast queries

**Phase 5**: Git/JJ integration for version control and sync

**Phase 6**: Security features including signed binaries and sandboxed plugins

**Phase 7**: AI integration with local LLMs and semantic search

**Phase 8**: Full API and optional web interface

## Open Source from Day One

Logseq Go is MIT licensed and open source. Every line of code is available on [GitHub](https://github.com/rehan/logseq-go). I believe in building in public and welcome contributions from the community.

## Try It Today

You can download and try Logseq Go right now:

```bash
git clone https://github.com/rehan/logseq-go.git
cd logseq-go

# Try the CLI
go run cmd/logseq-go/main.go testdata/pages

# Or launch the GUI
cd desktop/wails
wails dev
```

## Join the Journey

This is just the beginning. I'm building Logseq Go to be the knowledge management system I've always wanted, and I hope it can be that for you too.

- Star the project on [GitHub](https://github.com/rehan/logseq-go)
- Report issues or suggest features
- Contribute code or documentation
- Share your feedback and ideas

Together, we can build something amazing.

---

*What features would you like to see in a knowledge management system? Let me know in the [discussions](https://github.com/rehan/logseq-go/discussions)!*

</div>