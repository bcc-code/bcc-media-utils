package main

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var embededTemplates embed.FS

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

type avgData struct {
	Group string
	Value int64
}

type average struct {
	Total int64
	Count int
}

func (a *average) Add(val int64) {
	a.Total += val
	a.Count += 1
}

func (a *average) Get() float64 {
	if a.Count == 0 {
		return 0.0
	}

	return float64(a.Total) / float64(a.Count)
}

func main() {
	tmpl, err := template.ParseFS(embededTemplates, "templates/*.html")
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(cors.Default())
	r.SetHTMLTemplate(tmpl)
	cntChan := make(chan data, 10000)
	avgChan := make(chan avgData, 10000)
	counters := map[string](map[string]int){}
	averages := map[string]average{}
	percentCounters := map[string](map[string]float64){}

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

				percentCounters = percent

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

	r.GET("/average/:group/:value", func(c *gin.Context) {
		val, err := strconv.ParseInt(c.Param("value"), 10, 64)
		if err != nil || val < 0 || val > 18 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		key := avgData{
			Group: c.Param("group"),
			Value: val,
		}

		avgChan <- key
		c.Status(http.StatusAccepted)
	})

	// This is the endpoint that the clients will call to get the current counters
	r.GET("/show", func(c *gin.Context) {
		c.JSON(http.StatusOK, counters)
	})

	// This is the endpoint that the clients will call to get the current counters as percentages
	r.GET("/showpercent", func(c *gin.Context) {
		c.JSON(http.StatusOK, percentCounters)
	})

	r.GET("/showavg", func(c *gin.Context) {
		out := gin.H{}
		sumTotal := int64(0)
		sumCont := 0
		for k, v := range averages {
			out[k] = v.Get()
			sumTotal += v.Total
			sumCont += v.Count
		}

		out["total"] = 0
		if sumCont > 0 {
			out["total"] = float64(sumTotal) / float64(sumCont)
		}
		c.JSON(http.StatusOK, out)
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
		averages = map[string]average{}
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

	go func() {
		for key := range avgChan {
			if _, ok := averages[key.Group]; !ok {
				averages[key.Group] = average{}
			}

			a := averages[key.Group]
			a.Add(key.Value)
			averages[key.Group] = a
		}

	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
