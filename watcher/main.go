package main

import (
	"context"
	"flag"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/rs/zerolog"
	"github.com/samber/lo/parallel"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	watchDirsString := flag.String("dir", "", "directories to watch (comma-separated)")
	callbackUrlString := flag.String("callback", "", "callback url")

	flag.Parse()

	dirsToWatch := strings.Split(*watchDirsString, ",")

	ctx := context.Background()

	interval, err := strconv.Atoi(os.Getenv("WATCHER_INTERVAL"))
	if err != nil {
		interval = 10
	}

	parallel.ForEach(dirsToWatch, func(dir string, _ int) {
		w := newWatcher(dir, time.Second*time.Duration(interval), *callbackUrlString)
		w.run(ctx)
	})
}
