package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/service/auth/userAuth"
	"github.com/DenisquaP/yandex_gophermart/internal/service/balance"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"
	"github.com/DenisquaP/yandex_gophermart/internal/service/order"
)

type Service struct {
	auth    *userAuth.UserAuth
	balance *balance.Balance
	order   *order.Order
	token   *jwt.JWT
}

func NewService(store internal.DBStore) *internal.ServiceInterface {
	return &internal.ServiceInterface{
		AuthInterface:  userAuth.NewUserAuth(store),
		BalanceKeeper:  balance.NewBalance(store),
		OrderInterface: order.NewOrder(store),
		TokenInterface: jwt.NewJWT(),
	}
}
