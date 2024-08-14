package userAuth

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
)

type UserAuth struct {
	db internal.DBStore
}

func NewUserAuth(db internal.DBStore) *UserAuth {
	return &UserAuth{
		db: db,
	}
}
