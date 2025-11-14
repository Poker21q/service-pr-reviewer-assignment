package router

import (
	"net/http"

	"service-pr-reviewer-assignment/internal/handlers/team_add_post"
	"service-pr-reviewer-assignment/internal/service"
	"service-pr-reviewer-assignment/pkg/log"
)

func MustNew(
	logger *log.Logger,
	service *service.Service,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/team/add", team_add_post.NewHandler(logger, service))

	// default 404
	mux.Handle("/", http.NotFoundHandler())

	return mux
}
