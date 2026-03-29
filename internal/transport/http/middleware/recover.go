package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type errorResponse struct {
	Error string `json:"error"`
}

func Recover(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered", zap.Any("recover", rec))
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(errorResponse{Error: "internal"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
