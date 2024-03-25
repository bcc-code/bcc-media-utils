package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
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
	err := startReaper()
	if err != nil {
		errString := url.QueryEscape(err.Error())
		c.Redirect(http.StatusFound, "/status?err="+errString)
		return
	}

	c.Redirect(http.StatusFound, "/status")
}

func stopUI(c *gin.Context) {
	err := stopReaper()
	if err != nil {
		errString := url.QueryEscape(err.Error())
		c.Redirect(http.StatusFound, "/status?err="+errString)
		return
	}

	c.Redirect(http.StatusFound, "/status")
}
