package userAuth

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"go.uber.org/zap"
)

type UserAuth struct {
	db     internal.DBStore
	logger *zap.SugaredLogger
}

func NewUserAuth(db internal.DBStore, logger *zap.SugaredLogger) *UserAuth {
	return &UserAuth{
		db:     db,
		logger: logger,
	}
}
