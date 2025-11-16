package team_get

import (
	"context"
	"errors"
	"net/http"

	"service-pr-reviewer-assignment/internal/api/converters"
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/pkg/response"
	"service-pr-reviewer-assignment/internal/service/entities"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string)
	ErrorfContext(ctx context.Context, format string, args ...interface{})
	LogCtx(ctx context.Context, fields ...any) context.Context
}

type Service interface {
	GetTeam(
		ctx context.Context,
		teamName string,
	) (*entities.Team, error)
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

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		h.logger.ErrorfContext(ctx, "team_name parameter is required")
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "team_name parameter is required")
		return
	}

	ctx = h.logger.LogCtx(ctx, "team_name", teamName)

	team, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		h.logger.ErrorfContext(ctx, "get team failed: %v", err)

		var notFoundErr *entities.ErrTeamNotFound
		if errors.As(err, &notFoundErr) {
			response.Error(w, http.StatusNotFound, dto.NOTFOUND, notFoundErr.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	h.logger.InfoContext(ctx, "team retrieved successfully")
	resp := converters.TeamToDTO(team)
	response.OK(w, resp)
}
