package router

import (
	"context"
	"net/http"
	"sync/atomic"

	"service-pr-reviewer-assignment/internal/pkg/graceful_shutdown"

	"service-pr-reviewer-assignment/internal/api/handlers/healthcheck"
	"service-pr-reviewer-assignment/internal/api/handlers/pullrequest_create"
	"service-pr-reviewer-assignment/internal/api/handlers/pullrequest_merge"
	"service-pr-reviewer-assignment/internal/api/handlers/team_add"
	"service-pr-reviewer-assignment/internal/api/handlers/team_get"
	"service-pr-reviewer-assignment/internal/api/handlers/users_getreview"
	"service-pr-reviewer-assignment/internal/api/handlers/users_setisactive"
	"service-pr-reviewer-assignment/internal/pkg/panic_recover"
	"service-pr-reviewer-assignment/internal/pkg/request_logging_context"
	"service-pr-reviewer-assignment/internal/service"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/gorilla/mux"
)

func Must(
	isShuttingDown *atomic.Bool,
	ongoingCtx context.Context,
	logger *log.Logger,
	service *service.Service,
) http.Handler {
	router := mux.NewRouter()

	router.Use(panic_recover.Middleware(logger))
	router.Use(request_logging_context.Middleware)
	router.Use(graceful_shutdown.Middleware(isShuttingDown, ongoingCtx))

	router.Handle("/healthcheck", healthcheck.NewHandler(isShuttingDown, logger)).Methods(http.MethodHead)

	router.Handle("/team/add", team_add.NewHandler(logger, service)).Methods(http.MethodPost)
	router.Handle("/team/get", team_get.NewHandler(logger, service)).Methods(http.MethodGet)

	router.Handle("/users/getReview", users_getreview.NewHandler(logger, service)).Methods(http.MethodGet)
	router.Handle("/users/setIsActive", users_setisactive.NewHandler(logger, service)).Methods(http.MethodPost)

	router.Handle("/pullRequest/create", pullrequest_create.NewHandler(logger, service)).Methods(http.MethodPost)
	router.Handle("/pullRequest/merge", pullrequest_merge.NewHandler(logger, service)).Methods(http.MethodPost)

	router.NotFoundHandler = http.NotFoundHandler()

	return router
}
