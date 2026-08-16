---
title: 'Encoding house conventions into AI agent skills'
date: '2026-08-15'
lang: en
description: 'How I make coding agents produce code that looks like ours: per-repo AGENTS.md files for project facts, and reusable skills that capture house conventions, from a French conventional-commit workflow to a 200-line skill that scaffolds a complete Go microservice.'
ogDescription: 'Per-repo AGENTS.md for project facts, reusable skills for house conventions: making coding agents write code that looks like ours.'
keywords: AI agents, Claude Code, skills, AGENTS.md, conventions, Go, developer experience, LLM
image: https://rvier.fr/images/agent-skills-conventions.png
summary: 'Coding agents write plausible but foreign code by default. Per-repo AGENTS.md files and reusable skills, including one that scaffolds an entire Go microservice to house conventions, fix that.'
---

A coding agent with no context writes *plausible* code: idiomatic for the language, aligned with whatever the training distribution considers normal. The problem is that a codebase with fifteen years of history is not the training distribution. We have a house style: a specific config library, a specific logging wrapper, a specific way every service boots, lints, versions, ships. Code that ignores it is not wrong, it is *foreign*, and foreign code is expensive to review and expensive to live with.

<img src="../images/agent-skills-conventions.png" alt="AGENTS.md and skills feeding a coding agent that produces uniform services" loading="lazy" width="1200" height="627">

So the question that actually matters with coding agents is not "can it code?" but "can it code *like us*?". After a year of daily use, my answer is two mechanisms with a clean separation of concerns: per-repository context files, and reusable skills.

## AGENTS.md: the facts of one repository

Every active repo gets an `AGENTS.md` (or `CLAUDE.md`): the file agents read before touching anything. Mine converged on three ingredients:

- **The commands.** How to build, test, lint, run one single test. An agent that knows `flutter test --name "pattern"` exists will run it; one that does not will "verify" by re-reading its own diff.
- **The architecture in ten lines.** For my Flutter app Planova: clean architecture, Provider for state, structured Markdown files as the only storage. Enough to keep a change in the right layer.
- **The feature inventory.** The OhayōDojo file lists every implemented feature, from HTMX search behavior to which entities support inline editing. This is the part people skip, and it is the part that prevents the two classic agent failures: re-implementing something that exists, and "simplifying away" a behavior that users depend on.

The discipline is keeping it *facts only*. No aspirations, no roadmap, no style essays: those age badly, and an agent trusts a stale document more confidently than a human would.

## Skills: workflows that cross repositories

Repo files cannot hold what is shared across thirty repos. For that I use skills: named, self-contained instruction files an agent loads when the task matches. Three of mine, smallest to largest:

**`gitcommit`** encodes our commit ritual: Conventional Commits, description in French, strictly one line, show the message and get confirmation, never push. Trivial? It is. But it is exactly the kind of low-stakes convention that erodes when every commit message is negotiated from scratch. Since the skill exists, every commit in every repo looks the same.

**`release`** encodes releasing: collect commits since the last tag, *estimate the bump* (breaking marker or API removal means major, features mean minor, the rest is patch), update the changelog, commit, tag. The interesting part is that the estimation heuristics are written down: the agent applies the same judgment call I would, because I had to articulate it once.

**`yacast-go-service`** is the heavyweight: a ~200-line skill (plus reference files and templates) that scaffolds a complete Go microservice to house conventions. It opens with the non-negotiables as a table: configor with an env-var prefix, gorilla/mux behind a middleware chain, sqlx with separate read and write connection strings, our zap wrapper injected through request context, juju/errors annotations, Prometheus at `/metrics`, version extracted from the changelog and injected via ldflags, multi-stage Dockerfile to scratch, GitLab CI, Ansible plus systemd for deployment. Then the canonical package layout, then a numbered scaffolding procedure: copy, rename module, set env prefix, trim config, fix app name and metrics namespace.

Its first rule is the one I would give a junior engineer: **"don't invent, copy an existing service and rename."** That sentence does more work than everything after it.

## What writing the skill actually does

The uncomfortable discovery: before the skill, our "house conventions" were not written anywhere. They lived in the heads of the people who had copied the last service, in the diff patterns of code review. Writing `yacast-go-service` forced the conventions to become explicit, including the judgment calls nobody had ever stated out loud (which packages are mandatory, what a service must expose to be deployable, what the read/write DB split is *for*).

That gives skills a strange double life: they are documentation that executes. The same file onboards a human perfectly well, but unlike the wiki page we never wrote, it is exercised constantly, so rot is detected. When a convention changes, the skill fails visibly on the next scaffold, and gets fixed the way broken code gets fixed, not the way stale docs get ignored.

The framing that stuck with me: **an agent is a permanently new team member with perfect recall of what you wrote and zero memory of what you did not.** Every convention that lives only in heads is invisible to it. Skills are just onboarding documentation with a reader that never skips a paragraph, never assumes, and shows up to work already having read everything.

## Splitting the two correctly

The boundary that took me longest to find, stated plainly:

- **AGENTS.md holds what is true about this repo**: commands, architecture, inventory. It answers "where am I?"
- **Skills hold what is true about how we work**: rituals, scaffolds, conventions. They answer "how do we do this here?"

Mixing them fails in both directions. House style pasted into thirty AGENTS.md files drifts thirty ways; repo facts hoisted into a skill make it wrong everywhere except one place.

## Takeaways

- The default output of a coding agent is plausible, not yours. Closing that gap is a writing problem, not a model problem.
- Keep repo files factual: commands, architecture, feature inventory. The inventory prevents both re-implementation and accidental simplification.
- Promote anything cross-repo into a skill, even tiny rituals like commit messages. Small conventions erode first.
- Write the junior-engineer rule at the top: "don't invent, copy and rename" outperforms pages of specification.
- A skill is documentation that executes: it rots visibly, which is more than can be said for any wiki.
- Whatever you never wrote down, the agent does not know. That was always true of new colleagues too; agents just removed the option of hallway osmosis.
