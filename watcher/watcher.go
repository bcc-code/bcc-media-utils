package main

import (
	"context"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/samber/lo"
	"net/http"
	"os"
	"strings"
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

func listDirectory(path string) []string {
	var result []string
	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if !entry.IsDir() {
			result = append(result, entry.Name())
		} else {
			result = append(result, listDirectory(path+"/"+entry.Name())...)
		}
	}
	return lo.Map(result, func(i string, _ int) string {
		return path + "/" + i
	})
}

func (w *watcher) doWatch() {
	log.L.Debug().Str("path", w.path).Msg("Iterating through folder")
	files := listDirectory(w.path)

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
			w.fileUpdated(file)
		}
	}
	w.lastUpdated = time.Now()
}

func (w *watcher) fileUpdated(file string) {
	log.L.Debug().Str("file", file).Msg("File updated!")

	if w.callbackUrl != "" {
		callbackUrl := strings.ReplaceAll(w.callbackUrl, "{file}", file)

		_, err := http.Get(callbackUrl)
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
