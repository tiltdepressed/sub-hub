package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func MountSwagger(r chi.Router, openAPIPath string) {
	r.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, req *http.Request) {
		path := openAPIPath
		if !filepath.IsAbs(path) {
			path = filepath.Clean(path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, "openapi spec not found\n")
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/openapi.yaml"),
	))
}
