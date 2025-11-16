package pullrequest_create

import (
	"context"
	"encoding/json"
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
	CreatePullRequestAndAssignReviewers(
		ctx context.Context,
		pullRequestID uuid.UUID,
		pullRequestName string,
		authorID uuid.UUID,
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

	var req dto.CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorfContext(ctx, "decode body failed: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "decode body failed"+err.Error())
		return
	}

	ctx = h.logger.LogCtx(ctx,
		"pull_request_id", req.PullRequestId,
		"pull_request_name", req.PullRequestName,
		"author_id", req.AuthorId,
	)

	prID := req.PullRequestId
	authorID := req.AuthorId

	pullRequest, reviewerIDs, err := h.service.CreatePullRequestAndAssignReviewers(ctx, prID, req.PullRequestName, authorID)
	if err != nil {
		h.logger.ErrorfContext(ctx, "create pull request failed: %v", err)

		var prExists *entities.ErrPullRequestAlreadyExists
		if errors.As(err, &prExists) {
			response.Error(w, http.StatusConflict, dto.PREXISTS, prExists.Error())
			return
		}

		var userNotFound *entities.ErrUserNotFound
		if errors.As(err, &userNotFound) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, userNotFound.Error())
			return
		}

		var pullRequestNameValidation *entities.ErrPullRequestNameValidation
		if errors.As(err, &pullRequestNameValidation) {
			response.Error(w, http.StatusNotFound, dto.BADREQUEST, pullRequestNameValidation.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	h.logger.InfoContext(ctx, "pull request created successfully")
	resp := converters.PullRequestToDTO(pullRequest, reviewerIDs)
	response.OK(w, resp)
}
