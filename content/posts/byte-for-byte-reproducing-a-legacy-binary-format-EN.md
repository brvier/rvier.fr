---
title: 'Byte-for-byte: reproducing an undocumented legacy binary format in Go'
date: '2026-06-05'
lang: en
description: 'Rewriting a legacy media pipeline in Go meant reproducing an undocumented binary sidecar format byte-for-byte: reverse-engineering floor(min(mean(|s16|)/256, 127)) per 10 ms window, a fractional carry accumulator, and the unit tests that prove conformance.'
ogDescription: Reverse-engineering an undocumented audio profile format, including its fractional carry bug-for-bug, and proving conformance with byte-identical tests.
keywords: Go, Golang, reverse engineering, legacy, binary format, migration, ffmpeg, PCM
summary: Reverse-engineering an undocumented audio profile format, fractional carry included, bug-for-bug, and proving conformance with byte-identical tests.
---

Every rewrite of a legacy system eventually hits the same wall: a file format, a wire message or a checksum that no one documented, that the original author left the company with, and that a dozen downstream consumers depend on *exactly as it is*. This is the story of one of those: a small audio "profile" sidecar file that I had to reproduce byte-for-byte while rewriting a twenty-year-old media conversion daemon in Go.

## The format nobody wrote down

The legacy pipeline transcodes broadcast recordings, and next to every converted media file it emits a small `.ap` ("audio profile") sidecar. Downstream tools use it to draw waveforms and to seek interesting audio regions without decoding the media. There is no header, no magic number, no spec, just a raw stream of bytes, one per 10 milliseconds of source audio.

Since consumers compare and cache these files, "close enough" was not an option: the Go rewrite had to produce the same bytes as the old C binary, on every input, forever. The old binary's source was available but barely readable; the format itself lived implicitly in one function. Reverse-engineered and written down, it comes out as:

```
byte = floor( min( mean(|s16 sample|) / 256 , 127 ) )
```

For each 10 ms window: take the interleaved signed 16-bit samples, average their absolute values, divide by 256, clamp to 127, floor to an integer. One byte out. Simple, and completely undiscoverable without reading the original code, because several nearby variants (RMS instead of mean absolute value, clamping at 128, rounding instead of flooring) produce plausible-looking waveforms that are subtly wrong.

## The devil is the fractional carry

The real trap was the window length. A 10 ms window at 44,100 Hz mono is 441 samples, fine. But at stereo 44,100 Hz it is 882, and at some rate/channel combinations the ideal window is *fractional*. The legacy binary did not resample or accumulate floats: it kept an integer accumulator of thousandths (a variable charmingly named `lProfilArrondi`), added the fractional part every window, and stole one extra sample whenever the accumulator crossed 1000.

Reproducing this exactly matters: get it wrong and the windows slowly drift against the original, so the two files agree at the start and diverge after a few seconds, the worst kind of wrong, the kind a quick eyeball check does not catch. The Go version reproduces the accumulator as-is:

```
spw := float64(channels*rate) / 100.0    // samples per 10 ms window
base := int(spw)                         // whole samples per window
frac := int((spw - float64(base)) * 1000) // thousandths, carried
carry := 0

nextWindow := func() int {
    n := base
    carry += frac
    if carry > 999 {
        n++
        carry -= 1000
    }
    return n
}
```

Is an integer thousandths accumulator the way *I* would distribute fractional samples? No. Is it the way the format works? Yes, so that is what the code does, with a comment pointing at the legacy function it mirrors. When you reproduce a format, faithfulness beats elegance every time; the moment to be clever was twenty years ago.

Same reasoning for the edges: the trailing partial window is silently discarded, because that is what the binary did. An "improvement" here would be a regression.

## Proving it: byte-identical tests

A rewrite like this is only as trustworthy as its proof of conformance. Two layers of tests carry that weight.

First, unit tests lock down every property of the algorithm with hand-computable cases, constant amplitude, clamping, and the discarded trailing window:

```
// 1 second of stereo 48 kHz at amplitude 25600
// -> 100 windows, each byte = 25600/256 = 100.
in := pcm(48000*2, 25600)
Compute(&out, bytes.NewReader(in), 48000, 2)
// assert: exactly 100 bytes, all equal to 100

// Full-scale 16-bit (32767) -> 32767/256 = 127.99 -> clamped to 127.

// 661 mono samples at 44100 (1.5 windows) -> 1 byte emitted,
// trailing 220 samples discarded, like the binary.
```

Second, the decisive one, golden-file comparisons: run the legacy binary and the Go implementation on the same real-world recordings, and `cmp` the outputs. Not "similar", not "within tolerance": identical, or the build fails. The first run of that comparison is also the best reverse-engineering tool you have, every divergence points at an assumption you got wrong, and the offset of the first differing byte usually tells you *which* one (an off-by-one drifts, a formula error is wrong from byte zero).

The input side is kept deterministic by feeding both implementations the same decoded PCM: the Go version shells out to ffmpeg (`-f s16le`, source sample rate and channel count) exactly as the old pipeline decoded, so the comparison isolates the profile computation from any decoder drift.

## Takeaways for your next legacy rewrite

- **The format is the behaviour, not the intent.** Write down the formula the old code actually computes, including its rounding, clamping and edge-case quirks, before touching the new implementation. Bug-for-bug compatibility is a feature.
- **Name the ghost.** Keep references to the legacy code in comments (`lProfilArrondi` lives on in ours). The next maintainer needs to know which decisions are load-bearing compatibility and which are free to change.
- **Byte-identical or it did not happen.** Golden files from the real legacy binary, compared with `cmp`, are cheap and end every argument. Tolerance-based comparisons hide drift.
- **Watch for slow divergence.** The dangerous bugs are the ones that agree for the first N windows. Test on long inputs and on rate/channel combinations with fractional windows, not just the convenient ones.

The Go daemon has since replaced the legacy binary in production. The `.ap` files it writes are indistinguishable from the old ones, which is precisely the point: the best migration is the one nobody downstream can detect.
