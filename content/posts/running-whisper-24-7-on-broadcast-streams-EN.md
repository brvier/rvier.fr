---
title: 'Running Whisper 24/7: transcribing French TV and radio on a fleet of GPUs'
date: '2026-07-24'
lang: en
description: 'Experience report from a production speech-to-text pipeline: WhisperX and pyannote diarization on multi-GPU workers, dynamic GPU rebalancing, crash-only CUDA error handling, and audio that is never as clean as the demo.'
ogDescription: 'WhisperX and pyannote diarization on multi-GPU workers: dynamic GPU rebalancing, crash-only CUDA error handling, and broadcast audio that is never clean.'
keywords: Whisper, WhisperX, speech-to-text, speaker diarization, pyannote, CUDA, GPU, ffmpeg, Python
image: https://rvier.fr/images/whisper-broadcast-pipeline.png
summary: 'WhisperX and pyannote diarization on multi-GPU workers, 24/7: dynamic GPU rebalancing between transcription and diarization, crash-only CUDA error handling, and silence injection for missing audio.'
---

In [the pgvector post](deduplicating-speakers-with-pgvector-and-hnsw-EN.html) I explained how we deduplicate thousands of speaker voices in PostgreSQL. This post is about the upstream side: where those transcripts and voice embeddings come from. Our platform transcribes French TV and radio streams continuously, with word-level timestamps, speaker diarization and speaker identification. The workers are Python, PyTorch and ffmpeg, running unattended on multi-GPU machines.

<img src="../images/whisper-broadcast-pipeline.png" alt="Broadcast STT pipeline: extraction threads, WhisperX GPUs, diarization GPUs, result threads" loading="lazy" width="1200" height="627">

Running Whisper on a file is a tutorial. Running it around the clock on live broadcast audio, on machines nobody logs into, is a different exercise: the interesting problems are not in the model, they are in everything around it.

## The unit of work: five minutes, plus padding

The stream is cut into tasks of five minutes of audio per channel. Each worker pulls `(mediaid, timestamp)` jobs from a small queue API, so adding capacity means starting one more machine: no scheduler, no assignment logic, workers just pull when they have room.

Each task actually downloads 15 seconds *before* and *after* its window. Sentences do not respect five-minute boundaries, and both Whisper and the diarizer behave badly on words cut in half. The padded edges are transcribed, then trimmed when results are stored, and a check on existing results makes reprocessing idempotent: a re-queued task that was already done is acknowledged and skipped.

## One machine, three kinds of work

A worker host is a single Python program shaped as a chain of queues, because a GPU box runs three very different workloads at once:

```
extraction (threads, network + ffmpeg)
   └─▶ whisperx queue ─▶ WhisperX (1 process per GPU)
          └─▶ diarization queue ─▶ diarization + embeddings (1 process per GPU)
                 └─▶ gender queue ─▶ segmentation (thread)
                        └─▶ result queue ─▶ XML export + indexing (threads)
```

I/O-bound stages (chunk download, ffmpeg conversion, result upload) are plain threads. GPU stages are real processes, one per GPU, started with the `spawn` multiprocessing context, each loading its own models. The extraction stage throttles itself when the transcription queue grows past a threshold, so a slow GPU stage back-pressures the whole pipeline instead of filling the disk with WAV files.

## GPU placement, and stealing GPUs at runtime

On an 8-GPU machine the static assignment interleaves the two heavy stages: odd GPUs run WhisperX (large-v3, float16), even GPUs run pyannote diarization plus embedding extraction. Interleaving matters on dual-socket machines: neighbouring GPUs often share a PCIe switch and thermal envelope, and pairing a transcription GPU with a diarization GPU spreads the load better than putting all of one workload on one side.

The static split is a guess, and the real ratio between transcription time and diarization time depends on the content (talk radio and music-heavy stations diarize very differently). So the supervisor loop rebalances at runtime by watching the queues:

- diarization queue above 100 items: terminate one WhisperX process, restart that GPU as a diarization worker;
- transcription queue loaded while diarization idles: give the GPU back.

It is crude, a thermostat rather than a scheduler, but it converges in a couple of cycles and it removed the manual per-host tuning we did at the beginning.

## Crash-only CUDA error handling

The single most valuable design decision: **a GPU process never tries to recover from a CUDA error in-process**. When an exception contains `CUDA error`, `device-side assert` or `out of memory`, the process re-queues its current job (with a bounded retry count), logs, and calls `sys.exit(1)`.

We tried the alternative first. After certain CUDA failures the context is poisoned: subsequent inferences return garbage or hang, and `torch.cuda.empty_cache()` incantations only hide it. A dead process is a clean state; a "recovered" one is a question mark.

Restart responsibility is layered:

- the supervisor loop checks every 30 seconds that each GPU process is alive, and respawns dead ones on the same GPU;
- supervisord restarts the whole worker if the main process dies;
- the deploy script kills lingering GPU processes and reboots the host, because a fresh CUDA state at deploy time costs one minute and debugging a half-freed GPU costs an afternoon.

Jobs whose retry budget is exhausted are marked with a per-stage error status in the queue API (download error, transcription error, diarization error...), which makes "what is failing where" a one-query question.

## Broadcast audio is never clean

The audio arrives as MPEG-TS chunks from our capture platform's storage service. Two realities of broadcast capture shaped the extraction stage:

**Chunks go missing.** Storage rebalancing, capture hiccups: sometimes a chunk is not there. Skipping it would silently shift every later timestamp in the window, and timestamps are the product (they anchor words to the broadcast timeline). So a missing chunk is replaced by a *generated silence chunk* of the same duration (`ffmpeg -f lavfi -i anullsrc`), and the job aborts only past a missing-chunk budget. The transcript loses the missing words; every other word keeps its exact position.

**PTS are discontinuous.** Concatenating TS chunks naively confuses decoders, because presentation timestamps jump across chunk boundaries. The ffmpeg *concat demuxer* (a file list and `-f concat`, rather than binary concatenation) rewrites a coherent timeline before the 16 kHz mono-pipeline WAV conversion.

Chunk downloads themselves are written to a temp file, fsynced, then renamed: the same atomic-write habit as everywhere else in the platform, because a truncated TS file produces the weirdest downstream errors.

## The Whisper settings that matter here

Model choice and decoding options are ordinary (large-v3, float16, batch inference, beam 5). Three settings earned their place in production on broadcast content:

- `condition_on_previous_text: False`. Conditioning helps coherence on clean speech and *amplifies* hallucination loops on jingles, music beds and silence. On broadcast audio the trade is not close: a transcript that repeats "merci d'avoir regardé" forty times is worse than slightly less fluent punctuation.
- An aggressive VAD (pyannote, low onset/offset thresholds) in front of the model. Whisper hallucinates on non-speech; the cheapest fix is to not show it non-speech.
- Word-level alignment as a separate pass (WhisperX). Whisper's own segment timestamps are too coarse to attribute words to diarized speaker turns; the forced-alignment pass is what makes "who said this word" possible at all.

## From voices to identities

Diarization (pyannote) answers "speaker A, speaker B" within one five-minute window. To turn that into *persons*, the diarization stage also extracts one embedding per detected speaker: segments shorter than 5 seconds are skipped (short crops crash the embedding window and produce noisy vectors anyway), the rest are averaged. The embedding is sent to the speaker API, which does the pgvector similarity search from the previous post and returns a stable speaker UUID, either matched or newly created.

Two failure-handling details:

- Diarized segments that end up with no speaker label inherit the label of the nearest-in-time labelled segment. Crude, measurably better than dropping the words.
- If the speaker API is unreachable after retries with backoff, the worker generates deterministic emergency UUIDs and ships the transcript anyway. Speaker identity can be reconciled later; a hole in the 24/7 transcript timeline cannot.

## Takeaways

- The model is the easy part. The engineering is in the pipeline around it: queues, padding, retries, and the assumption that every stage will fail.
- Let GPU processes die. Crash-only handling of CUDA errors plus layered restarts (supervisor loop, supervisord, reboot at deploy) beats any in-process recovery we tried.
- Rebalance GPUs with a thermostat, not a scheduler: watch the queues, move one GPU at a time.
- Never let missing input shift your timeline: inject silence, keep the timestamps honest.
- Degrade rather than stop: emergency speaker IDs, nearest-speaker fallback, silence budgets. A 24/7 pipeline's first duty is to keep flowing.
