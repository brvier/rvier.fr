---
title: Gjallar
section: opensource
weight: 197
image: ./images/gjallar.png
alt: Gjallar, a KISS monitoring service in Go
link: https://github.com/brvier/Gjallar
linkText: View Project
---

A KISS monitoring service in Go: one static binary, one YAML file, one SQLite database, zero CGO. Monitors HTTP (status, body regex, TLS expiry), PostgreSQL, Oracle, Redis, Elasticsearch freshness, Prometheus metrics and ICMP ping, with a black and red HTMX status page and history. Alerts through any shoutrrr URL (Telegram, ntfy, Slack, email...), Free Mobile SMS or Signal, after a failure threshold and again on recovery. Warnings for degraded-but-up states, SIGHUP hot-reload with validation, incidents that survive restarts. Runs in production at Yacast.
