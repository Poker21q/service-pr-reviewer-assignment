package request_logging_context

import (
	"net/http"

	"service-pr-reviewer-assignment/pkg/log"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = log.WithContext(ctx,
			"path", r.URL.Path,
			"method", r.Method,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
