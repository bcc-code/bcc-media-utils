package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
)

// fileState tracks the last observed size/mtime of a file and how many
// consecutive polls it has remained unchanged.
type fileState struct {
	size        int64
	modTime     time.Time
	stableTicks int
}

// waitingWatcher is a file watcher that waits for a file to stop being written to
// before sending a notification.
type waitingWatcher struct {
	path          string
	interval      time.Duration
	stableTicks   int                   // consecutive unchanged polls required before reporting
	missingTicks  int                   // consecutive missing polls before a reported file is forgotten
	tracked       map[string]*fileState // observed, not yet reported
	filesReported map[string]int        // reported file -> consecutive missing-tick count
	callbackUrl   string
	store         reportedStore
}

func (w *waitingWatcher) doWatch(ctx context.Context) {
	files, err := filepath.Glob(w.path)
	if err != nil {
		log.L.Error().Err(err).Str("path", w.path).Send()
		return
	}

	seen := map[string]struct{}{}
	for _, file := range files {
		seen[file] = struct{}{}
		if _, reported := w.filesReported[file]; reported {
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
		size := stats.Size()
		if size == 0 {
			continue
		}
		st, ok := w.tracked[file]
		if !ok {
			w.tracked[file] = &fileState{size: size, modTime: stats.ModTime()}
			log.L.Debug().Str("file", file).Int64("size", size).Msg("Tracking new file")
			continue
		}
		if size != st.size || !stats.ModTime().Equal(st.modTime) {
			st.size = size
			st.modTime = stats.ModTime()
			st.stableTicks = 0
			log.L.Debug().Str("file", file).Int64("size", size).Msg("File still changing")
			continue
		}
		st.stableTicks++
		if st.stableTicks < w.stableTicks {
			log.L.Debug().Str("file", file).Int("stableTicks", st.stableTicks).Msg("File unchanged")
			continue
		}
		log.L.Info().Str("file", stats.Name()).Int64("size", size).Time("modTime", stats.ModTime()).Msg("File stable, reporting")
		if err := postCallback(w.callbackUrl, file, stats); err != nil {
			// Keep the file tracked so the next tick retries the callback.
			log.L.Error().Err(err).Str("file", stats.Name()).Msg("Callback failed, will retry next tick")
			continue
		}
		delete(w.tracked, file)
		if err := w.store.Add(ctx, file); err != nil {
			// The in-memory map stays authoritative; worst case is a duplicate
			// callback after a restart, which receivers must tolerate anyway.
			log.L.Error().Err(err).Str("file", file).Msg("Failed to persist reported file")
		}
		w.filesReported[file] = 0
	}

	// Drop tracked files that vanished before ever stabilizing.
	for file := range w.tracked {
		if _, ok := seen[file]; !ok {
			delete(w.tracked, file)
			log.L.Debug().Str("file", file).Msg("Tracked file disappeared, dropping")
		}
	}

	// A reported file is only forgotten (and thus eligible for re-reporting)
	// after it has been confirmed missing for several consecutive polls, so a
	// transient NFS blip or rename-in-place cannot cause a duplicate callback.
	for file, missing := range w.filesReported {
		_, err := os.Stat(file)
		switch {
		case err == nil:
			w.filesReported[file] = 0
		case os.IsNotExist(err):
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
		default:
			log.L.Warn().Err(err).Str("file", file).Msg("stat failed for reported file")
		}
	}
}

func (w *waitingWatcher) Run(ctx context.Context) {
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
