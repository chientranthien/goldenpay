package metric

import (
	"net/http"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricDefaultHost = ":5050"
)

func ServeDefault() error {
	return Serve(metricDefaultHost)
}
func Serve(host string) error {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(host, nil)
	if err != nil {
		common.L().Errorw("metricListenAndServeErr", "host", host, "err", err)
	}

	return nil
}
