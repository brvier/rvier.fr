---
title: QuickBorgmatic
section: opensource
weight: 215
image: ./images/quickborgmatic.png
alt: QuickBorgmatic panel showing per-repository backup freshness
link: https://github.com/brvier/QuickBorgmatic
linkText: View Project
---

An Omarchy shell bar-widget plugin (Quickshell) that watches the freshness of borgmatic backups on a remote server. The bar icon turns urgent as soon as the newest archive of any repository is older than the stale threshold, and the popup lists every repository with its last backup date and age. A repository borgmatic cannot list (wrong passphrase, missing repo, host down) gets its own highlighted row with borg's message while the others keep refreshing. Published on the Omarchy plugin marketplace.
