package user

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
)

type User struct {
	db    internal.DBStore
	token internal.TokenInterface
}

func NewUserAuth(db internal.DBStore, token internal.TokenInterface) *User {
	return &User{
		db:    db,
		token: token,
	}
}
