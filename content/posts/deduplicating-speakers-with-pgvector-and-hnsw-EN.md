---
title: 'Vector search in PostgreSQL: deduplicating thousands of speakers with pgvector and HNSW'
date: '2026-07-10'
lang: en
description: How I replaced an O(N²) cosine similarity self-join that timed out with an HNSW index and LATERAL k-nearest-neighbour queries in pgvector, to deduplicate thousands of voice embeddings in PostgreSQL.
ogDescription: 'From an O(N²) cosine self-join that timed out to HNSW + LATERAL k-nearest-neighbour queries: deduplicating voice embeddings in PostgreSQL.'
keywords: PostgreSQL, pgvector, HNSW, vector search, embeddings, speaker diarization
image: https://rvier.fr/images/pgvector-hnsw-dedup.png
summary: 'From an O(N²) cosine self-join that timed out to HNSW + LATERAL k-nearest-neighbour queries: deduplicating voice embeddings without leaving PostgreSQL.'
---

Most pgvector articles are about RAG and LLM chatbots. This one is not. It is about a very concrete, non-LLM problem: identifying *who is speaking* on TV and radio, at scale, with nothing more exotic than PostgreSQL.

## The problem

I work on a media analytics platform that runs speech-to-text with speaker diarization on broadcast streams, 24/7. Each diarized segment produces a voice embedding: a 256-dimension vector that characterizes a voice. Two segments from the same person should produce vectors that are close in cosine similarity; two different people should not.

The goal is a `speakers` table with one row per real-world person. Each speaker accumulates *fingerprints* (one embedding per segment), and carries a *centroid*: the average of its active fingerprint embeddings. When a new segment arrives, we match it against existing centroids; above a similarity threshold (0.71 for us) it is attached to the existing speaker, otherwise a new speaker is created.

The schema looks like this:

```
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE speakers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    embedding VECTOR(256),   -- centroid of active fingerprints
    status SMALLINT NOT NULL DEFAULT 0
);

CREATE TABLE speaker_fingerprints (
    id UUID PRIMARY KEY,
    speaker_id UUID NOT NULL REFERENCES speakers(id) ON DELETE CASCADE,
    embedding VECTOR(256),
    status SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX ON speakers USING hnsw (embedding vector_cosine_ops);
CREATE INDEX ON speaker_fingerprints USING hnsw (embedding vector_cosine_ops);
```

Threshold-based matching is never perfect. Diarization is noisy, voices change with audio quality, and a borderline segment sometimes creates a brand new speaker for a person who already exists. Over months of continuous ingestion, you accumulate near-duplicate speakers. So we need a periodic deduplication job: find every pair of speakers whose centroids are almost identical, and merge them.

## The naive version: an O(N²) self-join

The first implementation was the obvious one. Compare every speaker with every other speaker, take the best pair above the merge threshold:

```
SELECT a.id, b.id,
       1 - (a.embedding <=> b.embedding) AS sim
  FROM speakers a
  JOIN speakers b ON a.id < b.id
 WHERE a.status >= 0 AND b.status >= 0
   AND 1 - (a.embedding <=> b.embedding) >= 0.92
 ORDER BY sim DESC
 LIMIT 1;
```

This works fine in the demo, with a hundred speakers. With thousands of speakers it is a disaster, and it is important to understand *why*: **the HNSW index is not used at all**. pgvector's indexes only accelerate one query shape: `ORDER BY embedding <=> $1 LIMIT k`. A distance expression inside a `WHERE` clause or a join predicate is evaluated the hard way: a sequential scan computing N×(N−1)/2 cosine distances. At 5,000 speakers that is more than 12 million 256-dimension distance computations for a single query. Ours started hitting the statement timeout, and to make it worse, the job called this query once *per merge*.

## The fix: ask the index for k neighbours per speaker

The insight is to stop asking "which pairs are above the threshold?" (a question the index cannot answer) and instead ask, for each speaker, "what are your k nearest neighbours?", which is exactly the query shape HNSW is built for. In SQL this is a `CROSS JOIN LATERAL`:

```
WITH raw_pairs AS (
    SELECT DISTINCT ON (LEAST(a.id, n.id), GREATEST(a.id, n.id))
           LEAST(a.id, n.id)    AS sid_a,
           GREATEST(a.id, n.id) AS sid_b,
           1 - (a.embedding <=> n.embedding) AS similarity
      FROM speakers a
      CROSS JOIN LATERAL (
          SELECT id, embedding
            FROM speakers
           WHERE status >= 0
             AND embedding IS NOT NULL
             AND id <> a.id
           ORDER BY embedding <=> a.embedding   -- HNSW kicks in here
           LIMIT 5                                -- k neighbours per speaker
      ) n
     WHERE a.status >= 0 AND a.embedding IS NOT NULL
       AND 1 - (a.embedding <=> n.embedding) >= 0.92
)
SELECT * FROM raw_pairs ORDER BY similarity DESC;
```

Three details matter here:

- The inner query is `ORDER BY … LIMIT k`, so each of the N outer rows costs one indexed k-NN lookup instead of a full scan. The overall cost goes from O(N²) to roughly O(N·log N).
- k-NN returns both directions of every pair (A finds B, then B finds A). `DISTINCT ON (LEAST(id1, id2), GREATEST(id1, id2))` collapses each unordered pair to a single row.
- The threshold filter is still there, but it now filters a small candidate set (N×k rows) instead of driving the join.

A small k (5 in our case) is enough. If a speaker happens to have more than k mergeable duplicates, the surplus is simply not visible in this pass, and that is fine, because merging runs in rounds.

## Merging in rounds, not in one pass

You cannot just take the candidate list and merge every pair in it. Each merge moves the surviving speaker's centroid (it is recomputed as the average of the combined fingerprints), so a pair that was above the threshold when the list was built may no longer be a real duplicate two merges later. The opposite also happens: a merge can pull a centroid closer to a third speaker.

So the bulk merge is a loop:

- Fetch all candidate pairs with the indexed query above, sorted by similarity descending.
- Greedily walk the list and merge, *skipping any pair that involves a speaker already touched in this round*. Every merge in a round is therefore between two untouched speakers, computed from fresh centroids.
- When the list is exhausted, start a new round: re-fetch candidates with the drifted centroids.
- Stop when a round performs zero merges (with a hard cap on rounds as a safety net against pathological oscillation).

Each merge keeps the speaker that has the most fingerprints and folds the smaller one into it, then the centroid is recomputed. This is also where the missed "surplus" duplicates from the k-NN limit get caught: after round one merged the closest pairs, the survivors find their remaining duplicates in round two.

## Keeping centroids honest

One last piece makes the whole thing stable. A centroid computed from noisy fingerprints drifts, and a merge can import a few bad segments. After every recomputation, each fingerprint is re-evaluated against the new centroid: active fingerprints that fall below an outlier threshold are demoted, and, the interesting part, previously demoted fingerprints whose similarity is now *above* the threshold are re-promoted, because removing bad samples sharpens the centroid and can bring borderline segments back in. This demote/promote pass loops up to three times until the centroid converges. All of it is plain SQL: `AVG(embedding)` with a `FILTER` clause, and two `UPDATE` statements using the `<=>` operator.

## Takeaways

- pgvector's HNSW index accelerates exactly one query shape: `ORDER BY embedding <=> $1
            LIMIT k`. If your distance expression lives in a `WHERE` clause or a join condition, you are doing a sequential scan, whatever indexes exist.
- `CROSS JOIN LATERAL` is the bridge: it turns "compare everything with everything" into "one indexed k-NN lookup per row".
- HNSW is approximate. Combined with a small k, a single pass can miss pairs, design the surrounding process (rounds, periodic re-runs) so that misses are caught later instead of pretending the index is exhaustive.
- When entities are mutable aggregates (centroids), never batch-apply decisions computed from stale state. Greedy non-overlapping merges per round keep every decision based on fresh data, at the cost of a few extra query rounds, which are cheap now that they are indexed.

No vector database, no extra infrastructure: PostgreSQL, one extension, and the right query shape.
