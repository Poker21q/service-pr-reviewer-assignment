package storage

import (
	"github.com/Masterminds/squirrel"
)

type Storage struct {
	querier     querier
	stmtBuilder squirrel.StatementBuilderType
}

func Must(querier querier) *Storage {
	return &Storage{
		querier:     querier,
		stmtBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
