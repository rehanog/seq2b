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

<p>But somewhere along the way, Logseq lost the plot. What started as a focused outliner began sprouting features like a startup trying to justify its Series B funding. Database integration, complex queries, plugins for everything, AI features that felt bolted on rather than thoughtfully integrated. Each release seemed to drift further from the core promise: a simple, fast way to capture and connect your thoughts.</p>

<p>Performance suffered. The desktop app, already burdened by Electron, became sluggish. Pages with a few hundred blocks would stutter. Search felt like swimming through treacle. The very tool meant to capture the speed of thought was actively hampering it. And then there were the sync issues—files getting corrupted, changes lost between devices, the sort of data integrity problems that make you question whether you can trust your second brain to remember anything at all.</p>

<p>Support, meanwhile, seemed to evaporate. Bug reports would sit for months. Basic functionality would break between versions. The community grew increasingly frustrated, and the developers increasingly distant. It felt like watching a promising project slowly suffocate under its own ambitions.</p>

<p>I looked elsewhere, naturally. Obsidian was fast but felt too rigid. Roam was clever but lived in the cloud. Silver Bullet caught my attention—lovely interface, proper Markdown support, and genuinely innovative ideas about extensibility. But there was a catch. Silver Bullet relies on browser storage (IndexedDB) for offline operation. Your notes live in the browser's database, syncing back to a server when connected. It's clever engineering, but it felt wrong for something as important as a second brain. Browsers crash, storage gets cleared, and I've lost enough work to Chrome's occasional tantrums to trust it with my life's work. Call me old-fashioned, but I want my notes on the file system, where I can see them, back them up, and know they'll outlive whatever application I happen to be using this decade.</p>

<p>I'd been toying with the idea of writing my own replacement for months. The usual programmer's hubris: "How hard could it be?" But between work, life, and the general entropy of existence, it remained a vague someday project. Then Claude Code arrived, and I started seq2b as a toy project—a way to play with this new tool that promised to accelerate development. Write a simple Markdown parser, they said. Parse some blocks, maybe implement basic linking. A weekend experiment at most.</p>

<p>Two days later, I had a working CLI tool, a desktop GUI with proper block indentation, bidirectional linking, and HTML export. Two days. The velocity was startling. Not just the code generation, but the thoughtful architecture discussions, the test writing, the debugging sessions. It was like having a very competent junior developer who never got tired and occasionally suggested better approaches than the one I'd planned.</p>

<p>For the first time in twenty years—literally—I felt the joy of coding again. My day job has been more managerial for a long time, and while there's satisfaction in that work, it's not the same as staying up late watching something grow line by line. But here I was, past midnight, adding features and fixing bugs with the sort of focused enthusiasm I thought I'd left behind in my twenties. Except now, with these new tools, things were happening at light speed. Ideas could be tested, implemented, and refined in hours rather than weeks.</p>

<p>So here we are. Another note-taking app, because apparently the world needed one more. But this one has a different philosophy: performance first (built in Go, not Electron), files not databases (your notes as Markdown files on your file system), focus over features (an outliner with linking, no bells, no whistles, no blockchain integration), and most importantly, data integrity above all else. Your thoughts are precious cargo; they deserve better than to be casualties of feature creep or sync conflicts.</p>

<p>Will it succeed? Probably not in any conventional sense. The note-taking space is crowded, and most people are perfectly happy with their current tools. But it exists for the same reason people build custom keyboards or grow their own tomatoes: because sometimes the thing you want doesn't exist, and making it yourself is both practical and oddly satisfying. Besides, if nothing else, it's been an excellent excuse to thoroughly test Claude Code's capabilities. The future of programming might be more conversational than we imagined.</p>

</div>