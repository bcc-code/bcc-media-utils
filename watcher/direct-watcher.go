package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
)

// directWatcher reports new files immediately, without waiting for them to
// stop changing. Files already present at startup are not reported.
type directWatcher struct {
	path          string
	interval      time.Duration
	missingTicks  int            // consecutive missing polls before a reported file is forgotten
	filesReported map[string]int // reported file -> consecutive missing-tick count
	callbackUrl   string
	store         reportedStore
}

func (w *directWatcher) doWatch(ctx context.Context) {
	files, err := filepath.Glob(w.path)
	if err != nil {
		log.L.Error().Err(err).Str("path", w.path).Send()
		return
	}

	seen := map[string]struct{}{}
	for _, file := range files {
		seen[file] = struct{}{}
		if _, reported := w.filesReported[file]; reported {
			w.filesReported[file] = 0
			continue
		}
		stats, err := os.Stat(file)
		if err != nil {
			log.L.Warn().Err(err).Str("file", file).Msg("stat failed, skipping this tick")
			continue
		}
		if stats.IsDir() || strings.HasPrefix(stats.Name(), ".") {
			continue
		}
		log.L.Info().Str("file", stats.Name()).Int64("size", stats.Size()).Msg("New file, reporting")
		if err := postCallback(w.callbackUrl, file, stats); err != nil {
			// Not marked as reported, so the next tick retries the callback.
			log.L.Error().Err(err).Str("file", stats.Name()).Msg("Callback failed, will retry next tick")
			continue
		}
		if err := w.store.Add(ctx, file); err != nil {
			// The in-memory map stays authoritative; worst case is a duplicate
			// callback after a restart, which receivers must tolerate anyway.
			log.L.Error().Err(err).Str("file", file).Msg("Failed to persist reported file")
		}
		w.filesReported[file] = 0
	}

	// Forget reported files only after several consecutive missing polls, so a
	// transient NFS blip cannot cause a duplicate callback while genuinely
	// deleted files do not grow the map forever.
	for file, missing := range w.filesReported {
		if _, ok := seen[file]; ok {
			continue
		}
		missing++
		if missing >= w.missingTicks {
			delete(w.filesReported, file)
			if err := w.store.Remove(ctx, file); err != nil {
				log.L.Error().Err(err).Str("file", file).Msg("Failed to remove reported file from store")
			}
			log.L.Info().Str("file", file).Msg("Reported file gone, will report again if recreated")
		} else {
			w.filesReported[file] = missing
		}
	}
}

func (w *directWatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.doWatch(ctx)
		case <-ctx.Done():
			return
		}
	}
}
