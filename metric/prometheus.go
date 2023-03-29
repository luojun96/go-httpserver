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

func NewTimer() *ExectionTime {
	return NewExecutionTime(functionLatency)
}

func NewExecutionTime(histo *prometheus.HistogramVec) *ExectionTime {
	return &ExectionTime{
		histo: histo,
		start: time.Now(),
		last:  time.Now(),
	}
}

func (t *ExectionTime) Observe() {
	t.histo.WithLabelValues("total").Observe(time.Since(t.start).Seconds())
}
