package querier

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier struct {
	pool   *pgxpool.Pool
	getter *pgxv5.CtxGetter
}

func MustNew(pool *pgxpool.Pool, getter *pgxv5.CtxGetter) *Querier {
	return &Querier{
		pool:   pool,
		getter: getter,
	}
}

func (q *Querier) get(ctx context.Context) pgxv5.Tr {
	return q.getter.DefaultTrOrDB(ctx, q.pool)
}

func (q *Querier) Exec(
	ctx context.Context,
	sql string,
	arguments ...interface{},
) (pgconn.CommandTag, error) {
	ex := q.get(ctx)
	return ex.Exec(ctx, sql, arguments...)
}

func (q *Querier) Query(
	ctx context.Context,
	sql string,
	args ...interface{},
) (pgx.Rows, error) {
	ex := q.get(ctx)
	return ex.Query(ctx, sql, args...)
}

func (q *Querier) QueryRow(
	ctx context.Context,
	sql string,
	args ...interface{},
) pgx.Row {
	ex := q.get(ctx)
	return ex.QueryRow(ctx, sql, args...)
}
