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

	"github.com/bcc-code/bccm-utils/watcher/db"
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

func envStr(name string, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}

func newWatcher(ctx context.Context, path string, interval time.Duration, callbackUrl string, noWait bool, queries *db.Queries) (Watcher, error) {
	// Validate the pattern and the directory it points into before starting.
	if _, err := filepath.Glob(path); err != nil {
		return nil, err
	}
	if _, err := os.Stat(strings.Split(path, "*")[0]); err != nil {
		return nil, err
	}

	store := newStore(queries, path)
	firstRun, err := store.Register(ctx)
	if err != nil {
		return nil, err
	}
	reported, err := store.Load(ctx)
	if err != nil {
		return nil, err
	}

	missingTicks := envInt("WATCHER_MISSING_TICKS", 3)

	if noWait {
		log.L.Info().
			Str("path", path).
			Dur("interval", interval).
			Int("reportedFromDB", len(reported)).
			Msgf("Creating new no-wait watcher for %s", path)
		if firstRun {
			// Everything present on the very first run counts as already
			// reported. On later restarts the persisted set decides instead, so
			// files that appeared while the watcher was down are reported.
			files, err := filepath.Glob(path)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				if err := store.Add(ctx, file); err != nil {
					return nil, err
				}
				reported[file] = 0
			}
		}
		return &directWatcher{
			interval:      interval,
			path:          path,
			missingTicks:  missingTicks,
			callbackUrl:   callbackUrl,
			filesReported: reported,
			store:         store,
		}, nil
	}

	stableTicks := envInt("WATCHER_STABLE_TICKS", 3)

	log.L.Info().
		Str("path", path).
		Dur("interval", interval).
		Int("stableTicks", stableTicks).
		Int("missingTicks", missingTicks).
		Int("reportedFromDB", len(reported)).
		Msgf("Creating new watcher for %s", path)
	return &waitingWatcher{
		interval:      interval,
		path:          path,
		stableTicks:   stableTicks,
		missingTicks:  missingTicks,
		tracked:       map[string]*fileState{},
		filesReported: reported,
		callbackUrl:   callbackUrl,
		store:         store,
	}, nil
}

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	watchDirsString := flag.String("dir", "", "directories to watch (comma-separated)")
	callbackUrlString := flag.String("callback", "", "callback url")
	noWaitBool := flag.Bool("no-wait", false, "do not wait for file to finish changing")
	dbPath := flag.String("db", envStr("WATCHER_DB_PATH", "watcher.db"), "path to the sqlite state database")

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

	// Stop cleanly on SIGINT/SIGTERM (e.g. pod shutdown).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// The reported-files set persists here so a restart neither re-reports
	// files already delivered nor misses files that appeared while down.
	sqldb, err := db.Open(ctx, *dbPath)
	if err != nil {
		log.L.Fatal().Err(err).Str("db", *dbPath).Msg("Failed to open state database")
	}
	defer sqldb.Close()
	queries := db.New(sqldb)

	// Create all watchers up front so a bad path or pattern fails the whole
	// process at startup instead of panicking inside a goroutine later.
	var watchers []Watcher
	for _, dir := range dirsToWatch {
		w, err := newWatcher(ctx, dir, time.Second*time.Duration(interval), *callbackUrlString, *noWaitBool, queries)
		if err != nil {
			log.L.Fatal().Err(err).Str("dir", dir).Msg("Failed to create watcher")
		}
		watchers = append(watchers, w)
	}

	parallel.ForEach(watchers, func(w Watcher, _ int) {
		w.Run(ctx)
	})
}
