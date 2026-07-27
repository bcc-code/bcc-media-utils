# Watcher

Iterates through specified folders for file changes. 

Whenever a change is detected, a webhook is triggered.

## Configuration

- `WATCHER_INTERVAL` — poll interval in seconds (default `10`).
- `WATCHER_STABLE_TICKS` — number of consecutive polls a file's size and mtime
  must remain unchanged before it is considered done and reported (default `3`,
  i.e. 30 seconds of quiet at the default interval). Only applies to the
  default (waiting) mode.
- `WATCHER_MISSING_TICKS` — number of consecutive polls an already-reported
  file must be confirmed missing before it is forgotten and becomes eligible
  for re-reporting (default `3`). This prevents transient NFS errors or
  rename-in-place from causing duplicate callbacks, while a genuinely deleted
  and recreated file is still reported again.

## Delivery semantics

Callbacks are delivered at-least-once: if the callback endpoint is unreachable
or returns a non-2xx status, the file is not marked as reported and the
notification is retried on the next poll. The receiver should therefore treat
callbacks idempotently (the reported `path`, `size` and `updatedAt` can be used
for deduplication).

## Cache

Under certain circumstances the results of the `File.Stat()` call are cached on
the OS/FS level for a relatively long time. This is usually desirable to improve
the performance but in case of a watcher like this it is not desirable as we will
not pick up on any changes.

In case of monitoring the files on an NFS mounted volume, it seems that it is
enough to mount the file system using the `noac` (no attribute cache) option.
Testing in our system shows a cache time of ~4 seconds somewhere in the chain,
but that is sufficiently short for the use case. More information about this at
https://stackoverflow.com/a/35162336/556085

In order to avoid issues it is recommended to validate the `mtime` and file size
when receiving the callback.
