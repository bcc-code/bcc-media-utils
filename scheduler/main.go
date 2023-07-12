package main

import (
	"crypto/tls"
	"fmt"
	"github.com/bcc-code/bccm-utils/scheduler/cantemo"
	"github.com/bcc-code/bccm-utils/scheduler/jobs"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	client := cantemo.New(os.Getenv("CANTEMO_URL"), os.Getenv("CANTEMO_AUTH_TOKEN"))

	queue := jobs.NewQueue()

	r := gin.Default()
	r.GET("jobs", func(ctx *gin.Context) {
		j, _ := queue.GetJobs()
		ctx.JSON(200, j)
	})
	r.POST("jobs", jobPostHandler(client, queue))

	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.L.Panic().Err(err).Send()
	}
}

func jobPostHandler(client *cantemo.Client, queue *jobs.Queue) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := ctx.Request.Body
		q, _ := io.ReadAll(body)
		ids, err := getItemIDsFromSearchQuery(client, string(q))
		if err != nil {
			log.L.Error().Err(err).Send()
			ctx.Status(500)
			return
		}

		iterateThroughIDs(client, ids)

		ctx.JSON(200, ids)
	}
}
