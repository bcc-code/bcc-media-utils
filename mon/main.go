package main

import (
	"context"
	"fmt"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/durationpb"
)

func readTimeSeriesValue(projectID, serviceName, metricType string) (map[string]float64, error) {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	startTime := time.Now().UTC().Add(time.Minute * -3).Unix()
	endTime := time.Now().UTC().Unix()

	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   "projects/" + projectID,
		Filter: fmt.Sprintf("metric.type=\"%s\" resource.labels.service_name=\"%s\" ", metricType, serviceName),
		Aggregation: &monitoringpb.Aggregation{
			AlignmentPeriod: &durationpb.Duration{Seconds: 60},
			GroupByFields:   []string{"labels.response_code_class"},
		},
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{Seconds: startTime},
			EndTime:   &timestamp.Timestamp{Seconds: endTime},
		},
	}
	iter := c.ListTimeSeries(ctx, req)

	vals := make(map[string]float64)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read time series value, %w ", err)
		}
		vals[resp.Metric.Labels["response_code"]] = float64(resp.Points[0].GetValue().GetInt64Value()) / 60.0
	}

	return vals, nil
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	serviceName := os.Getenv("SERVICE_NAME")

	r := gin.Default()
	r.GET("/stats", func(c *gin.Context) {
		if c.Query("key") != os.Getenv("API_KEY") {
			c.AbortWithStatus(401)
		}

		data, err := readTimeSeriesValue(projectID, serviceName, "run.googleapis.com/request_count")
		if err != nil {
			panic(err)
		}

		c.JSON(200, data)
	})
}
