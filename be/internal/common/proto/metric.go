package proto

import (
	"github.com/chientranthien/goldenpay/internal/common/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	defaultLabels = []string{"api", "code"}

	serverRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_server_request_count",
		Help: "keep track of request count which can be used to keep track of QPS",
	}, defaultLabels)

	serverRequestLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_server_latency",
		Buckets: metric.DefaultMillisBucket,
	}, defaultLabels)

	clientRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_client_request_count",
		Help: "keep track of request count which can be used to keep track of QPS",
	}, defaultLabels)

	clientRequestLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_client_latency",
		Buckets: metric.DefaultMillisBucket,
	}, defaultLabels)
)
