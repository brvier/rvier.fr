---
title: 'One Go binary, one YAML file, one SQLite database: why I wrote my own monitoring tool'
date: '2026-07-09'
lang: en
description: 'Why I wrote Gjallar, a KISS monitoring service in ~3,400 lines of Go: zero CGO thanks to pure-Go drivers (pgx, go-ora, pro-bing), a lock-free alert state machine, and SIGHUP hot-reload with validation.'
ogDescription: 'Gjallar: a KISS monitoring service in ~3,400 lines of Go. Zero CGO, a lock-free alert state machine, SIGHUP hot-reload with validation.'
keywords: Go, Golang, monitoring, SQLite, CGO, pgx, go-ora, self-hosted
summary: 'Gjallar, a KISS monitoring service in ~3,400 lines of Go: zero CGO thanks to pure-Go drivers, a lock-free alert state machine, and SIGHUP hot-reload with validation.'
---

I had a few dozen things to watch: HTTP endpoints, PostgreSQL and Oracle databases, Redis, Elasticsearch indexes that have to stay fresh, machines that should answer ping, a handful of Prometheus metrics. I wanted to be told on Telegram and by SMS when one of them goes down, and when it comes back.

The usual answer is Prometheus with Alertmanager, Grafana and a few exporters, or a Node app in a container. For a few dozen checks that means running a second distributed system in order to know whether the first one is up. I went looking for something simpler and did not find it, and nothing I came across queried Oracle without the Instant Client installed somewhere.

So I wrote [Gjallar](https://github.com/brvier/Gjallar). One static binary, one YAML config file, one SQLite file, and a black and red status page with history, refreshed by HTMX. About 3,400 lines of Go, MIT licensed.

## Zero CGO, on purpose

The whole tool builds with `CGO_ENABLED=0`:

```
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w"
```

Every dependency that used to bind to a C library now has a pure-Go replacement, and they are good:

- [pgx](https://github.com/jackc/pgx) for PostgreSQL, so no libpq;
- [go-ora](https://github.com/sijms/go-ora) for Oracle, so no Oracle Instant Client. That one alone justified the project. Anyone who has deployed the Oracle client on a minimal box knows why;
- [pro-bing](https://github.com/prometheus-community/pro-bing) for ICMP echo, privileged or unprivileged;
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) for storage, SQLite transpiled to pure Go;
- Redis needs no driver at all. The check speaks the protocol directly: TCP connect, optional `AUTH`, `PING`, expect `+PONG`.

What comes out is a 36 MB binary, most of it the SQLite and Oracle drivers. I cross-compile it on my laptop with `GOOS`/`GOARCH` and deploy it with `scp`. Nothing else to install on the target.

## A lock-free alert pipeline

Every monitor spends most of its life waiting on the network, so the tool is concurrent by nature. That is usually where a side project grows its first mutex. Gjallar has no lock around its state, because of how the pipeline is shaped:

```
one goroutine per monitor ──▶ results channel ──▶ single consumer
                                                  (state machine + SQLite writes)
```

Each monitor runs its check loop in its own goroutine and sends `check.Result` values into a shared channel. A single consumer goroutine owns everything downstream: the up/down state machine, incident rows, history writes. Only that goroutine ever touches the state map and the database handle, so there is nothing to lock, and SQLite, which dislikes concurrent writers, gets exactly one.

The per-monitor state is small:

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

At startup, each monitor's state is seeded from any open incident found in SQLite. A restart while something is down neither re-fires the DOWN alert nor loses the recovery notification, which makes deploying a new version during an outage uneventful. Notifications are sent from their own goroutines with a 15 second timeout, so a slow SMTP server or a rate-limited Telegram API cannot back-pressure the pipeline. Alerts only fire after N consecutive failures, and an optional `realert` interval reminds me while an incident stays open.

## Configuration

Everything lives in one YAML file, with defaults, named notifiers and monitor groups:

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

`systemctl reload gjallar` sends SIGHUP and applies the new config, but only once it has been fully validated. A broken YAML logs the error and keeps the running configuration alive, instead of taking the monitoring down with it.

`${VAR}` expansion covers secrets, with a clear startup failure when a referenced variable is undefined, and a bare `$` (in a `~ ^OPEN$` regex rule, say) is left untouched. A `-check` flag validates a config file without starting anything, so CI can lint it before it reaches the server.

## What it does not do

Everything runs as one process on one machine, driven by one file. Nothing to deploy on the monitored hosts, and the status page is public or nothing, since there is no notion of a user. Time-series dashboards and plugins are out of scope. History is pruned after a configurable retention, 30 days by default, so the SQLite file stays small. Where a simple mechanism already exists, Gjallar delegates to it: systemd for the service lifecycle, shoutrrr URLs for the twenty notification backends I will never use.

I hold on to those limits. The monitoring tools I abandoned over the years all grew the same way, into platforms that eventually needed monitoring of their own. Keeping the state in one SQLite file and the behaviour in one YAML file is what stops that from happening again.

Code and documentation: [github.com/brvier/Gjallar](https://github.com/brvier/Gjallar).
