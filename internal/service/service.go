package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"go.uber.org/zap"
)

type Service struct {
	repository internal.DBStore
	logger     *zap.SugaredLogger
}

func NewService(store internal.DBStore, logger *zap.SugaredLogger) *Service {
	return &Service{
		repository: store,
		logger:     logger,
	}
}
