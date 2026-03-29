package middleware

import (
	"net/http"

	"sub-hub/internal/observability/tracing"
)

func Tracing(operation string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return tracing.WrapHTTPHandler(next, operation)
	}
}
