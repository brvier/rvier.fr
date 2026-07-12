---
title: 'tmpfs and /dev/shm: the ramdisk, the forgotten optimization'
date: '2026-06-19'
lang: en
description: 'A RAM-backed filesystem ships with every Linux machine and almost nobody uses it. Three real production uses of tmpfs and /dev/shm: temp files for third-party binaries, in-RAM indexes shared between processes, and monitoring via statfs.'
ogDescription: Three real production uses of tmpfs and /dev/shm, and the pitfalls to know about.
keywords: Linux, tmpfs, /dev/shm, ramdisk, performance, Go
summary: 'Three real production uses of the ramdisk: temp files for third-party binaries, in-RAM indexes shared between processes, monitoring via statfs, and the pitfalls to know about.'
---

When a service is slow, we optimize the SQL query, add a Redis cache, profile the code. We almost always forget that every Linux machine ships with an entirely RAM-backed filesystem, already mounted, requiring zero installation: `tmpfs`. I use it in production in three different services, for three different reasons. Here is a quick tour.

## A 30-second refresher

`tmpfs` is a filesystem whose data lives in the page cache: reads and writes at RAM speed, no SSD wear, contents gone at reboot. Every Linux distribution already mounts one at `/dev/shm`, sized by default at half the RAM. You can also mount your own, with a bounded size:

```
# /etc/fstab
tmpfs  /mnt/ramdisk  tmpfs  size=8G,mode=1777  0  0
```

The point of the bound: unlike `/dev/shm`, which is shared with everything on the machine, a dedicated mount guarantees that a runaway temp file will not eat all the memory.

## Use 1: third-party binaries that demand local files

Our audio fingerprinting workers drive a binary that only knows how to work on local file paths, no URLs, no stdin streaming. The master files, however, arrive over HTTP. The naive solution downloads them to `/tmp` on disk; on 256-core machines processing dozens of tasks in parallel, the writes plus immediate re-reads saturate the I/O for files whose lifetime is a few seconds.

Downloading to a ramdisk (`/mnt/ramdisk`, a configurable path in the worker) removes the problem entirely: the file only ever exists in RAM, the third-party binary accesses it like any other file, and everything is deleted when the task ends. No change to the binary, no library: just a different path in the config. It is the perfect bridge between "my tool wants a file" and "I refuse to pay for the disk".

## Use 2: in-RAM indexes shared between processes, with file semantics

Our audio recognition platform loads fingerprint catalogs (several gigabytes) into `/dev/shm`. The latency win is obvious, but the most elegant benefit is elsewhere: the RAM becomes *discoverable*. A worker starting up does not need a central registry to know which catalogs are already loaded on its machine, it lists the directory:

```
func ScanLoadedBuckets(shmPath string) []string {
    entries, _ := os.ReadDir(shmPath)
    // each subdirectory = one catalog loaded in RAM
    ...
}
```

The whole file toolbox works on shared memory: `ls` to debug, `du` to measure, `rm` to unload an index. Try doing that with a raw POSIX shared memory segment.

## Use 3: monitor it like a disk

A full ramdisk behaves like a full disk, except it is your RAM. So our workers expose `/dev/shm` usage in their stats endpoint, and it is three lines of Go:

```
var stat syscall.Statfs_t
syscall.Statfs("/dev/shm", &stat)
total := stat.Blocks * uint64(stat.Bsize)
free := stat.Bavail * uint64(stat.Bsize)
```

Those values feed our metrics and fire an alert before the next index load fails. One `statfs`, no agent, no dependency.

## The pitfalls

- **It is RAM.** A full tmpfs plus memory-hungry processes = OOM killer. Always bound the mount size, and size it for the real workload, not the nominal case.
- **Nothing survives a reboot.** Only put reconstructible data there: caches, temporary downloads, reloadable indexes.
- **Swap exists.** Under memory pressure, tmpfs pages can be swapped out, your "RAM" becomes a disk again, only worse. Monitor it, and keep headroom.
- **systemd may isolate you.** With `PrivateTmp=true`, the `/tmp` your service sees is not the one other services see. For a ramdisk shared between services, use a dedicated mount rather than `/tmp`.
- **Clean up.** A file forgotten in tmpfs is a silent memory leak. `defer os.Remove(...)` per task, plus a periodic safety net that purges orphaned files.

A ramdisk replaces neither an application cache nor a real database: it is a utility tool. But when the problem is "high-frequency temporary files" or "a third-party binary that wants local paths", it is often the simplest solution, the fastest to put in place, and the one nobody thinks of.
