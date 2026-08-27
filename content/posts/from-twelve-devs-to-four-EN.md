---
title: 'From twelve devs to four: the LLM quality drop that never came'
date: '2026-08-27'
lang: en
description: 'Everyone predicted AI-generated code would tank software quality. From inside a team that went from twelve devs to four, the honest assessment: for most of what we build, the generated code is as good or better, and that is exactly why it will spread everywhere.'
ogDescription: 'A team that went from twelve devs to four. Quality held. That is the problem, and why AI code will spread everywhere.'
keywords: LLM, AI, software quality, code generation, team size, developer jobs, productivity
summary: 'My team went from twelve to one lead, one dev, a data integrator, and an apprentice, through departures simply not replaced. Except for a few architecture-heavy pipelines, the AI-generated code is as good as or better than what an average dev writes. Devs will remain, like farm workers remained, but LLMs are the tractor.'
---

There is a popular storyline about LLMs and code: quality is collapsing, technical debt is piling up everywhere, and in a few years someone will have to pay for all this generated slop. Maybe. But from where I sit, the story looks different, and less comfortable.

## The numbers first

Where I work, the dev team used to be twelve people, including two senior leads. Today we are one lead, one dev, one data integrator, and one apprentice. There was no layoff plan, no dramatic all-hands. It happened progressively: people left, and they simply were not replaced. The workload did not disappear. AI helps.

## The quality question

Let's be honest about what most of our code actually is. Yes, we have a few genuinely touchy processes that demand solid architecture skills: the [speech-to-text pipeline](running-whisper-24-7-on-broadcast-streams-EN.html), the [OCR](ocr-llm-enrichment-broadcast-transcripts-EN.html), the [24/7 video capture platform](recording-broadcast-24-7-capture-platform-architecture-EN.html), and the audio recognition platform. Those still need someone who knows what they are doing, and the LLM is an assistant there, not an author.

But the bulk of what we ship is not that. It is data entry and data presentation tools: forms, tables, imports, exports, dashboards. And on that kind of code, comparing what the LLM generates to what an average dev writes, the generated code is often simply better. Not better than what our two former senior leads would have written. Better than the realistic alternative.

That is the part the quality-collapse storyline gets wrong. It compares AI code to some ideal artisan codebase. The honest comparison is AI code versus average-dev code under deadline, and the AI wins that one more often than we like to admit.

## The managerial view

Now put on the manager's hat, and again, let's not kid ourselves. If the overall quality is the same and the team is a third of the size, the costs drop massively. A company that does this can undercut a company that does not. Our competitors face the same arithmetic. So this is not a technology adoption question anymore, it is basic competition: AI-generated code will spread everywhere, not because anyone loves it, but because nobody can afford to be the last one paying twelve salaries for what four people plus an LLM deliver.

## Tractors

Farm workers did not disappear. There are still people driving the machines, fixing them, deciding what to plant and when. But the tractor divided their numbers by an order of magnitude, and no amount of debate about the craft of hand plowing changed that.

There will always be devs. Someone has to own the architecture, review what ships, handle the touchy pipelines, and know when the generated code is wrong. But today's LLMs are the tractor, and I am watching the head count adjust accordingly, one unreplaced departure at a time.
