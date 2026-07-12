---
title: 'Go + PocketBase + Templ + HTMX: a real SaaS without a JavaScript framework'
date: '2026-06-27'
lang: en
description: 'Experience report: building a full club-management SaaS with Go, PocketBase as an embedded backend, Templ for typed HTML and HTMX for interactivity, 671 hx- attributes, no JavaScript framework, one language, one binary.'
ogDescription: 'Building a full CRUD SaaS with the Go hypermedia stack: PocketBase, Templ, HTMX. What works, what hurts.'
keywords: Go, Golang, HTMX, PocketBase, Templ, hypermedia, SaaS, Tailwind, daisyUI
summary: 'Experience report from a real CRUD SaaS: 671 hx- attributes, typed HTML components, no API layer, no duplicated state, and an honest list of what hurts.'
---

For the last year and a half I have been building and operating a club-management SaaS, members, licenses, grades, competitions, payments, PDF attestations, federation CSV imports. A real application with real users, forms everywhere, and lists that need filtering, sorting and bulk editing. It contains **no JavaScript framework at all**: no React, no Vue, no build pipeline for a SPA, no JSON API layer. This is what the stack looks like, and what I learned.

## The stack, and what each piece actually does

- **[PocketBase](https://pocketbase.io/), used as a Go library**: not as the standalone app most people know. Imported as a framework, it provides the embedded SQLite database, migrations, authentication and session handling, while all routes and business logic are my own Go code. You get the "batteries" of Rails or Django, compiled into a single binary with your app.
- **[Templ](https://templ.guide/)**: HTML components as Go functions, compiled and type-checked. A component takes typed parameters, and a typo in a field name is a compile error, not a blank page in production.
- **[HTMX](https://htmx.org/)**: every interaction is a small HTML fragment swap. The server renders HTML; the browser replaces a node. That's the whole architecture.
- **[hyperscript](https://hyperscript.org/)** for the residual client-side micro-interactions (toggle this, focus that) that don't deserve a server round-trip.
- **Tailwind CSS + daisyUI** for styling, compiled at build time.

One language, one process, one binary to deploy, one SQLite file to back up.

## What "HTMX at scale" actually looks like

People sometimes picture HTMX as a toy for a search box. As I write this, the codebase contains **671 `hx-*` attributes**. Their distribution says a lot about how the pattern scales:

```
171  hx-target      128  hx-swap       90  hx-get
 77  hx-push-url     49  hx-post       34  hx-trigger
 29  hx-include      19  hx-confirm    18  hx-delete
 14  hx-patch        12  hx-on          7  hx-indicator
```

- `hx-target`/`hx-swap` dominate: almost everything is "fetch this fragment, swap that node". Instant search is `hx-get` + `hx-trigger="keyup changed delay:300ms"` swapping a table body.
- `hx-push-url` (77 uses) is the unsung hero: fragment navigation that still updates the URL, so every filtered list and detail view is bookmarkable and shareable. You get SPA-feeling navigation *and* working browser history for free.
- `hx-confirm` + `hx-delete`/`hx-patch` replace an entire modal library: proper HTTP verbs, one attribute for the confirmation.

The pattern that surprised me most in practice: **context preservation**. The same member edit form is reachable from a member sheet, a competition page, or an award list, and after submit or cancel, the user must return where they came from. With hypermedia this is just a parameter carried in the fragment URLs, rendered server-side into the form's `hx-post` target. No client router, no state store, the state *is* the HTML.

## What it removes

The honest pitch for this stack is not what it adds but what it deletes:

- **No API layer.** No JSON serializers, no DTOs, no versioning, no frontend/backend contract drift. The server renders what the user sees, full stop.
- **No duplicated state.** The eternal SPA problem, cache invalidation between client store and server, cannot exist, because there is no client store.
- **No split toolchain.** One `go build` (plus `templ generate` and the Tailwind CLI, both fast and boring). CI is short. Dependency updates are Go module updates.
- **Compile-time UI safety.** Renaming a model field breaks every template that uses it, at compile time. Refactoring a server-rendered UI in a typed language is dramatically calmer than grepping JSX.

## What hurts (honesty section)

- **You write your own components.** daisyUI covers the CSS, but there is no ecosystem of ready-made data grids or date pickers wired for HTMX. For a CRUD app that trade is fine; for a Figma-grade interactive product it would not be.
- **Some JavaScript survives.** Rich interactions (a box-drawing selector, drag-and-drop) still need real JS islands. The discipline is keeping them islands.
- **PocketBase coupling.** Using it as a library means its upgrade cadence is your framework's upgrade cadence, and major version bumps have real migration cost.
- **Fragment discipline.** As fragments multiply, you need conventions for what returns a full page versus a partial (the `HX-Request` header is your friend), or the routing layer gets messy.

## Verdict

For the vast category of software that is actually built in the world, admin tools, B2B SaaS, dashboards, back-offices, this stack is a quiet superpower for a solo developer or a small team. The entire application is one Go program you can hold in your head; the frontend cannot drift from the backend because there is no frontend in the SPA sense; and the user experience, thanks to HTMX's targeted swaps, feels every bit as responsive as the React equivalent that would have taken three times as long to build.

The hypermedia approach is not a nostalgia trip. It is what removing accidental complexity looks like.
