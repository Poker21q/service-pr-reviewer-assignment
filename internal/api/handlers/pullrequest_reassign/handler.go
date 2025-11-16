package pullrequest_reassign

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
	ReassignReviewer(
		ctx context.Context,
		pullRequestID uuid.UUID,
		oldReviewerID uuid.UUID,
	) (pr *entities.PullRequest, prReviewerIDs []uuid.UUID, newReviewerID uuid.UUID, err error)
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

	var req dto.ReassignPullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WarnfContext(ctx, "decode body failed: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "decode body failed"+err.Error())
		return
	}

	ctx = log.WithContext(ctx,
		"pull_request_id", req.PullRequestId,
		"old_user_id", req.OldUserId,
	)

	prID := req.PullRequestId
	oldReviewerID := req.OldUserId

	pullRequest, prReviewerIDs, newReviewerID, err := h.service.ReassignReviewer(ctx, prID, oldReviewerID)
	if err != nil {
		h.logger.ErrorfContext(ctx, "reassign reviewer failed: %v", err)

		var prNotFound *entities.ErrPullRequestNotFound
		if errors.As(err, &prNotFound) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, prNotFound.Error())
			return
		}

		var userNotFound *entities.ErrUserNotFound
		if errors.As(err, &userNotFound) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, userNotFound.Error())
			return
		}

		var prMerged *entities.ErrPullRequestAlreadyMerged
		if errors.As(err, &prMerged) {
			response.Error(w, http.StatusConflict, dto.PRMERGED, prMerged.Error())
			return
		}

		var notAssigned *entities.ErrReviewerNotAssigned
		if errors.As(err, &notAssigned) {
			response.Error(w, http.StatusConflict, dto.NOTASSIGNED, notAssigned.Error())
			return
		}

		noReplacement := entities.ErrNoReplacementCandidate
		if errors.Is(err, noReplacement) {
			response.Error(w, http.StatusConflict, dto.NOCANDIDATE, noReplacement.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	ctx = log.WithContext(ctx,
		"new_reviewer_id", newReviewerID.String(),
	)

	h.logger.InfoContext(ctx, "reviewer reassigned successfully")
	resp := converters.ReassignResponseToDTO(pullRequest, prReviewerIDs, newReviewerID)
	response.OK(w, resp)
}
