package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
)

type directWatcher struct {
	path          string
	interval      time.Duration
	filesReported map[string]struct{}
	callbackUrl   string
}

func (w *directWatcher) doWatch() {
	files, err := filepath.Glob(w.path)
	if err != nil {
		log.L.Error().Err(err).Send()
		return
	}

	for _, f := range files {
		if _, found := w.filesReported[f]; !found {
			stats, err := os.Stat(f)
			if err != nil {
				log.L.Error().Err(err).Send()
				return
			}
			if _, found := w.filesReported[f]; found {
				continue
			}

			if stats.IsDir() || strings.HasPrefix(stats.Name(), ".") {
				continue
			}

			w.fileUpdated(f, stats)
			w.filesReported[f] = struct{}{}
		}
	}
}

func (w *directWatcher) Run(ctx context.Context) {
	files, err := filepath.Glob(w.path)
	if err != nil {
		log.L.Error().Err(err).Send()
		return
	}

	for _, file := range files {
		w.filesReported[file] = struct{}{}
	}

	ticker := time.NewTicker(w.interval)
	for {
		select {
		case <-ticker.C:
			w.doWatch()
		case <-ctx.Done():
			return
		}
	}
}

// TODO: This is a direct copy of the waitingWatcher.fileUpdated method. It should be refactored to be shared.
func (w *directWatcher) fileUpdated(path string, file os.FileInfo) {
	log.L.Debug().Str("file", file.Name()).Msg("File updated!")

	if w.callbackUrl != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.L.Error().Err(err).Send()
		}

		str, _ := json.Marshal(callbackRequest{
			Name:      file.Name(),
			Size:      file.Size(),
			Path:      absPath,
			UpdatedAt: file.ModTime(),
		})

		_, err = http.Post(w.callbackUrl, "application/json", bytes.NewReader(str))
		if err != nil {
			log.L.Error().Err(err).Send()
		}
	}
}
