package auth

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
)

type UserAuth struct {
	db    internal.DBStore
	token internal.TokenInterface
}

func NewUserAuth(db internal.DBStore, token internal.TokenInterface) *UserAuth {
	return &UserAuth{
		db:    db,
		token: token,
	}
}
