package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"sub-hub/internal/observability/metrics"
	"sub-hub/internal/transport/http/handlers"
	appmw "sub-hub/internal/transport/http/middleware"
	"sub-hub/internal/usecase/subscriptions"
)

type RouterDeps struct {
	Log         *zap.Logger
	Health      handlers.HealthChecker
	OpenAPIPath string
	Subs        *subscriptions.Service
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	reg, httpMetrics := metrics.NewRegistry()

	r.Use(middleware.RealIP)
	r.Use(appmw.Recover(deps.Log))
	r.Use(appmw.RequestID())
	r.Use(appmw.Logging(deps.Log))
	r.Use(appmw.Timeout(15 * time.Second))
	r.Use(appmw.Metrics(httpMetrics))

	r.Get("/healthz", handlers.Healthz(deps.Health))
	r.Get("/readyz", handlers.Readyz(deps.Health))

	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	handlers.MountSwagger(r, deps.OpenAPIPath)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/", handlers.V1(deps.Log, deps.Subs))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
