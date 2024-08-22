package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/service/balance"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"
	"github.com/DenisquaP/yandex_gophermart/internal/service/order"
	"github.com/DenisquaP/yandex_gophermart/internal/service/user"
)

func NewService(store internal.DBStore, accrualSystemAddress string) *internal.Service {
	token := jwt.NewJWT()
	return &internal.Service{
		AuthInterface:  user.NewUserAuth(store, token),
		BalanceKeeper:  balance.NewBalance(store),
		OrderInterface: order.NewOrder(store, accrualSystemAddress),
		TokenInterface: token,
	}
}
