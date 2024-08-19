package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/service/balance"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"
	"github.com/DenisquaP/yandex_gophermart/internal/service/order"
	"github.com/DenisquaP/yandex_gophermart/internal/service/user"
)

type Service struct {
	auth    *user.User
	balance *balance.Balance
	order   *order.Order
	token   *jwt.JWT
}

func NewService(store internal.DBStore) *internal.ServiceInterface {
	token := jwt.NewJWT()
	return &internal.ServiceInterface{
		AuthInterface:  user.NewUserAuth(store, token),
		BalanceKeeper:  balance.NewBalance(store),
		OrderInterface: order.NewOrder(store),
		TokenInterface: token,
	}
}
