package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func status(c *gin.Context) {
	status := ReaperStatus{}
	status.Recording = false

	if reaperProcess != nil && reaperProcess.ProcessState == nil {
		status.ProcessState = "Running"
		status.Recording = true
	} else {
		status.ProcessState = "Not Running"
	}

	c.HTML(http.StatusOK, "status.gohtml", status)
}

func startUI(c *gin.Context) {
	sessionID := uuid.New().String()
	session := &RecordingSession{
		ID:        sessionID,
		Timestamp: time.Now(),
		Status:    "Recording",
	}
	currentSessionID = sessionID
	sessions[currentSessionID] = session

	err := startReaper()
	if err != nil {
		errString := url.QueryEscape(err.Error())
		c.Redirect(http.StatusFound, "/status?err="+errString)
		return
	}

	c.Redirect(http.StatusFound, "/status")
}

func stopUI(c *gin.Context) {
	fileList := listFiles(MediaGlob)
	diff, _ := lo.Difference(fileList, mediaList)

	if session, exists := sessions[currentSessionID]; exists {
		session.FileDiff = diff
		session.Status = "Stopped"
	}

	lastDiff = diff

	err := stopReaper()
	if err != nil {
		errString := url.QueryEscape(err.Error())
		c.Redirect(http.StatusFound, "/status?err="+errString)
		return
	}

	c.Redirect(http.StatusFound, "/status")
}
