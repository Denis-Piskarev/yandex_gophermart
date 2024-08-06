package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	repository repository.DBStore
	logger     *zap.SugaredLogger
}

func NewService(store repository.DBStore, logger *zap.SugaredLogger) *Service {
	return &Service{
		repository: store,
		logger:     logger,
	}
}
