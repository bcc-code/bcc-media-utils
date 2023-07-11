package main

import (
	"context"
	"flag"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/rs/zerolog"
	"github.com/samber/lo/parallel"
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

	parallel.ForEach(dirsToWatch, func(dir string, _ int) {
		w := newWatcher(dir, time.Second*5, *callbackUrlString)
		w.run(ctx)
	})
}
