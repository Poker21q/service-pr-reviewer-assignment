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

	"service-pr-reviewer-assignment/internal/app/router"

	"service-pr-reviewer-assignment/pkg/querier"
	"service-pr-reviewer-assignment/pkg/tx"

	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/internal/app/postgres"
	"service-pr-reviewer-assignment/internal/service"
	"service-pr-reviewer-assignment/internal/storage"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/avito-tech/go-transaction-manager/pgxv5"
)

const (
	shutdownPeriod      = 15 * time.Second
	shutdownHardPeriod  = 3 * time.Second
	readinessDrainDelay = 5 * time.Second
)

func Run(ctx context.Context, cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := log.Must()

	pg, err := postgres.NewConnPool(ctx, cfg.Postgres, logger)
	if err != nil {
		return fmt.Errorf("postgres new conn pool: %w", err)
	}
	defer pg.Close()

	txManager := tx.Must(pg)
	querier := querier.Must(pg, pgxv5.DefaultCtxGetter)
	storage := storage.Must(querier)
	service := service.Must(storage, txManager)

	var isShuttingDown atomic.Bool
	serverCtx, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.Must(&isShuttingDown, serverCtx, logger, service),
		BaseContext: func(net.Listener) context.Context {
			return serverCtx
		},
	}

	serverErrors := make(chan error, 1)
	go runServer(ctx, logger, server, cfg.Server.Port, serverErrors)

	if err := waitForShutdown(ctx, serverErrors, &isShuttingDown); err != nil {
		return fmt.Errorf("wait for shutdown: %w", err)
	}

	if err := shutdownServer(ctx, logger, server); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}

func runServer(ctx context.Context, logger *log.Logger, server *http.Server, port string, errorsCh chan<- error) {
	logger.InfofContext(ctx, "server listening on :%s", port)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errorsCh <- fmt.Errorf("server listen and serve: %w", err)
		return
	}

	close(errorsCh)
}

func waitForShutdown(ctx context.Context, serverErrors <-chan error, isShuttingDown *atomic.Bool) error {
	select {
	case <-ctx.Done():
		return nil
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	}
}

func shutdownServer(ctx context.Context, logger *log.Logger, server *http.Server) error {
	logger.InfoContext(ctx, "shutdown signal received")

	time.Sleep(readinessDrainDelay)
	logger.InfoContext(ctx, "draining ongoing requests...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.InfofContext(ctx, "graceful shutdown failed: %v", err)
		time.Sleep(shutdownHardPeriod)
		return fmt.Errorf("server shutdown: %w", err)
	}

	logger.InfoContext(ctx, "server stopped gracefully")
	return nil
}
