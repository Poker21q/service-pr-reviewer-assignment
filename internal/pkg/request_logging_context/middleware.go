package request_logging_context

import (
	"context"
	"net/http"
)

type Logger interface {
	LogCtx(ctx context.Context, fields ...any) context.Context
}

func Middleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = logger.LogCtx(ctx,
				"path", r.URL.Path,
				"method", r.Method,
			)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
