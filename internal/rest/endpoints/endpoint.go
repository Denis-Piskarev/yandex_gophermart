package endpoints

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"go.uber.org/zap"
)

type Endpoints struct {
	logger   *zap.SugaredLogger
	services internal.ServiceInterface
}

func NewEndpoints(logger *zap.SugaredLogger, services internal.ServiceInterface) *Endpoints {
	return &Endpoints{
		logger:   logger,
		services: services,
	}
}
