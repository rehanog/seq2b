---
layout: default
title: "Why Another Note-Taking App?"
date: 2025-01-11
author: Rehan
categories: reflection
---

<h1>Why Another Note-Taking App?</h1>

<div class="blog-post-meta">
  January 11, 2025 • Rehan
</div>

<div class="container" style="padding: 2rem 0;">

<p>I've been using Logseq for the better part of two years. Not because it was perfect, but because it got two fundamental things right: the outliner approach to thinking, and the local-file-first philosophy. Your thoughts, stored as plain text files on your own machine. Radical concept, apparently.</p>

<p>But somewhere along the way, Logseq lost the plot.</p>

<h2>The Feature Creep Problem</h2>

<p>What started as a focused outliner began sprouting features like a startup trying to justify its Series B funding. Database integration, complex queries, plugins for everything, AI features that felt bolted on rather than thoughtfully integrated. Each release seemed to drift further from the core promise: a simple, fast way to capture and connect your thoughts.</p>

<p>Performance suffered. The desktop app, already burdened by Electron, became sluggish. Pages with a few hundred blocks would stutter. Search felt like swimming through treacle. The very tool meant to capture the speed of thought was actively hampering it.</p>

<h2>The Alternatives Weren't Quite Right</h2>

<p>I looked elsewhere, naturally. Obsidian was fast but felt too rigid. Roam was clever but lived in the cloud. Silver Bullet caught my attention—lovely interface, proper Markdown support, and genuinely innovative ideas about extensibility.</p>

<p>But there was a catch. Silver Bullet relies on browser storage (IndexedDB) for offline operation. Your notes live in the browser's database, syncing back to a server when connected. It's clever engineering, but it felt wrong for something as important as a second brain. Browsers crash, storage gets cleared, and I've lost enough work to Chrome's occasional tantrums to trust it with my life's work.</p>

<p>Call me old-fashioned, but I want my notes on the file system, where I can see them, back them up, and know they'll outlive whatever application I happen to be using this decade.</p>

<h2>The Claude Code Catalyst</h2>

<p>I'd been toying with the idea of writing my own replacement for months. The usual programmer's hubris: "How hard could it be?" But between work, life, and the general entropy of existence, it remained a vague someday project.</p>

<p>Then Claude Code arrived.</p>

<p>I started seq2b as a toy project—a way to play with this new tool that promised to accelerate development. Write a simple Markdown parser, they said. Parse some blocks, maybe implement basic linking. A weekend experiment at most.</p>

<p>Three weeks later, I had a working CLI tool, a desktop GUI with proper block indentation, bidirectional linking, and HTML export. The velocity was startling. Not just the code generation, but the thoughtful architecture discussions, the test writing, the debugging sessions. It was like having a very competent junior developer who never got tired and occasionally suggested better approaches than the one I'd planned.</p>

<h2>Why seq2b</h2>

<p>So here we are. Another note-taking app, because apparently the world needed one more. But this one has a different philosophy:</p>

<ul>
<li><strong>Performance first</strong>: Built in Go, not Electron. No excuses for sluggish performance.</li>
<li><strong>Files, not databases</strong>: Your notes as Markdown files on your file system. Own your data completely.</li>
<li><strong>Focus over features</strong>: An outliner with linking. No bells, no whistles, no blockchain integration.</li>
<li><strong>Cross-platform</strong>: Native desktop today, mobile tomorrow. Same parser, different interfaces.</li>
</ul>

<p>Will it succeed? Probably not in any conventional sense. The note-taking space is crowded, and most people are perfectly happy with their current tools. But it exists for the same reason people build custom keyboards or grow their own tomatoes: because sometimes the thing you want doesn't exist, and making it yourself is both practical and oddly satisfying.</p>

<p>Besides, if nothing else, it's been an excellent excuse to thoroughly test Claude Code's capabilities. The future of programming might be more conversational than we imagined.</p>

</div>