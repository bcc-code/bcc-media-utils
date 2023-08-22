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
	"github.com/samber/lo"
)

type watcher struct {
	path            string
	interval        time.Duration
	recentlyUpdated []string
	lastUpdated     time.Time
	fileSizes       map[string]int64
	filesReported   []string
	callbackUrl     string
}

func newWatcher(path string, interval time.Duration, callbackUrl string) *watcher {
	return &watcher{
		interval:    interval,
		path:        path,
		lastUpdated: time.Now(),
		callbackUrl: callbackUrl,
		fileSizes:   map[string]int64{},
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
		if stats.IsDir() || strings.HasPrefix(stats.Name(), ".") || lo.Contains(w.filesReported, file) {
			continue
		}
		size := stats.Size()
		oldSize, ok := w.fileSizes[file]
		if !ok {
			w.fileSizes[file] = size
			continue
		}
		if oldSize < size {
			if !lo.Contains(w.recentlyUpdated, file) {
				w.recentlyUpdated = append(w.recentlyUpdated, file)
			}
			w.fileSizes[file] = size
			continue
		}
		w.recentlyUpdated = lo.Filter(w.recentlyUpdated, func(i string, _ int) bool {
			return i != file
		})
		delete(w.fileSizes, file)
		w.filesReported = append(w.filesReported, file)
		w.fileUpdated(file, stats)
	}

	for _, file := range w.filesReported {
		_, err := os.Stat(file)
		if err == nil {
			continue
		}
		if os.IsNotExist(err) {
			w.filesReported = lo.Filter(w.filesReported, func(i string, _ int) bool {
				return i != file
			})
			continue
		}
		log.L.Error().Err(err).Send()
		return
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
