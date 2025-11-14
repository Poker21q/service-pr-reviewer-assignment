package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"service-pr-reviewer-assignment/internal/pkg/querier"

	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/internal/app/postgres"
	"service-pr-reviewer-assignment/internal/app/router"
	"service-pr-reviewer-assignment/internal/pkg/tx"
	"service-pr-reviewer-assignment/internal/service"
	"service-pr-reviewer-assignment/internal/storage"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/avito-tech/go-transaction-manager/pgxv5"
)

func Run(cfg *config.Config) error {
	const (
		shutdownPeriod      = 15 * time.Second
		shutdownHardPeriod  = 3 * time.Second
		readinessDrainDelay = 5 * time.Second
	)

	var isShuttingDown atomic.Bool

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := log.MustNewLogger()

	pg, err := postgres.NewConnPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("cannot create postgres pool: %w", err)
	}
	defer pg.Close()

	txManager := tx.MustNewManager(pg)
	q := querier.MustNew(pg, pgxv5.DefaultCtxGetter)
	st := storage.MustNew(q)
	svc := service.MustNew(st, txManager)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router.MustNew(logger, svc),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.InfoContext(ctx, "server listening on :"+cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
			return
		}

		close(serverErrors)
	}()

	select {
	case <-ctx.Done():
		logger.InfoContext(ctx, "shutdown signal received")

	case err := <-serverErrors:
		if err != nil {
			return fmt.Errorf("server failed: %w", err)
		}
		logger.InfoContext(ctx, "server stopped normally")
		return nil
	}

	isShuttingDown.Store(true)
	logger.InfoContext(ctx, "draining ongoing requests...")
	time.Sleep(readinessDrainDelay)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.InfoContext(ctx, fmt.Sprintf("graceful shutdown failed: %v", err))
		time.Sleep(shutdownHardPeriod)
	}

	logger.InfoContext(ctx, "server stopped gracefully")
	return nil
}
