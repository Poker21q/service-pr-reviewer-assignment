package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries = 5
	baseDelay  = 500 * time.Millisecond
	maxJitter  = 200 * time.Millisecond
)

func NewConnPool(ctx context.Context, cfg config.Postgres, logger *log.Logger) (*pgxpool.Pool, error) {
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
	pgxCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	for attempt := 0; attempt < maxRetries; attempt++ {
		pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				return pool, nil
			}
		}

		logger.Errorf("attempt %d: failed to connect to database: %v", attempt+1, err)

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
