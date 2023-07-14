package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bcc-code/bccm-utils/scheduler/cantemo"
	"github.com/bcc-code/bccm-utils/scheduler/jobs"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	client := cantemo.New(os.Getenv("CANTEMO_URL"), os.Getenv("CANTEMO_AUTH_TOKEN"))

	q := jobs.NewQueue(5, func(id string) error {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		return nil
	})

	go q.Run(context.Background())

	r := gin.Default()
	r.GET("jobs", func(ctx *gin.Context) {
		ctx.JSON(200, q)
	})
	r.POST("jobs", jobPostHandler(client, q))

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

		for _, id := range ids {
			queue.Queue(id)
		}

		ctx.JSON(200, ids)
	}
}

func getItemIDsFromSearchQuery(client *cantemo.Client, query string) ([]string, error) {
	res, err := client.Search().Put(query, 1)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, r := range res.Results {
		ids = append(ids, r.VidispineId)
	}

	for res.Page < res.Pages {
		res, err = client.Search().Put(query, res.Page+1)
		if err != nil {
			log.L.Error().Err(err).Send()
		}
		for _, r := range res.Results {
			ids = append(ids, r.VidispineId)
		}
	}

	return ids, nil
}
