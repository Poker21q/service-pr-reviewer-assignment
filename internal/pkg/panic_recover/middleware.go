package panic_recover

import (
	"context"
	"net/http"
	"runtime/debug"

	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/pkg/response"
)

type Logger interface {
	ErrorfContext(ctx context.Context, format string, args ...interface{})
}

func Middleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					ctx := r.Context()
					logger.ErrorfContext(
						ctx,
						"panic recovered: %v\nstack trace: %s",
						err,
						string(debug.Stack()),
					)

					response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
