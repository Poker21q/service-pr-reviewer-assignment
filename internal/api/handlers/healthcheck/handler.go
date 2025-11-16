package healthcheck

import (
	"context"
	"net/http"
	"sync/atomic"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string)
}

type Handler struct {
	isShuttingDown *atomic.Bool
	logger         Logger
}

func NewHandler(isShuttingDown *atomic.Bool, logger Logger) *Handler {
	return &Handler{
		isShuttingDown: isShuttingDown,
		logger:         logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.isShuttingDown.Load() {
		h.logger.InfoContext(ctx, "shutting down")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	h.logger.InfoContext(ctx, "ok")
	w.WriteHeader(http.StatusOK)
}
