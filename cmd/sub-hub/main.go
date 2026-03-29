package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sub-hub/internal/bootstrap"
	"sub-hub/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "config load error:", err)
		os.Exit(1)
	}

	app, err := bootstrap.New(ctx, cfg, bootstrap.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "bootstrap error:", err)
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "run error:", err)
		os.Exit(1)
	}
}
