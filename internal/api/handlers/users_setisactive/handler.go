package users_setisactive

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
	SetUserActiveStatus(
		ctx context.Context,
		userID uuid.UUID,
		isActive bool,
	) (user *entities.User, err error)
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

	var req dto.SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorfContext(ctx, "decode body failed: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "decode body failed")
		return
	}

	ctx = h.logger.LogCtx(ctx,
		"user_id", req.UserId,
		"is_active", req.IsActive,
	)

	userID := req.UserId

	user, err := h.service.SetUserActiveStatus(ctx, userID, req.IsActive)
	if err != nil {
		h.logger.ErrorfContext(ctx, "set user active status failed: %v", err)

		var userNotFound *entities.ErrUserNotFound
		if errors.As(err, &userNotFound) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, userNotFound.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	h.logger.InfoContext(ctx, "user active status updated successfully")
	resp := converters.UserToDTO(user)
	response.OK(w, resp)
}
