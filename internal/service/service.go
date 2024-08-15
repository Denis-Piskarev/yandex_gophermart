package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/service/auth"
	"github.com/DenisquaP/yandex_gophermart/internal/service/balance"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"
	"github.com/DenisquaP/yandex_gophermart/internal/service/order"
)

type Service struct {
	auth    *auth.UserAuth
	balance *balance.Balance
	order   *order.Order
	token   *jwt.JWT
}

func NewService(store internal.DBStore) *internal.ServiceInterface {
	token := jwt.NewJWT()
	return &internal.ServiceInterface{
		AuthInterface:  auth.NewUserAuth(store, token),
		BalanceKeeper:  balance.NewBalance(store),
		OrderInterface: order.NewOrder(store),
		TokenInterface: token,
	}
}
