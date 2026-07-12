---
title: Media Insights Orchestrator (Yacast)
section: professional
weight: 40
image: ./images/yacast-media-insights-orchestrator.svg
alt: Media Insights Orchestrator, task distribution to GPU workers
stack: Go / PostgreSQL / pgvector
---

Design and development of the orchestration service at the heart of the Media Insights platform: distributes Speech-to-Text, OCR and AI enrichment tasks to a fleet of GPU workers, consolidates their results, and manages speaker identification at scale using vector embeddings with pgvector/HNSW similarity search.
