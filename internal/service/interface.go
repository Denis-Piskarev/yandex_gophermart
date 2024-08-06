package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal/service/auth"
	"github.com/DenisquaP/yandex_gophermart/internal/service/balance"
	"github.com/DenisquaP/yandex_gophermart/internal/service/orders"
)

// Manager - интерфейс слоя service
type Manager interface {
	auth.Manager
	orders.Manager
	balance.Manager
}
