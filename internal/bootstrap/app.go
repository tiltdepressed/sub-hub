package bootstrap

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"sub-hub/internal/config"
	"sub-hub/internal/observability/logging"
	"sub-hub/internal/observability/tracing"
	"sub-hub/internal/repository/pg"
	transporthttp "sub-hub/internal/transport/http"
	"sub-hub/internal/usecase/subscriptions"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

type App struct {
	log             *zap.Logger
	httpServer      *transporthttp.Server
	readiness       *Health
	shutdownTracing func(context.Context) error
}

func New(ctx context.Context, cfg config.Config, build BuildInfo) (*App, error) {
	log, err := logging.New(cfg.Log, cfg.Env, build.Version, build.Commit)
	if err != nil {
		return nil, err
	}

	shutdownTracing, err := tracing.Init(ctx, cfg.Tracing, log)
	if err != nil {
		return nil, err
	}

	dbPool, err := pg.New(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	subsRepo := pg.NewSubscriptions(dbPool)
	subsSvc := subscriptions.New(subsRepo)

	health := NewHealth(dbPool, cfg.DB.HealthcheckPing)

	router := transporthttp.NewRouter(transporthttp.RouterDeps{
		Log:         log,
		Health:      health,
		OpenAPIPath: cfg.OpenAPI.Path,
		Subs:        subsSvc,
	})

	httpServer := transporthttp.NewServer(cfg.HTTP, router, log)

	return &App{
		log:             log,
		httpServer:      httpServer,
		readiness:       health,
		shutdownTracing: shutdownTracing,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.httpServer.Start(); err != nil {
		return err
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shutdownErr error
	shutdownErr = errorsJoin(shutdownErr, a.httpServer.Shutdown(shutdownCtx))
	shutdownErr = errorsJoin(shutdownErr, a.shutdownTracing(shutdownCtx))
	shutdownErr = errorsJoin(shutdownErr, a.readiness.Close())

	if shutdownErr != nil {
		a.log.Error("shutdown finished with error", zap.Error(shutdownErr))
		return shutdownErr
	}

	a.log.Info("shutdown complete")
	return nil
}

func errorsJoin(existing, next error) error {
	if existing == nil {
		return next
	}
	if next == nil {
		return existing
	}
	return fmt.Errorf("%w; %v", existing, next)
}
