package main

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var (
	ReaperAddress string
	reaperProcess *exec.Cmd
	mediaList     []string
	lock          sync.Mutex
)

//go:embed templates/*.gohtml
var templateFS embed.FS

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

	router.Group("ui").
		GET("/start", startUI).
		GET("/stop", stopUI)

	router.Run(":8081")
}

const MediaGlob = "D:\\ReaperMedia\\*.wav"

func files(c *gin.Context) {
	fileList := listFiles(MediaGlob)
	diff, _ := lo.Difference(fileList, mediaList)

	c.JSON(200, diff)
}

type ReaperStatus struct {
	Heading      string
	ProcessState string
	Recording    bool
}

func start(c *gin.Context) {
	err := startReaper()

	if err == errAlreadyStarted {
		c.String(http.StatusConflict, "Reaper already started")
		return
	}

	if err == errUnknownOS {
		c.String(http.StatusInternalServerError, "Unknown operating system")
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to start Reaper: %v", err)
		return
	}

	c.String(http.StatusOK, "Reaper started")
}

func stop(c *gin.Context) {
	fileList := listFiles(MediaGlob)
	diff, _ := lo.Difference(fileList, mediaList)

	err := stopReaper()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to stop Reaper: %v", err)
		return
	}

	c.JSON(http.StatusOK, diff)
}
