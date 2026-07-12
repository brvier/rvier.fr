---
title: 'One Go binary, one YAML file, one SQLite database: why I wrote my own monitoring tool'
date: '2026-07-09'
lang: en
description: 'Why I wrote Gjallar, a KISS monitoring service in ~3,400 lines of Go: zero CGO thanks to pure-Go drivers (pgx, go-ora, pro-bing), a lock-free alert state machine, and SIGHUP hot-reload with validation.'
ogDescription: 'Gjallar: a KISS monitoring service in ~3,400 lines of Go. Zero CGO, a lock-free alert state machine, SIGHUP hot-reload with validation.'
keywords: Go, Golang, monitoring, SQLite, CGO, pgx, go-ora, self-hosted
summary: 'Gjallar, a KISS monitoring service in ~3,400 lines of Go: zero CGO thanks to pure-Go drivers, a lock-free alert state machine, and SIGHUP hot-reload with validation.'
---

I needed to watch a fleet of heterogeneous services: HTTP endpoints, PostgreSQL databases, a few Oracle instances, Redis, Elasticsearch indexes that must stay fresh, machines that should answer ping, and some Prometheus metrics. And I needed to be told, on Telegram, by SMS, on Signal, when something goes down, and when it comes back.

The classic answer is a monitoring *platform*: Prometheus plus Alertmanager plus Grafana plus a handful of exporters, or a container running a Node.js app with a database. All great tools. But for a few dozen checks, I did not want to operate a second distributed system just to know whether the first one is up. And none of the lightweight options could query Oracle without me installing the Oracle client libraries somewhere.

So I wrote [Gjallar](https://github.com/brvier/Gjallar): a KISS monitoring service. **One static binary, one YAML config file, one SQLite file.** A black-and-red status page with history, HTMX-refreshed. About 3,400 lines of Go. MIT licensed.

## Zero CGO, on purpose

The whole tool builds with `CGO_ENABLED=0`:

```
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w"
```

That is only possible because every dependency that would traditionally bind to a C library has a pure-Go replacement nowadays, and they are excellent:

- [pgx](https://github.com/jackc/pgx) for PostgreSQL: no libpq;
- [go-ora](https://github.com/sijms/go-ora) for Oracle: *no Oracle Instant Client*, which alone justified the project. If you have ever deployed the Oracle client on a minimal box, you know;
- [pro-bing](https://github.com/prometheus-community/pro-bing) for ICMP echo, privileged or unprivileged;
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) for storage: SQLite transpiled to pure Go, no libsqlite3;
- Redis needs no driver at all: the check speaks the protocol directly: TCP connect, optional `AUTH`, `PING`, expect `+PONG`.

The result is a single self-contained binary (about 36 MB, most of it the SQLite and Oracle drivers) that cross-compiles from my laptop to any target with `GOOS`/`GOARCH`, and deploys with `scp`. No Docker, no package manager, no shared libraries, no "works on my machine".

## A lock-free alert pipeline

Monitoring tools are naturally concurrent, every monitor waits on the network most of the time, and concurrency is where side projects usually grow their first mutex jungle. Gjallar has no locks around its state at all, because of how the pipeline is shaped:

```
one goroutine per monitor ──▶ results channel ──▶ single consumer
                                                  (state machine + SQLite writes)
```

Each monitor runs its check loop in its own goroutine and sends `check.Result` values into a shared channel. A single consumer goroutine owns everything downstream: the up/down state machine, incident rows, and history writes. Since only one goroutine ever touches the state map and the database connection, there is nothing to lock, and SQLite, which dislikes concurrent writers, gets exactly one.

The per-monitor state is small and explicit:

```
type monitorState struct {
    down         bool
    consecFails  int
    downSince    time.Time
    lastNotified time.Time
    threshold    int           // consecutive failures before DOWN fires
    realert      time.Duration // reminder interval while down; 0 = disabled
    notifiers    []string
}
```

Two design points earned their keep in production:

- **State survives restarts.** At startup, each monitor's state is seeded from any open incident in SQLite. A restart while something is down neither re-fires the DOWN alert nor misses the recovery notification. Deploying a new version during an outage is a non-event.
- **Notifications are dispatched asynchronously.** The consumer must never block: a slow SMTP server or a rate-limited Telegram API cannot back-pressure the whole pipeline. Sends go out in their own goroutines with a 15-second timeout.
- **Alerts fire after N consecutive failures**, not on the first blip, no flapping noise, and an optional `realert` interval reminds you while an incident stays open.

## Configuration that respects operations

Everything lives in one YAML file, with defaults, named notifiers, and monitor groups:

```
defaults:
  interval: 60s
  timeout: 10s
  failure_threshold: 3
  alerts: [ops-telegram]

alerts:
  ops-telegram:
    url: "telegram://TOKEN@telegram?chats=123456789"

monitors:
  - name: app-db
    type: postgres
    dsn: "postgres://monitor:${PG_PASSWORD}@db1:5432/app"
    query: "SELECT count(*) FROM jobs WHERE status = 'stuck'"
    rule: "== 0"
```

Three small features make it pleasant to operate:

- **Hot reload on SIGHUP**: `systemctl reload gjallar` applies the new config, but only after it has been *fully validated*. A broken YAML keeps the running configuration alive and logs the error, instead of taking the monitoring down with it. Your watcher should be the last thing that dies from a typo.
- **`${VAR}` environment expansion** for secrets, with a clear startup failure if a referenced variable is undefined, and a bare `$` (say, in a `~ ^OPEN$` regex rule) left untouched.
- **A `-check` flag** for dry-run validation, so CI can lint the config before it ever reaches the server.

## What it deliberately does not do

No clustering, no agents, no plugin system, no time-series dashboards, no user accounts. History is pruned after a configurable retention (30 days by default) so the SQLite file stays small forever. If a need is served well by an existing simple mechanism, systemd for the service lifecycle, shoutrrr URLs for the twenty notification services I will never use, Gjallar delegates instead of reimplementing.

This is the part I would defend the hardest. Every monitoring tool I have abandoned over the years died of the same disease: it slowly became a platform, and one day the monitoring needed monitoring. A tool whose whole state fits in one SQLite file and whose whole behaviour fits in one YAML file is a tool you still understand at 3 a.m., eighteen months after you wrote it.

Code and documentation: [github.com/brvier/Gjallar](https://github.com/brvier/Gjallar).
