package pullrequest_merge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"service-pr-reviewer-assignment/internal/api/converters"
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/pkg/response"
	"service-pr-reviewer-assignment/internal/service/entities"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/google/uuid"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string)
	WarnfContext(ctx context.Context, format string, args ...interface{})
	ErrorfContext(ctx context.Context, format string, args ...interface{})
}

type Service interface {
	MergePullRequestAndGetReviewers(
		ctx context.Context,
		pullRequestID uuid.UUID,
	) (pr *entities.PullRequest, reviewerIDs []uuid.UUID, err error)
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

	var req dto.MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WarnfContext(ctx, "decode body failed: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "decode body failed"+err.Error())
		return
	}

	ctx = log.WithContext(ctx,
		"pull_request_id", req.PullRequestId,
	)

	prID := req.PullRequestId

	pullRequest, reviewerIDs, err := h.service.MergePullRequestAndGetReviewers(ctx, prID)
	if err != nil {
		h.logger.ErrorfContext(ctx, "merge pull request failed: %v", err)

		var prNotFound *entities.ErrPullRequestNotFound
		if errors.As(err, &prNotFound) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, prNotFound.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	h.logger.InfoContext(ctx, "pull request merged successfully")
	resp := converters.PullRequestToDTO(pullRequest, reviewerIDs)
	response.OK(w, resp)
}
