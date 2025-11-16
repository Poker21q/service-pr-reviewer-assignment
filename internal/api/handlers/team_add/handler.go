package team_add

import (
	"context"
	"encoding/json"
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
	CreateTeam(
		ctx context.Context,
		teamName string,
		users []entities.User,
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

	var req dto.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorfContext(ctx, "decode body failed: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "decode body failed"+err.Error())
		return
	}

	ctx = h.logger.LogCtx(ctx,
		"team_name", req.TeamName,
		"members", req.Members,
	)

	teamName, users, err := converters.TeamFromDTO(req)
	if err != nil {
		h.logger.ErrorfContext(ctx, "convert dto to entities failed: %v", err)
		response.Error(w, http.StatusBadRequest, dto.BADREQUEST, "invalid team data")
		return
	}

	team, err := h.service.CreateTeam(ctx, teamName, users)
	if err != nil {
		h.logger.ErrorfContext(ctx, "create team failed: %v", err)

		var dupIDs *entities.ErrDuplicateUserIDs
		if errors.As(err, &dupIDs) {
			response.Error(w, http.StatusBadRequest, dto.DUPLICATEUSERID, dupIDs.Error())
			return
		}

		var teamExists *entities.ErrTeamAlreadyExists
		if errors.As(err, &teamExists) {
			response.Error(w, http.StatusConflict, dto.TEAMEXISTS, teamExists.Error())
			return
		}

		var userNameValidation *entities.ErrUserNameValidation
		if errors.As(err, &userNameValidation) {
			response.Error(w, http.StatusNotFound, dto.BADREQUEST, userNameValidation.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, dto.INTERNALERROR, "internal server error")
		return
	}

	h.logger.InfoContext(ctx, "team created successfully")
	resp := converters.TeamToDTO(team)
	response.OK(w, resp)
}
