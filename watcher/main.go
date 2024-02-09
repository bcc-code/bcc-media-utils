package main

import (
	"context"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/rs/zerolog"
	"github.com/samber/lo/parallel"
)

type Watcher interface {
	Run(ctx context.Context)
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

	log.L.Info().Str("path", path).Dur("interval", interval).Msgf("Creating new watcher for %s", path)
	return &waitingWatcher{
		interval:    interval,
		path:        path,
		lastUpdated: time.Now(),
		callbackUrl: callbackUrl,
		fileSizes:   map[string]int64{},
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

	interval, err := strconv.Atoi(os.Getenv("WATCHER_INTERVAL"))
	if err != nil {
		interval = 10
	}

	parallel.ForEach(dirsToWatch, func(dir string, _ int) {
		var w Watcher
		w, err := newWatcher(dir, time.Second*time.Duration(interval), *callbackUrlString, *noWaitBool)
		if err != nil {
			panic(err)
		}
		w.Run(ctx)
	})
}
