package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"sub-hub/internal/usecase/subscriptions"
)

type errorResponse struct {
	Error string `json:"error"`
}

func Router(log *zap.Logger, subs *subscriptions.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"pong": "ok"})
	})

	r.Mount("/subscriptions", SubscriptionsRouter(log, subs))

	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
