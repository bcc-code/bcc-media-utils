package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func validateKey(in string) bool {
	if in == "" {
		return false
	}
	return key == in
}

var key = ""

type data struct {
	Group    string
	Question string
}

func main() {
	r := gin.New()
	r.LoadHTMLGlob("templates/*")
	r.Use(cors.Default())
	cntChan := make(chan data, 10000)
	counters := map[string](map[string]int){}

	key = os.Getenv("KEY")

	bgCtx := context.Background()
	client, err := firestore.NewClient(bgCtx, "paske-2023")

	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)

	fsNeedsUpdate := true

	// This goroutine will update the firestore every second if there are changes to the counters
	// It writes the percentages and not the raw counts
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				if !fsNeedsUpdate {
					continue
				}
				fsNeedsUpdate = false
				coll := client.Collection("questions")

				// Calc totals
				totals := map[string]int{}
				for k, vals := range counters {
					if _, ok := totals[k]; !ok {
						totals[k] = 0
					}

					for _, val := range vals {
						totals[k] += val
					}

				}

				// Calc %
				percent := map[string](map[string]float64){}
				for k, vals := range counters {
					if _, ok := percent[k]; !ok {
						percent[k] = map[string]float64{}
					}

					for k2, val := range vals {
						if totals[k] == 0 {
							continue
						}
						percent[k][k2] = (float64(val) / float64(totals[k])) * 100
					}
				}

				for k, v := range percent {
					_, err := coll.Doc(k).Set(bgCtx, v)
					if err != nil {
						println(err.Error())
					} else {
						println("Updated FS")
					}
				}
			}
		}
	}()

	// This is the endpoint that the clients will call to increment the counter
	r.GET("/count/:group/:question", func(c *gin.Context) {
		key := data{
			Group:    c.Param("group"),
			Question: c.Param("question"),
		}
		cntChan <- key

		if !fsNeedsUpdate {
			println("setting needs update")
			fsNeedsUpdate = true
		}

		c.Status(http.StatusAccepted)
	})

	// This is the endpoint that the clients will call to get the current counters
	r.GET("/show", func(c *gin.Context) {
		c.JSON(http.StatusOK, counters)
	})

	// This renders a simple html page that can be used to reset the counters if one has the key
	r.GET("/reset", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reset.html", gin.H{})
	})

	// This is the endpoint that the clients will call to reset the counters
	// It is protected by a key that is set in the environment variable KEY
	r.POST("/reset", func(c *gin.Context) {
		if !validateKey(c.PostForm("key")) {
			c.Status(http.StatusUnauthorized)
			return
		}

		counters = map[string](map[string]int){}
		c.JSON(http.StatusOK, gin.H{"reset": "ok"})
	})

	go func() {
		for key := range cntChan {
			if _, ok := counters[key.Group]; !ok {
				counters[key.Group] = map[string]int{}
			}

			if _, ok := counters[key.Group][key.Question]; !ok {
				counters[key.Group][key.Question] = 0
			}

			counters[key.Group][key.Question] += 1
		}

	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
