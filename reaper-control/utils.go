package main

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func isReaperOn(timeout time.Duration) bool {
	timeoutChan := time.After(timeout)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeoutChan:
			return false
		case <-tick:
			resp, err := http.Get(ReaperAddress + "/_/40667;TRANSPORT;")
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return true
				}
			}
		}
	}
}

func stopProcess(process *os.Process) error {
	if runtime.GOOS == "windows" {
		err := process.Kill()
		return err
	}

	return process.Signal(syscall.SIGINT)
}

func listFiles(pattern string) []string {
	files, _ := filepath.Glob(pattern)
	return files
}
