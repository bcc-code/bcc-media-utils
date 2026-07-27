package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/rs/zerolog"
	"github.com/samber/lo/parallel"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type Watcher interface {
	Run(ctx context.Context)
}

func envInt(name string, def int) int {
	v, err := strconv.Atoi(os.Getenv(name))
	if err != nil || v < 1 {
		return def
	}
	return v
}

func newWatcher(path string, interval time.Duration, callbackUrl string, noWait bool) (Watcher, error) {
	// Validate the pattern and the directory it points into before starting.
	if _, err := filepath.Glob(path); err != nil {
		return nil, err
	}
	if _, err := os.Stat(strings.Split(path, "*")[0]); err != nil {
		return nil, err
	}

	missingTicks := envInt("WATCHER_MISSING_TICKS", 3)

	if noWait {
		log.L.Info().Str("path", path).Dur("interval", interval).Msgf("Creating new no-wait watcher for %s", path)
		return &directWatcher{
			interval:      interval,
			path:          path,
			missingTicks:  missingTicks,
			callbackUrl:   callbackUrl,
			filesReported: map[string]int{},
		}, nil
	}

	stableTicks := envInt("WATCHER_STABLE_TICKS", 3)

	log.L.Info().
		Str("path", path).
		Dur("interval", interval).
		Int("stableTicks", stableTicks).
		Int("missingTicks", missingTicks).
		Msgf("Creating new watcher for %s", path)
	return &waitingWatcher{
		interval:      interval,
		path:          path,
		stableTicks:   stableTicks,
		missingTicks:  missingTicks,
		tracked:       map[string]*fileState{},
		filesReported: map[string]int{},
		callbackUrl:   callbackUrl,
	}, nil
}

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	watchDirsString := flag.String("dir", "", "directories to watch (comma-separated)")
	callbackUrlString := flag.String("callback", "", "callback url")
	noWaitBool := flag.Bool("no-wait", false, "do not wait for file to finish changing")

	flag.Parse()

	var dirsToWatch []string
	for _, dir := range strings.Split(*watchDirsString, ",") {
		if dir = strings.TrimSpace(dir); dir != "" {
			dirsToWatch = append(dirsToWatch, dir)
		}
	}
	if len(dirsToWatch) == 0 {
		log.L.Fatal().Msg("No directories to watch, pass -dir")
	}
	if *callbackUrlString == "" {
		log.L.Fatal().Msg("No callback url, pass -callback")
	}

	interval := envInt("WATCHER_INTERVAL", 10)

	// Create all watchers up front so a bad path or pattern fails the whole
	// process at startup instead of panicking inside a goroutine later.
	var watchers []Watcher
	for _, dir := range dirsToWatch {
		w, err := newWatcher(dir, time.Second*time.Duration(interval), *callbackUrlString, *noWaitBool)
		if err != nil {
			log.L.Fatal().Err(err).Str("dir", dir).Msg("Failed to create watcher")
		}
		watchers = append(watchers, w)
	}

	// Stop cleanly on SIGINT/SIGTERM (e.g. pod shutdown).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	parallel.ForEach(watchers, func(w Watcher, _ int) {
		w.Run(ctx)
	})
}
