package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/samber/lo"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type watcher struct {
	path            string
	interval        time.Duration
	recentlyUpdated []string
	lastUpdated     time.Time
	callbackUrl     string
}

func newWatcher(path string, interval time.Duration, callbackUrl string) *watcher {
	return &watcher{
		interval:    interval,
		path:        path,
		lastUpdated: time.Now(),
		callbackUrl: callbackUrl,
	}
}

func (w *watcher) doWatch() {
	log.L.Debug().Str("path", w.path).Msg("Iterating through folder")
	files, err := filepath.Glob(w.path)
	if err != nil {
		log.L.Error().Err(err).Send()
		return
	}

	for _, file := range files {
		stats, err := os.Stat(file)
		if err != nil {
			log.L.Error().Err(err).Send()
			return
		}
		if stats.ModTime().After(w.lastUpdated) {
			if !lo.Contains(w.recentlyUpdated, file) {
				w.recentlyUpdated = append(w.recentlyUpdated, file)
			}
			continue
		}
		if lo.Contains(w.recentlyUpdated, file) {
			w.recentlyUpdated = lo.Filter(w.recentlyUpdated, func(i string, _ int) bool {
				return i != file
			})
			w.fileUpdated(file, stats)
		}
	}
	w.lastUpdated = time.Now()
}

type callbackRequest struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	UpdatedAt time.Time `json:"updatedAt"`
	Size      int64     `json:"size"`
}

func (w *watcher) fileUpdated(path string, file os.FileInfo) {
	log.L.Debug().Str("file", file.Name()).Msg("File updated!")

	if w.callbackUrl != "" {
		str, _ := json.Marshal(callbackRequest{
			Name:      file.Name(),
			Size:      file.Size(),
			Path:      path,
			UpdatedAt: file.ModTime(),
		})
		_, err := http.Post(w.callbackUrl, "application/json", bytes.NewReader(str))
		if err != nil {
			log.L.Error().Err(err).Send()
		}
	}
}

func (w *watcher) run(ctx context.Context) {
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
