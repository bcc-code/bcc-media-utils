package main

import (
	"errors"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var (
	errAlreadyStarted      = errors.New("Reaper already started")
	errUnknownOS           = errors.New("Unknown operating system")
	errFailedToStart       = errors.New("Failed to start Reaper")
	errRecodingFailed      = errors.New("Failed to start recording")
	errRecordingStopFailed = errors.New("Failed to stop recording")
	errFailedToStop        = errors.New("Failed to stop Reaper")
	errReaperNotStarted    = errors.New("Reaper not started")
)

func startReaper() error {
	lock.Lock()
	defer lock.Unlock()

	var cmd *exec.Cmd

	if reaperProcess != nil && !reaperProcess.ProcessState.Exited() {
		return errAlreadyStarted
	}

	mediaList = listFiles(MediaGlob)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("C:\\Program Files\\REAPER (x64)\\reaper.exe")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("/Applications/REAPER.app/Contents/MacOS/REAPER")
	}

	if cmd == nil {
		return errUnknownOS
	}

	err := cmd.Start()
	if err != nil {
		return errFailedToStart
	}

	reaperProcess = cmd

	if isReaperOn(40*time.Second) == false {
		return errRecodingFailed
	}

	_, err = http.Get(ReaperAddress + "/_/1013;TRANSPORT")
	if err != nil {
		return errRecodingFailed
	}

	return nil
}

func stopReaper() error {
	lock.Lock()
	defer lock.Unlock()

	if reaperProcess == nil {
		return errReaperNotStarted
	}

	_, err := http.Get(ReaperAddress + "/_/40667;TRANSPORT")
	if err != nil {
		return errRecordingStopFailed
	}

	if err := stopProcess(reaperProcess.Process); err != nil {
		return errFailedToStop
	}

	reaperProcess = nil

	return nil
}
