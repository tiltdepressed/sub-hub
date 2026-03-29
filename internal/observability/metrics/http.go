package metrics

import (
	"net/http"
	"time"
)

type statusCapturingWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCapturingWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCapturingWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func InstrumentHTTP(m *HTTPMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		cw := &statusCapturingWriter{ResponseWriter: w}
		next.ServeHTTP(cw, r)

		status := cw.Status()
		path := r.URL.Path
		m.RequestsTotal.WithLabelValues(r.Method, path, http.StatusText(status)).Inc()
		m.Duration.WithLabelValues(r.Method, path, http.StatusText(status)).Observe(time.Since(start).Seconds())
	})
}
