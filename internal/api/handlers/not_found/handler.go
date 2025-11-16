package not_found

import (
	"net/http"

	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/pkg/response"
)

type Logger interface {
	Error(msg string)
}

type Handler struct {
	logger Logger
}

func NewHandler(logger Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Error("resource not found")
	response.Error(w, http.StatusNotFound, dto.NOTFOUND, "resource not found")
}
