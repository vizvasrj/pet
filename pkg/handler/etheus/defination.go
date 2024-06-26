package etheus

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestCounter           *prometheus.CounterVec
	RequestDurationHistogram *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	requestDurationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code"},
	)

	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDurationHistogram)

	return &Metrics{
		RequestCounter:           requestCounter,
		RequestDurationHistogram: requestDurationHistogram,
	}
}
