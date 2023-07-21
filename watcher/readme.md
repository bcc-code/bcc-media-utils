# Watcher

Iterates through specified folders for file changes. 

Whenever a change is detected, a webhook is triggered.

## Cache

Under certain circumstances the results of the `FileStat()` call are cached on
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
