package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/ansel1/merry/v2"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func importIntoBQ(source string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return merry.Wrap(err)
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference(source)
	gcsRef.SourceFormat = bigquery.JSON
	gcsRef.Schema = schema
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteAppend

	job, err := loader.Run(ctx)
	if err != nil {
		return merry.Wrap(err)
	}

	log.L.Info().
		Str("JobID", job.ID()).
		Msg("Started import job")

	status, err := job.Wait(ctx)
	if err != nil {
		return merry.Wrap(err)
	}

	return merry.Wrap(status.Err())
}

type NPAWWebhookData struct {
	Title        string `json:"title"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	URL          string `json:"url"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func handleRequest(c *gin.Context) {
	webhookData := &NPAWWebhookData{}
	err := c.BindJSON(webhookData)

	if err != nil {
		log.L.Error().
			Err(err).
			Msg("NPAW reported error in report")
		return
	}

	if webhookData.ErrorCode != "" {
		log.L.Error().
			Str("NPAW Error code", webhookData.ErrorCode).
			Str("NPAW Error message", webhookData.ErrorMessage).
			Msg("NPAW reported error in report")
		return
	}

	source := fmt.Sprintf("gs://%s/Report-bccmedia-%s-%s-%s.json.gz", gcpBucketID, webhookData.Title, webhookData.StartDate, webhookData.EndDate)
	err = importIntoBQ(source)
	if err != nil {
		log.L.Error().
			Err(err).
			Msg("Import failed")
		return
	}

	// Send heartbeat to Betteruptime
	if heartbeatURL != "" {
		_, _ = http.Get(heartbeatURL)
	}
}

var schema bigquery.Schema
var (
	projectID    string
	datasetID    string
	tableID      string
	gcpBucketID  string
	npawID       string
	heartbeatURL string
)

func main() {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)

	dat, err := os.ReadFile("schema.json")
	if err != nil {
		log.L.Panic().Err(err).Msg("Unable to read schema file")
	}

	schemaR, err := bigquery.SchemaFromJSON(dat)
	if err != nil {
		log.L.Panic().Err(err).Msg("Unable to parse schema")
	}

	schema = schemaR

	projectID = os.Getenv("GCP_PROJECT")
	if projectID == "" {
		log.L.Panic().Msg("Empty GCP_PROJECT")
	}

	datasetID = os.Getenv("DATASET")
	if datasetID == "" {
		panic("Empty DATASET")
	}

	tableID = os.Getenv("TABLE")
	if tableID == "" {
		log.L.Panic().Msg("Empty TABLE")
	}

	gcpBucketID = os.Getenv("BUCKET")
	if gcpBucketID == "" {
		log.L.Panic().Msg("Empty BUCKET")
	}

	npawID = os.Getenv("NPAWID")
	if npawID == "" {
		log.L.Panic().Msg("Empty NPAWID")
	}

	heartbeatURL = os.Getenv("HEARTBEATURL")

	c := gin.Default()
	c.POST("/", handleRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	c.Run(":" + port)
}
