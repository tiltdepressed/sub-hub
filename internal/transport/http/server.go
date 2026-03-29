package http

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"sub-hub/internal/config"
)

type Server struct {
	log *zap.Logger
	srv *http.Server
}

func NewServer(cfg config.HTTPConfig, handler http.Handler, log *zap.Logger) *Server {
	return &Server{
		log: log,
		srv: &http.Server{
			Addr:         cfg.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	if s.srv == nil {
		return errors.New("http server is nil")
	}

	go func() {
		s.log.Info("http server listening", zap.String("addr", s.srv.Addr))
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server error", zap.Error(err))
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}
