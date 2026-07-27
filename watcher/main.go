package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	_, err := os.Stat(strings.Split(path, "*")[0])
	if err != nil {
		return nil, err
	}

	if noWait {
		log.L.Info().Str("path", path).Dur("interval", interval).Msgf("Creating new no-wait watcher for %s", path)
		return &directWatcher{
			interval:      interval,
			path:          path,
			callbackUrl:   callbackUrl,
			filesReported: make(map[string]struct{}),
		}, nil
	}

	stableTicks := envInt("WATCHER_STABLE_TICKS", 3)
	missingTicks := envInt("WATCHER_MISSING_TICKS", 3)

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

	dirsToWatch := strings.Split(*watchDirsString, ",")

	ctx := context.Background()

	interval := envInt("WATCHER_INTERVAL", 10)

	parallel.ForEach(dirsToWatch, func(dir string, _ int) {
		var w Watcher
		w, err := newWatcher(dir, time.Second*time.Duration(interval), *callbackUrlString, *noWaitBool)
		if err != nil {
			panic(err)
		}
		w.Run(ctx)
	})
}
