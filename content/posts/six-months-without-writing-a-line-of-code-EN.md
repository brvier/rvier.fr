---
title: 'Six months without writing a line of code'
date: '2026-11-23'
lang: en
translation: six-mois-sans-ecrire-une-ligne-de-code-FR
description: 'For six months I have not written code myself, except for the fun of a puzzle or an assembly debugging session. How I got there, from line completion to an agent that I hand the architecture to, the habits LLMs made me pick up that actually improve the code, and what they now do with my tickets, my email and Slack.'
ogDescription: 'From line completion to an agent coding under my architecture: six months without typing code, and the code is better than before.'
image: https://rvier.fr/images/six-months-without-code.png
keywords: LLM, AI, coding agents, code completion, software quality, tests, changelog, refactoring, software architecture, support, productivity, developer
summary: 'Six months without typing a line of code, puzzles and assembly aside. What I handed over little by little, from line completion to an agent coding under my architecture, the habits that improved my code in three months, and the same principle applied to support tickets, email and Slack.'
---

I have not written a line of code in six months. Well, almost: I still write some for fun, when I feel like a puzzle, or to debug assembly. Everything else, at work and on my personal projects, is typed by an LLM.

<img src="../images/six-months-without-code.png" alt="A computer keyboard seen from above with every key dark except the space bar, lit orange under a finger and surrounded by dictation waves" loading="lazy" width="1200" height="627">

I did not get there out of conviction, and not overnight either. Little by little, it simply became much faster and more efficient than typing it myself. Every time the models cleared a new bar, and they have changed enormously, I handed over a little more.

## I hand over the end of the line

At first, the LLM completed the line I was on. I started it, it guessed the end. It was glorified autocomplete, but already effective: most of the time it predicted exactly what I was about to write.

## I hand over the next lines

Then it started suggesting the lines that came after, and that is when I began saving time for real.

I write a SQL query: it suggests the Go struct that receives the result, with the right fields and the right types, based on the query just above. I paste the JSON response of an API: it derives the matching struct, tags included. Nothing hard in that code. Just dozens of tedious lines where you make mistakes out of inattention.

## I hand over the function

Next, I stopped starting the lines. I described a small function, exactly: what it takes, what it returns, the edge cases. The LLM wrote it, I reviewed it.

The important word is "exactly". At that stage I was still the author: I knew what code I wanted, I just had someone else type it.

## I hand over the code, I keep the architecture

Today I no longer describe functions. I give the code architecture and the complete software architecture, and the LLM writes the rest: the functions, the tests, the migrations, the interface.

I give the architecture because that is where LLMs still make a lot of mistakes. Left alone, they pile things up: one abstraction too many, a useless dependency, two ways of doing the same thing three files apart. So my projects start with a document: the context, numbered principles, what is forbidden. "Code decides, the model writes." "Nothing goes out without a gesture from the user." I already wrote about this in [house conventions turned into skills](encoding-house-conventions-into-ai-agent-skills-EN): what is not written down does not exist for an agent.

That said, there is less and less to write down. At first I used a huge number of skills to help the LLM make its decisions. Today I keep only two or three, and they are business preferences: how to write changelogs and commit messages, plus two or three deployment architecture choices, Ansible for instance. For the rest, it no longer needs its hand held.

As for the architecture mistakes, they will not last. In two years this will surely no longer be the case. Maybe even in six months.

## Three months of habits that improve the code

I already wrote that on my team [quality did not drop](from-twelve-devs-to-four-EN). I will go further: for the past three months, my code has been better than before. Not because the model writes better than I do. Because everything I knew I should do, and did not do for lack of time, now costs nothing.

**A bug fixed, a test that keeps it from coming back.** A recent example: a timeline whose markers all piled up on the left edge of the page. The server's security policy was blocking `style` attributes written in the HTML, silently. The fix is ten lines. The test shipped with it walks through every page and every script and rejects any inline style: that bug will not come back, and neither will its cousins. My latest project has close to a hundred test files, one for every two code files. I have never had that ratio in my life.

**The commit does not go through if the tests do not pass.** A pre-commit hook with `gofmt`, `go vet` and the tests: everybody knows it should be done, nobody takes the hour to set it up. Now it is one sentence. And it serves the agent as much as it serves me: an LLM that can run the tests fixes its own mistakes, an LLM that cannot "verifies" by re-reading its own diff.

**A changelog that says what changed.** My commit messages often looked like "fix", "wip", "fix again". Today they describe the behavior: "a read abandoned by the client returns 499 instead of a 500". And the changelog is updated in the same commit. I do not write them, but I read them, and that has become my first code review: if the line does not describe what I asked for, something went sideways.

**The refactorings we kept postponing.** On [OhayōDojo](go-pocketbase-templ-htmx-EN), the "someday we should" list was emptied in a few days: dead code removed, a single amount formatter for pages, messages, PDFs and CSVs, shared plumbing for the three importers, centralized access control. Nothing hard, just long and boring, so never a priority. When a refactoring costs twenty minutes of review instead of two days of typing, it gets done.

**Throwaway tools to verify.** A fake LLM server and a fake IMAP server to test offline, a command that takes screenshots of the interface, a `--dry-run` on every scheduled job. I almost never built that kind of tooling, too expensive for what it returns. It now comes with the project.

**Throwing away without regret.** In early October I started a personal assistant: 34,000 lines of Go in a month. In use, the starting idea turned out to be wrong, and at the end of October I restarted from scratch. With 34,000 hand-written lines, I would have defended my architecture out of sheer pride. When code costs almost nothing, you stop getting attached to it, and you make better decisions.

## The bulk of the work: taking over legacy code

None of this happened on fresh, clean projects. The bulk of my work these past months has been taking over legacy code in a pretty sorry state.

The reasons are the usual ones. Several junior developers worked on it, then several contractors. The specifications were vague, or missing altogether. And the architecture was a problem everywhere: a multitude of micro-programs, each in charge of one task, each with its own spaghetti code. Several of them read and modified the same databases, the same tables, and they all depended on one another.

In short, an architecture picked because it was fashionable, by developers who did not really know what they were doing.

That is the job I carried out without typing code.

## And not only code

The same shift happened outside the editor. I now use an LLM for support tickets: it sorts them, dispatches 80% of the answers and drafts them. All I have left to do is approve. Same for email, Slack questions and requests from management.

It also writes everything that surrounds the code. It summarizes changes in commit messages and changelogs. It prepares my daily review for the daily meeting, from what was actually done the day before rather than from what I remember. And it produces documentation and release notes for users who are not technical: the same change, told without a function name or a migration number. That translation, I almost never did, and yet they are the ones using the software.

It is the same pattern as for code, applied elsewhere: the machine prepares, I decide. Nothing goes out before I have read it.

## I review everything, I correct less and less

Right now I still review all the generated code, and the tests generated with it. But I need to ask for changes less and less. Most of the time I read and I approve, without changing anything.

The logical consequence: I now use Claude Code or OpenCode in auto mode. The agent chains its actions without asking my permission at every step, and I review the result at the end.

One exception: when I hand it a large architecture, I switch to plan mode. I describe a very precise plan, and it only starts coding once that plan is settled.

And I am not necessarily talking about the most advanced models, like Opus 5.5: this is also true with more modest ones. To sort my email, for instance, the job is done by a 27-billion-parameter Qwen 3.8, served by vLLM on my own machine.

## What is still my job

I no longer code, but I have never done so much development. The work has simply moved.

Today I mostly spend my time bringing order to user requests, which arrive rough and jumbled, all the more so since there is no project manager. My role looks more and more like a product owner's, or rather a mix: project manager, software architect, a bit of team management, and still lead dev, to guide the teams and stop the architectural nonsense their LLM sometimes generates.

Because it is quite surprising to watch: an LLM is a powerful tool, and even the most powerful ones, like Fable, can be used very, very badly. In some hands they produce absurd results: disgusting architectures, filthy code that is unstable, impossible to evolve and poorly documented.

So I decide what gets built and what does not. I hold the architecture. I read the diffs, the changelogs and the tests. I say no, often: no to one dependency too many, no to the premature abstraction, no to the feature nobody asked for.

None of the habits in this list is new. They are the ones every book has recommended for twenty years and almost nobody keeps, because time always runs short. The LLM did not invent rigor, it made it free.

## By the way, this post

This post was written by an LLM. The outline, I made by hand. The facts, I dictated. But the sentences are its own.

I dictated all of it while walking on my treadmill, one finger on my laptop's space bar to talk to it. Of the whole keyboard, that is the only key that got used.

And when I miss the rest of the keyboard, there is still assembly.
