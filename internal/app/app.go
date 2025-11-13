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

	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/internal/app/postgres"
	"service-pr-reviewer-assignment/internal/service"
	"service-pr-reviewer-assignment/internal/storage"
	"service-pr-reviewer-assignment/pkg/log"
)

// Run — точка запуска приложения.
// Реализован graceful shutdown по рекомендациям разработчика victoriametrics
// https://victoriametrics.com/blog/go-graceful-shutdown/
func Run(cfg *config.Config) error {
	const (
		shutdownPeriod      = 15 * time.Second
		shutdownHardPeriod  = 3 * time.Second
		readinessDrainDelay = 5 * time.Second
	)

	var isShuttingDown atomic.Bool

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := log.NewLogger()

	pg, err := postgres.NewConnPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("cannot create postgres pool: %w", err)
	}
	defer pg.Close()

	st := storage.New(pg)
	svc := service.New(st)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: initRouter(svc, ctx),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	serverErrCh := make(chan error, 1)

	go func() {
		logger.Info(
			fmt.Sprintf("server listening on :%s", cfg.Server.Port),
		)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- fmt.Errorf("server error: %w", err)
			return
		}

		close(serverErrCh)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

	case err = <-serverErrCh:
		if err != nil {
			return fmt.Errorf("server failed: %w", err)
		}
		logger.Info("server stopped normally")
		return nil
	}

	isShuttingDown.Store(true)
	time.Sleep(readinessDrainDelay)
	logger.Info("draining ongoing requests...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Info(
			fmt.Sprintf("graceful shutdown failed: %v", err),
		)
		time.Sleep(shutdownHardPeriod)
	}

	logger.Info("server stopped gracefully")
	return nil
}

func initRouter(a, b interface{}) http.Handler {
	return nil
}
