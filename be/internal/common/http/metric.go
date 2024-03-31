package http

import (
	"github.com/chientranthien/goldenpay/internal/common/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	defaultLabels = []string{"api", "code"}

	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "keep track of request count which can be used to keep track of QPS",
	}, defaultLabels)

	requestLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_latency",
		Buckets: metric.DefaultMillisBucket,
	}, defaultLabels)
)
