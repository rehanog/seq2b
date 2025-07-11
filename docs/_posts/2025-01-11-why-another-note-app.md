---
layout: default
title: "Yet Another PKM app?"
date: 2025-01-11
author: Rehan
categories: reflection
---

# Yet Another PKM app?

<div class="blog-post-meta">
  January 11, 2025 • Rehan
</div>

I've been using Logseq for the better part of two years. Not because it was perfect, but because it got two fundamental things right: the outliner approach to thinking, and the local-file-first philosophy. Your thoughts, stored as plain text files on your own machine.

But somewhere along the way, Logseq lost the plot. What started as a focused outliner began sprouting features like a startup trying to justify its Series B funding. Instead of concentrating on speed, stability and fixing the numerous sync / data loss issues, the team focused on adding ever more features of questionable value, and then really went off the rails with a massive a pivot away from one of their core values to build a database version. Their rationale was that they needed that in order to produce a collaborative editing feature that most of their users don't need or want. And worst of all, they've gone super quiet and non-communicative while doing so. Quite sad for a team that initially built a lot of goodwill by releasing very regularly and being very transparent.

I have been looking for anything else that would satisfy my fairly simple requirements. Obsidian is not an outliner. Roam is not local. The most promising alternative seems to be Silver Bullet - the philosophy and vibe of the project seem great. But I just can't bring myself to trust browser storage (IndexedDB) for offline operation. Your notes live in the browser's database, syncing back to a server when connected. It's clever engineering, but it feels wrong for something as important as a second brain.

I've been toying with the idea of writing my own replacement for the better part of a year. But between work, life, and the general entropy of existence, it remained a vague someday project. Then Claude Code arrived, and I started seq2b as a toy project—a way to play with this new tool that promised to accelerate development.

Three days later, and I've got a desktop GUI with block indentation, bidirectional linking. Three days, just doing it in my spare time. The velocity has been bonkers. Got to say Claude Code has kind of blown my mind. It feels like I've jumped into the 22nd Century. It's like having a super competent junior developer who never gets tired and has a freakishly good knowledge about everything.

For the first time in twenty years—literally—I feel the joy of coding again. My day job has been more managerial for a long time, and while there's satisfaction in that work, it's not the same as staying up late watching something grow line by line. But I've been up past midnight, two nights in a row, adding features and fixing bugs with the sort of focused enthusiasm I thought I'd left behind in my twenties. Except now, with these new tools, things were happening at light speed. Ideas could be tested, implemented, and refined in hours rather than weeks.

So here we are. Another note-taking app, because apparently the world needed one more. But this one has a different philosophy: performance first (built in Go, not Electron), files not databases (your notes as Markdown files on your file system), focus over features (an outliner with linking, no bells, no whistles, no blockchain integration), and most importantly, data integrity above all else. Well actually the same philosophy that I thought Logseq had, but that they eem to have abandoned. Our thoughts are precious cargo; they deserve better than to be casualties of feature creep or sync conflicts.
