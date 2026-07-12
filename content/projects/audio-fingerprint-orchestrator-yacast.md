---
title: Audio Fingerprint Orchestrator (Yacast)
section: professional
weight: 80
image: ./images/yacast-audio-fingerprint.svg
alt: Audio Fingerprint Orchestrator, fingerprint catalog and index buckets
stack: Go / PostgreSQL / Linux
---

Development of the distributed system managing the audio fingerprint reference catalog: fingerprint creation, movement and deletion across index buckets, per-bucket operation serialization, idle-aware index flushing, and in-memory (tmpfs) caching on high-concurrency workers.
