package users_getreview

import (
	"context"
	"errors"
	"net/http"

	"service-pr-reviewer-assignment/internal/api/converters"
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/pkg/response"
	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/google/uuid"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string)
	ErrorfContext(ctx context.Context, format string, args ...interface{})
	LogCtx(ctx context.Context, fields ...any) context.Context
}

type Service interface {
	GetUserPullRequestReviewRequests(
		ctx context.Context,
		userID uuid.UUID,
	) ([]entities.PullRequest, error)
}

type Handler struct {
	logger  Logger
	service Service
}

func NewHandler(logger Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.ErrorfContext(ctx, "user_id parameter is required")
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "user_id parameter is required")
		return
	}

	ctx = h.logger.LogCtx(ctx, "user_id", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.ErrorfContext(ctx, "failed to parse user_id as uuid: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "invalid user_id format")
		return
	}

	pullRequests, err := h.service.GetUserPullRequestReviewRequests(ctx, userID)
	if err != nil {
		h.logger.ErrorfContext(ctx, "get user review requests failed: %v", err)

		var notFoundErr *entities.ErrUserNotFound
		if errors.As(err, &notFoundErr) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, notFoundErr.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	h.logger.InfoContext(ctx, "user review requests retrieved successfully")
	resp := converters.UserReviewRequestsToDTO(userID, pullRequests)
	response.OK(w, resp)
}
