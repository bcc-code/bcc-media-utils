package main

import (
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var ReaperAddress string

func main() {
	ReaperAddress = os.Getenv("REAPER_ADDRESS")
	if ReaperAddress == "" {
		ReaperAddress = "http://10.12.6.12:8080"
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	router.GET("/start", start)
	router.GET("/stop", stop)

	/// *** LEGACY CODE *** ///
	router.GET("/startReaper", startReaper)
	router.GET("/stopReaper", stopReaper)
	router.GET("/startReaperRecording", startReaperRecording)
	router.GET("/stopReaperRecording", stopReaperRecording)
	// *** END LEGACY CODE *** ///

	router.Run(":8081")
}

var reaperProcess *os.Process
var mediaList []string

const MediaGlob = "D:\\ReaperMedia\\*.wav"

func start(c *gin.Context) {
	var cmd *exec.Cmd

	mediaList = listFiles(MediaGlob)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("C:\\Program Files\\REAPER (x64)\\reaper.exe")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("/Applications/REAPER.app/Contents/MacOS/REAPER")
	}

	if cmd == nil {
		c.String(500, "Unsupported operating system")
		return
	}

	err := cmd.Start()
	if err != nil {
		c.String(500, "Failed to start Reaper: %v", err)
		return
	}
	reaperProcess = cmd.Process

	if isReaperOn(40*time.Second) == false {
		c.String(500, "Error requesting URL: %v", err)
		return
	}

	_, err = http.Get(ReaperAddress + "/_/1013;TRANSPORT")
	if err != nil {
		c.String(500, "Error requesting URL: %v", err)
		return
	}

	c.String(200, "Reaper started")
}

func stop(c *gin.Context) {
	_, err := http.Get(ReaperAddress + "/_/40667;TRANSPORT")
	if err != nil {
		c.String(500, "Error requesting URL: %v", err)
		return
	}

	fileList := listFiles(MediaGlob)
	diff, _ := lo.Difference(fileList, mediaList)

	if reaperProcess == nil {
		c.String(500, "Reaper not started")
		return
	}

	if err := stopProcess(reaperProcess); err != nil {
		c.String(500, "Failed to stop Reaper: %v", err)
		return
	}

	c.JSON(200, diff)
}

//// *** LEGACY CODE *** ////

func startReaper(c *gin.Context) {
	var cmd *exec.Cmd

	mediaList = listFiles(MediaGlob)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("C:\\Program Files\\REAPER (x64)\\reaper.exe")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("/Applications/REAPER.app/Contents/MacOS/REAPER")
	}

	if cmd == nil {
		c.String(500, "Unsupported operating system")
		return
	}

	err := cmd.Start()
	if err != nil {
		c.String(500, "Failed to start Reaper: %v", err)
		return
	}
	reaperProcess = cmd.Process

	if isReaperOn(40*time.Second) == false {
		c.String(500, "Error requesting URL: %v", err)
		return
	}

	c.String(200, "Reaper started")
}

func stopReaper(c *gin.Context) {
	if reaperProcess == nil {
		c.String(500, "Reaper not started")
		return
	}

	if err := stopProcess(reaperProcess); err != nil {
		c.String(500, "Failed to stop Reaper: %v", err)
		return
	}

	c.String(200, "Reaper stopped")
}

func startReaperRecording(c *gin.Context) {
	_, err := http.Get(ReaperAddress + "/_/1013;TRANSPORT")
	if err != nil {
		c.String(500, "Error requesting URL: %v", err)
		return
	}

	c.String(200, "Reaper recording started")
}

func stopReaperRecording(c *gin.Context) {
	_, err := http.Get(ReaperAddress + "/_/40667;TRANSPORT")
	if err != nil {
		c.String(500, "Error requesting URL: %v", err)
		return
	}
	c.String(200, "Reaper recording stopped")
}
