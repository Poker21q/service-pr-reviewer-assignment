package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnPool создает пул соединений PostgreSQL с retry и взвешенным джиттером.
func NewConnPool(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	const (
		maxRetries = 5
		baseDelay  = 500 * time.Millisecond
		maxJitter  = 200 * time.Millisecond
	)

	connStr := newDsn(cfg)
	pgxCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	pgxCfg.MaxConns = 12
	pgxCfg.MinConns = 2
	pgxCfg.MaxConnLifetime = 5 * time.Minute
	pgxCfg.MaxConnIdleTime = 30 * time.Minute
	pgxCfg.HealthCheckPeriod = time.Minute
	pgxCfg.ConnConfig.ConnectTimeout = 5 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
		if err == nil {
			pingErr := pool.Ping(ctx)
			if pingErr == nil {
				return pool, nil
			}

			err = pingErr
		}

		log.Warn(
			fmt.Sprintf("attempt %d: failed to connect to database: %s", attempt+1, err.Error()),
		)

		// Задержка растет линейно с номером попытки (baseDelay * (attempt+1))
		// плюс случайная вариативность jitter до maxJitter, чтобы избежать "эффекта лавины".
		//
		// | Попытка | baseDelay*(attempt+1) | jitter | sleep = baseDelay*(attempt+1)+jitter |
		// |---------|-----------------------|--------|-------------------------------------|
		// | 1       | 500ms                 | 100ms  | 600ms                               |
		// | 2       | 1000ms                | 50ms   | 1050ms                              |
		// | 3       | 1500ms                | 150ms  | 1650ms                              |
		// | 4       | 2000ms                | 20ms   | 2020ms                              |
		// | 5       | 2500ms                | 180ms  | 2680ms                              |
		jitter := time.Duration(rand.Int63n(int64(maxJitter)))
		time.Sleep(baseDelay*time.Duration(attempt+1) + jitter)
	}

	return nil, fmt.Errorf("all retries failed")
}

func newDsn(cfg config.Postgres) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}
