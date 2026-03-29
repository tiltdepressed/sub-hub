package bootstrap

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Health struct {
	db   *pgxpool.Pool
	ping time.Duration
}

func NewHealth(db *pgxpool.Pool, pingTimeout time.Duration) *Health {
	return &Health{db: db, ping: pingTimeout}
}

func (h *Health) Liveness(ctx context.Context) error {
	return nil
}

func (h *Health) Readiness(ctx context.Context) error {
	if h.db == nil {
		return nil
	}
	pingCtx, cancel := context.WithTimeout(ctx, h.ping)
	defer cancel()
	return h.db.Ping(pingCtx)
}

func (h *Health) Close() error {
	if h.db != nil {
		h.db.Close()
	}
	return nil
}
