package router

import (
	"github.com/DenisquaP/yandex_gophermart/internal/service"
	"go.uber.org/zap"
)

type Router struct {
	service service.Manager
	logger  *zap.SugaredLogger
}

func NewRouter(service service.Manager, logger *zap.SugaredLogger) *Router {
	return &Router{
		service: service,
		logger:  logger,
	}
}
