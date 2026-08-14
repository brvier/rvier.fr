---
title: 'Recording broadcast 24/7: the architecture of a capture platform'
date: '2026-08-14'
lang: en
description: 'The architecture behind years of continuous TV and radio capture: time-addressed immutable chunks, multicast as a decoupling layer, a control plane that can die without stopping the recording, tiered storage, and a repair layer that assumes nothing is ever perfect.'
ogDescription: 'Time-addressed immutable chunks, multicast decoupling, a control plane that can die safely, tiered storage, and a repair layer that assumes imperfection.'
keywords: broadcast, capture, DVB, ffmpeg, Go, multicast, MPEG-TS, Ceph, architecture, SRT
image: https://rvier.fr/images/hyperion-capture-architecture.png
summary: 'The capstone of the broadcast series: how the capture platform records TV and radio around the clock, with time-addressed chunks, multicast decoupling, tiered storage and a repair layer.'
---

Every post in this series so far consumed the same input: [Whisper transcription](running-whisper-24-7-on-broadcast-streams-EN.html), [on-screen OCR](ocr-llm-enrichment-broadcast-transcripts-EN.html), [speaker deduplication](deduplicating-speakers-with-pgvector-and-hnsw-EN.html), even the [byte-for-byte audio profile format](byte-for-byte-reproducing-a-legacy-binary-format-EN.html). This post is about where that input comes from: the platform that records French TV and radio channels continuously, from terrestrial DVB multiplexes, FM antennas, satellite, web radio streams and IP sources, and has been doing so for years.

<img src="../images/hyperion-capture-architecture.png" alt="Capture platform architecture: tuners, multicast, encoders, storage servers, tiered storage, control plane" loading="lazy" width="1200" height="627">

A capture platform has one non-negotiable property that shapes every other decision: **the signal does not wait**. A web service that is down retries later; a broadcast minute that was not recorded is gone forever. Everything below follows from that.

## One abstraction to rule the platform: the time-addressed chunk

The whole system agrees on a single data model: a recording is a sequence of immutable MPEG-TS chunks, addressed by `(media, timestamp)`. Not filenames, not playlists, not sessions: a channel identifier and a UTC time.

Every component speaks this address and nothing else. Encoders produce chunks and upload them. Storage servers store and serve them over HTTP, with the chunk duration in a response header so a consumer can walk the timeline chunk by chunk. The STT, OCR and fingerprinting workers from the previous posts fetch `(media, t)`, `(media, t + len)`, and so on. The address does not encode a location, so a chunk can live on any storage tier without anyone noticing. The purger applies retention by deleting address ranges.

Each channel is also recorded in several qualities, each with its own retention: *high*, *low* and *verylow* video profiles, plus a *raw* profile that keeps the audio only, in its native broadcast format, for later audio analysis. The expensive profiles live shorter lives than the cheap ones, which is how "keep everything" stays affordable. On top of this, the storage API exposes a `best` endpoint: give me this minute, in the best quality currently available. Consumers express intent; the storage layer resolves it. When a high-quality chunk is missing or already purged but a lower one exists, the pipeline keeps working instead of failing.

## The data plane: tuners, multicast, encoders

Acquisition and encoding are deliberately separate machines with a multicast network between them:

```
DVB / FM / satellite / web-radio tuners
        │  (UDP multicast, one group per stream)
        ▼
encoders (Go service driving ffmpeg, N per site)
        │  (HTTP, ordered chunk upload)
        ▼
storage servers ("blobbers") ─▶ hot and archive tiers
```

Multicast is the decoupling layer, and it is hard to overstate how much simplicity it buys. A tuner emits one stream; any number of consumers can join the group: the nominal encoder, a spare warming up, an engineer's ffprobe during an incident, all without the tuner knowing or caring. Failing an encoding over to another machine is joining a multicast group, not re-plumbing a source. The tuners themselves stay thin: drive the DVB frontend, forward transport streams, and export signal metrics (RSSI included) so that a degrading antenna shows up in monitoring before it becomes a hole in the archive.

The encoder is a Go supervisor around ffmpeg, and its changelog is a museum of everything that can go wrong between a live source and an immutable file. The lessons that stuck:

- **Upload chunks strictly in order per channel, concurrently across channels.** Downstream consumers walk the timeline; an out-of-order upload looks exactly like a hole.
- **Never upload a file that might still be written.** A minimum-age threshold before upload, because "the file exists" and "the file is complete" are different claims.
- **Read the file once, for both the checksum and the HTTP body.** An early version read it twice and could checksum a different byte stream than it sent. CRC first, then send the same bytes.
- **Classify your deaths.** Every way an encoding can fail (ffmpeg dying, stale segments, PTS jumping beyond plausible amplitude, timestamps from the future) sets an explicit error reason on the encoding record. "It stopped" is not a diagnosis.
- **Make probing a first-class operation.** A dedicated RPC probes a source (multicast, SRT, HTTP) without starting an encoding, so "is the source alive?" does not require touching production state.
- **Sign your work.** Each encoder writes its hostname into the MPEG-TS metadata of what it produces. When a bad chunk surfaces weeks later, you know which machine made it.

Video encoding itself runs on GPUs and dedicated transcoding hardware (H.264/H.265, deinterlacing included), the encoder service staying deliberately agnostic about which accelerator sits behind ffmpeg.

## The control plane can die; the recording cannot

Above the data plane sits a supervisor, Hyperion Core: the authenticated REST API used by internal tooling, a small public API exposing availability, Prometheus metrics, and gRPC links to every tuner, encoder and purger. It owns the source of truth in PostgreSQL: which channels exist, where they are encoded, what the retention rules are.

The design rule is that the control plane configures the data plane but never sits inside it. Chunks flow from encoders to storage directly; if the core is down, tuners keep tuning, encoders keep encoding, uploads keep uploading. You lose the ability to change things, not the recording itself. For a system whose product is "we did not miss it", this separation is the single most protective architectural decision.

## Storage: tiers behind one API

Behind the storage servers, storage is tiered: a hot Ceph cluster for the recent window that the analysis workers hammer, a cold Ceph cluster for the long tail, and some old NAS volumes left from before Ceph was installed, still quietly serving the chunks of their era. The purger enforces per-channel retention. Because everything is addressed by `(media, timestamp)`, the tiers are invisible to consumers: the same request serves from wherever the chunk lives.

Chunks *can* be migrated between tiers, but deliberately by hand: a CLI tool moves a channel and time range when there is a reason to (freeing a tier, consolidating an archive), rather than an automatic mover shuffling data in the background. And that is also what made the storage generations survivable: the pre-Ceph archives never had to be migrated on a deadline, they just stayed behind the same API while new writes went to Ceph, and get moved range by range when it is actually worth it.

Immutability is what keeps this boring. A chunk is written once and never modified, so replication, migration and caching never worry about coherence: the [tmpfs tricks from an earlier post](tmpfs-dev-shm-the-forgotten-ramdisk-optimization-EN.html) work precisely because a cached chunk can never be stale.

## Packing chunks into blobs

There is one more layer between "chunk" and "disk", and it exists because of arithmetic. Chunks are seconds long; multiply by two or three qualities, by every channel, by 24/7, by years, and you get hundreds of millions of small files. Filesystems hate that: inode pressure, directory listings that take minutes, backups and scrubs dominated by metadata rather than data. Small files are the classic way to kill a storage cluster with data that is, in total, not even that big.

So chunks are not stored as files. A storage daemon (affectionately named *Blobibloba*, hence "blobbers") appends them into **blob files**: one blob per channel, per quality, per fixed time window, the window being simply the chunk timestamp truncated. Next to each blob lives a small index file mapping chunk timestamp to offset and length, and each entry carries the chunk's duration, CRC32 and a discontinuity flag. Serving `(media, t)` is a lookup in the index and one positioned read; appending is sequential I/O, which is exactly what both spinning disks and Ceph like.

Two properties fall out of the time-windowed layout for free:

- **Locality matches access patterns.** Consumers read timelines, and a timeline is contiguous bytes in one blob, not a thousand file opens scattered across a cluster.
- **Retention becomes directory math.** The purger walks blob directories and deletes whole time windows past their per-quality retention (the high profile does not have to live as long as the low one). Deleting a day of a channel is removing a handful of blobs, not unlinking tens of thousands of files.

Above the blobs, the storage servers keep the routing map (which storage tier serves which channel and period, cached in Redis), while the blob index files themselves remain the source of truth for what exists: no separate chunk database to keep consistent with the bytes on disk. The storage servers expose the higher-level reads: single chunks, and streamed extracts of arbitrary time ranges that concatenate chunks on the fly with ffmpeg, cumulative PTS offsets across discontinuities included, cached with honest HTTP semantics (Range, ETag, `304 Not Modified`).

## The repair layer: assume imperfection

The uncomfortable truth of continuous capture is that something is always slightly broken: an antenna degrades, a multiplex hiccups, a machine reboots. The platform treats gaps as a normal operational object, not an exception:

- checkers walk the timeline and report discontinuities per channel, so a hole is detected in minutes, not discovered by a customer months later;
- a recovery tool re-imports missing ranges from secondary captures and backup receivers, and only fills actual holes unless explicitly forced;
- consumers degrade instead of stopping: the STT workers [inject silence for a missing chunk](running-whisper-24-7-on-broadcast-streams-EN.html) so the transcript timeline stays honest.

Detection, repair, and graceful degradation: three separate mechanisms, because no single one of them is reliable enough alone.

## Takeaways

- Pick one address for your data and make every component speak it. `(media, timestamp)` for an immutable chunk turned tiering, retention, repair and every analysis pipeline into independent problems.
- Multicast between acquisition and processing is the cheapest high-availability primitive there is: one producer, any number of consumers, zero coupling.
- Let consumers express intent (`best` quality), not storage paths. It is the difference between a degraded service and an outage.
- Keep the control plane out of the data path. Configuration can be down; capture cannot.
- Do the small-file arithmetic before your filesystem does it for you: pack chunks into time-windowed blobs with an index, and retention becomes deleting a few big files instead of millions of small ones.
- Chunks are immutable, uploads are ordered, files are checksummed from the bytes actually sent, and every failure mode has a name. Archive quality is the sum of these small disciplines.
- Build the repair layer on day one. A 24/7 platform is not one that never breaks; it is one that notices, repairs and degrades honestly.
