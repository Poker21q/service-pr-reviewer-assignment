package team_add_post

import (
	"context"
	"encoding/json"
	"net/http"

	"service-pr-reviewer-assignment/internal/handlers/converters"
	"service-pr-reviewer-assignment/pkg/log"

	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/service/entities"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string)
	WarnContext(ctx context.Context, msg string)
	ErrorContext(ctx context.Context, msg string)
}

type Service interface {
	CreateTeamByUsers(
		ctx context.Context,
		teamName string,
		users []entities.UserModify,
	) ([]entities.User, error)
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

	h.logger.InfoContext(ctx, "team_add_post begin")

	var req dto.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WarnContext(ctx, "decode body failed: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx = log.WithContext(ctx,
		log.NewField("team_name", req.TeamName),
		log.NewField("members_count", len(req.Members)),
	)

	teamName, users, err := converters.TeamFromDTO(req)
	if err != nil {
		h.logger.WarnContext(ctx, "convert dto to entities failed: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdUsers, err := h.service.CreateTeamByUsers(ctx, teamName, users)
	if err != nil {
		h.logger.ErrorContext(ctx, "service create team failed: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.InfoContext(ctx, "team created successfully")

	resp := converters.TeamToDTO(teamName, createdUsers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
