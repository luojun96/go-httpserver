package metric

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	functionLatency = CreateExecutionTimeMetric("httpserver", "function execution time")
)

func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	histo := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_time",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 20),
		},
		[]string{"step"},
	)
	return histo
}

type ExectionTime struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}
