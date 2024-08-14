package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
)

type Service struct {
	repository internal.DBStore
}

func NewService(store internal.DBStore) *Service {
	return &Service{
		repository: store,
	}
}
