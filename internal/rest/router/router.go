package router

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type Router struct {
	service internal.Manager
	logger  *zap.SugaredLogger
}

func NewRouter(service internal.Manager, logger *zap.SugaredLogger) *Router {
	return &Router{
		service: service,
		logger:  logger,
	}
}

func NewRouterWithMiddleware(logger *zap.SugaredLogger, service internal.Manager) (http.Handler, error) {
	r := chi.NewRouter()

	return r, nil
}
