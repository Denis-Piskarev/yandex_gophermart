package internal

import (
	"github.com/DenisquaP/yandex_gophermart/internal/service"
)

// ServiceInterface - интерфейс слоя service
type ServiceInterface interface {
	service.BalanceKeeper
	service.AuthInterface
	service.OrderInterface
}
