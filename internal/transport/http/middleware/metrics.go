package middleware

import (
	"net/http"

	"sub-hub/internal/observability/metrics"
)

func Metrics(m *metrics.HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return metrics.InstrumentHTTP(m, next)
	}
}
