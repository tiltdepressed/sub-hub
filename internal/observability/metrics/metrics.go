package metrics

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetrics struct {
	RequestsTotal *prometheus.CounterVec
	Duration      *prometheus.HistogramVec
}

func NewRegistry() (*prometheus.Registry, *HTTPMetrics) {
	r := prometheus.NewRegistry()
	m := &HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
	}

	r.MustRegister(m.RequestsTotal, m.Duration)
	return r, m
}
