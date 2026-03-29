package handlers

import (
	"net/http"

	"go.uber.org/zap"

	v1 "sub-hub/internal/transport/http/handlers/v1"
	"sub-hub/internal/usecase/subscriptions"
)

func V1(log *zap.Logger, svc *subscriptions.Service) http.Handler {
	return v1.Router(log, svc)
}
