package main

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type RecordingSession struct {
	ID        string
	Timestamp time.Time
	Recording bool
	FileDiff  []string
}

var (
	ReaperAddress    string
	reaperProcess    *exec.Cmd
	mediaList        []string
	lastDiff         []string
	lock             sync.Mutex
	sessions         = make(map[string]*RecordingSession)
	currentSessionID string
)

//go:embed templates/*.gohtml
var templateFS embed.FS

func sessionsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "sessions.gohtml", sessions)
}

func main() {
	ReaperAddress = os.Getenv("REAPER_ADDRESS")
	if ReaperAddress == "" {
		ReaperAddress = "http://10.12.6.12:8080"
	}

	parsedTemplates, err := template.ParseFS(templateFS, "templates/*.gohtml")
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.SetHTMLTemplate(parsedTemplates)

	router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	router.GET("/start", start)
	router.GET("/status", status)
	router.GET("/stop", stop)
	router.GET("/files", files)
	router.GET("/sessions", sessionsHandler)

	router.Group("ui").
		GET("/start", startUI).
		GET("/stop", stopUI)

	router.Run(":8081")
}

const MediaGlob = "D:\\ReaperMedia\\*.wav"

func files(c *gin.Context) {
	sessionID := c.Query("session_id")
	var diff []string

	if sessionID != "" {
		if session, exists := sessions[sessionID]; exists {
			diff = session.FileDiff
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
	} else {
		diff = lastDiff
	}

	if reaperProcess != nil && reaperProcess.ProcessState == nil {
		// Reaper is recording, check for new files
		fileList := listFiles(MediaGlob)
		diff, _ = lo.Difference(fileList, mediaList)
	}

	c.JSON(http.StatusOK, diff)
}

type ReaperStatus struct {
	Heading      string
	ProcessState string
	Recording    bool
}

func start(c *gin.Context) {
	sessionID := uuid.New().String()
	session := &RecordingSession{
		ID:        sessionID,
		Timestamp: time.Now(),
		Recording: true,
	}
	currentSessionID = sessionID
	sessions[currentSessionID] = session

	err := startReaper()

	if err == errAlreadyStarted {
		c.JSON(http.StatusConflict, gin.H{"error": "Reaper already started"})
		return
	}

	if err == errUnknownOS {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown operating system"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start Reaper", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaper started", "session_id": sessionID})
}

func stop(c *gin.Context) {
	fileList := listFiles(MediaGlob)
	diff, _ := lo.Difference(fileList, mediaList)

	err := stopReaper()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to stop Reaper: %v", err)
		return
	}

	if session, exists := sessions[currentSessionID]; exists {
		session.FileDiff = diff
		session.Recording = false
	}

	lastDiff = diff
	c.JSON(http.StatusOK, lastDiff)
}
